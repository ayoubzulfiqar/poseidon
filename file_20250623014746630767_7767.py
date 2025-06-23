import random
import string

def generate_password(
    length=12,
    use_lowercase=True,
    use_uppercase=True,
    use_digits=True,
    use_symbols=True,
    exclude_ambiguous=False,
    min_lowercase=0,
    min_uppercase=0,
    min_digits=0,
    min_symbols=0
):
    if length <= 0:
        raise ValueError("Password length must be a positive integer.")

    char_sets = []
    if use_lowercase:
        char_sets.append(string.ascii_lowercase)
    if use_uppercase:
        char_sets.append(string.ascii_uppercase)
    if use_digits:
        char_sets.append(string.digits)
    if use_symbols:
        char_sets.append(string.punctuation)

    if not char_sets:
        raise ValueError("At least one character type (lowercase, uppercase, digits, symbols) must be enabled.")

    all_characters_pool = "".join(char_sets)

    if exclude_ambiguous:
        ambiguous_chars = "lIO01"
        all_characters_pool = "".join(c for c in all_characters_pool if c not in ambiguous_chars)
        if not all_characters_pool:
             raise ValueError("No characters left after excluding ambiguous ones. Adjust complexity rules.")

    password_chars = []
    required_chars_count = 0

    if min_lowercase > 0:
        if not use_lowercase:
            raise ValueError("Cannot set min_lowercase if use_lowercase is False.")
        for _ in range(min_lowercase):
            password_chars.append(random.choice(string.ascii_lowercase))
        required_chars_count += min_lowercase

    if min_uppercase > 0:
        if not use_uppercase:
            raise ValueError("Cannot set min_uppercase if use_uppercase is False.")
        for _ in range(min_uppercase):
            password_chars.append(random.choice(string.ascii_uppercase))
        required_chars_count += min_uppercase

    if min_digits > 0:
        if not use_digits:
            raise ValueError("Cannot set min_digits if use_digits is False.")
        for _ in range(min_digits):
            password_chars.append(random.choice(string.digits))
        required_chars_count += min_digits

    if min_symbols > 0:
        if not use_symbols:
            raise ValueError("Cannot set min_symbols if use_symbols is False.")
        for _ in range(min_symbols):
            password_chars.append(random.choice(string.punctuation))
        required_chars_count += min_symbols

    if required_chars_count > length:
        raise ValueError("Sum of minimum character requirements exceeds total password length.")

    for _ in range(length - required_chars_count):
        password_chars.append(random.choice(all_characters_pool))

    random.shuffle(password_chars)

    return "".join(password_chars)

default_password = generate_password()
long_complex_password = generate_password(length=20, use_symbols=True, min_uppercase=3, min_digits=2, min_symbols=1)
numeric_only_password = generate_password(length=8, use_lowercase=False, use_uppercase=False, use_symbols=False, use_digits=True)
no_ambiguous_password = generate_password(length=15, exclude_ambiguous=True)
short_password = generate_password(length=6)

# Additional implementation at 2025-06-23 01:49:01
import random
import string

class PasswordGenerator:
    def __init__(self):
        self.default_length = 12
        self.default_include_upper = True
        self.default_include_lower = True
        self.default_include_digits = True
        self.default_include_symbols = True
        self.default_exclude_similar = False # e.g., l, 1, I, o, 0, O
        self.default_exclude_ambiguous = False # e.g., {}[]()/`~,;.<>
        self.generated_passwords = []

        # Define character sets for filtering
        self._lower_chars = set(string.ascii_lowercase)
        self._upper_chars = set(string.ascii_uppercase)
        self._digit_chars = set(string.digits)
        self._symbol_chars = set(string.punctuation)

        self.similar_chars = set("l1Io0O")
        self.ambiguous_chars = set("{}[]()/`~,;.<>") # Common ambiguous chars

    def _get_character_pool(self, include_upper, include_lower, include_digits, include_symbols, exclude_similar, exclude_ambiguous):
        char_pool_set = set()
        if include_lower:
            char_pool_set.update(self._lower_chars)
        if include_upper:
            char_pool_set.update(self._upper_chars)
        if include_digits:
            char_pool_set.update(self._digit_chars)
        if include_symbols:
            char_pool_set.update(self._symbol_chars)

        if not char_pool_set:
            raise ValueError("At least one character type (lowercase, uppercase, digits, symbols) must be included.")

        if exclude_similar:
            char_pool_set = char_pool_set - self.similar_chars
        if exclude_ambiguous:
            char_pool_set = char

