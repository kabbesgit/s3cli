package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage S3 stores (mc aliases)",
}

var addStoreCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new store (mc alias set)",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --name flag: %v\n", err)
			os.Exit(1)
		}
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --endpoint flag: %v\n", err)
			os.Exit(1)
		}
		accessKey, err := cmd.Flags().GetString("access-key")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --access-key flag: %v\n", err)
			os.Exit(1)
		}
		secretKey, err := cmd.Flags().GetString("secret-key")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --secret-key flag: %v\n", err)
			os.Exit(1)
		}

		if name == "" || endpoint == "" || accessKey == "" || secretKey == "" {
			fmt.Println("All fields are required")
			os.Exit(1)
		}

		mcCmd := exec.Command("mc", "alias", "set", name, endpoint, accessKey, secretKey)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc alias set error: %v\n", err)
		}
	},
}

var listStoresCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all stores (mc alias list)",
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := getCurrentStore()
		brief, err := cmd.Flags().GetBool("brief")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --brief flag: %v\n", err)
			os.Exit(1)
		}
		mcCmd := exec.Command("mc", "alias", "list")
		output, err := mcCmd.CombinedOutput()
		lines := strings.Split(string(output), "\n")
		if brief {
			var aliasName, url string
			for i := 0; i < len(lines); i++ {
				line := lines[i]
				if strings.TrimSpace(line) == "" {
					continue
				}
				if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
					aliasName = strings.TrimSpace(line)
					url = ""
					// Look ahead for URL line
					for j := i + 1; j < len(lines); j++ {
						if strings.HasPrefix(lines[j], "  URL") || strings.HasPrefix(lines[j], "\tURL") {
							parts := strings.SplitN(lines[j], ":", 2)
							if len(parts) == 2 {
								url = strings.TrimSpace(parts[1])
							}
							break
						}
						if !strings.HasPrefix(lines[j], " ") && !strings.HasPrefix(lines[j], "\t") {
							break
						}
					}
					prefix := "  "
					if current != "" && aliasName == current {
						prefix = "* "
					}
					fmt.Printf("%s%s\t%s\n", prefix, aliasName, url)
				}
			}
			if current == "" {
				fmt.Println("(No current store set)")
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "mc alias list error: %v\n", err)
			}
			return
		}
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			alias := strings.Fields(line)
			if len(alias) > 0 {
				prefix := "  "
				if current != "" && alias[0] == current {
					prefix = "* "
				}
				fmt.Printf("%s%s\n", prefix, line)
			}
		}
		if current == "" {
			fmt.Println("(No current store set)")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc alias list error: %v\n", err)
		}
	},
}

var deleteStoreCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a store (mc alias rm)",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --name flag: %v\n", err)
			os.Exit(1)
		}
		if name == "" {
			fmt.Println("Store name is required")
			os.Exit(1)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --force flag: %v\n", err)
			os.Exit(1)
		}
		if !force {
			fmt.Printf("Are you sure you want to delete the store '%s'? [y/N]: ", name)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		mcCmd := exec.Command("mc", "alias", "rm", name)
		output, err := mcCmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc alias rm error: %v\n", err)
		}
	},
}

var useStoreCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the current store",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Could not find config dir:", err)
			os.Exit(1)
		}
		path := fmt.Sprintf("%s/s3cli/current_store", configDir)
		os.MkdirAll(filepath.Dir(path), 0700)
		err = os.WriteFile(path, []byte(name), 0600)
		if err != nil {
			fmt.Println("Failed to set current store:", err)
			os.Exit(1)
		}
		fmt.Printf("Current store set to '%s'\n", name)
	},
}

var logoutStoreCmd = &cobra.Command{
	Use:   "logout",
	Short: "Unset the current store",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Could not find config dir:", err)
			os.Exit(1)
		}
		path := fmt.Sprintf("%s/s3cli/current_store", configDir)
		err = os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			fmt.Println("Failed to unset current store:", err)
			os.Exit(1)
		}
		fmt.Println("Current store unset.")
	},
}

func getCurrentStore() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/s3cli/current_store", configDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func init() {
	addStoreCmd.Flags().String("name", "", "Name of the store")
	addStoreCmd.Flags().String("endpoint", "", "Endpoint of the store")
	addStoreCmd.Flags().String("access-key", "", "Access key of the store")
	addStoreCmd.Flags().String("secret-key", "", "Secret key of the store")
	deleteStoreCmd.Flags().String("name", "", "Name of the store to delete")
	deleteStoreCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	listStoresCmd.Flags().BoolP("brief", "b", false, "Show only store name and URL")
	storeCmd.AddCommand(addStoreCmd, listStoresCmd, deleteStoreCmd, useStoreCmd, logoutStoreCmd)
	rootCmd.AddCommand(storeCmd)
}
