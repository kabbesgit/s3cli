package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kabbesgit/s3cli/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <store> <bucket>/<path>",
	Short: "Remove a file from a bucket/path in a store",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		storeName := args[0]
		bucketAndPath := args[1]

		// Split bucket and object key
		parts := strings.SplitN(bucketAndPath, "/", 2)
		bucket := parts[0]
		objectKey := ""
		if len(parts) > 1 {
			objectKey = parts[1]
		}
		if objectKey == "" {
			fmt.Println("You must specify a bucket and object key, e.g. mybucket/myfile.txt")
			os.Exit(1)
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Are you sure you want to remove '%s/%s'? [y/N]: ", bucket, objectKey)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		var store *config.Store
		for _, s := range cfg.Stores {
			if s.Name == storeName {
				store = &s
				break
			}
		}
		if store == nil {
			fmt.Printf("Store '%s' not found.\n", storeName)
			os.Exit(1)
		}

		minioClient, err := minio.New(store.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(store.AccessKey, store.SecretKey, ""),
			Secure: strings.HasPrefix(store.Endpoint, "https://"),
		})
		if err != nil {
			fmt.Println("Failed to create S3 client:", err)
			os.Exit(1)
		}

		ctx := context.Background()
		err = minioClient.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("Failed to remove object: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed '%s/%s'\n", bucket, objectKey)
	},
}

func init() {
	rmCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	rootCmd.AddCommand(rmCmd)
}
