package main

import (
	"bufio"
	"os"

	"gitlab.vlah.sh/intellistage/fintech/content-generator/config"
	"gitlab.vlah.sh/intellistage/fintech/content-generator/environment"
	"gitlab.vlah.sh/intellistage/fintech/content-generator/openai"
)

func main() {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	// Check for .env file
	environment.EnvFileCheck(reader)

	// Read config
	configuration := config.GetConfig(reader)

	// Read csv
	csvFileReader, headingRow := config.GetCsv(configuration)

	// Build prompt list
	promptList := openai.SetPromptList(csvFileReader, headingRow, configuration)

	newDirectoryName := environment.PrepareEnvironment()

	openai.RunOpenAI(promptList, newDirectoryName)
}
