package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var prometheusURL string
var Version = "v1"

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "Sentinel is a CLI tool for HPC cluster monitoring",
	Long:  `A high-performance monitoring tool that fetches metrics from Prometheus for Nodes, Kafka, and OpenSearch.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags are available to every subcommand
	rootCmd.PersistentFlags().StringVar(&prometheusURL, "url", "http://192.168.0.103:9090", "Prometheus server URL")
	rootCmd.Version = Version
}
