package main

import (
	"fmt"
	"strings"
)

func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var wrappedLines []string
	currentLine := ""
	currentLength := 0

	for _, word := range words {
		wordLen := len(word)

		if currentLength == 0 {
			currentLine = word
			currentLength = wordLen
		} else {
			if currentLength+1+wordLen <= maxWidth {
				currentLine += " " + word
				currentLength += 1 + wordLen
			} else {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = word
				currentLength = wordLen
			}
		}
	}

	if currentLine != "" {
		wrappedLines = append(wrappedLines, currentLine)
	}

	return strings.Join(wrappedLines, "\n")
}

func main() {
	sampleText := "This is a sample text that needs to be wrapped at word boundaries. It should demonstrate how the wrapping function works with various lengths of words and lines."
	maxWidth := 30

	fmt.Println("Original Text:")
	fmt.Println(sampleText)
	fmt.Printf("\nWrapped Text (Max Width: %d):\n", maxWidth)
	wrapped := wrapText(sampleText, maxWidth)
	fmt.Println(wrapped)

	fmt.Println("\n--- Another Example (Smaller Width) ---")
	sampleText2 := "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software."
	maxWidth2 := 20
	fmt.Println("Original Text:")
	fmt.Println(sampleText2)
	fmt.Printf("\nWrapped Text (Max Width: %d):\n", maxWidth2)
	wrapped2 := wrapText(sampleText2, maxWidth2)
	fmt.Println(wrapped2)

	fmt.Println("\n--- Example with long word ---")
	sampleText3 := "This is a supercalifragilisticexpialidocious word that should be on its own line if it's too long."
	maxWidth3 := 15
	fmt.Println("Original Text:")
	fmt.Println(sampleText3)
	fmt.Printf("\nWrapped Text (Max Width: %d):\n", maxWidth3)
	wrapped3 := wrapText(sampleText3, maxWidth3)
	fmt.Println(wrapped3)

	fmt.Println("\n--- Example with empty text ---")
	sampleText4 := ""
	maxWidth4 := 10
	fmt.Println("Original Text:")
	fmt.Println(sampleText4)
	fmt.Printf("\nWrapped Text (Max Width: %d):\n", maxWidth4)
	wrapped4 := wrapText(sampleText4, maxWidth4)
	fmt.Println(wrapped4)

	fmt.Println("\n--- Example with zero width ---")
	sampleText5 := "Hello World"
	maxWidth5 := 0
	fmt.Println("Original Text:")
	fmt.Println(sampleText5)
	fmt.Printf("\nWrapped Text (Max Width: %d):\n", maxWidth5)
	wrapped5 := wrapText(sampleText5, maxWidth5)
	fmt.Println(wrapped5)
}

