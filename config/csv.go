package config

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"gitlab.vlah.sh/intellistage/fintech/content-generator/reader"
)

func GetCsv(config Config) (*csv.Reader, []string) {

	args := os.Args
	var csvPath string
	var err error

	if (len(args) < 3) && (config.Data != "") {

		// Use data path from config to get csv rows
		csvPath = config.Data

	} else if config.Data == "" {

		// Ask user for path to csv
		fmt.Println("Please enter csv path:")
		csvPath, err = reader.Reader().ReadString('\n')
		if err != nil {
			fmt.Println("Csv path no good!\n\nClosing...")
			panic(err)
		}

	} else {

		// Set path as second named argument
		csvPath = args[2]
	}

	// Open the csv with determined path
	csvFile, err := os.Open(strings.TrimSpace(csvPath))
	if err != nil {
		fmt.Println("Unable to read csv from provided arguments!\n\nClosing...")
		panic(err)
	}

	csvFileReader := csv.NewReader(csvFile)

	headingRow, err := csvFileReader.Read()
	if err != nil {
		fmt.Println("Can't read csv with\n\nClosing...")
		panic(err)
	}

	return csvFileReader, headingRow
}
