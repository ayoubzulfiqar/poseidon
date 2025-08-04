package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
	"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
	"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla",
	"pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident",
	"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id",
	"est", "laborum",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomWord() string {
	return loremWords[rand.Intn(len(loremWords))]
}

func generateSentence(minWords, maxWords int) string {
	numWords := rand.Intn(maxWords-minWords+1) + minWords
	words := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		words[i] = getRandomWord()
	}
	sentence := strings.Join(words, " ")
	if len(sentence) > 0 {
		sentence = strings.ToUpper(string(sentence[0])) + sentence[1:]
	}
	return sentence + "."
}

func generateParagraph(minSentences, maxSentences, minWordsPerSentence, maxWordsPerSentence int) string {
	numSentences := rand.Intn(maxSentences-minSentences+1) + minSentences
	sentences := make([]string, numSentences)
	for i := 0; i < numSentences; i++ {
		sentences[i] = generateSentence(minWordsPerSentence, maxWordsPerSentence)
	}
	return strings.Join(sentences, " ")
}

func generateLoremIpsum(numParagraphs, minSentencesPerParagraph, maxSentencesPerParagraph, minWordsPerSentence, maxWordsPerSentence int) string {
	paragraphs := make([]string, numParagraphs)
	for i := 0; i < numParagraphs; i++ {
		paragraphs[i] = generateParagraph(minSentencesPerParagraph, maxSentencesPerParagraph, minWordsPerSentence, maxWordsPerSentence)
	}
	return strings.Join(paragraphs, "\n\n")
}

func main() {
	loremText := generateLoremIpsum(3, 3, 7, 5, 15)
	fmt.Println(loremText)
}

// Additional implementation at 2025-08-04 06:32:15
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// LoremIpsumGenerator generates random Lorem Ipsum text.
type LoremIpsumGenerator struct {
	source *rand.Rand
	words  []string
	isFirstSentence bool // Flag to ensure the first sentence starts with "Lorem ipsum..."
}

// NewLoremIpsumGenerator creates and returns a new LoremIpsumGenerator.
func NewLoremIpsumGenerator() *LoremIpsumGenerator {
	// Standard Lorem Ipsum words, including some less common ones for variety.
	// The order is not important as they are picked randomly.
	words := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
		"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
		"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
		"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
		"consequat", "duis", "aute", "irure", "reprehenderit", "in", "voluptate",
		"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
		"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
		"deserunt", "mollit", "anim", "id", "est", "laborum", "curabitur", "vestibulum",
		"sagittis", "felis", "vitae", "varius", "mauris", "eget", "dictum", "nunc",
		"auctor", "libero", "vel", "risus", "finibus", "porta", "suspendisse", "potenti",
		"integer", "nec", "odio", "praesent", "ultricies", "lectus", "quisque", "bibendum",
		"elementum", "imperdiet", "donec", "pharetra", "vulputate", "dignissim", "fusce",
		"euismod", "lacinia", "orci", "neque", "facilisis", "ultrices", "fringilla",
		"cursus", "malesuada", "tellus", "luctus", "vivamus", "accumsan", "placerat",
		"urna", "aenean", "pulvinar", "rutrum", "metus", "aliquam", "erat", "volutpat",
		"cras", "dapibus", "purus", "eu", "semper", "iaculis", "justo", "rhoncus",
		"nullam", "facilisi", "pellentesque", "habitant", "morbi", "tristique", "senectus",
		"netus", "fames", "turpis", "egestas", "class", "aptent", "taciti", "sociosqu",
		"ad", "litora", "torquent", "per", "conubia", "nostra", "inceptos", "himenaeos",
		"donec", "ac", "nibh", "congue", "massa", "laoreet", "scelerisque", "etiam",
		"porttitor", "nunc", "sed", "tincidunt", "venenatis", "phasellus", "gravida",
		"lacinia", "turpis", "ut", "fermentum", "ante", "velit", "auctor", "leo",
		"quis", "molestie", "nunc", "lacinia", "quis", "elit", "a", "maximus", "ligula",
		"vitae", "finibus", "nisl", "sed", "efficitur", "nunc", "vitae", "varius",
		"aliquam", "erat", "volutpat", "maecenas", "mattis", "vel", "nisl", "eget",
		"tincidunt", "curabitur", "at", "nisl", "eu", "nisl", "ultrices", "tincidunt",
		"sed", "id", "nisl", "vel", "nisl", "ultrices", "tincidunt", "sed", "id",
		"nisl", "vel", "nisl", "ultrices", "tincidunt", "sed", "id", "nisl",
	}

	return &LoremIpsumGenerator{
		source: rand.New(rand.NewSource(time.Now().UnixNano())),
		words:  words,
		isFirstSentence: true,
	}
}

