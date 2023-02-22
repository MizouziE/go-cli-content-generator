package main

import (
	"gitlab.vlah.sh/intellistage/fintech/content-generator/config"
	"gitlab.vlah.sh/intellistage/fintech/content-generator/environment"
	"gitlab.vlah.sh/intellistage/fintech/content-generator/openai"
)

func main() {

	// Check for .env file
	environment.EnvFileCheck()

	// Read config
	configuration := config.GetConfig()

	// Read csv
	csvFileReader, headingRow := config.GetCsv(configuration)

	// Build prompt list
	promptList := openai.SetPromptList(csvFileReader, headingRow, configuration)

	newDirectoryName := environment.PrepareEnvironment()

	openai.RunOpenAI(promptList, newDirectoryName)
}
