package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	dice              [5]int
	scoreCard         map[string]int
	categories        []string
	availableCategories map[string]bool
	reader            *bufio.Scanner
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	categories := []string{
		"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes",
		"Three of a Kind", "Four of a Kind", "Full House",
		"Small Straight", "Large Straight", "Yahtzee", "Chance",
	}
	available := make(map[string]bool)
	for _, cat := range categories {
		available[cat] = true
	}

	return &Game{
		scoreCard:         make(map[string]int),
		categories:        categories,
		availableCategories: available,
		reader:            bufio.NewScanner(os.Stdin),
	}
}

func (g *Game) rollDice(keep []int) {
	for i := 0; i < 5; i++ {
		kept := false
		for _, k := range keep {
			if i+1 == k {
				kept = true
				break
			}
		}
		if !kept {
			g.dice[i] = rand.Intn(6) + 1
		}
	}
}

func (g *Game) displayDice() {
	fmt.Print("Your dice: [ ")
	for i, d := range g.dice {
		fmt.Printf("%d", d)
		if i < len(g.dice)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println(" ]")
}

func getCounts(dice [5]int) map[int]int {
	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}
	return counts
}

func calculateCategoryScore(category string, dice [5]int) int {
	counts := getCounts(dice)
	sum := 0
	for _, d := range dice {
		sum += d
	}

	sortedDice := make([]int, 5)
	copy(sortedDice, dice[:])
	sort.Ints(sortedDice)

	switch category {
	case "Ones":
		return counts[1] * 1
	case "Twos":
		return counts[2] * 2
	case "Threes":
		return counts[3] * 3
	case "Fours":
		return counts[4] * 4
	case "Fives":
		return counts[5] * 5
	case "Sixes":
		return counts[6] * 6
	case "Three of a Kind":
		for _, count := range counts {
			if count >= 3 {
				return sum
			}
		}
		return 0
	case "Four of a Kind":
		for _, count := range counts {
			if count >= 4 {
				return sum
			}
		}
		return 0
	case "Full House":
		hasThree := false
		hasTwo := false
		for _, count := range counts {
			if count == 3 {
				hasThree = true
			} else if count == 2 {
				hasTwo = true
			}
		}
		if hasThree && hasTwo {
			return 25
		}
		return 0
	case "Small Straight":
		uniqueSorted := []int{}
		seen := make(map[int]bool)
		for _, d := range sortedDice {
			if !seen[d] {
				uniqueSorted = append(uniqueSorted, d)
				seen[d] = true
			}
		}
		s := ""
		for _, d := range uniqueSorted {
			s += strconv.Itoa(d)
		}
		if strings.Contains(s, "1234") || strings.Contains(s, "2345") || strings.Contains(s, "3456") {
			return 30
		}
		return 0
	case "Large Straight":
		if len(counts) != 5 {
			return 0
		}
		if (sortedDice[0] == 1 && sortedDice[1] == 2 && sortedDice[2] == 3 && sortedDice[3] == 4 && sortedDice[4] == 5) ||
			(sortedDice[0] == 2 && sortedDice[1] == 3 && sortedDice[2] == 4 && sortedDice[3] == 5 && sortedDice[4] == 6) {
			return 40
		}
		return 0
	case "Yahtzee":
		for _, count := range counts {
			if count == 5 {
				return 50
			}
		}
		return 0
	case "Chance":
		return sum
	default:
		return 0
	}
}

func (g *Game) displayScoreCard() {
	fmt.Println("\n--- Score Card ---")
	upperScore := 0
	lowerScore := 0
	for _, cat := range g.categories {
		score, ok := g.scoreCard[cat]
		status := " "
		if !g.availableCategories[cat] {
			status = "X"
		}
		fmt.Printf("%-20s [%s]: %d\n", cat, status, score)

		if cat == "Ones" || cat == "Twos" || cat == "Threes" || cat == "Fours" || cat == "Fives" || cat == "Sixes" {
			upperScore += score
		} else {
			lowerScore += score
		}
	}
	fmt.Printf("%-20s    : %d\n", "Upper Section Score", upperScore)
	upperBonus := 0
	if upperScore >= 63 {
		upperBonus = 35
	}
	fmt.Printf("%-20s    : %d\n", "Upper Section Bonus", upperBonus)
	fmt.Printf("%-20s    : %d\n", "Lower Section Score", lowerScore)
	fmt.Printf("%-20s    : %d\n", "TOTAL SCORE", upperScore+upperBonus+lowerScore)
	fmt.Println("------------------")
}

