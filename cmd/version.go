package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display the version of gotmpl and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gotmpl version: " + version) // You may want to replace this with a variable
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
