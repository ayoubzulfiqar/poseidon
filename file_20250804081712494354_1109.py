import json
import os

SAVE_FILE = "adventure_save.json"

rooms = {
    "start_room": {
        "name": "Old Study",
        "description": "You are in an old, dusty study. Bookshelves line the walls, and a flickering candle sits on a desk. A faint draft comes from the north.",
        "exits": {"north": "hallway"},
        "items": ["old key", "dusty book"],
    },
    "hallway": {
        "name": "Dark Hallway",
        "description": "A long, dark hallway stretches before you. To the south is the study. To the east, you hear dripping water.",
        "exits": {"south": "start_room", "east": "damp_cave"},
        "items": [],
    },
    "damp_cave": {
        "name": "Damp Cave",
        "description": "The air is thick with moisture, and water drips from the stalactites. A strange glow emanates from a crack in the wall to the west.",
        "exits": {"west": "hallway", "crack": "treasure_room"},
        "items": ["glowing orb"],
    },
    "treasure_room": {
        "name": "Treasure Room",
        "description": "You've stumbled upon a hidden chamber filled with glittering gold and ancient artifacts! You win!",
        "exits": {"back": "damp_cave"},
        "items": ["golden idol"],
    },
}

game_state = {
    "current_room": "start_room",
    "inventory": [],
    "game_over": False,
    "has_orb_for_crack": False,
}

def save_game():
    try:
        with open(SAVE_FILE, "w") as f:
            json.dump(game_state, f, indent=4)
        print("Game saved successfully!")
    except IOError:
        print("Error: Could not save game.")

def load_game():
    global game_state
    if os.path.exists(SAVE_FILE):
        try:
            with open(SAVE_FILE, "r") as f:
                loaded_state = json.load(f)
                if all(key in loaded_state for key in ["current_room", "inventory", "game_over", "has_orb_for_crack"]):
                    game_state = loaded_state
                    print("Game loaded successfully!")
                    return True
                else:
                    print("Error: Save file corrupted or invalid. Starting new game.")
                    return False
        except (IOError, json.JSONDecodeError):
            print("Error: Could not load game or save file is corrupted. Starting new game.")
            return False
    else:
        print("No saved game found. Starting new game.")
        return False

def display_room():
    room_id = game_state["current_room"]
    room = rooms[room_id]
    print(f"\n--- {room['name']} ---")
    print(room["description"])

    if room["items"]:
        print("You see:", ", ".join(room["items"]) + ".")
    else:
        print("There are no items here.")

    exits = [f"{direction} ({target_room})" for direction, target_room in room["exits"].items()]
    print("Exits:", ", ".join(exits) + ".")

def move_player(direction):
    room_id = game_state["current_room"]
    room = rooms[room_id]
    if direction in room["exits"]:
        next_room_id = room["exits"][direction]
        if next_room_id == "treasure_room" and not game_state["has_orb_for_crack"]:
            print("The crack is too small to fit through. Perhaps something could widen it?")
        else:
            game_state["current_room"] = next_room_id
            if next_room_id == "treasure_room":
                print("\nCongratulations! You found the treasure and won the game!")
                game_state["game_over"] = True
            display_room()
    else:
        print("You can't go that way.")

def take_item(item_name):
    room_id = game_state["current_room"]
    room = rooms[room_id]
    if item_name in room["items"]:
        room["items"].remove(item_name)
        game_state["inventory"].append(item_name)
        print(f"You picked up the {item_name}.")
        if item_name == "glowing orb":
            game_state["has_orb_for_crack"] = True
            print("The glowing orb seems to hum faintly in your hand.")
    else:
        print(f"There is no {item_name} here.")

def drop_item(item_name):
    if item_name in game_state["inventory"]:
        game_state["inventory"].remove(item_name)
        rooms[game_state["current_room"]]["items"].append(item_name)
        print(f"You dropped the {item_name}.")
        if item_name == "glowing orb":
            game_state["has_orb_for_crack"] = False
    else:
        print(f"You don't have a {item_name} in your inventory.")

def show_inventory():
    if game_state["inventory"]:
        print("Inventory:", ", ".join(game_state["inventory"]) + ".")
    else:
        print("Your inventory is empty.")

