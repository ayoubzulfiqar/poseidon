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

func rollDice(dice []int, keep []bool) {
	for i := range dice {
		if !keep[i] {
			dice[i] = rand.Intn(6) + 1
		}
	}
}

func printDice(dice []int) {
	fmt.Print("Your dice: [")
	for i, d := range dice {
		fmt.Printf("%d", d)
		if i < len(dice)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")
}

func chooseDiceToKeep(reader *bufio.Reader, dice []int) ([]bool, error) {
	keep := make([]bool, len(dice))
	fmt.Print("Enter indices of dice to keep (e.g., '1 3 5' for dice at index 0, 2, 4). Enter nothing to re-roll all: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return keep, nil
	}

	parts := strings.Fields(input)
	for _, p := range parts {
		idx, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid input '%s': not a number", p)
		}
		if idx < 1 || idx > len(dice) {
			return nil, fmt.Errorf("invalid index %d: must be between 1 and %d", idx, len(dice))
		}
		keep[idx-1] = true
	}
	return keep, nil
}

func calculateScore(dice []int) (string, int) {
	sort.Ints(dice)

	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}

	for _, count := range counts {
		if count == 5 {
			return "Yahtzee", 50
		}
	}

	hasThree := false
	hasTwo := false
	sumOfDice := 0
	for val, count := range counts {
		sumOfDice += val * count
		if count == 4 {
			return "Four of a Kind", sumOfDice
		}
		if count == 3 {
			hasThree = true
		}
		if count == 2 {
			hasTwo = true
		}
	}

	if hasThree && hasTwo {
		return "Full House", 25
	}
	if hasThree {
		return "Three of a Kind", sumOfDice
	}

	uniqueDice := make(map[int]bool)
	for _, d := range dice {
		uniqueDice[d] = true
	}
	uniqueSorted := make([]int, 0, len(uniqueDice))
	for k := range uniqueDice {
		uniqueSorted = append(uniqueSorted, k)
	}
	sort.Ints(uniqueSorted)

	if len(uniqueSorted) >= 5 {
		if (uniqueSorted[0] == 1 && uniqueSorted[1] == 2 && uniqueSorted[2] == 3 && uniqueSorted[3] == 4 && uniqueSorted[4] == 5) ||
			(uniqueSorted[0] == 2 && uniqueSorted[1] == 3 && uniqueSorted[2] == 4 && uniqueSorted[3] == 5 && uniqueSorted[4] == 6) {
			return "Large Straight", 40
		}
	}

	if len(uniqueSorted) >= 4 {
		for i := 0; i <= len(uniqueSorted)-4; i++ {
			if uniqueSorted[i+1] == uniqueSorted[i]+1 &&
				uniqueSorted[i+2] == uniqueSorted[i+1]+1 &&
				uniqueSorted[i+3] == uniqueSorted[i+2]+1 {
				return "Small Straight", 30
			}
		}
	}

	return "Chance", sumOfDice
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Yahtzee Dice Simulator!")
	fmt.Println("You have 3 rolls to get the best hand.")

	dice := make([]int, 5)
	keep := make([]bool, 5)

	for rollNum := 1; rollNum <= 3; rollNum++ {
		fmt.Printf("\n--- Roll %d ---\n", rollNum)
		rollDice(dice, keep)
		printDice(dice)

		if rollNum < 3 {
			for {
				selectedKeep, err := chooseDiceToKeep(reader, dice)
				if err != nil {
					fmt.Printf("Error: %s. Please try again.\n", err)
					continue
				}
				keep = selectedKeep
				break
			}
		}
	}

	fmt.Println("\n--- Final Hand ---")
	printDice(dice)
	category, score := calculateScore(dice)
	fmt.Printf("Your best hand is: %s (Score: %d)\n", category, score)
	fmt.Println("Thanks for playing!")
}

// Additional implementation at 2025-06-21 02:51:16
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	numDice  = 5
	maxRolls = 3
)

