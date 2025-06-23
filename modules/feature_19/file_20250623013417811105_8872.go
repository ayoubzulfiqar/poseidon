package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type ColumnConfig struct {
	Name  string `json:"name"`
	Width int    `json:"width"`
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config_file.json> <input_file.csv> <output_file.txt>\n", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	configData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file %s: %v\n", configFile, err)
		os.Exit(1)
	}

	var configs []ColumnConfig
	err = json.Unmarshal(configData, &configs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config JSON: %v\n", err)
		os.Exit(1)
	}

	csvFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input CSV file %s: %v\n", inputFile, err)
		os.Exit(1)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	fwFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output fixed-width file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	defer fwFile.Close()

	header, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			fmt.Fprintf(os.Stderr, "Input CSV file is empty: %s\n", inputFile)
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error reading CSV header: %v\n", err)
		os.Exit(1)
	}

	colIndexMap := make(map[string]int)
	for i, colName := range header {
		colIndexMap[colName] = i
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading CSV record: %v\n", err)
			os.Exit(1)
		}

		var sb strings.Builder
		for _, cfg := range configs {
			colIdx, ok := colIndexMap[cfg.Name]
			if !ok {
				fmt.Fprintf(os.Stderr, "Column '%s' not found in CSV header. Using empty string for this column.\n", cfg.Name)
				sb.WriteString(formatField("", cfg.Width))
				continue
			}

			if colIdx >= len(record) {
				fmt.Fprintf(os.Stderr, "Record has fewer columns than expected for column '%s'. Using empty string.\n", cfg.Name)
				sb.WriteString(formatField("", cfg.Width))
				continue
			}

			fieldValue := record[colIdx]
			formattedField := formatField(fieldValue, cfg.Width)
			sb.WriteString(formattedField)
		}
		_, err = fwFile.WriteString(sb.String() + "\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
			os.Exit(1)
		}
	}
}

func formatField(value string, width int) string {
	runes := []rune(value)
	currentLen := len(runes)

	if currentLen > width {
		return string(runes[:width])
	} else if currentLen < width {
		padding := width - currentLen
		var sb strings.Builder
		sb.WriteString(value)
		for i := 0; i < padding; i++ {
			sb.WriteRune(' ')
		}
		return sb.String()
	}
	return value
}

// Additional implementation at 2025-06-23 01:34:56
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func padString(s string, width int, padLeft bool, padChar rune) string {
	currentWidth := utf8.RuneCountInString(s)
	if currentWidth >= width {
		runes := []rune(s)
		return string(runes[:width])
	}
	padding := strings.Repeat(string(padChar), width-currentWidth)
	if padLeft {
		return padding + s
	}
	return s + padding
}

func parseWidths(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	widths := make([]int, len(parts))
	for i, p := range parts {
		width, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid width '%s': %w", p, err)
		}
		if width < 0 {
			return nil, fmt.Errorf("width cannot be negative: %d", width)
		}
		widths[i] = width
	}
	return widths, nil
}

