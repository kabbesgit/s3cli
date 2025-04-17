package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [store] <bucket>/<path>",
	Short: "Remove a file from a bucket/path in a store (via mc)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var storeName, bucketAndPath string
		if len(args) == 2 {
			storeName = args[0]
			bucketAndPath = args[1]
		} else if len(args) == 1 {
			var err error
			storeName, err = getCurrentStore()
			if err != nil || storeName == "" {
				fmt.Fprintln(os.Stderr, "No store specified and no current store set. Use 's3cli store use <name>' or provide a store argument.")
				os.Exit(1)
			}
			bucketAndPath = args[0]
		} else {
			fmt.Fprintln(os.Stderr, "Usage: s3cli rm [store] <bucket>/<path>")
			os.Exit(1)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --force flag: %v\n", err)
			os.Exit(1)
		}
		if !force {
			fmt.Printf("Are you sure you want to remove '%s/%s'? [y/N]: ", storeName, bucketAndPath)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		mcPath := fmt.Sprintf("%s/%s", storeName, bucketAndPath)
		mcCmd := exec.Command("mc", "rm", mcPath)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc rm error: %v\n", err)
		}
	},
}

func init() {
	rmCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	rootCmd.AddCommand(rmCmd)
}
