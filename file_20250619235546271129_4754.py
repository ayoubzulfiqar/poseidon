from cryptography.fernet import Fernet
import hashlib

class SimpleCryptoToolkit:
    def generate_fernet_key(self) -> bytes:
        """Generates a new Fernet encryption key."""
        return Fernet.generate_key()

    def encrypt_message(self, key: bytes, message: str) -> bytes:
        """Encrypts a string message using a Fernet key."""
        f = Fernet(key)
        encrypted_data = f.encrypt(message.encode('utf-8'))
        return encrypted_data

    def decrypt_message(self, key: bytes, encrypted_message: bytes) -> str:
        """Decrypts a Fernet-encrypted message using the corresponding key."""
        f = Fernet(key)
        decrypted_data = f.decrypt(encrypted_message)
        return decrypted_data.decode('utf-8')

    def hash_message(self, message: str) -> str:
        """Hashes a string message using SHA256."""
        sha256_hash = hashlib.sha256()
        sha256_hash.update(message.encode('utf-8'))
        return sha256_hash.hexdigest()

if __name__ == "__main__":
    toolkit = SimpleCryptoToolkit()

    # --- Demonstrate Fernet Symmetric Encryption/Decryption ---
    print("--- Fernet Symmetric Encryption/Decryption ---")
    fernet_key = toolkit.generate_fernet_key()
    print(f"Generated Fernet Key: {fernet_key.decode('utf-8')}")

    original_message = "This is a highly confidential message that needs to be encrypted."
    print(f"Original Message: {original_message}")

    encrypted_msg = toolkit.encrypt_message(fernet_key, original_message)
    print(f"Encrypted Message (bytes): {encrypted_msg}")

    decrypted_msg = toolkit.decrypt_message(fernet_key, encrypted_msg)
    print(f"Decrypted Message: {decrypted_msg}")

    if original_message == decrypted_msg:
        print("Fernet Encryption/Decryption successful!")
    else:
        print("Fernet Encryption/Decryption FAILED!")

    # --- Demonstrate SHA256 Hashing ---
    print("\n--- SHA256 Hashing ---")
    message_to_hash_1 = "Cryptography is fascinating!"
    print(f"Message to Hash 1: {message_to_hash_1}")
    hashed_output_1 = toolkit.hash_message(message_to_hash_1)
    print(f"SHA256 Hash 1: {hashed_output_1}")

    message_to_hash_2 = "Cryptography is fascinating!" # Same message
    print(f"Message to Hash 2: {message_to_hash_2}")
    hashed_output_2 = toolkit.hash_message(message_to_hash_2)
    print(f"SHA256 Hash 2: {hashed_output_2}")

    if hashed_output_1 == hashed_output_2:
        print("Hashing consistency confirmed for identical messages.")
    else:
        print("Hashing consistency FAILED for identical messages.")

    message_to_hash_3 = "Cryptography is fascinating!" # Different message (extra space)
    print(f"Message to Hash 3: {message_to_hash_3} ")
    hashed_output_3 = toolkit.hash_message(message_to_hash_3 + " ")
    print(f"SHA256 Hash 3: {hashed_output_3}")

    if hashed_output_1 != hashed_output_3:
        print("Hashing produces different outputs for different messages (even slight changes).")
    else:
        print("Hashing collision detected (should not happen for different messages).")

# Additional implementation at 2025-06-19 23:56:45
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.asymmetric import rsa, padding as asym_padding
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.serialization import load_pem_public_key, load_pem_private_key
import os
import base64

