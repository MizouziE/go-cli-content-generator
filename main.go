package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

// define types
type Config struct {
	Prompts []string
	Data    string
}

type Table struct {
	Headings []string
	Rows     []map[string]string
}

type Values struct {
	Animal, Mood, Description string
}

func main() {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	// Check for .env file
	envFileCheck(reader)

	fmt.Println("enter config yaml path:")
	configPath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Config path no good!\n\nClosing...")
		return
	}

	// Read config file
	configFile, err := os.ReadFile(strings.TrimSpace(configPath))
	if err != nil {
		fmt.Println("Unable to read config file\n\nClosing...")
		return
	}

	// This will be the config
	config := Config{}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Cannot parse config file\n\nClosing...")
		return
	}

	// Use data path from config to get csv rows
	csvFromConfig, err := os.Open(strings.TrimSpace(config.Data))
	if err != nil {
		fmt.Println("Unable to read csv from config's data path\n\nClosing...")
		fmt.Printf("Error: %v\n", err)
		return
	}

	csvFromConfigReader := csv.NewReader(csvFromConfig)

	headingRow, err := csvFromConfigReader.Read()
	if err != nil {
		fmt.Println("Can't read csv with\n\nClosing...")
		fmt.Println(err)
		return
	}

	// Make empty list
	var promptList []string

	for {
		row, err := csvFromConfigReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Iteration through csv broken!")
			return
		}

		substitutes := make(map[string]string)

		for i := range row {
			substitutes[headingRow[i]] = row[i]
		}

		for _, prompt := range config.Prompts {
			// Create a new template
			tmpl := template.Must(template.New("prompt").Parse(prompt))

			// Execute template
			promptAsBytes := new(bytes.Buffer)
			err := tmpl.Execute(promptAsBytes, substitutes)
			if err != nil {
				fmt.Println("Cannot execute template\n\nClosing...")
				fmt.Println(err)
				return
			}

			promptList = append(promptList, promptAsBytes.String())
		}
	}

	//load env variables - e.g. API Key
	errLoadDotEnv := godotenv.Load()
	if errLoadDotEnv != nil {
		fmt.Println("Error loading .env file\n\nClosing...")
		return
	}

	// Connect with API Key
	c := gogpt.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	ctx := context.Background()

	// Create directory for responses to prompts
	newDirectoryName := "storage/" + time.Now().Format("060102_150405")
	err = os.Mkdir(newDirectoryName, 0777)
	if err != nil {
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
		errWrite := os.WriteFile(strings.TrimSpace(newDirectoryName+"/output-"+strconv.Itoa(index+1)+".md"), []byte(resp.Choices[0].Text), 0777)
		if errWrite != nil {
			fmt.Println("Unable to write new file\n\nClosing...\n\nBecause: ", errWrite)
			return
		}
		bar.Add(1)
	}

	// Save used prompts to same directory
	listAsJSON, _ := json.Marshal(promptList)
	JSONfull := `{"prompts":` + string(listAsJSON) + `}`
	err = os.WriteFile(strings.TrimSpace(newDirectoryName+"/prompts.json"), []byte(JSONfull), 0777)
	if err != nil {
		fmt.Println("Unable to save prompts\n\nClosing...")
		return
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