def game():
    print("Welcome to the Python Text Adventure!")
    print("Type 'help' for commands.")

    while True:
        choice = input("Load saved game? (yes/no): ").lower().strip()
        if choice == "yes":
            if load_game():
                break
            else:
                print("Starting a new game.")
                break
        elif choice == "no":
            print("Starting a new game.")
            break
        else:
            print("Invalid choice. Please type 'yes' or 'no'.")

    display_room()

    while not game_state["game_over"]:
        command = input("\nWhat do you do? ").lower().strip().split(maxsplit=1)
        action = command[0]
        arg = command[1] if len(command) > 1 else ""

        if action == "go":
            move_player(arg)
        elif action == "take":
            take_item(arg)
        elif action == "drop":
            drop_item(arg)
        elif action == "inventory" or action == "inv":
            show_inventory()
        elif action == "look":
            display_room()
        elif action == "save":
            save_game()
        elif action == "load":
            load_game()
            display_room()
        elif action == "quit" or action == "exit":
            print("Exiting game. Goodbye!")
            game_state["game_over"] = True
        elif action == "help":
            print("\nAvailable commands:")
            print("  go [direction] - Move in a direction (e.g., 'go north')")
            print("  take [item]    - Pick up an item (e.g., 'take key')")
            print("  drop [item]    - Drop an item from your inventory")
            print("  inventory / inv - Show your inventory")
            print("  look           - Look around the current room again")
            print("  save           - Save your current game progress")
            print("  load           - Load a previously saved game")
            print("  quit / exit    - Exit the game")
        else:
            print("I don't understand that command. Type 'help' for a list of commands.")

if __name__ == "__main__":
    game()

# Additional implementation at 2025-08-04 08:18:18
import pickle
import os
import sys
import random

class Player:
    def __init__(self, name="Hero"):
        self.name = name
        self.inventory = []
        self.current_location_id = "start_room"
        self.health = 100
        self.strength = 10
        self.max_health = 100

    def __str__(self):
        return f"Name: {self.name}\nHealth: {self.health}/{self.max_health}\nStrength: {self.strength}\nInventory: {', '.join(self.inventory) if self.inventory else 'Empty'}"

class Location:
    def __init__(self, id, name, description, exits, items=None, challenge=None, visited=False):
        self.id = id
        self.name = name
        self.description = description
        self.exits = exits  # {'direction': 'location_id'}
        self.items = items if items is not None else []
        self.challenge = challenge  # {'type': 'combat', 'name': 'Monster', 'difficulty': 5, 'reward': 'item'} or None
        self.visited = visited

    def __str__(self):
        return f"Location: {self.name}\nDescription: {self.description}\nExits: {', '.join(self.exits.keys())}\nItems: {', '.join(self.items) if self.items else 'None'}"

