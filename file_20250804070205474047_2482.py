def start_game():
    print("You wake up in a dark, cold room. You have no memory of how you got here.")
    print("There are two doors in front of you: one on the 'left' and one on the 'right'.")
    choice = input("Which door do you choose? (left/right): ").lower().strip()
    if choice == "left":
        left_door_path()
    elif choice == "right":
        right_door_path()
    else:
        print("Invalid choice. You hesitate, and the darkness consumes you. Game Over.")

def left_door_path():
    print("\nYou open the left door and step into a long, dusty corridor.")
    print("A faint, unsettling sound echoes from further down the hall.")
    print("Do you 'investigate' the sound or 'continue' quietly past it?")
    choice = input("What do you do? (investigate/continue): ").lower().strip()
    if choice == "investigate":
        monster_lair()
    elif choice == "continue":
        library_room()
    else:
        print("Invalid choice. You stand frozen in indecision. Game Over.")

def right_door_path():
    print("\nYou open the right door and are greeted by a burst of light and the scent of fresh air.")
    print("You are in a vibrant, overgrown garden. A winding 'path' leads deeper, and a strange glowing 'plant' catches your eye.")
    print("Do you follow the 'path' or examine the 'plant'?")
    choice = input("What do you do? (path/plant): ").lower().strip()
    if choice == "path":
        ancient_portal()
    elif choice == "plant":
        magical_bloom()
    else:
        print("Invalid choice. You are overwhelmed by the beauty and lose your way. Game Over.")

def monster_lair():
    print("\nAs you investigate the sound, you stumble into a large cavern. A monstrous creature stirs in the shadows!")
    print("It roars and lunges towards you!")
    print("Do you 'fight' the monster or 'flee' back the way you came?")
    choice = input("What is your decision? (fight/flee): ").lower().strip()
    if choice == "fight":
        ending_hero()
    elif choice == "flee":
        ending_escape()
    else:
        print("Invalid choice. The monster devours you whole. Game Over.")

def library_room():
    print("\nYou quietly continue down the corridor and find a hidden door leading to a vast, ancient library.")
    print("Dusty tomes line the shelves. One particular 'book' glows faintly. You also see a small, dark 'exit' in the corner.")
    print("Do you 'read' the glowing book or 'search' for the exit?")
    choice = input("What do you do? (read/search): ").lower().strip()
    if choice == "read":
        ending_knowledge()
    elif choice == "search":
        ending_lost()
    else:
        print("Invalid choice. The library's silence becomes deafening, and you vanish without a trace. Game Over.")

def ancient_portal():
    print("\nYou follow the winding path through the garden. It leads to a clearing with a shimmering, ancient portal.")
    print("Strange symbols glow around its edges, beckoning you.")
    print("Do you 'enter' the portal or 'turn back' to the garden?")
    choice = input("What do you do? (enter/turn back): ").lower().strip()
    if choice == "enter":
        ending_new_world()
    elif choice == "turn back":
        ending_return_home()
    else:
        print("Invalid choice. The portal flickers and disappears, leaving you stranded. Game Over.")

def magical_bloom():
    print("\nYou approach the glowing plant. It's a beautiful, otherworldly bloom pulsating with soft light.")
    print("You feel a strange urge to 'touch' its petals, but also a sense of caution to 'observe' from afar.")
    print("Do you 'touch' the bloom or 'observe' it?")
    choice = input("What do you do? (touch/observe): ").lower().strip()
    if choice == "touch":
        ending_transformation()
    elif choice == "observe":
        ending_peaceful()
    else:
        print("Invalid choice. The plant's light intensifies, blinding you. Game Over.")

def ending_hero():
    print("\n--- Ending: The Hero's Triumph ---")
    print("You bravely fight the monster, using your wits and strength. After a grueling battle, you emerge victorious!")
    print("You find a hidden passage behind its lair, leading to freedom and a new life as a legendary hero.")
    print("Congratulations! You have found a heroic ending.")

def ending_escape():
    print("\n--- Ending: The Narrow Escape ---")
    print("You wisely choose to flee the monster. You run as fast as you can, finding a hidden crack in the wall.")
    print("You squeeze through, emerging into a vast, sunlit forest, leaving the horrors behind.")
    print("You have escaped, but the memory of the monster will forever haunt your dreams.")

