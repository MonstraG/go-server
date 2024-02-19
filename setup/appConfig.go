package setup

import (
	"encoding/json"
	"log"
	"os"
)

const configPath = "config.json"

type AppConfig struct {
	Host           string `json:"host"`
	DatabaseFolder string `json:"databaseFolder"`
}

func ReadConfig() AppConfig {
	configJson, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	var config AppConfig
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	return config
}
