package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Grammar struct {
		GrammarUrl    string `yaml:"url"`
		GrammarApiKey string `yaml:"api_key"`
		Host          string `yaml:"host"`
	} `yaml:"grammar"`

	BotToken string `yaml:"bot_token"`
}

const configPath = "config.yaml"

func ReadConfig() (Config, error) {
	var config Config

	// Open YAML file
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Decode YAML file to struct
	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return Config{}, err
		}
	}

	return config, nil
}
