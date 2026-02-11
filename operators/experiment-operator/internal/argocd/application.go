package argocd

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	experimentsv1alpha1 "github.com/illmadecoder/experiment-operator/api/v1alpha1"
	"github.com/illmadecoder/experiment-operator/internal/components"
)

var (
	// ArgoCD Application GVK
	applicationGVK = schema.GroupVersionKind{
		Group:   "argoproj.io",
		Version: "v1alpha1",
		Kind:    "Application",
	}
)

// ApplicationManager manages ArgoCD Applications
type ApplicationManager struct {
	client.Client
	Resolver              *components.Resolver
	TailscaleClientID     string
	TailscaleClientSecret string
}

// NewApplicationManager creates a new ApplicationManager
func NewApplicationManager(c client.Client) *ApplicationManager {
	return &ApplicationManager{
		Client:   c,
		Resolver: components.NewResolver(c),
	}
}

// CreateApplication creates an ArgoCD Application for a target
func (m *ApplicationManager) CreateApplication(ctx context.Context, experimentName string, target experimentsv1alpha1.Target, clusterServer string) error {
	log := log.FromContext(ctx)

	appName := fmt.Sprintf("%s-%s", experimentName, target.Name)

	// Build ArgoCD Application
	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)
	app.SetName(appName)
	app.SetNamespace("argocd")

	// Set labels
	app.SetLabels(map[string]string{
		"app.kubernetes.io/managed-by": "experiment-operator",
		"experiments.illm.io/experiment": experimentName,
		"experiments.illm.io/target":     target.Name,
	})

	// Resolve components using the resolver
	resolvedComponents, err := m.Resolver.ResolveComponents(ctx, target.Components)
	if err != nil {
		return fmt.Errorf("failed to resolve components: %w", err)
	}

	// Auto-inject observability components when enabled
	if target.Observability != nil && target.Observability.Enabled {
		obsRefs := ObservabilityComponentRefs(target.Observability, experimentName, m.TailscaleClientID, m.TailscaleClientSecret)
		obsResolved, obsErr := m.Resolver.ResolveComponents(ctx, obsRefs)
		if obsErr != nil {
			log.Error(obsErr, "Failed to resolve observability components — continuing without them")
		} else {
			resolvedComponents = append(resolvedComponents, obsResolved...)
		}
	}

	// Build sources from resolved components
	sources := []interface{}{}
	for _, resolved := range resolvedComponents {
		// Check if any source in this component uses $values references.
		// If so, we need to add ref: "values" to the git source (non-chart source).
		needsValuesRef := false
		for _, source := range resolved.Sources {
			if source.Helm != nil {
				for _, vf := range source.Helm.ValuesFiles {
					if strings.HasPrefix(vf, "$values") {
						needsValuesRef = true
						break
					}
				}
			}
		}

		for _, source := range resolved.Sources {
			argoSource := map[string]interface{}{
				"repoURL":        source.RepoURL,
				"targetRevision": source.TargetRevision,
			}

			// Use 'chart' for Helm repositories, 'path' for Git repositories
			if source.Chart != "" {
				argoSource["chart"] = source.Chart
			} else if needsValuesRef && source.Helm == nil {
				// This git source is used as a ref for $values — set ref instead of path
				// to avoid ArgoCD trying to deploy component.yaml as a manifest.
				argoSource["ref"] = "values"
			} else {
				argoSource["path"] = source.Path
			}

			// Add Helm configuration if present
			if source.Helm != nil {
				helmConfig := map[string]interface{}{}

				if source.Helm.ReleaseName != "" {
					helmConfig["releaseName"] = source.Helm.ReleaseName
				}

				if len(source.Helm.ValuesFiles) > 0 {
					// Convert []string to []interface{} for unstructured deep copy compatibility
					vf := make([]interface{}, len(source.Helm.ValuesFiles))
					for i, f := range source.Helm.ValuesFiles {
						vf[i] = f
					}
					helmConfig["valueFiles"] = vf
				}

				if len(source.Helm.Parameters) > 0 {
					helmParams := []interface{}{}
					for key, value := range source.Helm.Parameters {
						helmParams = append(helmParams, map[string]interface{}{
							"name":  key,
							"value": value,
						})
					}
					helmConfig["parameters"] = helmParams
				}

				if len(helmConfig) > 0 {
					argoSource["helm"] = helmConfig
				}
			}

			sources = append(sources, argoSource)
		}
	}

	// If no sources resolved, skip creating the application
	if len(sources) == 0 {
		log.Info("No components resolved for target, skipping application creation", "target", target.Name)
		return nil
	}

	// Pre-create the destination namespace with appropriate labels
	// This ensures PodSecurity labels are set before ArgoCD syncs workloads
	if err := m.ensureNamespace(ctx, experimentName); err != nil {
		log.Error(err, "Failed to ensure namespace", "namespace", experimentName)
		// Non-fatal: ArgoCD's CreateNamespace=true will create it without labels
	}

	// Build Application spec
	spec := map[string]interface{}{
		"project": "default",
		"sources":  sources,
		"destination": map[string]interface{}{
			"server":    clusterServer,
			"namespace": experimentName,
		},
		"syncPolicy": map[string]interface{}{
			"automated": map[string]interface{}{
				"prune":    true,
				"selfHeal": true,
			},
			"syncOptions": []interface{}{
				"CreateNamespace=true",
				"ServerSideApply=true",
			},
		},
	}

	if err := unstructured.SetNestedMap(app.Object, spec, "spec"); err != nil {
		return fmt.Errorf("failed to set application spec: %w", err)
	}

	// Check if application already exists
	existing := &unstructured.Unstructured{}
	existing.SetGroupVersionKind(applicationGVK)
	err = m.Get(ctx, client.ObjectKey{Name: appName, Namespace: "argocd"}, existing)
	if err == nil {
		// Application exists, update it
		existing.Object["spec"] = app.Object["spec"]
		if err := m.Update(ctx, existing); err != nil {
			return fmt.Errorf("failed to update application: %w", err)
		}
		log.Info("Updated ArgoCD Application", "name", appName)
		return nil
	}

	// Create new application
	if err := m.Create(ctx, app); err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	log.Info("Created ArgoCD Application", "name", appName, "target", target.Name)
	return nil
}

