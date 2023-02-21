package environment

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func EnvFileCheck(r *bufio.Reader) {
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

func PrepareEnvironment() string {
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
