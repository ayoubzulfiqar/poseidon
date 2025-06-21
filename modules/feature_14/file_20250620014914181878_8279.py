import random

def choose_word(word_list):
    return random.choice(word_list).lower()

def display_word(word, guessed_letters):
    displayed = ""
    for letter in word:
        if letter in guessed_letters:
            displayed += letter
        else:
            displayed += "_"
    return displayed

def draw_hangman(lives):
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
    return stages[6 - lives]

def play_hangman(custom_word_list=None):
    if custom_word_list is None:
        word_list = ["python", "programming", "hangman", "computer", "keyboard", "developer", "algorithm", "function", "variable", "internet"]
    else:
        word_list = custom_word_list

    word = choose_word(word_list)
    guessed_letters = set()
    incorrect_guesses = set()
    lives = 6

    print("Welcome to Hangman!")
    print(draw_hangman(lives))
    print(f"The word has {len(word)} letters.")

    while lives > 0:
        current_display = display_word(word, guessed_letters)
        print(f"\nWord: {current_display}")
        print(f"Incorrect guesses: {', '.join(sorted(list(incorrect_guesses)))}")
        print(f"Lives remaining: {lives}")

        if "_" not in current_display:
            print("\nCongratulations! You guessed the word!")
            print(f"The word was: {word.upper()}")
            break

        guess = input("Guess a letter: ").lower()

        if not guess.isalpha() or len(guess) != 1:
            print("Invalid input. Please enter a single letter.")
        elif guess in guessed_letters or guess in incorrect_guesses:
            print(f"You already guessed '{guess}'. Try a different letter.")
        elif guess in word:
            print(f"Good guess! '{guess}' is in the word.")
            guessed_letters.add(guess)
        else:
            print(f"Sorry, '{guess}' is not in the word.")
            incorrect_guesses.add(guess)
            lives -= 1
            print(draw_hangman(lives))

    else:
        print(draw_hangman(lives))
        print("\nGame Over! You ran out of lives.")
        print(f"The word was: {word.upper()}")

    play_again = input("\nDo you want to play again? (yes/no): ").lower()
    if play_again == "yes":
        play_hangman(custom_word_list)

if __name__ == "__main__":
    play_hangman()

# Additional implementation at 2025-06-20 01:49:40
import random
import sys

HANGMAN_PICS = ['''
  +---+
  |   |
      |
      |
      |
      |
=========''', '''
  +---+
  |   |
  O   |
      |
      |
      |
=========''', '''
  +---+
  |   |
  O   |
  |   |
      |
      |
=========''', '''
  +---+
  |   |
  O   |
 /|   |
      |
      |
=========''', '''
  +---+
  |   |
  O   |
 /|\  |
      |
      |
=========''', '''
  +---+
  |   |
  O   |
 /|\  |
 /    |
      |
=========''', '''
  +---+
  |   |
  O   |
 /|\  |
 / \  |
      |
=========''']

WORD_LISTS = {
    "animals": "ant baboon badger bat bear beaver camel cat clam cobra cougar coyote crow deer dog donkey duck eagle ferret fox frog goat goose gorilla grizzly hamster hawk lion lizard llama mole monkey moose mouse mule newt otter owl panda parrot pigeon python rabbit ram rat raven rhino salmon seal shark sheep skunk sloth snake spider stork swan tiger toad trout turkey turtle weasel whale wolf wombat zebra".split(),
    "fruits": "apple banana cherry date fig grape kiwi lemon mango melon orange peach pear plum raspberry strawberry watermelon".split(),
    "countries": "usa canada mexico brazil argentina uk france germany italy spain china india japan australia egypt southafrica nigeria kenya russia turkey greece sweden norway finland denmark poland austria switzerland belgium netherlands portugal ireland newzealand".split(),
    "sports": "football basketball soccer tennis baseball golf hockey volleyball boxing wrestling swimming cycling athletics gymnastics".split()
}

def get_random_word(word_list):
    """Returns a random word from the passed list of words."""
    word_index = random.randint(0, len(word_list) - 1)
    return word_list[word_index]

