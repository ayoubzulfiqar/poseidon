package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type GameState struct {
	CurrentLocationID string          `json:"current_location_id"`
	Inventory         []string        `json:"inventory"`
	VisitedLocations  map[string]bool `json:"visited_locations"`
}

type Location struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	Items       []string          `json:"items"`
}

var (
	gameMap   = make(map[string]Location)
	gameState GameState
	reader    *bufio.Reader
)

const saveFileName = "adventure_save.json"

func initGame() {
	gameMap["hallway"] = Location{
		ID:          "hallway",
		Name:        "Hallway",
		Description: "You are in a dimly lit hallway. There's a door to the north and a dusty portrait on the wall.",
		Exits:       map[string]string{"north": "kitchen", "east": "living_room"},
		Items:       []string{"dusty portrait"},
	}
	gameMap["kitchen"] = Location{
		ID:          "kitchen",
		Name:        "Kitchen",
		Description: "A small, somewhat messy kitchen. A sharp knife lies on the counter.",
		Exits:       map[string]string{"south": "hallway"},
		Items:       []string{"knife"},
	}
	gameMap["living_room"] = Location{
		ID:          "living_room",
		Name:        "Living Room",
		Description: "A cozy living room with a fireplace. A thick book rests on a coffee table.",
		Exits:       map[string]string{"west": "hallway"},
		Items:       []string{"book"},
	}

	gameState = GameState{
		CurrentLocationID: "hallway",
		Inventory:         []string{},
		VisitedLocations:  make(map[string]bool),
	}
	gameState.VisitedLocations[gameState.CurrentLocationID] = true

	reader = bufio.NewReader(os.Stdin)
}

func printMessage(msg string) {
	fmt.Println("\n" + msg)
}

func look() {
	currentLoc := gameMap[gameState.CurrentLocationID]
	printMessage(currentLoc.Name)
	printMessage(currentLoc.Description)

	if len(currentLoc.Items) > 0 {
		fmt.Print("You see: ")
		for i, item := range currentLoc.Items {
			fmt.Print(item)
			if i < len(currentLoc.Items)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println(".")
	}

	fmt.Print("Exits: ")
	exits := []string{}
	for dir := range currentLoc.Exits {
		exits = append(exits, dir)
	}
	if len(exits) > 0 {
		fmt.Println(strings.Join(exits, ", ") + ".")
	} else {
		fmt.Println("None.")
	}
}

func goDirection(direction string) {
	currentLoc := gameMap[gameState.CurrentLocationID]
	nextLocID, exists := currentLoc.Exits[direction]
	if !exists {
		printMessage("You can't go that way.")
		return
	}

	gameState.CurrentLocationID = nextLocID
	gameState.VisitedLocations[nextLocID] = true
	look()
}

func takeItem(itemName string) {
	currentLoc := gameMap[gameState.CurrentLocationID]
	itemFound := false
	newItemsInLocation := []string{}

	for _, item := range currentLoc.Items {
		if strings.EqualFold(item, itemName) {
			gameState.Inventory = append(gameState.Inventory, item)
			printMessage(fmt.Sprintf("You picked up the %s.", item))
			itemFound = true
		} else {
			newItemsInLocation = append(newItemsInLocation, item)
		}
	}

	if itemFound {
		currentLoc.Items = newItemsInLocation
		gameMap[gameState.CurrentLocationID] = currentLoc
	} else {
		printMessage(fmt.Sprintf("There is no %s here.", itemName))
	}
}

func showInventory() {
	if len(gameState.Inventory) == 0 {
		printMessage("Your inventory is empty.")
		return
	}
	printMessage("Your inventory:")
	for _, item := range gameState.Inventory {
		fmt.Println("- " + item)
	}
}

func saveGame() {
	data, err := json.MarshalIndent(gameState, "", "  ")
	if err != nil {
		printMessage(fmt.Sprintf("Error saving game: %v", err))
		return
	}

	err = os.WriteFile(saveFileName, data, 0644)
	if err != nil {
		printMessage(fmt.Sprintf("Error writing save file: %v", err))
		return
	}
	printMessage("Game saved successfully!")
}

func loadGame() {
	data, err := os.ReadFile(saveFileName)
	if err != nil {
		if os.IsNotExist(err) {
			printMessage("No saved game found.")
		} else {
			printMessage(fmt.Sprintf("Error loading game: %v", err))
		}
		return
	}

	var loadedState GameState
	err = json.Unmarshal(data, &loadedState)
	if err != nil {
		printMessage(fmt.Sprintf("Error parsing save file: %v", err))
		return
	}

	gameState = loadedState
	printMessage("Game loaded successfully!")
	look()
}

func handleCommand(command string) bool {
	parts := strings.Fields(strings.ToLower(command))
	if len(parts) == 0 {
		return true
	}

	verb := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	switch verb {
	case "go":
		if len(args) > 0 {
			goDirection(args[0])
		} else {
			printMessage("Go where?")
		}
	case "look":
		look()
	case "take":
		if len(args) > 0 {
			takeItem(strings.Join(args, " "))
		} else {
			printMessage("Take what?")
		}
	case "inventory", "inv":
		showInventory()
	case "save":
		saveGame()
	case "load":
		loadGame()
	case "quit", "exit":
		printMessage("Goodbye!")
		return false
	default:
		printMessage("I don't understand that command.")
	}
	return true
}

func main() {
	initGame()

	printMessage("Welcome to the Simple Text Adventure!")
	printMessage("Type 'help' for commands (though 'look', 'go <direction>', 'take <item>', 'inventory', 'save', 'load', 'quit' are the main ones).")
	look()

	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "help" {
			printMessage("Available commands:")
			printMessage("- look: Describe your current location.")
			printMessage("- go <direction>: Move in a direction (e.g., 'go north').")
			printMessage("- take <item>: Pick up an item (e.g., 'take knife').")
			printMessage("- inventory / inv: Show items in your inventory.")
			printMessage("- save: Save your current game progress.")
			printMessage("- load: Load a previously saved game.")
			printMessage("- quit / exit: End the game.")
			continue
		}

		if !handleCommand(input) {
			break
		}
	}
}

