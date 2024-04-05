package processing

import (
	"fmt"
	"os"

	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/utils"
	"cloudblockscli.com/internal/workload"

	// "os"
	"reflect"
	"strings"
	"testing"
)

const work_dir = "../../work/"

func TestProcess(t *testing.T) {
	configManager := config.NewConfigManager("../../config.json")
	tf := &workload.Workload{
		Name:        "test",
		Description: "A test workload",
		RunID:       "abc123",
		Cloudblock: workload.Cloudblock{
			Name:    "ec2",
			Version: "1.0.0",
			Runtime: "terraform",
		},
		Env: "test",
		Config: map[string]interface{}{
			"name":                        "helloworld",
			"subnet_id":                   "test",
			"instance_type":               "linux",
			"key_name":                    "andy",
			"ami":                         "test",
			"tags":                        map[string]string{"name": "test", "env": "test"},
			"associate_public_ip_address": true,
			"user_data":                   "../../script.sh",
		},
	}
	err := utils.CreateWorkloadDir(configManager)
	if err != nil {
		fmt.Printf("workload folder exists\n")
	}
	err = utils.CreateWorkDir(configManager, tf.GetUUID())
	if err != nil {
		fmt.Printf("Trouble creating work dir for %s: %s", tf.GetUUID(), err)
	}
	err = ProcessConfig(configManager, tf)
	if err != nil {
		t.Errorf("got %v want nil", err)
	}

	runDir := configManager.GetWorkDir() + "/" + tf.GetUUID() + "/" + tf.GetRunId()
	if _, err := os.Stat(runDir + "/main.tf"); os.IsNotExist(err) {
		t.Errorf("main.tf does not exist in the run directory")
	}
}


// tests to see if bool and int strings are transformed to bool and int
func TestTransformStringVars(t *testing.T) {
	variable := "true"
	TransformStringVars("", &variable)
	if variable != "true" {
		t.Errorf("got %s, want true", variable)
	}

	variable = "123"
	TransformStringVars("", &variable)
	if variable != "123" {
		t.Errorf("got %s, want 123", variable)
	}

	variable = "hello"
	TransformStringVars("", &variable)
	fmt.Printf("variable: %+v\n", variable)
	if variable != "\"hello\"" {
		t.Errorf("got %s, want \"hello\"", variable)
	}

}

func ParseTerraformMap(variable string, output interface{}) error {
	// Remove the surrounding curly braces and newline characters
	variable = strings.TrimSpace(variable)
	variable = strings.TrimPrefix(variable, "{")
	variable = strings.TrimSuffix(variable, "}")
	variable = strings.ReplaceAll(variable, "\n", "")

	// Split the variable into key-value pairs
	pairs := strings.Split(variable, " ")

	// Create a map to store the parsed key-value pairs
	parsedMap := make(map[string]string)

	// Iterate over the pairs and populate the map
	for _, pair := range pairs {
		if pair != "" {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"") // Remove the surrounding quotes
				parsedMap[key] = value
			}
		}
	}

	// Convert the parsedMap to the desired output type
	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() != reflect.Ptr || outputValue.Elem().Kind() != reflect.Map {
		return fmt.Errorf("output must be a pointer to a map")
	}
	outputValue.Elem().Set(reflect.ValueOf(parsedMap))

	return nil
}

// tests the ParseVariables function to see if it can turn vars from a variables.tf file into a json object
func TestParseVariables(t *testing.T) {
	variables := "../../modules/ec2/variables.tf"
	res := ParseVariables(variables)
	fmt.Printf("res: %+v\n", res)
	if res.Variables[0].Name != "ami" {
		t.Errorf("got %s, want helloworld", res.Variables[0].Name)
	}
}