// Additional implementation at 2025-06-17 23:27:57
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// wrapText wraps the given text at word boundaries to fit within the specified width.
// It handles existing newlines as paragraph breaks and applies indentation.
// If trimSpaces is true, multiple internal spaces are collapsed to single spaces
// and leading/trailing spaces of paragraphs are removed before wrapping.
func wrapText(text string, width int, indent string, trimSpaces bool) string {
	if width <= 0 {
		return text // No wrapping if width is invalid
	}

	var wrappedLines []string
	paragraphs := strings.Split(text, "\n")

	for _, para := range paragraphs {
		if trimSpaces {
			// Collapse multiple spaces into single spaces and remove leading/trailing
			fields := strings.Fields(para)
			para = strings.Join(fields, " ")
		}

		if len(para) == 0 {
			wrappedLines = append(wrappedLines, "") // Preserve empty lines as paragraph breaks
			continue
		}

		// Split the paragraph into words using any Unicode space as a delimiter.
		// This naturally handles multiple spaces between words by treating them as a single delimiter.
		words := strings.FieldsFunc(para, func(r rune) bool {
			return unicode.IsSpace(r)
		})

		if len(words) == 0 {
			// If the paragraph was just spaces or became empty after trimming
			wrappedLines = append(wrappedLines, indent)
			continue
		}

		currentLine := new(strings.Builder)
		currentLine.WriteString(indent)
		currentLineLen := len(indent)

		for i, word := range words {
			wordLen := len(word)

			// Handle words that are longer than the available line width (after indent).
			// Such words will occupy a full line by themselves, or be placed on a new line.
			if wordLen > width-len(indent) {
				if currentLineLen > len(indent) { // If current line has content, flush it first
					wrappedLines = append(wrappedLines, currentLine.String())
					currentLine.Reset()
					currentLine.WriteString(indent)
					currentLineLen = len(indent)
				}
				// Add the long word on its own line (or break it if it's the first word and too long)
				wrappedLines = append(wrappedLines, indent+word)
				currentLine.Reset() // Start a new line after the long word
				currentLine.WriteString(indent)
				currentLineLen = len(indent)
				continue
			}

			// Check if adding the next word (plus a space) exceeds the width.
			// If currentLineLen is just len(indent), it means it's the start of a new line,
			// so we don't add a space before the first word.
			if currentLineLen+wordLen+1 > width && currentLineLen > len(indent) {
				wrappedLines = append(wrappedLines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(indent)
				currentLineLen = len(indent)
			}

			// Add a space if it's not the very beginning of a line (after indent).
			if currentLineLen > len(indent) {
				currentLine.WriteRune(' ')
				currentLineLen++
			}
			currentLine.WriteString(word)
			currentLineLen += wordLen
		}
		// Add the last line of the paragraph if it has content
		if currentLine.Len() > len(indent) {
			wrappedLines = append(wrappedLines, currentLine.String())
		} else if len(para) > 0 { // If the paragraph was not empty but only contained spaces or was just the indent
			wrappedLines = append(wrappedLines, indent)
		}
	}

	return strings.Join(wrappedLines, "\n")
}

func main() {
	var width int
	var inputFile string
	var outputFile string
	var indentStr string
	var trimSpaces bool

	flag.IntVar(&width, "width", 80, "Maximum line width for wrapping")
	flag.IntVar(&width, "w", 80, "Maximum line width for wrapping (shorthand)")
	flag.StringVar(&inputFile, "input", "", "Input file path (default: stdin)")
	flag.StringVar(&inputFile, "i", "", "Input file path (shorthand)")
	flag.StringVar(&outputFile, "output", "", "Output file path (default: stdout)")
	flag.StringVar(&outputFile, "o", "", "Output file path (shorthand)")
	flag.StringVar(&indentStr, "indent", "", "String to prepend to each wrapped line (e.g., '    ' or '\\t')")
	flag.StringVar(&indentStr, "t", "", "String to prepend to each wrapped line (shorthand)")
	flag.BoolVar(&trimSpaces, "trim-spaces", false, "Collapse multiple spaces into single spaces before wrapping")
	flag.BoolVar(&trimSpaces, "s", false, "Collapse multiple spaces into single spaces before wrapping (shorthand)")

	flag.Parse()

	var inputReader io.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		inputReader = file
	} else {
		inputReader = os.Stdin
	}

	inputBytes, err := io.ReadAll(inputReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	inputText := string(inputBytes)

	wrappedText := wrapText(inputText, width, indentStr, trimSpaces)

	var outputWriter io.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		outputWriter = file
	} else {
		outputWriter = os.Stdout
	}

	writer := bufio.NewWriter(outputWriter)
	_, err = writer.WriteString(wrappedText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
	writer.Flush()
}

// Additional implementation at 2025-06-17 23:28:45
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// wrapText wraps the given text at word boundaries to fit within the specified width.
// It also applies an optional indent and prefix to each wrapped line.
// Existing newlines in the input text are preserved as paragraph breaks.
func wrapText(text string, width int, indent string, prefix string) string {
	if width <= 0 {
		return text // Cannot wrap to non-positive width, return original text
	}

	var wrappedLines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines) // Process line by line to respect existing newlines

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			wrappedLines = append(wrappedLines, "") // Preserve empty lines
			continue
		}

		words := strings.Fields(line) // Split by whitespace
		if len(words) == 0 {
			continue // Skip lines that are only whitespace after splitting
		}

		currentLineBuilder := new(strings.Builder)
		currentLineLength := 0

		// Add initial indent and prefix for the very first line of a paragraph
		currentLineBuilder.WriteString(indent)
		currentLineBuilder.WriteString(prefix)
		currentLineLength += len(indent) + len(prefix)

		for i, word := range words {
			wordLen := len(word)

			// Determine if a space is needed before the current word.
			// A space is needed if this is not the very first word on the logical line
			// (i.e., currentLineLength is greater than just the indent+prefix).
			needsSpace := currentLineLength > len(indent)+len(prefix)

			// Calculate potential length if word is added to current line
			potentialLength := currentLineLength
			if needsSpace {
				potentialLength++ // for the space
			}
			potentialLength += wordLen

			// Check if adding the word (plus a space if needed) would exceed the width.
			// If it's the first word on a new line and it's longer than the effective width,
			// it will still be placed on that line, potentially exceeding the width.
			// This is standard behavior for "word boundary" wrapping.
			if potentialLength > width && currentLineLength > len(indent)+len(prefix) {
				// Word doesn't fit, and there's content on the current line, so flush it.
				wrappedLines = append(wrappedLines, currentLineBuilder.String())
				currentLineBuilder.Reset()
				currentLineBuilder.WriteString(indent)
				currentLineBuilder.WriteString(prefix)
				currentLineLength = len(indent) + len(prefix)
				needsSpace = false // No space needed at the start of a new line
			}

			if needsSpace {
				currentLineBuilder.WriteByte(' ')
				currentLineLength++
			}
			currentLineBuilder.WriteString(word)
			currentLineLength += wordLen

			// If it's the last word in the paragraph, append the current line
			if i == len(words)-1 {
				wrappedLines = append(wrappedLines, currentLineBuilder.String())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		// In a real script, you might want to log this or handle it more gracefully.
		// For a simple utility, exiting is acceptable.
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	return strings.Join(wrappedLines, "\n")
}

func main() {
	width := flag.Int("w", 80, "Maximum line width for wrapping")
	indentStr := flag.String("i", "", "String to prepend as indentation to each wrapped line")
	prefixStr := flag.String("p", "", "String to prepend as a prefix after indentation to each wrapped line")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Wraps text read from stdin at word boundaries to a specified width.")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Read all input from stdin
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
		os.Exit(1)
	}
	inputText := string(inputBytes)

	// Wrap the text
	wrappedText := wrapText(inputText, *width, *indentStr, *prefixStr)

	// Print the wrapped text to stdout
	fmt.Print(wrappedText)
}