// DeleteApplication deletes an ArgoCD Application
func (m *ApplicationManager) DeleteApplication(ctx context.Context, experimentName string, targetName string) error {
	log := log.FromContext(ctx)

	appName := fmt.Sprintf("%s-%s", experimentName, targetName)

	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)
	app.SetName(appName)
	app.SetNamespace("argocd")

	if err := m.Delete(ctx, app); err != nil {
		if errors.IsNotFound(err) {
			log.Info("ArgoCD Application already deleted", "name", appName)
			return nil
		}
		return fmt.Errorf("failed to delete application: %w", err)
	}

	log.Info("Deleted ArgoCD Application", "name", appName)
	return nil
}

// IsApplicationHealthy checks if an Application is healthy
func (m *ApplicationManager) IsApplicationHealthy(ctx context.Context, experimentName string, targetName string) (bool, error) {
	appName := fmt.Sprintf("%s-%s", experimentName, targetName)

	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)

	if err := m.Get(ctx, client.ObjectKey{Name: appName, Namespace: "argocd"}, app); err != nil {
		return false, fmt.Errorf("failed to get application: %w", err)
	}

	// Check health status
	healthStatus, found, err := unstructured.NestedString(app.Object, "status", "health", "status")
	if err != nil || !found {
		return false, nil // Not ready yet
	}

	// Check sync status — multi-source apps report "Unknown" which is acceptable
	syncStatus, found, err := unstructured.NestedString(app.Object, "status", "sync", "status")
	if err != nil || !found {
		return false, nil
	}

	// Check for ComparisonError conditions — ArgoCD reports "Healthy" even when
	// manifest generation fails (since no resources exist to be unhealthy).
	conditions, condFound, _ := unstructured.NestedSlice(app.Object, "status", "conditions")
	if condFound {
		for _, c := range conditions {
			cond, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			condType, _, _ := unstructured.NestedString(cond, "type")
			if condType == "ComparisonError" {
				return false, nil
			}
		}
	}

	// Application is ready when core components are deployed. Complex stacks
	// (e.g. Mimir distributed mode) may stay "Degraded" due to resource pressure,
	// but monitoring services are typically available. Accept Healthy, Degraded,
	// or Progressing — only reject Missing (nothing deployed) or Suspended.
	acceptableHealth := healthStatus == "Healthy" || healthStatus == "Degraded" || healthStatus == "Progressing"
	// Multi-source ArgoCD apps report "Unknown" sync status, which is acceptable.
	acceptableSync := syncStatus == "Synced" || syncStatus == "Unknown" || syncStatus == "OutOfSync"
	return acceptableHealth && acceptableSync, nil
}

