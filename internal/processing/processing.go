package processing

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/utils"
	"cloudblockscli.com/internal/workload"
)

const templateFileName = "/module_template.tpl"

var configManger = config.NewConfigManager("config.json")

type TerraformVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type VariablesJSON struct {
	Variables []TerraformVariable `json:"variables"`
}

func ReadTemplate(configManager config.ConfigManager, moduleName string) ([]byte, error) {
	filePath := configManager.GetModulesDir() + "/" + moduleName + templateFileName
	mainTf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading main.tf: %v", err)
	}
	return mainTf, nil
}

func TransformStringVars(key string, variable *string) {
	if *variable == "true" || *variable == "false" {
		return
	} else if _, err := strconv.Atoi(*variable); err == nil {
		*variable = fmt.Sprintf("%s", *variable)
		return
	} else if key == "tags" {
		// Leave the tags as a JSON string
		return
	} else {
		*variable = fmt.Sprintf("\"%s\"", *variable)
	}
}

func ReplaceVariables(mainTf []byte, variables map[string]interface{}) []byte {
	for k, v := range variables {
		placeholder := "$" + strings.ToUpper(k)

		switch value := v.(type) {
		case string:
			if k != "tags" {
				value = fmt.Sprintf("\"%v\"", v)
			}
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(value))
		case bool:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strconv.FormatBool(value)))
		case float64:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strconv.FormatFloat(value, 'f', -1, 64)))
		case []interface{}:
			listValue := make([]string, len(value))
			for i, item := range value {
				switch itemValue := item.(type) {
				case string:
					listValue[i] = fmt.Sprintf("\"%v\"", itemValue)
				case map[string]interface{}:
					jsonBytes, err := json.Marshal(itemValue)
					if err == nil {
						listValue[i] = string(jsonBytes)
					}
				default:
					listValue[i] = fmt.Sprintf("%v", itemValue)
				}
			}
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(fmt.Sprintf("[%s]", strings.Join(listValue, ", "))))
		case map[string]interface{}:
			if k == "tags" {
				tagPairs := make([]string, 0, len(value))
				for tagKey, tagValue := range value {
					tagPairs = append(tagPairs, fmt.Sprintf("%s = \"%v\"", tagKey, tagValue))
				}
				mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(fmt.Sprintf("{\n%s\n}", strings.Join(tagPairs, ",\n"))))
			} else {
				jsonBytes, err := json.Marshal(value)
				if err == nil {
					mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), jsonBytes)
				}
			}
		default:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(fmt.Sprintf("%v", v)))
		}
	}
	return mainTf
}

func WriteMakeFile(configManager config.ConfigManager, mainTf []byte, workloadName string, runID string) error {
	if !utils.CheckWorkDir(configManager, workloadName) {
		err := utils.CreateWorkDir(configManager, workloadName)
		if err != nil {
			return fmt.Errorf("error creating work directory: %v", err)
		}
	}

	runDir := configManager.GetWorkDir() + "/" + workloadName + "/" + runID
	err := os.MkdirAll(runDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating run directory: %v", err)
	}

	filePath := runDir + "/Makefile"
	fmt.Printf("Writing main.make to %s\n", filePath)
	err = os.WriteFile(filePath, mainTf, 0644)
	if err != nil {
		return fmt.Errorf("error writing main.make: %v", err)
	}
	return nil
}

func WriteMainTf(configManager config.ConfigManager, mainTf []byte, workloadName string, runID string) error {
	if !utils.CheckWorkDir(configManager, workloadName) {
		err := utils.CreateWorkDir(configManager, workloadName)
		if err != nil {
			return fmt.Errorf("error creating work directory: %v", err)
		}
	}

	runDir := configManager.GetWorkDir() + "/" + workloadName + "/" + runID
	err := os.MkdirAll(runDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating run directory: %v", err)
	}

	filePath := runDir + "/main.tf"
	fmt.Printf("Writing main.tf to %s\n", filePath)
	err = os.WriteFile(filePath, mainTf, 0644)
	if err != nil {
		return fmt.Errorf("error writing main.tf: %v", err)
	}
	return nil
}

func ProcessConfig(configManager config.ConfigManager, wl *workload.Workload) error {
	moduleName := wl.GetModuleName()
	workloadName := wl.GetUUID()
	runID := wl.GetRunId()

	moduleConfig, err := configManager.GetModuleConfig(moduleName)
	if err != nil {
		return err
	}

	switch moduleConfig.Runtime {
	case "terraform":
		return processTerraformConfig(configManager, wl, moduleConfig, workloadName, runID)
	case "cmd":
		return processCMDConfig(configManager, wl, moduleConfig, workloadName, runID)
	default:
		return fmt.Errorf("unsupported runtime: %s", moduleConfig.Runtime)
	}
}

