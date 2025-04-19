package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotmpl",
	Short: "A tool for generating files from templates",
	Long: `gotmpl is a command-line tool that helps you generate files from templates.
It supports various template formats and provides a simple way to manage your templates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no arguments are provided, show help
		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		// If arguments are provided but no subcommand, assume "gen" command
		// We pass control to gen command with the same args
		return genCmd.RunE(genCmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
