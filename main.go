package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/u-root/u-root/pkg/cp"
)

func main() {
	// Check for .env file
	if _, err := os.Stat(".env"); err != nil {
		// Create copy of example file if .env does not
		cp.Copy(".env.example", ".env")
		return
	}

	//load env variables - e.g. API Key
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file\n\nClosing...")
		return
	}

	// Connect with API Key
	c := gogpt.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Ada,
		MaxTokens: 5,
		Prompt:    "Are we connected?",
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return
	}

	// Return answer to prompt
	fmt.Println(resp.Choices[0].Text)

}
