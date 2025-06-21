def run_length_encode(data):
    if not data:
        return []

    encoded_data = []
    current_char = data[0]
    count = 1

    for i in range(1, len(data)):
        if data[i] == current_char:
            count += 1
        else:
            encoded_data.append((current_char, count))
            current_char = data[i]
            count = 1
    
    encoded_data.append((current_char, count))
    return encoded_data

# Additional implementation at 2025-06-21 03:35:42
class RLECompressor:
    def compress(self, data):
        if not data:
            return []
        
        compressed_data = []
        count = 1
        
        for i in range(1, len(data)):
            if data[i] == data[i-1]:
                count += 1
            else:
                compressed_data.append((count, data[i-1]))
                count = 1
        
        compressed_data.append((count, data[-1]))
        
        return compressed_data

    def decompress(self, compressed_data):
        if not compressed_data:
            return "" # Default to empty string for common RLE use case
        
        decompressed_items = []
        for count, item in compressed_data:
            decompressed_items.extend([item] * count)
            
        all_single_char_strings = True
        for _, item in compressed_data:
            if not (isinstance(item, str) and len(item) == 1):
                all_single_char_strings = False
                break
        
        if all_single_char_strings:
            return "".join(decompressed_items)
        else:
            return decompressed_items

    def get_compression_ratio(self, original_data, compressed_data):
        original_length = len(original_data)
        
        if original_length == 0:
            return 1.0 
        
        compressed_length = len(compressed_data)
        
        return compressed_length / original_length

    def is_compressible(self, data):
        if not data or len(data) < 2:
            return False
        
        for i in range(1, len(data)):
            if data[i] == data[i-1]:
                return True
        return False

    def get_expanded_size(self, compressed_data):
        total_count = 0
        for count, _ in compressed_data:
            total_count += count
        return total_count

    def get_compression_savings(self, original_data, compressed_data):
        original_size = len(original_data)
        compressed_size_in_items = len(compressed_data)
        
        if original_size == 0:
            return 0.0
        
        return (original_size - compressed_size_in_items) / original_size

    def __call__(self, data):
        return self.compress(data)

# Additional implementation at 2025-06-21 03:36:58


# Additional implementation at 2025-06-21 03:38:25
class RLECompressor:
    def encode(self, data):
        if not data:
            return []

        encoded_data = []
        count = 1
        current_char = data[0]

        for i in range(1, len(data)):
            if data[i] == current_char:
                count += 1
            else:
                encoded_data.append((count, current_char))
                current_char = data[i]
                count = 1
        encoded_data.append((count, current_char))
        return encoded_data

    def decode(self, encoded_data):
        if not encoded_data:
            return ""

        decoded_string_parts = []
        for count, char in encoded_data:
            decoded_string_parts.append(char * count)
        return "".join(decoded_string_parts)