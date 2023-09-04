package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogDir string `yaml:"log_dir"`
	Port   int    `yaml:"port"`
}

var C Config

func Load(yamlFile string) {

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading yaml file: %s, error: %v \n\n", yamlFile, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &C)
	if err != nil {
		fmt.Printf("Error parsing yaml data: %s, error: %v \n\n", yamlFile, err)
		os.Exit(1)
	}
}
