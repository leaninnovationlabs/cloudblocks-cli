/*
Copyright © 2024 Andrew Fiala andrew.f@leaninnovationlabs.com
*/
package main

import (
	"cloudblockscli.com/cmd"
	"cloudblockscli.com/internal/config"
	"cloudblockscli.com/internal/utils"
	"cloudblockscli.com/internal/workload"
)

func main() {
	configManager := config.NewConfigManager(config.ConfigFile)
	utils.CreateWorkloadDir(configManager)
	workload.InitializeWorkloadList(configManager)
	cmd.Execute()
}