// rollDie generates a random value for a single die (1-6).
func rollDie() int {
	return rand.Intn(6) + 1
}

// rollAllDice rolls all 5 dice initially.
func rollAllDice() []int {
	dice := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		dice[i] = rollDie()
	}
	return dice
}

// displayDice prints the current state of the dice.
func displayDice(dice []int) {
	fmt.Println("Current Dice:")
	for i, d := range dice {
		fmt.Printf("[%d] %d ", i+1, d)
	}
	fmt.Println()
}

// chooseDiceToReroll prompts the user to select which dice to reroll.
// Returns a slice of indices (0-based) of dice to reroll.
func chooseDiceToReroll() []int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter dice numbers to reroll (e.g., 1 3 5), or press Enter to keep all: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return []int{} // No dice to reroll
	}

	parts := strings.Fields(input)
	var toReroll []int
	for _, p := range parts {
		num, err := strconv.Atoi(p)
		if err == nil && num >= 1 && num <= numDice {
			toReroll = append(toReroll, num-1) // Convert to 0-based index
		} else {
			fmt.Printf("Invalid input: %s. Please enter numbers between 1 and %d.\n", p, numDice)
			return chooseDiceToReroll() // Re-prompt on invalid input
		}
	}
	return toReroll
}

// rerollSelectedDice rerolls the dice at the specified indices.
func rerollSelectedDice(dice []int, indices []int) {
	for _, idx := range indices {
		if idx >= 0 && idx < numDice {
			dice[idx] = rollDie()
		}
	}
}

// simulateYahtzeeTurn simulates a single turn of Yahtzee.
func simulateYahtzeeTurn() {
	fmt.Println("--- Starting Yahtzee Turn ---")
	dice := rollAllDice()
	displayDice(dice)

	for rollNum := 1; rollNum < maxRolls; rollNum++ {
		fmt.Printf("Roll %d of %d\n", rollNum+1, maxRolls)
		indicesToReroll := chooseDiceToReroll()
		if len(indicesToReroll) == 0 {
			fmt.Println("Keeping all dice.")
			break // Player chose to keep all dice, end reroll phase
		}
		rerollSelectedDice(dice, indicesToReroll)
		displayDice(dice)
	}

	fmt.Println("--- Final Dice for this Turn ---")
	displayDice(dice)
	fmt.Println("Time to score!")
	// In a full game, scoring logic would go here based on the final dice.
	// For this simulator, we just show the final dice.
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	simulateYahtzeeTurn()
}

// Additional implementation at 2025-06-21 02:52:29
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

func rollDice() []int {
	dice := make([]int, 5)
	for i := 0; i < 5; i++ {
		dice[i] = rand.Intn(6) + 1
	}
	return dice
}

func rerollDice(currentDice []int, indicesToReroll []int) []int {
	newDice := make([]int, len(currentDice))
	copy(newDice, currentDice)

	for _, index := range indicesToReroll {
		if index >= 1 && index <= 5 {
			newDice[index-1] = rand.Intn(6) + 1
		}
	}
	return newDice
}

func countFrequencies(dice []int) map[int]int {
	counts := make(map[int]int)
	for _, die := range dice {
		counts[die]++
	}
	return counts
}