func main() {
	var (
		inputFile   string
		outputFile  string
		delimiter   string
		hasHeader   bool
		widthsStr   string
		padCharStr  string
		padLeft     bool
	)

	flag.StringVar(&inputFile, "input", "", "Input CSV file path (required)")
	flag.StringVar(&outputFile, "output", "", "Output fixed-width file path (required)")
	flag.StringVar(&delimiter, "delimiter", ",", "CSV delimiter character")
	flag.BoolVar(&hasHeader, "header", true, "Set to false if the CSV does not have a header row")
	flag.StringVar(&widthsStr, "widths", "", "Comma-separated list of fixed widths for columns (e.g., \"10,20,5\"). If empty, widths are auto-detected.")
	flag.StringVar(&padCharStr, "padchar", " ", "Character used for padding (default is space)")
	flag.BoolVar(&padLeft, "padleft", false, "Pad on the left side of the column (default is right padding)")

	flag.Parse()

	if inputFile == "" || outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Both -input and -output flags are required.\n")
		flag.Usage()
		os.Exit(1)
	}

	if len(delimiter) != 1 {
		fmt.Fprintf(os.Stderr, "Error: -delimiter must be a single character.\n")
		os.Exit(1)
	}

	if len(padCharStr) != 1 {
		fmt.Fprintf(os.Stderr, "Error: -padchar must be a single character.\n")
		os.Exit(1)
	}
	padChar := rune(padCharStr[0])

	inF, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file '%s': %v\n", inputFile, err)
		os.Exit(1)
	}
	defer inF.Close()

	csvReader := csv.NewReader(inF)
	csvReader.Comma = rune(delimiter[0])

	allRecords, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV file: %v\n", err)
		os.Exit(1)
	}

	if len(allRecords) == 0 {
		fmt.Println("Input CSV file is empty. No output generated.")
		os.Exit(0)
	}

	var columnWidths []int
	if widthsStr != "" {
		columnWidths, err = parseWidths(widthsStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing widths: %v\n", err)
			os.Exit(1)
		}
		// Validate that provided widths match the number of columns in the first record
		// This is a basic check, a more robust one might check all records for consistency.
		if len(allRecords[0]) > 0 && len(columnWidths) != len(allRecords[0]) {
			fmt.Fprintf(os.Stderr, "Error: Number of specified widths (%d) does not match number of columns in the first record (%d).\n", len(columnWidths), len(allRecords[0]))
			os.Exit(1)
		}
	} else {
		// Auto-detect widths
		maxCols := 0
		for _, record := range allRecords {
			if len(record) > maxCols {
				maxCols = len(record)
			}
		}
		columnWidths = make([]int, maxCols)

		for _, record := range allRecords {
			for i, field := range record {
				if i < len(columnWidths) {
					width := utf8.RuneCountInString(field)
					if width > columnWidths[i] {
						columnWidths[i] = width
					}
				}
			}
		}
	}

	outF, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file '%s': %v\n", outputFile, err)
		os.Exit(1)
	}
	defer outF.Close()

	startIndex := 0
	if hasHeader {
		// Write header
		if len(allRecords[0]) > len(columnWidths) {
			fmt.Fprintf(os.Stderr, "Warning: Header row has more columns (%d) than auto-detected/specified widths (%d). Extra columns will not be padded.\n", len(allRecords[0]), len(columnWidths))
		}
		var paddedHeaderFields []string
		for i, field := range allRecords[0] {
			if i < len(columnWidths) {
				paddedHeaderFields = append(paddedHeaderFields, padString(field, columnWidths[i], padLeft, padChar))
			} else {
				paddedHeaderFields = append(paddedHeaderFields, field)
			}
		}
		_, err = fmt.Fprintln(outF, strings.Join(paddedHeaderFields, ""))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing header to output file: %v\n", err)
			os.Exit(1)
		}
		startIndex = 1
	}

	// Write data rows
	for i := startIndex; i < len(allRecords); i++ {
		record := allRecords[i]
		if len(record) > len(columnWidths) {
			fmt.Fprintf(os.Stderr, "Warning: Data row %d has more columns (%d) than auto-detected/specified widths (%d). Extra columns will not be padded.\n", i+1, len(record), len(columnWidths))
		}
		var paddedFields []string
		for j, field := range record {
			if j < len(columnWidths) {
				paddedFields = append(paddedFields, padString(field, columnWidths[j], padLeft, padChar))
			} else {
				paddedFields = append(paddedFields, field)
			}
		}
		_, err = fmt.Fprintln(outF, strings.Join(paddedFields, ""))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing data row %d to output file: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully converted '%s' to fixed-width format in '%s'.\n", inputFile, outputFile)
}

// Additional implementation at 2025-06-23 01:36:19
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ColumnConfig struct {
	CSVIndex    int    // 0-based index of the column in the input CSV
	Width       int    // Desired fixed width for this column
	Alignment   string // "left" or "right"
	PaddingChar rune   // Character to use for padding
}

