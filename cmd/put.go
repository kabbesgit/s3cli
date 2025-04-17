package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put [store] <bucket>/<path> <localfile>",
	Short: "Upload a local file to a bucket/path in a store (via mc)",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		var storeName, bucketAndPath, localFile string
		if len(args) == 3 {
			storeName = args[0]
			bucketAndPath = args[1]
			localFile = args[2]
		} else if len(args) == 2 {
			var err error
			storeName, err = getCurrentStore()
			if err != nil || storeName == "" {
				fmt.Fprintln(os.Stderr, "No store specified and no current store set. Use 's3cli store use <name>' or provide a store argument.")
				os.Exit(1)
			}
			bucketAndPath = args[0]
			localFile = args[1]
		} else {
			fmt.Fprintln(os.Stderr, "Usage: s3cli put [store] <bucket>/<path> <localfile>")
			os.Exit(1)
		}

		// Compose mc alias/bucket/path
		mcPath := fmt.Sprintf("%s/%s", storeName, bucketAndPath)

		mcCmd := exec.Command("mc", "cp", localFile, mcPath)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc cp error: %v\n", err)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get [store] <bucket>/<path> <localfile>",
	Short: "Download a file from a bucket in a store (via mc)",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		var storeName, bucketAndPath, localFile string
		if len(args) == 3 {
			storeName = args[0]
			bucketAndPath = args[1]
			localFile = args[2]
		} else if len(args) == 2 {
			var err error
			storeName, err = getCurrentStore()
			if err != nil || storeName == "" {
				fmt.Fprintln(os.Stderr, "No store specified and no current store set. Use 's3cli store use <name>' or provide a store argument.")
				os.Exit(1)
			}
			bucketAndPath = args[0]
			localFile = args[1]
		} else {
			fmt.Fprintln(os.Stderr, "Usage: s3cli get [store] <bucket>/<path> <localfile>")
			os.Exit(1)
		}

		mcPath := fmt.Sprintf("%s/%s", storeName, bucketAndPath)
		mcCmd := exec.Command("mc", "cp", mcPath, localFile)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc get error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(getCmd)
}
