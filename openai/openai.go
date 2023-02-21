package openai

import (
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

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/schollz/progressbar/v3"
	"gitlab.vlah.sh/intellistage/fintech/content-generator/config"
)

func SetPromptList(c *csv.Reader, hr []string, config config.Config) []string {
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

func RunOpenAI(p []string, dir string) {
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
