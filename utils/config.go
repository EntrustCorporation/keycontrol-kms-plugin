package utils

import (
	"encoding/json"
	"io"
	"kms-plugin/models"
	"log"
	"os"
)

func LoadConfig(configPath string) models.Config {
	configJson, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer configJson.Close()

	jsonVal, err := io.ReadAll(configJson)
	if err != nil {
		log.Fatal(err)
	}
	var config models.Config
	err = json.Unmarshal(jsonVal, &config)
	if err != nil {
	}
	return config
}