def ending_knowledge():
    print("\n--- Ending: The Scholar's Enlightenment ---")
    print("You open the glowing book. Its pages are filled with ancient knowledge and forgotten spells.")
    print("You spend years in the library, absorbing wisdom, becoming a powerful sage. You never leave, but you find true purpose.")
    print("You have found enlightenment, but at the cost of your freedom.")

def ending_lost():
    print("\n--- Ending: Lost in the Labyrinth ---")
    print("You search frantically for the exit, but the library seems to shift around you.")
    print("Every path leads back to where you started. You are forever lost within its endless shelves, a prisoner of knowledge.")
    print("You are lost, with no hope of escape.")

def ending_new_world():
    print("\n--- Ending: A New Beginning ---")
    print("You step through the shimmering portal. The world on the other side is unlike anything you've ever seen.")
    print("A vibrant, alien landscape unfolds before you, full of new possibilities and adventures.")
    print("You have journeyed to a new world, ready to forge a new destiny.")

def ending_return_home():
    print("\n--- Ending: The Comfort of Familiarity ---")
    print("You decide against the portal, feeling a strong pull to return to what you know.")
    print("You retrace your steps through the garden, finding a hidden path that leads you back to your own world, safe and sound.")
    print("You have returned home, valuing the peace of your own world above all else.")

def ending_transformation():
    print("\n--- Ending: The Bloom's Embrace ---")
    print("You touch the glowing bloom. A surge of energy flows through you, transforming your very being.")
    print("You become one with the garden, a guardian spirit, forever bound to its beauty and magic.")
    print("You are no longer human, but something more, something ancient.")

def ending_peaceful():
    print("\n--- Ending: The Observer's Peace ---")
    print("You observe the magical bloom from a distance, appreciating its beauty without interfering.")
    print("The garden becomes your sanctuary. You live out your days in peaceful contemplation, finding serenity in nature's wonders.")
    print("You have found peace, living a quiet life in harmony with the magical garden.")

start_game()

# Additional implementation at 2025-08-04 07:07:48
import time
import sys

game_state = {
    "current_location": "start",
    "inventory": [],
    "player_name": "",
    "has_torch_lit": False,
    "has_met_elder": False,
    "has_key": False,
    "has_artifact": False,
    "game_over": False,
    "ending_message": ""
}

def print_slowly(text, delay=0.03):
    for char in text:
        sys.stdout.write(char)
        sys.stdout.flush()
        time.sleep(delay)
    sys.stdout.write("\n")

def add_item(item):
    if item not in game_state["inventory"]:
        game_state["inventory"].append(item)
        print_slowly(f"You picked up the {item}.")
    else:
        print_slowly(f"You already have the {item}.")

def remove_item(item):
    if item in game_state["inventory"]:
        game_state["inventory"].remove(item)
        print_slowly(f"You used the {item}.")
    else:
        print_slowly(f"You don't have the {item}.")

def has_item(item):
    return item in game_state["inventory"]

def display_inventory():
    if not game_state["inventory"]:
        print_slowly("Your inventory is empty.")
    else:
        print_slowly("Inventory:")
        for item in game_state["inventory"]:
            print_slowly(f"- {item}")

def get_player_choice(choices_dict):
    while True:
        for key, value in choices_dict.items():
            print_slowly(f"{key}. {value}")
        print_slowly("What do you do? (Type a number, or 'look', 'inventory', 'quit')")
        choice = input("> ").strip().lower()

        if choice == "look":
            scenes[game_state["current_location"]](re_describe=True)
            continue
        elif choice == "inventory":
            display_inventory()
            continue
        elif choice == "quit":
            game_state["game_over"] = True
            game_state["ending_message"] = "You decided to quit the adventure. Perhaps another time."
            return None
        elif choice.isdigit() and int(choice) in choices_dict:
            return int(choice)
        else:
            print_slowly("Invalid choice. Please try again.")

