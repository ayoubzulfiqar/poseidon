package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func wrapText(text string, lineWidth int) string {
	if lineWidth <= 0 {
		return text
	}
	if text == "" {
		return ""
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var builder strings.Builder
	currentLineLength := 0
	firstWordOnLine := true

	for _, word := range words {
		wordLen := len(word)

		if wordLen > lineWidth {
			if !firstWordOnLine {
				builder.WriteString("\n")
			}
			builder.WriteString(word)
			builder.WriteString("\n")
			currentLineLength = 0
			firstWordOnLine = true
			continue
		}

		if currentLineLength == 0 {
			builder.WriteString(word)
			currentLineLength = wordLen
			firstWordOnLine = false
		} else if currentLineLength+1+wordLen <= lineWidth {
			builder.WriteString(" ")
			builder.WriteString(word)
			currentLineLength += 1 + wordLen
		} else {
			builder.WriteString("\n")
			builder.WriteString(word)
			currentLineLength = wordLen
			firstWordOnLine = false
		}
	}

	return builder.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter text to wrap (press Enter twice to finish):\n")
	var inputLines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		if line == "" {
			break
		}
		inputLines = append(inputLines, line)
	}
	textToWrap := strings.Join(inputLines, " ")

	fmt.Print("Enter line width: ")
	widthStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading width:", err)
		os.Exit(1)
	}
	widthStr = strings.TrimSpace(widthStr)
	lineWidth, err := strconv.Atoi(widthStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid width. Please enter a number.", err)
		os.Exit(1)
	}

	wrappedText := wrapText(textToWrap, lineWidth)
	fmt.Println("\nWrapped Text:")
	fmt.Println(wrappedText)
}

// Additional implementation at 2025-06-19 23:53:21
package main

import (
	"fmt"
	"regexp"
	"strings"
)

// WrapText wraps the given text at word boundaries to fit within the specified line width.
// It treats blocks of text separated by one or more blank lines as separate paragraphs.
// An optional prefix string can be added to the beginning of each wrapped line.
// Words longer than the effective line width (lineWidth - len(prefix)) will be placed on their own line,
// potentially exceeding the specified lineWidth.
func WrapText(text string, lineWidth int, prefix string) string {
	if lineWidth <= 0 {
		return text // Return original text if width is non-positive
	}

	// Normalize input: trim leading/trailing whitespace from the whole text
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	// Split text into paragraphs based on one or more blank lines.
	// A blank line is a line containing only whitespace characters.
	// \R matches any Unicode newline sequence. \s* matches any whitespace.
	paragraphSeparator := regexp.MustCompile(`\R\s*\R+`)
	paragraphs := paragraphSeparator.Split(text, -1)

	var allWrappedLines []string
	for i, p := range paragraphs {
		// Trim whitespace from each paragraph content
		p = strings.TrimSpace(p)
		if p == "" {
			continue // Skip empty paragraphs that might result from splitting
		}

		// Replace internal newlines and multiple spaces with a single space within a paragraph.
		// This ensures that a paragraph like "Line1\nLine2" is treated as "Line1 Line2" for wrapping.
		internalWhitespace := regexp.MustCompile(`\s+`)
		normalizedParagraph := internalWhitespace.ReplaceAllString(p, " ")

		wrappedParagraphLines := wrapSingleParagraph(normalizedParagraph, lineWidth, prefix)
		allWrappedLines = append(allWrappedLines, wrappedParagraphLines...)

		// Add an empty line between paragraphs, but not after the last one
		if i < len(paragraphs)-1 && len(wrappedParagraphLines) > 0 {
			allWrappedLines = append(allWrappedLines, "")
		}
	}

	return strings.Join(allWrappedLines, "\n")
}

