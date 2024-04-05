package executors

import (
	"context"
	"testing"
	"time"
)

func TestCheckMainFile(t *testing.T) {
	// Test checking if main.tf exists
	val := CheckMainFile("test", "abc123")
	if val.Success != false {
		t.Errorf("got %v want nil", val)
	}
}

func TestRunTfInit(t *testing.T) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Specify the workload directory and run ID for testing
	workloadDir := "../../work/test"
	runID := "abc123"

	// Call the RunTfInit function with the test context, workload directory, and run ID
	result := RunTfInit(ctx, workloadDir, runID)

	// Assert the expected behavior and output
	if !result.Success {
		t.Errorf("Expected success, but got error: %v", result.Error)
	}
}

func TestRunTfApply(t *testing.T) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Specify the workload directory and run ID for testing
	workloadDir := "../../work/test"
	runID := "abc123"

	// Call the RunTfApply function with the test context, workload directory, and run ID
	result := RunTfApply(ctx, workloadDir, runID)

	// Assert the expected behavior and output
	if !result.Success {
		t.Errorf("Expected success, but got error: %v", result.Error)
	}
}

func TestRunTfDestroy(t *testing.T) {
	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Specify the workload directory and run ID for testing
	workloadDir := "../../work/test"
	runID := "abc123"

	// Call the RunTfDestroy function with the test context, workload directory, and run ID
	result := RunTfDestroy(ctx, workloadDir, runID)

	// Assert the expected behavior and output
	if !result.Success {
		t.Errorf("Expected success, but got error: %v", result.Error)
	}
}
