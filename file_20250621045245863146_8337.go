package main

import (
	"fmt"
	"strings"
)

type TrieNode struct {
	children  map[rune]*TrieNode
	isEndOfWord bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, found := node.children[char]; !found {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[char]
	}
	node.isEndOfWord = true
}

func (t *Trie) SearchPrefix(prefix string) []string {
	node := t.root
	for _, char := range prefix {
		if _, found := node.children[char]; !found {
			return []string{} // Prefix not found
		}
		node = node.children[char]
	}

	var results []string
	t.collectWords(node, prefix, &results)
	return results
}

func (t *Trie) collectWords(node *TrieNode, currentWord string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, currentWord)
	}

	for char, childNode := range node.children {
		t.collectWords(childNode, currentWord+string(char), results)
	}
}

func main() {
	trie := NewTrie()

	words := []string{"apple", "app", "apricot", "apply", "banana", "band", "bat", "cat", "car", "cart"}
	for _, word := range words {
		trie.Insert(word)
	}

	fmt.Println("Words inserted:", strings.Join(words, ", "))
	fmt.Println()

	prefixes := []string{"ap", "ban", "a", "c", "z", "appl"}

	for _, prefix := range prefixes {
		completions := trie.SearchPrefix(prefix)
		fmt.Printf("Completions for '%s': %v\n", prefix, completions)
	}
}

// Additional implementation at 2025-06-21 04:53:37
package trie

type TrieNode struct {
	children map[rune]*TrieNode
	isEndOfWord bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEndOfWord: false,
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, found := node.children[char]; !found {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEndOfWord = true
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, char := range word {
		if _, found := node.children[char]; !found {
			return false
		}
		node = node.children[char]
	}
	return node.isEndOfWord
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for _, char := range prefix {
		if _, found := node.children[char]; !found {
			return false
		}
		node = node.children[char]
	}
	return true
}

func (t *Trie) Autocomplete(prefix string) []string {
	var suggestions []string
	node := t.root

	for _, char := range prefix {
		if _, found := node.children[char]; !found {
			return []string{}
		}
		node = node.children[char]
	}

	t.collectWords(node, prefix, &suggestions)
	return suggestions
}

func (t *Trie) collectWords(node *TrieNode, currentWord string, suggestions *[]string) {
	if node.isEndOfWord {
		*suggestions = append(*suggestions, currentWord)
	}

	for char, childNode := range node.children {
		t.collectWords(childNode, currentWord+string(char), suggestions)
	}
}

// Additional implementation at 2025-06-21 04:54:45
package main

type Node struct {
	children map[rune]*Node
	isEndOfWord bool
}

type Trie struct {
	root *Node
}

func NewNode() *Node {
	return &Node{
		children: make(map[rune]*Node),
	}
}

func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, found := node.children[char]; !found {
			node.children[char] = NewNode()
		}
		node = node.children[char]
	}
	node.isEndOfWord = true
}

func (t *Trie) Autocomplete(prefix string) []string {
	node := t.root
	for _, char := range prefix {
		if _, found := node.children[char]; !found {
			return nil
		}
		node = node.children[char]
	}

	var results []string
	t.collectWords(node, prefix, &results)
	return results
}

func (t *Trie) collectWords(node *Node, currentWord string, words *[]string) {
	if node.isEndOfWord {
		*words = append(*words, currentWord)
	}

	for char, childNode := range node.children {
		t.collectWords(childNode, currentWord+string(char), words)
	}
}