// wrapSingleParagraph wraps a single paragraph of text.
func wrapSingleParagraph(paragraph string, lineWidth int, prefix string) []string {
	var lines []string
	words := strings.Fields(paragraph) // Splits by one or more whitespace characters

	effectiveWidth := lineWidth - len(prefix)
	if effectiveWidth <= 0 {
		// If prefix is too long or lineWidth is too small, each word gets its own line with prefix.
		// These lines will exceed lineWidth.
		for _, word := range words {
			lines = append(lines, prefix+word)
		}
		return lines
	}

	currentLine := ""
	for _, word := range words {
		if currentLine == "" {
			// First word on the line
			if len(word) > effectiveWidth {
				// Word itself is longer than effective width, put it on its own line.
				// This line will exceed the effectiveWidth and thus lineWidth.
				lines = append(lines, prefix+word)
				currentLine = "" // Start fresh for the next word
			} else {
				currentLine = word
			}
		} else {
			// Not the first word, check if it fits with a space
			if len(currentLine)+1+len(word) > effectiveWidth {
				// Word doesn't fit, add currentLine to lines and start new line with word
				lines = append(lines, prefix+currentLine)
				if len(word) > effectiveWidth {
					// New word is also too long, put it on its own line.
					lines = append(lines, prefix+word)
					currentLine = "" // Start fresh
				} else {
					currentLine = word
				}
			} else {
				// Word fits, add it to currentLine
				currentLine += " " + word
			}
		}
	}

	// Add any remaining text in currentLine
	if currentLine != "" {
		lines = append(lines, prefix+currentLine)
	}

	return lines
}

func main() {
	// Example 1: Basic wrapping
	text1 := "This is a long sentence that needs to be wrapped at word boundaries. It should fit within the specified line width."
	fmt.Println("--- Example 1: Basic Wrapping (width 40) ---")
	fmt.Println(WrapText(text1, 40, ""))
	fmt.Println("\n--- Example 1: Basic Wrapping (width 20) ---")
	fmt.Println(WrapText(text1, 20, ""))

	// Example 2: Wrapping with a prefix
	text2 := "This paragraph will be wrapped with a prefix. It demonstrates how the prefix affects the available space for the text content on each line."
	fmt.Println("\n--- Example 2: Wrapping with Prefix (width 50, prefix '> ') ---")
	fmt.Println(WrapText(text2, 50, "> "))

	// Example 3: Multiple paragraphs and varying newlines
	text3 := `This is the first paragraph. It has multiple sentences and should be wrapped independently.

This is the second paragraph. It also has several sentences.

	This is a third paragraph, which might have leading/trailing whitespace that should be trimmed.
It also demonstrates how lines within a logical paragraph are joined before wrapping.


This is the fourth paragraph, separated by more newlines.`
	fmt.Println("\n--- Example 3: Multiple Paragraphs (width 60, prefix '  ') ---")
	fmt.Println(WrapText(text3, 60, "  "))

	// Example 4: Word longer than line width
	text4 := "This is a verylongwordthatwillnotfitontheline and then some more text."
	fmt.Println("\n--- Example 4: Long Word (width 20) ---")
	fmt.Println(WrapText(text4, 20, ""))

	// Example 5: Empty string and zero/negative width
	fmt.Println("\n--- Example 5: Empty String (width 30) ---")
	fmt.Println(WrapText("", 30, ""))
	fmt.Println("\n--- Example 5: Zero Width (width 0) ---")
	fmt.Println(WrapText(text1, 0, ""))
	fmt.Println("\n--- Example 5: Negative Width (width -5) ---")
	fmt.Println(WrapText(text1, -5, ""))

	// Example 6: Prefix longer than line width
	text6 := "Short text."
	fmt.Println("\n--- Example 6: Prefix Longer Than Width (width 5, prefix '----------') ---")
	fmt.Println(WrapText(text6, 5, "----------"))

	// Example 7: Text with only spaces/newlines
	text7 := "   \n\n  \t  "
	fmt.Println("\n--- Example 7: Only Whitespace (width 30) ---")
	fmt.Println(WrapText(text7, 30, ""))

	// Example 8: Paragraphs with only one word
	text8 := "Word1\n\nWord2\n\nWord3"
	fmt.Println("\n--- Example 8: Single Word Paragraphs (width 10) ---")
	fmt.Println(WrapText(text8, 10, ""))
}

