package config

import (
    "os"
    "path/filepath"
    "testing"
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "fmt"
)

func setupTestConfigManager(t *testing.T) (ConfigManager, string) {
    tempDir := t.TempDir()
    configFile := filepath.Join(tempDir, "config_test.json")
    return NewConfigManager(configFile), tempDir
}

func TestInitializeConfig(t *testing.T) {
    configManager, _ := setupTestConfigManager(t)

    err := configManager.InitializeConfig()
    assert.NoError(t, err)

    assert.True(t, configManager.IsInitialized())
    assert.NotEmpty(t, configManager.GetWorkDir())
    assert.NotEmpty(t, configManager.GetModulesDir())
    assert.NotEmpty(t, configManager.GetRootPath())
}

func TestSaveAndLoadConfig(t *testing.T) {
    configManager, _ := setupTestConfigManager(t)

    config := Config{
        Initialized:       true,
        Workloaddirectory: "./test_work",
        Modulesdirectory:  "./test_modules",
        RootPath:          "./test_root",
        Env: map[string]EnvConfig{
            "dev": {
                Bucket: "test_bucket",
                Region: "test_region",
            },
        },
    }

    err := configManager.SaveConfig(config)
    assert.NoError(t, err)

    loadedConfig, err := configManager.LoadConfig()
    assert.NoError(t, err)
    assert.Equal(t, config, loadedConfig)
}

func TestEnvironmentOperations(t *testing.T) {
    configManager, _ := setupTestConfigManager(t)

    err := configManager.AddEnvironment("prod", "prod_bucket", "prod_region")
    assert.NoError(t, err)

    err = configManager.UpdateEnvironment("prod", "new_prod_bucket", "")
    assert.NoError(t, err)

    environments, err := configManager.ListEnvironments()
    assert.NoError(t, err)
    assert.Contains(t, environments, EnvInfo{Name: "prod", Bucket: "new_prod_bucket", Region: "prod_region"})

    err = configManager.DeleteEnvironment("prod")
    assert.NoError(t, err)
}

func TestCloudblockOperations(t *testing.T) {
    configManager, _ := setupTestConfigManager(t)

    err := configManager.InitializeCloudblocksList()
    assert.NoError(t, err)

    cloudblocks := []CloudblockConfig{
        {Name: "cloudblock1", Version: "v1"},
        {Name: "cloudblock2", Version: "v2"},
    }
    err = configManager.UpdateCloudblocksList(cloudblocks)
    assert.NoError(t, err)

    loadedCloudblocks, err := configManager.LoadModulesList()
    assert.NoError(t, err)
    assert.Equal(t, cloudblocks, loadedCloudblocks)

    cloudblock, err := configManager.GetCloudblockByName("cloudblock1")
    assert.NoError(t, err)
    assert.Equal(t, "cloudblock1", cloudblock.Name)
    assert.Equal(t, "v1", cloudblock.Version)

    err = configManager.DeleteCloudblock("cloudblock1")
    assert.NoError(t, err)
}

func TestGetModuleConfig(t *testing.T) {
    configManager, tempDir := setupTestConfigManager(t)

    moduleConfig := ModuleConfig{
        Name:    "test_module",
        Runtime: "test_runtime",
        Params: []ModuleParam{
            {Name: "param1", Type: "string", Description: "Test param 1"},
        },
        Actions: []ModuleAction{
            {Name: "action1", Description: "Test action 1", Params: []string{"param1"}},
        },
    }

    moduleDir := filepath.Join(tempDir, "modules", "test_module")
    err := os.MkdirAll(moduleDir, os.ModePerm)
    assert.NoError(t, err)

    file, err := os.Create(filepath.Join(moduleDir, "module.json"))
    assert.NoError(t, err)
    defer file.Close()

    err = configManager.SaveConfig(Config{
        Modulesdirectory: filepath.Join(tempDir, "modules"),
    })
    assert.NoError(t, err)

    err = json.NewEncoder(file).Encode(moduleConfig)
    assert.NoError(t, err)

    loadedModuleConfig, err := configManager.GetModuleConfig("test_module")
    assert.NoError(t, err)
    assert.Equal(t, moduleConfig, loadedModuleConfig)
}
