package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RootDir string `toml:"root_dir"`
}

func GetConfig() (config Config, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}
	configPath := filepath.Join(configDir, "forge", "config.toml")
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	_, err = toml.Decode(string(fileContent), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
