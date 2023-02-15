package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
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

	// Ask user for prompt
	fmt.Println("Please provide a prompt:")
	userPrompt, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Could not retrieve prompt from user.\n\nClosing...")
		return
	}

	// Connect with API Key
	c := gogpt.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 100,
		Prompt:    userPrompt,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return
	}

	// Return answer to prompt
	fmt.Println(resp.Choices[0].Text)

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
