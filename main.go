/*
Copyright © 2024 Andrew Fiala andrew.f@leaninnovationlabs.com
*/
package main

import (
	"cloudblockscli.com/cmd"
	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/utils"
	"cloudblockscli.com/internal/workload"
	// "os"
	// "fmt"
)

func main() {
	utils.CreateWorkloadDir(config.NewConfigManager("config.json"))
	workload.InitializeWorkloadList(config.NewConfigManager("config.json"))
	cmd.Execute()
}
