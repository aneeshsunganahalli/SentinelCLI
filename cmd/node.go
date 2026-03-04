package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/aneeshsunganahalli/SentintelCLI/internal"
)

var history string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Monitor HPC node performance",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Initialize the Prometheus Client
		client, err := internal.NewPromClient(prometheusURL)
		if err != nil {
			log.Fatalf("Failed to connect to Prometheus: %v", err)
		}

		// 2. Fetch the metrics using the history flag
		fmt.Printf("Fetching metrics with lookback: %s...\n", history)
		stats, err := client.GetNodeStats(history)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// 3. Initialize tabwriter to handle column alignment
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "NODE\tCPU (%)\tMEM (%)\tDISK (%)\t")
		fmt.Fprintln(w, "----\t-------\t-------\t--------\t")

		for _, s := range stats {
			fmt.Fprintf(w, "%s\t%.2f%%\t%.2f%%\t%.2f%%\t\n",
				s.NodeName, s.CPUUsage, s.MemUsage, s.DiskUsage)
		}
		w.Flush()
	},
}

func init() {
	// Add the node command to the root
	rootCmd.AddCommand(nodeCmd)

	// Define the history flag (local to this command)
	nodeCmd.Flags().StringVarP(&history, "history", "t", "5m", "Time range to look back (e.g., 1h, 1d, 30s)")
}
