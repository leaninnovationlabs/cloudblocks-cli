package config

import (
    "os"
    "path/filepath"
)

var (
    InstallDir string
    ConfigFile string
)

func init() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        panic("Error getting user home directory: " + err.Error())
    }

    InstallDir = filepath.Join(homeDir, ".cloudblocks")
    ConfigFile = filepath.Join(InstallDir, "config.json")
}