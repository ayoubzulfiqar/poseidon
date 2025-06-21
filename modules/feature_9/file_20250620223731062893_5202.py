import collections
import heapq
class HuffmanNode:
    def __init__(self, char, freq):
        self.char = char
        self.freq = freq
        self.left = None
        self.right = None
    def __lt__(self, other):
        return self.freq < other.freq
def _build_frequency_table(text):
    return collections.Counter(text)
def _build_huffman_tree(freq_table):
    priority_queue = []
    for char, freq in freq_table.items():
        heapq.heappush(priority_queue, HuffmanNode(char, freq))
    if not priority_queue:
        return None
    if len(priority_queue) == 1:
        node = heapq.heappop(priority_queue)
        dummy_node = HuffmanNode(None, 0)
        merged = HuffmanNode(None, node.freq + dummy_node.freq)
        merged.left = node
        merged.right = dummy_node
        heapq.heappush(priority_queue, merged)
    while len(priority_queue) > 1:
        node1 = heapq.heappop(priority_queue)
        node2 = heapq.heappop(priority_queue)
        merged = HuffmanNode(None, node1.freq + node2.freq)
        merged.left = node1
        merged.right = node2
        heapq.heappush(priority_queue, merged)
    return heapq.heappop(priority_queue)
def _generate_huffman_codes(node, current_code, codes):
    if node is None:
        return
    if node.char is not None:
        codes[node.char] = current_code
        return
    _generate_huffman_codes(node.left, current_code + '0', codes)
    _generate_huffman_codes(node.right, current_code + '1', codes)
def _encode_text(text, codes):
    encoded_bits = []
    for char in text:
        encoded_bits.append(codes[char])
    return "".join(encoded_bits)
def _decode_bits(encoded_bits, tree_root):
    if tree_root is None:
        return ""
    decoded_text = []
    current_node = tree_root
    for bit in encoded_bits:
        if bit == '0':
            current_node = current_node.left
        else:
            current_node = current_node.right
        if current_node.char is not None:
            decoded_text.append(current_node.char)
            current_node = tree_root
    return "".join(decoded_text)
class HuffmanCompressorDecompressor:
    def compress(self, text):
        if not text:
            return "", None
        freq_table = _build_frequency_table(text)
        tree_root = _build_huffman_tree(freq_table)
        codes = {}
        _generate_huffman_codes(tree_root, "", codes)
        encoded_bits = _encode_text(text, codes)
        return encoded_bits, tree_root
    def decompress(self, encoded_bits, tree_root):
        if not encoded_bits or tree_root is None:
            return ""
        return _decode_bits(encoded_bits, tree_root)

# Additional implementation at 2025-06-20 22:38:18
import heapq
import collections
import struct
import sys

class HuffmanNode:
    def __init__(self, char, freq, left=None, right=None):
        self.char = char
        self.freq = freq
        self.left = left
        self.right = right

    def __lt__(self, other):
        return self.freq < other.freq

    def __eq__(self, other):
        return self.freq == other.freq

