package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for gotmpl.

To load completions:

Bash:
  $ source <(gotmpl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gotmpl completion bash > /etc/bash_completion.d/gotmpl
  # macOS:
  $ gotmpl completion bash > /usr/local/etc/bash_completion.d/gotmpl

Zsh:
  $ source <(gotmpl completion zsh)

  # To load completions for each session, execute once:
  $ gotmpl completion zsh > "${fpath[1]}/_gotmpl"`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		default:
			fmt.Printf("Unsupported shell type: %s\n", args[0])
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
