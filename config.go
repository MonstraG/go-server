package main

import (
	"encoding/json"
	"log"
	"os"
)

const configPath = "config.json"

type Config struct {
	Host string `json:"host"`
}

func ReadConfig() Config {
	configJson, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	var config Config
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	return config
}
