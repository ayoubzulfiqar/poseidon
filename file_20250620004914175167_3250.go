package main

import (
	"fmt"
	"sort"
)

type TrieNode struct {
	children    map[rune]*TrieNode
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

func (t *Trie) StartsWith(prefix string) []string {
	node := t.root
	for _, char := range prefix {
		if _, found := node.children[char]; !found {
			return []string{}
		}
		node = node.children[char]
	}

	var results []string
	t.collectWords(node, prefix, &results)
	sort.Strings(results)
	return results
}

func (t *Trie) collectWords(node *TrieNode, currentWord string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, currentWord)
	}

	for char, child := range node.children {
		t.collectWords(child, currentWord+string(char), results)
	}
}

func main() {
	trie := NewTrie()

	words := []string{"apple", "app", "apricot", "banana", "band", "apply", "cat", "car"}
	for _, word := range words {
		trie.Insert(word)
	}

	fmt.Println("Search 'apple':", trie.Search("apple"))
	fmt.Println("Search 'app':", trie.Search("app"))
	fmt.Println("Search 'apricot':", trie.Search("apricot"))
	fmt.Println("Search 'appl':", trie.Search("appl"))
	fmt.Println("Search 'orange':", trie.Search("orange"))
	fmt.Println("Search 'apply':", trie.Search("apply"))

	fmt.Println("\nWords starting with 'ap':", trie.StartsWith("ap"))
	fmt.Println("Words starting with 'app':", trie.StartsWith("app"))
	fmt.Println("Words starting with 'ban':", trie.StartsWith("ban"))
	fmt.Println("Words starting with 'b':", trie.StartsWith("b"))
	fmt.Println("Words starting with 'c':", trie.StartsWith("c"))
	fmt.Println("Words starting with 'z':", trie.StartsWith("z"))
	fmt.Println("Words starting with '':", trie.StartsWith(""))
}

// Additional implementation at 2025-06-20 00:49:59
package main

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a node in the Trie
type Node struct {
	children    map[rune]*Node
	isEndOfWord bool
}

// NewNode creates a new Trie node
func NewNode() *Node {
	return &Node{
		children: make(map[rune]*Node),
	}
}

// Trie represents the prefix tree
type Trie struct {
	root *Node
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

// Insert adds a word to the Trie
func (t *Trie) Insert(word string) {
	curr := t.root
	for _, char := range strings.ToLower(word) { // Store words in lowercase for case-insensitive search
		if _, ok := curr.children[char]; !ok {
			curr.children[char] = NewNode()
		}
		curr = curr.children[char]
	}
	curr.isEndOfWord = true
}

// SearchPrefix finds the node corresponding to the given prefix
func (t *Trie) SearchPrefix(prefix string) *Node {
	curr := t.root
	for _, char := range strings.ToLower(prefix) {
		if _, ok := curr.children[char]; !ok {
			return nil // Prefix not found
		}
		curr = curr.children[char]
	}
	return curr // Returns the node at the end of the prefix
}

// Autocomplete returns all words that start with the given prefix
func (t *Trie) Autocomplete(prefix string) []string {
	node := t.SearchPrefix(prefix)
	if node == nil {
		return []string{} // No words found for this prefix
	}

	var results []string
	t.collectWords(node, strings.ToLower(prefix), &results)
	sort.Strings(results) // Sort results alphabetically
	return results
}

// collectWords is a helper function to recursively collect all words from a given node
func (t *Trie) collectWords(node *Node, currentWord string, words *[]string) {
	if node.isEndOfWord {
		*words = append(*words, currentWord)
	}

	for char, child := range node.children {
		t.collectWords(child, currentWord+string(char), words)
	}
}

// Contains checks if a word exists in the Trie
func (t *Trie) Contains(word string) bool {
	node := t.SearchPrefix(word)
	return node != nil && node.isEndOfWord
}

// Delete removes a word from the Trie
// Returns true if the word was found and deleted, false otherwise.
func (t *Trie) Delete(word string) bool {
	path := make([]*Node, 0)
	curr := t.root
	for _, char := range strings.ToLower(word) {
		if _, ok := curr.children[char]; !ok {
			return false // Word not found
		}
		path = append(path, curr)
		curr = curr.children[char]
	}

	if !curr.isEndOfWord {
		return false // Word not found (only a prefix)
	}

	curr.isEndOfWord = false // Mark as not end of word

	// Backtrack and remove nodes if they are no longer part of any word
	for i := len(path) - 1; i >= 0; i-- {
		parent := path[i]
		char := rune(strings.ToLower(word)[i])
		child := parent.children[char]

		// If the child node is not an end of word and has no other children, delete it
		if !child.isEndOfWord && len(child.children) == 0 {
			delete(parent.children, char)
		} else {
			break // Stop if the node is still part of another word or has children
		}
	}
	return true
}

func main() {
	trie := NewTrie()

	wordsToInsert := []string{"apple", "apricot", "application", "apply", "banana", "band", "cat", "car", "cart", "dog", "door", "data"}
	for _, word := range wordsToInsert {
		trie.Insert(word)
	}

	fmt.Println("Autocomplete for 'ap':", trie.Autocomplete("ap"))
	fmt.Println("Autocomplete for 'app':", trie.Autocomplete("app"))
	fmt.Println("Autocomplete for 'b':", trie.Autocomplete("b"))
	fmt.Println("Autocomplete for 'ca':", trie.Autocomplete("ca"))
	fmt.Println("Autocomplete for 'd':", trie.Autocomplete("d"))
	fmt.Println("Autocomplete for 'xyz':", trie.Autocomplete("xyz"))
	fmt.Println("Autocomplete for 'a':", trie.Autocomplete("a"))

	fmt.Println("\nChecking word existence:")
	fmt.Println("Contains 'apple':", trie.Contains("apple"))
	fmt.Println("Contains 'app':", trie.Contains("app")) // "app" was not inserted as a full word
	fmt.Println("Contains 'banana':", trie.Contains("banana"))
	fmt.Println("Contains 'apricot':", trie.Contains("apricot"))
	fmt.Println("Contains 'appl':", trie.Contains("appl"))

	fmt.Println("\nDeleting words:")
	fmt.Println("Delete 'apple':", trie.Delete("apple"))
	fmt.Println("Contains 'apple' after deletion:", trie.Contains("apple"))
	fmt.Println("Autocomplete for 'ap' after 'apple' deletion:", trie.Autocomplete("ap"))

	fmt.Println("Delete 'app':", trie.Delete("app")) // Should return false as 'app' was not a full word
	fmt.Println("Autocomplete for 'app' after 'app' deletion attempt:", trie.Autocomplete("app"))

	fmt.Println("Delete 'application':", trie.Delete("application"))
	fmt.Println("Autocomplete for 'ap' after 'application' deletion:", trie.Autocomplete("ap"))

	fmt.Println("Delete 'data':", trie.Delete("data"))
	fmt.Println("Autocomplete for 'd' after 'data' deletion:", trie.Autocomplete("d"))

	fmt.Println("Delete 'dog':", trie.Delete("dog"))
	fmt.Println("Autocomplete for 'd' after 'dog' deletion:", trie.Autocomplete("d"))

	fmt.Println("Delete 'door':", trie.Delete("door"))
	fmt.Println("Autocomplete for 'd' after 'door' deletion:", trie.Autocomplete("d"))
}