func ReplaceVariablesForCmd(mainTf []byte, variables map[string]interface{}) []byte {
	for k, v := range variables {
		placeholder := "$" + strings.ToUpper(k)

		switch value := v.(type) {
		case string:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(value))
		case bool:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strconv.FormatBool(value)))
		case float64:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strconv.FormatFloat(value, 'f', -1, 64)))
		case []interface{}:
			listValue := make([]string, len(value))
			for i, item := range value {
				listValue[i] = fmt.Sprintf("%v", item)
			}
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strings.Join(listValue, " ")))
		case map[string]interface{}:
			mapValue := make([]string, 0, len(value))
			for mapKey, mapVal := range value {
				mapValue = append(mapValue, fmt.Sprintf("%s=%v", mapKey, mapVal))
			}
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(strings.Join(mapValue, " ")))
		default:
			mainTf = bytes.ReplaceAll(mainTf, []byte(placeholder), []byte(fmt.Sprintf("%v", v)))
		}
	}
	return mainTf
}

func processCMDConfig(configManager config.ConfigManager, wl *workload.Workload, moduleConfig config.ModuleConfig, workloadName, runID string) error {
	mainTf, err := ReadTemplate(configManager, moduleConfig.Name)
	if err != nil {
		return err
	}

	variables := wl.GetVariables()
	updatedMake := ReplaceVariablesForCmd(mainTf, variables)
	fmt.Printf("Updated Makefile: %s\n", updatedMake)

	err = WriteMakeFile(configManager, updatedMake, workloadName, runID)
	if err != nil {
		return err
	}

	return nil
}

func addSourceBlock(configManager config.ConfigManager, wl *workload.Workload, mainTf []byte) []byte {
	moduleDir := configManager.GetModulesDir()
	cloudblockName := wl.GetModuleName()
	sourceBlock := fmt.Sprintf(`"%s/%s"`, moduleDir, cloudblockName)
	return bytes.ReplaceAll(mainTf, []byte("$MODULES_SOURCE"), []byte(sourceBlock))
}

func processTerraformConfig(configManager config.ConfigManager, wl *workload.Workload, moduleConfig config.ModuleConfig, workloadName, runID string) error {
	mainTf, err := ReadTemplate(configManager, moduleConfig.Name)
	if err != nil {
		return err
	}

	variables := wl.GetVariables()
	updatedMainTf := ReplaceVariables(mainTf, variables)
	fmt.Printf("Updated main.tf: %s\n", updatedMainTf)

	updatedMainTf = addSourceBlock(configManager, wl, updatedMainTf)
	updatedMainTf = AddBackendBlock(configManager, wl, updatedMainTf)

	err = WriteMainTf(configManager, updatedMainTf, workloadName, runID)
	if err != nil {
		return err
	}

	return nil
}

func AddBackendBlock(configManager config.ConfigManager, wl *workload.Workload, mainTf []byte) []byte {
	env := wl.GetEnv()
	bucket := configManager.GetBucketByEnv(env)
	region := configManager.GetRegionByEnv(env)

	backendBlock := []byte(fmt.Sprintf(`terraform {
        backend "s3" {
            bucket = "%s"
            key    = "tfstate/%s.tfstate"
            region = "%s"
        }
    }
    `, bucket, wl.GetUUID(), region))

	return bytes.Join([][]byte{backendBlock, mainTf}, []byte("\n\n"))
}

func ParseVariables(filepath string) VariablesJSON {
	fileContent := readFileContent(filepath)
	variableBlocks := extractVariableBlocks(fileContent)

	var variables []TerraformVariable
	for _, block := range variableBlocks {
		variable := parseVariableBlock(block)
		variables = append(variables, variable)
	}

	fmt.Printf("*** Variables: %+v\n", variables)

	return VariablesJSON{Variables: variables}
}

func readFileContent(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	return string(content)
}

func extractVariableBlocks(fileContent string) []string {
	variableRegex := regexp.MustCompile(`variable\s+"([^"]+)"\s+{`)
	matches := variableRegex.FindAllStringSubmatch(fileContent, -1)

	var variableBlocks []string
	for _, match := range matches {
		startIndex := strings.Index(fileContent, match[0])
		endIndex := findBlockEnd(fileContent, startIndex)
		block := fileContent[startIndex:endIndex]
		variableBlocks = append(variableBlocks, block)
	}

	return variableBlocks
}

func parseVariableBlock(block string) TerraformVariable {
	nameRegex := regexp.MustCompile(`variable\s+"([^"]+)"`)
	descriptionRegex := regexp.MustCompile(`description\s+=\s+"([^"]+)"`)
	typeRegex := regexp.MustCompile(`type\s+=\s+([^\s]+)`)

	nameMatch := nameRegex.FindStringSubmatch(block)
	descriptionMatch := descriptionRegex.FindStringSubmatch(block)
	typeMatch := typeRegex.FindStringSubmatch(block)

	name := ""
	description := ""
	varType := ""

	if len(nameMatch) > 1 {
		name = nameMatch[1]
	}
	if len(descriptionMatch) > 1 {
		description = descriptionMatch[1]
	}
	if len(typeMatch) > 1 {
		varType = typeMatch[1]
	}

	return TerraformVariable{
		Name:        name,
		Description: description,
		Type:        varType,
	}
}

func findBlockEnd(content string, startIndex int) int {
	braceCount := 1
	for i := startIndex + 1; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				return i + 1
			}
		}
	}
	return len(content)
}
