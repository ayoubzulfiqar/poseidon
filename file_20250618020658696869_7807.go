package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"isActive"`
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateUser(id int) User {
	return User{
		ID:       id,
		Name:     generateRandomString(rand.Intn(10) + 5),
		Email:    fmt.Sprintf("%s@example.com", generateRandomString(rand.Intn(8)+4)),
		Age:      rand.Intn(60) + 18,
		IsActive: rand.Float64() < 0.7,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numRecords := 100
	outputFile := "test_data.json"

	users := make([]User, numRecords)
	for i := 0; i < numRecords; i++ {
		users[i] = generateUser(i + 1)
	}

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
}

// Additional implementation at 2025-06-18 02:07:38
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Person struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	IsAdult   bool      `json:"isAdult"`
	Weight    float64   `json:"weight"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      []string  `json:"tags"`
}

func generateName() string {
	firstNames := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi"}
	lastNames := []string{"Smith", "Jones", "Williams", "Brown", "Davis", "Miller", "Wilson", "Moore"}
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}

func generateEmail(name string) string {
	domains := []string{"example.com", "test.org", "mail.net", "domain.io"}
	nameParts := []rune(name)
	emailPrefix := ""
	for _, r := range nameParts {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			emailPrefix += string(r)
		}
	}
	return fmt.Sprintf("%s%d@%s", emailPrefix, rand.Intn(1000), domains[rand.Intn(len(domains))])
}

func generateTags() []string {
	possibleTags := []string{"tech", "sport", "food", "travel", "art", "music", "science", "nature"}
	numTags := rand.Intn(3) + 1
	tags := make([]string, numTags)
	for i := 0; i < numTags; i++ {
		tags[i] = possibleTags[rand.Intn(len(possibleTags))]
	}
	return tags
}

func generatePerson(id int) Person {
	name := generateName()
	age := rand.Intn(80) + 1
	weight := 40.0 + rand.Float64()*80.0
	createdAt := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)

	return Person{
		ID:        id,
		Name:      name,
		Age:       age,
		IsAdult:   age >= 18,
		Weight:    weight,
		Email:     generateEmail(name),
		CreatedAt: createdAt,
		Tags:      generateTags(),
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numRecords := 10

	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil && val > 0 {
			numRecords = val
		} else {
			fmt.Fprintf(os.Stderr, "Usage: %s [number_of_records]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Invalid number of records: %s. Using default %d.\n", os.Args[1], numRecords)
		}
	}

	data := make([]Person, numRecords)
	for i := 0; i < numRecords; i++ {
		data[i] = generatePerson(i + 1)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}