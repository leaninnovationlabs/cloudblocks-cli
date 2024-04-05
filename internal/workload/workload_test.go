package workload

import (
	"cloudblockscli.com/internal/config"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var configManager = config.NewConfigManager("../../config.json")

func TestWorkload_GetRunId(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetRunID",
			m: &Workload{
				Name:  "name",
				RunID: "abc123",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				RunID:       tt.m.RunID,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetRunId(); got != tt.want {
				t.Errorf("Workload.GetRunId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetModuleName(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetModuleName",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetModuleName(); got != tt.want {
				t.Errorf("Workload.GetModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetRuntime(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetRuntime",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "ec2",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "terraform",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetRuntime(); got != tt.want {
				t.Errorf("Workload.GetModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetModuleStruct(t *testing.T) {
	moduleName := "ec2"
	expectedModule := Module{
		Name:    "ec2",
		Runtime: "terraform",
	}

	module, err := getModuleStruct(moduleName)
	// fmt.Println(module)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(module.Runtime, expectedModule.Runtime) {
		t.Errorf("Expected module %+v, but got %+v", expectedModule, module)
	}
}

func TestGetModuleFilePath(t *testing.T) {
	moduleName := "ec2"
	expectedPath := filepath.Join(configManager.GetModulesDir(), moduleName, "module.json")

	path := getModuleFilePath(moduleName)
	fmt.Printf("path: %s\n", path)

	if path != expectedPath {
		t.Errorf("Expected path %s, but got %s", expectedPath, path)
	}
}

func TestReadModuleFile(t *testing.T) {
	moduleName := "ec2"
	moduleFile := getModuleFilePath(moduleName)

	content, err := readModuleFile(moduleFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("Expected non-empty module file content")
	}
}

func TestUnmarshalModule(t *testing.T) {
	jsonData := []byte(`{
		"name": "ec2",
		"runtime": "terraform"
	}`)

	expectedModule := Module{
		Name:    "ec2",
		Runtime: "terraform",
	}

	module, err := unmarshalModule(jsonData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(module, expectedModule) {
		t.Errorf("Expected module %+v, but got %+v", expectedModule, module)
	}
}

func TestWorkload_GetRuntime_EC2(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetRuntime with EC2 module",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "ec2",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "terraform",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetRuntime(); got != tt.want {
				t.Errorf("Workload.GetModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetVersion(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetVersion",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "1.2.3",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "1.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetVersion(); got != tt.want {
				t.Errorf("Workload.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetUUID(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetUUID",
			m: &Workload{
				Name: "name",
				UUID: "uuid",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "uuid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				UUID:        tt.m.UUID,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetUUID(); got != tt.want {
				t.Errorf("Workload.GetUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetEnv(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want string
	}{
		{
			name: "Test GetEnv",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{},
			},
			want: "env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetEnv(); got != tt.want {
				t.Errorf("Workload.GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkload_GetVariables(t *testing.T) {
	tests := []struct {
		name string
		m    *Workload
		want map[string]interface{}
	}{
		{
			name: "Test GetVariables",
			m: &Workload{
				Name: "name",
				Cloudblock: Cloudblock{
					Name:    "name",
					Version: "version",
				},
				Description: "description",
				Env:         "env",
				Config:      map[string]interface{}{"ami": "ami-12345678", "instance_type": "t2.micro"},
			},
			want: map[string]interface{}{"ami": "ami-12345678", "instance_type": "t2.micro"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Workload{
				Name:        tt.m.Name,
				Cloudblock:  tt.m.Cloudblock,
				Description: tt.m.Description,
				Env:         tt.m.Env,
				Config:      tt.m.Config,
			}
			if got := m.GetVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Workload.GetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListWorkloads(t *testing.T) {
	// configManager.InitializeConfig()

	// Create a sample workloads.json file
	jsonData, _ := GenerateEmptyWorkloadList()
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	err := os.WriteFile(workloadsFile, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to create sample workloads.json file: %v", err)
	}
	defer os.Remove(workloadsFile)

	// Call the ListWorkloads function
	list, err := ListWorkloads(configManager)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the returned workload list
	expectedList := WorkloadList{List: []ListElement{}}
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("Expected list %+v, but got %+v", expectedList, list)
	}
}

func TestGenerateEmptyWorkloadList(t *testing.T) {
	// Call the GenerateWorkloadList function
	jsonData, err := GenerateEmptyWorkloadList()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that the returned JSON data is not nil
	if jsonData == nil {
		t.Error("Expected non-nil JSON data, but got nil")
	}

	// Unmarshal the JSON data into a WorkloadList
	var list WorkloadList
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		t.Errorf("Error unmarshaling JSON data: %v", err)
	}

	// Verify that the unmarshaled list matches the expected list
	expectedList := WorkloadList{List: []ListElement{}}
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("Expected list %+v, but got %+v", expectedList, list)
	}
}

func TestCheckWorkloadList(t *testing.T) {
	configManager.InitializeConfig()

	// Call the WriteWorkloadList function
	err := WriteWorkloadList(configManager)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer os.Remove(filepath.Join(configManager.GetWorkDir(), "workloads.json"))

	// Check if the workloads.json file exists
	exists := CheckWorkloadList(configManager)
	if !exists {
		t.Error("Expected workloads.json file to exist, but got false")
	}
}

func TestWriteWorkloads(t *testing.T) {
	configManager.InitializeConfig()

	// Create a sample WorkloadList
	sampleList := WorkloadList{
		List: []ListElement{
			{Name: "test1", UUID: "9324"},
			{Name: "test2", UUID: "38242"},
		},
	}
	// Call the WriteWorkloads function
	err := WriteWorkloads(configManager, sampleList)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer os.Remove(filepath.Join(configManager.GetWorkDir(), "workloads.json"))

	// Check if the workloads.json file exists
	exists := CheckWorkloadList(configManager)
	if !exists {
		t.Error("Expected workloads.json file to exist, but got false")
	}

	// Read the contents of the workloads.json file
	workloadsFile := filepath.Join(configManager.GetWorkDir(), "workloads.json")
	jsonData, err := os.ReadFile(workloadsFile)
	if err != nil {
		t.Errorf("Error reading workloads.json file: %v", err)
	}

	// Unmarshal the JSON data
	var list WorkloadList
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		t.Errorf("Error unmarshaling JSON data: %v", err)
	}

	// Compare the unmarshaled list with the sample list
	if !reflect.DeepEqual(list, sampleList) {
		t.Errorf("Expected list %+v, but got %+v", sampleList, list)
	}

	DeleteWorkload(configManager, "9324")
	DeleteWorkload(configManager, "38242")
}

func TestAddWorkloadToList(t *testing.T) {
	configManager.InitializeConfig()

	err := InitializeWorkloadList(configManager)
	if err != nil {
		t.Errorf("Problem making workloads.json!\n")
	}
	defer os.Remove(filepath.Join(configManager.GetWorkDir(), "workloads.json"))

	wl := Workload{Name: "test2", UUID: "1234"}
	err = AddWorkload(configManager, &wl)
	if err != nil {
		t.Errorf("Error adding workload: %v", err)
	}

	list, _ := ListWorkloads(configManager)
	fmt.Printf("list: %+v\n", list)
	expectedList := WorkloadList{List: []ListElement{{Name: "test2", UUID: wl.UUID}}}
	if !reflect.DeepEqual(list.List, expectedList.List) {
		t.Errorf("Expected list %+v, but got %+v", expectedList, list)
	}
}

func TestDeleteWorkloadFromList(t *testing.T) {
	configManager.InitializeConfig()

	err := InitializeWorkloadList(configManager)
	if err != nil {
		t.Errorf("Problem making workloads.json!\n")
	}
	defer os.Remove(filepath.Join(configManager.GetWorkDir(), "workloads.json"))

	wl := Workload{Name: "test2", UUID: "1234"}
	err = AddWorkload(configManager, &wl)
	if err != nil {
		t.Errorf("Failed to add workload: %v", err)
	}

	// Verify that the workload exists in the list
	list, _ := ListWorkloads(configManager)
	if !CheckDuplicates(list.List, wl.UUID) {
		t.Errorf("Workload %s does not exist in the list", wl.Name)
	}

	// Delete the workload from the list
	err = DeleteWorkload(configManager, wl.UUID)
	if err != nil {
		t.Errorf("Failed to delete workload: %v", err)
	}

	// Verify that the workload is deleted from the list
	list, _ = ListWorkloads(configManager)
	expectedList := WorkloadList{List: []ListElement{}}
	if !reflect.DeepEqual(list.List, expectedList.List) {
		t.Errorf("Expected list %+v, but got %+v", expectedList, list)
	}
}
