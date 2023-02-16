package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/schollz/progressbar/v3"
)

func main() {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	// Check for .env file
	envFileCheck(reader)

	//load env variables - e.g. API Key
	errLoadDotEnv := godotenv.Load()
	if errLoadDotEnv != nil {
		fmt.Println("Error loading .env file\n\nClosing...")
		return
	}

	// Ask for path to input file
	fmt.Println("Please provide relative path to input csv file:")
	csvFilePath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to accept file path.\n\nClosing...")
		return
	}

	// Read csv file rows
	csvFile, err := os.ReadFile(strings.TrimSpace(csvFilePath))
	if err != nil {
		fmt.Println("Unable to read given file path.\n\nBeacause of:...")
		fmt.Println(err)
		return
	}

	csvReader := csv.NewReader(strings.NewReader(string(csvFile)))

	rows, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Can't read csv!")
		return
	}

	// Create a list of prompts to run
	var promptList []string

	for _, row := range rows {
		prompt := fmt.Sprintf("Write me a 150 word story about a %s %s that is %s", row[0], row[1], row[2])
		promptList = append(promptList, prompt)
	}

	// Connect with API Key
	c := gogpt.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	ctx := context.Background()

	// Create directory for responses to prompts
	newDirectoryName := "storage/" + time.Now().Format("060102_150405")
	errDir := os.Mkdir(newDirectoryName, 0777)
	if errDir != nil {
		fmt.Println("Cannot create directory\n\nClosing...")
		return
	}

	// Run each prompt from the list thru openai API
	bar := progressbar.Default(int64((len(promptList))))
	for index, prompt := range promptList {

		req := gogpt.CompletionRequest{
			Model:     gogpt.GPT3TextDavinci003,
			MaxTokens: 300,
			Prompt:    prompt,
		}
		resp, err := c.CreateCompletion(ctx, req)
		if err != nil {
			return
		}

		// Write output to file
		errWrite := os.WriteFile(strings.TrimSpace(newDirectoryName+"/story-"+strconv.Itoa(index+1)+".md"), []byte(resp.Choices[0].Text), 0777)
		if errWrite != nil {
			fmt.Println("Unable to write new file\n\nClosing...\n\nBecause: ", errWrite)
			return
		}
		bar.Add(1)
	}
}

func envFileCheck(r *bufio.Reader) {
	if _, err := os.Stat(".env"); err != nil {
		// Create copy of example file if .env does not
		file, err := os.Create(".env")
		if err != nil {
			fmt.Println("Could not create .env file\n\nClosing...")
		}

		defer file.Close()

		// Ask user to input key
		fmt.Println("Please enter OpenAI API key below:")
		apiKey, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Could not retrieve API key.\n\nClosing...")
			return
		}

		w := bufio.NewWriter(file)
		dotEnvInput := "OPEN_AI_API_KEY=" + apiKey
		_, errWriter := w.WriteString(dotEnvInput)
		if errWriter != nil {
			fmt.Println("Failed to write new .env file\n\nClosing...")
			w.Flush()
			return
		}

		w.Flush()

		// Give feedback
		fmt.Println("API Key added, thank you")

		// Exit
		return
	}

	// Greet user
	fmt.Println("Welcome back!")
}
