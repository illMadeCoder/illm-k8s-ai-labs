package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/illmadecoder/labctl/internal/k8s"
	"github.com/spf13/cobra"
)

var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig <experiment-name> [target]",
	Short: "Extract kubeconfig from experiment to local file",
	Long: `Reads the kubeconfig secret for an experiment's target cluster and writes
it to ~/.illmlab/kubeconfig-<experiment>-<target>.

If no target is specified, extracts kubeconfigs for all targets.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runKubeconfig,
}

func runKubeconfig(cmd *cobra.Command, args []string) error {
	name := args[0]
	targetFilter := ""
	if len(args) > 1 {
		targetFilter = args[1]
	}

	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("cannot connect to hub cluster: %w", err)
	}

	exp, err := client.GetExperiment(cmd.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot get experiment %q: %w", name, err)
	}

	if len(exp.KubeconfigSecrets) == 0 {
		return fmt.Errorf("experiment %q has no kubeconfig secrets (is spec.tutorial.exposeKubeconfig enabled?)", name)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}
	labDir := filepath.Join(homeDir, ".illmlab")
	if err := os.MkdirAll(labDir, 0700); err != nil {
		return fmt.Errorf("cannot create lab directory: %w", err)
	}

	for targetName, secretName := range exp.KubeconfigSecrets {
		if targetFilter != "" && targetName != targetFilter {
			continue
		}

		data, err := client.GetSecretData(cmd.Context(), exp.Namespace, secretName, "kubeconfig")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot read kubeconfig for target %q: %v\n", targetName, err)
			continue
		}

		outPath := filepath.Join(labDir, fmt.Sprintf("kubeconfig-%s-%s", name, targetName))
		if err := os.WriteFile(outPath, data, 0600); err != nil {
			return fmt.Errorf("cannot write kubeconfig: %w", err)
		}

		fmt.Printf("Wrote kubeconfig for target %q to %s\n", targetName, outPath)
		fmt.Printf("  export KUBECONFIG=%s\n", outPath)
	}

	return nil
}
