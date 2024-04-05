package utils

import (
	"cloudblockscli.com/internal/config"
	"fmt"
	// "os"
	"testing"
)

// test MUST use config file created on program startup

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	if uuid == "" {
		t.Errorf("got %s want not empty", uuid)
	}
	fmt.Printf("UUID: %+v\n", uuid)
}

func TestCheckModulesDirectory(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	fmt.Printf("configManager: %+v\n", configManager.GetModulesDir())
	// configManager.InitializeConfig()
	val := CheckModulesDirectory(configManager)
	fmt.Printf("val: %v\n", val)
	if !val {
		t.Errorf("expected modules directory to exist")
	}
}

func TestCreateWorkDir(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	configManager.InitializeConfig()

	// Test creating a directory
	err := CreateWorkDir(configManager, "test999")
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
	// delete the directory
	err = DeleteWorkDir(configManager, "test999")
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}

	// Test creating a directory that already exists
	CreateWorkDir(configManager, "pizza")
	err2 := CreateWorkDir(configManager, "pizza")
	if err2 == nil {
		t.Errorf("Expected error, but got nil")
	}
	// delete the directory
	DeleteWorkDir(configManager, "pizza")
}

func TestDeleteWorkDir(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	configManager.InitializeConfig()

	CreateWorkDir(configManager, "test")
	// Test deleting a directory
	err := DeleteWorkDir(configManager, "test")
	if err != nil {
		t.Errorf("Expected nil, but got %v\n", err)
	}
	CreateWorkDir(configManager, "test2")
	// Test deleting a directory that does not exist
	err = DeleteWorkDir(configManager, "test2")
	if err != nil {
		t.Errorf("Error deleting directory: %s\n", err)
	}
}

func TestCheckWorkDir(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	configManager.InitializeConfig()

	err := CreateWorkDir(configManager, "test888")
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
	// Check if the test directory exists
	exists := CheckWorkDir(configManager, "test888")
	if !exists {
		t.Errorf("Expected directory to exist, but got false")
	}
	err2 := DeleteWorkDir(configManager, "test888")
	if err2 != nil {
		t.Errorf("Expected nil, but got %v", err)
	}

	// Check if a non-existent directory exists
	exists = CheckWorkDir(configManager, "non_existant")
	if exists {
		t.Errorf("Expected directory to not exist, but got true")
	}
}

func TestCheckWorkloadDir(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	configManager.InitializeConfig()

	// check if workload directory exists or not
	val := CheckWorkloadDir(configManager)
	if !val {
		t.Errorf("expected workload directory to exist")
	}
}

func TestCreateWorkloadDir(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	configManager.InitializeConfig()

	// Test creating the workload directory
	err := CreateWorkloadDir(configManager)
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}

	// Test creating the workload directory when it already exists
	err = CreateWorkloadDir(configManager)
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}

	// Clean up the created workload directory
	// err = os.RemoveAll(configManager.GetWorkDir())
	// if err != nil {
	// 	t.Errorf("Error cleaning up workload directory: %v", err)
	// }
}