func (g *Game) PlayTurn() {
	fmt.Println("\n--- New Turn ---")
	g.rollDice([]int{})

	for rollNum := 1; rollNum <= 3; rollNum++ {
		g.displayDice()
		if rollNum < 3 {
			fmt.Print("Which dice to keep? (e.g., 1 3 5 to keep dice at positions 1, 3, 5. Enter nothing to re-roll all): ")
			g.reader.Scan()
			input := strings.TrimSpace(g.reader.Text())
			if input == "" {
				g.rollDice([]int{})
			} else {
				parts := strings.Fields(input)
				keepIndices := []int{}
				for _, p := range parts {
					idx, err := strconv.Atoi(p)
					if err == nil && idx >= 1 && idx <= 5 {
						keepIndices = append(keepIndices, idx)
					}
				}
				g.rollDice(keepIndices)
			}
		}
	}
	g.displayDice()

	potentialScores := make(map[string]int)
	fmt.Println("\nPotential Scores:")
	availableCount := 0
	for _, cat := range g.categories {
		if g.availableCategories[cat] {
			score := calculateCategoryScore(cat, g.dice)
			potentialScores[cat] = score
			fmt.Printf("  %s: %d\n", cat, score)
			availableCount++
		}
	}

	if availableCount == 0 {
		fmt.Println("No available categories left. Game should be over.")
		return
	}

	for {
		fmt.Print("Choose a category to score (e.g., 'Chance'): ")
		g.reader.Scan()
		choice := strings.TrimSpace(g.reader.Text())

		if _, ok := g.availableCategories[choice]; ok && g.availableCategories[choice] {
			g.scoreCard[choice] = potentialScores[choice]
			g.availableCategories[choice] = false
			fmt.Printf("Scored %d in %s.\n", potentialScores[choice], choice)
			break
		} else {
			fmt.Println("Invalid or already used category. Please choose from available categories.")
		}
	}
}

func (g *Game) calculateFinalScore() int {
	upperScore := 0
	lowerScore := 0
	for _, cat := range g.categories {
		score, ok := g.scoreCard[cat]
		if !ok {
			score = 0
		}

		if cat == "Ones" || cat == "Twos" || cat == "Threes" || cat == "Fours" || cat == "Fives" || cat == "Sixes" {
			upperScore += score
		} else {
			lowerScore += score
		}
	}

	upperBonus := 0
	if upperScore >= 63 {
		upperBonus = 35
	}

	return upperScore + upperBonus + lowerScore
}

func main() {
	game := NewGame()
	fmt.Println("Welcome to Yahtzee!")

	for i := 0; i < 13; i++ {
		if len(game.availableCategories) == 0 {
			break
		}
		game.PlayTurn()
		game.displayScoreCard()
	}

	fmt.Println("\n--- Game Over ---")
	fmt.Printf("Final Score: %d\n", game.calculateFinalScore())
	fmt.Println("Thanks for playing!")
}

// Additional implementation at 2025-06-23 02:18:34
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Dice represents the 5 dice in the game.
type Dice []int

// ScoreCardEntry represents a category's score and whether it's been used.
type ScoreCardEntry struct {
	Score int
	Used  bool
}

// ScoreCard holds all categories and their states.
type ScoreCard map[string]ScoreCardEntry

// Category names
const (
	Ones          = "Ones"
	Twos          = "Twos"
	Threes        = "Threes"
	Fours         = "Fours"
	Fives         = "Fives"
	Sixes         = "Sixes"
	ThreeOfAKind  = "Three of a Kind"
	FourOfAKind   = "Four of a Kind"
	FullHouse     = "Full House"
	SmallStraight = "Small Straight"
	LargeStraight = "Large Straight"
	Yahtzee       = "Yahtzee"
	Chance        = "Chance"
)

