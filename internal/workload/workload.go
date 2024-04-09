package workload

import (
	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ListElement struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	LatestRunID string `json:"latest_run_id"`
}

type WorkloadList struct {
	List []ListElement `json:"workloads"`
}

type Workload struct {
	Name        string                 `json:"name"`
	UUID        string                 `json:"uuid"`
	RunID       string                 `json:"run_id"`
	Cloudblock  Cloudblock             `json:"cloudblock"`
	Taget       string                 `json:"target"`
	Description string                 `json:"description"`
	Env         string                 `json:"env"`
	Config      map[string]interface{} `json:"config"`
}

type Module struct {
	Name    string `json:"name"`
	Runtime string `json:"runtime"`
}

type Cloudblock struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Runtime string `json:"runtime"`
}

type TerraformModule interface {
	GetModuleName() string
	GetSource() string
	GetVariables() map[string]string
}

func (m *Workload) GetModuleName() string {
	return m.Cloudblock.Name
}

func (m *Workload) GetTarget() string {
	return m.Taget
}

func (m *Workload) GetRunId() string {
	return m.RunID
}

func (m *Workload) GetRuntime() string {
	moduleStruct, err := getModuleStruct(m.Cloudblock.Name)
	if err != nil {
		return ""
	}
	return moduleStruct.Runtime
}

func (m *Workload) GetEnv() string {
	return m.Env
}

func (m *Workload) GetVersion() string {
	return m.Cloudblock.Version
}

func (m *Workload) GetUUID() string {
	return m.UUID
}

func (m *Workload) GetLatestRunID(configManager config.ConfigManager) string {
	list, err := ListWorkloads(configManager)
	if err != nil {
		return ""
	}

	for _, v := range list.List {
		if v.UUID == m.UUID {
			return v.LatestRunID
		}
	}

	return ""
}

// func (m *Workload) GetSource(configManager config.ConfigManager) (string, error) {
// 	cloudblocksManager := config.NewCloudblocksManager(filepath.Join(configManager.GetModulesDir(), "cloudblocks.json"))
// 	if cloudblock, err := cloudblocksManager.GetCloudblockByName(m.Cloudblock.Name); err == nil {
// 		return cloudblock.Source, nil
// 	}
// 	return "", fmt.Errorf("module not found")
// }

func (m *Workload) GetVariables() map[string]interface{} {
	variables := make(map[string]interface{})

	for key, value := range m.Config {
		switch v := value.(type) {
		case string:
			variables[key] = v
		case map[string]interface{}:
			jsonValue, err := json.Marshal(v)
			if err == nil {
				variables[key] = string(jsonValue)
			}
		default:
			variables[key] = v
		}
	}

	return variables
}

func GenerateEmptyWorkloadList() ([]byte, error) {
	list := WorkloadList{List: []ListElement{}}
	jsonData, err := json.Marshal(list)
	if err != nil {
		return nil, fmt.Errorf("error marshaling workload list: %v", err)
	}
	return jsonData, nil
}

func CheckWorkloadList(configManager config.ConfigManager) bool {
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	_, err := os.Stat(workloadsFile)
	return !os.IsNotExist(err)
}

func UpdateLatestRunID(configManager config.ConfigManager, uuid string, runID string) error {
	list, err := ListWorkloads(configManager)
	if err != nil {
		return err
	}

	found := false
	for i, v := range list.List {
		if v.UUID == uuid {
			list.List[i].LatestRunID = runID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("workload with UUID %s not found", uuid)
	}

	err = WriteWorkloads(configManager, list)
	if err != nil {
		return err
	}

	return nil
}

func WriteWorkloadList(configManager config.ConfigManager) error {
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	if CheckWorkloadList(configManager) {
		return nil
	}

	if !utils.CheckWorkloadDir(configManager) {
		err := utils.CreateWorkloadDir(configManager)
		if err != nil {
			return fmt.Errorf("error creating workload directory: %v", err)
		}
	}

	jsonData, err := GenerateEmptyWorkloadList()
	if err != nil {
		return err
	}

	err = os.WriteFile(workloadsFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error creating workloads.json: %v", err)
	}

	return nil
}

func WriteWorkloads(configManager config.ConfigManager, list WorkloadList) error {
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	jsonData, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("error marshaling workload list: %v", err)
	}

	err = os.WriteFile(workloadsFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing workloads.json: %v", err)
	}

	return nil
}

