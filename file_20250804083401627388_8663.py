class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end_of_word = False

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end_of_word = True

    def _find_node(self, prefix: str) -> TrieNode:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return None
            node = node.children[char]
        return node

    def _collect_words(self, node: TrieNode, current_prefix: str, results: list) -> None:
        if node.is_end_of_word:
            results.append(current_prefix)

        for char, child_node in node.children.items():
            self._collect_words(child_node, current_prefix + char, results)

    def autocomplete(self, prefix: str) -> list[str]:
        prefix_node = self._find_node(prefix)
        if not prefix_node:
            return []

        results = []
        self._collect_words(prefix_node, prefix, results)
        return results

# Additional implementation at 2025-08-04 08:34:34
class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end_of_word = False

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end_of_word = True

    def search(self, word: str) -> bool:
        node = self.root
        for char in word:
            if char not in node.children:
                return False
            node = node.children[char]
        return node.is_end_of_word

    def starts_with(self, prefix: str) -> bool:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return False
            node = node.children[char]
        return True

    def _find_words_from_node(self, node: TrieNode, current_prefix: str, words: list) -> None:
        if node.is_end_of_word:
            words.append(current_prefix)

        for char, child_node in node.children.items():
            self._find_words_from_node(child_node, current_prefix + char, words)

    def autocomplete(self, prefix: str) -> list[str]:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return []

            node = node.children[char]

        suggestions = []
        self._find_words_from_node(node, prefix, suggestions)
        return suggestions