# Additional implementation at 2025-06-23 01:50:23
import random
import string

def generate_password(length, include_lowercase, include_uppercase, include_digits, include_symbols):
    char_sets = []
    if include_lowercase:
        char_sets.append(string.ascii_lowercase)
    if include_uppercase:
        char_sets.append(string.ascii_uppercase)
    if include_digits:
        char_sets.append(string.digits)
    if include_symbols:
        char_sets.append(string.punctuation)

    if not char_sets:
        raise ValueError("At least one character type must be selected.")

    all_chars = "".join(char_sets)

    password_chars = []
    for char_set in char_sets:
        password_chars.append(random.choice(char_set))

    if length < len(password_chars):
        raise ValueError(f"Password length ({length}) is too short to include all selected character types ({len(password_chars)} required).")

    for _ in range(length - len(password_chars)):
        password_chars.append(random.choice(all_chars))

    random.shuffle(password_chars)

    return "".join(password_chars)

if __name__ == "__main__":
    print("Welcome to the Adjustable Password Generator!")

    while True:
        try:
            min_len_val = 8
            max_len_val = 128
            length_input = input(f"Enter desired password length (min {min_len_val}, max {max_len_val}): ")
            length = int(length_input)
            if not (min_len_val <= length <= max_len_val):
                print(f"Length must be between {min_len_val} and {max_len_val}.")
                continue
            break
        except ValueError:
            print("Invalid input. Please enter a number.")

    while True:
        use_lowercase = input("Include lowercase letters? (y/n): ").lower() == 'y'
        use_uppercase = input("Include uppercase letters? (y/n): ").lower() == 'y'
        use_digits = input("Include numbers? (y/n): ").lower() == 'y'
        use_symbols = input("Include symbols? (y/n): ").lower() == 'y'

        if not (use_lowercase or use_uppercase or use_digits or use_symbols):
            print("You must select at least one character type.")
        else:
            break

    while True:
        try:
            min_passwords = 1
            max_passwords = 10
            num_passwords_input = input(f"How many passwords do you want to generate? ({min_passwords}-{max_passwords}): ")
            num_passwords = int(num_passwords_input)
            if not (min_passwords <= num_passwords <= max_passwords):
                print(f"Please enter a number between {min_passwords} and {max_passwords}.")
                continue
            break
        except ValueError:
            print("Invalid input. Please enter a number.")

    print("\nGenerating passwords...")
    for i in range(num_passwords):
        try:
            password = generate_password(length, use_lowercase, use_uppercase, use_digits, use_symbols)
            print(f"Password {i+1}: {password}")
        except ValueError as e:
            print(f"Error generating password: {e}")
            break

# Additional implementation at 2025-06-23 01:50:58
import random
import string