class Game:
    def __init__(self):
        self.player = Player()
        self.locations = {}
        self._create_world()
        self.current_location = self.locations[self.player.current_location_id]
        self.game_state_file = "adventure_save.pkl"
        self.game_over = False
        self.game_won = False

    def _create_world(self):
        self.locations["start_room"] = Location(
            id="start_room",
            name="Start Room",
            description="You are in a dimly lit room. There's a dusty old table and a flickering torch on the wall.",
            exits={"north": "forest_path", "east": "storage_room"},
            items=["rusty key"]
        )
        self.locations["forest_path"] = Location(
            id="forest_path",
            name="Forest Path",
            description="A winding path through a dense forest. Sunlight barely penetrates the canopy.",
            exits={"south": "start_room", "north": "dark_cave", "west": "river_bank"}
        )
        self.locations["dark_cave"] = Location(
            id="dark_cave",
            name="Dark Cave Entrance",
            description="The entrance to a dark, ominous cave. A chilling wind blows from within.",
            exits={"south": "forest_path", "east": "treasure_room_locked"},
            challenge={'type': 'combat', 'name': 'Giant Spider', 'difficulty': 15, 'reward': 'shiny sword'}
        )
        self.locations["treasure_room_locked"] = Location(
            id="treasure_room_locked",
            name="Blocked Passage",
            description="A narrow passage blocked by a heavy, ornate door. It seems to require a special key.",
            exits={"west": "dark_cave"}
        )
        self.locations["treasure_room"] = Location(
            id="treasure_room",
            name="Treasure Room",
            description="You've entered a magnificent chamber filled with glittering gold and ancient artifacts! A pedestal in the center holds a glowing orb.",
            exits={"west": "dark_cave"},
            items=["glowing orb"]
        )
        self.locations["river_bank"] = Location(
            id="river_bank",
            name="River Bank",
            description="The gentle sound of a flowing river fills the air. A small, rickety bridge crosses to the other side.",
            exits={"east": "forest_path", "north": "old_cabin"},
            items=["healing potion"]
        )
        self.locations["old_cabin"] = Location(
            id="old_cabin",
            name="Old Cabin",
            description="An abandoned cabin, its windows boarded up. The air inside is stale and dusty.",
            exits={"south": "river_bank"},
            items=["old map"]
        )
        self.locations["storage_room"] = Location(
            id="storage_room",
            name="Storage Room",
            description="A small, cluttered storage room. Boxes are stacked high, and cobwebs hang everywhere.",
            exits={"west": "start_room"},
            items=["torch"]
        )

    def _get_location(self, location_id):
        return self.locations.get(location_id)

    def display_location(self):
        loc = self.current_location
        if not loc.visited:
            print(f"\n--- {loc.name} ---")
            print(loc.description)
            loc.visited = True
        else:
            print(f"\n--- {loc.name} ---")
            print("You are back here.")

        if loc.items:
            print(f"You see: {', '.join(loc.items)}.")
        else:
            print("There are no items here.")

        print(f"Exits: {', '.join(loc.exits.keys())}.")

    def handle_input(self, command):
        command_parts = command.lower().split(maxsplit=1)
        action = command_parts[0]
        arg = command_parts[1] if len(command_parts) > 1 else None

        if action == "go":
            self.move(arg)
        elif action == "take":
            self.take_item(arg)
        elif action == "use":
            self.use_item(arg)
        elif action == "inventory" or action == "i":
            self.check_inventory()
        elif action == "look":
            self.display_location()
        elif action == "status" or action == "s":
            print(self.player)
        elif action == "save":
            self.save_game()
        elif action == "quit":
            self.quit_game()
        else:
            print("Invalid command. Try 'go [direction]', 'take [item]', 'use [item]', 'inventory', 'status', 'save', 'quit'.")

    def move(self, direction):
        if direction in self.current_location.exits:
            next_location_id = self.current_location.exits[direction]

            if next_location_id == "treasure_room_locked":
                if "rusty key" in self.player.inventory:
                    print("You use the rusty key to unlock the heavy door!")
                    # Permanently change the exit for this game instance
                    self.current_location.exits[direction] = "treasure_room"
                    next_location_id = "treasure_room"
                    self.player.inventory.remove("rusty key")
                else:
                    print("The door is locked. You need a key.")
                    return

            self.current_location = self._get_location(next_location_id)
            self.player.current_location_id = self.current_location.id
            self.display_location()
            self._handle_challenge()
        else:
            print("You can't go that way.")

    def take_item(self, item_name):
        if item_name in self.current_location.items:
            self.player.inventory.append(item_name)
            self.current_location.items.remove(item_name)
            print(f"You took the {item_name}.")
        else:
            print(f"There is no {item_name} here.")

    def use_item(self, item_name):
        if item_name not in self.player.inventory:
            print(f"You don't have a {item_name}.")
            return

        if item_name == "healing potion":
            heal_amount = 30
            self.player.health = min(self.player.max_health, self.player.health + heal_amount)
            self.player.inventory.remove(item_name)
            print(f"You drank the healing potion and recovered {heal_amount} health. Your health is now {self.player.health}.")
        elif item_name == "glowing orb":
            print("You hold the glowing orb. A warm light envelops you, and you feel a sense of profound accomplishment.")
            print("Congratulations! You have found the ancient artifact and won the game!")
            self.game_won = True
        elif item_name == "torch":
            print("You light the torch. The room brightens slightly, but nothing new is revealed here.")
        elif item_name == "old map":
            print("You unfold the old map. It shows a crude drawing of the surrounding area, but it's too faded to be very useful.")
        else:
            print(f"You can't seem to use the {item_name} in any meaningful way right now.")

    def check_inventory(self):
        if self.player.inventory:
            print(f"Your inventory: {', '.join(self.player.inventory)}")
        else:
            print("Your inventory is empty.")

    def _handle_challenge(self):
        loc = self.current_location
        if loc.challenge and loc.challenge['type'] == 'combat':
            challenge_name = loc.challenge['name']
            difficulty = loc.challenge['difficulty']
            reward = loc.challenge['reward']

            print(f"\n--- CHALLENGE: {challenge_name} ---")
            print(f"A {challenge_name} blocks your path! It looks formidable.")
            print(f"Your strength: {self.player.strength}. Required strength: {difficulty}.")

            if self.player.strength >= difficulty:
                print(f"You bravely confront the {challenge_name} and, with your superior strength, defeat it!")
                if reward and reward not in self.player.inventory:
                    self.player.inventory.append(reward)
                    print(f"You found a {reward}!")
                loc.challenge = None
            else:
                damage = random.randint(difficulty - self.player.strength, difficulty)
                self.player.health -= damage
                print(f"The {challenge_name} is too strong! You take {damage} damage and are forced to retreat!")
                print(f"Your health is now {self.player.health}.")
                if self.player.health <= 0:
                    print("Your health has dropped to zero. You collapse...")
                    self.game_over = True
                else:
                    print("You stumble back to the previous location.")
                    if self.current_location.id == "dark_cave":
                        self.current_location = self.locations["forest_path"]
                        self.player.current_location_id = "forest_path"
                        self.display_location()
                    else:
                        print("You are too disoriented to know where you are going.")

    def save_game(self):
        try:
            with open(self.game_state_file, 'wb') as f:
                pickle.dump(self, f)
            print("Game saved successfully!")
        except Exception as e:
            print(f"Error saving game: {e}")

    @classmethod
    def load_game(cls, filename):
        try:
            with open(filename, 'rb') as f:
                game_instance = pickle.load(f)
            print("Game loaded successfully!")
            return game_instance
        except FileNotFoundError:
            print("No saved game found.")
            return None
        except Exception as e:
            print(f"Error loading game: {e}")
            return None

    def quit_game(self):
        print("Thanks for playing!")
        self.game_over = True

    def run_game(self):
        print("Welcome to the Python Text Adventure!")
        if not self.game_over and not self.game_won:
            self.display_location()

        while not self.game_over and not self.game_won:
            command = input("\nWhat do you do? ").strip()
            if not command:
                continue
            self.handle_input(command)

            if self.player.health <= 0:
                self.game_over = True
                print("\n--- GAME OVER ---")
                print("You have perished in the adventure.")
            elif self.game_won:
                print("\n--- YOU WIN! ---")


