package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [store] <bucket>/<prefix>",
	Short: "List contents of a bucket or prefix in a store (via mc)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var storeName, bucketAndPrefix string
		if len(args) == 2 {
			storeName = args[0]
			bucketAndPrefix = args[1]
		} else if len(args) == 1 {
			var err error
			storeName, err = getCurrentStore()
			if err != nil || storeName == "" {
				fmt.Fprintln(os.Stderr, "No store specified and no current store set. Use 's3cli store use <name>' or provide a store argument.")
				os.Exit(1)
			}
			bucketAndPrefix = args[0]
		} else {
			fmt.Fprintln(os.Stderr, "Usage: s3cli ls [store] <bucket>/<prefix>")
			os.Exit(1)
		}

		// Compose mc alias/bucket path
		mcPath := fmt.Sprintf("%s/%s", storeName, bucketAndPrefix)

		mcCmd := exec.Command("mc", "ls", mcPath)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc ls error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
