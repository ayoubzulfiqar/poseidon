package main

// Match checks if the text matches the pattern.
// It supports literal characters, '.' (any character), '*' (zero or more of the preceding character),
// '^' (start of text), and '$' (end of text).
func Match(pattern, text string) bool {
	// If the pattern starts with '^', it must match at the beginning of the text.
	if len(pattern) > 0 && pattern[0] == '^' {
		return matchHere(pattern[1:], text)
	}
	// Otherwise, try to match the pattern at any position in the text.
	// The loop includes len(text) to allow matching an empty pattern against an empty string
	// or a pattern like ".*" against an empty string.
	for i := 0; i <= len(text); i++ {
		if matchHere(pattern, text[i:]) {
			return true
		}
	}
	return false
}

// matchHere attempts to match the pattern starting from the beginning of the text.
func matchHere(pattern, text string) bool {
	// If pattern is empty, we've successfully matched everything.
	if len(pattern) == 0 {
		return true
	}

	// If the next character in pattern is '*' (e.g., 'c*').
	if len(pattern) >= 2 && pattern[1] == '*' {
		// Try to match zero or more occurrences of pattern[0]
		// followed by the rest of the pattern (pattern[2:]).
		return matchStar(pattern[0], pattern[2:], text)
	}

	// If the pattern starts with '$', it must match only if the text is exhausted.
	if len(pattern) > 0 && pattern[0] == '$' {
		return len(text) == 0
	}

	// If text is empty but pattern is not (and not a special case like '$' or '*'),
	// then no match.
	if len(text) == 0 {
		return false
	}

	// If the current character matches (either literal or '.')
	// then try to match the rest of the pattern and text.
	if pattern[0] == '.' || pattern[0] == text[0] {
		return matchHere(pattern[1:], text[1:])
	}

	// No match for the current character.
	return false
}

// matchStar handles the 'c*' pattern. It tries to match zero or more 'c's.
func matchStar(c byte, pattern, text string) bool {
	// Try to match the rest of the pattern (pattern) with the current text.
	// This covers the zero-occurrence case for 'c*'.
	for {
		if matchHere(pattern, text) {
			return true
		}
		// If text is exhausted or 'c' doesn't match the current text character,
		// then we can't consume more 'c's.
		if len(text) == 0 || (c != '.' && c != text[0]) {
			break
		}
		// Consume one character from text and try again.
		text = text[1:]
	}
	return false
}

// Additional implementation at 2025-08-04 08:31:51


// Additional implementation at 2025-08-04 08:32:48
