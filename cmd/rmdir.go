package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rmdirCmd = &cobra.Command{
	Use:   "rmdir [store] <bucket>/<prefix>",
	Short: "Remove a folder (prefix) and all its contents from a bucket for a store (via mc)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var storeName, prefix string
		if len(args) == 2 {
			storeName = args[0]
			prefix = args[1]
		} else if len(args) == 1 {
			var err error
			storeName, err = getCurrentStore()
			if err != nil || storeName == "" {
				fmt.Fprintln(os.Stderr, "No store specified and no current store set. Use 's3cli store use <name>' or provide a store argument.")
				os.Exit(1)
			}
			prefix = args[0]
		} else {
			fmt.Fprintln(os.Stderr, "Usage: s3cli rmdir [store] <bucket>/<prefix>")
			os.Exit(1)
		}

		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --force flag: %v\n", err)
			os.Exit(1)
		}
		if !force {
			fmt.Printf("Are you sure you want to remove all objects with prefix '%s/%s'? [y/N]: ", storeName, prefix)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		mcPath := fmt.Sprintf("%s/%s", storeName, prefix)
		mcCmd := exec.Command("mc", "rm", "--recursive", "--force", mcPath)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc rmdir error: %v\n", err)
		}
	},
}

func init() {
	rmdirCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	rootCmd.AddCommand(rmdirCmd)
}