var allCategories = []string{
	Ones, Twos, Threes, Fours, Fives, Sixes,
	ThreeOfAKind, FourOfAKind, FullHouse, SmallStraight, LargeStraight, Yahtzee, Chance,
}

// init initializes the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// rollDice rolls the specified number of dice.
func rollDice(num int) Dice {
	dice := make(Dice, num)
	for i := 0; i < num; i++ {
		dice[i] = rand.Intn(6) + 1 // Dice values 1-6
	}
	return dice
}

// printDice displays the current dice values.
func printDice(dice Dice) {
	fmt.Print("Your dice: [ ")
	for _, d := range dice {
		fmt.Printf("%d ", d)
	}
	fmt.Println("]")
}

// getIntInput reads an integer from the user.
func getIntInput(prompt string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid input: please enter a number")
	}
	return val, nil
}

// getLineInput reads a line of string input from the user.
func getLineInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// countFrequencies counts the occurrences of each die value.
func countFrequencies(dice Dice) map[int]int {
	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}
	return counts
}

// sumDice sums all values in a Dice slice.
func sumDice(dice Dice) int {
	sum := 0
	for _, d := range dice {
		sum += d
	}
	return sum
}

// calculateCategoryScore calculates the score for a given category.
func calculateCategoryScore(dice Dice, category string) int {
	counts := countFrequencies(dice)
	sortedDice := make(Dice, len(dice))
	copy(sortedDice, dice)
	sort.Ints(sortedDice)

	switch category {
	case Ones:
		return counts[1] * 1
	case Twos:
		return counts[2] * 2
	case Threes:
		return counts[3] * 3
	case Fours:
		return counts[4] * 4
	case Fives:
		return counts[5] * 5
	case Sixes:
		return counts[6] * 6
	case ThreeOfAKind:
		for _, count := range counts {
			if count >= 3 {
				return sumDice(dice)
			}
		}
		return 0
	case FourOfAKind:
		for _, count := range counts {
			if count >= 4 {
				return sumDice(dice)
			}
		}
		return 0
	case FullHouse:
		hasTwo := false
		hasThree := false
		for _, count := range counts {
			if count == 2 {
				hasTwo = true
			}
			if count == 3 {
				hasThree = true
			}
		}
		if hasTwo && hasThree {
			return 25
		}
		return 0
	case SmallStraight:
		// Check for sequences of 4: 1-2-3-4, 2-3-4-5, 3-4-5-6
		// Use a set to remove duplicates for straight checking
		uniqueDice := make(map[int]bool)
		for _, d := range dice {
			uniqueDice[d] = true
		}
		if len(uniqueDice) < 4 {
			return 0
		}
		uniqueSorted := make([]int, 0, len(uniqueDice))
		for k := range uniqueDice {
			uniqueSorted = append(uniqueSorted, k)
		}
		sort.Ints(uniqueSorted)

		// Check for 1-2-3-4, 2-3-4-5, 3-4-5-6 within the unique sorted dice
		for i := 0; i <= len(uniqueSorted)-4; i++ {
			if uniqueSorted[i+1] == uniqueSorted[i]+1 &&
				uniqueSorted[i+2] == uniqueSorted[i]+2 &&
				uniqueSorted[i+3] == uniqueSorted[i]+3 {
				return 30
			}
		}
		return 0
	case LargeStraight:
		// Check for 1-2-3-4-5 or 2-3-4-5-6
		is1to5 := true
		is2to6 := true
		for i := 0; i < 5; i++ {
			if sortedDice[i] != i+1 {
				is1to5 = false
			}
			if sortedDice[i] != i+2 {
				is2to6 = false
			}
		}
		if is1to5 || is2to6 {
			return 40
		}
		return 0
	case Yahtzee:
		for _, count := range counts {
			if count == 5 {
				return 50
			}
		}
		return 0
	case Chance:
		return sumDice(dice)
	default:
		return 0 // Should not happen
	}
}