// Additional implementation at 2025-06-17 23:29:34
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func wrapParagraph(paragraph string, width int) string {
	var wrapped strings.Builder
	words := strings.Fields(paragraph)

	if len(words) == 0 {
		return ""
	}

	currentLineLength := 0
	firstWordOnLine := true

	for _, word := range words {
		wordLen := len(word)

		if firstWordOnLine {
			wrapped.WriteString(word)
			currentLineLength = wordLen
			firstWordOnLine = false
		} else if currentLineLength+1+wordLen <= width { // +1 for the space
			wrapped.WriteString(" ")
			wrapped.WriteString(word)
			currentLineLength += 1 + wordLen
		} else {
			wrapped.WriteString("\n")
			wrapped.WriteString(word)
			currentLineLength = wordLen
		}
	}
	return wrapped.String()
}

func main() {
	widthPtr := flag.Int("width", 80, "Line width for wrapping text")
	flag.Parse()

	var paragraphs []string
	var currentParagraph strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" { // Blank line indicates a paragraph break
			if currentParagraph.Len() > 0 {
				paragraphs = append(paragraphs, currentParagraph.String())
				currentParagraph.Reset()
			}
			continue // Skip multiple consecutive blank lines
		}

		if currentParagraph.Len() > 0 {
			currentParagraph.WriteString(" ") // Join lines within a paragraph with a space
		}
		currentParagraph.WriteString(trimmedLine)
	}

	// Add the last paragraph if any content was accumulated
	if currentParagraph.Len() > 0 {
		paragraphs = append(paragraphs, currentParagraph.String())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	var wrappedOutput strings.Builder
	for i, p := range paragraphs {
		if i > 0 {
			wrappedOutput.WriteString("\n\n") // Separate wrapped paragraphs with double newline
		}
		wrappedOutput.WriteString(wrapParagraph(p, *widthPtr))
	}

	fmt.Print(wrappedOutput.String())
}