// Additional implementation at 2025-06-19 23:54:53
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

// wrapText wraps the given text at word boundaries.
// It respects maxWidth and applies the given indent string to each wrapped line.
// If preserveBlanks is true, empty lines in the input are preserved in the output.
func wrapText(text string, maxWidth int, indent string, preserveBlanks bool) string {
	var wrappedLines []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// Handle blank lines
		if strings.TrimSpace(line) == "" {
			if preserveBlanks {
				wrappedLines = append(wrappedLines, "") // Add an empty string for a blank line
			}
			continue
		}

		// Split the line into words using any whitespace as a delimiter
		words := strings.FieldsFunc(line, func(r rune) bool {
			return unicode.IsSpace(r)
		})

		if len(words) == 0 {
			continue // Should not happen if TrimSpace check is effective, but good for safety
		}

		currentLineWords := []string{}
		currentLineLength := 0 // Length includes indent and words, but not trailing space

		for _, word := range words {
			wordLen := len(word)

			// Calculate length if this word were added to the current line
			// +1 for the space between words
			potentialLineLength := currentLineLength
			if len(currentLineWords) > 0 { // If not the first word on the line, account for space
				potentialLineLength += 1
			}
			potentialLineLength += wordLen

			// If the word fits on the current line (considering indent and spaces)
			if len(currentLineWords) == 0 { // First word on a new line
				currentLineWords = append(currentLineWords, word)
				currentLineLength = len(indent) + wordLen
			} else if potentialLineLength <= maxWidth {
				currentLineWords = append(currentLineWords, word)
				currentLineLength = potentialLineLength
			} else { // Word does not fit, start a new line
				wrappedLines = append(wrappedLines, indent+strings.Join(currentLineWords, " "))
				currentLineWords = []string{word}
				currentLineLength = len(indent) + wordLen
			}
		}
		// Add the last accumulated line for the current input line
		if len(currentLineWords) > 0 {
			wrappedLines = append(wrappedLines, indent+strings.Join(currentLineWords, " "))
		}
	}

	return strings.Join(wrappedLines, "\n")
}

func main() {
	var width int
	var indent string
	var inputFile string
	var preserveBlanks bool

	flag.IntVar(&width, "w", 80, "Maximum line width for wrapping")
	flag.IntVar(&width, "width", 80, "Maximum line width for wrapping (long form)")
	flag.StringVar(&indent, "i", "", "Indentation string for wrapped lines (e.g., '  ' for two spaces)")
	flag.StringVar(&indent, "indent", "", "Indentation string for wrapped lines (long form)")
	flag.StringVar(&inputFile, "f", "", "Input file to read from (default: stdin)")
	flag.StringVar(&inputFile, "file", "", "Input file to read from (long form)")
	flag.BoolVar(&preserveBlanks, "p", false, "Preserve blank lines from input")
	flag.BoolVar(&preserveBlanks, "preserve-blanks", false, "Preserve blank lines from input (long form)")

	flag.Parse()

	var reader io.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	// Use ScanLines to read input line by line, preserving original line breaks
	scanner.Split(bufio.ScanLines)

	var inputBuilder strings.Builder
	for scanner.Scan() {
		inputBuilder.WriteString(scanner.Text())
		inputBuilder.WriteString("\n") // Re-add the newline character that ScanLines removes
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// The inputBuilder will have an extra newline at the end if the input was not empty.
	// This is fine for strings.Split, as it will result in an empty string at the end of the slice,
	// which is correctly handled by the wrapText logic for blank lines.
	inputContent := inputBuilder.String()

	wrappedOutput := wrapText(inputContent, width, indent, preserveBlanks)
	fmt.Print(wrappedOutput)
}