# Additional implementation at 2025-08-04 08:19:33
import json
import sys
import os

game_state = {
    'current_location': 'forest_path',
    'inventory': [],
    'health': 100,
    'flags': {
        'has_key': False,
        'door_unlocked': False,
        'monster_defeated': False
    }
}

locations = {
    'forest_path': {
        'name': 'Forest Path',
        'description': 'You are on a winding path through a dense forest. Sunlight barely penetrates the canopy.',
        'exits': {'north': 'clearing', 'east': 'riverbank'},
        'items': ['old_map'],
        'events': []
    },
    'clearing': {
        'name': 'Forest Clearing',
        'description': 'A small clearing opens up. In the center, a gnarled old tree stands.',
        'exits': {'south': 'forest_path', 'west': 'cave_entrance'},
        'items': ['healing_herb'],
        'events': []
    },
    'riverbank': {
        'name': 'Riverbank',
        'description': 'A narrow river flows gently here. The water looks clear and inviting.',
        'exits': {'west': 'forest_path', 'north': 'waterfall'},
        'items': ['shiny_pebble'],
        'events': []
    },
    'waterfall': {
        'name': 'Behind the Waterfall',
        'description': 'You\'ve managed to squeeze behind the roaring waterfall. A hidden alcove is here.',
        'exits': {'south': 'riverbank'},
        'items': ['rusty_key'],
        'events': []
    },
    'cave_entrance': {
        'name': 'Cave Entrance',
        'description': 'A dark, foreboding cave mouth yawns before you. A strange symbol is carved above it.',
        'exits': {'east': 'clearing', 'north': 'dark_cave'},
        'items': [],
        'events': []
    },
    'dark_cave': {
        'name': 'Dark Cave',
        'description': 'It\'s pitch black in here. You can hear dripping water and a faint growling sound.',
        'exits': {'south': 'cave_entrance'},
        'items': [],
        'events': ['monster_encounter'],
        'requires_flag': {'flag': 'door_unlocked', 'value': True, 'message': 'The cave entrance is blocked by a magical barrier. Perhaps the key has something to do with it?'}
    },
    'treasure_room': {
        'name': 'Treasure Room',
        'description': 'You\'ve found it! A small chamber filled with glittering gold and ancient artifacts.',
        'exits': {'south': 'dark_cave'},
        'items': ['ancient_artifact'],
        'events': []
    }
}

