package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illmadecoder/labctl/internal/k8s"
	"github.com/illmadecoder/labctl/internal/tui"
	"github.com/illmadecoder/labctl/internal/tutorial"
	"github.com/spf13/cobra"
)

var tutorialCmd = &cobra.Command{
	Use:   "tutorial <experiment-name>",
	Short: "Run interactive tutorial for an experiment",
	Long: `Connects to the hub cluster, reads the Experiment CR status, extracts
kubeconfigs and service endpoints, loads the tutorial.yaml, and launches
an interactive terminal UI.

The experiment must be in Running phase. Use 'task hub:up -- <name>' first.`,
	Args: cobra.ExactArgs(1),
	RunE: runTutorial,
}

func runTutorial(cmd *cobra.Command, args []string) error {
	name := args[0]
	ctx := cmd.Context()

	// Find experiments directory and tutorial file
	expDir := findExperimentsDir()
	tutorialPath := filepath.Join(expDir, name, "tutorial.yaml")
	if _, err := os.Stat(tutorialPath); err != nil {
		return fmt.Errorf("no tutorial.yaml found at %s", tutorialPath)
	}

	// Load and parse tutorial
	tut, err := tutorial.Load(tutorialPath)
	if err != nil {
		return fmt.Errorf("cannot load tutorial: %w", err)
	}

	// Try to connect to hub cluster and get experiment status
	var expInfo *k8s.ExperimentInfo
	client, clientErr := k8s.NewClient()
	if clientErr == nil {
		expInfo, err = client.GetExperiment(ctx, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot read experiment CR: %v\n", err)
			fmt.Fprintf(os.Stderr, "Running in offline mode (no live services).\n\n")
		}
	} else {
		fmt.Fprintf(os.Stderr, "Warning: cannot connect to hub cluster: %v\n", clientErr)
		fmt.Fprintf(os.Stderr, "Running in offline mode.\n\n")
	}

	// Extract kubeconfig if available
	var kubeconfigPath string
	if expInfo != nil && len(expInfo.KubeconfigSecrets) > 0 {
		homeDir, _ := os.UserHomeDir()
		labDir := filepath.Join(homeDir, ".illmlab")
		os.MkdirAll(labDir, 0700)

		for targetName, secretName := range expInfo.KubeconfigSecrets {
			data, err := client.GetSecretData(ctx, expInfo.Namespace, secretName, "kubeconfig")
			if err != nil {
				continue
			}
			outPath := filepath.Join(labDir, fmt.Sprintf("kubeconfig-%s-%s", name, targetName))
			if err := os.WriteFile(outPath, data, 0600); err != nil {
				continue
			}
			if kubeconfigPath == "" {
				kubeconfigPath = outPath
			}
		}
	}

	// Build TUI model
	model := tui.NewModel(tut, expInfo, kubeconfigPath)

	// Launch bubbletea
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
