package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/illmadecoder/labctl/internal/k8s"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <experiment-name>",
	Short: "Show experiment status including tutorial services",
	Args:  cobra.ExactArgs(1),
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	name := args[0]

	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("cannot connect to hub cluster: %w", err)
	}

	exp, err := client.GetExperiment(cmd.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot get experiment %q: %w", name, err)
	}

	fmt.Printf("Experiment: %s\n", exp.Name)
	fmt.Printf("Phase:      %s\n", exp.Phase)
	fmt.Printf("TTL:        %d days\n", exp.TTLDays)

	if exp.CompletionMode != "" {
		fmt.Printf("Mode:       %s\n", exp.CompletionMode)
	}

	fmt.Println()

	// Targets
	if len(exp.Targets) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TARGET\tCLUSTER\tPHASE\tENDPOINT")
		for _, t := range exp.Targets {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.Name, t.ClusterName, t.Phase, t.Endpoint)
		}
		w.Flush()
		fmt.Println()
	}

	// Tutorial services
	if len(exp.Services) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tENDPOINT\tREADY")
		for _, s := range exp.Services {
			fmt.Fprintf(w, "%s\t%s\t%v\n", s.Name, s.Endpoint, s.Ready)
		}
		w.Flush()
		fmt.Println()
	}

	// Kubeconfig secrets
	if len(exp.KubeconfigSecrets) > 0 {
		fmt.Println("Kubeconfig Secrets:")
		for target, secret := range exp.KubeconfigSecrets {
			fmt.Printf("  %s: %s\n", target, secret)
		}
	}

	return nil
}
