package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Store struct {
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type Config struct {
	Stores []Store `json:"stores"`
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "s3cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "s3cli_config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil // No config yet
	} else if err != nil {
		return nil, err
	}
	defer file.Close()
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
