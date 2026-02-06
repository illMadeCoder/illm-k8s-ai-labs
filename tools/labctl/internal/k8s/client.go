package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var experimentGVR = schema.GroupVersionResource{
	Group:    "experiments.illm.io",
	Version:  "v1alpha1",
	Resource: "experiments",
}

const defaultNamespace = "experiments"

// Client provides access to the hub cluster.
type Client struct {
	clientset *kubernetes.Clientset
	dynamic   dynamic.Interface
}

// NewClient creates a Client using the current kubeconfig context.
func NewClient() (*Client, error) {
	config, err := buildConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot build kubeconfig: %w", err)
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create dynamic client: %w", err)
	}

	return &Client{clientset: cs, dynamic: dyn}, nil
}

func buildConfig() (*rest.Config, error) {
	// Try in-cluster first
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// ExperimentInfo holds extracted experiment information for the TUI.
type ExperimentInfo struct {
	Name              string
	Namespace         string
	Phase             string
	TTLDays           int
	CompletionMode    string
	Targets           []TargetInfo
	Services          []ServiceInfo
	KubeconfigSecrets map[string]string
}

// TargetInfo holds target status.
type TargetInfo struct {
	Name        string
	ClusterName string
	Phase       string
	Endpoint    string
}

// ServiceInfo holds discovered service info.
type ServiceInfo struct {
	Name     string
	Endpoint string
	Ready    bool
}

// GetExperiment reads an Experiment CR and extracts relevant info.
func (c *Client) GetExperiment(ctx context.Context, name string) (*ExperimentInfo, error) {
	obj, err := c.dynamic.Resource(experimentGVR).Namespace(defaultNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	info := &ExperimentInfo{
		Name:      name,
		Namespace: defaultNamespace,
	}

	// Phase
	info.Phase, _, _ = unstructured.NestedString(obj.Object, "status", "phase")

	// TTL
	ttl, _, _ := unstructured.NestedInt64(obj.Object, "spec", "ttlDays")
	info.TTLDays = int(ttl)

	// Completion mode
	info.CompletionMode, _, _ = unstructured.NestedString(obj.Object, "spec", "workflow", "completion", "mode")

	// Target statuses
	targets, _, _ := unstructured.NestedSlice(obj.Object, "status", "targets")
	for _, t := range targets {
		tm, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		ti := TargetInfo{}
		ti.Name, _, _ = unstructured.NestedString(tm, "name")
		ti.ClusterName, _, _ = unstructured.NestedString(tm, "clusterName")
		ti.Phase, _, _ = unstructured.NestedString(tm, "phase")
		ti.Endpoint, _, _ = unstructured.NestedString(tm, "endpoint")
		info.Targets = append(info.Targets, ti)
	}

	// Tutorial status - services
	services, _, _ := unstructured.NestedSlice(obj.Object, "status", "tutorialStatus", "services")
	for _, s := range services {
		sm, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		si := ServiceInfo{}
		si.Name, _, _ = unstructured.NestedString(sm, "name")
		si.Endpoint, _, _ = unstructured.NestedString(sm, "endpoint")
		si.Ready, _, _ = unstructured.NestedBool(sm, "ready")
		info.Services = append(info.Services, si)
	}

	// Tutorial status - kubeconfig secrets
	kcSecrets, _, _ := unstructured.NestedStringMap(obj.Object, "status", "tutorialStatus", "kubeconfigSecrets")
	if len(kcSecrets) > 0 {
		info.KubeconfigSecrets = kcSecrets
	}

	return info, nil
}

// ListExperimentPhases returns a map of experiment name -> phase.
func (c *Client) ListExperimentPhases(ctx context.Context) (map[string]string, error) {
	list, err := c.dynamic.Resource(experimentGVR).Namespace(defaultNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, item := range list.Items {
		phase, _, _ := unstructured.NestedString(item.Object, "status", "phase")
		result[item.GetName()] = phase
	}
	return result, nil
}

// GetSecretData reads a specific key from a Secret.
func (c *Client) GetSecretData(ctx context.Context, namespace, name, key string) ([]byte, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, ok := secret.Data[key]
	if !ok {
		// Try common alternatives
		for _, altKey := range []string{"value", "kubeconfig", "config"} {
			if d, ok := secret.Data[altKey]; ok {
				return d, nil
			}
		}
		return nil, fmt.Errorf("key %q not found in secret %s/%s", key, namespace, name)
	}
	return data, nil
}
