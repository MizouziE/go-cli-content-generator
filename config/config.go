package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// define types
type Config struct {
	Prompts []string
	Data    string
}

func GetConfig(r *bufio.Reader) Config {
	fmt.Println("enter config yaml path:")
	configPath, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Config path no good!\n\nClosing...")
		panic(err)
	}

	// Read config file
	configFile, err := os.ReadFile(strings.TrimSpace(configPath))
	if err != nil {
		fmt.Println("Unable to read config file\n\nClosing...")
		panic(err)
	}

	// This will be the config
	config := Config{}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Cannot parse config file\n\nClosing...")
		panic(err)
	}

	return config
}
