package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aneeshsunganahalli/SentinelCLI/internal"
	"github.com/spf13/cobra"
)

var history string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Monitor HPC node performance",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := internal.NewPromClient(prometheusURL)
		if err != nil {
			log.Fatalf("Failed to connect to Prometheus: %v", err)
		}

		stats, err := client.GetNodeCPUStats(history)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("\033[1m--- Sentinel Compute Profile (Lookback: %s) ---\033[0m\n\n", history)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintf(w, "NODE\tCPU UTIL (%%)\tLOAD (1m)\tTHROTTLES\tSTATUS\n")
		fmt.Fprintf(w, "----\t------------\t---------\t---------\t------\n")

		for _, s := range stats {
			status := "\033[32mOK\033[0m" // Green OK
			if s.Utilization > 90 || s.Load1m > 30 {
				status = "\033[31mHIGH LOAD\033[0m" // Red HIGH LOAD
			}
			if s.Throttles > 0 {
				status = "\033[33mTHROTTLING\033[0m" // Yellow THROTTLING
			}

			fmt.Fprintf(w, "%s\t%.2f%%\t%.2f\t%.0f\t%s\n",
				s.NodeName, s.Utilization, s.Load1m, s.Throttles, status)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&history, "history", "t", "5m", "Time range to look back (e.g., 1h, 1d, 30s)")
}
