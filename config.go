package main

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Apis struct {
		ScstApi string `yaml:"scst_api"`
		ZfsApi  string `yaml:"zfs_api"`
	} `yaml:"apis"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
