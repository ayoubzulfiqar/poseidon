package main

import (
	"fmt"
	"strconv"
	"strings"
)

func rleCompress(s string) string {
	if len(s) == 0 {
		return ""
	}

	var builder strings.Builder
	count := 1
	for i := 1; i <= len(s); i++ {
		if i == len(s) || s[i] != s[i-1] {
			builder.WriteString(strconv.Itoa(count))
			builder.WriteByte(s[i-1])
			count = 1
		} else {
			count++
		}
	}
	return builder.String()
}

func main() {
	testStrings := []string{
		"AAABBC",
		"WWWWWWWWWWWWBWWWWWWWWWWWWBBB",
		"A",
		"",
		"ABBC",
		"HHHHHHHHEEEEEELLLLLOOOOO",
		"GoGoGo",
		"111223",
	}

	for _, s := range testStrings {
		compressed := rleCompress(s)
		fmt.Printf("Original: \"%s\" -> Compressed: \"%s\"\n", s, compressed)
	}
}

// Additional implementation at 2025-06-18 00:29:36