func InitializeWorkloadList(configManager config.ConfigManager) error {
	if CheckWorkloadList(configManager) {
		return nil
	}

	if !utils.CheckWorkloadDir(configManager) {
		err := utils.CreateWorkloadDir(configManager)
		if err != nil {
			return fmt.Errorf("error creating workload directory: %v", err)
		}
	}

	emptyList := WorkloadList{List: []ListElement{}}
	err := WriteWorkloads(configManager, emptyList)
	if err != nil {
		return err
	}

	return nil
}

func ReadWorkloads(configManager config.ConfigManager) ([]byte, error) {
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	workloads, err := os.ReadFile(workloadsFile)
	if err != nil {
		return nil, fmt.Errorf("error reading workloads.json: %v", err)
	}
	return workloads, nil
}

func ListWorkloads(configManager config.ConfigManager) (WorkloadList, error) {
	wl, err := ReadWorkloads(configManager)
	if err != nil {
		fmt.Printf("Error reading workloads.json: %v\n", err)
		return WorkloadList{}, err
	}

	var list WorkloadList
	err = json.Unmarshal(wl, &list)
	if err != nil {
		return WorkloadList{}, fmt.Errorf("error unmarshaling workloads.json: %v", err)
	}

	return list, nil
}

func AddWorkload(configManager config.ConfigManager, wl *Workload) error {
	list, err := ListWorkloads(configManager)
	if err != nil {
		return err
	}

	if CheckDuplicates(list.List, wl.UUID) {
		return fmt.Errorf("workload already exists")
	}

	newElement := ListElement{Name: wl.Name, UUID: wl.UUID, LatestRunID: wl.RunID}
	list.List = append(list.List, newElement)
	err = WriteWorkloads(configManager, list)
	if err != nil {
		return err
	}

	return nil
}

func DeleteWorkload(configManager config.ConfigManager, uuid string) error {
	list, err := ListWorkloads(configManager)
	if err != nil {
		return err
	}

	if !CheckDuplicates(list.List, uuid) {
		return fmt.Errorf("workload does not exist")
	}

	for i, v := range list.List {
		if v.UUID == uuid {
			list.List = append(list.List[:i], list.List[i+1:]...)
			break
		}
	}

	err = WriteWorkloads(configManager, list)
	if err != nil {
		return err
	}

	return nil
}

func CheckDuplicates(list []ListElement, value string) bool {
	for _, v := range list {
		if v.UUID == value {
			return true
		}
	}
	return false
}

func ParseWorkloadJSON(jsonString string) (Workload, error) {
	// Replace single quotes with double quotes
	jsonString = strings.ReplaceAll(jsonString, "'", "\"")

	// Remove newlines and whitespace characters
	jsonString = strings.ReplaceAll(jsonString, "\n", "")
	jsonString = strings.ReplaceAll(jsonString, "\r", "")
	jsonString = strings.TrimSpace(jsonString)

	var wl Workload
	err := json.Unmarshal([]byte(jsonString), &wl)
	if err != nil {
		return Workload{}, fmt.Errorf("error parsing workload JSON: %v", err)
	}

	if wl.Name == "" || wl.UUID == "" || wl.Cloudblock.Name == "" {
		return Workload{}, fmt.Errorf("invalid workload JSON: missing required fields")
	}

	return wl, nil
}

func getModuleStruct(cloudblockName string) (Module, error) {
	moduleFile := getModuleFilePath(cloudblockName)
	module, err := readModuleFile(moduleFile)
	if err != nil {
		return Module{}, err
	}
	return unmarshalModule(module)
}

func getModuleFilePath(cloudblockName string) string {
	cfmgr := config.NewConfigManager("../../config.json")
	mdir := cfmgr.GetModulesDir()
	return filepath.Join(mdir, cloudblockName, "module.json")
}

func readModuleFile(moduleFile string) ([]byte, error) {
	return os.ReadFile(moduleFile)
}

func unmarshalModule(module []byte) (Module, error) {
	var moduleStruct Module
	err := json.Unmarshal(module, &moduleStruct)
	return moduleStruct, err
}
