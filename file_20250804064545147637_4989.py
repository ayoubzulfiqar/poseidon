import random

def get_hangman_drawing(incorrect_guesses):
    stages = [
        """
           -----
           |   |
               |
               |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
               |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
           |   |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|   |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|\\  |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|\\  |
          /    |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|\\  |
          / \\  |
               |
        ---------
        """
    ]
    return stages[incorrect_guesses]

def play_hangman(word_list=None):
    if word_list is None:
        word_list = ["PYTHON", "PROGRAMMING", "COMPUTER", "ALGORITHM", "DEVELOPER", "KEYBOARD", "MONITOR", "SOFTWARE", "HARDWARE", "INTERNET"]
    
    word = random.choice(word_list).upper()
    guessed_letters = set()
    incorrect_guesses = 0
    max_incorrect_guesses = 6

    print("Welcome to Hangman!")
    print("Try to guess the word.")

    while True:
        display_word = ""
        for letter in word:
            if letter in guessed_letters:
                display_word += letter + " "
            else:
                display_word += "_ "

        print(get_hangman_drawing(incorrect_guesses))
        print(f"Word: {display_word}")
        print(f"Guessed letters: {', '.join(sorted(list(guessed_letters)))}")
        print(f"Lives left: {max_incorrect_guesses - incorrect_guesses}")

        if "_" not in display_word:
            print("\nCongratulations! You guessed the word!")
            print(f"The word was: {word}")
            break

        if incorrect_guesses >= max_incorrect_guesses:
            print("\nYou ran out of lives! Game Over!")
            print(f"The word was: {word}")
            break

        guess = input("Guess a letter: ").upper()

        if not guess.isalpha() or len(guess) != 1:
            print("Invalid input. Please enter a single letter.")
            continue

        if guess in guessed_letters:
            print(f"You already guessed '{guess}'. Try again.")
            continue

        guessed_letters.add(guess)

        if guess in word:
            print(f"Good guess! '{guess}' is in the word.")
        else:
            print(f"Sorry, '{guess}' is not in the word.")
            incorrect_guesses += 1

if __name__ == "__main__":
    # Example 1: Play with default word list
    play_hangman()

    # Example 2: Play with a custom word list
    # custom_words = ["APPLE", "BANANA", "CHERRY", "DATE", "ELDERBERRY"]
    # play_hangman(custom_words)

# Additional implementation at 2025-08-04 06:46:32
import random
import os

HANGMAN_PICS = [
    """
  +---+
  |   |
      |
      |
      |
      |
=========""",
    """
  +---+
  |   |
  O   |
      |
      |
      |
=========""",
    """
  +---+
  |   |
  O   |
  |   |
      |
      |
=========""",
    """
  +---+
  |   |
  O   |
 /|   |
      |
      |
=========""",
    """
  +---+
  |   |
  O   |
 /|\\  |
      |
      |
=========""",
    """
  +---+
  |   |
  O   |
 /|\\  |
 /    |
      |
=========""",
    """
  +---+
  |   |
  O   |
 /|\\  |
 / \\  |
      |
========="""
]

DEFAULT_WORDS = [
    "python", "hangman", "programming", "computer", "keyboard",
    "monitor", "developer", "algorithm", "variable", "function",
    "string", "integer", "boolean", "loop", "condition",
    "module", "library", "framework", "database", "internet"
]

wins = 0
losses = 0

def load_word_list(filename):
    try:
        with open(filename, 'r') as f:
            words = [word.strip().lower() for word in f if word.strip().isalpha()]
        if not words:
            return DEFAULT_WORDS
        return words
    except FileNotFoundError:
        return DEFAULT_WORDS

def get_random_word(word_list):
    return random.choice(word_list)

def display_game_state(word, guessed_letters, incorrect_guesses_count):
    os.system('cls' if os.name == 'nt' else 'clear')
    print(HANGMAN_PICS[incorrect_guesses_count])
    print()

    display_word = ""
    for letter in word:
        if letter in guessed_letters:
            display_word += letter + " "
        else:
            display_word += "_ "
    print(display_word)
    print()

    incorrect_guesses = sorted([g for g in guessed_letters if g not in word])
    print("Incorrect guesses:", " ".join(incorrect_guesses))
    print()