// displayScoreCard prints the current state of the score card.
func displayScoreCard(scoreCard ScoreCard) {
	fmt.Println("\n--- Score Card ---")
	upperScore := 0
	for _, cat := range []string{Ones, Twos, Threes, Fours, Fives, Sixes} {
		entry := scoreCard[cat]
		status := " "
		if entry.Used {
			status = "X"
			upperScore += entry.Score
		}
		fmt.Printf("%-15s: [%s] %d\n", cat, status, entry.Score)
	}
	fmt.Printf("Upper Section Total: %d\n", upperScore)
	upperBonus := 0
	if upperScore >= 63 {
		upperBonus = 35
	}
	fmt.Printf("Upper Section Bonus: %d\n", upperBonus)

	fmt.Println("--------------------")
	for _, cat := range []string{ThreeOfAKind, FourOfAKind, FullHouse, SmallStraight, LargeStraight, Yahtzee, Chance} {
		entry := scoreCard[cat]
		status := " "
		if entry.Used {
			status = "X"
		}
		fmt.Printf("%-15s: [%s] %d\n", cat, status, entry.Score)
	}
	fmt.Println("--------------------")
}

// playTurn manages a single player's turn.
func playTurn(turn int, scoreCard ScoreCard, yahtzeeBonus *int) {
	fmt.Printf("\n--- Turn %d ---\n", turn)
	displayScoreCard(scoreCard)

	currentDice := rollDice(5)
	printDice(currentDice)

	for rollNum := 1; rollNum < 3; rollNum++ {
		fmt.Printf("Roll %d of 3.\n", rollNum+1)
		fmt.Println("Enter dice to re-roll (e.g., '1 3 5' to re-roll 1st, 3rd, 5th dice). Enter '0' to keep all.")
		input := getLineInput("Dice to re-roll (space separated numbers): ")
		if input == "0" {
			break // Player wants to keep all dice
		}

		parts := strings.Fields(input)
		reRollIndices := make(map[int]bool)
		for _, p := range parts {
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 1 || idx > 5 {
				fmt.Println("Invalid input. Please enter numbers between 1 and 5.")
				continue
			}
			reRollIndices[idx-1] = true // Convert to 0-indexed
		}

		newDice := make(Dice, 0, 5)
		for i := 0; i < 5; i++ {
			if reRollIndices[i] {
				newDice = append(newDice, rand.Intn(6)+1)
			} else {
				newDice = append(newDice, currentDice[i])
			}
		}
		currentDice = newDice
		printDice(currentDice)
	}

	fmt.Println("\n--- Final Dice for this Turn ---")
	printDice(currentDice)

	// Check for Yahtzee bonus condition
	isCurrentRollYahtzee := calculateCategoryScore(currentDice, Yahtzee) == 50
	yahtzeeCategoryEntry := scoreCard[Yahtzee]

	if isCurrentRollYahtzee && yahtzeeCategoryEntry.Used && yahtzeeCategoryEntry.Score == 50 {
		fmt.Println("YAHTZEE BONUS! You get 100 points!")
		*yahtzeeBonus += 100
		fmt.Println("You must now use this Yahtzee as a Joker to score in another open category.")
		// The joker rule means it can score in any category, but if the corresponding upper section is open, it must go there.
		// For simplicity, we'll let the user pick, and it scores normally.
	}

	fmt.Println("\n--- Potential Scores ---")
	availableCategories := make([]string, 0)
	for _, cat := range allCategories {
		entry := scoreCard[cat]
		if !entry.Used {
			score := calculateCategoryScore(currentDice, cat)
			fmt.Printf("%-15s: %d\n", cat, score)
			availableCategories = append(availableCategories, cat)
		}
	}

	if len(availableCategories) == 0 {
		fmt.Println("No available categories left. This shouldn't happen in a 13-turn game.")
		return
	}

	chosenCategory := ""
	for {
		fmt