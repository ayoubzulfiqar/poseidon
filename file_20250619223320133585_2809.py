from cryptography.fernet import Fernet
import hashlib

def generate_fernet_key():
    return Fernet.generate_key()

def encrypt_message(key, message):
    f = Fernet(key)
    encrypted_message = f.encrypt(message.encode())
    return encrypted_message

def decrypt_message(key, encrypted_message):
    f = Fernet(key)
    decrypted_message = f.decrypt(encrypted_message).decode()
    return decrypted_message

def hash_message(message, algorithm='sha256'):
    if algorithm == 'sha256':
        hasher = hashlib.sha256()
    elif algorithm == 'md5':
        hasher = hashlib.md5()
    elif algorithm == 'sha512':
        hasher = hashlib.sha512()
    else:
        raise ValueError("Unsupported hashing algorithm.")
    
    hasher.update(message.encode())
    return hasher.hexdigest()

def verify_hash(message, known_hash, algorithm='sha256'):
    calculated_hash = hash_message(message, algorithm)
    return calculated_hash == known_hash

if __name__ == "__main__":
    fernet_key = generate_fernet_key()
    print(f"Generated Fernet Key: {fernet_key.decode()}")

    original_message = "This is a secret message for encryption."
    print(f"Original Message: {original_message}")

    encrypted_data = encrypt_message(fernet_key, original_message)
    print(f"Encrypted Message: {encrypted_data.decode()}")

    decrypted_data = decrypt_message(fernet_key, encrypted_data)
    print(f"Decrypted Message: {decrypted_data}")

    print(f"Decryption successful: {original_message == decrypted_data}")
    print("-" * 40)

    message_to_hash = "Hello, world! This is a test for hashing."
    print(f"Message to Hash: {message_to_hash}")

    hashed_value = hash_message(message_to_hash, 'sha256')
    print(f"SHA256 Hash: {hashed_value}")

    is_verified = verify_hash(message_to_hash, hashed_value, 'sha256')
    print(f"Hash verification successful (same message): {is_verified}")

    different_message_2 = "Hello, world! This is a different message."
    is_verified_false = verify_hash(different_message_2, hashed_value, 'sha256')
    print(f"Hash verification successful (different message): {is_verified_false}")

    print("-" * 40)

    message_to_hash_md5 = "Another message for MD5."
    print(f"Message to Hash (MD5): {message_to_hash_md5}")

    hashed_value_md5 = hash_message(message_to_hash_md5, 'md5')
    print(f"MD5 Hash: {hashed_value_md5}")

    is_verified_md5 = verify_hash(message_to_hash_md5, hashed_value_md5, 'md5')
    print(f"MD5 Hash verification successful: {is_verified_md5}")
    print("-" * 40)

# Additional implementation at 2025-06-19 22:34:05
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes, hmac
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey, RSAPublicKey
import os
import base64

