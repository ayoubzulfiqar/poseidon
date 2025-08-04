import hashlib
import zlib

def _to_bytes(data):
    if isinstance(data, str):
        return data.encode('utf-8')
    return data

def calculate_md5(data):
    data_bytes = _to_bytes(data)
    return hashlib.md5(data_bytes).hexdigest()

def calculate_sha1(data):
    data_bytes = _to_bytes(data)
    return hashlib.sha1(data_bytes).hexdigest()

def calculate_sha256(data):
    data_bytes = _to_bytes(data)
    return hashlib.sha256(data_bytes).hexdigest()

def calculate_sha512(data):
    data_bytes = _to_bytes(data)
    return hashlib.sha512(data_bytes).hexdigest()

def calculate_crc32(data):
    data_bytes = _to_bytes(data)
    return f"{zlib.crc32(data_bytes) & 0xFFFFFFFF:08x}"

def validate_checksum(data, expected_checksum, algorithm):
    calculated_checksum = None
    algorithm = algorithm.lower()

    if algorithm == 'md5':
        calculated_checksum = calculate_md5(data)
    elif algorithm == 'sha1':
        calculated_checksum = calculate_sha1(data)
    elif algorithm == 'sha256':
        calculated_checksum = calculate_sha256(data)
    elif algorithm == 'sha512':
        calculated_checksum = calculate_sha512(data)
    elif algorithm == 'crc32':
        calculated_checksum = calculate_crc32(data)
    
    return calculated_checksum == expected_checksum

# Additional implementation at 2025-08-04 08:15:07
import hashlib
import zlib
import os

class ChecksumValidator:
    """
    A class to calculate and validate checksums for files and strings
    using various algorithms like MD5, SHA (SHA1, SHA256, SHA512), and CRC32.
    """

    _CHUNK_SIZE = 65536  # 64KB

    def __init__(self):
        self._hashlib_algorithms = {
            "md5": hashlib.md5,
            "sha1": hashlib.sha1,
            "sha256": hashlib.sha256,
            "sha512": hashlib.sha512,
        }
        self._supported_algorithms = list(self._hashlib_algorithms.keys()) + ["crc32"]

    def _get_hashlib_hasher(self, algorithm: str):
        """Returns a hashlib hasher object or raises an error."""
        algo_lower = algorithm.lower()
        if algo_lower in self._hashlib_algorithms:
            return self._hashlib_algorithms[algo_lower]()
        else:
            raise ValueError(f"Unsupported hashlib algorithm: {algorithm}")

    def _calculate_hashlib_checksum_from_bytes(self, data_bytes: bytes, algorithm: str) -> str:
        """Calculates checksum for hashlib algorithms from a bytes object."""
        hasher = self._get_hashlib_hasher(algorithm)
        hasher.update(data_bytes)
        return hasher.hexdigest()

    def _calculate_crc32_checksum_from_bytes(self, data_bytes: bytes) -> str:
        """Calculates CRC32 checksum from a bytes object."""
        crc_val = zlib.crc32(data_bytes)
        return f"{crc_val & 0xFFFFFFFF:08x}" # Format as 8-digit unsigned hex

    def calculate_checksum_from_file(self, file_path: str, algorithm: str) -> str:
        """
        Calculates the checksum of a file.

        Args:
            file_path (str): The path to the file.
            algorithm (str): The hashing algorithm (e.g., 'md5', 'sha256', 'crc32').

        Returns:
            str: The hexadecimal checksum.

        Raises:
            FileNotFoundError: If the file does not exist.
            ValueError: If an unsupported algorithm is specified.
            IOError: If there's an error reading the file.
        """
        algo_lower = algorithm.lower()
        if algo_lower not in self._supported_algorithms:
            raise ValueError(f"Unsupported algorithm: {algorithm}. Supported algorithms are: {', '.join(self._supported_algorithms)}")

        if not os.path.exists(file_path):
            raise FileNotFoundError(f"File not found: {file_path}")

        try:
            with open(file_path, 'rb') as f:
                if algo_lower == "crc32":
                    crc_val = 0
                    while True:
                        chunk = f.read(self._CHUNK_SIZE)
                        if not chunk:
                            break
                        crc_val = zlib.crc32(chunk, crc_val)
                    return f"{crc_val & 0xFFFFFFFF:08x}"
                else:
                    hasher = self._get_hashlib_hasher(algo_lower)
                    while True:
                        chunk = f.read(self._CHUNK_SIZE)
                        if not chunk:
                            break
                        hasher.update(chunk)
                    return hasher.hexdigest()
        except Exception as e:
            raise IOError(f"Error reading file for checksum calculation: {e}")

    def calculate_checksum_from_string(self, data_string: str, algorithm: str) -> str:
        """
        Calculates the checksum of a string.

        Args:
            data_string (str): The input string.
            algorithm (str): The hashing algorithm (e.g., 'md5', 'sha256', 'crc32').

        Returns:
            str: The hexadecimal checksum.

        Raises:
            ValueError: If an unsupported algorithm is specified.
        """
        algo_lower = algorithm.lower()
        if algo_lower not in self._supported_algorithms:
            raise ValueError(f"Unsupported algorithm: {algorithm}. Supported algorithms are: {', '.join(self._supported_algorithms)}")

        data_bytes = data_string.encode('utf-8')

        if algo_lower == "crc32":
            return self._calculate_crc32_checksum_from_bytes(data_bytes)
        else:
            return self._calculate_hashlib_checksum_from_bytes(data_bytes, algo_lower)

    def validate_checksum(self, data_source: str, expected_checksum: str, algorithm: str, is_file: bool = True) -> bool:
        """
        Validates a given checksum against a calculated one.

        Args:
            data_source (str): The file path or string content to validate.
            expected_checksum (str): The checksum to validate against.
            algorithm (str): The hashing algorithm used (e.g., 'md5', 'sha256', 'crc32').
            is_file (bool): True if data_source is a file path, False if it's a string.

        Returns:
            bool: True if the calculated checksum matches the expected checksum, False otherwise.

        Raises:
            FileNotFoundError: If is_file is True and the file does not exist.
            ValueError: If an unsupported algorithm is specified.
            IOError: If there's an error reading the file (when is_file is True).
        """
        calculated_checksum = ""
        if is_file:
            calculated_checksum = self.calculate_checksum_from_file(data_source, algorithm)
        else:
            calculated_checksum = self.calculate_checksum_from_string(data_source, algorithm)

        return calculated_checksum.lower() == expected_checksum.lower()

