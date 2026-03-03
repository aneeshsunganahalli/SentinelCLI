package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/vinyas-bharadwaj/sentinel/internal"
)

var watch bool

var serviceCmd = &cobra.Command{
	Use:   "service [kafka|opensearch|logstash]",
	Short: "Monitor specific cluster services",
	Args:  cobra.ExactArgs(1), // Ensure user provides one service name
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		// Define a function that clears the screen and prints the table
		render := func() {
			client, _ := internal.NewPromClient(prometheusURL)
			stats, err := client.GetServiceStats(serviceName)

			// 2. Bold Header with timestamp

			if err != nil {
				fmt.Printf("\033[31mError fetching metrics: %v\033[0m\n", err)
				return
			}

			// 3. Setup Tabwriter with padding for alignment
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintf(w, "ENTITY\tVALUE\tUNIT\tSTATUS\tDESCRIPTION\n")
			fmt.Fprintf(w, "------\t-----\t----\t------\t-----------\n")

			for _, s := range stats {
				// Simple Color Logic
				statusColor := "\033[32m" // Green
				if s.Status == "WARN" || s.Status == "Warning" {
					statusColor = "\033[33m" // Yellow
				} else if s.Status == "CRIT" || s.Status == "ALERT" {
					statusColor = "\033[31m" // Red
				}

				fmt.Fprintf(w, "%s\t%.2f\t%s\t%s%s\033[0m\t%s\n",
					s.MetricName, s.Value, s.Unit, statusColor, s.Status, s.Description)
			}

			fmt.Println()
			w.Flush()
		}

		// Initial render
		fmt.Printf("\033[1m--- Sentinel Live Watch: %s | %s ---\033[0m\n",
			serviceName, time.Now().Format("15:04:05"))
		fmt.Println("Press Ctrl+C to exit")
		render()

		// If watch is enabled, enter the loop
		if watch {
			ticker := time.NewTicker(5 * time.Second)
			for range ticker.C {
				render()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	// Add the watch flag
	serviceCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Refresh the metrics every 5 seconds")
}
