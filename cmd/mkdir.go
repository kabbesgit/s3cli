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

var mkdirCmd = &cobra.Command{
	Use:   "mkdir <store> <bucket>/<prefix>",
	Short: "Create a folder (prefix) in a bucket for a store",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		storeName := args[0]
		bucketAndPrefix := args[1]

		// Split bucket and prefix
		parts := strings.SplitN(bucketAndPrefix, "/", 2)
		bucket := parts[0]
		prefix := ""
		if len(parts) > 1 {
			prefix = parts[1]
		}
		if prefix == "" {
			fmt.Println("You must specify a bucket and prefix, e.g. mybucket/myfolder/")
			os.Exit(1)
		}
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
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
		// Create a zero-byte object with the prefix (ending with /)
		_, err = minioClient.PutObject(ctx, bucket, prefix, strings.NewReader(""), 0, minio.PutObjectOptions{})
		if err != nil {
			fmt.Printf("Failed to create folder: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created folder '%s/%s'\n", bucket, prefix)
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}