// randomInt generates a random integer between min (inclusive) and max (inclusive).
func (lig *LoremIpsumGenerator) randomInt(min, max int) int {
	return lig.source.Intn(max-min+1) + min
}

// capitalize capitalizes the first letter of a string.
func (lig *LoremIpsumGenerator) capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// generateWord returns a single random Lorem Ipsum word.
func (lig *LoremIpsumGenerator) generateWord() string {
	return lig.words[lig.source.Intn(len(lig.words))]
}

// generateSentence generates a single sentence with a random number of words.
// minWords and maxWords define the inclusive range for the sentence length.
func (lig *LoremIpsumGenerator) generateSentence(minWords, maxWords int) string {
	var sb strings.Builder
	numWords := lig.randomInt(minWords, maxWords)

	if lig.isFirstSentence {
		sb.WriteString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
		lig.isFirstSentence = false
		return sb.String()
	}

	for i := 0; i < numWords; i++ {
		word := lig.generateWord()
		if i == 0 {
			sb.WriteString(lig.capitalize(word))
		} else {
			// Randomly add a comma
			if lig.source.Intn(10) == 0 && i > 2 && i < numWords-1 { // 10% chance, not too early or late
				sb.WriteString(", ")
				sb.WriteString(word)
			} else {
				sb.WriteString(" ")
				sb.WriteString(word)
			}
		}
	}
	sb.WriteString(".")
	return sb.String()
}

// generateParagraph generates a single paragraph with a random number of sentences.
// minSentences and maxSentences define the inclusive range for the number of sentences.
// minWordsPerSentence and maxWordsPerSentence define the inclusive range for sentence length.
func (lig *LoremIpsumGenerator) generateParagraph(minSentences, maxSentences, minWordsPerSentence, maxWordsPerSentence int) string {
	var sb strings.Builder
	numSentences := lig.randomInt(minSentences, maxSentences)

	for i := 0; i < numSentences; i++ {
		sentence := lig.generateSentence(minWordsPerSentence, maxWordsPerSentence)
		sb.WriteString(sentence)
		if i < numSentences-1 {
			sb.WriteString(" ") // Space between sentences
		}
	}
	return sb.String()
}

// Generate generates a specified number of Lorem Ipsum paragraphs.
// It returns the generated text as a single string.
func (lig *LoremIpsumGenerator) Generate(numParagraphs, minWordsPerSentence, maxWordsPerSentence, minSentencesPerParagraph, maxSentencesPerParagraph int) string {
	var sb strings.Builder
	lig.isFirstSentence = true // Reset for new generation call

	for i := 0; i < numParagraphs; i++ {
		paragraph := lig.generateParagraph(minSentencesPerParagraph, maxSentencesPerParagraph, minWordsPerSentence, maxWordsPerSentence)
		sb.WriteString(paragraph)
		if i < numParagraphs-1 {
			sb.WriteString("\n\n") // Double newline between paragraphs
		}
	}
	return sb.String()
}

// GenerateHTML generates a specified number of Lorem Ipsum paragraphs wrapped in HTML <p> tags.
// It returns the generated HTML string.
func (lig *LoremIpsumGenerator) GenerateHTML(numParagraphs, minWordsPerSentence, maxWordsPerSentence, minSentencesPerParagraph, maxSentencesPerParagraph int) string {
	var sb strings.Builder
	lig.isFirstSentence = true // Reset for new generation call

	for i := 0; i < numParagraphs; i++ {
		paragraph := lig.generateParagraph(minSentencesPerParagraph, maxSentencesPerParagraph, minWordsPerSentence, maxWordsPerSentence)
		sb.WriteString("<p>")
		sb.WriteString(paragraph)
		sb.WriteString("</p>")
		if i < numParagraphs-1 {
			sb.WriteString("\n") // Newline between <p> tags for readability
		}
	}
	return sb.String()
}