# Additional implementation at 2025-08-04 08:16:16
import hashlib
import zlib
import os

class ChecksumValidator:
    def __init__(self):
        pass

    def calculate_checksum(self, filepath, algorithm, buffer_size=65536):
        if not os.path.exists(filepath):
            return None

        if algorithm in ['md5', 'sha1', 'sha256', 'sha512']:
            hasher = None
            if algorithm == 'md5':
                hasher = hashlib.md5()
            elif algorithm == 'sha1':
                hasher = hashlib.sha1()
            elif algorithm == 'sha256':
                hasher = hashlib.sha256()
            elif algorithm == 'sha512':
                hasher = hashlib.sha512()
            
            if hasher is None:
                raise ValueError(f"Unsupported hash algorithm: {algorithm}")

            try:
                with open(filepath, 'rb') as f:
                    while True:
                        chunk = f.read(buffer_size)
                        if not chunk:
                            break
                        hasher.update(chunk)
                return hasher.hexdigest()
            except IOError:
                return None
        elif algorithm == 'crc32':
            crc = 0
            try:
                with open(filepath, 'rb') as f:
                    while True:
                        chunk = f.read(buffer_size)
                        if not chunk:
                            break
                        crc = zlib.crc32(chunk, crc)
                return f"{crc & 0xFFFFFFFF:08x}"
            except IOError:
                return None
        else:
            raise ValueError(f"Unsupported checksum algorithm: {algorithm}")

    def validate_checksum(self, filepath, algorithm, expected_checksum, buffer_size=65536):
        calculated_checksum = self.calculate_checksum(filepath, algorithm, buffer_size)
        if calculated_checksum is None:
            return False
        return calculated_checksum.lower() == expected_checksum.lower()