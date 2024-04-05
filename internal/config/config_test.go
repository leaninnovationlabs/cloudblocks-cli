package config

import (
"github.com/stretchr/testify/assert"
"os"
"testing"
)

func TestConfigManager(t *testing.T) {
	configFile := "config_test.json"
	configManager := NewConfigManager(configFile)

	// Test InitializeConfig
	err := configManager.InitializeConfig()
	assert.NoError(t, err)

	// Test SaveConfig
	config := Config{
		Initialized:       true,
		Workloaddirectory: "./test_work",
		Modulesdirectory:  "./test_modules",
		RootPath:          "./test_root",
		Env: map[string]EnvConfig{
			"dev": {
				Bucket: "govpdfsandy",
				Region: "us-east-1",
			},
		},
	}
	err = configManager.SaveConfig(config)
	assert.NoError(t, err)

	// Test LoadConfig
	loadedConfig, err := configManager.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, config, loadedConfig)

	// Test GetWorkDir
	workDir := configManager.GetWorkDir()
	assert.Equal(t, "./test_work", workDir)

	// Test GetModulesDir
	modulesDir := configManager.GetModulesDir()
	assert.Equal(t, "./test_modules", modulesDir)

	// Test GetRootPath
	rootPath := configManager.GetRootPath()
	assert.Equal(t, "./test_root", rootPath)

	// Test IsInitialized
	initialized := configManager.IsInitialized()
	assert.True(t, initialized)

	// Test GetBucketByEnv
	bucket := configManager.GetBucketByEnv("dev")
	assert.Equal(t, "govpdfsandy", bucket)

	// Test GetRegionByEnv
	region := configManager.GetRegionByEnv("dev")
	assert.Equal(t, "us-east-1", region)

	// Clean up the test config file
	err = os.Remove(configFile)
	assert.NoError(t, err)

}

func TestCloudblocksManager(t *testing.T) {
	cloudblocksFile := "cloudblocks_test.json"
	cloudblocksManager := NewCloudblocksManager(cloudblocksFile)

	// Test InitializeCloudblocksList
	err := cloudblocksManager.InitializeCloudblocksList()
	assert.NoError(t, err)

	// Test VerifyCloudblocksList
	exists := cloudblocksManager.VerifyCloudblocksList()
	assert.True(t, exists)

	// Test UpdateCloudblocksList
	cloudblocks := []CloudblockConfig{
		{Name: "cloudblock1", Version: "v1"},
		{Name: "cloudblock2", Version: "v2"},
	}
	err = cloudblocksManager.UpdateCloudblocksList(cloudblocks)
	assert.NoError(t, err)

	// Test GetCloudblockByName
	cloudblock, err := cloudblocksManager.GetCloudblockByName("cloudblock1")
	assert.NoError(t, err)
	assert.Equal(t, "cloudblock1", cloudblock.Name)
	assert.Equal(t, "v1", cloudblock.Version)

	// Test DeleteCloudblock
	err = cloudblocksManager.DeleteCloudblock("cloudblock1")
	assert.NoError(t, err)

	// Verify cloudblock is deleted
	_, err = cloudblocksManager.GetCloudblockByName("cloudblock1")
	assert.Error(t, err)

	// Clean up the test cloudblocks file
	err = os.Remove(cloudblocksFile)
	assert.NoError(t, err)

}