def display_board(hangman_pics, missed_letters, correct_letters, secret_word):
    """Displays the current state of the Hangman game."""
    print(hangman_pics[len(missed_letters)])
    print()

    print('Missed letters:', end=' ')
    for letter in missed_letters:
        print(letter, end=' ')
    print()

    blanks = '_' * len(secret_word)

    for i in range(len(secret_word)): # Replace blanks with correctly guessed letters.
        if secret_word[i] in correct_letters:
            blanks = blanks[:i] + secret_word[i] + blanks[i+1:]

    for letter in blanks: # Show the secret word with spaces in between each letter.
        print(letter, end=' ')
    print()

def get_guess(already_guessed):
    """Ensures the player enters a single letter that hasn't been guessed before."""
    while True:
        print('Guess a letter.')
        guess = input().lower()
        if len(guess) != 1:
            print('Please enter a single letter.')
        elif guess not in 'abcdefghijklmnopqrstuvwxyz':
            print('Please enter a LETTER.')
        elif guess in already_guessed:
            print('You have already guessed that letter. Choose again.')
        else:
            return guess

def play_again():
    """Asks the player if they want to play again."""
    print('Do you want to play again? (yes or no)')
    return input().lower().startswith('y')

def select_word_category():
    """Allows the player to select a word category."""
    print("Available word categories:")
    categories = list(WORD_LISTS.keys())
    for i, category in enumerate(categories):
        print(f"{i+1}. {category.capitalize()}")

    while True:
        try:
            choice = int(input(f"Enter the number of your chosen category (1-{len(categories)}): "))
            if 1 <= choice <= len(categories):
                return WORD_LISTS[categories[choice-1]]
            else:
                print("Invalid choice. Please enter a number within the range.")
        except ValueError:
            print("Invalid input. Please enter a number.")

def main_game_loop():
    """Main function to run the Hangman game."""
    print('H A N G M A N')

    while True:
        chosen_word_list = select_word_category()
        secret_word = get_random_word(chosen_word_list)
        missed_letters = ''
        correct_letters = ''
        game_is_done = False

        while not game_is_done:
            display_board(HANGMAN_PICS, missed_letters, correct_letters, secret_word)

            guess = get_guess(missed_letters + correct_letters)

            if guess in secret_word:
                correct_letters += guess

                # Check if the player has won
                found_all_letters = True
                for i in range(len(secret_word)):
                    if secret_word[i] not in correct_letters:
                        found_all_letters = False
                        break
                if found_all_letters:
                    print(f'Yes! The secret word is "{secret_word}"! You have won!')
                    game_is_done = True
            else:
                missed_letters += guess

                # Check if player has guessed too many times
                if len(missed_letters) == len(HANGMAN_PICS) - 1:
                    display_board(HANGMAN_PICS, missed_letters, correct_letters, secret_word)
                    print(f'You have run out of guesses!\nThe word was "{secret_word}"')
                    game_is_done = True

        if not play_again():
            sys.exit()

main_game_loop()

# Additional implementation at 2025-06-20 01:50:15
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

DEFAULT_WORD_LISTS = {
    "animals": ["cat", "dog", "elephant", "giraffe", "zebra", "penguin", "octopus", "squirrel"],
    "fruits": ["apple", "banana", "orange", "grape", "strawberry", "blueberry", "pineapple", "mango"],
    "countries": ["canada", "brazil", "japan", "australia", "germany", "india", "egypt", "mexico"]
}

def load_word_list(filepath):
    """Loads words from a specified file, one word per line."""
    try:
        with open(filepath, 'r') as f:
            words = [word.strip().lower() for word in f if word.strip()]
        if not words:
            print(f"Warning: No valid words found in '{filepath}'.")
        return words
    except FileNotFoundError:
        print(f"Error: File '{filepath}' not found.")
        return []
    except Exception as e:
        print(f"An error occurred while reading the file: {e}")
        return []

def get_random_word(word_list):
    """Selects a random word from the given list."""
    import random
    if not word_list:
        return None
    return random.choice(word_list)

