package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "s3cli",
	Short: "A CLI tool to manage S3-compatible stores",
	Long:  `A CLI tool to manage S3-compatible stores, supporting multiple stores and credentials.`,
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Prints a greeting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from s3cli!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