// GetApplicationComponents returns the list of deployed components
func (m *ApplicationManager) GetApplicationComponents(ctx context.Context, experimentName string, targetName string) ([]string, error) {
	appName := fmt.Sprintf("%s-%s", experimentName, targetName)

	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)

	if err := m.Get(ctx, client.ObjectKey{Name: appName, Namespace: "argocd"}, app); err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Get resources from status
	resources, found, err := unstructured.NestedSlice(app.Object, "status", "resources")
	if err != nil || !found {
		return []string{}, nil
	}

	components := []string{}
	for _, res := range resources {
		resource, ok := res.(map[string]interface{})
		if !ok {
			continue
		}

		kind, _, _ := unstructured.NestedString(resource, "kind")
		name, _, _ := unstructured.NestedString(resource, "name")
		if kind != "" && name != "" {
			components = append(components, fmt.Sprintf("%s/%s", kind, name))
		}
	}

	return components, nil
}

// ObservabilityComponentRefs returns ComponentRefs for the observability stack
// based on the target's ObservabilitySpec.
func ObservabilityComponentRefs(obs *experimentsv1alpha1.ObservabilitySpec, experimentName, tsClientID, tsClientSecret string) []experimentsv1alpha1.ComponentRef {
	refs := []experimentsv1alpha1.ComponentRef{
		// VictoriaMetrics egress service (always needed)
		{Config: "metrics-egress"},
		// Metrics agent with experiment name as external label
		{
			App: "metrics-agent",
			Params: map[string]string{
				"alloy.extraEnv[0].value": experimentName,
			},
		},
	}

	// Tailscale operator for mesh transport
	if obs.Transport == "tailscale" {
		ref := experimentsv1alpha1.ComponentRef{
			App: "tailscale-operator",
		}
		if tsClientID != "" && tsClientSecret != "" {
			ref.Params = map[string]string{
				"oauth.clientId":     tsClientID,
				"oauth.clientSecret": tsClientSecret,
			}
		}
		refs = append(refs, ref)
	}

	return refs
}

// LayerName constants for layered deployment ordering.
const (
	LayerInfra    = "infra"
	LayerObs      = "obs"
	LayerWorkload = "workload"
)

// ClassifiedComponents holds component refs split into deployment layers.
type ClassifiedComponents struct {
	Infra    []experimentsv1alpha1.ComponentRef
	Obs      []experimentsv1alpha1.ComponentRef
	Workload []experimentsv1alpha1.ComponentRef
}

// ClassifyComponents splits user and observability components into deployment layers.
// Infra: tailscale-operator. Obs: metrics-egress, metrics-agent. Workload: user components.
// When observability is disabled (no obsRefs), everything goes into Workload.
func ClassifyComponents(userComponents []experimentsv1alpha1.ComponentRef, obsComponents []experimentsv1alpha1.ComponentRef) ClassifiedComponents {
	result := ClassifiedComponents{
		Workload: userComponents,
	}

	for _, ref := range obsComponents {
		switch {
		case ref.App == "tailscale-operator":
			result.Infra = append(result.Infra, ref)
		case ref.App == "metrics-agent" || ref.Config == "metrics-egress":
			result.Obs = append(result.Obs, ref)
		default:
			// Unknown obs component — put in obs layer
			result.Obs = append(result.Obs, ref)
		}
	}

	return result
}

