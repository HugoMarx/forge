package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RootDir string `toml:"root_dir"`
}

func GetConfig() (config Config, err error) {
	fileContent, err := os.ReadFile("../config.toml")
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config: %w", err)
	}
	_, err = toml.Decode(string(fileContent), &config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode config: %w", err)
	}

	return config, nil
}
