import sys

def _caesar_transform(text, shift, encrypt=True):
    result = ""
    for char in text:
        if 'a' <= char <= 'z':
            base = ord('a')
            if encrypt:
                transformed_char = chr(((ord(char) - base + shift) % 26) + base)
            else:
                transformed_char = chr(((ord(char) - base - shift + 26) % 26) + base)
            result += transformed_char
        elif 'A' <= char <= 'Z':
            base = ord('A')
            if encrypt:
                transformed_char = chr(((ord(char) - base + shift) % 26) + base)
            else:
                transformed_char = chr(((ord(char) - base - shift + 26) % 26) + base)
            result += transformed_char
        else:
            result += char
    return result

def caesar_encrypt(text, shift):
    return _caesar_transform(text, shift, encrypt=True)

def caesar_decrypt(text, shift):
    return _caesar_transform(text, shift, encrypt=False)

def run_cryptography_toolkit():
    while True:
        print("\n--- Simple Cryptography Toolkit ---")
        print("1. Encrypt Message (Caesar Cipher)")
        print("2. Decrypt Message (Caesar Cipher)")
        print("3. Exit")
        
        choice = input("Enter your choice: ")
        
        if choice == '1':
            message = input("Enter the message to encrypt: ")
            try:
                shift = int(input("Enter the shift value (integer): "))
                encrypted_message = caesar_encrypt(message, shift)
                print(f"Encrypted message: {encrypted_message}")
            except ValueError:
                print("Invalid shift value. Please enter an integer.")
        elif choice == '2':
            message = input("Enter the message to decrypt: ")
            try:
                shift = int(input("Enter the shift value (integer): "))
                decrypted_message = caesar_decrypt(message, shift)
                print(f"Decrypted message: {decrypted_message}")
            except ValueError:
                print("Invalid shift value. Please enter an integer.")
        elif choice == '3':
            print("Exiting toolkit. Goodbye!")
            sys.exit()
        else:
            print("Invalid choice. Please select a valid option (1, 2, or 3).")

if __name__ == "__main__":
    run_cryptography_toolkit()

# Additional implementation at 2025-06-23 02:20:16
import os
import base64
import hashlib
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend

