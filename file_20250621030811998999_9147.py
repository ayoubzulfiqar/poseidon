def create_text_histogram(text):
    if not isinstance(text, str):
        print("Input must be a string.")
        return

    char_counts = {}
    for char in text:
        char_counts[char] = char_counts.get(char, 0) + 1

    sorted_chars = sorted(char_counts.keys())

    for char in sorted_chars:
        count = char_counts[char]
        print(f"{char}: {'*' * count}")

if __name__ == "__main__":
    sample_text_1 = "hello world"
    create_text_histogram(sample_text_1)

    print("\n---")

    sample_text_2 = "Python Programming"
    create_text_histogram(sample_text_2)

    print("\n---")

    sample_text_3 = "aaaaabbc"
    create_text_histogram(sample_text_3)

    print("\n---")

    sample_text_4 = ""
    create_text_histogram(sample_text_4)

    print("\n---")

    sample_text_5 = "The quick brown fox jumps over the lazy dog."
    create_text_histogram(sample_text_5)

# Additional implementation at 2025-06-21 03:09:09
import collections
import string

def create_histogram(text, case_sensitive=False, include_whitespace=False, include_punctuation=False, include_numbers=False, sort_by_frequency=True, show_percentage=True, bar_chart=True, bar_char='#', bar_max_length=50):
    """
    Generates a character frequency histogram from the given text.

    Args:
        text (str): The input string to analyze.
        case_sensitive (bool): If True, 'a' and 'A' are counted separately. If False, they are treated as the same (converted to lowercase).
        include_whitespace (bool): If True, spaces, tabs, newlines are included in the count.
        include_punctuation (bool): If True, punctuation characters (as defined by string.punctuation) are included in the count.
        include_numbers (bool): If True, digit characters ('0'-'9') are included in the count.
        sort_by_frequency (bool): If True, characters are sorted by their frequency (descending). If False, by character (ascending).
        show_percentage (bool): If True, displays the percentage frequency alongside the count.
        bar_chart (bool): If True, displays a simple text-based bar chart.
        bar_char (str): The character to use for the bar chart.
        bar_max_length (int): The maximum length of the bar chart for the most frequent character.
    """
    if not text:
        print("Input text is empty.")
        return

    char_counts = collections.defaultdict(int)
    total_chars_counted = 0

    for char in text:
        original_char = char

        processed_char = char.lower() if not case_sensitive else char

        is_alpha = processed_char.isalpha()
        is_space = original_char.isspace()
        is_punct = original_char in string.punctuation
        is_digit = original_char.isdigit()

        if is_alpha:
            char_counts[processed_char] += 1
            total_chars_counted += 1
        elif include_whitespace and is_space:
            char_counts[original_char] += 1
            total_chars_counted += 1
        elif include_punctuation and is_punct:
            char_counts[original_char] += 1
            total_chars_counted += 1
        elif include_numbers and is_digit:
            char_counts[original_char] += 1
            total_chars_counted += 1

    if not char_counts:
        print("No countable characters found based on current settings.")
        return

    if sort_by_frequency:
        sorted_chars = sorted(char_counts.items(), key=lambda item: item[1], reverse=True)
    else:
        sorted_chars = sorted(char_counts.items(), key=lambda item: item[0])

    max_count = max(count for char, count in sorted_chars) if sorted_chars else 0

    print("\n--- Character Frequency Histogram ---")
    print(f"Total characters analyzed: {total_chars_counted}")
    print(f"Unique characters: {len(sorted_chars)}\n")

    for char, count in sorted_chars:
        display_char = repr(char) if char.isspace() or not char.isprintable() else char
        
        output_line = f"'{display_char}': {count:4d}"

        if show_percentage and total_chars_counted > 0:
            percentage = (count / total_chars_counted) * 100
            output_line += f" ({percentage:6.2f}%)"
        
        if bar_chart:
            if max_count > 0:
                bar_length = int((count / max_count) * bar_max_length)
                output_line += f" | {bar_char * bar_length}"
            else:
                output_line += " |"

        print(output_line)

if __name__ == "__main__":
    sample_text = """
    This is a sample text to demonstrate the character frequency histogram.
    It includes various characters, such as letters (both uppercase and lowercase),
    spaces, punctuation marks like periods, commas, exclamation points!
    Numbers like 123 and symbols like @#$% are also present.
    Let's see how it handles different cases.
    """

    print("--- Default Histogram (case-insensitive, alpha only, sorted by frequency, with percentage and bars) ---")
    create_histogram(sample_text)

    print("\n--- Case-sensitive, include all (whitespace, punctuation, numbers), sorted by character, no percentage, no bars ---")
    create_histogram(sample_text, case_sensitive=True, include_whitespace=True, include_punctuation=True, include_numbers=True, sort_by_frequency=False, show_percentage=False, bar_chart=False)

    print("\n--- Only whitespace, punctuation, and numbers, sorted by frequency, with percentage ---")
    create_histogram(sample_text, include_whitespace=True, include_punctuation=True, include_numbers=True, sort_by_frequency=True, show_percentage=True, bar_chart=True, bar_char='*')

    print("\n--- Short text, custom bar char, max bar length 20 ---")
    create_histogram("Python is fun and powerful!", bar_char='=', bar_max_length=20, include_punctuation=True)

    print("\n--- Empty text ---")
    create_histogram("")

    print("\n--- Text with only ignored characters (numbers, punctuation, whitespace ignored) ---")
    create_histogram("123!@#   ", include_numbers=False, include_punctuation=False, include_whitespace=False)

    print("\n--- Text with only numbers (counted) ---")
    create_histogram("1234567890", include_numbers=True)

    print("\n--- Text with mixed characters, focusing on specific inclusions ---")
    create_histogram("Hello World! 123 ABC", case_sensitive=False, include_whitespace=True, include_punctuation=True, include_numbers=True)