def play_game(current_word_list):
    global wins, losses
    secret_word = get_random_word(current_word_list)
    guessed_letters = set()
    incorrect_guesses_count = 0
    max_incorrect_guesses = len(HANGMAN_PICS) - 1

    while True:
        display_game_state(secret_word, guessed_letters, incorrect_guesses_count)

        if incorrect_guesses_count >= max_incorrect_guesses:
            print("You lost! The word was:", secret_word)
            losses += 1
            break

        word_guessed = True
        for letter in secret_word:
            if letter not in guessed_letters:
                word_guessed = False
                break
        if word_guessed:
            print("Congratulations! You guessed the word:", secret_word)
            wins += 1
            break

        guess = input("Guess a letter: ").lower()

        if not guess.isalpha() or len(guess) != 1:
            print("Please enter a single letter.")
            continue

        if guess in guessed_letters:
            print("You already guessed that letter.")
            continue

        guessed_letters.add(guess)

        if guess not in secret_word:
            incorrect_guesses_count += 1

def main_menu():
    global wins, losses
    current_word_list = DEFAULT_WORDS

    while True:
        os.system('cls' if os.name == 'nt' else 'clear')
        print("Hangman Game")
        print("1. Play Game")
        print("2. Choose Word List")
        print("3. View Stats")
        print("4. Exit")
        print("Current word list source: " + ("Default" if current_word_list is DEFAULT_WORDS else "Custom File"))
        choice = input("Enter your choice: ")

        if choice == '1':
            play_game(current_word_list)
            input("Press Enter to continue...")
        elif choice == '2':
            filename = input("Enter filename for word list (e.g., words.txt): ")
            new_list = load_word_list(filename)
            if new_list is DEFAULT_WORDS:
                print("File not found or empty. Using default word list.")
            else:
                current_word_list = new_list
                print("Word list loaded successfully.")
            input("Press Enter to continue...")
        elif choice == '3':
            print(f"Games Won: {wins}")
            print(f"Games Lost: {losses}")
            input("Press Enter to continue...")
        elif choice == '4':
            print("Thanks for playing!")
            break
        else:
            print("Invalid choice. Please try again.")
            input("Press Enter to continue...")

if __name__ == "__main__":
    main_menu()

# Additional implementation at 2025-08-04 06:47:22
import random
import os

def clear_screen():
    # Clears the console screen for a cleaner display
    os.system('cls' if os.name == 'nt' else 'clear')

def display_hangman(lives):
    # ASCII art for the hangman stages
    stages = [
        """
           -----
           |   |
           O   |
          /|\\  |
          / \\  |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|\\  |
          /    |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|   |
          /    |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
          /|   |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
           |   |
               |
               |
        ---------
        """,
        """
           -----
           |   |
           O   |
               |
               |
               |
        ---------
        """,
        """
           -----
           |   |
               |
               |
               |
               |
        ---------
        """
    ]
    return stages[lives]

def choose_word(word_categories, category_name):
    # Selects a random word from the specified category
    if category_name in word_categories:
        return random.choice(word_categories[category_name]).upper()
    else:
        # Fallback if category is not found (should be handled by input validation)
        print(f"Error: Category '{category_name}' not found. Using a random word from all categories.")
        all_words = [word for sublist in word_categories.values() for word in sublist]
        return random.choice(all_words).upper()

def get_display_word(word, guessed_letters):
    # Creates the current display of the word (e.g., "_ A P P _ E")
    display = ""
    for letter in word:
        if letter in guessed_letters:
            display += letter + " "
        else:
            display += "_ "
    return display.strip()

def get_player_guess(guessed_letters_so_far):
    # Prompts player for a guess and validates input
    while True:
        guess = input("Guess a letter: ").upper()
        if len(guess) != 1 or not guess.isalpha():
            print("Invalid input. Please enter a single letter.")
        elif guess in guessed_letters_so_far:
            print(f"You've already guessed '{guess}'. Try a different letter.")
        else:
            return guess