// HasLayers returns true if infra or obs layers have components (i.e., layering is needed).
func (c ClassifiedComponents) HasLayers() bool {
	return len(c.Infra) > 0 || len(c.Obs) > 0
}

// layerAppName returns the ArgoCD Application name for a given layer.
// Workload layer uses {exp}-{target} for backward compatibility.
func layerAppName(experimentName, targetName, layer string) string {
	if layer == LayerWorkload {
		return fmt.Sprintf("%s-%s", experimentName, targetName)
	}
	return fmt.Sprintf("%s-%s-%s", experimentName, targetName, layer)
}

// CreateLayeredApplication creates an ArgoCD Application for a specific deployment layer.
func (m *ApplicationManager) CreateLayeredApplication(ctx context.Context, experimentName string, target experimentsv1alpha1.Target, clusterServer string, layer string, componentRefs []experimentsv1alpha1.ComponentRef) error {
	log := log.FromContext(ctx)

	appName := layerAppName(experimentName, target.Name, layer)

	// Resolve components
	resolvedComponents, err := m.Resolver.ResolveComponents(ctx, componentRefs)
	if err != nil {
		return fmt.Errorf("failed to resolve %s layer components: %w", layer, err)
	}

	// Build sources from resolved components
	sources := []interface{}{}
	for _, resolved := range resolvedComponents {
		needsValuesRef := false
		for _, source := range resolved.Sources {
			if source.Helm != nil {
				for _, vf := range source.Helm.ValuesFiles {
					if strings.HasPrefix(vf, "$values") {
						needsValuesRef = true
						break
					}
				}
			}
		}

		for _, source := range resolved.Sources {
			argoSource := map[string]interface{}{
				"repoURL":        source.RepoURL,
				"targetRevision": source.TargetRevision,
			}

			if source.Chart != "" {
				argoSource["chart"] = source.Chart
			} else if needsValuesRef && source.Helm == nil {
				argoSource["ref"] = "values"
			} else {
				argoSource["path"] = source.Path
			}

			if source.Helm != nil {
				helmConfig := map[string]interface{}{}
				if source.Helm.ReleaseName != "" {
					helmConfig["releaseName"] = source.Helm.ReleaseName
				}
				if len(source.Helm.ValuesFiles) > 0 {
					vf := make([]interface{}, len(source.Helm.ValuesFiles))
					for i, f := range source.Helm.ValuesFiles {
						vf[i] = f
					}
					helmConfig["valueFiles"] = vf
				}
				if len(source.Helm.Parameters) > 0 {
					helmParams := []interface{}{}
					for key, value := range source.Helm.Parameters {
						helmParams = append(helmParams, map[string]interface{}{
							"name":  key,
							"value": value,
						})
					}
					helmConfig["parameters"] = helmParams
				}
				if len(helmConfig) > 0 {
					argoSource["helm"] = helmConfig
				}
			}

			sources = append(sources, argoSource)
		}
	}

	if len(sources) == 0 {
		log.Info("No components resolved for layer, skipping", "target", target.Name, "layer", layer)
		return nil
	}

	if err := m.ensureNamespace(ctx, experimentName); err != nil {
		log.Error(err, "Failed to ensure namespace", "namespace", experimentName)
	}

	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)
	app.SetName(appName)
	app.SetNamespace("argocd")
	app.SetLabels(map[string]string{
		"app.kubernetes.io/managed-by":  "experiment-operator",
		"experiments.illm.io/experiment": experimentName,
		"experiments.illm.io/target":     target.Name,
		"experiments.illm.io/layer":      layer,
	})

	spec := map[string]interface{}{
		"project": "default",
		"sources": sources,
		"destination": map[string]interface{}{
			"server":    clusterServer,
			"namespace": experimentName,
		},
		"syncPolicy": map[string]interface{}{
			"automated": map[string]interface{}{
				"prune":    true,
				"selfHeal": true,
			},
			"syncOptions": []interface{}{
				"CreateNamespace=true",
				"ServerSideApply=true",
			},
		},
	}

	if err := unstructured.SetNestedMap(app.Object, spec, "spec"); err != nil {
		return fmt.Errorf("failed to set application spec: %w", err)
	}

	existing := &unstructured.Unstructured{}
	existing.SetGroupVersionKind(applicationGVK)
	err = m.Get(ctx, client.ObjectKey{Name: appName, Namespace: "argocd"}, existing)
	if err == nil {
		existing.Object["spec"] = app.Object["spec"]
		if err := m.Update(ctx, existing); err != nil {
			return fmt.Errorf("failed to update %s layer application: %w", layer, err)
		}
		log.Info("Updated layered ArgoCD Application", "name", appName, "layer", layer)
		return nil
	}

	if err := m.Create(ctx, app); err != nil {
		return fmt.Errorf("failed to create %s layer application: %w", layer, err)
	}

	log.Info("Created layered ArgoCD Application", "name", appName, "layer", layer, "target", target.Name)
	return nil
}