items = {
    'old_map': {
        'name': 'Old Map',
        'description': 'A tattered map showing various paths and landmarks. It seems incomplete.',
        'can_pickup': True,
        'use_effect': 'You study the map. It seems to point towards a hidden path behind the waterfall.'
    },
    'healing_herb': {
        'name': 'Healing Herb',
        'description': 'A fragrant herb known for its restorative properties.',
        'can_pickup': True,
        'use_effect': 'You consume the healing herb. You feel a surge of vitality!'
    },
    'shiny_pebble': {
        'name': 'Shiny Pebble',
        'description': 'Just a smooth, shiny pebble. Nothing special.',
        'can_pickup': True,
        'use_effect': 'You toss the pebble. It bounces harmlessly.'
    },
    'rusty_key': {
        'name': 'Rusty Key',
        'description': 'An old, rusty key. It looks like it might open something important.',
        'can_pickup': True,
        'use_effect': 'You hold the key. It feels warm in your hand. Perhaps it unlocks the magical barrier at the cave entrance?'
    },
    'ancient_artifact': {
        'name': 'Ancient Artifact',
        'description': 'A glowing artifact of immense power. You feel a strange energy emanating from it.',
        'can_pickup': True,
        'use_effect': 'You hold the artifact. The air crackles around you. You feel incredibly powerful!'
    }
}

SAVE_FILE = 'adventure_save.json'

def save_game():
    try:
        with open(SAVE_FILE, 'w') as f:
            json.dump(game_state, f, indent=4)
        print("Game saved successfully!")
    except IOError:
        print("Error: Could not save game.")

def load_game():
    global game_state
    try:
        with open(SAVE_FILE, 'r') as f:
            game_state = json.load(f)
        print("Game loaded successfully!")
        display_location()
    except FileNotFoundError:
        print("No saved game found.")
    except json.JSONDecodeError:
        print("Error: Corrupted save file.")
    except IOError:
        print("Error: Could not load game.")

def display_location():
    current_loc_id = game_state['current_location']
    loc = locations[current_loc_id]
    print(f"\n--- {loc['name']} ---")
    print(loc['description'])

    if loc['items']:
        print("You see the following items:")
        for item_id in loc['items']:
            print(f"- {items[item_id]['name']}")

    print("Exits:")
    for direction, target_loc_id in loc['exits'].items():
        print(f"- {direction.capitalize()} to {locations[target_loc_id]['name']}")
    print(f"Health: {game_state['health']}/100")

def move_player(direction):
    current_loc_id = game_state['current_location']
    loc = locations[current_loc_id]
    
    if direction in loc['exits']:
        target_loc_id = loc['exits'][direction]
        target_loc = locations[target_loc_id]

        if 'requires_flag' in target_loc:
            flag_name = target_loc['requires_flag']['flag']
            flag_value = target_loc['requires_flag']['value']
            message = target_loc['requires_flag']['message']
            if game_state['flags'].get(flag_name) != flag_value:
                print(message)
                return

        game_state['current_location'] = target_loc_id
        print(f"You go {direction}.")
        run_location_events()
        display_location()
    else:
        print("You can't go that way.")

