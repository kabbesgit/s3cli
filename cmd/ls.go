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

var lsCmd = &cobra.Command{
	Use:   "ls <store> <bucket>/<prefix>",
	Short: "List contents of a bucket or prefix in a store",
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
			Recursive: false,
		})

		for object := range objectCh {
			if object.Err != nil {
				fmt.Println("Error listing objects:", object.Err)
				os.Exit(1)
			}
			fmt.Printf("%s\t%d bytes\n", object.Key, object.Size)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
