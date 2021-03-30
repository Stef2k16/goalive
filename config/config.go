package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Notification    Notification `yaml:"notification"`
	URL             []string     `yaml:"url"`
	PollingInterval int          `yaml:"pollingInterval"`
	LogFile         string       `yaml:"logFile"`
}

type Notification struct {
	Client  string `yaml:"client"`
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

// New returns a decoded Config from a YAML file at the given path.
func New(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	if err = file.Close(); err != nil {
		return Config{}, err
	}

	return config, nil
}
