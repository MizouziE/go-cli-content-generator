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

func main() {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	// Check for .env file
	envFileCheck(reader)

	// Read config
	config := getConfig(reader)

	// Use data path from config to get csv rows
	csvFile, err := os.Open(strings.TrimSpace(config.Data))
	if err != nil {
		fmt.Println("Unable to read csv from config's data path\n\nClosing...")
		panic(err)
	}

	csvFileReader := csv.NewReader(csvFile)

	headingRow, err := csvFileReader.Read()
	if err != nil {
		fmt.Println("Can't read csv with\n\nClosing...")
		panic(err)
	}

	// Build prompt list
	promptList := setPromptList(csvFileReader, headingRow, config)

	newDirectoryName := prepareEnvironment()

	runOpenAI(promptList, newDirectoryName)
}

func envFileCheck(r *bufio.Reader) {
	if _, err := os.Stat(".env"); err != nil {
		// Create copy of example file if .env does not
		file, err := os.Create(".env")
		if err != nil {
			fmt.Println("Could not create .env file\n\nClosing...")
			panic(err)
		}

		defer file.Close()

		// Ask user to input key
		fmt.Println("Please enter OpenAI API key below:")
		apiKey, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Could not retrieve API key.\n\nClosing...")
			panic(err)
		}

		w := bufio.NewWriter(file)
		dotEnvInput := "OPEN_AI_API_KEY=" + apiKey
		_, errWriter := w.WriteString(dotEnvInput)
		if errWriter != nil {
			fmt.Println("Failed to write new .env file\n\nClosing...")
			w.Flush()
			panic(err)
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

func getConfig(r *bufio.Reader) Config {
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

func setPromptList(c *csv.Reader, hr []string, config Config) []string {
	// Make empty list
	var promptList []string

	for {
		row, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Iteration through csv broken!")
			panic(err)
		}

		substitutes := make(map[string]string)

		for i := range row {
			substitutes[hr[i]] = row[i]
		}

		for _, prompt := range config.Prompts {
			// Create a new template
			tmpl := template.Must(template.New("prompt").Parse(prompt))

			// Execute template
			promptAsBytes := new(bytes.Buffer)
			err := tmpl.Execute(promptAsBytes, substitutes)
			if err != nil {
				fmt.Println("Cannot execute template\n\nClosing...")
				panic(err)
			}

			promptList = append(promptList, promptAsBytes.String())
		}

	}
	return promptList

}

func prepareEnvironment() string {
	//load env variables - e.g. API Key
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file\n\nClosing...")
		panic(err)
	}

	// Create directory for responses to prompts
	newDirectoryName := "storage/" + time.Now().Format("060102_150405")
	err = os.Mkdir(newDirectoryName, 0777)
	if err != nil {
		fmt.Println("Cannot create directory\n\nClosing...")
		panic(err)
	}

	return newDirectoryName
}

func runOpenAI(p []string, dir string) {
	// Connect with API Key
	c := gogpt.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	ctx := context.Background()

	// Run each prompt from the list thru openai API
	bar := progressbar.Default(int64((len(p))))
	for index, prompt := range p {

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
		err = os.WriteFile(strings.TrimSpace(dir+"/output-"+strconv.Itoa(index+1)+".md"), []byte(resp.Choices[0].Text), 0777)
		if err != nil {
			fmt.Println("Unable to write new file\n\nClosing...\n\nBecause: ", err)
			return
		}
		bar.Add(1)
	}

	// Save used prompts to same directory
	listAsJSON, _ := json.Marshal(p)
	JSONfull := `{"prompts":` + string(listAsJSON) + `}`
	err := os.WriteFile(strings.TrimSpace(dir+"/prompts.json"), []byte(JSONfull), 0777)
	if err != nil {
		fmt.Println("Unable to save prompts\n\nClosing...")
		return
	}
}