def play_hangman():
    # Define customizable word lists by category
    word_lists = {
        "Animals": ["elephant", "giraffe", "kangaroo", "penguin", "zebra", "octopus", "squirrel", "chimpanzee"],
        "Fruits": ["apple", "banana", "cherry", "grape", "mango", "orange", "strawberry", "pineapple"],
        "Countries": ["canada", "brazil", "germany", "japan", "mexico", "australia", "india", "france", "egypt"],
        "Sports": ["football", "basketball", "tennis", "swimming", "volleyball", "cycling", "badminton", "baseball"],
        "Professions": ["doctor", "teacher", "engineer", "artist", "chef", "pilot", "firefighter", "programmer"]
    }

    print("Welcome to Hangman!")
    print("Choose a word category:")
    # Display categories to the user
    category_names = list(word_lists.keys())
    for i, category in enumerate(category_names):
        print(f"{i+1}. {category}")

    chosen_category = ""
    while chosen_category not in word_lists:
        category_input = input("Enter the number or name of the category: ").strip()
        if category_input.isdigit() and 1 <= int(category_input) <= len(category_names):
            chosen_category = category_names[int(category_input) - 1]
        elif category_input.title() in word_lists: # Allow case-insensitive category name
            chosen_category = category_input.title()
        else:
            print("Invalid category selection. Please try again.")

    word_to_guess = choose_word(word_lists, chosen_category)
    lives = 6 # Number of incorrect guesses allowed
    guessed_letters = set() # Stores all letters guessed by the player

    clear_screen() # Clear screen before starting the game loop

    while lives > 0:
        current_display = get_display_word(word_to_guess, guessed_letters)
        print(display_hangman(lives))
        print(f"Category: {chosen_category}")
        print(f"Word: {current_display}")
        print(f"Lives remaining: {lives}")
        print(f"Guessed letters: {', '.join(sorted(list(guessed_letters))) if guessed_letters else 'None'}")
        print("-" * 30)

        # Check for win condition
        if "_" not in current_display:
            print("\n" + "=" * 30)
            print(f"CONGRATULATIONS! You guessed the word: {word_to_guess}")
            print("YOU WIN!")
            print("=" * 30 + "\n")
            break

        guess = get_player_guess(guessed_letters)
        guessed_letters.add(guess) # Add the new guess to the set of guessed letters

        if guess in word_to_guess:
            print(f"Good guess! '{guess}' is in the word.")
        else:
            print(f"Sorry, '{guess}' is not in the word.")
            lives -= 1 # Decrement lives for incorrect guess
        
        clear_screen() # Clear screen for the next turn

    # Game over condition (ran out of lives)
    if lives == 0:
        print(display_hangman(lives))
        print("\n" + "=" * 30)
        print("GAME OVER! You ran out of lives.")
        print(f"The word was: {word_to_guess}")
        print("=" * 30 + "\n")

    # Ask to play again
    play_again = input("Do you want to play again? (yes/no): ").lower()
    if play_again == 'yes':
        play_hangman() # Restart the game
    else:
        print("Thanks for playing Hangman! Goodbye!")

# Start the game
play_hangman()

# Additional implementation at 2025-08-04 06:48:10
import random
import os

