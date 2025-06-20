import hashlib
import zlib
import os

def _calculate_md5(filepath):
    hasher = hashlib.md5()
    with open(filepath, 'rb') as f:
        while True:
            chunk = f.read(4096)
            if not chunk:
                break
            hasher.update(chunk)
    return hasher.hexdigest()

def _calculate_sha(filepath, algorithm):
    hasher = None
    if algorithm == 'sha1':
        hasher = hashlib.sha1()
    elif algorithm == 'sha256':
        hasher = hashlib.sha256()
    elif algorithm == 'sha512':
        hasher = hashlib.sha512()
    else:
        return None

    with open(filepath, 'rb') as f:
        while True:
            chunk = f.read(4096)
            if not chunk:
                break
            hasher.update(chunk)
    return hasher.hexdigest()

def _calculate_crc32(filepath):
    crc = 0
    with open(filepath, 'rb') as f:
        while True:
            chunk = f.read(4096)
            if not chunk:
                break
            crc = zlib.crc32(chunk, crc)
    return f'{crc & 0xFFFFFFFF:08x}'

def validate_checksum(filepath, expected_checksum, checksum_type):
    if not os.path.exists(filepath):
        return False

    calculated_checksum = None
    checksum_type_lower = checksum_type.lower()

    if checksum_type_lower == 'md5':
        calculated_checksum = _calculate_md5(filepath)
    elif checksum_type_lower.startswith('sha'):
        calculated_checksum = _calculate_sha(filepath, checksum_type_lower)
    elif checksum_type_lower == 'crc32':
        calculated_checksum = _calculate_crc32(filepath)
    else:
        return False

    if calculated_checksum is None:
        return False

    return calculated_checksum == expected_checksum.lower()

# Additional implementation at 2025-06-20 01:20:17
import hashlib
import zlib
import os

class ChecksumValidator:
    def __init__(self):
        self._supported_hash_algorithms = {
            "md5": hashlib.md5,
            "sha1": hashlib.sha1,
            "sha256": hashlib.sha256,
            "sha512": hashlib.sha512,
        }
        self._supported_crc_algorithms = {
            "crc32": zlib.crc32,
        }
        self.supported_algorithms = list(self._supported_hash_algorithms.keys()) + list(self._supported_crc_algorithms.keys())

    def _get_hash_object(self, algorithm: str):
        algo_lower = algorithm.lower()
        if algo_lower in self._supported_hash_algorithms:
            return self._supported_hash_algorithms[algo_lower]()
        raise ValueError(f"Unsupported hash algorithm: {algorithm}")

    def _calculate_hash_checksum(self, data_bytes: bytes, algorithm: str) -> str:
        hasher = self._get_hash_object(algorithm)
        hasher.update(data_bytes)
        return hasher.hexdigest()

    def _calculate_crc_checksum(self, data_bytes: bytes, algorithm: str) -> str:
        algo_lower = algorithm.lower()
        if algo_lower == "crc32":
            crc_val = zlib.crc32(data_bytes) & 0xFFFFFFFF
            return f"{crc_val:08x}"
        raise ValueError(f"Unsupported CRC algorithm: {algorithm}")

    def calculate_string_checksum(self, data: str, algorithm: str) -> str:
        """
        Calculates the checksum for a given string using the specified algorithm.
        """
        data_bytes = data.encode('utf-8')
        algo_lower = algorithm.lower()
        if algo_lower in self._supported_hash_algorithms:
            return self._calculate_hash_checksum(data_bytes, algorithm)
        elif algo_lower in self._supported_crc_algorithms:
            return self._calculate_crc_checksum(data_bytes, algorithm)
        else:
            raise ValueError(f"Unsupported algorithm: {algorithm}. Supported algorithms are: {', '.join(self.supported_algorithms)}")

    def calculate_file_checksum(self, filepath: str, algorithm: str, chunk_size: int = 8192) -> str:
        """
        Calculates the checksum for a given file using the specified algorithm.
        Reads the file in chunks to handle large files efficiently.
        """
        if not os.path.exists(filepath):
            raise FileNotFoundError(f"File not found: {filepath}")

        algo_lower = algorithm.lower()

        if algo_lower in self._supported_hash_algorithms:
            hasher = self._get_hash_object(algorithm)
            with open(filepath, 'rb') as f:
                while chunk := f.read(chunk_size):
                    hasher.update(chunk)
            return hasher.hexdigest()
        elif algo_lower in self._supported_crc_algorithms:
            crc_val = 0
            with open(filepath, 'rb') as f:
                while chunk := f.read(chunk_size):
                    crc_val = zlib.crc32(chunk, crc_val)
            return f"{(crc_val & 0xFFFFFFFF):08x}"
        else:
            raise ValueError(f"Unsupported algorithm: {algorithm}. Supported algorithms are: {', '.join(self.supported_algorithms)}")

    def validate_string_checksum(self, data: str, expected_checksum: str, algorithm: str) -> bool:
        """
        Validates if the calculated checksum of a string matches the expected checksum.
        """
        calculated_checksum = self.calculate_string_checksum(data, algorithm)
        return calculated_checksum.lower() == expected_checksum.lower()

    def validate_file_checksum(self, filepath: str, expected_checksum: str, algorithm: str) -> bool:
        """
        Validates if the calculated checksum of a file matches the expected checksum.
        """
        calculated_checksum = self.calculate_file_checksum(filepath, algorithm)
        return calculated_checksum.lower() == expected_checksum.lower()