func main() {
	generator := NewLoremIpsumGenerator()

	// Example 1: Generate 3 paragraphs of plain text
	fmt.Println("--- Plain Text (3 Paragraphs) ---")
	plainText := generator.Generate(
		3,    // Number of paragraphs
		5, 15, // Min/Max words per sentence
		3, 7, // Min/Max sentences per paragraph
	)
	fmt.Println(plainText)
	fmt.Println("\n-----------------------------------\n")

	// Example 2: Generate 2 paragraphs of HTML text
	fmt.Println("--- HTML Text (2 Paragraphs) ---")
	htmlText := generator.GenerateHTML(
		2,    // Number of paragraphs
		8, 20, // Min/Max words per sentence
		4, 8, // Min/Max sentences per paragraph
	)
	fmt.Println(htmlText)
	fmt.Println("\n-----------------------------------\n")

	// Example 3: Generate a short single paragraph
	fmt.Println("--- Short Single Paragraph ---")
	shortParagraph := generator.Generate(
		1,    // Number of paragraphs
		3, 8, // Min/Max words per sentence
		1, 2, // Min/Max sentences per paragraph
	)
	fmt.Println(shortParagraph)
	fmt.Println("\n-----------------------------------\n")

	// Example 4: Generate a long single paragraph
	fmt.Println("--- Long Single Paragraph ---")
	longParagraph := generator.Generate(
		1,     // Number of paragraphs
		10, 30, // Min/Max words per sentence
		8, 12, // Min/Max sentences per paragraph
	)
	fmt.Println(longParagraph)
	fmt.Println("\n-----------------------------------\n")
}

// Additional implementation at 2025-08-04 06:33:24
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "sed", "do",
	"eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua", "ut",
	"enim", "ad", "minim", "veniam", "quis", "nostrud", "exercitation", "ullamco", "laboris",
	"nisi", "ut", "aliquip", "ex", "ea", "commodo", "consequat", "duis", "aute", "irure",
	"dolor", "in", "reprehenderit", "in", "voluptate", "velit", "esse", "cillum", "dolore",
	"eu", "fugiat", "nulla", "pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non",
	"proident", "sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id",
	"est", "laborum",
}

var punctuations = []string{".", "?", "!"}

type LoremGenerator struct {
	randSource *rand.Rand
}

func NewLoremGenerator() *LoremGenerator {
	return &LoremGenerator{
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (lg *LoremGenerator) getRandomWord() string {
	return loremWords[lg.randSource.Intn(len(loremWords))]
}

func (lg *LoremGenerator) generateSentence(minWords, maxWords int) string {
	if minWords <= 0 {
		minWords = 1
	}
	if maxWords < minWords {
		maxWords = minWords
	}

	numWords := lg.randSource.Intn(maxWords-minWords+1) + minWords
	words := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		words[i] = lg.getRandomWord()
	}

	if len(words) > 0 {
		words[0] = strings.ToUpper(string(words[0][0])) + words[0][1:]
	}

	sentence := strings.Join(words, " ")
	sentence += punctuations[lg.randSource.Intn(len(punctuations))]
	return sentence
}

func (lg *LoremGenerator) generateParagraph(minSentences, maxSentences, minWords, maxWords int) string {
	if minSentences <= 0 {
		minSentences = 1
	}
	if maxSentences < minSentences {
		maxSentences = minSentences
	}

	numSentences := lg.randSource.Intn(maxSentences-minSentences+1) + minSentences
	sentences := make([]string, numSentences)

	if lg.randSource.Intn(100) < 30 { // Approximately 30% chance to start with the classic phrase
		if numSentences > 0 {
			sentences[0] = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
		}
		for i := 1; i < numSentences; i++ {
			sentences[i] = lg.generateSentence(minWords, maxWords)
		}
	} else {
		for i := 0; i < numSentences; i++ {
			sentences[i] = lg.generateSentence(minWords, maxWords)
		}
	}

	return strings.Join(sentences, " ")
}

func (lg *LoremGenerator) GenerateLoremIpsum(paragraphs, minSentences, maxSentences, minWords, maxWords int) string {
	if paragraphs <= 0 {
		return ""
	}

	result := make([]string, paragraphs)
	for i := 0; i < paragraphs; i++ {
		result[i] = lg.generateParagraph(minSentences, maxSentences, minWords, maxWords)
	}
	return strings.Join(result, "\n\n")
}

func main() {
	generator := NewLoremGenerator()

	// Example 1: Generate 3 paragraphs with varying sentence and word counts
	fmt.Println("--- Example 1: 3 Paragraphs ---")
	loremText1 := generator.GenerateLoremIpsum(3, 3, 7, 5, 15)
	fmt.Println(loremText1)

	// Example 2: Generate 1 paragraph with specific sentence and word counts
	fmt.Println("\n--- Example 2: 1 Paragraph ---")
	loremText2 := generator.GenerateLoremIpsum(1, 2, 4, 8, 20)
	fmt.Println(loremText2)

	// Example 3: Generate 5 paragraphs with different ranges
	fmt.Println("\n--- Example 3: 5 Paragraphs ---")
	loremText3 := generator.GenerateLoremIpsum(5, 4, 8, 6, 12)
	fmt.Println(loremText3)
}