type Config struct {
	InputFile      string
	OutputFile     string
	SkipHeader     bool
	Columns        []ColumnConfig
	DefaultPadding rune
	DefaultAlign   string
}

func parseColumnConfigs(s string, defaultPadding rune, defaultAlign string) ([]ColumnConfig, error) {
	if s == "" {
		return nil, fmt.Errorf("column configuration string cannot be empty")
	}

	parts := strings.Split(s, ",")
	configs := make([]ColumnConfig, len(parts))

	for i, part := range parts {
		subParts := strings.Split(part, ":")
		if len(subParts) < 2 || len(subParts) > 4 {
			return nil, fmt.Errorf("invalid column definition '%s'. Expected format 'index:width[:alignment[:paddingChar]]'", part)
		}

		idx, err := strconv.Atoi(subParts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid column index '%s' in '%s': %w", subParts[0], part, err)
		}
		if idx < 0 {
			return nil, fmt.Errorf("column index cannot be negative: %d", idx)
		}

		width, err := strconv.Atoi(subParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid column width '%s' in '%s': %w", subParts[1], part, err)
		}
		if width <= 0 {
			return nil, fmt.Errorf("column width must be positive: %d", width)
		}

		configs[i].CSVIndex = idx
		configs[i].Width = width
		configs[i].Alignment = defaultAlign
		configs[i].PaddingChar = defaultPadding

		if len(subParts) >= 3 {
			align := strings.ToLower(subParts[2])
			if align != "left" && align != "right" {
				return nil, fmt.Errorf("invalid alignment '%s' in '%s'. Must be 'left' or 'right'", subParts[2], part)
			}
			configs[i].Alignment = align
		}

		if len(subParts) == 4 {
			if len(subParts[3]) != 1 {
				return nil, fmt.Errorf("padding character must be a single character in '%s'", part)
			}
			configs[i].PaddingChar = rune(subParts[3][0])
		}
	}
	return configs, nil
}

func padAndTruncate(s string, width int, align string, padChar rune) string {
	runes := []rune(s)
	currentWidth := utf8.RuneCountInString(s)

	if currentWidth == width {
		return s
	}

	if currentWidth > width {
		if width >= 3 {
			return string(runes[:width-3]) + "..."
		}
		return string(runes[:width])
	}

	padding := strings.Repeat(string(padChar), width-currentWidth)
	if align == "left" {
		return s + padding
	}
	return padding + s
}

func convertCSVToFixedWidth(cfg Config) error {
	inFile, err := os.Open(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file '%s': %w", cfg.InputFile, err)
	}
	defer inFile.Close()

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFile, err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	if cfg.SkipHeader {
		_, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("input file '%s' is empty or only contains a header", cfg.InputFile)
			}
			return fmt.Errorf("failed to read header from '%s': %w", cfg.InputFile, err)
		}
	}

	lineNum := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading CSV record at line %d: %v", lineNum+1, err)
			continue
		}
		lineNum++

		var fixedWidthLine strings.Builder
		for _, colCfg := range cfg.Columns {
			val := ""
			if colCfg.CSVIndex < len(record) {
				val = record[colCfg.CSVIndex]
			} else {
				log.Printf("Warning: Row %d, column index %d out of bounds. Using empty string.", lineNum, colCfg.CSVIndex)
			}
			formattedVal := padAndTruncate(val, colCfg.Width, colCfg.Alignment, colCfg.PaddingChar)
			fixedWidthLine.WriteString(formattedVal)
		}
		fixedWidthLine.WriteString("\n")

		_, err = writer.WriteString(fixedWidthLine.String())
		if err != nil {
			return fmt.Errorf("failed to write to output file at line %d: %w", lineNum, err)
		}
	}

	return nil
}

