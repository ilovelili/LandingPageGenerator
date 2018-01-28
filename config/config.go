// Package config configuration related
package config

import (
	"encoding/json"
	"os"
	"path"
)

// GetConfig get config defined in config.json
func GetConfig() (config *Config, err error) {
	pwd, _ := os.Getwd()
	path := path.Join(pwd, "config.json")
	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		return
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return
}

// FTP ftp entry
type FTP struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Config config entry
type Config struct {
	FTP `json:"ftp"`
}
