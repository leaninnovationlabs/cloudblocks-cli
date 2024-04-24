package cli

import (
	// "bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/executors"
	"cloudblockscli.com/internal/processing"
	"cloudblockscli.com/internal/utils"
	"cloudblockscli.com/internal/workload"

	// "github.com/bytedance/sonic"
	"os"

	"github.com/spf13/cobra"

	// "strings"
	"sync"
	"time"
)

// internal/cli.go

func ExecuteCommand(cmd *cobra.Command, args []string, wg *sync.WaitGroup, resultCh chan<- executors.ExecutorOutput) {
	defer wg.Done()

	configManager := config.NewConfigManager("config.json")
	if !configManager.IsInitialized() {
		fmt.Println("cloudblocks environment not initialized.")
		os.Exit(1)
	}

	var wl workload.Workload
	var err error

	// Check if the --file flag is provided
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		// Read the JSON from the file
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			os.Exit(1)
		}
		err = json.Unmarshal(jsonData, &wl)
		if wl.GetRunId() == "" {
			wl.RunID = wl.Name + "_" +  utils.GenerateRandNum()
		}
		if wl.GetUUID() == "" {
			wl.UUID = wl.Name + "_" +  utils.GenerateRandNum()
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} else if len(args) == 1 {
		// Parse the JSON workload from the command line argument
		err = json.Unmarshal([]byte(args[0]), &wl)
		if wl.GetRunId() == "" {
			wl.RunID = utils.GenerateUUID()
		}
		if wl.GetUUID() == "" {
			wl.UUID = utils.GenerateUUID()
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Please provide a JSON workload either as an argument or using the --file flag.")
		os.Exit(1)
	}

	err = workload.AddWorkload(configManager, &wl)
	if err != nil {
		fmt.Println("Error adding workload:", err)
	}

	if !utils.CheckWorkDir(configManager, wl.UUID) {
		utils.CreateWorkDir(configManager, wl.UUID)
	}

	// process config
	fmt.Println("Getting module config for:", wl)

	err = processing.ProcessConfig(configManager, &wl)
	if err != nil {
		fmt.Println("Error processing config:", err)
		os.Exit(1)
	}

	// Get the module runtime
	fmt.Println("Getting module config for:", wl.GetModuleName())
	moduleConfig, err := configManager.GetModuleConfig(wl.GetModuleName())
	if err != nil {
		fmt.Println("Error getting module config:", err)
		os.Exit(1)
	}

	// Create the executor input
	input := executors.ExecutorInput{WorkloadName: wl.Name, UUID: wl.UUID, RunID: wl.GetRunId()}

	// Execute the workload based on the module runtime
	var res executors.ExecutorOutput
	switch moduleConfig.Runtime {
	case "terraform":
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Second)
		defer cancel()
		res, err = executors.Execute(ctx, input, configManager.GetWorkDir())
	case "cmd":
		workloadDir := filepath.Join(configManager.GetWorkDir(), wl.UUID, wl.GetRunId())
		target := wl.GetAction()
		err = executors.ExecuteMakefile(workloadDir, target)
		if err != nil {
			res = executors.ExecutorOutput{Success: false, Error: err}
		} else {
			res = executors.ExecutorOutput{Success: true}
		}
	default:
		fmt.Printf("Unsupported runtime: %s\n", moduleConfig.Runtime)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error executing workload:", err)
		fmt.Printf("Success: %t\nError: %s\n", res.Success, res.Error)
		os.Exit(1)
	}

	// Update the latest run ID in workloads.json
	err = workload.UpdateLatestRunID(configManager, wl.UUID, wl.GetRunId())
	if err != nil {
		fmt.Println("Error updating latest run ID:", err)
		os.Exit(1)
	}

	wl.UpdateStatus(configManager, "executed")

	// Send the result to the channel
	resultCh <- res
}