def scene_start(re_describe=False):
    if not re_describe:
        print_slowly("Welcome, adventurer, to the Whispering Woods!")
        game_state["player_name"] = input("What is your name? ").strip()
        print_slowly(f"Greetings, {game_state['player_name']}! You find yourself at the edge of a dense, ancient forest.")
        print_slowly("A narrow, overgrown path stretches before you, disappearing into the shadows.")
    
    choices = {
        1: "Follow the path into the forest.",
        2: "Look around the immediate area."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        game_state["current_location"] = "forest_path"
    elif choice == 2:
        print_slowly("You see tall, gnarled trees, their branches intertwining overhead. The air is cool and smells of damp earth and pine.")
        scene_start(re_describe=True)

def scene_forest_path(re_describe=False):
    if not re_describe:
        print_slowly("You walk deeper into the forest. The path becomes less distinct, and the trees grow thicker.")
        print_slowly("You come to a fork. To your left, the path seems to lead into a dark, shadowy area. To your right, it follows a faint sound of rushing water.")
    
    choices = {
        1: "Go left, towards the shadows.",
        2: "Go right, towards the sound of water.",
        3: "Go back to the start."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        game_state["current_location"] = "dark_cave_entrance"
    elif choice == 2:
        game_state["current_location"] = "river_bank"
    elif choice == 3:
        game_state["current_location"] = "start"

def scene_dark_cave_entrance(re_describe=False):
    if not re_describe:
        print_slowly("You approach a gaping maw in the earth – the entrance to a cave. It's pitch black inside, and a cold, damp air emanates from it.")
        if not has_item("torch"):
            print_slowly("Near the entrance, you notice a discarded, unlit torch.")
            add_item("torch")
    
    choices = {
        1: "Enter the cave.",
        2: "Go back to the forest path."
    }
    
    if has_item("torch") and not game_state["has_torch_lit"]:
        choices[3] = "Try to light the torch."

    choice = get_player_choice(choices)

    if choice == 1:
        if game_state["has_torch_lit"]:
            game_state["current_location"] = "deep_cave"
        else:
            print_slowly("It's too dark to enter without a light source. You might get lost or worse.")
            scene_dark_cave_entrance(re_describe=True)
    elif choice == 2:
        game_state["current_location"] = "forest_path"
    elif choice == 3 and has_item("torch") and not game_state["has_torch_lit"]:
        print_slowly("You rub the torch against a rough rock, and after a few attempts, a small flame flickers to life, casting dancing shadows.")
        game_state["has_torch_lit"] = True
        scene_dark_cave_entrance(re_describe=True)

def scene_deep_cave(re_describe=False):
    if not re_describe:
        print_slowly("With your lit torch, you venture deeper into the cave. The air grows heavy, and strange rock formations loom around you.")
        print_slowly("You hear a faint dripping sound and see a narrow passage ahead.")
        if not has_item("key"):
            print_slowly("You notice a glinting object on the ground – a small, ornate key.")
            add_item("key")
    
    choices = {
        1: "Continue through the narrow passage.",
        2: "Go back to the cave entrance."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        print_slowly("The passage leads to a dead end. You hear a low growl...")
        game_state["game_over"] = True
        game_state["ending_message"] = "You stumbled into a creature's lair and became its next meal. (Bad Ending: Lost in the Dark)"
    elif choice == 2:
        game_state["current_location"] = "dark_cave_entrance"

def scene_river_bank(re_describe=False):
    if not re_describe:
        print_slowly("You arrive at the bank of a wide, fast-flowing river. The water is murky, and the current looks strong.")
        print_slowly("On the opposite bank, you can barely make out what looks like a small settlement.")
    
    choices = {
        1: "Attempt to swim across the river.",
        2: "Follow the river downstream.",
        3: "Go back to the forest path."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        print_slowly("You bravely jump into the cold water. The current immediately pulls you under...")
        game_state["game_over"] = True
        game_state["ending_message"] = "The river's powerful current overwhelmed you. (Bad Ending: Drowned)"
    elif choice == 2:
        game_state["current_location"] = "village_outskirts"
    elif choice == 3:
        game_state["current_location"] = "forest_path"

def scene_village_outskirts(re_describe=False):
    if not re_describe:
        print_slowly("You follow the river for a while, and soon the faint settlement becomes a clear village.")
        print_slowly("Simple wooden houses are nestled among the trees. Smoke rises from a few chimneys.")
    
    choices = {
        1: "Enter the village.",
        2: "Stay hidden and observe.",
        3: "Go back to the river bank."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        game_state["current_location"] = "village_square"
    elif choice == 2:
        print_slowly("You observe for a while. The villagers seem peaceful, going about their daily lives. Nothing suspicious.")
        scene_village_outskirts(re_describe=True)
    elif choice == 3:
        game_state["current_location"] = "river_bank"

def scene_village_square(re_describe=False):
    if not re_describe:
        print_slowly("You enter the village square. A few villagers look up, but none seem alarmed. A friendly-looking elder approaches you.")
        if not game_state["has_met_elder"]:
            print_slowly(f"'Welcome, stranger,' says the elder. 'My name is Elara. What brings you to our humble village?'")
            game_state["has_met_elder"] = True
        else:
            print_slowly("Elara smiles. 'Welcome back, friend.'")
    
    choices = {
        1: "Ask about the village.",
        2: "Ask about any local legends or quests.",
        3: "Decide to settle down here."
    }
    
    if has_item("key"):
        choices[4] = "Show Elara the key you found."

    choice = get_player_choice(choices)

    if choice == 1:
        print_slowly("Elara tells you about the village's peaceful history, their reliance on the river, and their simple way of life.")
        scene_village_square(re_describe=True)
    elif choice == 2:
        print_slowly("Elara's eyes twinkle. 'There is an old tale... of a hidden chamber beneath the ancient oak, guarded by a riddle. It's said to hold a relic of great power.'")
        game_state["current_location"] = "ancient_oak"
    elif choice == 3:
        game_state["game_over"] = True
        game_state["ending_message"] = "You decided to abandon your adventurous life and settle down in the peaceful village. You lived a quiet, contented life. (Neutral Ending: A Simple Life)"
    elif choice == 4 and has_item("key"):
        print_slowly("You show Elara the ornate key. Her eyes widen.")
        print_slowly("'This! This is the Key of Whispers, said to unlock the chamber beneath the Ancient Oak! You must be destined for this quest!'")
        game_state["has_key"] = True
        scene_village_square(re_describe=True)

def scene_ancient_oak(re_describe=False):
    if not re_describe:
        print_slowly("You arrive at the Ancient Oak, a colossal tree whose branches seem to touch the sky. Its roots twist and turn, forming natural alcoves.")
        print_slowly("You notice a faint inscription on one of the largest roots.")
    
    choices = {
        1: "Examine the inscription.",
        2: "Go back to the village square."
    }
    
    if game_state["has_key"]:
        choices[3] = "Look for a lock or hidden entrance."

    choice = get_player_choice(choices)

    if choice == 1:
        print_slowly("The inscription reads: 'I have cities, but no houses; forests, but no trees; and water, but no fish. What am I?'")
        print_slowly("You ponder the riddle...")
        riddle_answer = input("Your answer: ").strip().lower()
        if riddle_answer == "a map":
            print_slowly("A low rumble echoes from beneath the oak. A section of the root system slides away, revealing a dark passage.")
            game_state["current_location"] = "hidden_chamber"
        else:
            print_slowly("Nothing happens. The inscription remains silent.")
            scene_ancient_oak(re_describe=True)
    elif choice == 2:
        game_state["current_location"] = "village_square"
    elif choice == 3 and game_state["has_key"]:
        print_slowly("You carefully examine the roots. Behind a thick curtain of ivy, you find a small, almost invisible keyhole.")
        print_slowly("You insert the Key of Whispers. With a soft click, a section of the root system slides away, revealing a dark passage.")
        game_state["current_location"] = "hidden_chamber"

def scene_hidden_chamber(re_describe=False):
    if not re_describe:
        print_slowly("You descend into the hidden chamber. It's surprisingly well-preserved, illuminated by a soft, ethereal glow.")
        print_slowly("In the center, on a stone pedestal, rests a shimmering artifact.")
        if not game_state["has_artifact"]:
            print_slowly("You reach out and take the Ancient Artifact.")
            add_item("Ancient Artifact")
            game_state["has_artifact"] = True
    
    choices = {
        1: "Take the artifact and return to the village.",
        2: "Examine the chamber further (nothing else here)."
    }
    choice = get_player_choice(choices)

    if choice == 1:
        game_state["game_over"] = True
        game_state["ending

# Additional implementation at 2025-08-04 07:09:13


# Additional implementation at 2025-08-04 07:10:46
game_state = {
    "current_room": "foyer",
    "inventory": [],
    "flags": {
        "has_read_book": False,
        "basement_lit": False
    }
}

rooms = {
    "foyer": {
        "description": "You find yourself in a dimly lit foyer. The front door is heavy and locked. A grand staircase leads upstairs, and a dark archway opens to your left.",
        "choices": {
            "go upstairs": {"next_room": "upstairs_landing"},
            "go through archway": {"next_room": "living_room"},
            "try front door": {"action": "check_front_door"}
        },
        "items_in_room": []
    },
    "upstairs_landing": {
        "description": "You are on a landing at the top of the stairs. There's a door to the north and another to the east. A faint light comes from the east.",
        "choices": {
            "go north": {"next_room": "study"},
            "go east": {"next_room": "bedroom"},
            "go downstairs": {"next_room": "foyer"}
        },
        "items_in_room": []
    },
    "living_room": {
        "description": "A spacious living room with dusty furniture. A fireplace dominates one wall. There's a small, ornate key on the mantelpiece. A door leads to the kitchen.",
        "choices": {
            "go to kitchen": {"next_room": "kitchen"},
            "go back to foyer": {"next_room": "foyer"}
        },
        "items_in_room": ["ornate key"]
    },
    "kitchen": {
        "description": "The kitchen is surprisingly clean, but cold. A back door is visible, and a rug in the corner looks like it might be hiding something.",
        "choices": {
            "open back door": {"next_room": "garden"},
            "look under rug": {"action": "reveal_trapdoor"},
            "go back to living room": {"next_room": "living_room"}
        },
        "items_in_room": []
    },
    "garden": {
        "description": "You are in an overgrown garden. A high wall surrounds it, and a rusty gate is to the east. There's a thick rope coiled near a broken fountain.",
        "choices": {
            "try garden gate": {"action": "check_garden_gate"},
            "go back to kitchen": {"next_room": "kitchen"}
        },
        "items_in_room": ["thick rope"]
    },
    "study": {
        "description": "A quiet study filled with books. A large desk sits in the center. On the desk, you see an old, leather-bound book and a folded map.",
        "choices": {
            "read old book": {"action": "read_book"},
            "go back to landing": {"next_room": "upstairs_landing"}
        },
        "items_in_room": ["old book", "folded map"]
    },
    "bedroom": {
        "description": "A dusty bedroom. The bed is unmade, and a small, flickering flashlight lies on the nightstand.",
        "choices": {
            "go back to landing": {"next_room": "upstairs_landing"}
        },
        "items_in_room": ["flashlight"]
    },
    "basement": {
        "description": "It's pitch black down here. You can hear scurrying noises. Without a light, you're completely lost and stumble into a deep pit.",
        "choices": {}, # No choices, leads to an ending if entered without light
        "items_in_room": []
    },
    "secret_passage": {
        "description": "Following the map, you found a hidden passage behind a loose brick in the study! It leads outside to a secluded forest path.",
        "choices": {}, # Leads to an ending
        "items_in_room": []
    }
}

endings = {
    "front_door_escape": "You used the ornate key to unlock the heavy front door! You step out into the crisp night air, finally free. (Ending 1/5: Freedom!)",
    "garden_escape": "You used the thick rope to scale the high garden wall! You land softly on the other side and run into the darkness. (Ending 2/5: Over the Wall!)",
    "basement_trap": "You descended into the dark basement without a light. You stumbled, fell into a deep pit, and couldn't get out. Your adventure ends here. (Ending 3/5: Trapped!)",
    "secret_escape": "Following the cryptic map, you discovered a hidden passage! It led you out of the house and into a secluded forest, far from the mysterious mansion. (Ending 4/5: The Secret Path!)",
    "give_up": "Overwhelmed by the house's mysteries, you simply sit down and wait. Your adventure ends here, in quiet despair. (Ending 5/5: Despair.)"
}

def display_room():
    room = rooms[game_state["current_room"]]
    print("\n" + "="*40)
    print(room["description"])
    if room["items_in_room"]:
        print(f"You see: {', '.join(room['items_in_room'])}.")
    print("="*40)

    print("\nWhat do you do?")
    for i, choice_text in enumerate(room["choices"]):
        print(f"{i+1}. {choice_text.capitalize()}")
    print("Type 'inventory' to see your items.")
    print("Type 'look' to re-describe the room.")
    print("Type 'quit' to give up.")

def process_input(command):
    global game_state
    current_room_data = rooms[game_state["current_room"]]

    if command == "look":
        display_room()
        return True
    elif command == "inventory":
        if game_state["inventory"]:
            print(f"Your inventory: {', '.join(game_state['inventory'])}.")
        else:
            print("Your inventory is empty.")
        return True
    elif command.startswith("take "):
        item_to_take = command[5:].strip().lower()
        if item_to_take in current_room_data["items_in_room"]:
            game_state["inventory"].append(item_to_take)
            current_room_data["items_in_room"].remove(item_to_take)
            print(f"You took the {item_to_take}.")

            if item_to_take == "folded map" and game_state["flags"]["has_read_book"]:
                if "follow map to secret passage" not in rooms["study"]["choices"]:
                    rooms["study"]["choices"]["follow map to secret passage"] = {"next_room": "secret_passage", "item_required": "folded map", "flag_required": "has_read_book"}
                    print("You now understand the map's cryptic markings!")
        else:
            print(f"There's no {item_to_take} here to take.")
        return True
    elif command == "quit":
        print(endings["give_up"])
        return False

    choices_list = list(current_room_data["choices"].keys())
    try:
        choice_index = int(command) - 1
        if 0 <= choice_index < len(choices_list):
            chosen_text = choices_list[choice_index]
            choice_data = current_room_data["choices"][chosen_text]

            if "item_required" in choice_data and choice_data["item_required"] not in game_state["inventory"]:
                print(f"You need the {choice_data['item_required']} to do that.")
                return True

            if "flag_required" in choice_data and not game_state["flags"].get(choice_data["flag_required"], False):
                print("You haven't discovered how to do that yet.")
                return True

            if "action" in choice_data:
                return handle_action(choice_data["action"])
            elif "next_room" in choice_data:
                game_state["current_room"] = choice_data["next_room"]
                return True
        else:
            print("Invalid choice number.")
    except ValueError:
        print("Invalid command or choice. Please enter a number or a command like 'take [item]', 'inventory', 'look', 'quit'.")
    return True

def handle_action(action_name):
    global game_state

    if action_name == "check_front_door":
        if "ornate key" in game_state["inventory"]:
            print(endings["front_door_escape"])
            return False
        else:
            print("The front door is locked tight. You need a key.")
            return True
    elif action_name == "reveal_trapdoor":
        if "open trapdoor" not in rooms["kitchen"]["choices"]:
            rooms["kitchen"]["choices"]["open trapdoor"] = {"action": "descend_basement"}
            print("You pull back the rug and reveal a hidden trapdoor!")
        else:
            print("The trapdoor is already revealed.")
        return True
    elif action_name == "check_garden_gate":
        if "thick rope" in game_state["inventory"]:
            print(endings["garden_escape"])
            return False
        else:
            print("The garden wall is too high to climb without assistance.")
            return True
    elif action_name == "read_book":
        if not game_state["flags"]["has_read_book"]:
            game_state["flags"]["has_read_book"] = True
            print("You open the old book. It's filled with cryptic symbols and a faint, almost invisible, drawing of a hidden passage.")
            print("You now understand the significance of the 'folded map'.")
            if "folded map" in game_state["inventory"]:
                if "follow map to secret passage" not in rooms["study"]["choices"]:
                    rooms["study"]["choices"]["follow map to secret passage"] = {"next_room": "secret_passage", "item_required": "folded map", "flag_required": "has_read_book"}
                    print("A new option appears!")
        else:
            print("You've already read this book. It's still cryptic.")
        return True
    elif action_name == "descend_basement":
        if "flashlight" not in game_state["inventory"]:
            print(endings["basement_trap"])
            return False
        else:
            game_state["current_room"] = "basement"
            rooms["basement"]["description"] = "The basement is dusty but navigable with your flashlight. You see old crates and a faint light coming from a small grate high on the wall."
            rooms["basement"]["choices"] = {
                "examine crates": {"action": "examine_crates"},
                "try to reach grate": {"action": "reach_grate