// Additional implementation at 2025-06-23 00:56:13
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Room struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"` // direction -> roomID
}

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CanTake     bool   `json:"canTake"`
}

type GameState struct {
	PlayerLocation string                     `json:"playerLocation"`
	Inventory      map[string]bool            `json:"inventory"` // itemID -> present
	RoomItems      map[string]map[string]bool `json:"roomItems"` // roomID -> itemID -> present
	VisitedRooms   map[string]bool            `json:"visitedRooms"` // roomID -> visited
}

var rooms = map[string]Room{
	"start_room": {
		ID:          "start_room",
		Name:        "Dusty Old Study",
		Description: "You are in a dusty old study. Bookshelves line the walls, filled with ancient tomes. A flickering gas lamp illuminates a worn wooden desk.",
		Exits:       map[string]string{"north": "hallway", "east": "library"},
	},
	"hallway": {
		ID:          "hallway",
		Name:        "Dimly Lit Hallway",
		Description: "A long, dimly lit hallway stretches before you. Cobwebs hang from the ceiling.",
		Exits:       map[string]string{"south": "start_room", "west": "kitchen"},
	},
	"kitchen": {
		ID:          "kitchen",
		Name:        "Grubby Kitchen",
		Description: "The kitchen is surprisingly clean, but smells faintly of old spices. A single, sharp knife lies on the counter.",
		Exits:       map[string]string{"east": "hallway"},
	},
	"library": {
		ID:          "library",
		Name:        "Grand Library",
		Description: "An enormous library, filled with countless books. A strange, glowing orb rests on a pedestal in the center.",
		Exits:       map[string]string{"west": "start_room"},
	},
}

