package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/yourusername/s3cli/config"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage S3 stores",
}

var addStoreCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new store",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		accessKey, _ := cmd.Flags().GetString("access-key")
		secretKey, _ := cmd.Flags().GetString("secret-key")

		if name == "" || endpoint == "" || accessKey == "" || secretKey == "" {
			fmt.Println("All fields are required")
			os.Exit(1)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		for _, s := range cfg.Stores {
			if s.Name == name {
				fmt.Println("A store with that name already exists.")
				os.Exit(1)
			}
		}

		// Validate store by trying to list buckets
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false, // set to true if using https
		})
		if err != nil {
			fmt.Println("Failed to create S3 client:", err)
			os.Exit(1)
		}
		ctx := context.Background()
		_, err = minioClient.ListBuckets(ctx)
		if err != nil {
			fmt.Println("Failed to validate S3 credentials or endpoint:", err)
			os.Exit(1)
		}

		cfg.Stores = append(cfg.Stores, config.Store{
			Name:      name,
			Endpoint:  endpoint,
			AccessKey: accessKey,
			SecretKey: secretKey,
		})

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}

		fmt.Printf("Store '%s' added successfully.\n", name)
	},
}

var listStoresCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stores",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		if len(cfg.Stores) == 0 {
			fmt.Println("No stores configured.")
			return
		}
		for _, s := range cfg.Stores {
			fmt.Printf("Name: %s\n  Endpoint: %s\n  Access Key: %s\n", s.Name, s.Endpoint, s.AccessKey)
		}
	},
}

var deleteStoreCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a store",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Println("Store name is required")
			os.Exit(1)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		found := false
		newStores := make([]config.Store, 0, len(cfg.Stores))
		for _, s := range cfg.Stores {
			if s.Name == name {
				found = true
				continue
			}
			newStores = append(newStores, s)
		}
		if !found {
			fmt.Println("No store found with that name.")
			os.Exit(1)
		}
		cfg.Stores = newStores
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}
		fmt.Printf("Store '%s' deleted successfully.\n", name)
	},
}

func init() {
	addStoreCmd.Flags().String("name", "", "Name of the store")
	addStoreCmd.Flags().String("endpoint", "", "Endpoint of the store")
	addStoreCmd.Flags().String("access-key", "", "Access key of the store")
	addStoreCmd.Flags().String("secret-key", "", "Secret key of the store")
	deleteStoreCmd.Flags().String("name", "", "Name of the store to delete")
	storeCmd.AddCommand(addStoreCmd, listStoresCmd, deleteStoreCmd)
	rootCmd.AddCommand(storeCmd)
}