class CryptoToolkit:
    def __init__(self):
        pass

    def generate_fernet_key(self) -> bytes:
        """Generates a new Fernet key for symmetric encryption."""
        return Fernet.generate_key()

    def encrypt_fernet(self, key: bytes, data: bytes) -> bytes:
        """Encrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        return f.encrypt(data)

    def decrypt_fernet(self, key: bytes, encrypted_data: bytes) -> bytes:
        """Decrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        try:
            return f.decrypt(encrypted_data)
        except Exception as e:
            raise ValueError(f"Fernet decryption failed: {e}")

    def hash_data(self, data: bytes, algorithm: str = 'sha256') -> str:
        """Hashes data using the specified algorithm (SHA256, SHA512, MD5)."""
        if algorithm == 'sha256':
            return hashlib.sha256(data).hexdigest()
        elif algorithm == 'sha512':
            return hashlib.sha512(data).hexdigest()
        elif algorithm == 'md5':
            return hashlib.md5(data).hexdigest()
        else:
            raise ValueError("Unsupported hash algorithm. Choose 'sha256', 'sha512', or 'md5'.")

    def generate_rsa_keys(self, key_size: int = 2048):
        """Generates a new RSA private and public key pair."""
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=key_size,
            backend=default_backend()
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def encrypt_rsa(self, public_key, data: bytes) -> bytes:
        """Encrypts data using RSA public key (OAEP padding).
        Note: RSA encryption is limited to small amounts of data (e.g., symmetric keys).
        """
        return public_key.encrypt(
            data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def decrypt_rsa(self, private_key, encrypted_data: bytes) -> bytes:
        """Decrypts data using RSA private key (OAEP padding)."""
        try:
            return private_key.decrypt(
                encrypted_data,
                padding.OAEP(
                    mgf=padding.MGF1(algorithm=hashes.SHA256()),
                    algorithm=hashes.SHA256(),
                    label=None
                )
            )
        except Exception as e:
            raise ValueError(f"RSA decryption failed: {e}")

    def sign_data(self, private_key, data: bytes) -> bytes:
        """Signs data using RSA private key (PSS padding)."""
        return private_key.sign(
            data,
            padding.PSS(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )

    def verify_signature(self, public_key, data: bytes, signature: bytes) -> bool:
        """Verifies a signature using RSA public key (PSS padding)."""
        try:
            public_key.verify(
                signature,
                data,
                padding.PSS(
                    mgf=padding.MGF1(algorithm=hashes.SHA256()),
                    salt_length=padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            return True
        except Exception:
            return False

    def derive_key_from_password(self, password: str, salt: bytes = None, iterations: int = 480000, length: int = 32) -> tuple[bytes, bytes]:
        """Derives a cryptographic key from a password using PBKDF2HMAC.
        Returns the derived key and the salt used.
        """
        if salt is None:
            salt = os.urandom(16) # Generate a new salt if not provided

        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=length,
            salt=salt,
            iterations=iterations,
            backend=default_backend()
        )
        key = kdf.derive(password.encode('utf-8'))
        return key, salt

    def encrypt_file_fernet(self, key: bytes, input_filepath: str, output_filepath: str):
        """Encrypts a file using Fernet symmetric encryption."""
        f = Fernet(key)
        with open(input_filepath, 'rb') as file:
            file_data = file.read()
        encrypted_data = f.encrypt(file_data)
        with open(output_filepath, 'wb') as file:
            file.write(encrypted_data)

    def decrypt_file_fernet(self, key: bytes, input_filepath: str, output_filepath: str):
        """Decrypts a file using Fernet symmetric encryption."""
        f = Fernet(key)
        with open(input_filepath, 'rb') as file:
            encrypted_data = file.read()
        try:
            decrypted_data = f.decrypt(encrypted_data)
            with open(output_filepath, 'wb') as file:
                file.write(decrypted_data)
        except Exception as e:
            raise ValueError(f"File decryption failed: {e}")

if __name__ == "__main__":
    toolkit = CryptoToolkit()

    # --- Fernet Symmetric Encryption/Decryption ---
    fernet_key = toolkit.generate_fernet_key()
    original_message_fernet = b"This is a secret message using Fernet!"
    print(f"Original Fernet message: {original_message_fernet}")

    encrypted_fernet_message = toolkit.encrypt_fernet(fernet_key, original_message_fernet)
    print(f"Encrypted Fernet message: {encrypted_fernet_message}")

    decrypted_fernet_message = toolkit.decrypt_fernet(fernet_key, encrypted_fernet_message)
    print(f"Decrypted Fernet message: {decrypted_fernet_message}")
    print(f"Fernet decryption successful: {original_message_fernet == decrypted_fernet_message}\n")

    # --- Hashing ---
    data_to_hash = b"Hello, cryptography!"
    print(f"Data to hash: {data_to_hash}")

    sha256_hash = toolkit.hash_data(data_to_hash, 'sha256')
    print(f"SHA256 Hash: {sha256_hash}")

    sha512_hash = toolkit.hash_data(data_to_hash, 'sha512')
    print(f"SHA512 Hash: {sha512_hash}")

    md5_hash = toolkit.hash_data(data_to_hash, 'md5')
    print(f"MD5 Hash: {md5_hash}\n")

    # --- RSA Asymmetric Encryption/Decryption and Signing ---
    private_key_rsa, public_key_rsa = toolkit.generate_rsa_keys()
    original_rsa_message = b"This is a small secret for RSA encryption."
    print(f"Original RSA message: {original_rsa_message}")

    try:
        encrypted_rsa_message = toolkit.encrypt_rsa(public_key_rsa, original_rsa_message)
        print(f"Encrypted RSA message: {encrypted_rsa_message}")

        decrypted_rsa_message = toolkit.decrypt_rsa(private_key_rsa, encrypted_rsa_message)
        print(f"Decrypted RSA message: {decrypted_rsa_message}")
        print(f"RSA decryption successful: {original_rsa_message == decrypted_rsa_message}\n")
    except ValueError as e:
        print(f"RSA encryption/decryption failed (message too long?): {e}\n")

    data_to_sign = b"This data needs to be signed and verified."
    print(f"Data to sign: {data_to_sign}")

    signature = toolkit.sign_data(private_key_rsa, data_to_sign)
    print(f"Signature: {signature}")

    is_valid = toolkit.verify_signature(public_key_rsa, data_to_sign, signature)
    print(f"Signature valid: {is_valid}")

    tampered_data = b"This data needs to be signed and verified. TAMPERED!"
    is_valid_tampered = toolkit.verify_signature(public_key_rsa, tampered_data, signature)
    print(f"Signature valid with tampered data: {is_valid_tampered}\n")

    # --- PBKDF2 Key Derivation ---
    password = "mysecretpassword123"
    print(f"Password: {password}")

    derived_key_1, salt_1 = toolkit.derive_key_from_password(password)
    print(f"Derived Key 1 (base64): {base64.urlsafe_b64encode(derived_key_1)}")
    print(f"Salt 1 (base64): {base64.urlsafe_b64encode(salt_1)}")

    derived_key_2, _ = toolkit.derive_key_from_password(password, salt=salt_1)
    print(f"Derived Key 2 (base64, same salt): {base64.urlsafe_b64encode(derived_key_2)}")
    print(f"Keys match: {derived_key_1 == derived_key_2}\n")

    # --- File Encryption/Decryption (Fernet) ---
    dummy_file_name = "test_file.txt"
    encrypted_file_name = "test_file.enc"
    decrypted_file_name = "test_file.dec"

    with open(dummy_file_name, "w") as f:
        f.write("This is some test content for file encryption.\n")
        f.write("It has multiple lines to ensure full file handling.\n")
        f.write("And a final line.\n")
    print(f"Created dummy file: {dummy_file_name}")

    file_fernet_key = toolkit.generate_fernet_key()

    toolkit.encrypt_file_fernet(file_fernet_key, dummy_file_name, encrypted_file_name)
    print(f"Encrypted '{dummy_file_name}' to '{encrypted_file_name}'")

    toolkit.decrypt_file_fernet(file_fernet_key, encrypted_file_name, decrypted_file_name)
    print(f"Decrypted '{encrypted_file_name}' to '{decrypted_file_name}'")

    with open(dummy_file_name, 'rb') as f_orig, open(decrypted_file_name, 'rb') as f_dec:
        original_content = f_orig.read()
        decrypted_content = f_dec.read()
        print(f"File content matches: {original_content == decrypted_content}")

    os.remove(dummy_file_name)
    os.remove(encrypted_file_name)
    os.remove(decrypted_file_name)
    print("Cleaned up dummy files.")

# Additional implementation at 2025-06-23 02:21:17
import os
import base64
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import utils as asymmetric_utils

class CryptoToolkit:
    def __init__(self):
        pass

    # --- Symmetric Encryption (Fernet) ---
    def generate_fernet_key(self) -> bytes:
        """Generates a Fernet key."""
        return Fernet.generate_key()

    def encrypt_fernet(self, key: bytes, data: bytes) -> bytes:
        """Encrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        return f.encrypt(data)

    def decrypt_fernet(self, key: bytes, encrypted_data: bytes) -> bytes:
        """Decrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        return f.decrypt(encrypted_data)

    # --- Hashing ---
    def hash_data(self, data: bytes, algorithm: str = 'sha256') -> bytes:
        """
        Hashes data using specified algorithm (SHA256 or SHA512).
        Returns the hexadecimal digest.
        """
        digest = hashes.Hash(self._get_hash_algorithm(algorithm), backend=default_backend())
        digest.update(data)
        return digest.finalize()

    def _get_hash_algorithm(self, algorithm_name: str):
        """Helper to get hash algorithm object."""
        if algorithm_name.lower() == 'sha256':
            return hashes.SHA256()
        elif algorithm_name.lower() == 'sha512':
            return hashes.SHA512()
        else:
            raise ValueError("Unsupported hash algorithm. Choose 'sha256' or 'sha512'.")

    # --- Password Hashing (KDF - PBKDF2) ---
    def derive_key_pbkdf2(self, password: bytes, salt: bytes = None, iterations: int = 480000, key_length: int = 32) -> tuple[bytes, bytes]:
        """
        Derives a key from a password using PBKDF2HMAC.
        If no salt is provided, a new one is generated.
        Returns (derived_key, salt).
        """
        if salt is None:
            salt = os.urandom(16) # 16 bytes is standard for PBKDF2 salt

        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=key_length,
            salt=salt,
            iterations=iterations,
            backend=default_backend()
        )
        derived_key = kdf.derive(password)
        return derived_key, salt

    # --- Asymmetric Encryption (RSA) ---
    def generate_rsa_keys(self, key_size: int = 2048) -> tuple[rsa.RSAPrivateKey, rsa.RSAPublicKey]:
        """Generates RSA private and public keys."""
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=key_size,
            backend=default_backend()
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def encrypt_rsa(self, public_key: rsa.RSAPublicKey, data: bytes) -> bytes:
        """Encrypts data using RSA public key."""
        return public_key.encrypt(
            data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def decrypt_rsa(self, private_key: rsa.RSAPrivateKey, encrypted_data: bytes) -> bytes:
        """Decrypts data using RSA private key."""
        return private_key.decrypt(
            encrypted_data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def sign_data(self, private_key: rsa.RSAPrivateKey, data: bytes) -> bytes:
        """Signs data using RSA private key."""
        signer = private_key.signer(
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        signer.update(data)
        return signer.finalize()

    def verify_signature(self, public_key: rsa.RSAPublicKey, data: bytes, signature: bytes) -> bool:
        """Verifies a signature using RSA public key."""
        try:
            verifier = public_key.verifier(
                signature,
                padding.PSS(
                    mgf=padding.MGF1(hashes.SHA256()),
                    salt_length=padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            verifier.update(data)
            verifier.verify()
            return True
        except Exception: # InvalidSignature or other exceptions
            return False

    # --- Key Serialization (for RSA keys) ---
    def serialize_private_key(self, private_key: rsa.RSAPrivateKey, password: bytes = None) -> bytes:
        """Serializes a private key to PEM format."""
        encryption_algorithm = serialization.NoEncryption()
        if password:
            encryption_algorithm = serialization.BestAvailableEncryption(password)

        return private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=encryption_algorithm
        )

    def deserialize_private_key(self, pem_data: bytes, password: bytes = None) -> rsa.RSAPrivateKey:
        """Deserializes a private key from PEM format."""
        return serialization.load_pem_private_key(
            pem_data,
            password=password,
            backend=default_backend()
        )

    def serialize_public_key(self, public_key: rsa.RSAPublicKey) -> bytes:
        """Serializes a public key to PEM format."""
        return public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )

    def deserialize_public_key(self, pem_data: bytes) -> rsa.RSAPublicKey:
        """Deserializes a public key from PEM format."""
        return serialization.load_pem_public_key(
            pem_data,
            backend=default_backend()
        )

    # --- Utility Functions ---
    def base64_encode(self, data: bytes) -> bytes:
        """Encodes bytes data to Base64 (URL-safe)."""
        return base64.urlsafe_b64encode(data)

    def base64_decode(self, encoded_data: bytes) -> bytes:
        """Decodes Base64 (URL-safe) encoded bytes data."""
        return base64.urlsafe_b64decode(encoded_data)

    # --- File Encryption/Decryption (Fernet) ---
    def encrypt_file_fernet(self, key: bytes, input_filepath: str, output_filepath: str):
        """Encrypts a file using Fernet."""
        f = Fernet(key)
        with open(input_filepath, 'rb') as f_in:
            file_data = f_in.read()
        encrypted_data = f.encrypt(file_data)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(encrypted_data)

    def decrypt_file_fernet(self, key: bytes, input_filepath: str, output_filepath: str):
        """Decrypts a file using Fernet."""
        f = Fernet(key)
        with open(input_filepath, 'rb') as f_in:
            encrypted_data = f_in.read()
        decrypted_data = f.decrypt(encrypted_data)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(decrypted_data)