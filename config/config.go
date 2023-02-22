package config

import (
	"fmt"
	"os"
	"strings"

	"gitlab.vlah.sh/intellistage/fintech/content-generator/reader"
	"gopkg.in/yaml.v3"
)

// define types
type Config struct {
	Prompts []string
	Data    string
}

func GetConfig() Config {
	args := os.Args
	var configPath string
	var err error

	if len(args) < 2 {

		fmt.Println("Please enter config yaml path:")
		configPath, err = reader.Reader().ReadString('\n')
		if err != nil {
			fmt.Println("Config path no good!\n\nClosing...")
			panic(err)
		}

	} else {

		configPath = os.Args[1]

	}

	// Read config file
	configFile, err := os.ReadFile(strings.TrimSpace(configPath))
	if err != nil {
		fmt.Println("Config path not retrievable. Are you sure it is relative?\n\nClosing...")
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
