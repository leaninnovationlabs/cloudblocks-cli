/*
Copyright © 2024 Andrew Fiala andrew.f@leaninnovationlabs.com
*/

package cmd

import (
	"cloudblockscli.com/internal/cli"
	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/executors"
	"cloudblockscli.com/internal/utils"
	// "cloudblockscli.com/internal/workload"
	// "context"
	// "encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sync"
	// "time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudblocks",
	Short: "A CLI tool for managing cloudblocks",
	Long:  `cloudblocks is a command-line interface tool for managing and interacting with cloudblocks.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cloudblocks",
	Long:  `Initializes the cloudblocks environment by creating the work and modules directories and the config.json file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the workdir and modulesdir flags
		workDir, _ := cmd.Flags().GetString("workdir")
		modulesDir, _ := cmd.Flags().GetString("modulesdir")

		// create config manager
		configManager := config.NewConfigManager("config.json")
		err := configManager.InitializeConfig()
		if err != nil {
			fmt.Println("Error initializing config:", err)
			os.Exit(1)
		}

		// update work and modules directories
		cfg, err := configManager.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		cwd, _ := os.Getwd()
		if workDir == "" {
			workDir = cwd + "/" + "work"
		}
		if modulesDir == "" {
			modulesDir = cwd + "/" + "modules"
		}
		cfg.Workloaddirectory = workDir
		cfg.Modulesdirectory = modulesDir
		err = configManager.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}

		err = utils.CreateWorkloadDir(configManager)
		if err != nil {
			fmt.Println("workdir: ", configManager.GetWorkDir())
			os.Exit(1)
		}

		err = configManager.InitializeCloudblocksList()
		if err != nil {
			fmt.Println("Error initializing cloudblocks list:", err)
			os.Exit(1)
		}

		fmt.Printf("Cloudblocks environment initialized successfully.\n")
		fmt.Printf("Work directory: %s\nModules directory: %s\n", workDir, modulesDir)
	},
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage cloudblocks environments",
	Long:  "Commands for managing cloudblocks environments.",
}

var envAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new environment",
	Long:  "Adds a new environment to the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		name, _ := cmd.Flags().GetString("name")
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")

		err := configManager.AddEnvironment(name, bucket, region)
		if err != nil {
			fmt.Println("Error adding environment:", err)
			os.Exit(1)
		}

		fmt.Printf("Environment '%s' added successfully.\n", name)
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing environment",
	Long:  "Updates an existing environment in the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		name, _ := cmd.Flags().GetString("name")
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")

		err := configManager.UpdateEnvironment(name, bucket, region)
		if err != nil {
			fmt.Println("Error updating environment:", err)
			os.Exit(1)
		}

		fmt.Printf("Environment '%s' updated successfully.\n", name)
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available environments",
	Long:  "Lists the available environments from the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		environments, err := configManager.ListEnvironments()
		if err != nil {
			fmt.Println("Error listing environments:", err)
			os.Exit(1)
		}

		if len(environments) == 0 {
			fmt.Println("No environments found.")
		} else {
			fmt.Println("Available environments:")
			for _, env := range environments {
				fmt.Printf("- Name: %s, Bucket: %s, Region: %s\n", env.Name, env.Bucket, env.Region)
			}
		}
	},
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment",
	Long:  "Deletes an environment from the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		name, _ := cmd.Flags().GetString("name")

		err := configManager.DeleteEnvironment(name)
		if err != nil {
			fmt.Println("Error deleting environment:", err)
			os.Exit(1)
		}

		fmt.Printf("Environment '%s' deleted successfully.\n", name)
	},
}

var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage cloudblock modules",
	Long:  "Commands for managing cloudblock modules.",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new cloudblock module",
	Long:  "Adds a new cloudblock module to the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		// Get the name and version from the command line flags
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")

		// Create a new cloudblock config
		cloudblock := config.CloudblockConfig{
			Name:    name,
			Version: version,
		}

		// Load the existing config
		cfg, err := configManager.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		// Append the new cloudblock to the list
		cfg.ModulesList = append(cfg.ModulesList, cloudblock)

		// Save the updated config
		err = configManager.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}

		fmt.Println("Cloudblock module added successfully.")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing cloudblock module",
	Long:  "Updates an existing cloudblock module in the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		// Get the name and version from the command line flags
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")

		// Load the existing config
		cfg, err := configManager.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		// Find the cloudblock with the specified name and update its version
		updated := false
		for i, cloudblock := range cfg.ModulesList {
			if cloudblock.Name == name {
				cfg.ModulesList[i].Version = version
				updated = true
				break
			}
		}

		if !updated {
			fmt.Printf("Cloudblock module with name '%s' not found.\n", name)
			os.Exit(1)
		}

		// Save the updated config
		err = configManager.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}

		fmt.Printf("Cloudblock module '%s' updated successfully.\n", name)
	},
}

var deleteModuleCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cloudblock module",
	Long:  "Deletes a cloudblock module from the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		configManager := config.NewConfigManager("config.json")

		// Get the name from the command line flag
		name, _ := cmd.Flags().GetString("name")

		// Delete the cloudblock
		err := configManager.DeleteCloudblock(name)
		if err != nil {
			fmt.Println("Error deleting cloudblock:", err)
			os.Exit(1)
		}

		fmt.Printf("Cloudblock module '%s' deleted successfully.\n", name)
	},
}