func main() {
	var (
		inputFile      string
		outputFile     string
		skipHeader     bool
		columnDefs     string
		defaultPadding string
		defaultAlign   string
	)

	flag.StringVar(&inputFile, "in", "", "Path to the input CSV file (required)")
	flag.StringVar(&outputFile, "out", "", "Path to the output fixed-width file (required)")
	flag.BoolVar(&skipHeader, "skip-header", false, "Skip the first row (header) of the CSV input")
	flag.StringVar(&columnDefs, "cols", "", "Comma-separated column definitions. Format: 'index:width[:alignment[:paddingChar]]'. E.g., '0:10:left: ,1:20:right:-,2:5:left:0'. (required)")
	flag.StringVar(&defaultPadding, "pad-char", " ", "Default padding character (single character). Overridden by column-specific padding.")
	flag.StringVar(&defaultAlign, "align", "left", "Default alignment ('left' or 'right'). Overridden by column-specific alignment.")

	flag.Parse()

	if inputFile == "" || outputFile == "" || columnDefs == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(defaultPadding) != 1 {
		log.Fatalf("Error: --pad-char must be a single character.")
	}
	if defaultAlign != "left" && defaultAlign != "right" {
		log.Fatalf("Error: --align must be 'left' or 'right'.")
	}

	cfg := Config{
		InputFile:      inputFile,
		OutputFile:     outputFile,
		SkipHeader:     skipHeader,
		DefaultPadding: rune(defaultPadding[0]),
		DefaultAlign:   defaultAlign,
	}

	var err error
	cfg.Columns, err = parseColumnConfigs(columnDefs, cfg.DefaultPadding, cfg.DefaultAlign)
	if err != nil {
		log.Fatalf("Error parsing column definitions: %v", err)
	}

	log.Printf("Starting conversion from '%s' to '%s'...", cfg.InputFile, cfg.OutputFile)
	if err := convertCSVToFixedWidth(cfg); err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}

	log.Println("Conversion completed successfully.")
}

// Additional implementation at 2025-06-23 01:37:25
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ColumnSpec struct {
	SourceName string
	Width      int
	Alignment  string // "left" or "right"
}

type Config struct {
	InputFile   string
	OutputFile  string
	ColumnSpecs []ColumnSpec
	SkipHeader  bool
	PaddingChar rune
}

func parseColumnSpecs(specStr string) ([]ColumnSpec, error) {
	var specs []ColumnSpec
	parts := strings.Split(specStr, ",")
	for _, part := range parts {
		subParts := strings.Split(part, ":")
		if len(subParts) < 2 || len(subParts) > 3 {
			return nil, fmt.Errorf("invalid column spec format: %s. Expected SourceName:Width[:Alignment]", part)
		}

		name := subParts[0]
		width, err := strconv.Atoi(subParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid width for column %s: %w", name, err)
		}
		if width <= 0 {
			return nil, fmt.Errorf("width for column %s must be positive", name)
		}

		alignment := "left" // Default alignment
		if len(subParts) == 3 {
			align := strings.ToLower(subParts[2])
			if align != "left" && align != "right" {
				return nil, fmt.Errorf("invalid alignment for column %s: %s. Expected 'left' or 'right'", name, subParts[2])
			}
			alignment = align
		}

		specs = append(specs, ColumnSpec{
			SourceName: name,
			Width:      width,
			Alignment:  alignment,
		})
	}
	return specs, nil
}

func padAndTruncate(s string, width int, alignment string, padChar rune) string {
	runes := []rune(s)
	currentWidth := utf8.RuneCountInString(s)

	if currentWidth > width {
		// Truncate
		return string(runes[:width])
	}

	// Pad
	padding := strings.Repeat(string(padChar), width-currentWidth)
	if alignment == "right" {
		return padding + s
	}
	// Default to left alignment
	return s + padding
}

