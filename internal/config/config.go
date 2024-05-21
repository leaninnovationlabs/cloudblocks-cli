package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Initialized       bool                 `json:"initialized"`
	Workloaddirectory string               `json:"workdir"`
	Modulesdirectory  string               `json:"modulesdir"`
	RootPath          string               `json:"rootpath"`
	Env               map[string]EnvConfig `json:"env"`
	ModulesList       []CloudblockConfig   `json:"modules"`
}

type ConfigManager interface {
	LoadConfig() (Config, error)
	SaveConfig(config Config) error
	InitializeConfig() error
	GetWorkDir() string
	GetModulesDir() string
	GetRootPath() string
	IsInitialized() bool
	GetBucketByEnv(string) string
	GetRegionByEnv(string) string
	UpdateEnvironment(string, string, string) error
	AddEnvironment(string, string, string) error
	DeleteEnvironment(string) error
	ListEnvironments() ([]EnvInfo, error)
	InitializeCloudblocksList() error
	UpdateCloudblocksList(cloudblocks []CloudblockConfig) error
	DeleteCloudblock(name string) error
	GetCloudblockByName(name string) (CloudblockConfig, error)
	GetModuleConfig(moduleName string) (ModuleConfig, error)
	LoadModulesList() ([]CloudblockConfig, error)
}

type ConfigManagerImpl struct {
	configFile string
}

type CloudblockConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ModuleConfig struct {
	Name    string         `json:"name"`
	Runtime string         `json:"runtime"`
	Params  []ModuleParam  `json:"params"`
	Actions []ModuleAction `json:"actions"`
}

type ModuleParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ModuleAction struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Params      []string `json:"params"`
}

type EnvConfig struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type EnvInfo struct {
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type ModulesList struct {
	Cloudblocks []CloudblockConfig `json:"modules"`
}

func NewConfigManager(configFile string) ConfigManager {
	return &ConfigManagerImpl{
		configFile: configFile,
	}
}

func (cm *ConfigManagerImpl) LoadConfig() (Config, error) {
	file, err := os.Open(cm.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, return default config
			path, _ := os.Getwd()
			return Config{
				Initialized:       false,
				Workloaddirectory: path + "/work",
				Modulesdirectory:  path + "/modules",
				RootPath:          path,
				Env:               map[string]EnvConfig{},
				ModulesList:       []CloudblockConfig{},
			}, nil
		}
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cm *ConfigManagerImpl) GetModuleConfig(moduleName string) (ModuleConfig, error) {
    modulesDir := cm.GetModulesDir()
    moduleConfigPath := filepath.Join(modulesDir, moduleName, "module.json")
    fmt.Println(moduleConfigPath)
    fmt.Println(moduleName)

    file, err := os.Open(moduleConfigPath)
    if err != nil {
        return ModuleConfig{}, fmt.Errorf("error opening module.json for module %s: %v", moduleName, err)
    }
    defer file.Close()

    var moduleConfig ModuleConfig
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&moduleConfig)
    if err != nil {
        return ModuleConfig{}, fmt.Errorf("error decoding module.json for module %s: %v", moduleName, err)
    }

    return moduleConfig, nil
}

func (cm *ConfigManagerImpl) GetBucketByEnv(env string) string {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return ""
	}
	envConfig, ok := config.Env[env]
	if !ok {
		fmt.Printf("Environment %s not found in config\n", env)
		return ""
	}
	return envConfig.Bucket
}

func (cm *ConfigManagerImpl) GetRegionByEnv(env string) string {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return ""
	}
	envConfig, ok := config.Env[env]
	if !ok {
		fmt.Printf("Environment %s not found in config\n", env)
		return ""
	}
	return envConfig.Region
}

func (cm *ConfigManagerImpl) AddEnvironment(name, bucket, region string) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	if _, ok := config.Env[name]; ok {
		return fmt.Errorf("environment '%s' already exists", name)
	}

	config.Env[name] = EnvConfig{
		Bucket: bucket,
		Region: region,
	}

	return cm.SaveConfig(config)
}

func (cm *ConfigManagerImpl) ListEnvironments() ([]EnvInfo, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return nil, err
	}

	var environments []EnvInfo
	for name, envConfig := range config.Env {
		environments = append(environments, EnvInfo{
			Name:   name,
			Bucket: envConfig.Bucket,
			Region: envConfig.Region,
		})
	}

	return environments, nil
}

func (cm *ConfigManagerImpl) UpdateEnvironment(name, bucket, region string) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	envConfig, ok := config.Env[name]
	if !ok {
		return fmt.Errorf("environment '%s' not found", name)
	}

	if bucket != "" {
		envConfig.Bucket = bucket
	}
	if region != "" {
		envConfig.Region = region
	}

	config.Env[name] = envConfig

	return cm.SaveConfig(config)
}

func (cm *ConfigManagerImpl) DeleteEnvironment(name string) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	if _, ok := config.Env[name]; !ok {
		return fmt.Errorf("environment '%s' not found", name)
	}

	delete(config.Env, name)

	return cm.SaveConfig(config)
}

func (cm *ConfigManagerImpl) SaveConfig(config Config) error {
	file, err := os.OpenFile(cm.configFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ConfigManagerImpl) InitializeConfig() error {
    config, err := cm.LoadConfig()
    if err != nil {
        return err
    }

    if !config.Initialized {
        config.Initialized = true
        config.RootPath = InstallDir
        config.Workloaddirectory = filepath.Join(InstallDir, "work")
        config.Modulesdirectory = filepath.Join(InstallDir, "modules")
        config.Env = map[string]EnvConfig{
            "dev": {
                Bucket: "",
                Region: "",
            },
        }
        config.ModulesList = []CloudblockConfig{}
        return cm.SaveConfig(config)
    }

    return nil
}

func (cm *ConfigManagerImpl) GetWorkDir() string {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return ""
	}
	return config.Workloaddirectory
}

func (cm *ConfigManagerImpl) GetModulesDir() string {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return ""
	}
	return config.Modulesdirectory
}

func (cm *ConfigManagerImpl) GetRootPath() string {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return ""
	}
	return config.RootPath
}

func (cm *ConfigManagerImpl) IsInitialized() bool {
	config, err := cm.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return false
	}
	return config.Initialized
}

func (cm *ConfigManagerImpl) InitializeCloudblocksList() error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	if config.ModulesList == nil {
		config.ModulesList = []CloudblockConfig{}
		return cm.SaveConfig(config)
	}

	return nil
}

func (cm *ConfigManagerImpl) UpdateCloudblocksList(cloudblocks []CloudblockConfig) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	config.ModulesList = cloudblocks
	return cm.SaveConfig(config)
}

func (cm *ConfigManagerImpl) DeleteCloudblock(name string) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	found := false
	for i, cloudblock := range config.ModulesList {
		if cloudblock.Name == name {
			config.ModulesList = append(config.ModulesList[:i], config.ModulesList[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("cloudblock with name '%s' not found", name)
	}

	return cm.SaveConfig(config)
}

func (cm *ConfigManagerImpl) GetCloudblockByName(name string) (CloudblockConfig, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return CloudblockConfig{}, err
	}

	for _, cloudblock := range config.ModulesList {
		if cloudblock.Name == name {
			return cloudblock, nil
		}
	}

	return CloudblockConfig{}, fmt.Errorf("cloudblock with name '%s' not found", name)
}

func (cm *ConfigManagerImpl) LoadModulesList() ([]CloudblockConfig, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.ModulesList, nil
}
