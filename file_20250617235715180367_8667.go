package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

func generatePassword(length int, includeUppercase, includeLowercase, includeNumbers, includeSymbols bool) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	var charPool []rune
	var guaranteedChars []rune

	if includeLowercase {
		charPool = append(charPool, []rune(lowercaseChars)...)
		guaranteedChars = append(guaranteedChars, []rune(lowercaseChars)[rand.Intn(len(lowercaseChars))])
	}
	if includeUppercase {
		charPool = append(charPool, []rune(uppercaseChars)...)
		guaranteedChars = append(guaranteedChars, []rune(uppercaseChars)[rand.Intn(len(uppercaseChars))])
	}
	if includeNumbers {
		charPool = append(charPool, []rune(numberChars)...)
		guaranteedChars = append(guaranteedChars, []rune(numberChars)[rand.Intn(len(numberChars))])
	}
	if includeSymbols {
		charPool = append(charPool, []rune(symbolChars)...)
		guaranteedChars = append(guaranteedChars, []rune(symbolChars)[rand.Intn(len(symbolChars))])
	}

	if len(charPool) == 0 {
		return "", fmt.Errorf("no character types selected for password generation")
	}

	if length < len(guaranteedChars) {
		rand.Shuffle(len(guaranteedChars), func(i, j int) {
			guaranteedChars[i], guaranteedChars[j] = guaranteedChars[j], guaranteedChars[i]
		})
		return string(guaranteedChars[:length]), nil
	}

	password := make([]rune, length)

	copy(password, guaranteedChars)

	for i := len(guaranteedChars); i < length; i++ {
		password[i] = charPool[rand.Intn(len(charPool))]
	}

	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var lengthStr string
	var includeUppercaseStr, includeLowercaseStr, includeNumbersStr, includeSymbolsStr string

	fmt.Print("Enter password length: ")
	fmt.Scanln(&lengthStr)

	length, err := strconv.Atoi(strings.TrimSpace(lengthStr))
	if err != nil {
		fmt.Println("Invalid length. Please enter a number.")
		return
	}

	fmt.Print("Include uppercase letters? (y/n): ")
	fmt.Scanln(&includeUppercaseStr)
	includeUppercase := strings.ToLower(strings.TrimSpace(includeUppercaseStr)) == "y"

	fmt.Print("Include lowercase letters? (y/n): ")
	fmt.Scanln(&includeLowercaseStr)
	includeLowercase := strings.ToLower(strings.TrimSpace(includeLowercaseStr)) == "y"

	fmt.Print("Include numbers? (y/n): ")
	fmt.Scanln(&includeNumbersStr)
	includeNumbers := strings.ToLower(strings.TrimSpace(includeNumbersStr)) == "y"

	fmt.Print("Include symbols? (y/n): ")
	fmt.Scanln(&includeSymbolsStr)
	includeSymbols := strings.ToLower(strings.TrimSpace(includeSymbolsStr)) == "y"

	password, err := generatePassword(length, includeUppercase, includeLowercase, includeNumbers, includeSymbols)
	if err != nil {
		fmt.Println("Error generating password:", err)
		return
	}

	fmt.Println("Generated Password:", password)
}