func processCSV(cfg Config) error {
	var input io.Reader
	if cfg.InputFile == "" || cfg.InputFile == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("failed to open input file %s: %w", cfg.InputFile, err)
		}
		defer f.Close()
		input = f
	}

	var output io.Writer
	if cfg.OutputFile == "" || cfg.OutputFile == "-" {
		output = os.Stdout
	} else {
		f, err := os.Create(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", cfg.OutputFile, err)
		}
		defer f.Close()
		output = f
	}

	csvReader := csv.NewReader(input)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	header, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("input CSV is empty")
		}
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create a mapping from source column name to its index in the CSV header
	headerMap := make(map[string]int)
	for i, colName := range header {
		headerMap[colName] = i
	}

	// Validate that all specified source columns exist in the CSV header
	for _, spec := range cfg.ColumnSpecs {
		if _, ok := headerMap[spec.SourceName]; !ok {
			return fmt.Errorf("source column '%s' not found in CSV header", spec.SourceName)
		}
	}

	outputWriter := bufio.NewWriter(output)
	defer outputWriter.Flush()

	if !cfg.SkipHeader {
		var headerLine strings.Builder
		for _, spec := range cfg.ColumnSpecs {
			headerLine.WriteString(padAndTruncate(spec.SourceName, spec.Width, spec.Alignment, cfg.PaddingChar))
		}
		_, err := outputWriter.WriteString(headerLine.String() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		var line strings.Builder
		for _, spec := range cfg.ColumnSpecs {
			colIndex, ok := headerMap[spec.SourceName]
			var value string
			if ok && colIndex < len(record) {
				value = record[colIndex]
			}
			line.WriteString(padAndTruncate(value, spec.Width, spec.Alignment, cfg.PaddingChar))
		}
		_, err = outputWriter.WriteString(line.String() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func main() {
	var (
		inputFile   string
		outputFile  string
		columnSpecs string
		skipHeader  bool
		paddingChar string
	)

	// Simple command-line argument parsing for demonstration
	// A more robust solution would use the 'flag' package or similar.
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-i", "--input":
			if i+1 < len(args) {
				inputFile = args[i+1]
				i++
			}
		case "-o", "--output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-c", "--columns":
			if i+1 < len(args) {
				columnSpecs = args[i+1]
				i++
			}
		case "-s", "--skip-header":
			skipHeader = true
		case "-p", "--padding-char":
			if i+1 < len(args) {
				paddingChar = args[i+1]
				i++
			}
		case "-h", "--help":
			fmt.Println("Usage: go run main.go [OPTIONS]")
			fmt.Println("Converts CSV to fixed-width columns.")
			fmt.Println("\nOptions:")
			fmt.Println("  -i, --input <file>        Input CSV file (default: stdin)")
			fmt.Println("  -o, --output <file>       Output fixed-width file (default: stdout)")
			fmt.Println("  -c, --columns <spec>      Comma-separated column specifications.")
			fmt.Println("                            Format: 'SourceName:Width[:Alignment]'")
			fmt.Println("                            Example: 'Name:20:left,Age:5:right,City:15'")
			fmt.Println("  -s, --skip-header         Do not include header in output.")
			fmt.Println("  -p, --padding-char <char> Character to use for padding (default: space).")
			fmt.Println("  -h, --help                Show this help message.")
			os.Exit(0)
		}
	}

	if columnSpecs == "" {
		fmt.Fprintln(os.Stderr, "Error: Column specifications are required. Use -c or --columns.")
		fmt.Fprintln(os.Stderr, "Use -h or --help for usage information.")
		os.Exit(1)
	}

	parsedSpecs, err := parseColumnSpecs(columnSpecs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing column specifications: %v\n", err)
		os.Exit(1)
	}

	padRune := ' ' // Default padding character
	if paddingChar != "" {
		r, size := utf8.DecodeRuneInString(paddingChar)
		if size != len(paddingChar) || r == utf8.RuneError {
			fmt.Fprintf(os.Stderr, "Error: Padding character must be a single valid UTF-8 character.\n")
			os.Exit(1)
		}
		padRune = r
	}

	config := Config{
		InputFile:   inputFile,
		OutputFile:  outputFile,
		ColumnSpecs: parsedSpecs,
		SkipHeader:  skipHeader,
		PaddingChar: padRune,
	}

	if err := processCSV(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing CSV: %v\n", err)
		os.Exit(1)
	}
}