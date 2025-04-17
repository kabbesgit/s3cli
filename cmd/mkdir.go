package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir [store] <bucket>/<prefix>",
	Short: "Create a folder (prefix) in a bucket for a store (via mc)",
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
			fmt.Fprintln(os.Stderr, "Usage: s3cli mkdir [store] <bucket>/<prefix>")
			os.Exit(1)
		}

		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		mcPath := fmt.Sprintf("%s/%s", storeName, prefix)
		mcCmd := exec.Command("mc", "cp", "/dev/null", mcPath)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc mkdir error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}