class CryptoToolkit:
    def __init__(self):
        pass

    def generate_fernet_key(self) -> bytes:
        """Generates a new Fernet key."""
        return Fernet.generate_key()

    def fernet_encrypt(self, data: bytes, key: bytes) -> bytes:
        """Encrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        return f.encrypt(data)

    def fernet_decrypt(self, encrypted_data: bytes, key: bytes) -> bytes:
        """Decrypts data using Fernet symmetric encryption."""
        f = Fernet(key)
        return f.decrypt(encrypted_data)

    def hash_data(self, data: bytes) -> bytes:
        """Hashes data using SHA256."""
        digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
        digest.update(data)
        return digest.finalize()

    def generate_rsa_key_pair(self, key_size: int = 2048) -> tuple[RSAPrivateKey, RSAPublicKey]:
        """Generates an RSA private and public key pair."""
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=key_size,
            backend=default_backend()
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def rsa_encrypt(self, data: bytes, public_key: RSAPublicKey) -> bytes:
        """Encrypts data using RSA public key with OAEP padding."""
        return public_key.encrypt(
            data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def rsa_decrypt(self, encrypted_data: bytes, private_key: RSAPrivateKey) -> bytes:
        """Decrypts data using RSA private key with OAEP padding."""
        return private_key.decrypt(
            encrypted_data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def rsa_sign(self, data: bytes, private_key: RSAPrivateKey) -> bytes:
        """Signs data using RSA private key with PSS padding."""
        return private_key.sign(
            data,
            padding.PSS(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )

    def rsa_verify(self, data: bytes, signature: bytes, public_key: RSAPublicKey) -> bool:
        """Verifies an RSA signature."""
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

    def derive_key_from_password(self, password: str, salt: bytes, iterations: int = 100000, key_length: int = 32) -> bytes:
        """Derives a cryptographic key from a password using PBKDF2HMAC."""
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=key_length,
            salt=salt,
            iterations=iterations,
            backend=default_backend()
        )
        return kdf.derive(password.encode('utf-8'))

    def generate_hmac(self, data: bytes, key: bytes) -> bytes:
        """Generates an HMAC for data using a shared secret key."""
        h = hmac.HMAC(key, hashes.SHA256(), backend=default_backend())
        h.update(data)
        return h.finalize()

    def verify_hmac(self, data: bytes, hmac_value: bytes, key: bytes) -> bool:
        """Verifies an HMAC against data and a shared secret key."""
        h = hmac.HMAC(key, hashes.SHA256(), backend=default_backend())
        h.update(data)
        try:
            h.verify(hmac_value)
            return True
        except Exception:
            return False

    def b64_encode(self, data: bytes) -> bytes:
        """Encodes bytes data to Base64 (URL-safe)."""
        return base64.urlsafe_b64encode(data)

    def b64_decode(self, encoded_data: bytes) -> bytes:
        """Decodes Base64 (URL-safe) encoded bytes data."""
        return base64.urlsafe_b64decode(encoded_data)

# Additional implementation at 2025-06-19 22:34:46
import os
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
import bcrypt

class CryptoToolkit:
    def __init__(self):
        pass

    def generate_fernet_key(self):
        return Fernet.generate_key()

    def encrypt_fernet(self, data: bytes, key: bytes) -> bytes:
        f = Fernet(key)
        return f.encrypt(data)

    def decrypt_fernet(self, encrypted_data: bytes, key: bytes) -> bytes:
        f = Fernet(key)
        return f.decrypt(encrypted_data)

    def generate_rsa_keys(self, public_exponent=65537, key_size=2048):
        private_key = rsa.generate_private_key(
            public_exponent=public_exponent,
            key_size=key_size,
            backend=default_backend()
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def encrypt_rsa(self, data: bytes, public_key) -> bytes:
        return public_key.encrypt(
            data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def decrypt_rsa(self, encrypted_data: bytes, private_key) -> bytes:
        return private_key.decrypt(
            encrypted_data,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def sign_rsa(self, data: bytes, private_key) -> bytes:
        signer = private_key.signer(
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        signer.update(data)
        return signer.finalize()

    def verify_rsa_signature(self, data: bytes, signature: bytes, public_key) -> bool:
        verifier = public_key.verifier(
            signature,
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        verifier.update(data)
        try:
            verifier.verify()
            return True
        except Exception:
            return False

    def hash_sha256(self, data: bytes) -> bytes:
        digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
        digest.update(data)
        return digest.finalize()

    def hash_password_bcrypt(self, password: str) -> bytes:
        hashed = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
        return hashed

    def verify_password_bcrypt(self, password: str, hashed_password: bytes) -> bool:
        try:
            return bcrypt.checkpw(password.encode('utf-8'), hashed_password)
        except ValueError:
            return False

    def encrypt_file_fernet(self, input_filepath: str, output_filepath: str, key: bytes):
        f = Fernet(key)
        with open(input_filepath, 'rb') as f_in:
            original_data = f_in.read()
        encrypted_data = f.encrypt(original_data)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(encrypted_data)

    def decrypt_file_fernet(self, input_filepath: str, output_filepath: str, key: bytes):
        f = Fernet(key)
        with open(input_filepath, 'rb') as f_in:
            encrypted_data = f_in.read()
        decrypted_data = f.decrypt(encrypted_data)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(decrypted_data)

    def save_private_key(self, private_key, filepath: str, password: str = None):
        if password:
            encryption_algorithm = serialization.BestAvailableEncryption(password.encode('utf-8'))
        else:
            encryption_algorithm = serialization.NoEncryption()

        pem = private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=encryption_algorithm
        )
        with open(filepath, 'wb') as f:
            f.write(pem)

    def load_private_key(self, filepath: str, password: str = None):
        with open(filepath, 'rb') as f:
            pem = f.read()
        return serialization.load_pem_private_key(
            pem,
            password=password.encode('utf-8') if password else None,
            backend=default_backend()
        )

    def save_public_key(self, public_key, filepath: str):
        pem = public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
        with open(filepath, 'wb') as f:
            f.write(pem)

    def load_public_key(self, filepath: str):
        with open(filepath, 'rb') as f:
            pem = f.read()
        return serialization.load_pem_public_key(
            pem,
            backend=default_backend()
        )

# Additional implementation at 2025-06-19 22:35:23
import os
import hashlib
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import padding as sym_padding
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding as asym_padding
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend

class CryptoToolkit:
    def __init__(self):
        self.aes_key = None
        self.rsa_private_key = None
        self.rsa_public_key = None

    def generate_aes_key(self, key_size_bits=256):
        """Generates a random AES key."""
        if key_size_bits not in [128, 192, 256]:
            raise ValueError("AES key size must be 128, 192, or 256 bits.")
        self.aes_key = os.urandom(key_size_bits // 8)
        return self.aes_key

    def set_aes_key(self, key: bytes):
        """Sets the AES key."""
        if len(key) * 8 not in [128, 192, 256]:
            raise ValueError("AES key must be 16, 24, or 32 bytes (128, 192, or 256 bits).")
        self.aes_key = key

    def encrypt_aes(self, plaintext: bytes) -> bytes:
        """Encrypts data using AES in CBC mode with PKCS7 padding."""
        if not self.aes_key:
            raise ValueError("AES key not set. Generate or set a key first.")

        iv = os.urandom(16)
        cipher = Cipher(algorithms.AES(self.aes_key), modes.CBC(iv), backend=default_backend())
        encryptor = cipher.encryptor()

        padder = sym_padding.PKCS7(algorithms.AES.block_size).padder()
        padded_data = padder.update(plaintext) + padder.finalize()

        ciphertext = encryptor.update(padded_data) + encryptor.finalize()
        return iv + ciphertext

    def decrypt_aes(self, ciphertext_with_iv: bytes) -> bytes:
        """Decrypts data using AES in CBC mode with PKCS7 padding."""
        if not self.aes_key:
            raise ValueError("AES key not set. Generate or set a key first.")
        if len(ciphertext_with_iv) < 16:
            raise ValueError("Ciphertext too short to contain IV.")

        iv = ciphertext_with_iv[:16]
        ciphertext = ciphertext_with_iv[16:]

        cipher = Cipher(algorithms.AES(self.aes_key), modes.CBC(iv), backend=default_backend())
        decryptor = cipher.decryptor()

        padded_plaintext = decryptor.update(ciphertext) + decryptor.finalize()

        unpadder = sym_padding.PKCS7(algorithms.AES.block_size).unpadder()
        plaintext = unpadder.update(padded_plaintext) + unpadder.finalize()
        return plaintext

    def hash_sha256(self, data: bytes) -> bytes:
        """Computes the SHA256 hash of the given data."""
        return hashlib.sha256(data).digest()

    def hash_sha512(self, data: bytes) -> bytes:
        """Computes the SHA512 hash of the given data."""
        return hashlib.sha512(data).digest()

    def generate_rsa_key_pair(self, public_exponent=65537, key_size=2048):
        """Generates a new RSA private and public key pair."""
        self.rsa_private_key = rsa.generate_private_key(
            public_exponent=public_exponent,
            key_size=key_size,
            backend=default_backend()
        )
        self.rsa_public_key = self.rsa_private_key.public_key()
        return self.rsa_private_key, self.rsa_public_key

    def get_rsa_public_key(self):
        """Returns the current RSA public key."""
        if not self.rsa_public_key:
            raise ValueError("RSA public key not generated. Generate a key pair first.")
        return self.rsa_public_key

    def get_rsa_private_key(self):
        """Returns the current RSA private key."""
        if not self.rsa_private_key:
            raise ValueError("RSA private key not generated. Generate a key pair first.")
        return self.rsa_private_key

    def encrypt_rsa(self, plaintext: bytes, public_key) -> bytes:
        """Encrypts data using RSA with OAEP padding."""
        if len(plaintext) > (public_key.key_size // 8) - 42:
            raise ValueError(f"Plaintext too long for RSA key size {public_key.key_size}. Max length for OAEP is approx {public_key.key_size // 8 - 42} bytes.")
        return public_key.encrypt(
            plaintext,
            asym_padding.OAEP(
                mgf=asym_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def decrypt_rsa(self, ciphertext: bytes) -> bytes:
        """Decrypts data using RSA with OAEP padding."""
        if not self.rsa_private_key:
            raise ValueError("RSA private key not set. Generate or load a key first.")
        return self.rsa_private_key.decrypt(
            ciphertext,
            asym_padding.OAEP(
                mgf=asym_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def sign_data(self, data: bytes) -> bytes:
        """Signs data using the RSA private key with PSS padding."""
        if not self.rsa_private_key:
            raise ValueError("RSA private key not set. Generate or load a key first.")
        return self.rsa_private_key.sign(
            data,
            asym_padding.PSS(
                mgf=asym_padding.MGF1(hashes.SHA256()),
                salt_length=asym_padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )

    def verify_signature(self, data: bytes, signature: bytes, public_key) -> bool:
        """Verifies a signature using the RSA public key with PSS padding."""
        try:
            public_key.verify(
                signature,
                data,
                asym_padding.PSS(
                    mgf=asym_padding.MGF1(hashes.SHA256()),
                    salt_length=asym_padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            return True
        except Exception:
            return False

    def save_rsa_private_key(self, filename: str, password: bytes = None):
        """Saves the RSA private key to a file, optionally encrypted."""
        if not self.rsa_private_key:
            raise ValueError("RSA private key not generated. Generate a key pair first.")

        if password:
            encryption_algorithm = serialization.BestAvailableEncryption(password)
        else:
            encryption_algorithm = serialization.NoEncryption()

        pem = self.rsa_private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=encryption_algorithm
        )
        with open(filename, 'wb') as f:
            f.write(pem)

    def load_rsa_private_key(self, filename: str, password: bytes = None):
        """Loads an RSA private key from a file."""
        with open(filename, 'rb') as f:
            pem = f.read()
        self.rsa_private_key = serialization.load_pem_private_key(
            pem,
            password=password,
            backend=default_backend()
        )
        self.rsa_public_key = self.rsa_private_key.public_key()
        return self.rsa_private_key

    def save_rsa_public_key(self, filename: str):
        """Saves the RSA public key to a file."""
        if not self.rsa_public_key:
            raise ValueError("RSA public key not generated. Generate a key pair first.")
        pem = self.rsa_public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
        with open(filename, 'wb') as f:
            f.write(pem)

    def load_rsa_public_key(self, filename: str):
        """Loads an RSA public key from a file."""
        with open(filename, 'rb') as f:
            pem = f.read()
        public_key = serialization.load_pem_public_key(
            pem,
            backend=default_backend()
        )
        return public_key