// Additional implementation at 2025-06-17 23:58:03
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"flag"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>/?`~"
	ambiguousChars = "lLiIoO0"
)

type PasswordConfig struct {
	Length            int
	IncludeUppercase  bool
	IncludeLowercase  bool
	IncludeDigits     bool
	IncludeSymbols    bool
	ExcludeAmbiguous  bool
	ExcludeCharacters string
}

func getRandomIndex(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be positive")
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}
	return int(nBig.Int64()), nil
}

func filterChars(source string, exclude string) string {
	if exclude == "" {
		return source
	}
	var b strings.Builder
	for _, r := range source {
		if !strings.ContainsRune(exclude, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func GeneratePassword(config PasswordConfig) (string, error) {
	if config.Length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	var charPoolBuilder strings.Builder
	var requiredChars []rune

	if config.IncludeLowercase {
		chars := lowercaseChars
		if config.ExcludeAmbiguous {
			chars = filterChars(chars, ambiguousChars)
		}
		chars = filterChars(chars, config.ExcludeCharacters)
		if len(chars) == 0 {
			return "", fmt.Errorf("no lowercase characters available after exclusions")
		}
		charPoolBuilder.WriteString(chars)
		idx, err := getRandomIndex(len(chars))
		if err != nil { return "", err }
		requiredChars = append(requiredChars, rune(chars[idx]))
	}
	if config.IncludeUppercase {
		chars := uppercaseChars
		if config.ExcludeAmbiguous {
			chars = filterChars(chars, ambiguousChars)
		}
		chars = filterChars(chars, config.ExcludeCharacters)
		if len(chars) == 0 {
			return "", fmt.Errorf("no uppercase characters available after exclusions")
		}
		charPoolBuilder.WriteString(chars)
		idx, err := getRandomIndex(len(chars))
		if err != nil { return "", err }
		requiredChars = append(requiredChars, rune(chars[idx]))
	}
	if config.IncludeDigits {
		chars := digitChars
		if config.ExcludeAmbiguous {
			chars = filterChars(chars, ambiguousChars)
		}
		chars = filterChars(chars, config.ExcludeCharacters)
		if len(chars) == 0 {
			return "", fmt.Errorf("no digit characters available after exclusions")
		}
		charPoolBuilder.WriteString(chars)
		idx, err := getRandomIndex(len(chars))
		if err != nil { return "", err }
		requiredChars = append(requiredChars, rune(chars[idx]))
	}
	if config.IncludeSymbols {
		chars := symbolChars
		if config.ExcludeAmbiguous {
			chars = filterChars(chars, ambiguousChars)
		}
		chars = filterChars(chars, config.ExcludeCharacters)
		if len(chars) == 0 {
			return "", fmt.Errorf("no symbol characters available after exclusions")
		}
		charPoolBuilder.WriteString(chars)
		idx, err := getRandomIndex(len(chars))
		if err != nil { return "", err }
		requiredChars = append(requiredChars, rune(chars[idx]))
	}

	charPool := []rune(charPoolBuilder.String())
	if len(charPool) == 0 {
		return "", fmt.Errorf("no character types selected for password generation")
	}

	if config.Length < len(requiredChars) {
		return "", fmt.Errorf("password length (%d) is too short to include all required character types (%d)", config.Length, len(requiredChars))
	}

	passwordRunes := make([]rune, config.Length)

	for i, r := range requiredChars {
		passwordRunes[i] = r
	}

	for i := len(requiredChars); i < config.Length; i++ {
		idx, err := getRandomIndex(len(charPool))
		if err != nil {
			return "", err
		}
		passwordRunes[i] = charPool[idx]
	}

	for i := len(passwordRunes) - 1; i > 0; i-- {
		j, err := getRandomIndex(i + 1)
		if err != nil {
			return "", err
		}
		passwordRunes[i], passwordRunes[j] = passwordRunes[j], passwordRunes[i]
	}

	return string(passwordRunes), nil
}

func main() {
	length := flag.Int("length", 16, "Length of the password")
	uppercase := flag.Bool("uppercase", true, "Include uppercase characters")
	lowercase := flag.Bool("lowercase", true, "Include lowercase characters")
	digits := flag.Bool("digits", true, "Include digits")
	symbols := flag.Bool("symbols", true, "Include symbols")
	excludeAmbiguous := flag.Bool("exclude-ambiguous", false, "Exclude ambiguous characters (l, I, 1, O, 0)")
	excludeChars := flag.String("exclude-chars", "", "Characters to explicitly exclude from the password")
	count := flag.Int("count", 1, "Number of passwords to generate")

	flag.Parse()

	config := PasswordConfig{
		Length:            *length,
		IncludeUppercase:  *uppercase,
		IncludeLowercase:  *lowercase,
		IncludeDigits:     *digits,
		IncludeSymbols:    *symbols,
		ExcludeAmbiguous:  *excludeAmbiguous,
		ExcludeCharacters: *excludeChars,
	}

	for i := 0; i < *count; i++ {
		password, err := GeneratePassword(config)
		if err != nil {
			fmt.Printf("Error generating password: %v\n", err)
			return
		}
		fmt.Println(password)
	}
}

// Additional implementation at 2025-06-17 23:59:26
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>/?`~"
	ambiguousChars = "l1O0" // Characters that can be easily confused
)