class HuffmanCompressor:
    def __init__(self):
        self.huffman_codes = {}
        self.huffman_tree_root = None

    def _calculate_frequencies(self, text):
        return collections.Counter(text)

    def _build_huffman_tree(self, frequencies):
        priority_queue = []
        for char, freq in frequencies.items():
            heapq.heappush(priority_queue, HuffmanNode(char, freq))

        if not priority_queue:
            return None
        if len(priority_queue) == 1:
            node = heapq.heappop(priority_queue)
            # Create a dummy root for single character to ensure a '0' code
            # This makes the _generate_huffman_codes logic consistent.
            root = HuffmanNode(None, node.freq, left=node)
            return root

        while len(priority_queue) > 1:
            left = heapq.heappop(priority_queue)
            right = heapq.heappop(priority_queue)
            merged_node = HuffmanNode(None, left.freq + right.freq, left, right)
            heapq.heappush(priority_queue, merged_node)

        return heapq.heappop(priority_queue)

    def _generate_huffman_codes(self, node, current_code=""):
        if node is None:
            return

        if node.char is not None:
            self.huffman_codes[node.char] = current_code
            return

        self._generate_huffman_codes(node.left, current_code + "0")
        self._generate_huffman_codes(node.right, current_code + "1")

    def _text_to_bit_string(self, text):
        encoded_bits = []
        for char in text:
            if char not in self.huffman_codes:
                raise ValueError(f"Character '{char}' not found in Huffman codes. Text might contain characters not present during tree building.")
            encoded_bits.append(self.huffman_codes[char])
        return "".join(encoded_bits)

    def _bit_string_to_bytes(self, bit_string):
        padding_needed = 8 - (len(bit_string) % 8)
        if padding_needed == 8:
            padding_needed = 0
        
        padded_bit_string = bit_string + '0' * padding_needed
        
        byte_array = bytearray()
        for i in range(0, len(padded_bit_string), 8):
            byte = int(padded_bit_string[i:i+8], 2)
            byte_array.append(byte)
            
        return bytes(byte_array), padding_needed

    def _bytes_to_bit_string(self, byte_array, padding_bits):
        bit_string = ""
        for byte in byte_array:
            bit_string += bin(byte)[2:].zfill(8)
        
        if padding_bits > 0:
            bit_string = bit_string[:-padding_bits]
            
        return bit_string

    def compress(self, text):
        if not text:
            return b'', {}, 0

        frequencies = self._calculate_frequencies(text)
        self.huffman_tree_root = self._build_huffman_tree(frequencies)
        self._generate_huffman_codes(self.huffman_tree_root)

        # Special handling for single character text: ensure it gets a code
        if len(frequencies) == 1 and not self.huffman_codes:
            char = list(frequencies.keys())[0]
            self.huffman_codes[char] = '0' # Assign '0' as code for single character

        bit_string = self._text_to_bit_string(text)
        compressed_bytes, padding_bits = self._bit_string_to_bytes(bit_string)
        
        return compressed_bytes, self.huffman_codes, padding_bits

    def decompress(self, compressed_bytes, huffman_codes, padding_bits):
        if not compressed_bytes and not huffman_codes:
            return ""

        # Rebuild the decoding tree from the provided Huffman codes
        decoding_tree_root = HuffmanNode(None, 0) 
        for char, code in huffman_codes.items():
            current_node = decoding_tree_root
            for bit in code:
                if bit == '0':
                    if current_node.left is None:
                        current_node.left = HuffmanNode(None, 0)
                    current_node = current_node.left
                else: # bit == '1'
                    if current_node.right is None:
                        current_node.right = HuffmanNode(None, 0)
                    current_node = current_node.right
            current_node.char = char

        bit_string = self._bytes_to_bit_string(compressed_bytes, padding_bits)
        
        decoded_text = []
        current_node = decoding_tree_root
        
        for bit in bit_string:
            if bit == '0':
                current_node = current_node.left
            else: # bit == '1'
                current_node = current_node.right

            if current_node.char is not None:
                decoded_text.append(current_node.char)
                current_node = decoding_tree_root

        return "".join(decoded_text)

original_text = "this is an example of a huffman coding text compressor and decompressor"
# original_text = "AAAAABBBCCDE"
# original_text = "A"
# original_text = ""

compressor = HuffmanCompressor()

compressed_data, codes, padding = compressor.compress(original_text)

print(f"Original text length: {len(original_text)} characters")
print(f"Original text size (approx): {sys.getsizeof(original_text)} bytes")
print(f"Compressed data length: {len(compressed_data)} bytes")
print(f"Compressed data size (approx): {sys.getsizeof(compressed_data)} bytes")
print(f"Huffman Codes: {codes}")
print(f"Padding bits: {padding}")

decompressed_text = compressor.decompress(compressed_data, codes, padding)

print(f"Decompressed text: {decompressed_text}")
print(f"Original and decompressed text match: {original_text == decompressed_text}")

print("\n--- Testing with a larger text ---")
long_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." * 10
compressor_long = HuffmanCompressor()
compressed_long_data, codes_long, padding_long = compressor_long.compress(long_text)

print(f"Original long text length: {len(long_text)} characters")
print(f"Original long text size (approx): {sys.getsizeof(long_text)} bytes")
print(f"Compressed long data length: {len(compressed_long_data)} bytes")
print(f"Compressed long data size (approx): {sys.getsizeof(compressed_long_data)} bytes")

decompressed_long_text = compressor_long.decompress(compressed_long_data, codes_long, padding_long)
print(f"Original and decompressed long text match: {long_text == decompressed_long_text}")

# Additional implementation at 2025-06-20 22:39:39
import heapq
import collections
import pickle
import struct

class Node:
    def __init__(self, char, freq, left=None, right=None):
        self.char = char
        self.freq = freq
        self.left = left
        self.right = right

    def __lt__(self, other):
        return self.freq < other.freq