func calculateScore(dice []int, category string) int {
	counts := countFrequencies(dice)
	sumAllDice := 0
	for _, die := range dice {
		sumAllDice += die
	}

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
				return sumAllDice
			}
		}
		return 0
	case "Four of a Kind":
		for _, count := range counts {
			if count >= 4 {
				return sumAllDice
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
		sortedUniqueDice := getSortedUniqueDice(dice)
		if containsSequence(sortedUniqueDice, []int{1, 2, 3, 4}) ||
			containsSequence(sortedUniqueDice, []int{2, 3, 4, 5}) ||
			containsSequence(sortedUniqueDice, []int{3, 4, 5, 6}) {
			return 30
		}
		return 0
	case "Large Straight":
		sortedUniqueDice := getSortedUniqueDice(dice)
		if containsSequence(sortedUniqueDice, []int{1, 2, 3, 4, 5}) ||
			containsSequence(sortedUniqueDice, []int{2, 3, 4, 5, 6}) {
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
		return sumAllDice
	default:
		return 0
	}
}

func getSortedUniqueDice(dice []int) []int {
	unique := make(map[int]bool)
	for _, d := range dice {
		unique[d] = true
	}
	sortedUnique := make([]int, 0, len(unique))
	for d := range unique {
		sortedUnique = append(sortedUnique, d)
	}
	sort.Ints(sortedUnique)
	return sortedUnique
}

func containsSequence(source, sequence []int) bool {
	if len(source) < len(sequence) {
		return false
	}
	for i := 0; i <= len(source)-len(sequence); i++ {
		match := true
		for j := 0; j < len(sequence); j++ {
			if source[i+j] != sequence[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func printDice(dice []int) {
	fmt.Print("Your dice: [")
	for i, d := range dice {
		fmt.Printf("%d", d)
		if i < len(dice)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")
}

func printAvailableScores(dice []int, usedCategories map[string]bool) {
	fmt.Println("\n--- Available Scores ---")
	categories := []string{
		"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes",
		"Three of a Kind", "Four of a Kind", "Full House",
		"Small Straight", "Large Straight", "Yahtzee", "Chance",
	}

	for _, cat := range categories {
		if !usedCategories[cat] {
			score := calculateScore(dice, cat)
			fmt.Printf("%-18s: %d\n", cat, score)
		} else {
			fmt.Printf("%-18s: (USED)\n", cat)
		}
	}
	fmt.Println("------------------------")
}

func getPlayerInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func playYahtzeeTurn(playerScores map[string]int, usedCategories map[string]bool) {
	fmt.Println("\n--- Starting New Turn ---")
	currentDice := rollDice()
	printDice(currentDice)

	rerollsLeft := 2
	for rerollsLeft > 0 {
		input := getPlayerInput(fmt.Sprintf("Rerolls left: %d. Enter dice to reroll (e.g., '1 3 5') or 'n' to keep: ", rerollsLeft))
		if strings.ToLower(input) == "n" {
			break
		}

		parts := strings.Fields(input)
		indicesToReroll := []int{}
		for _, p := range parts {
			idx, err := strconv.Atoi(p)
			if err == nil && idx >= 1 && idx <= 5 {
				indicesToReroll = append(indicesToReroll, idx)
			}
		}

		if len(indicesToReroll) > 0 {
			currentDice = rerollDice(currentDice, indicesToReroll)
			printDice(currentDice)
			rerollsLeft--
		} else {
			fmt.Println("Invalid input. Please enter 1-5 or 'n'.")
		}
	}

	printAvailableScores(currentDice, usedCategories)

	validCategoriesMap := make(map[string]bool)
	for _, cat := range []string{
		"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes",
		"Three of a Kind", "Four of a Kind", "Full House",
		"Small Straight", "Large Straight", "Yahtzee", "Chance",
	} {
		validCategoriesMap[cat] = true
	}

	for {
		chosenCategory := getPlayerInput("Choose a category to score (e.g., 'Yahtzee'): ")
		if !validCategoriesMap[chosenCategory] {
			fmt.Println("Invalid category name. Please choose from the list.")
			continue
		}
		if usedCategories[chosenCategory] {
			fmt.Println("Category already used. Please choose another one.")
			continue
		}

		score := calculateScore(currentDice, chosenCategory)
		playerScores[chosenCategory] = score
		usedCategories[chosenCategory] = true
		fmt.Printf("Scored %d points for %s.\n", score, chosenCategory)
		break
	}
}

func calculateTotalScore(scores map[string]int) int {
	total := 0
	upperSectionScore := 0
	upperSectionCategories := []string{"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes"}

	for cat, score := range scores {
		total += score
		for _, upperCat := range upperSectionCategories {
			if cat == upperCat {
				upperSectionScore += score
				break
			}
		}
	}

	if upperSectionScore >= 63 {
		fmt.Println("Upper section bonus: +35 points!")
		total += 35
	}
	return total
}

func main() {
	rand.Seed(time.Now().UnixNano())

	playerScores := make(map[string]int)
	usedCategories := make(map[string]bool)

	allCategories := []string{
		"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes",
		"Three of a Kind", "Four of a Kind", "Full House",
		"Small Straight", "Large Straight", "Yahtzee", "Chance",
	}
	for _, cat := range allCategories {
		usedCategories[cat] = false
	}

	totalTurns := 13
	for turn := 1; turn <= totalTurns; turn++ {
		fmt.Printf("\n===== Turn %d/%d =====\n", turn, totalTurns)
		playYahtzeeTurn(playerScores, usedCategories)

		fmt.Println("\n--- Current Scorecard ---")
		for _, cat := range allCategories {
			score, ok := playerScores[cat]
			if ok {
				fmt.Printf("%-18s: %d\n", cat, score)
			} else {
				fmt.Printf("%-18s: -\n", cat)
			}
		}
		fmt.Printf("Total Score: %d\n", calculateTotalScore(playerScores))

		if turn < totalTurns {
			_ = getPlayerInput("Press Enter to continue to next turn...")
		}
	}

	fmt.Println("\n===== Game Over! =====")
	fmt.Println("Final Scorecard:")
	for _, cat := range allCategories {
		score, ok := playerScores[cat]
		if ok {
			fmt.Printf("%-18s:

// Additional implementation at 2025-06-21 02:53:43
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
type Dice [5]int

// ScoreCategory represents the different scoring categories in Yahtzee.
type ScoreCategory string

const (
	Aces         ScoreCategory = "Aces"
	Twos         ScoreCategory = "Twos"
	Threes       ScoreCategory = "Threes"
	Fours        ScoreCategory = "Fours"
	Fives        ScoreCategory = "Fives"
	Sixes        ScoreCategory = "Sixes"
	ThreeOfAKind ScoreCategory = "Three of a Kind"
	FourOfAKind  ScoreCategory = "Four of a Kind"
	FullHouse    ScoreCategory = "Full House"
	SmallStraight ScoreCategory = "Small Straight"
	LargeStraight ScoreCategory = "Large Straight"
	Yahtzee      ScoreCategory = "Yahtzee"
	Chance       ScoreCategory = "Chance"
)

// AllCategories lists all available scoring categories.
var AllCategories = []ScoreCategory{
	Aces, Twos, Threes, Fours, Fives, Sixes,
	ThreeOfAKind, FourOfAKind, FullHouse, SmallStraight, LargeStraight, Yahtzee, Chance,
}

// ScoreCard holds the scores for each category and tracks used categories.
type ScoreCard struct {
	Scores      map[ScoreCategory]int
	Used        map[ScoreCategory]bool
	YahtzeeBonus int // Tracks additional Yahtzee bonuses
}

// NewScoreCard initializes a new scorecard.
func NewScoreCard() *ScoreCard {
	sc := &ScoreCard{
		Scores: make(map[ScoreCategory]int),
		Used:   make(map[ScoreCategory]bool),
	}
	for _, cat := range AllCategories {
		sc.Scores[cat] = 0 // Initialize all scores to 0
		sc.Used[cat] = false
	}
	return sc
}

// Roll rolls all 5 dice.
func (d *Dice) Roll() {
	for i := range d {
		d[i] = rand.Intn(6) + 1 // Generates a number between 1 and 6
	}
}

// Reroll rerolls specific dice based on their 0-indexed positions.
func (d *Dice) Reroll(indices []int) {
	for _, idx := range indices {
		if idx >= 0 && idx < 5 {
			d[idx] = rand.Intn(6) + 1
		}
	}
}

// calculateCounts returns a map of dice value counts.
func (d Dice) calculateCounts() map[int]int {
	counts := make(map[int]int)
	for _, die := range d {
		counts[die]++
	}
	return counts
}

// calculateScore calculates the potential score for a given category and dice.
func (sc *ScoreCard) calculateScore(dice Dice, category ScoreCategory) int {
	counts := dice.calculateCounts()
	score := 0

	switch category {
	case Aces:
		score = counts[1] * 1
	case Twos:
		score = counts[2] * 2
	case Threes:
		score = counts[3] * 3
	case Fours:
		score = counts[4] * 4
	case Fives:
		score = counts[5] * 5
	case Sixes:
		score = counts[6] * 6
	case ThreeOfAKind:
		for _, count := range counts {
			if count >= 3 {
				score = dice[0] + dice[1] + dice[2] + dice[3] + dice[4]
				break
			}
		}
	case FourOfAKind:
		for _, count := range counts {
			if count >= 4 {
				score = dice[0] + dice[1] + dice[2] + dice[3] + dice[4]
				break
			}
		}
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
			score = 25
		}
	case SmallStraight:
		sortedDice := make([]int, 5)
		copy(sortedDice, dice[:])
		sort.Ints(sortedDice)

		uniqueDice := make(map[int]bool)
		for _, dVal := range sortedDice {
			uniqueDice[dVal] = true
		}
		
		if len(uniqueDice) >= 4 {
			if (uniqueDice[1] && uniqueDice[2] && uniqueDice[3] && uniqueDice[4]) ||
				(uniqueDice[2] && uniqueDice[3] && uniqueDice[4] && uniqueDice[5]) ||
				(uniqueDice[3] && uniqueDice[4] && uniqueDice[5] && uniqueDice[6]) {
				score = 30
			}
		}
	case LargeStraight:
		sortedDice := make([]int, 5)
		copy(sortedDice, dice[:])
		sort.Ints(sortedDice)

		if (sortedDice[0] == 1 && sortedDice[1] == 2 && sortedDice[2] == 3 && sortedDice[3] == 4 && sortedDice[4] == 5) ||
			(sortedDice[0] == 2 && sortedDice[1] == 3 && sortedDice[2] == 4 && sortedDice[3] == 5 && sortedDice[4] == 6) {
			score = 40
		}
	case Yahtzee:
		for _, count := range counts {
			if count == 5 {
				score = 50
				break
			}
		}
	case Chance:
		score = dice[0] + dice[1] + dice[2] + dice[3] + dice[4]
	}
	return score
}

// recordScore records the score for a chosen category.
// Returns true if successful, false if category is already used.
func (sc *ScoreCard) recordScore(dice Dice, category ScoreCategory) bool {
	if sc.Used[category] {
		return false // Category already used
	}

	potentialScore := sc.calculateScore(dice, category)

	// Simplified Yahtzee bonus: if Yahtzee is already scored 50 and another Yahtzee is rolled
	if category == Yahtzee && potentialScore == 50 {
		if sc.Used[Yahtzee] && sc.Scores[Yahtzee] == 50 {
			sc.YahtzeeBonus += 100 // Add 100 points for each additional Yahtzee
		}
	}

	sc.Scores[category] = potentialScore
	sc.Used[category] = true
	return true
}

// GetTotalScore calculates the final score including bonuses.
func (sc *ScoreCard) GetTotalScore() int {
	total := 0
	upperScore := 0
	for _, cat := range []ScoreCategory{Aces, Twos, Threes, Fours, Fives, Sixes} {
		total += sc.Scores[cat]
		upperScore += sc.Scores[cat]
	}

	if upperScore >= 63 {
		total += 35 // Upper section bonus
	}

	for _, cat := range []ScoreCategory{ThreeOfAKind, FourOfAKind, FullHouse, SmallStraight, LargeStraight, Yahtzee, Chance} {
		total += sc.Scores[cat]
	}

	total += sc.YahtzeeBonus // Add Yahtzee bonuses

	return total
}

// Game represents the state of a Yahtzee game.
type Game struct {
	Dice     Dice
	ScoreCard *ScoreCard
	Round    int
	Reader   *bufio.Reader
}

// NewGame initializes a new Yahtzee game.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	return &Game{
		Dice:     Dice{},
		ScoreCard: NewScoreCard(),
		Round:    0,
		Reader:   bufio.NewReader(os.Stdin),
	}
}

// displayDice prints the current state of the dice.
func (g *Game) displayDice() {
	fmt.Print("Current Dice: [")
	for i, d := range g.Dice {
		fmt.Printf("%d", d)
		if i < 4 {
			fmt.Print(", ")
		}
	}
	fmt.Println("]")
}

// displayScoreCard prints the current scorecard.
func (g *Game) displayScoreCard() {
	fmt.Println("\n--- Scorecard ---")
	upperScore := 0
	for _, cat := range []ScoreCategory{Aces, Twos, Threes, Fours, Fives, Sixes} {
		status := " (Available)"
		if g.ScoreCard.Used[cat] {
			status = ""
		}
		fmt.Printf("%-15s: %2d%s\n", cat, g.ScoreCard.Scores[cat], status)
		upperScore += g.ScoreCard.Scores[cat]
	}
	fmt.Printf("Upper Total      : %2d\n", upperScore)
	if upperScore >= 63 {
		fmt.Println("Upper Bonus      : 35 (Achieved!)")
	} else {
		fmt.Printf("Upper Bonus      : 0 (Need %d more for 35 bonus)\n", 63-upperScore)
	}

	fmt.Println("-----------------")
	for _, cat := range []ScoreCategory{ThreeOfAKind, FourOfAKind, FullHouse, SmallStraight, LargeStraight, Yahtzee, Chance} {
		status := " (Available)"
		if g.ScoreCard.Used[cat] {
			status = ""
		}
		fmt.Printf("%-15s: %2d%s\n", cat, g.ScoreCard.Scores[cat], status)
	}
	if g.ScoreCard.YahtzeeBonus > 0 {
		fmt.Printf("Yahtzee Bonus    : %d\n", g.ScoreCard.YahtzeeBonus)
	}
	fmt.Println("-----------------")
	fmt.Printf("Grand Total      : %2d\n", g.ScoreCard.GetTotalScore())
	fmt.Println("-----------------")
}

// getPlayerInput reads a line from stdin.
func (g *Game) getPlayerInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := g.Reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// playTurn manages a single turn of the game.
func (g *Game) playTurn() {
	g.Round++
	fmt.Printf("\n--- Round %d/13 ---\n", g.Round)

	g.Dice.Roll()
	g.displayDice()

	rollsLeft := 2
	for rollsLeft > 0 {
		input := g.getPlayerInput(fmt.Sprintf("Reroll? (y/n, %d rolls left): ", rollsLeft))
		if strings.ToLower(input) != "y" {
			break
		}

		fmt.Println("Enter dice numbers to reroll (1-5, e.g., '1 3 5'). Press Enter to skip:")
		input = g.getPlayerInput("Dice to reroll: ")
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		var indicesToReroll []int
		for _, p := range parts {
			idx, err := strconv.Atoi(p)
			if err == nil && idx >= 1 && idx <= 5 {
				indicesToReroll = append(indicesToReroll, idx-1) // Convert to 0-indexed
			} else {
				fmt.Printf("Invalid input: '%s'. Please enter numbers between 1 and 5.\n", p)
			}
		}
		if len(indicesToReroll) > 0 {
			g.Dice.Reroll(indicesToReroll)
			g.displayDice()
			rollsLeft--
		} else {
			fmt.Println("No valid dice selected for reroll.")
		}
	}

	g.displayScoreCard()

	// Choose category to score
	for {
		fmt.Println("\nAvailable categories to score:")
		var availableCategories []ScoreCategory
		for _, cat := range AllCategories {
			if !g.ScoreCard.Used[cat] {
				potentialScore :=