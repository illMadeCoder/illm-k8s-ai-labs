package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/illmadecoder/labctl/internal/k8s"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available experiments and their tutorial status",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	// Find experiments directory
	expDir := findExperimentsDir()

	// List local experiments with tutorial.yaml
	entries, err := os.ReadDir(expDir)
	if err != nil {
		return fmt.Errorf("cannot read experiments directory: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EXPERIMENT\tTUTORIAL\tCR STATUS")

	// Try to connect to hub cluster for status
	client, clientErr := k8s.NewClient()
	var experiments map[string]string
	if clientErr == nil {
		experiments, _ = client.ListExperimentPhases(cmd.Context())
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "_template" {
			continue
		}
		name := entry.Name()

		// Check for tutorial.yaml
		hasTutorial := "no"
		if _, err := os.Stat(filepath.Join(expDir, name, "tutorial.yaml")); err == nil {
			hasTutorial = "yes"
		}

		// Check CR status
		crStatus := "-"
		if experiments != nil {
			if phase, ok := experiments[name]; ok {
				crStatus = phase
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", name, hasTutorial, crStatus)
	}

	w.Flush()
	return nil
}

func findExperimentsDir() string {
	// Walk up from CWD looking for experiments/
	dir, _ := os.Getwd()
	for {
		candidate := filepath.Join(dir, "experiments")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "experiments"
}