class HuffmanCoder:
    def __init__(self):
        self.huffman_tree = None
        self.char_codes = {}
        self.reverse_char_codes = {}

    @staticmethod
    def _calculate_frequencies(text):
        return collections.Counter(text)

    def _build_huffman_tree(self, frequencies):
        priority_queue = [Node(char, freq) for char, freq in frequencies.items()]
        heapq.heapify(priority_queue)

        if len(priority_queue) == 1:
            self.huffman_tree = priority_queue[0]
            return self.huffman_tree

        while len(priority_queue) > 1:
            left = heapq.heappop(priority_queue)
            right = heapq.heappop(priority_queue)
            merged_node = Node(None, left.freq + right.freq, left, right)
            heapq.heappush(priority_queue, merged_node)

        self.huffman_tree = priority_queue[0] if priority_queue else None
        return self.huffman_tree

    def _generate_codes(self, node, current_code=""):
        if node is None:
            return

        if node.char is not None:
            if self.huffman_tree == node and not node.left and not node.right:
                self.char_codes[node.char] = "0"
                self.reverse_char_codes["0"] = node.char
            else:
                self.char_codes[node.char] = current_code
                self.reverse_char_codes[current_code] = node.char
            return

        self._generate_codes(node.left, current_code + "0")
        self._generate_codes(node.right, current_code + "1")

    def compress(self, text):
        if not text:
            return b'', None, 0

        frequencies = self._calculate_frequencies(text)
        if not frequencies:
            return b'', None, 0

        self._build_huffman_tree(frequencies)
        self.char_codes = {}
        self.reverse_char_codes = {}
        self._generate_codes(self.huffman_tree)

        encoded_bits = "".join(self.char_codes[char] for char in text)

        padding_bits = 8 - (len(encoded_bits) % 8)
        if padding_bits == 8:
            padding_bits = 0
        encoded_bits += '0' * padding_bits

        byte_array = bytearray()
        for i in range(0, len(encoded_bits), 8):
            byte = int(encoded_bits[i:i+8], 2)
            byte_array.append(byte)

        return bytes(byte_array), self.huffman_tree, padding_bits

    def decompress(self, encoded_bytes, huffman_tree, padding_bits):
        if not encoded_bytes:
            return ""

        bit_string = "".join(f'{byte:08b}' for byte in encoded_bytes)
        if padding_bits > 0:
            bit_string = bit_string[:-padding_bits]

        current_node = huffman_tree
        decoded_text = []

        for bit in bit_string:
            if bit == '0':
                current_node = current_node.left
            else:
                current_node = current_node.right

            if current_node.char is not None:
                decoded_text.append(current_node.char)
                current_node = huffman_tree

        return "".join(decoded_text)

    def compress_to_file(self, input_filepath, output_filepath):
        with open(input_filepath, 'r', encoding='utf-8') as f:
            text = f.read()

        encoded_bytes, tree, padding_bits = self.compress(text)

        with open(output_filepath, 'wb') as f:
            f.write(struct.pack('B', padding_bits))
            pickle.dump(tree, f)
            f.write(encoded_bytes)

    def decompress_from_file(self, input_filepath, output_filepath):
        with open(input_filepath, 'rb') as f:
            padding_bits = struct.unpack('B', f.read(1))[0]
            huffman_tree = pickle.load(f)
            encoded_bytes = f.read()

        decoded_text = self.decompress(encoded_bytes, huffman_tree, padding_bits)

        with open(output_filepath, 'w', encoding='utf-8') as f:
            f.write(decoded_text)

if __name__ == "__main__":
    import os

    coder = HuffmanCoder()

    original_text = "this is an example for huffman coding. this is a test text. a simple test to demonstrate the functionality."
    print("Original text length:", len(original_text), "characters")

    encoded_data, tree, padding = coder.compress(original_text)
    print("Encoded bytes length:", len(encoded_data), "bytes")
    print("Padding bits:", padding)

    decoded_text = coder.decompress(encoded_data, tree, padding)
    print("Decoded text length:", len(decoded_text), "characters")
    print("Decoded text matches original:", original_text == decoded_text)

    input_file = "input.txt"
    compressed_file = "compressed.bin"
    decompressed_file = "decompressed.txt"

    with open(input_file, "w", encoding="utf-8") as f:
        f.write(original_text)

    print(f"\nCompressing '{input_file}' to '{compressed_file}'...")
    coder.compress_to_file(input_file, compressed_file)
    print("Compression complete.")

    print(f"Decompressing '{compressed_file}' to '{decompressed_file}'...")
    coder.decompress_from_file(compressed_file, decompressed_file)
    print("Decompression complete.")

    with open(decompressed_file, "r", encoding="utf-8") as f:
        reconstructed_text = f.read()

    print("Reconstructed text matches original (from file):", original_text == reconstructed_text)

    os.remove(input_file)
    os.remove(compressed_file)
    os.remove(decompressed_file)

    print("\nTesting with single character text:")
    single_char_text = "aaaaa"
    encoded_single, tree_single, padding_single = coder.compress(single_char_text)
    decoded_single = coder.decompress(encoded_single, tree_single, padding_single)
    print(f"Original: '{single_char_text}', Decoded: '{decoded_single}'")
    print("Decoded single char text matches original:", single_char_text == decoded_single)

    print("\nTesting with empty text:")
    empty_text = ""
    encoded_empty, tree_empty, padding_empty = coder.compress(empty_text)
    decoded_empty = coder.decompress(encoded_empty, tree_empty, padding_empty)
    print(f"Original: '{empty_text}', Decoded: '{decoded_empty}'")
    print("Decoded empty text matches original:", empty_text == decoded_empty)