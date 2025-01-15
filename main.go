package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	// Input folder containing CSV files
	inputFolder := "./csv_files"

	// Output folder for JSON files
	outputFolder := "./output_json_files"
	err := os.MkdirAll(outputFolder, os.ModePerm) // Create the output folder if it doesn't exist
	if err != nil {
		fmt.Println("Error creating output folder:", err)
		return
	}

	// Read all files in the input folder
	files, err := ioutil.ReadDir(inputFolder)
	if err != nil {
		fmt.Println("Error reading input folder:", err)
		return
	}

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Process each CSV file concurrently
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			csvFilePath := filepath.Join(inputFolder, file.Name())
			jsonFilePath := filepath.Join(outputFolder, file.Name()+".json")

			wg.Add(1) // Increment WaitGroup counter
			go func(csvPath, jsonPath string) {
				defer wg.Done() // Decrement WaitGroup counter when done
				err := processCSVFile(csvPath, jsonPath)
				if err != nil {
					fmt.Printf("Error processing file %s: %v\n", csvPath, err)
				} else {
					fmt.Printf("Successfully converted %s to %s\n", csvPath, jsonPath)
				}
			}(csvFilePath, jsonFilePath)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All files processed.")
}

// Function to process a single CSV file
func processCSVFile(csvFilePath, jsonFilePath string) error {
	// Open the CSV file
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer csvFile.Close()

	// Read the CSV file
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// Convert CSV records to JSON
	if len(records) < 1 {
		return fmt.Errorf("no data found in CSV file")
	}

	headers := records[0]
	var jsonData []map[string]string

	for _, row := range records[1:] { // Skip headers
		if len(row) != len(headers) {
			fmt.Println("Skipping row with mismatched column count")
			continue
		}

		recordMap := make(map[string]string)
		for i, value := range row {
			recordMap[headers[i]] = value
		}
		jsonData = append(jsonData, recordMap)
	}

	// Write the JSON data to a file
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(jsonData); err != nil {
		return fmt.Errorf("error encoding JSON data: %w", err)
	}

	return nil
}
