package argocd

import (
	"context"

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

// RegisterClusterAndCreateApps registers a cluster and creates apps for all components
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

// DeleteClusterAndApps unregisters a cluster and deletes all its applications
func (c *Client) DeleteClusterAndApps(ctx context.Context, experimentName string, targets []experimentsv1alpha1.Target, clusterNames []string) error {
	// Delete all applications
	for _, target := range targets {
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