// PasswordRules defines the complexity requirements for a password.
type PasswordRules struct {
	Length          int
	IncludeUppercase bool
	IncludeLowercase bool
	IncludeDigits    bool
	IncludeSymbols   bool
	NoAmbiguous      bool // Exclude characters like 'l', '1', 'O', '0'
}

// GeneratePassword generates a password based on the given rules.
func GeneratePassword(rules PasswordRules) (string, error) {
	if rules.Length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	var charSetBuilder strings.Builder
	var requiredChars []rune

	if rules.IncludeLowercase {
		charSetBuilder.WriteString(lowercaseChars)
		requiredChars = append(requiredChars, []rune(lowercaseChars)[rand.Intn(len(lowercaseChars))])
	}
	if rules.IncludeUppercase {
		charSetBuilder.WriteString(uppercaseChars)
		requiredChars = append(requiredChars, []rune(uppercaseChars)[rand.Intn(len(uppercaseChars))])
	}
	if rules.IncludeDigits {
		charSetBuilder.WriteString(digitChars)
		requiredChars = append(requiredChars, []rune(digitChars)[rand.Intn(len(digitChars))])
	}
	if rules.IncludeSymbols {
		charSetBuilder.WriteString(symbolChars)
		requiredChars = append(requiredChars, []rune(symbolChars)[rand.Intn(len(symbolChars))])
	}

	if charSetBuilder.Len() == 0 {
		return "", fmt.Errorf("at least one character type (lowercase, uppercase, digits, symbols) must be selected")
	}

	allowedChars := []rune(charSetBuilder.String())

	if rules.NoAmbiguous {
		filteredChars := make([]rune, 0, len(allowedChars))
		ambiguousMap := make(map[rune]bool)
		for _, r := range ambiguousChars {
			ambiguousMap[r] = true
		}
		for _, r := range allowedChars {
			if !ambiguousMap[r] {
				filteredChars = append(filteredChars, r)
			}
		}
		allowedChars = filteredChars
		if len(allowedChars) == 0 {
			return "", fmt.Errorf("no characters left after excluding ambiguous ones with current rules")
		}
	}

	if len(requiredChars) > rules.Length {
		return "", fmt.Errorf("required character types exceed specified password length")
	}

	password := make([]rune, rules.Length)
	// Place required characters at random positions
	rand.Shuffle(len(requiredChars), func(i, j int) {
		requiredChars[i], requiredChars[j] = requiredChars[j], requiredChars[i]
	})

	// Keep track of used positions
	usedPositions := make(map[int]bool)
	for _, char := range requiredChars {
		pos := rand.Intn(rules.Length)
		for usedPositions[pos] {
			pos = rand.Intn(rules.Length)
		}
		password[pos] = char
		usedPositions[pos] = true
	}

	// Fill the remaining positions
	for i := 0; i < rules.Length; i++ {
		if !usedPositions[i] {
			password[i] = allowedChars[rand.Intn(len(allowedChars))]
		}
	}

	return string(password), nil
}