var items = map[string]Item{
	"key": {
		ID:          "key",
		Name:        "small key",
		Description: "A small, tarnished brass key. It looks like it might open something important.",
		CanTake:     true,
	},
	"knife": {
		ID:          "knife",
		Name:        "sharp knife",
		Description: "A very sharp kitchen knife. Useful for cutting, or perhaps for self-defense.",
		CanTake:     true,
	},
	"orb": {
		ID:          "orb",
		Name:        "glowing orb",
		Description: "A pulsating, ethereal orb. It hums with a faint energy. This must be the artifact!",
		CanTake:     true,
	},
	"book": {
		ID:          "book",
		Name:        "ancient book",
		Description: "An ancient, leather-bound book. Its pages are brittle and filled with unreadable script.",
		CanTake:     false,
	},
}

var currentGameState GameState

const saveFileName = "adventure_save.json"

func InitializeGame() {
	currentGameState = GameState{
		PlayerLocation: "start_room",
		Inventory:      make(map[string]bool),
		RoomItems:      make(map[string]map[string]bool),
		VisitedRooms:   make(map[string]bool),
	}

	currentGameState.RoomItems["start_room"] = map[string]bool{"key": true, "book": true}
	currentGameState.RoomItems["kitchen"] = map[string]bool{"knife": true}
	currentGameState.RoomItems["library"] = map[string]bool{"orb": true}

	currentGameState.VisitedRooms[currentGameState.PlayerLocation] = true
}

func SaveGame() {
	data, err := json.MarshalIndent(currentGameState, "", "  ")
	if err != nil {
		fmt.Println("Error saving game:", err)
		return
	}

	err = ioutil.WriteFile(saveFileName, data, 0644)
	if err != nil {
		fmt.Println("Error writing save file:", err)
		return
	}
	fmt.Println("Game saved successfully!")
}

func LoadGame() bool {
	data, err := ioutil.ReadFile(saveFileName)
	if err != nil {
		fmt.Println("No saved game found or error reading file:", err)
		return false
	}

	err = json.Unmarshal(data, &currentGameState)
	if err != nil {
		fmt.Println("Error loading game:", err)
		return false
	}
	fmt.Println("Game loaded successfully!")
	return true
}

func describeCurrentLocation() {
	room, exists := rooms[currentGameState.PlayerLocation]
	if !exists {
		fmt.Println("Error: Player is in an unknown location.")
		return
	}

	fmt.Println("\n---", room.Name, "---")
	fmt.Println(room.Description)

	if len(currentGameState.RoomItems[room.ID]) > 0 {
		fmt.Print("You see: ")
		first := true
		for itemID := range currentGameState.RoomItems[room.ID] {
			if currentGameState.RoomItems[room.ID][itemID] {
				if !first {
					fmt.Print(", ")
				}
				fmt.Print(items[itemID].Name)
				first = false
			}
		}
		fmt.Println(".")
	}

	if len(room.Exits) > 0 {
		fmt.Print("Exits: ")
		first := true
		for dir := range room.Exits {
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(dir)
			first = false
		}
		fmt.Println(".")
	}
}

func handleGoCommand(direction string) {
	room, exists := rooms[currentGameState.PlayerLocation]
	if !exists {
		fmt.Println("Error: Player is in an unknown location.")
		return
	}

	nextRoomID, hasExit := room.Exits[direction]
	if !hasExit {
		fmt.Println("You can't go that way.")
		return
	}

	currentGameState.PlayerLocation = nextRoomID
	currentGameState.VisitedRooms[currentGameState.PlayerLocation] = true
	describeCurrentLocation()
}

func handleLookCommand(target string) {
	if target == "" {
		describeCurrentLocation()
		return
	}

	roomItems, roomHasItems := currentGameState.RoomItems[currentGameState.PlayerLocation]
	if roomHasItems {
		for itemID := range roomItems {
			if roomItems[itemID] && strings.EqualFold(items[itemID].Name, target) {
				fmt.Println(items[itemID].Description)
				return
			}
		}
	}

	for itemID := range currentGameState.Inventory {
		if currentGameState.Inventory[itemID] && strings.EqualFold(items[itemID].Name, target) {
			fmt.Println(items[itemID].Description)
			return
		}
	}

	fmt.Println("You don't see a '" + target + "' here or in your inventory.")
}