class CryptoToolkit:
    def __init__(self):
        self.backend = default_backend()

    def generate_aes_key(self, key_size_bits=256):
        """Generates a random AES key."""
        if key_size_bits not in [128, 192, 256]:
            raise ValueError("AES key size must be 128, 192, or 256 bits.")
        return os.urandom(key_size_bits // 8)

    def aes_gcm_encrypt(self, key, plaintext):
        """Encrypts data using AES-256 GCM mode."""
        iv = os.urandom(12)  # 96-bit IV for GCM
        cipher = Cipher(algorithms.AES(key), modes.GCM(iv), backend=self.backend)
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(plaintext) + encryptor.finalize()
        tag = encryptor.tag
        return iv + ciphertext + tag  # Concatenate IV, ciphertext, and tag

    def aes_gcm_decrypt(self, key, encrypted_data):
        """Decrypts data using AES-256 GCM mode."""
        if len(encrypted_data) < 28: # 12 (IV) + 16 (tag)
            raise ValueError("Encrypted data is too short for AES GCM decryption.")
        iv = encrypted_data[:12]
        ciphertext = encrypted_data[12:-16]  # GCM tag is 16 bytes
        tag = encrypted_data[-16:]

        cipher = Cipher(algorithms.AES(key), modes.GCM(iv, tag), backend=self.backend)
        decryptor = cipher.decryptor()
        return decryptor.update(ciphertext) + decryptor.finalize()

    def hash_sha256(self, data):
        """Computes the SHA256 hash of the given data."""
        digest = hashes.Hash(hashes.SHA256(), backend=self.backend)
        digest.update(data)
        return digest.finalize()

    def hash_sha512(self, data):
        """Computes the SHA512 hash of the given data."""
        digest = hashes.Hash(hashes.SHA512(), backend=self.backend)
        digest.update(data)
        return digest.finalize()

    def generate_salt(self, length=16):
        """Generates a random salt for key derivation functions."""
        return os.urandom(length)

    def derive_key_pbkdf2(self, password, salt, iterations=310000, key_length=32, hash_algorithm=hashes.SHA256()):
        """Derives a key from a password using PBKDF2HMAC."""
        kdf = PBKDF2HMAC(
            algorithm=hash_algorithm,
            length=key_length,
            salt=salt,
            iterations=iterations,
            backend=self.backend
        )
        return kdf.derive(password)

    def generate_rsa_key_pair(self, public_exponent=65537, key_size=2048):
        """Generates an RSA private and public key pair."""
        private_key = rsa.generate_private_key(
            public_exponent=public_exponent,
            key_size=key_size,
            backend=self.backend
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def serialize_private_key(self, private_key, password=None):
        """Serializes an RSA private key to PEM format."""
        if password:
            encryption_algorithm = serialization.BestAvailableEncryption(password)
        else:
            encryption_algorithm = serialization.NoEncryption()
        return private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=encryption_algorithm
        )

    def serialize_public_key(self, public_key):
        """Serializes an RSA public key to PEM format."""
        return public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )

    def load_private_key_from_pem(self, pem_data, password=None):
        """Loads an RSA private key from PEM formatted bytes."""
        return load_pem_private_key(pem_data, password=password, backend=self.backend)

    def load_public_key_from_pem(self, pem_data):
        """Loads an RSA public key from PEM formatted bytes."""
        return load_pem_public_key(pem_data, backend=self.backend)

    def rsa_encrypt(self, public_key, plaintext):
        """Encrypts data using RSA with OAEP padding. Suitable for small data like symmetric keys."""
        return public_key.encrypt(
            plaintext,
            asym_padding.OAEP(
                mgf=asym_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def rsa_decrypt(self, private_key, ciphertext):
        """Decrypts data using RSA with OAEP padding."""
        return private_key.decrypt(
            ciphertext,
            asym_padding.OAEP(
                mgf=asym_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def rsa_sign(self, private_key, data_to_sign):
        """Signs data using RSA with PSS padding."""
        signer = private_key.signer(
            asym_padding.PSS(
                mgf=asym_padding.MGF1(hashes.SHA256()),
                salt_length=asym_padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        signer.update(data_to_sign)
        return signer.finalize()

    def rsa_verify_signature(self, public_key, data_signed, signature):
        """Verifies an RSA signature with PSS padding."""
        verifier = public_key.verifier(
            signature,
            asym_padding.PSS(
                mgf=asym_padding.MGF1(hashes.SHA256()),
                salt_length=asym_padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        verifier.update(data_signed)
        try:
            verifier.verify()
            return True
        except Exception:
            return False

    def generate_secure_random_bytes(self, length):
        """Generates cryptographically secure random bytes."""
        return os.urandom(length)

    def base64_encode(self, data_bytes):
        """Encodes bytes to Base64."""
        return base64.b64encode(data_bytes)

    def base64_decode(self, base64_bytes):
        """Decodes Base64 bytes."""
        return base64.b64decode(base64_bytes)

    def encrypt_file(self, input_filepath, output_filepath, key):
        """Encrypts a file using AES-256 GCM."""
        with open(input_filepath, 'rb') as f_in:
            plaintext = f_in.read()
        encrypted_data = self.aes_gcm_encrypt(key, plaintext)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(encrypted_data)

    def decrypt_file(self, input_filepath, output_filepath, key):
        """Decrypts a file using AES-256 GCM."""
        with open(input_filepath, 'rb') as f_in:
            encrypted_data = f_in.read()
        decrypted_data = self.aes_gcm_decrypt(key, encrypted_data)
        with open(output_filepath, 'wb') as f_out:
            f_out.write(decrypted_data)

# Additional implementation at 2025-06-19 23:57:27
import os
import base64
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import rsa, padding as asymmetric_padding
from cryptography.hazmat.primitives.kdf.scrypt import Scrypt

class CryptoToolkit:
    def __init__(self):
        self.backend = default_backend()

    def generate_aes_key(self, key_size_bits=256):
        """Generates a random AES key."""
        if key_size_bits not in [128, 192, 256]:
            raise ValueError("AES key size must be 128, 192, or 256 bits.")
        return os.urandom(key_size_bits // 8)

    def aes_encrypt(self, key, plaintext):
        """Encrypts data using AES in GCM mode."""
        if not isinstance(key, bytes) or len(key) * 8 not in [128, 192, 256]:
            raise ValueError("Key must be bytes of 16, 24, or 32 length for AES-128, 192, 256 respectively.")
        if not isinstance(plaintext, bytes):
            plaintext = plaintext.encode('utf-8')

        iv = os.urandom(12) # GCM recommended IV size is 12 bytes
        cipher = Cipher(algorithms.AES(key), modes.GCM(iv), backend=self.backend)
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(plaintext) + encryptor.finalize()
        tag = encryptor.tag
        return iv + ciphertext + tag

    def aes_decrypt(self, key, encrypted_data):
        """Decrypts data using AES in GCM mode."""
        if not isinstance(key, bytes) or len(key) * 8 not in [128, 192, 256]:
            raise ValueError("Key must be bytes of 16, 24, or 32 length for AES-128, 192, 256 respectively.")
        if not isinstance(encrypted_data, bytes):
            raise ValueError("Encrypted data must be bytes.")

        iv_len = 12
        tag_len = 16
        if len(encrypted_data) < iv_len + tag_len:
            raise ValueError("Encrypted data too short to contain IV and tag.")

        iv = encrypted_data[:iv_len]
        ciphertext = encrypted_data[iv_len:-tag_len]
        tag = encrypted_data[-tag_len:]

        cipher = Cipher(algorithms.AES(key), modes.GCM(iv, tag), backend=self.backend)
        decryptor = cipher.decryptor()
        try:
            plaintext = decryptor.update(ciphertext) + decryptor.finalize()
            return plaintext
        except Exception as e:
            raise ValueError("Decryption failed, possibly due to invalid key or corrupted data/tag.") from e

    def hash_data(self, data, algorithm="SHA256"):
        """Hashes data using the specified algorithm (e.g., SHA256, SHA512)."""
        if not isinstance(data, bytes):
            data = data.encode('utf-8')

        if algorithm == "SHA256":
            digest = hashes.Hash(hashes.SHA256(), backend=self.backend)
        elif algorithm == "SHA512":
            digest = hashes.Hash(hashes.SHA512(), backend=self.backend)
        else:
            raise ValueError(f"Unsupported hash algorithm: {algorithm}")

        digest.update(data)
        return digest.finalize()

    def base64_encode(self, data):
        """Encodes bytes data to Base64 string."""
        if not isinstance(data, bytes):
            data = data.encode('utf-8')
        return base64.urlsafe_b64encode(data).decode('utf-8')

    def base64_decode(self, encoded_data):
        """Decodes Base64 string to bytes data."""
        if not isinstance(encoded_data, str):
            raise ValueError("Encoded data must be a string.")
        return base64.urlsafe_b64decode(encoded_data.encode('utf-8'))

    def generate_rsa_key_pair(self, public_exponent=65537, key_size=2048):
        """Generates an RSA private and public key pair."""
        private_key = rsa.generate_private_key(
            public_exponent=public_exponent,
            key_size=key_size,
            backend=self.backend
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def rsa_encrypt(self, public_key, plaintext):
        """Encrypts data using RSA public key (OAEP padding)."""
        if not isinstance(plaintext, bytes):
            plaintext = plaintext.encode('utf-8')
        ciphertext = public_key.encrypt(
            plaintext,
            asymmetric_padding.OAEP(
                mgf=asymmetric_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
        return ciphertext

    def rsa_decrypt(self, private_key, ciphertext):
        """Decrypts data using RSA private key (OAEP padding)."""
        plaintext = private_key.decrypt(
            ciphertext,
            asymmetric_padding.OAEP(
                mgf=asymmetric_padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
        return plaintext

    def rsa_sign(self, private_key, data):
        """Signs data using RSA private key (PSS padding)."""
        if not isinstance(data, bytes):
            data = data.encode('utf-8')
        signature = private_key.sign(
            data,
            asymmetric_padding.PSS(
                mgf=asymmetric_padding.MGF1(hashes.SHA256()),
                salt_length=asymmetric_padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        return signature

    def rsa_verify(self, public_key, data, signature):
        """Verifies a signature using RSA public key."""
        if not isinstance(data, bytes):
            data = data.encode('utf-8')
        try:
            public_key.verify(
                signature,
                data,
                asymmetric_padding.PSS(
                    mgf=asymmetric_padding.MGF1(hashes.SHA256()),
                    salt_length=asymmetric_padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            return True
        except Exception:
            return False

    def hash_password(self, password):
        """Hashes a password using Scrypt."""
        if not isinstance(password, bytes):
            password = password.encode('utf-8')
        salt = os.urandom(16)
        kdf = Scrypt(
            salt=salt,
            length=32,
            n=2**14,
            r=8,
            p=1,
            backend=self.backend
        )
        hashed_password = kdf.derive(password)
        return salt + hashed_password

    def verify_password(self, password, hashed_password_with_salt):
        """Verifies a password against a Scrypt hash."""
        if not isinstance(password, bytes):
            password = password.encode('utf-8')
        if not isinstance(hashed_password_with_salt, bytes):
            raise ValueError("Hashed password (with salt) must be bytes.")

        salt_len = 16
        derived_key_len = 32
        if len(hashed_password_with_salt) != salt_len + derived_key_len:
            raise ValueError("Invalid hashed password format.")

        salt = hashed_password_with_salt[:salt_len]
        stored_derived_key = hashed_password_with_salt[salt_len:]

        kdf = Scrypt(
            salt=salt,
            length=derived_key_len,
            n=2**14,
            r=8,
            p=1,
            backend=self.backend
        )
        try:
            kdf.verify(password, stored_derived_key)
            return True
        except Exception:
            return False

    def derive_key_pbkdf2(self, password, salt=None, iterations=100000, key_length=32):
        """Derives a key from a password using PBKDF2-HMAC-SHA256."""
        if not isinstance(password, bytes):
            password = password.encode('utf-8')
        if salt is None:
            salt = os.urandom(16)

        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=key_length,
            salt=salt,
            iterations=iterations,
            backend=self.backend
        )
        derived_key = kdf.derive(password)
        return derived_key, salt

# Additional implementation at 2025-06-19 23:58:23
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
import bcrypt

class CryptoToolkit:
    def __init__(self):
        pass

    def generate_fernet_key(self) -> bytes:
        """Generates a Fernet key."""
        return Fernet.generate_key()

    def encrypt_fernet(self, data: bytes, key: bytes) -> bytes:
        """Encrypts data using Fernet."""
        f = Fernet(key)
        return f.encrypt(data)

    def decrypt_fernet(self, encrypted_data: bytes, key: bytes) -> bytes:
        """Decrypts data using Fernet."""
        f = Fernet(key)
        return f.decrypt(encrypted_data)

    def hash_sha256(self, data: bytes) -> str:
        """Hashes data using SHA256 and returns hex digest."""
        digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
        digest.update(data)
        return digest.finalize().hex()

    def generate_rsa_key_pair(self, key_size: int = 2048):
        """Generates an RSA private and public key pair."""
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=key_size,
            backend=default_backend()
        )
        public_key = private_key.public_key()
        return private_key, public_key

    def serialize_private_key(self, private_key, password: bytes = None) -> bytes:
        """Serializes a private key to PEM format."""
        if password:
            encryption_algorithm = serialization.BestAvailableEncryption(password)
        else:
            encryption_algorithm = serialization.NoEncryption()
        return private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=encryption_algorithm
        )

    def serialize_public_key(self, public_key) -> bytes:
        """Serializes a public key to PEM format."""
        return public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )

    def load_private_key_from_pem(self, pem_data: bytes, password: bytes = None):
        """Loads a private key from PEM format."""
        return serialization.load_pem_private_key(
            pem_data,
            password=password,
            backend=default_backend()
        )

    def load_public_key_from_pem(self, pem_data: bytes):
        """Loads a public key from PEM format."""
        return serialization.load_pem_public_key(
            pem_data,
            backend=default_backend()
        )

    def encrypt_rsa(self, public_key, plaintext: bytes) -> bytes:
        """Encrypts data using RSA public key."""
        return public_key.encrypt(
            plaintext,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def decrypt_rsa(self, private_key, ciphertext: bytes) -> bytes:
        """Decrypts data using RSA private key."""
        return private_key.decrypt(
            ciphertext,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )

    def sign_rsa(self, private_key, data: bytes) -> bytes:
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

    def verify_rsa_signature(self, public_key, data: bytes, signature: bytes) -> bool:
        """Verifies an RSA signature using public key."""
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

    def hash_password_bcrypt(self, password: str) -> bytes:
        """Hashes a password using bcrypt."""
        hashed = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
        return hashed

    def verify_password_bcrypt(self, password: str, hashed_password: bytes) -> bool:
        """Verifies a password against a bcrypt hash."""
        try:
            return bcrypt.checkpw(password.encode('utf-8'), hashed_password)
        except ValueError:
            return False

    def encrypt_file_fernet(self, filepath: str, key: bytes, output_filepath: str = None):
        """Encrypts a file using Fernet."""
        if output_filepath is None:
            output_filepath = filepath + ".enc"
        
        f = Fernet(key)
        with open(filepath, 'rb') as file:
            original_data = file.read()
        
        encrypted_data = f.encrypt(original_data)
        
        with open(output_filepath, 'wb') as file:
            file.write(encrypted_data)
        
        return output_filepath

    def decrypt_file_fernet(self, filepath: str, key: bytes, output_filepath: str = None):
        """Decrypts a file using Fernet."""
        if output_filepath is None:
            if filepath.endswith(".enc"):
                output_filepath = filepath[:-4]
            else:
                output_filepath = filepath + ".dec"
        
        f = Fernet(key)
        with open(filepath, 'rb') as file:
            encrypted_data = file.read()
        
        decrypted_data = f.decrypt(encrypted_data)
        
        with open(output_filepath, 'wb') as file:
            file.write(decrypted_data)
        
        return output_filepath