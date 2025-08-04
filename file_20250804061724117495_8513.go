package main

import "fmt"

// isMatch checks if the text matches the pattern.
func isMatch(pattern, text string) bool {
	if len(pattern) > 0 && pattern[0] == '^' {
		return matchHere(pattern[1:], text)
	}
	for i := 0; i <= len(text); i++ {
		if matchHere(pattern, text[i:]) {
			return true
		}
	}
	return false
}

// matchHere attempts to match the pattern starting at the beginning of the text.
func matchHere(pattern, text string) bool {
	if len(pattern) == 0 {
		return true
	}

	if len(pattern) == 1 && pattern[0] == '$' {
		return len(text) == 0
	}

	if len(pattern) >= 2 {
		switch pattern[1] {
		case '*':
			return matchStar(pattern[0], pattern[2:], text)
		case '+':
			return matchPlus(pattern[0], pattern[2:], text)
		case '?':
			return matchQuestion(pattern[0], pattern[2:], text)
		}
	}

	if len(text) > 0 && (pattern[0] == '.' || pattern[0] == text[0]) {
		return matchHere(pattern[1:], text[1:])
	}

	return false
}

// matchStar handles 'c*'
func matchStar(c byte, pattern, text string) bool {
	for i := 0; i <= len(text); i++ {
		if matchHere(pattern, text[i:]) {
			return true
		}
		if i == len(text) || (c != '.' && c != text[i]) {
			break
		}
	}
	return false
}

// matchPlus handles 'c+'
func matchPlus(c byte, pattern, text string) bool {
	if len(text) == 0 || (c != '.' && c != text[0]) {
		return false
	}
	return matchStar(c, pattern, text[1:])
}

// matchQuestion handles 'c?'
func matchQuestion(c byte, pattern, text string) bool {
	if matchHere(pattern, text) { // Match zero occurrences
		return true
	}
	if len(text) > 0 && (c == '.' || c == text[0]) { // Match one occurrence
		return matchHere(pattern, text[1:])
	}
	return false
}

func main() {
	fmt.Println("--- Literal Matches ---")
	fmt.Println("Match 'abc' in 'abc':", isMatch("abc", "abc"))
	fmt.Println("Match 'abc' in 'xabc':", isMatch("abc", "xabc"))
	fmt.Println("Match 'abc' in 'abx':", isMatch("abc", "abx"))
	fmt.Println("Match 'a' in 'b':", isMatch("a", "b"))
	fmt.Println("Match '' in 'abc':", isMatch("", "abc"))

	fmt.Println("\n--- Dot (.) Matches ---")
	fmt.Println("Match 'a.c' in 'abc':", isMatch("a.c", "abc"))
	fmt.Println("Match 'a.c' in 'axc':", isMatch("a.c", "axc"))
	fmt.Println("Match 'a.c' in 'ab':", isMatch("a.c", "ab"))
	fmt.Println("Match '...' in 'abc':", isMatch("...", "abc"))

	fmt.Println("\n--- Star (*) Matches ---")
	fmt.Println("Match 'a*b' in 'b':", isMatch("a*b", "b"))
	fmt.Println("Match 'a*b' in 'ab':", isMatch("a*b", "ab"))
	fmt.Println("Match 'a*b' in 'aaab':", isMatch("a*b", "aaab"))
	fmt.Println("Match '.*' in 'abc':", isMatch(".*", "abc"))
	fmt.Println("Match 'a.*b' in 'axbyb':", isMatch("a.*b", "axbyb"))
	fmt.Println("Match 'a*b' in 'acb':", isMatch("a*b", "acb"))

	fmt.Println("\n--- Plus (+) Matches ---")
	fmt.Println("Match 'a+b' in 'b':", isMatch("a+b", "b"))
	fmt.Println("Match 'a+b' in 'ab':", isMatch("a+b", "ab"))
	fmt.Println("Match 'a+b' in 'aaab':", isMatch("a+b", "aaab"))
	fmt.Println("Match '.+' in 'abc':", isMatch(".+", "abc"))
	fmt.Println("Match '.+' in '':", isMatch(".+", ""))
	fmt.Println("Match 'a.+b' in 'axbyb':", isMatch("a.+b", "axbyb"))

	fmt.Println("\n--- Question (?) Matches ---")
	fmt.Println("Match 'a?b' in 'b':", isMatch("a?b", "b"))
	fmt.Println("Match 'a?b' in 'ab':", isMatch("a?b", "ab"))
	fmt.Println("Match 'a?b' in 'aaab':", isMatch("a?b", "aaab"))
	fmt.Println("Match '.?' in 'a':", isMatch(".?", "a"))
	fmt.Println("Match '.?' in '':", isMatch(".?", ""))

	fmt.Println("\n--- Anchors (^) ($) Matches ---")
	fmt.Println("Match '^abc' in 'abc':", isMatch("^abc", "abc"))
	fmt.Println("Match '^abc' in 'xabc':", isMatch("^abc", "xabc"))
	fmt.Println("Match 'abc$' in 'abc':", isMatch("abc$", "abc"))
	fmt.Println("Match 'abc$' in 'abcx':", isMatch("abc$", "abcx"))
	fmt.Println("Match '^a.c$' in 'axc':", isMatch("^a.c$", "axc"))
	fmt.Println("Match '^a.c$' in 'axcx':", isMatch("^a.c$", "axcx"))
	fmt.Println("Match '^a.c$' in 'xaxc':", isMatch("^a.c$", "xaxc"))
	fmt.Println("Match '^$' in '':", isMatch("^$", ""))
	fmt.Println("Match '^$' in 'a':", isMatch("^$", "a"))

	fmt.Println("\n--- Combined Patterns ---")
	fmt.Println("Match 'a.*c' in 'abxxc':", isMatch("a.*c", "abxxc"))
	fmt.Println("Match 'a.+c' in 'ac':", isMatch("a.+c", "ac"))
	fmt.Println("Match 'a.+c' in 'axc':", isMatch("a.+c", "axc"))
	fmt.Println("Match 'a?b?c?' in 'abc':", isMatch("a?b?c?", "abc"))
	fmt.Println("Match 'a?b?c?' in 'ac':", isMatch("a?b?c?", "ac"))
	fmt.Println("Match 'a?b?c?' in 'x':", isMatch("a?b?c?", "x"))
	fmt.Println("Match 'a?b?c?' in '':", isMatch("a?b?c?", ""))
}

// Additional implementation at 2025-08-04 06:18:29


// Additional implementation at 2025-08-04 06:19:34


// Additional implementation at 2025-08-04 06:21:08
