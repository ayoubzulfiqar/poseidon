def generate_acronym(phrase):
    """
    Generates an acronym from a given phrase.

    Args:
        phrase (str): The input phrase.

    Returns:
        str: The generated acronym.
    """
    words = phrase.split()
    acronym_letters = []
    for word in words:
        if word:
            acronym_letters.append(word[0].upper())
    return "".join(acronym_letters)

# Additional implementation at 2025-06-22 23:44:11
import re

def generate_acronym(phrase, exclude_words=None, ignore_case=True, separator='', make_uppercase=True):
    if not phrase:
        return ""

    if exclude_words is None:
        exclude_words_set = set()
    else:
        exclude_words_set = {word.lower() for word in exclude_words} if ignore_case else set(exclude_words)

    words = re.findall(r'\b\w+\b', phrase)

    acronym_letters = []
    for word in words:
        word_to_check = word.lower() if ignore_case else word

        if word_to_check and word_to_check not in exclude_words_set:
            acronym_letters.append(word[0])

    acronym = separator.join(acronym_letters)

    return acronym.upper() if make_uppercase else acronym

if __name__ == "__main__":
    print(generate_acronym("National Aeronautics and Space Administration"))
    print(generate_acronym("The United States of America", exclude_words=["the", "of"]))
    print(generate_acronym("Hyper Text Markup Language", separator='.'))
    print(generate_acronym("Frequently Asked Questions", exclude_words=["asked"]))
    print(generate_acronym("World Health Organization", make_uppercase=False))
    print(generate_acronym("Don't Forget To Bring A Towel!", exclude_words=["to", "a"]))
    print(generate_acronym(""))
    print(generate_acronym("The of a an", exclude_words=["the", "of", "a", "an"]))

# Additional implementation at 2025-06-22 23:45:18
import re

def generate_acronym(phrase, ignore_case=True, ignore_common_stopwords=True, custom_stopwords=None, min_word_length=1, remove_punctuation=True):
    """
    Generates an acronym from a given phrase with various customization options.

    Args:
        phrase (str): The input string from which to generate the acronym.
        ignore_case (bool): If True, treats all words as lowercase for processing
                            (the resulting acronym will be uppercase). Defaults to True.
        ignore_common_stopwords (bool): If True, ignores a predefined list of common
                                        English stop words. Defaults to True.
        custom_stopwords (list): An optional list of additional words to ignore.
                                 These are added to common stopwords if
                                 ignore_common_stopwords is True. Defaults to None.
        min_word_length (int): Only includes words with length greater than or
                               equal to this value. Defaults to 1.
        remove_punctuation (bool): If True, strips non-alphabetic characters from
                                   words before processing. Defaults to True.

    Returns:
        str: The generated acronym in uppercase.
    """

    common_stopwords = {
        "a", "an", "the", "and", "or", "but", "for", "nor", "so", "yet",
        "of", "in", "on", "at", "by", "with", "from", "about", "as", "to",
        "is", "are", "was", "were", "be", "been", "being", "has", "have", "had",
        "do", "does", "did", "not", "no", "yes", "it", "its", "itself", "he",
        "him", "his", "himself", "she", "her", "hers", "herself", "we", "us",
        "our", "ours", "ourselves", "you", "your", "yours", "yourself",
        "yourselves", "they", "them", "their", "theirs", "themselves", "this",
        "that", "these", "those", "i", "me", "my", "mine", "myself", "which",
        "what", "who", "whom", "whose", "where", "when", "why", "how", "all",
        "any", "both", "each", "few", "more", "most", "other", "some", "such",
        "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very",
        "s", "t", "can", "will", "just", "don", "should", "now"
    }

    stopwords_set = set()
    if ignore_common_stopwords:
        stopwords_set.update(common_stopwords)
    if custom_stopwords:
        if ignore_case:
            stopwords_set.update(word.lower() for word in custom_stopwords)
        else:
            stopwords_set.update(custom_stopwords)

    words = re.split(r'[\s-]+', phrase)

    acronym_letters = []
    for word in words:
        if remove_punctuation:
            cleaned_word = re.sub(r'[^a-zA-Z]', '', word)
        else:
            cleaned_word = word

        if not cleaned_word:
            continue

        processed_word = cleaned_word.lower() if ignore_case else cleaned_word

        if processed_word in stopwords_set:
            continue

        if len(cleaned_word) < min_word_length:
            continue

        acronym_letters.append(cleaned_word[0].upper())

    return "".join(acronym_letters)

if __name__ == "__main__":
    result1 = generate_acronym("North Atlantic Treaty Organization")
    print(result1)

    result2 = generate_acronym("The quick brown fox jumps over the lazy dog", ignore_common_stopwords=True)
    print(result2)

    result3 = generate_acronym("Self-Contained Underwater Breathing Apparatus", remove_punctuation=True)
    print(result3)

    result4 = generate_acronym("Artificial Intelligence", ignore_case=False)
    print(result4)

    result5 = generate_acronym("Frequently Asked Questions", min_word_length=3)
    print(result5)

    result6 = generate_acronym("HyperText Markup Language", custom_stopwords=["markup"])
    print(result6)

    result7 = generate_acronym("This is a Test Phrase", ignore_common_stopwords=False)
    print(result7)

    result8 = generate_acronym("A very long and complicated phrase for testing", min_word_length=4)
    print(result8)

    result9 = generate_acronym("Don't Forget About Hyphenated-Words", remove_punctuation=True)
    print(result9)

    result10 = generate_acronym("Project Management Institute", custom_stopwords=["institute"])
    print(result10)

# Additional implementation at 2025-06-22 23:45:59
import re

_COMMON_STOP_WORDS = {
    "a", "an", "the", "and", "or", "but", "nor", "for", "yet", "so",
    "of", "in", "on", "at", "by", "with", "from", "about", "as", "into",
    "through", "during", "before", "after", "above", "below", "to", "up",
    "down", "out", "off", "over", "under", "again", "further", "then",
    "once", "here", "there", "when", "where", "why", "how", "all", "any",
    "both", "each", "few", "more", "most", "other", "some", "such", "no",
    "not", "only", "own", "same", "too", "very", "s", "t", "can", "will",
    "just", "don", "should", "now", "d", "ll", "m", "o", "re", "ve", "y",
    "ain", "aren", "couldn", "didn", "doesn", "hadn", "hasn", "haven",
    "isn", "ma", "mightn", "mustn", "needn", "shan", "shouldn", "wasn",
    "weren", "won", "wouldn"
}

def generate_acronym(phrase, ignore_common_words=True, custom_ignore_words=None, include_numbers=False, max_length=None):
    if not isinstance(phrase, str) or not phrase:
        return ""

    ignore_words_set = set()
    if ignore_common_words:
        ignore_words_set.update(_COMMON_STOP_WORDS)
    if custom_ignore_words:
        ignore_words_set.update(word.lower() for word in custom_ignore_words)

    cleaned_phrase = re.sub(r'[^\w\s-]', '', phrase)
    words = cleaned_phrase.split()

    acronym_letters = []
    for word in words:
        if not word:
            continue

        lower_word = word.lower()
        if lower_word in ignore_words_set:
            continue

        sub_words = word.split('-')
        for sub_word in sub_words:
            if not sub_word:
                continue

            first_char = sub_word[0]
            if first_char.isalpha():
                acronym_letters.append(first_char.upper())
            elif include_numbers and first_char.isdigit():
                acronym_letters.append(first_char)

    acronym = "".join(acronym_letters)

    if max_length is not None and isinstance(max_length, int) and max_length >= 0:
        acronym = acronym[:max_length]

    return acronym