func DeleteCommand(cmd *cobra.Command, args []string, wg *sync.WaitGroup, resultCh chan<- executors.ExecutorOutput) {
	defer wg.Done()

	configManager := config.NewConfigManager("config.json")
	if !configManager.IsInitialized() {
		fmt.Println("cloudblocks environment not initialized.")
		os.Exit(1)
	}

	var wl workload.Workload
	var err error
	
	Name, _ := cmd.Flags().GetString("name")
	if Name != "" {
		if workload.GetWorkloadUUIDByName(configManager, Name) != "" {
			wl.Name = Name
			wl.UUID = workload.GetWorkloadUUIDByName(configManager, Name)
		} else {
			fmt.Println("No workload with this name exists")
			os.Exit(1)
		}
	}

	// Check if the --file flag is provided
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		// Read the JSON from the file
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			os.Exit(1)
		}
		err = json.Unmarshal(jsonData, &wl)
		if wl.UUID == "" {
			if workload.GetWorkloadUUIDByName(configManager, wl.Name) != "" {
				wl.UUID = workload.GetWorkloadUUIDByName(configManager, wl.Name)
			} else {
				fmt.Println("No workload with this name exists!")
				os.Exit(1)
			}
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} else if len(args) == 1 {
		// Parse the JSON workload from the command line argument
		err = json.Unmarshal([]byte(args[0]), &wl)
		if wl.UUID == "" {
			if workload.GetWorkloadUUIDByName(configManager, wl.Name) != "" {
				wl.UUID = workload.GetWorkloadUUIDByName(configManager, wl.Name)
			} else {
				fmt.Println("No workload with this name exists!")
				os.Exit(1)
			}
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} 

	if !utils.CheckWorkDir(configManager, wl.UUID) {
		fmt.Printf("%s\n", configManager.GetWorkDir()+"/"+wl.UUID)
		fmt.Println("Workload directory not found.")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	dir := configManager.GetWorkDir() + "/" + wl.UUID

	input := executors.ExecutorInput{WorkloadName: wl.Name, UUID: wl.UUID, RunID: wl.GetLatestRunID(configManager)}

	res, err := executors.Delete(ctx, configManager, input, dir)
	if err != nil {
		fmt.Println("Error deleting workload:", err)
		os.Exit(1)
	}

	// err = workload.DeleteWorkload(configManager, wl.UUID)
	// if err != nil {
	// 	fmt.Println("Error deleting workload:", err)
	// 	os.Exit(1)
	// }

	wl.UpdateStatus(configManager, "deleted")

	// delete workload directory
	// Send the result to the channel
	resultCh <- res
}

func ListCommand(cmd *cobra.Command, wg *sync.WaitGroup, resultCh chan<- executors.ExecutorOutput) {
	defer wg.Done()

	configManager := config.NewConfigManager("config.json")
	if !configManager.IsInitialized() {
		fmt.Println("cloudblocks environment not initialized.")
		resultCh <- executors.ExecutorOutput{Success: false, Error: fmt.Errorf("cloudblocks environment not initialized")}
		os.Exit(1)
	}

	// Load the cloudblocks list
	cloudblocksList, err := configManager.LoadModulesList()
	if err != nil {
		fmt.Println("Error loading cloudblocks list:", err)
		resultCh <- executors.ExecutorOutput{Success: false, Error: err}
		os.Exit(1)
	}

	// Convert the cloudblocks list to JSON
	jsonData, err := json.MarshalIndent(cloudblocksList, "", "  ")
	if err != nil {
		fmt.Println("Error converting cloudblocks list to JSON:", err)
		resultCh <- executors.ExecutorOutput{Success: false, Error: err}
		os.Exit(1)
	}

	// Print the JSON data
	fmt.Println(string(jsonData))

	// Send the result to the channel
	resultCh <- executors.ExecutorOutput{Success: true}
}

func DryRunCommand(cmd *cobra.Command, args []string, wg *sync.WaitGroup, resultCh chan<- executors.ExecutorOutput) {
	defer wg.Done()

	configManager := config.NewConfigManager("config.json")
	if !configManager.IsInitialized() {
		fmt.Println("cloudblocks environment not initialized.")
		os.Exit(1)
	}

	var wl workload.Workload
	var err error
	// var jsonData []byte

	// Check if the --file flag is provided
	// Check if the --file flag is provided
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		// Read the JSON from the file
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			os.Exit(1)
		}
		err = json.Unmarshal(jsonData, &wl)
		if wl.GetRunId() == "" {
			wl.RunID = utils.GenerateUUID()
		}
		if wl.GetUUID() == "" {
			wl.UUID = utils.GenerateUUID()
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} else if len(args) == 1 {
		// Parse the JSON workload from the command line argument
		err = json.Unmarshal([]byte(args[0]), &wl)
		if wl.GetRunId() == "" {
			wl.RunID = utils.GenerateUUID()
		}
		if wl.GetUUID() == "" {
			wl.UUID = utils.GenerateUUID()
		}
		if err != nil {
			fmt.Println("Error parsing workload JSON:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Please provide a JSON workload either as an argument or using the --file flag.")
		os.Exit(1)
	}

	err = workload.AddWorkload(configManager, &wl)
	if err != nil {
		fmt.Println("Error adding workload:", err)
	}

	if !utils.CheckWorkDir(configManager, wl.UUID) {
		utils.CreateWorkDir(configManager, wl.UUID)
	}

	// process config
	err = processing.ProcessConfig(configManager, &wl)
	if err != nil {
		fmt.Println("Error processing config:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Second)
	defer cancel()
	// Execute the workload
	input := executors.ExecutorInput{WorkloadName: wl.Name, UUID: wl.UUID, RunID: wl.GetRunId()}

	res, err := executors.DryRun(ctx, input, configManager.GetWorkDir())
	if err != nil {
		fmt.Println("Error executing workload:", err)
		fmt.Printf("Success: %t\nError: %s\n", res.Success, res.Error)
		os.Exit(1)
	}

	// Update the latest run ID in workloads.json
	err = workload.UpdateLatestRunID(configManager, wl.UUID, wl.GetRunId())
	if err != nil {
		fmt.Println("Error updating latest run ID:", err)
		os.Exit(1)
	}

	// Send the result to the channel
	resultCh <- res
}

func ListWorkloadsCommand(cmd *cobra.Command, wg *sync.WaitGroup, resultCh chan<- executors.ExecutorOutput) {
	defer wg.Done()

	configManager := config.NewConfigManager("config.json")
	if !configManager.IsInitialized() {
		resultCh <- executors.ExecutorOutput{Success: false, Error: fmt.Errorf("cloudblocks environment not initialized")}
		os.Exit(1)
	}
	list, err := workload.ListWorkloads(configManager)
	if err != nil {
		resultCh <- executors.ExecutorOutput{Success: false, Error: err}
		os.Exit(1)
	}

	if len(list.List) == 0 {
		fmt.Println("No workloads found.")
	} else {
		fmt.Println("workloads:")
		for _, workload := range list.List {
			fmt.Println(workload)
		}
	}
	resultCh <- executors.ExecutorOutput{Success: true}
}
