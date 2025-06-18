package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: csvtojson <input.csv> <output.json>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]
	jsonFilePath := os.Args[2]

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	header, err := reader.Read()
	if err == io.EOF {
		log.Fatalf("CSV file is empty: %s", csvFilePath)
	}
	if err != nil {
		log.Fatalf("Error reading CSV header: %v", err)
	}

	var records []map[string]string

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV row: %v", err)
		}

		record := make(map[string]string)
		for i, h := range header {
			if i < len(row) {
				record[h] = row[i]
			} else {
				record[h] = "" // Handle cases where a row might have fewer columns than header
			}
		}
		records = append(records, record)
	}

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ") // For pretty-printing JSON

	if err := encoder.Encode(records); err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	fmt.Printf("Successfully converted '%s' to '%s'\n", csvFilePath, jsonFilePath)
}

// Additional implementation at 2025-06-18 00:48:50
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func convertValue(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	lowerS := strings.ToLower(s)
	if lowerS == "null" || lowerS == "none" || lowerS == "n/a" {
		return nil
	}
	return s
}

func main() {
	reader := csv.NewReader(os.Stdin)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			fmt.Fprintln(os.Stderr, "Error: Empty CSV input.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error reading CSV headers: %v\n", err)
		os.Exit(1)
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	var records []map[string]interface{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV row: %v\n", err)
			os.Exit(1)
		}

		record := make(map[string]interface{})
		for i, header := range headers {
			if i < len(row) {
				record[header] = convertValue(row[i])
			} else {
				record[header] = nil
			}
		}
		records = append(records, record)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(records); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-18 00:49:50
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	inputPath := flag.String("i", "", "Input CSV file path (defaults to stdin)")
	outputPath := flag.String("o", "", "Output JSON file path (defaults to stdout)")
	delimiterStr := flag.String("d", ",", "CSV delimiter character (single character)")
	hasHeader := flag.Bool("header", true, "Whether the CSV has a header row (first row used as JSON keys)")
	prettyPrint := flag.Bool("pretty", true, "Pretty print JSON output")

	flag.Parse()

	if len(*delimiterStr) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Delimiter must be a single character.\n")
		os.Exit(1)
	}
	delimiter := rune((*delimiterStr)[0])

	var reader io.Reader
	if *inputPath == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(*inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	var writer io.Writer
	if *outputPath == "" {
		writer = os.Stdout
	} else {
		file, err := os.Create(*outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	err := convertCSVtoJSON(reader, writer, delimiter, *hasHeader, *prettyPrint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during conversion: %v\n", err)
		os.Exit(1)
	}
}

func convertCSVtoJSON(reader io.Reader, writer io.Writer, delimiter rune, hasHeader bool, prettyPrint bool) error {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimiter
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	var headers []string
	if hasHeader {
		var err error
		headers, err = csvReader.Read()
		if err == io.EOF {
			// Empty CSV or only header row, output empty JSON array
			if prettyPrint {
				_, err = writer.Write([]byte("[]\n"))
			} else {
				_, err = writer.Write([]byte("[]"))
			}
			return err
		}
		if err != nil {
			return fmt.Errorf("failed to read header row: %w", err)
		}
	}

	var records []map[string]string

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %w", err)
		}

		record := make(map[string]string)
		if hasHeader {
			for i, value := range row {
				if i < len(headers) {
					record[headers[i]] = value
				} else {
					// If row has more fields than header, use generic names for extra fields
					record[fmt.Sprintf("field_%d", i+1)] = value
				}
			}
			// If row has fewer fields than header, add empty strings for missing header fields
			for i := len(row); i < len(headers); i++ {
				record[headers[i]] = ""
			}
		} else {
			// No header, use field_1, field_2, etc. as keys
			for i, value := range row {
				record[fmt.Sprintf("field_%d", i+1)] = value
			}
		}
		records = append(records, record)
	}

	var jsonData []byte
	var err error
	if prettyPrint {
		jsonData, err = json.MarshalIndent(records, "", "  ")
	} else {
		jsonData, err = json.Marshal(records)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to output: %w", err)
	}

	if prettyPrint {
		_, err = writer.Write([]byte("\n")) // Add a newline for pretty output
	}
	return err
}