// IsLayerHealthy checks if an ArgoCD Application for a specific layer is healthy.
func (m *ApplicationManager) IsLayerHealthy(ctx context.Context, experimentName, targetName, layer string) (bool, error) {
	appName := layerAppName(experimentName, targetName, layer)

	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(applicationGVK)

	if err := m.Get(ctx, client.ObjectKey{Name: appName, Namespace: "argocd"}, app); err != nil {
		if errors.IsNotFound(err) {
			return false, nil // App not created yet
		}
		return false, fmt.Errorf("failed to get %s layer application: %w", layer, err)
	}

	healthStatus, found, err := unstructured.NestedString(app.Object, "status", "health", "status")
	if err != nil || !found {
		return false, nil
	}

	syncStatus, found, err := unstructured.NestedString(app.Object, "status", "sync", "status")
	if err != nil || !found {
		return false, nil
	}

	// Check for ComparisonError conditions
	conditions, condFound, _ := unstructured.NestedSlice(app.Object, "status", "conditions")
	if condFound {
		for _, c := range conditions {
			cond, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			condType, _, _ := unstructured.NestedString(cond, "type")
			if condType == "ComparisonError" {
				return false, nil
			}
		}
	}

	acceptableHealth := healthStatus == "Healthy" || healthStatus == "Degraded" || healthStatus == "Progressing"
	acceptableSync := syncStatus == "Synced" || syncStatus == "Unknown" || syncStatus == "OutOfSync"
	return acceptableHealth && acceptableSync, nil
}

// DeleteLayeredApplications deletes all layer applications for a target (infra, obs, workload).
func (m *ApplicationManager) DeleteLayeredApplications(ctx context.Context, experimentName, targetName string) error {
	log := log.FromContext(ctx)

	for _, layer := range []string{LayerInfra, LayerObs, LayerWorkload} {
		appName := layerAppName(experimentName, targetName, layer)

		app := &unstructured.Unstructured{}
		app.SetGroupVersionKind(applicationGVK)
		app.SetName(appName)
		app.SetNamespace("argocd")

		if err := m.Delete(ctx, app); err != nil {
			if errors.IsNotFound(err) {
				continue // Already deleted or never created
			}
			log.Error(err, "Failed to delete layered application", "name", appName, "layer", layer)
		} else {
			log.Info("Deleted layered ArgoCD Application", "name", appName, "layer", layer)
		}
	}

	return nil
}

// ensureNamespace creates the experiment namespace with appropriate labels
func (m *ApplicationManager) ensureNamespace(ctx context.Context, name string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by":           "experiment-operator",
				"pod-security.kubernetes.io/enforce":     "privileged",
				"pod-security.kubernetes.io/enforce-version": "latest",
			},
		},
	}

	existing := &corev1.Namespace{}
	err := m.Get(ctx, client.ObjectKey{Name: name}, existing)
	if err == nil {
		// Namespace exists, ensure labels
		if existing.Labels == nil {
			existing.Labels = make(map[string]string)
		}
		existing.Labels["pod-security.kubernetes.io/enforce"] = "privileged"
		existing.Labels["pod-security.kubernetes.io/enforce-version"] = "latest"
		return m.Update(ctx, existing)
	}

	return m.Create(ctx, ns)
}
