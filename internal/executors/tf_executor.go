package executors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"cloudblockscli.com/internal/config"
	// "cloudblockscli.com/internal/utils"
	// "cloudblockscli.com/internal/workload"
)

// ExecutorInput represents the input for the Terraform executor.
type ExecutorInput struct {
	WorkloadName string
	UUID         string
	RunID        string
}

// ExecutorOutput represents the output of the Terraform executor.
type ExecutorOutput struct {
	Success bool
	Error   error
}

func CheckMainFile(workloadDir string, runID string) ExecutorOutput {
	mainTFFile := filepath.Join(workloadDir, runID, "main.tf")
	fmt.Println(mainTFFile)
	_, err := os.Stat(mainTFFile)
	if os.IsNotExist(err) {
		return ExecutorOutput{Success: false, Error: fmt.Errorf("main.tf file not found for workload: %s and run ID: %s", workloadDir, runID)}
	}
	return ExecutorOutput{Success: true}
}

func CopyMainFile(workloadDir string, runID string, newRunID string) error {
	mainTFFile := filepath.Join(workloadDir, runID, "main.tf")
	newMainTFFile := filepath.Join(workloadDir, newRunID, "main.tf")

	sourceFile, err := os.Open(mainTFFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(newMainTFFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create the destination file for writing
	destinationFile, err := os.Create(newMainTFFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destinationFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

func createLogFile(workloadDir string, uuid string) (*os.File, error) {
	logsDir := filepath.Join(workloadDir, "logs")
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	logFile := filepath.Join(logsDir, uuid+".log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return file, nil
}

func RunTfDestroy(ctx context.Context, workloadDir string, input ExecutorInput) ExecutorOutput {
	logFile, err := createLogFile(workloadDir, input.RunID)
	if err != nil {
		return ExecutorOutput{Success: false, Error: err}
	}
	defer logFile.Close()

	destroyCmd := exec.CommandContext(ctx, "terraform", "destroy", "-auto-approve")
	destroyCmd.Dir = filepath.Join(workloadDir, input.RunID)
	destroyCmd.Stdout = logFile
	destroyCmd.Stderr = logFile
	err = destroyCmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform destroy cancelled: %v", err)}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform destroy timed out: %v", err)}
		} else {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("failed to run terraform destroy: %v", err)}
		}
	}
	return ExecutorOutput{Success: true}
}

func RunTfInit(ctx context.Context, workloadDir string, input ExecutorInput) ExecutorOutput {
	logFile, err := createLogFile(workloadDir, input.RunID)
	if err != nil {
		return ExecutorOutput{Success: false, Error: err}
	}
	defer logFile.Close()

	initCmd := exec.CommandContext(ctx, "terraform", "init")
	initCmd.Dir = filepath.Join(workloadDir, input.RunID)
	initCmd.Stdout = logFile
	initCmd.Stderr = logFile
	err = initCmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform init cancelled: %v", err)}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform init timed out: %v", err)}
		} else {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("failed to run terraform init: %v", err)}
		}
	}
	return ExecutorOutput{Success: true}
}

func RunTfPlan(ctx context.Context, workloadDir string, input ExecutorInput) ExecutorOutput {
	logFile, err := createLogFile(workloadDir, input.RunID)
	if err != nil {
		return ExecutorOutput{Success: false, Error: err}
	}
	defer logFile.Close()

	applyCmd := exec.CommandContext(ctx, "terraform", "plan")
	applyCmd.Dir = filepath.Join(workloadDir, input.RunID)
	applyCmd.Stdout = logFile
	applyCmd.Stderr = logFile
	err = applyCmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform plan cancelled: %v", err)}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform plan timed out: %v", err)}
		} else {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("failed to run terraform plan: %v", err)}
		}
	}
	return ExecutorOutput{Success: true}
}

func RunTfApply(ctx context.Context, workloadDir string, input ExecutorInput) ExecutorOutput {
	logFile, err := createLogFile(workloadDir, input.RunID)
	if err != nil {
		return ExecutorOutput{Success: false, Error: err}
	}
	defer logFile.Close()

	applyCmd := exec.CommandContext(ctx, "terraform", "apply", "-auto-approve")
	applyCmd.Dir = filepath.Join(workloadDir, input.RunID)
	applyCmd.Stdout = logFile
	applyCmd.Stderr = logFile
	err = applyCmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform apply cancelled: %v", err)}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("terraform apply timed out: %v", err)}
		} else {
			return ExecutorOutput{Success: false, Error: fmt.Errorf("failed to run terraform apply: %v", err)}
		}
	}
	return ExecutorOutput{Success: true}
}

// Runs terraform plan to try the configuration without actually creating any resources.
func DryRun(ctx context.Context, input ExecutorInput, workdir string) (ExecutorOutput, error) {
	workloadDir := filepath.Join(workdir, input.UUID)
	result := CheckMainFile(workloadDir, input.RunID)

	// Check if the main.tf file exists
	if !result.Success {
		return result, result.Error
	}

	// Run "terraform init"
	result = RunTfInit(ctx, workloadDir, input)
	if !result.Success {
		return result, result.Error
	}

	// Run "terraform plan"
	result = RunTfPlan(ctx, workloadDir, input)
	if !result.Success {
		return result, result.Error
	}

	return ExecutorOutput{Success: true}, nil
}

// Execute runs the Terraform commands for the specified workload.
func Execute(ctx context.Context, input ExecutorInput, workdir string) (ExecutorOutput, error) {
	workloadDir := filepath.Join(workdir, input.UUID)
	fmt.Println(workloadDir)
	result := CheckMainFile(workloadDir, input.RunID)

	// Check if the main.tf file exists
	if !result.Success {
		return result, result.Error
	}

	// Run "terraform init"
	result = RunTfInit(ctx, workloadDir, input)
	if !result.Success {
		return result, result.Error
	}

	// Run "terraform apply"
	result = RunTfApply(ctx, workloadDir, input)
	if !result.Success {
		return result, result.Error
	}

	return ExecutorOutput{Success: true}, nil
}

func Delete(ctx context.Context, cfmgr config.ConfigManager, input ExecutorInput, workdir string) (ExecutorOutput, error) {
	fmt.Println(workdir)
	result := CheckMainFile(workdir, input.RunID)

	// Check if the main.tf file exists
	if !result.Success {
		return result, result.Error
	}

	// Run "terraform destroy"
	result = RunTfDestroy(ctx, workdir, input)
	if !result.Success {
		return result, result.Error
	}

	return ExecutorOutput{Success: true}, nil
}