def take_item(item_name):
    current_loc_id = game_state['current_location']
    loc = locations[current_loc_id]
    
    found_item_id = None
    for item_id in loc['items']:
        if items[item_id]['name'].lower() == item_name.lower():
            found_item_id = item_id
            break
    
    if found_item_id:
        if items[found_item_id]['can_pickup']:
            game_state['inventory'].append(found_item_id)
            loc['items'].remove(found_item_id)
            print(f"You pick up the {items[found_item_id]['name']}.")
            if found_item_id == 'rusty_key':
                game_state['flags']['has_key'] = True
        else:
            print(f"You can't pick up the {items[found_item_id]['name']}.")
    else:
        print(f"There's no {item_name} here.")

def use_item(item_name):
    found_item_id = None
    for item_id in game_state['inventory']:
        if items[item_id]['name'].lower() == item_name.lower():
            found_item_id = item_id
            break
    
    if found_item_id:
        item_data = items[found_item_id]
        print(item_data['use_effect'])
        
        if found_item_id == 'healing_herb':
            game_state['health'] = min(100, game_state['health'] + 25)
            game_state['inventory'].remove(found_item_id)
            print(f"Your health is now {game_state['health']}.")
        elif found_item_id == 'rusty_key':
            if game_state['current_location'] == 'cave_entrance':
                game_state['flags']['door_unlocked'] = True
                print("You insert the rusty key into a hidden slot near the cave entrance. There's a click, and the magical barrier shimmers and disappears!")
                game_state['inventory'].remove(found_item_id)
            else:
                print("There's nothing to use the key on here.")
        elif found_item_id == 'ancient_artifact':
            print("The artifact pulses with power. You feel invincible!")
        else:
            print("It doesn't seem to do anything special right now.")
    else:
        print(f"You don't have a {item_name}.")

def show_inventory():
    if not game_state['inventory']:
        print("Your inventory is empty.")
        return
    
    print("Your Inventory:")
    for item_id in game_state['inventory']:
        print(f"- {items[item_id]['name']}: {items[item_id]['description']}")

def look_around():
    display_location()

def run_location_events():
    current_loc_id = game_state['current_location']
    loc = locations[current_loc_id]

    if 'events' in loc:
        for event in loc['events']:
            if event == 'monster_encounter' and not game_state['flags']['monster_defeated']:
                print("\nA monstrous creature lunges from the shadows!")
                game_state['health'] -= 40
                print(f"You take damage! Your health is now {game_state['health']}.")
                if game_state['health'] <= 0:
                    print("You succumb to your wounds. Game Over!")
                    sys.exit()
                else:
                    print("You manage to fight it off, but it was a tough battle.")
                    game_state['flags']['monster_defeated'] = True

def handle_input(command):
    command_parts = command.lower().split(maxsplit=1)
    verb = command_parts[0]
    
    if verb == 'go':
        if len(command_parts) > 1:
            move_player(command_parts[1])
        else:
            print("Go where? (e.g., go north)")
    elif verb == 'take':
        if len(command_parts) > 1:
            take_item(command_parts[1])
        else:
            print("Take what? (e.g., take map)")
    elif verb == 'use':
        if len(command_parts) > 1:
            use_item(command_parts[1])
        else:
            print("Use what? (e.g., use herb)")
    elif verb == 'look':
        look_around()
    elif verb == 'inventory' or verb == 'inv':
        show_inventory()
    elif verb == 'save':
        save_game()
    elif verb == 'load':
        load_game()
    elif verb == 'quit' or verb == 'exit':
        print("Exiting game. Goodbye!")
        sys.exit()
    else:
        print("I don't understand that command.")

def game_loop():
    print("Welcome to the Python Text Adventure!")
    print("Type 'help' for a list of commands.")
    
    if os.path.exists(SAVE_FILE):
        print("A saved game exists. Do you want to load it? (yes/no)")
        choice = input("> ").lower()
        if choice == 'yes':
            load_game()
        else:
            display_location()
    else:
        display_location()

    while True:
        if game_state['current_location'] == 'treasure_room' and 'ancient_artifact' in game_state['inventory']:
            print("\nCongratulations! You have found the Ancient Artifact and completed your quest!")
            print("Thanks for playing!")
            break

        command = input("\nWhat do you do? ").strip()
        if not command:
            continue
        
        if command.lower() == 'help':
            print("\nAvailable commands:")
            print("- go [direction] (e.g., go north)")
            print("- take [item] (e.g., take map)")
            print("- use [item] (e.g., use herb)")
            print("- look (to see your surroundings again)")
            print("- inventory (or inv, to see what you're carrying)")
            print("- save (to save your progress)")
            print("- load (to load a saved game)")
            print("- quit (to exit the game)")
        else:
            handle_input(command)

        if game_state['health'] <= 0:
            print("Your health has dropped to zero. You have perished!")
            print("Game Over!")
            break