var listWorkloadsCmd = &cobra.Command{
	Use:   "list-workloads",
	Short: "List cloudblock workloads",
	Long:  "Lists the available workloads.",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		resultCh := make(chan executors.ExecutorOutput)

		wg.Add(1)
		go cli.ListWorkloadsCommand(cmd, &wg, resultCh)

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			// Handle the result
			if !res.Success {
				fmt.Println(res)
				os.Exit(1)
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available cloudblocks",
	Long:  "Lists the available cloudblocks from the config.json file.",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		resultCh := make(chan executors.ExecutorOutput)

		wg.Add(1)
		go cli.ListCommand(cmd, &wg, resultCh)

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			// Handle the result
			if !res.Success {
				fmt.Println(res)
				os.Exit(1)
			}
		}
	},
}

// cmd/root.go

var executeCmd = &cobra.Command{
	Use:   "execute [workload JSON]",
	Short: "Execute a workload",
	Long:  `Executes a workload by processing user config into a main.tf file, creates directory, and finally runs terraform init and apply.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		resultCh := make(chan executors.ExecutorOutput)

		var jsonString string
		if len(args) == 1 {
			jsonString = args[0]
		} else {
			filePath, _ := cmd.Flags().GetString("file")
			if filePath != "" {
				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading file: %v\n", err)
					os.Exit(1)
				}
				jsonString = string(content)
			} else {
				fmt.Println("Please provide either a JSON string as an argument or a file path using the --file flag.")
				os.Exit(1)
			}
		}

		wg.Add(1)
		go cli.ExecuteCommand(cmd, []string{jsonString}, &wg, resultCh)

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			// Handle the result
			fmt.Printf("Execution Result:\nSuccess: %t\nError: %v\n", res.Success, res.Error)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [workload JSON]",
	Short: "Delete a workload",
	Long:  `Runs terraform destroy and then deletes the workload directory and its contents.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		resultCh := make(chan executors.ExecutorOutput)

		wg.Add(1)
		go cli.DeleteCommand(cmd, args, &wg, resultCh)

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			// Handle the result
			fmt.Printf("Execution Result:\nSuccess: %t\nError: %v\n", res.Success, res.Error)
		}
	},
}

var dryRunCmd = &cobra.Command{
	Use:   "dry-run [workload JSON]",
	Short: "Perform a dry run of a workload",
	Long:  `Performs a dry run of a workload by processing user config into a main.tf file, creates directory, and runs terraform init and plan.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		resultCh := make(chan executors.ExecutorOutput)

		var jsonString string
		if len(args) == 1 {
			jsonString = args[0]
		} else {
			filePath, _ := cmd.Flags().GetString("file")
			if filePath != "" {
				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading file: %v\n", err)
					os.Exit(1)
				}
				jsonString = string(content)
			} else {
				fmt.Println("Please provide either a JSON string as an argument or a file path using the --file flag.")
				os.Exit(1)
			}
		}

		wg.Add(1)
		go cli.DryRunCommand(cmd, []string{jsonString}, &wg, resultCh)

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			// Handle the result
			fmt.Printf("Dry Run Result:\nSuccess: %t\nError: %v\n", res.Success, res.Error)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(modulesCmd)
	rootCmd.AddCommand(listWorkloadsCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(dryRunCmd)
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envUpdateCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envListCmd)
	executeCmd.Flags().String("file", "", "Path to the JSON workload file")
	deleteCmd.Flags().String("file", "", "Path to the JSON workload file")
	initCmd.Flags().String("workdir", "", "Path to the work directory")
	initCmd.Flags().String("modulesdir", "", "Path to the modules directory")

	modulesCmd.AddCommand(addCmd)
	modulesCmd.AddCommand(deleteModuleCmd)
	modulesCmd.AddCommand(updateCmd)
	addCmd.Flags().String("name", "", "Name of the cloudblock module")
	addCmd.Flags().String("version", "", "Version of the cloudblock module")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("version")

	// addCmd.Flags().String("cloudblock", "", "JSON cloudblock configuration")
	// deletecloudblockCmd.Flags().String("cloudblock", "", "JSON cloudblock configuration")

	envAddCmd.Flags().String("name", "", "Name of the environment")
	envAddCmd.Flags().String("bucket", "", "S3 bucket for the environment")
	envAddCmd.Flags().String("region", "", "AWS region for the environment")
	envAddCmd.MarkFlagRequired("name")
	envAddCmd.MarkFlagRequired("bucket")
	envAddCmd.MarkFlagRequired("region")

	envUpdateCmd.Flags().String("name", "", "Name of the environment")
	envUpdateCmd.Flags().String("bucket", "", "S3 bucket for the environment")
	envUpdateCmd.Flags().String("region", "", "AWS region for the environment")
	envUpdateCmd.MarkFlagRequired("name")

	envDeleteCmd.Flags().String("name", "", "Name of the environment")
	envDeleteCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("cloudblock")
	updateCmd.Flags().String("name", "", "Name of the cloudblock module")
	updateCmd.Flags().String("version", "", "New version of the cloudblock module")
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("version")

	deleteModuleCmd.Flags().String("name", "", "Name of the cloudblock module")
	deleteModuleCmd.MarkFlagRequired("name")

	dryRunCmd.Flags().String("file", "", "Path to the JSON workload file")
}
