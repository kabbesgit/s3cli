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

var rmdirCmd = &cobra.Command{
	Use:   "rmdir <store> <bucket>/<prefix>",
	Short: "Remove a folder (prefix) and all its contents from a bucket for a store",
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

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Are you sure you want to remove all objects with prefix '%s/%s'? [y/N]: ", bucket, prefix)
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
		objectCh := minioClient.ListObjects(ctx, bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		})

		var objectsToDelete []minio.ObjectInfo
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println("Error listing objects:", object.Err)
				os.Exit(1)
			}
			objectsToDelete = append(objectsToDelete, object)
		}

		if len(objectsToDelete) == 0 {
			fmt.Printf("No objects found with prefix '%s' in bucket '%s'.\n", prefix, bucket)
			return
		}

		// Remove all objects with the prefix
		for _, obj := range objectsToDelete {
			err := minioClient.RemoveObject(ctx, bucket, obj.Key, minio.RemoveObjectOptions{})
			if err != nil {
				fmt.Printf("Failed to remove object '%s': %v\n", obj.Key, err)
			} else {
				fmt.Printf("Removed '%s/%s'\n", bucket, obj.Key)
			}
		}
		fmt.Printf("Successfully removed folder '%s/%s' and all its contents.\n", bucket, prefix)
	},
}

func init() {
	rmdirCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	rootCmd.AddCommand(rmdirCmd)
}
