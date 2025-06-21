package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"isActive"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomUser(id int) User {
	nameSuffix := randomString(5)
	emailSuffix := randomString(8)

	return User{
		ID:       id,
		Name:     fmt.Sprintf("User_%s", nameSuffix),
		Email:    fmt.Sprintf("user_%s@example.com", emailSuffix),
		Age:      rand.Intn(63) + 18,
		IsActive: rand.Intn(2) == 1,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numRecords := 10
	outputFile := ""

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n", "--num":
			if i+1 < len(args) {
				n, err := strconv.Atoi(args[i+1])
				if err == nil && n > 0 {
					numRecords = n
				} else {
					fmt.Fprintln(os.Stderr, "Warning: Invalid number of records, using default 10.")
				}
				i++
			}
		case "-o", "--output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-h", "--help":
			fmt.Println("Usage: go run main.go [-n <num_records>] [-o <output_file>]")
			fmt.Println("  -n, --num    Number of records to generate (default: 10)")
			fmt.Println("  -o, --output Output file path (default: stdout)")
			os.Exit(0)
		}
	}

	users := make([]User, numRecords)
	for i := 0; i < numRecords; i++ {
		users[i] = generateRandomUser(i + 1)
	}

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully generated %d records to %s\n", numRecords, outputFile)
	} else {
		fmt.Println(string(jsonData))
	}
}

// Additional implementation at 2025-06-21 02:38:24
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Field struct {
	Name string
	Type string
}

type Config struct {
	NumRecords     int
	OutputFilePath string
	Schema         []Field
	PrettyPrint    bool
}

func generateRandomValue(fieldType string, r *rand.Rand) interface{} {
	switch fieldType {
	case "string":
		length := r.Intn(11) + 5
		b := make([]byte, length)
		for i := range b {
			b[i] = byte(r.Intn(26) + 'a')
		}
		return string(b)
	case "int":
		return r.Intn(10000)
	case "float":
		return r.Float64() * 1000.0
	case "bool":
		return r.Intn(2) == 0
	case "timestamp":
		now := time.Now()
		oneYearAgo := now.AddDate(-1, 0, 0)
		diff := now.Unix() - oneYearAgo.Unix()
		randomUnix := oneYearAgo.Unix() + r.Int63n(diff)
		return time.Unix(randomUnix, 0).Format(time.RFC3339)
	default:
		return nil
	}
}

func generateData(cfg Config) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, cfg.NumRecords)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	if len(cfg.Schema) == 0 {
		return nil, fmt.Errorf("schema cannot be empty")
	}

	for i := 0; i < cfg.NumRecords; i++ {
		record := make(map[string]interface{})
		for _, field := range cfg.Schema {
			record[field.Name] = generateRandomValue(field.Type, r)
		}
		data[i] = record
	}
	return data, nil
}

func saveToFile(data []map[string]interface{}, cfg Config) error {
	var jsonData []byte
	var err error

	if cfg.PrettyPrint {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = os.WriteFile(cfg.OutputFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file %s: %w", cfg.OutputFilePath, err)
	}
	return nil
}

func splitAndTrim(s, sep string) []string {
	var result []string
	parts := strings.Split(s, sep)
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	var numRecords int
	var outputPath string
	var schemaStr string
	var prettyPrint bool

	flag.IntVar(&numRecords, "n", 10, "Number of records to generate")
	flag.StringVar(&outputPath, "o", "test_data.json", "Output file path")
	flag.StringVar(&schemaStr, "s", "id:int,name:string,age:int,isActive:bool,createdAt:timestamp", "Schema definition (e.g., 'id:int,name:string')")
	flag.BoolVar(&prettyPrint, "pretty", false, "Pretty print JSON output")

	flag.Parse()

	if numRecords <= 0 {
		log.Fatalf("Number of records must be positive.")
	}

	var schema []Field
	fields := splitAndTrim(schemaStr, ",")
	for _, fieldDef := range fields {
		parts := splitAndTrim(fieldDef, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid schema format: %s. Expected 'name:type'.", fieldDef)
		}
		schema = append(schema, Field{Name: parts[0], Type: parts[1]})
	}

	if len(schema) == 0 {
		log.Fatalf("Schema cannot be empty. Please define at least one field.")
	}

	cfg := Config{
		NumRecords:     numRecords,
		OutputFilePath: outputPath,
		Schema:         schema,
		PrettyPrint:    prettyPrint,
	}

	data, err := generateData(cfg)
	if err != nil {
		log.Fatalf("Error generating data: %v", err)
	}

	err = saveToFile(data, cfg)
	if err != nil {
		log.Fatalf("Error saving data: %v", err)
	}

	fmt.Printf("Successfully generated %d records to %s\n", numRecords, outputPath)
}

// Additional implementation at 2025-06-21 02:39:06


// Additional implementation at 2025-06-21 02:39:58
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type DataGenerator interface {
	Generate() map[string]interface{}
	Headers() []string
}

type PersonGenerator struct {
	randSource *rand.Rand
}

func NewPersonGenerator(seed int64) *PersonGenerator {
	return &PersonGenerator{
		randSource: rand.New(rand.NewSource(seed)),
	}
}

func (pg *PersonGenerator) Generate() map[string]interface{} {
	id := pg.randSource.Intn(1000000) + 1
	nameLength := pg.randSource.Intn(10) + 5
	name := generateRandomString(pg.randSource, nameLength)
	age := pg.randSource.Intn(80) + 18
	isActive := pg.randSource.Intn(2) == 1
	createdAt := generateRandomTime(pg.randSource, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())

	return map[string]interface{}{
		"ID":        id,
		"Name":      name,
		"Age":       age,
		"IsActive":  isActive,
		"CreatedAt": createdAt,
	}
}

func (pg *PersonGenerator) Headers() []string {
	return []string{"ID", "Name", "Age", "IsActive", "CreatedAt"}
}

func generateRandomString(r *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomTime(r *rand.Rand, start, end time.Time) time.Time {
	min := start.Unix()
	max := end.Unix()
	delta := max - min
	sec := r.Int63n(delta) + min
	return time.Unix(sec, 0).UTC()
}

type OutputFormat string

const (
	JSON OutputFormat = "json"
	CSV  OutputFormat = "csv"
)

func generateData(n int, generator DataGenerator) []map[string]interface{} {
	records := make([]map[string]interface{}, n)
	for i := 0; i < n; i++ {
		records[i] = generator.Generate()
	}
	return records
}

func writeJSON(w io.Writer, data []map[string]interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func writeCSV(w io.Writer, data []map[string]interface{}, headers []string) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing CSV headers: %w", err)
	}

	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			val := record[header]
			switch v := val.(type) {
			case int:
				row[i] = strconv.Itoa(v)
			case string:
				row[i] = v
			case bool:
				row[i] = strconv.FormatBool(v)
			case time.Time:
				row[i] = v.Format(time.RFC3339)
			default:
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %w", err)
		}
	}
	return nil
}

func main() {
	numRecords := 5
	seed := time.Now().UnixNano()
	outputFormat := JSON

	generator := NewPersonGenerator(seed)

	data := generateData(numRecords, generator)

	outputWriter := os.Stdout

	var err error
	switch outputFormat {
	case JSON:
		err = writeJSON(outputWriter, data)
	case CSV:
		err = writeCSV(outputWriter, data, generator.Headers())
	default:
		err = fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating or writing data: %v\n", err)
		os.Exit(1)
	}
}