package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aneeshsunganahalli/SentinelCLI/internal"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service [kafka|opensearch|logstash]",
	Short: "Monitor CPU load for specific cluster services",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		client, err := internal.NewPromClient(prometheusURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mFailed to connect to Prometheus: %v\033[0m\n", err)
			os.Exit(1)
		}

		stats, err := client.GetServiceCPUStats(serviceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			return
		}

		fmt.Printf("\033[1m--- Sentinel Service Compute Profile: %s ---\033[0m\n\n", serviceName)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintf(w, "ENTITY\tLOAD (1m)\tDESCRIPTION\n")
		fmt.Fprintf(w, "------\t---------\t-----------\n")

		for _, s := range stats {
			fmt.Fprintf(w, "%s\t%.2f\t%s\n", s.EntityName, s.LoadValue, s.Description)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
