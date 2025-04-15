package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kabbesgit/s3cli/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put <store> <bucket>/<path> <localfile>",
	Short: "Upload a local file to a bucket/path in a store",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		storeName := args[0]
		bucketAndPath := args[1]
		localFile := args[2]

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

		file, err := os.Open(localFile)
		if err != nil {
			fmt.Printf("Failed to open local file '%s': %v\n", localFile, err)
			os.Exit(1)
		}
		defer file.Close()

		ctx := context.Background()
		contentType := "application/octet-stream"
		ext := strings.ToLower(filepath.Ext(localFile))
		if ext == ".txt" {
			contentType = "text/plain"
		} // (add more types as needed)

		_, err = minioClient.FPutObject(ctx, bucket, objectKey, localFile, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			fmt.Printf("Failed to upload file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully uploaded '%s' to '%s/%s'\n", localFile, bucket, objectKey)
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}
