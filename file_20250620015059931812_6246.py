import random

def display_hangman(incorrect_guesses):
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
          /|\\  |
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
    print(stages[6 - incorrect_guesses])

def play_hangman(word_list):
    word = random.choice(word_list).upper()
    guessed_letters = set()
    incorrect_guesses = 0
    max_incorrect_guesses = 6

    print("Welcome to Hangman!")
    print("Try to guess the word.")

    while incorrect_guesses < max_incorrect_guesses:
        display_hangman(incorrect_guesses)

        display_word = ""
        for letter in word:
            if letter in guessed_letters:
                display_word += letter + " "
            else:
                display_word += "_ "
        print("\nWord: " + display_word)

        if "_" not in display_word:
            print("\nCongratulations! You guessed the word: " + word)
            break

        print("Guessed letters: " + ", ".join(sorted(list(guessed_letters))))
        print(f"Incorrect guesses remaining: {max_incorrect_guesses - incorrect_guesses}")

        guess = input("Guess a letter: ").upper()

        if not guess.isalpha() or len(guess) != 1:
            print("Invalid input. Please enter a single letter.")
            continue

        if guess in guessed_letters:
            print("You already guessed that letter. Try again.")
            continue

        guessed_letters.add(guess)

        if guess in word:
            print("Good guess!")
        else:
            print("That letter is not in the word.")
            incorrect_guesses += 1
    else:
        display_hangman(incorrect_guesses)
        print("\nGame Over! You ran out of guesses.")
        print("The word was: " + word)

if __name__ == "__main__":
    words = [
        "python", "programming", "computer", "science", "developer",
        "algorithm", "function", "variable", "keyboard", "monitor",
        "internet", "software", "hardware", "application", "database",
        "framework", "library", "module", "terminal", "debugger",
        "challenge", "hangman", "customizable", "functional", "complete"
    ]
    play_hangman(words)

# Additional implementation at 2025-06-20 01:52:06