if __name__ == "__main__":
    game_loop()

# Additional implementation at 2025-08-04 08:21:03
import json
import sys
import os

# --- Game Data ---
rooms = {
    "antechamber": {
        "description": "You are in a dimly lit antechamber. A dusty old rug covers most of the floor.",
        "exits": {"north": "hallway"},
        "initial_items": ["rusty key"]
    },
    "hallway": {
        "description": "A long, narrow hallway. Cobwebs hang from the ceiling.",
        "exits": {"south": "antechamber", "east": "study", "west": "kitchen"},
        "initial_items": []
    },
    "study": {
        "description": "This study is filled with ancient books. A sturdy wooden door is to the north.",
        "exits": {"west": "hallway"},
        "locked_exit_info": {"north": {"room": "library", "key_needed": "ornate key", "state_key": "study_north_unlocked"}},
        "initial_items": []
    },
    "kitchen": {
        "description": "A grimy kitchen. A strange, gurgling sound comes from the corner.",
        "exits": {"east": "hallway"},
        "initial_items": ["apple"],
        "npc_info": {"name": "Goblin", "item_needed": "apple", "item_given": "ornate key", "dialogues": ["Grrr... hungry...", "You gave me food! Here, take this.", "Leave me alone, I'm eating."]},
        "enemy_info": {"name": "Slime", "initial_health": 30, "attack_damage": 5, "state_key": "slime_defeated"}
    },
    "library": {
        "description": "You've entered a vast library. Bookshelves reach to the high ceiling.",
        "exits": {"south": "study"},
        "initial_items": ["ancient scroll"]
    }
}

# --- Game State (dynamic, will be saved/loaded) ---
game_state = {}

def initialize_game_state():
    global game_state
    game_state = {
        "current_room": "antechamber",
        "inventory": [],
        "player_health": 100,
        "game_over": False,
        "win_condition_met": False,
        "turns": 0,
        "room_states": {}, # To store dynamic room data like items, enemy states, door states
        "npc_states": {},
        "enemy_healths": {} # To store current health of specific enemies
    }

    # Populate initial room states
    for room_name, room_data in rooms.items():
        game_state["room_states"][room_name] = {
            "items": list(room_data.get("initial_items", [])), # Make a copy
            "exits_unlocked": {} # For locked doors
        }
        if "locked_exit_info" in room_data:
            for exit_dir, info in room_data["locked_exit_info"].items():
                game_state["room_states"][room_name]["exits_unlocked"][exit_dir] = False

        if "npc_info" in room_data:
            npc_name = room_data["npc_info"]["name"]
            game_state["npc_states"][npc_name] = {
                "met": False,
                "item_given": False,
                "dialogue_index": 0
            }
        if "enemy_info" in room_data:
            enemy_name = room_data["enemy_info"]["name"]
            game_state["enemy_healths"][enemy_name] = room_data["enemy_info"]["initial_health"]
            game_state["room_states"][room_name][room_data["enemy_info"]["state_key"]] = False # False means not defeated

# --- Save/Load Functions ---
SAVE_FILE = "adventure_save.json"

def save_game():
    try:
        with open(SAVE_FILE, "w") as f:
            json.dump(game_state, f, indent=4)
        print("Game saved successfully!")
    except IOError:
        print("Error: Could not save game.")

def load_game():
    global game_state
    if not os.path.exists(SAVE_FILE):
        print("No saved game found. Starting new game.")
        initialize_game_state()
        return False
    try:
        with open(SAVE_FILE, "r") as f:
            game_state = json.load(f)
        print("Game loaded successfully!")
        return True
    except (IOError, json.JSONDecodeError):
        print("Error: Could not load game. Save file might be corrupted. Starting new game.")
        initialize_game_state()
        return False

# --- Game Logic Functions ---
def display_room():
    current_room_name = game_state["current_room"]
    room_data = rooms[current_room_name]
    room_state = game_