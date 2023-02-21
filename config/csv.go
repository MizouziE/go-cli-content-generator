package config

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func GetCsv(config Config) (*csv.Reader, []string) {

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

	return csvFileReader, headingRow
}