func handleTakeCommand(itemName string) {
	roomID := currentGameState.PlayerLocation
	roomItems, roomHasItems := currentGameState.RoomItems[roomID]

	if !roomHasItems || len(roomItems) == 0 {
		fmt.Println("There's nothing to take here.")
		return
	}

	found := false
	for itemID := range roomItems {
		if roomItems[itemID] && strings.EqualFold(items[itemID].Name, itemName) {
			if items[itemID].CanTake {
				currentGameState.Inventory[itemID] = true
				delete(currentGameState.RoomItems[roomID], itemID)
				fmt.Println("You took the", items[itemID].Name + ".")
				found = true
				break
			} else {
				fmt.Println("You can't take the", items[itemID].Name + ".")
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("You don't see a '" + itemName + "' here.")
	}
}

func handleDropCommand(itemName string) {
	found := false
	for itemID := range currentGameState.Inventory {
		if currentGameState.Inventory[itemID] && strings.EqualFold(items[itemID].Name, itemName) {
			delete(currentGameState.Inventory, itemID)

			if _, ok := currentGameState.RoomItems[currentGameState.PlayerLocation]; !ok {
				currentGameState.RoomItems[currentGameState.PlayerLocation] = make(map[string]bool)
			}
			currentGameState.RoomItems[currentGameState.PlayerLocation][itemID] = true

			fmt.Println("You dropped the", items[itemID].Name + ".")
			found = true
			break
		}
	}

	if !found {
		fmt.Println("You don't have a '" + itemName + "' to drop.")
	}
}

func handleInventoryCommand() {
	if len(currentGameState.Inventory) == 0 {
		fmt.Println("Your inventory is empty.")
		return
	}

	fmt.Println("Your inventory:")
	for itemID := range currentGameState.Inventory {
		if currentGameState.Inventory[itemID] {
			fmt.Println("- " + items[itemID].Name)
		}
	}
}

func handleHelpCommand() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  go <direction> (e.g., go north)")
	fmt.Println("  look           (describe current room)")
	fmt.Println("  look <item>    (describe an item)")
	fmt.Println("  take <item>    (pick up an item)")
	fmt.Println("  drop <item>    (put down an item)")
	fmt.Println("  inventory      (list items in your inventory)")
	fmt.Println("  save           (save your game)")
	fmt.Println("  load           (load your game)")
	fmt.Println("  help           (show this list)")
	fmt.Println("  quit           (exit the game)")
}

func checkWinCondition() bool {
	if currentGameState.PlayerLocation == "library" {
		if currentGameState.Inventory["orb"] {
			fmt.Println("\n--- CONGRATULATIONS! ---")
			fmt.Println("You have found the glowing orb and brought it back to the library!")
			fmt.Println("You have successfully completed your quest!")
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Go Text Adventure!")
	fmt.Println("Type 'load' to load a saved game, or press Enter to start a new game.")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "load" {
		if !LoadGame() {
			fmt.Println("Starting a new game...")
			InitializeGame()
		}
	} else {
		fmt.Println("Starting a new game...")
		InitializeGame()
	}

	describeCurrentLocation()
	handleHelpCommand()

	for {
		if checkWinCondition() {
			break
		}

		fmt.Print("\nWhat do you want to do? > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(strings.ToLower(input))

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = strings.Join(parts[1:], " ")
		}

		switch command {
		case "go":
			if arg == "" {
				fmt.Println("Go where? (e.g., go north)")
			} else {
				handleGoCommand(arg)
			}
		case "look":
			handleLookCommand(arg)
		case "take":
			if arg == "" {
				fmt.Println("Take what? (e.g., take key)")
			} else {
				handleTakeCommand(arg)
			}
		case "drop":
			if arg == "" {
				fmt.Println("Drop what? (e.g., drop key)")
			} else {
				handleDropCommand(arg)
			}
		case "inventory", "inv":
			handleInventoryCommand()
		case "save":
			SaveGame()
		case "load":
			LoadGame()
			describeCurrentLocation()
		case "help":
			handleHelpCommand()
		case "quit", "exit":
			fmt.Println("Thanks for playing!")
			return
		default:
			fmt.Println("I don't understand that command. Type 'help' for a list of commands.")
		}
	}
}