class PasswordGenerator:
    def __init__(self):
        self.rules = {
            'length': 12,
            'include_uppercase': True,
            'include_lowercase': True,
            'include_digits': True,
            'include_symbols': True,
            'exclude_ambiguous': False, # e.g., l, 1, I, o, 0, O
            'no_consecutive_repeats': False,
        }
        self.char_sets = {
            'lowercase': string.ascii_lowercase,
            'uppercase': string.ascii_uppercase,
            'digits': string.digits,
            'symbols': string.punctuation,
        }
        self.ambiguous_chars = "lIO01" # Common ambiguous characters

    def set_rules(self, **kwargs):
        """
        Set or update the password generation rules.
        Example: generator.set_rules(length=16, include_symbols=False)
        """
        for key, value in kwargs.items():
            if key in self.rules:
                self.rules[key] = value
            else:
                print(f"Warning: Unknown rule '{key}'. Ignoring.")

    def _get_character_pool(self):
        """Constructs the character pool based on current rules."""
        pool = []
        required_pools = []

        if self.rules['include_lowercase']:
            pool.extend(list(self.char_sets['lowercase']))
            required_pools.append(list(self.char_sets['lowercase']))
        if self.rules['include_uppercase']:
            pool.extend(list(self.char_sets['uppercase']))
            required_pools.append(list(self.char_sets['uppercase']))
        if self.rules['include_digits']:
            pool.extend(list(self.char_sets['digits']))
            required_pools.append(list(self.char_sets['digits']))
        if self.rules['include_symbols']:
            pool.extend(list(self.char_sets['symbols']))
            required_pools.append(list(self.char_sets['symbols']))

        # Ensure at least one character type is selected
        if not pool:
            raise ValueError("At least one character type (lowercase, uppercase, digits, symbols) must be included.")

        # Remove duplicates from pool if any (though unlikely with string constants)
        pool = list(set(pool))

        if self.rules['exclude_ambiguous']:
            pool = [char for char in pool if char not in self.ambiguous_chars]
            # Filter required pools as well
            required_pools = [[char for char in rp if char not in self.ambiguous_chars] for rp in required_pools]
            # Remove empty required pools if all chars were ambiguous
            required_pools = [rp for rp in required_pools if rp]

        if not pool:
            raise ValueError("Character pool is empty after applying rules (e.g., excluding all available characters).")

        return pool, required_pools

    def generate_password(self):
        """Generates a single password based on the current rules."""
        length = self.rules['length']
        if length <= 0:
            raise ValueError("Password length must be greater than 0.")

        try:
            all_chars, required_pools = self._get_character_pool()
        except ValueError as e:
            return f"Error: {e}"

        password_chars = []
        
        # Ensure minimum length for guaranteed characters doesn't exceed total length
        if len(required_pools) > length:
            return "Error: Number of required character types exceeds the specified password length."

        # Guarantee at least one character from each required category
        # Make a copy of required_pools to modify
        temp_required_pools = [list(rp) for rp in required_pools]
        
        for pool in temp_required_pools:
            if pool: # Ensure pool is not empty after filtering
                char = random.choice(pool)
                password_chars.append(char)
        
        # Fill the remaining length
        last_char = None
        for _ in range(length - len(password_chars)):
            char = random.choice(all_chars)
            if self.rules['no_consecutive_repeats']:
                # Keep picking until it's not the same as the last char
                # This loop can potentially run forever if all_chars only has one character
                # or if the last char is the only option left.
                # We need to ensure there are other options.
                attempts = 0
                max_attempts = 100 # Prevent infinite loop for very small pools
                while char == last_char and len(all_chars) > 1 and attempts < max_attempts:
                    char = random.choice(all_chars)
                    attempts += 1
                if attempts == max_attempts:
                    # Fallback: if we can't find a non-repeating char, just use it.
                    # This makes the rule "best effort" for extreme cases.
                    pass 
            password_chars.append(char)
            last_char = char

        random.shuffle(password_chars) # Shuffle to randomize positions of guaranteed chars
        return "".join(password_chars)

    def generate_multiple_passwords(self, count):
        """Generates multiple passwords."""
        if not isinstance(count, int) or count <= 0:
            raise ValueError("Count must be a positive integer.")
        
        passwords = []
        for _ in range(count):
            passwords.append(self.generate_password())
        return passwords

    def check_password_strength(self, password):
        """
        Evaluates the strength of a given password.
        Returns a dictionary with score and qualitative assessment.
        """
        score = 0
        feedback = []

        # Length
        length = len(password)
        if length < 8:
            feedback.append("Too short (min 8 recommended).")
        elif length < 12:
            score += 1
            feedback.append("Good length, but longer is better.")
        else:
            score += 2
            feedback.append("Excellent length.")

        # Character types
        has_lower = any(c.islower() for c in password)
        has_upper = any(c.isupper() for c in password)
        has_digit = any(c.isdigit() for c in password)
        has_symbol = any(c in string.punctuation for c in password)

        char_types_count = sum([has_lower, has_upper, has_digit, has_symbol])

        if has_lower: score += 1
        if has_upper: score += 1
        if has_digit: score += 1
        if has_symbol: score += 1

        if char_types_count < 3:
            feedback.append("Use a mix of character types (lowercase, uppercase, digits, symbols).")
        elif char_types_count == 3:
            feedback.append("Good mix of character types.")
        else:
            feedback.append("Excellent mix of character types.")

        # Consecutive characters (simple check)
        for i in range(len(password) - 1):
            if password[i] == password[i+1]:
                score -= 1 # Penalize for consecutive repeats
                feedback.append("Avoid consecutive repeating characters.")
                break # Only penalize once for this

        # Qualitative assessment
        if score < 3:
            strength = "Very Weak"
        elif score < 5:
            strength = "Weak"
        elif score < 7:
            strength = "Moderate"
        elif score < 9:
            strength = "Strong"
        else:
            strength = "Very Strong"

        return {
            "score": score,
            "strength": strength,
            "feedback": feedback
        }