# Additional implementation at 2025-06-21 03:09:54
import collections

def create_histogram(text, case_sensitive=False, include_alphanumeric_only=False, sort_by_frequency=True, bar_character='*'):
    """
    Generates and prints a character frequency histogram from a given text.

    Args:
        text (str): The input text to analyze.
        case_sensitive (bool): If True, 'a' and 'A' are counted separately.
                               If False, they are treated as the same (converted to lowercase).
        include_alphanumeric_only (bool): If True, only counts letters and numbers.
                                           Punctuation, spaces, etc., are ignored.
        sort_by_frequency (bool): If True, characters are sorted by their frequency (descending).
                                  If False, characters are sorted alphabetically.
        bar_character (str): The character used to draw the histogram bars. Must be a single character.
    """
    if not isinstance(text, str):
        raise TypeError("Input 'text' must be a string.")
    if not isinstance(bar_character, str) or len(bar_character) != 1:
        raise ValueError("Input 'bar_character' must be a single character string.")

    char_counts = collections.defaultdict(int)

    for char in text:
        if include_alphanumeric_only and not char.isalnum():
            continue

        if not case_sensitive:
            char = char.lower()

        char_counts[char] += 1

    if not char_counts:
        print("No characters to display in histogram (text might be empty or filtered out).")
        return

    # Sort the characters based on specified criteria
    if sort_by_frequency:
        # Sort by count descending, then by character ascending for tie-breaking
        sorted_items = sorted(char_counts.items(), key=lambda item: (-item[1], item[0]))
    else:
        # Sort by character ascending
        sorted_items = sorted(char_counts.items(), key=lambda item: item[0])

    print("\n--- Character Frequency Histogram ---")
    # Determine max character representation length for alignment
    max_char_repr_len = 0
    for char, _ in sorted_items:
        # Use repr() for non-printable characters or spaces to make them visible
        display_char = repr(char) if not char.isprintable() or char == ' ' else char
        max_char_repr_len = max(max_char_repr_len, len(display_char))

    for char, count in sorted_items:
        display_char = repr(char) if not char.isprintable() or char == ' ' else char
        # Pad character representation for alignment
        print(f"'{display_char.ljust(max_char_repr_len)}': {count: <4} {bar_character * count}")
    print("-----------------------------------")

if __name__ == "__main__":
    sample_text1 = "Hello World! This is a test. 123 ABC."
    sample_text2 = "Python programming is fun and powerful."
    sample_text3 = "aaaaabbcdeff"
    sample_text4 = ""
    sample_text5 = "   " # Only spaces
    sample_text6 = "Mix of\nNewlines and\tTabs!"

    print("--- Example 1: Default (case-insensitive, all chars, freq sort) ---")
    create_histogram(sample_text1)

    print("\n--- Example 2: Case-sensitive, all chars, alpha sort ---")
    create_histogram(sample_text1, case_sensitive=True, sort_by_frequency=False)

    print("\n--- Example 3: Case-insensitive, alphanumeric only, freq sort ---")
    create_histogram(sample_text2, include_alphanumeric_only=True)

    print("\n--- Example 4: Case-sensitive, alphanumeric only, freq sort, custom bar ---")
    create_histogram(sample_text2, case_sensitive=True, include_alphanumeric_only=True, bar_character='#')

    print("\n--- Example 5: Simple text, default ---")
    create_histogram(sample_text3)

    print("\n--- Example 6: Empty string ---")
    create_histogram(sample_text4)

    print("\n--- Example 7: Only spaces, alphanumeric only (should be empty) ---")
    create_histogram(sample_text5, include_alphanumeric_only=True)

    print("\n--- Example 8: Only spaces, all chars ---")
    create_histogram(sample_text5, include_alphanumeric_only=False)

    print("\n--- Example 9: Text with newlines and tabs ---")
    create_histogram(sample_text6)

    # Example with error handling for invalid inputs
    try:
        print("\n--- Error Test: Invalid text type ---")
        create_histogram(123)
    except TypeError as e:
        print(f"Caught expected error: {e}")

    try:
        print("\n--- Error Test: Invalid bar_character length ---")
        create_histogram("test", bar_character="**")
    except ValueError as e:
        print(f"Caught expected error: {e}")