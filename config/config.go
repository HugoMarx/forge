package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Config Configuration

type Configuration struct {
	RootDir string `toml:"root_dir"`
}

func Load() error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	Config = config
	return nil
}

func GetConfig() (config Configuration, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Configuration{}, err
	}
	configPath := filepath.Join(configDir, "forge", "config.toml")
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		return Configuration{}, nil
	}
	meta, err := toml.Decode(string(fileContent), &config)
	if err != nil {
		return Configuration{}, err
	}

	if meta.IsDefined("root_dir") {
		config.RootDir = os.ExpandEnv(config.RootDir)
	}

	return config, nil
}