def main():
    generator = PasswordGenerator()
    
    while True:
        print("\n--- Password Generator Menu ---")
        print("1. Generate Password(s)")
        print("2. View Current Rules")
        print("3. Adjust Rules")
        print("4. Check Password Strength")
        print("5. Exit")
        
        choice = input("Enter your choice: ")
        
        if choice == '1':
            try:
                num_passwords_str = input("How many passwords to generate? (default 1): ")
                num_passwords = int(num_passwords_str) if num_passwords_str else 1
                
                if num_passwords <= 0:
                    print("Please enter a positive number.")
                    continue

                passwords = generator.generate_multiple_passwords(num_passwords)
                print("\nGenerated Passwords:")
                for i, pwd in enumerate(passwords):
                    print(f"{i+1}. {pwd}")
            except ValueError as e:
                print(f"Invalid input: {e}")
            except Exception as e:
                print(f"An error occurred: {e}")

        elif choice == '2':
            print("\n--- Current Password Generation Rules ---")
            for rule, value in generator.rules.items():
                print(f"{rule.replace('_', ' ').title()}: {value}")

        elif choice == '3':
            print("\n--- Adjust Rules ---")
            print("Enter 'y' for True, 'n' for False, or a number for length.")
            print("Press Enter to keep current value.")
            
            new_rules = {}
            
            try:
                length_str = input(f"Length (current: {generator.rules['length']}): ")
                if length_str:
                    new_rules['length'] = int(length_str)
                    if new_rules['length'] <= 0:
                        print("Length must be a positive integer. Not updated.")
                        del new_rules['length']
            except ValueError:
                print("Invalid length. Must be a number. Not updated.")

            bool_rules = [
                'include_uppercase', 'include_lowercase', 
                'include_digits', 'include_symbols', 
                'exclude_ambiguous', 'no_consecutive_repeats'
            ]
            
            for rule_name in bool_rules:
                current_value = generator.rules[rule_name]
                prompt = f"{rule_name.replace('_', ' ').title()} (current: {current_value} [y/n]): "
                user_input = input(prompt).lower()
                if user_input == 'y':
                    new_rules[rule_name] = True
                elif user_input == 'n':
                    new_rules[rule_name] = False
                elif user_input == '':
                    pass # Keep current value
                else:
                    print(f"Invalid input for {rule_name}. Keeping current value.")
            
            generator.set_rules(**new_rules)
            print("Rules updated.")

        elif choice == '4':
            password_to_check = input("Enter password to check strength: ")
            if not password_to_check:
                print("No password entered.")
                continue
            strength_info = generator.check_password_strength(password_to_check)
            print("\n--- Password Strength Report ---")
            print(f"Password: {password_to_check}")
            print(f"Strength: {strength_info['strength']} (Score: {strength_info['score']})")
            print("Feedback:")
            for feedback_item in strength_info['feedback']:
                print(f"- {feedback_item}")

        elif choice == '5':
            print("Exiting Password Generator. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()