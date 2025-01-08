package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HasuraEndpoint string `json:"HasuraEndpoint"`
	AdminSecret    string `json:"AdminSecret"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
