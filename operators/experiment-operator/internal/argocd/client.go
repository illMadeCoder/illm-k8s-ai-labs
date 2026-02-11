package argocd

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	experimentsv1alpha1 "github.com/illmadecoder/experiment-operator/api/v1alpha1"
)

// Client provides ArgoCD integration
type Client struct {
	client.Client
	AppManager *ApplicationManager
}

// NewClient creates a new ArgoCD client
func NewClient(c client.Client, opts ...ClientOption) *Client {
	am := NewApplicationManager(c)
	for _, opt := range opts {
		opt(am)
	}
	return &Client{
		Client:     c,
		AppManager: am,
	}
}

// ClientOption configures the ArgoCD client.
type ClientOption func(*ApplicationManager)

// WithTailscaleOAuth sets Tailscale OAuth credentials for target cluster observability.
func WithTailscaleOAuth(clientID, clientSecret string) ClientOption {
	return func(am *ApplicationManager) {
		am.TailscaleClientID = clientID
		am.TailscaleClientSecret = clientSecret
	}
}

// RegisterClusterAndCreateApps registers a cluster and creates apps for all components.
// For targets with observability enabled, deploys infra+obs layers first (workload deferred to reconcileReady).
// For targets without observability, deploys a single workload application.
func (c *Client) RegisterClusterAndCreateApps(ctx context.Context, experimentName string, target experimentsv1alpha1.Target, clusterName string, kubeconfig []byte, server string) error {
	// Register cluster with ArgoCD
	if err := RegisterCluster(ctx, c.Client, clusterName, kubeconfig, server); err != nil {
		return err
	}

	// Create ArgoCD Application for this target
	if err := c.AppManager.CreateApplication(ctx, experimentName, target, server); err != nil {
		return err
	}

	return nil
}

// RegisterClusterAndCreateLayeredApps registers a cluster and creates infra+obs layer apps.
// Workload layer is deferred to reconcileReady after infra+obs are healthy.
func (c *Client) RegisterClusterAndCreateLayeredApps(ctx context.Context, experimentName string, target experimentsv1alpha1.Target, clusterName string, kubeconfig []byte, server string, classified ClassifiedComponents) ([]string, error) {
	// Register cluster with ArgoCD
	if err := RegisterCluster(ctx, c.Client, clusterName, kubeconfig, server); err != nil {
		return nil, err
	}

	var deployedLayers []string

	// Deploy infra layer (e.g., tailscale-operator)
	if len(classified.Infra) > 0 {
		if err := c.AppManager.CreateLayeredApplication(ctx, experimentName, target, server, LayerInfra, classified.Infra); err != nil {
			return deployedLayers, fmt.Errorf("create infra layer: %w", err)
		}
		deployedLayers = append(deployedLayers, LayerInfra)
	}

	// Deploy obs layer (e.g., metrics-egress, metrics-agent) simultaneously
	// Alloy's remote-write retries until Tailscale tunnel is up
	if len(classified.Obs) > 0 {
		if err := c.AppManager.CreateLayeredApplication(ctx, experimentName, target, server, LayerObs, classified.Obs); err != nil {
			return deployedLayers, fmt.Errorf("create obs layer: %w", err)
		}
		deployedLayers = append(deployedLayers, LayerObs)
	}

	return deployedLayers, nil
}

// DeleteClusterAndApps unregisters a cluster and deletes all its applications.
// Handles both layered and non-layered applications.
func (c *Client) DeleteClusterAndApps(ctx context.Context, experimentName string, targets []experimentsv1alpha1.Target, clusterNames []string) error {
	// Delete all applications (layered + legacy single-app)
	for _, target := range targets {
		// Delete layered apps (infra, obs, workload)
		c.AppManager.DeleteLayeredApplications(ctx, experimentName, target.Name)
		// Also try deleting the legacy single app name in case it exists
		if err := c.AppManager.DeleteApplication(ctx, experimentName, target.Name); err != nil {
			// Log but continue with other deletions
			continue
		}
	}

	// Unregister clusters
	for _, clusterName := range clusterNames {
		if err := UnregisterCluster(ctx, c.Client, clusterName); err != nil {
			// Log but continue
			continue
		}
	}

	return nil
}