def display_game_state(incorrect_guesses, display_word, guessed_letters):
    """Prints the current state of the Hangman game."""
    print(HANGMAN_PICS[incorrect_guesses])
    print("\nWord: " + " ".join(display_word))
    print(f"Guessed letters: {', '.join(sorted(list(guessed_letters)))}")
    print(f"Lives left: {len(HANGMAN_PICS) - 1 - incorrect_guesses}")

def get_guess(guessed_letters):
    """Gets and validates a single letter guess from the player."""
    while True:
        guess = input("Guess a letter: ").lower()
        if len(guess) != 1:
            print("Please enter a single letter.")
        elif not guess.isalpha():
            print("Please enter a letter (a-z).")
        elif guess in guessed_letters:
            print("You already guessed that letter. Try again.")
        else:
            return guess

def play_hangman(word_list):
    """Plays a single round of Hangman."""
    if not word_list:
        print("No words available to play with. Please choose a valid word list.")
        return False # Indicate game could not be played

    word = get_random_word(word_list)
    if word is None:
        print("Could not select a word. Exiting game round.")
        return False

    display_word = ['_'] * len(word)
    guessed_letters = set()
    incorrect_guesses = 0
    max_incorrect_guesses = len(HANGMAN_PICS) - 1

    print("\n--- Starting New Game ---")
    print(f"The word has {len(word)} letters.")

    while True:
        display_game_state(incorrect_guesses, display_word, guessed_letters)

        if "_" not in display_word:
            print("\nCongratulations! You guessed the word:")
            print(" ".join(display_word))
            return True # Player won

        if incorrect_guesses >= max_incorrect_guesses:
            print(HANGMAN_PICS[max_incorrect_guesses])
            print("\nGame Over! You ran out of lives.")
            print(f"The word was: {word.upper()}")
            return False # Player lost

        guess = get_guess(guessed_letters)
        guessed_letters.add(guess)

        if guess in word:
            print(f"Good guess! '{guess}' is in the word.")
            for i, char in enumerate(word):
                if char == guess:
                    display_word[i] = guess
        else:
            print(f"Sorry, '{guess}' is not in the word.")
            incorrect_guesses += 1

def main():
    """Main function to run the Hangman game with customizable word lists and statistics."""
    total_wins = 0
    total_losses = 0

    while True:
        print("\n--- Hangman Game Menu ---")
        print("1. Play with default word categories")
        print("2. Load words from a file")
        print("3. View statistics")
        print("4. Exit")

        choice = input("Enter your choice (1-4): ")

        current_word_list = []
        if choice == '1':
            print("\nAvailable Categories:")
            for i, category in enumerate(DEFAULT_WORD_LISTS.keys()):
                print(f"{i+1}. {category.capitalize()}")
            
            category_choice = input("Enter category number or name: ").lower()
            
            selected_category = None
            if category_choice.isdigit():
                idx = int(category_choice) - 1
                if 0 <= idx < len(DEFAULT_WORD_LISTS):
                    selected_category = list(DEFAULT_WORD_LISTS.keys())[idx]
            else:
                if category_choice in DEFAULT_WORD_LISTS:
                    selected_category = category_choice
            
            if selected_category:
                current_word_list = DEFAULT_WORD_LISTS[selected_category]
                print(f"Selected category: {selected_category.capitalize()}")
            else:
                print("Invalid category choice. Using 'animals' by default.")
                current_word_list = DEFAULT_WORD_LISTS["animals"]

        elif choice == '2':
            filepath = input("Enter the path to your word list file (e.g., my_words.txt): ")
            current_word_list = load_word_list(filepath)
            if not current_word_list:
                print("Could not load words from file. Returning to main menu.")
                continue

        elif choice == '3':
            print(f"\n--- Game Statistics ---")
            print(f"Total Wins: {total_wins}")
            print(f"Total Losses: {total_losses}")
            continue

        elif choice == '4':
            print("Thanks for playing Hangman!")
            break

        else:
            print("Invalid choice. Please enter a number between 1 and 4.")
            continue

        if current_word_list:
            game_result = play_hangman(current_word_list)
            if game_result is True:
                total_wins += 1
            elif game_result is False:
                total_losses += 1
            # If game_result is None (no words), stats are not updated.

if __name__ == "__main__":
    main()