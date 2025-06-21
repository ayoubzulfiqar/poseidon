def generate_acronym(phrase):
    words = phrase.split()
    acronym = ""
    for word in words:
        if word:
            acronym += word[0].upper()
    return acronym

# Additional implementation at 2025-06-21 01:07:01
import re

def generate_acronym(phrase: str) -> str:
    """
    Generates an acronym from a given phrase.

    The function processes the input phrase by:
    1. Ignoring common short words (e.g., "a", "the", "and").
    2. Handling various delimiters (spaces, hyphens, underscores, punctuation)
       by splitting the phrase into words based on non-alphanumeric characters.
    3. Ensuring the resulting acronym is in uppercase.
    """
    if not phrase:
        return ""

    # Common words to ignore (case-insensitive)
    ignore_words = {
        "a", "an", "the", "and", "or", "of", "in", "on", "at", "for", "with",
        "is", "are", "was", "were", "be", "been", "being", "to", "from", "by",
        "as", "it", "its", "if", "but", "not", "no", "so", "up", "down", "out",
        "off", "over", "under", "through", "about", "against", "among", "around",
        "before", "behind", "below", "beside", "between", "beyond", "during",
        "except", "inside", "into", "near", "outside", "past", "since", "than",
        "until", "upon", "within", "without", "etc"
    }

    # Split the phrase by any sequence of characters that are not letters or numbers.
    # This effectively handles spaces, hyphens, underscores, and other punctuation
    # as word delimiters.
    words = re.split(r'[^a-zA-Z0-9]+', phrase)

    acronym_letters = []
    for word in words:
        # Clean and normalize the word: remove leading/trailing whitespace and convert to lowercase.
        cleaned_word = word.strip().lower()
        # If the cleaned word is not empty and not in the ignore list,
        # take its first letter.
        if cleaned_word and cleaned_word not in ignore_words:
            acronym_letters.append(cleaned_word[0])

    # Join the collected letters and convert the final acronym to uppercase.
    return "".join(acronym_letters).upper()