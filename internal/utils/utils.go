package utils

import (
    "cloudblockscli.com/internal/config"
    cryptoRand "crypto/rand"
    mathRand "math/rand"
    "strconv"
    "time"
    "fmt"
    "os"
)

// func LoadEnv(configManager config.ConfigManager) error {
// 	err := godotenv.Load(configManager.GetRootPath() + "/.env")
// 	if err != nil {
// 		return fmt.Errorf("error loading .env file: %v", err)
// 	}
// 	return nil
// }

func GenerateUUID() string {
    b := make([]byte, 16)
    _, err := cryptoRand.Read(b)
    if err != nil {
        return ""
    }
    uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return uuid
}

func GenerateRandNum() string {
    mathRand.Seed(time.Now().UnixNano())

    // Generate a random number between 0 and 999
    randomNum := mathRand.Intn(1000)

    // Convert the random number to a string
    randomNumStr := strconv.Itoa(randomNum)

    // Pad the string with leading zeros if necessary
    for len(randomNumStr) < 3 {
        randomNumStr = "0" + randomNumStr
    }

    return randomNumStr
}

func CheckModulesDirectory(configManager config.ConfigManager) bool {
	path := configManager.GetModulesDir()
	fmt.Printf("path: %s\n", path)

	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CheckWorkDir(configManager config.ConfigManager, wlname string) bool {
	path := configManager.GetWorkDir() + "/" + wlname
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DeleteWorkDir(configManager config.ConfigManager, name string) error {
	if !CheckWorkDir(configManager, name) {
		return fmt.Errorf("directory does not exist")
	}

	path := configManager.GetWorkDir() + "/" + name
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("error deleting directory: %v", err)
	}
	return nil
}
func CreateWorkDir(configManager config.ConfigManager, name string) error {
	if CheckWorkDir(configManager, name) {
		return fmt.Errorf("directory already exists")
	}

	path := configManager.GetWorkDir() + "/" + name
	err := os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	return nil
}

func CheckWorkloadDir(configManager config.ConfigManager) bool {
	_, err := os.Stat(configManager.GetWorkDir())
	return !os.IsNotExist(err)
}

func CreateWorkloadDir(configManager config.ConfigManager) error {
	if !CheckWorkloadDir(configManager) {
		err := os.Mkdir(configManager.GetWorkDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating work directory: %v", err)
		}
		fmt.Println("work/ directory created")
	}
	return nil
}
