package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var port *int

var rootCmd = &cobra.Command{
	Use:   "precrux",
	Short: "A precrux setup helper tool",
}

func init() {
	port = rootCmd.PersistentFlags().IntP("port", "p", 7777, "incoming port number.  Example: 7777")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra will print the error
		os.Exit(1)
	}
}