// EstimateStrength provides a very basic entropy estimation for a password.
// Entropy (bits) = log2(charset_size ^ length)
func EstimateStrength(password string) float64 {
	if len(password) == 0 {
		return 0.0
	}

	// Determine the character set used in the password to estimate its size
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSymbol := false

	for _, r := range password {
		switch {
		case strings.ContainsRune(lowercaseChars, r):
			hasLower = true
		case strings.ContainsRune(uppercaseChars, r):
			hasUpper = true
		case strings.ContainsRune(digitChars, r):
			hasDigit = true
		case strings.ContainsRune(symbolChars, r):
			hasSymbol = true
		}
	}

	charsetSize := 0
	if hasLower {
		charsetSize += len(lowercaseChars)
	}
	if hasUpper {
		charsetSize += len(uppercaseChars)
	}
	if hasDigit {
		charsetSize += len(digitChars)
	}
	if hasSymbol {
		charsetSize += len(symbolChars)
	}

	if charsetSize == 0 {
		return 0.0
	}

	return float64(len(password)) * math.Log2(float64(charsetSize))
}

// GetStrengthDescription provides a human-readable description of password strength.
func GetStrengthDescription(entropy float64) string {
	if entropy < 28 {
		return "Very Weak (easily guessable)"
	} else if entropy < 36 {
		return "Weak (vulnerable to brute-force)"
	} else if entropy < 60 {
		return "Moderate (might be cracked by dedicated attackers)"
	} else if entropy < 80 {
		return "Strong (good for most purposes)"
	} else if entropy < 128 {
		return "Very Strong (highly resistant to brute-force)"
	} else {
		return "Excellent (extremely resistant to brute-force)"
	}
}

// parseBoolInput parses a string to a boolean, handling common affirmative/negative inputs.
func parseBoolInput(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "y" || s == "yes" || s == "true" || s == "1"
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Go Password Generator")
	fmt.Println("---------------------")

	var rules PasswordRules
	var numPasswords int = 1

	// Get password length
	for {
		fmt.Print("Enter desired password length (e.g., 12): ")
		input, _ := reader.ReadString('\n')
		length, err := strconv.Atoi(strings.TrimSpace(input))
		if err == nil && length > 0 {
			rules.Length = length
			break
		}
		fmt.Println("Invalid length. Please enter a positive number.")
	}

	// Get character type preferences
	fmt.Print("Include uppercase letters? (y/n): ")
	input, _ := reader.ReadString('\n')
	rules.IncludeUppercase = parseBoolInput(input)

	fmt.Print("Include lowercase letters? (y/n): ")
	input, _ = reader.ReadString('\n')
	rules.IncludeLowercase = parseBoolInput(input)

	fmt.Print("Include digits? (y/n): ")
	input, _ = reader.ReadString('\n')
	rules.IncludeDigits = parseBoolInput(input)

	fmt.Print("Include symbols? (y/n): ")
	input, _ = reader.ReadString('\n')
	rules.IncludeSymbols = parseBoolInput(input)

	fmt.Print("Exclude ambiguous characters (l, 1, O, 0)? (y/n): ")
	input, _ = reader.ReadString('\n')
	rules.NoAmbiguous = parseBoolInput(input)

	// Get number of passwords to generate (batch generation)
	for {
		fmt.Print("How many passwords to generate? (default: 1): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			numPasswords = 1
			break
		}
		count, err := strconv.Atoi(input)
		if err == nil && count > 0 {
			numPasswords = count
			break
		}
		fmt.Println("Invalid number. Please enter a positive number or leave blank for 1.")
	}

	fmt.Println("\nGenerated Passwords:")
	fmt.Println("--------------------")

	for i := 0; i < numPasswords; i++ {
		password, err := GeneratePassword(rules)
		if err != nil {
			fmt.Printf("Error generating password %d: %v\n", i+1, err)
			continue
		}
		entropy := EstimateStrength(password)
		strengthDesc := GetStrengthDescription(entropy)
		fmt.Printf("Password %d: %s (Entropy: %.2f bits, Strength: %s)\n", i+1, password, entropy, strengthDesc)
	}

	fmt.Println("\nGeneration complete.")
}