class HangmanGame:
    def __init__(self):
        self.word_lists = {
            "easy": ["apple", "banana", "grape", "melon", "peach", "lemon", "kiwi", "plum", "pear", "berry"],
            "medium": ["computer", "keyboard", "monitor", "software", "program", "internet", "network", "browser", "desktop", "laptop"],
            "hard": ["xylophone", "juxtapose", "quizzical", "rhythm", "sphinx", "awkward", "gossip", "haphazard", "jinx", "zigzag"],
            "animals": ["elephant", "giraffe", "hippopotamus", "kangaroo", "rhinoceros", "chimpanzee", "crocodile", "octopus", "penguin", "squirrel"],
            "countries": ["america", "brazil", "canada", "denmark", "egypt", "france", "germany", "india", "japan", "mexico"]
        }
        self.selected_word_list_name = None
        self.secret_word = ""
        self.guessed_letters = set()
        self.incorrect_guesses = 0
        self.max_incorrect_guesses = 7 # Corresponds to 7 stages of hangman
        self.display_word = []

    def _clear_screen(self):
        os.system('cls' if os.name == 'nt' else 'clear')

    def _get_hangman_stage(self):
        stages = [
            """
               -----
               |   |
               O   |
              /|\\  |
              / \\  |
                   |
            ---------
            """,
            """
               -----
               |   |
               O   |
              /|\\  |
              /    |
                   |
            ---------
            """,
            """
               -----
               |   |
               O   |
              /|   |
                   |
                   |
            ---------
            """,
            """
               -----
               |   |
               O   |
                   |
                   |
                   |
            ---------
            """,
            """
               -----
               |   |
                   |
                   |
                   |
                   |
            ---------
            """,
            """
               -----
                   |
                   |
                   |
                   |
                   |
            ---------
            """,
            """
               -----
                   |
                   |
                   |
                   |
                   |
            ---------
            """,
            """
               -----
                   |
                   |
                   |
                   |
                   |
            ---------
            """
        ]
        return stages[self.max_incorrect_guesses - self.incorrect_guesses]

    def _initialize_game(self):
        self._clear_screen()
        print("Welcome to Hangman!")
        print("Available word lists:")
        for i, name in enumerate(self.word_lists.keys()):
            print(f"  {i+1}. {name.capitalize()}")

        while True:
            try:
                choice = input("Enter the number of the word list you want to play with: ")
                choice_index = int(choice) - 1
                list_names = list(self.word_lists.keys())
                if 0 <= choice_index < len(list_names):
                    self.selected_word_list_name = list_names[choice_index]
                    break
                else:
                    print("Invalid choice. Please enter a valid number.")
            except ValueError:
                print("Invalid input. Please enter a number.")

        self.secret_word = random.choice(self.word_lists[self.selected_word_list_name]).upper()
        self.guessed_letters = set()
        self.incorrect_guesses = 0
        self.display_word = ["_" for _ in self.secret_word]
        self._clear_screen()

    def _display_game_state(self):
        self._clear_screen()
        print(self._get_hangman_stage())
        print("\nWord: " + " ".join(self.display_word))
        print(f"Guessed letters: {', '.join(sorted(list(self.guessed_letters)))}")
        print(f"Incorrect guesses remaining: {self.max_incorrect_guesses - self.incorrect_guesses}")

    def _get_player_guess(self):
        while True:
            guess = input("Guess a letter: ").upper()
            if not guess.isalpha():
                print("Invalid input. Please enter a letter.")
            elif len(guess) != 1:
                print("Invalid input. Please enter only one letter.")
            elif guess in self.guessed_letters:
                print(f"You already guessed '{guess}'. Try again.")
            else:
                self.guessed_letters.add(guess)
                return guess

    def _process_guess(self, guess):
        if guess in self.secret_word:
            print(f"Good guess! '{guess}' is in the word.")
            for i, letter in enumerate(self.secret_word):
                if letter == guess:
                    self.display_word[i] = guess
        else:
            print(f"Sorry, '{guess}' is not in the word.")
            self.incorrect_guesses += 1

    def _check_game_over(self):
        if "_" not in self.display_word:
            self._display_game_state()
            print("\nCongratulations! You guessed the word!")
            print(f"The word was: {self.secret_word}")
            return True
        elif self.incorrect_guesses >= self.max_incorrect_guesses:
            self._display_game_state()
            print("\nGame Over! You ran out of guesses.")
            print(f"The word was: {self.secret_word}")
            return True
        return False

    def play_game(self):
        while True:
            self._initialize_game()
            game_over = False
            while not game_over:
                self._display_game_state()
                guess = self._get_player_guess()
                self._process_guess(guess)
                game_over = self._check_game_over()

            play_again = input("\nDo you want to play again? (yes/no): ").lower()
            if play_again != 'yes':
                print("Thanks for playing Hangman!")
                break

if __name__ == "__main__":
    game = HangmanGame()
    game.play_game()