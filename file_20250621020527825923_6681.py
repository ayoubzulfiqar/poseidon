import json
import random
import string
import os

STORAGE_FILE = 'urls.json'
SHORT_CODE_LENGTH = 6

def _load_urls():
    if os.path.exists(STORAGE_FILE):
        with open(STORAGE_FILE, 'r') as f:
            try:
                return json.load(f)
            except json.JSONDecodeError:
                return {}
    return {}

def _save_urls(data):
    with open(STORAGE_FILE, 'w') as f:
        json.dump(data, f, indent=4)

def _generate_short_code(existing_codes):
    while True:
        code = ''.join(random.choices(string.ascii_letters + string.digits, k=SHORT_CODE_LENGTH))
        if code not in existing_codes:
            return code

def shorten_url(long_url):
    urls = _load_urls()
    for short_code, existing_long_url in urls.items():
        if existing_long_url == long_url:
            return short_code
    
    short_code = _generate_short_code(urls.keys())
    urls[short_code] = long_url
    _save_urls(urls)
    return short_code

def get_long_url(short_code):
    urls = _load_urls()
    return urls.get(short_code)

# Additional implementation at 2025-06-21 02:05:58
import json
import os
import random
import string
import datetime

class URLShortenerService:
    def __init__(self, storage_file='urls.json', base_domain='http://localhost:8000/'):
        self.storage_file = storage_file
        self.base_domain = base_domain
        self.urls = self._load_urls()
        self.short_code_length = 6

    def _load_urls(self):
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r') as f:
                    return json.load(f)
            except json.JSONDecodeError:
                return {}
        return {}

    def _save_urls(self):
        with open(self.storage_file, 'w') as f:
            json.dump(self.urls, f, indent=4)

    def _generate_short_code(self):
        characters = string.ascii_letters + string.digits
        while True:
            short_code = ''.join(random.choice(characters) for _ in range(self.short_code_length))
            if short_code not in self.urls:
                return short_code

    def shorten_url(self, long_url, custom_code=None, expires_at=None):
        if not long_url or not (long_url.startswith('http://') or long_url.startswith('https://')):
            raise ValueError("Invalid URL. Must start with http:// or https://")

        if custom_code:
            if not all(c.isalnum() for c in custom_code):
                raise ValueError("Custom short code must be alphanumeric.")
            if custom_code in self.urls:
                raise ValueError(f"Custom short code '{custom_code}' is already taken.")
            short_code = custom_code
        else:
            for code, data in self.urls.items():
                if data['long_url'] == long_url and not data.get('custom', False):
                    return self.base_domain + code
            short_code = self._generate_short_code()

        expiry_date_str = None
        if expires_at:
            if isinstance(expires_at, datetime.datetime):
                expiry_date_str = expires_at.isoformat()
            elif isinstance(expires_at, str):
                try:
                    datetime.datetime.fromisoformat(expires_at)
                    expiry_date_str = expires_at
                except ValueError:
                    raise ValueError("Invalid expires_at format. Use ISO format (YYYY-MM-DDTHH:MM:SS) or datetime object.")
            else:
                raise ValueError("expires_at must be a datetime object or an ISO formatted string.")

        self.urls[short_code] = {
            'long_url': long_url,
            'created_at': datetime.datetime.now().isoformat(),
            'clicks': 0,
            'expires_at': expiry_date_str,
            'custom': custom_code is not None
        }
        self._save_urls()
        return self.base_domain + short_code

    def get_long_url(self, short_code):
        data = self.urls.get(short_code)
        if not data:
            return None

        if data.get('expires_at'):
            expiry_time = datetime.datetime.fromisoformat(data['expires_at'])
            if datetime.datetime.now() > expiry_time:
                del self.urls[short_code]
                self._save_urls()
                return None

        data['clicks'] += 1
        self._save_urls()
        return data['long_url']

    def get_short_url_info(self, short_code):
        data = self.urls.get(short_code)
        if data:
            info = data.copy()
            info['short_url'] = self.base_domain + short_code
            return info
        return None

    def list_all_urls(self):
        all_info = []
        for short_code, data in self.urls.items():
            info = data.copy()
            info['short_code'] = short_code
            info['short_url'] = self.base_domain + short_code
            all_info.append(info)
        return all_info

    def delete_short_url(self, short_code):
        if short_code in self.urls:
            del self.urls[short_code]
            self._save_urls()
            return True
        return False

    def update_expiration(self, short_code, new_expires_at):
        data = self.urls.get(short_code)
        if not data:
            return False

        expiry_date_str = None
        if new_expires_at:
            if isinstance(new_expires_at, datetime.datetime):
                expiry_date_str = new_expires_at.isoformat()
            elif isinstance(new_expires_at, str):
                try:
                    datetime.datetime.fromisoformat(new_expires_at

# Additional implementation at 2025-06-21 02:06:50
import json
import os
import random
import string

class URLShortener:
    """
    A URL shortening service that stores mappings in a JSON file.
    """

    def __init__(self, storage_file='url_mappings.json'):
        self.storage_file = storage_file
        self.mappings = {}
        self._load_mappings()

    def _load_mappings(self):
        """Loads URL mappings from the storage file."""
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r') as f:
                    self.mappings = json.load(f)
            except json.JSONDecodeError:
                # Handle case where file is empty or contains invalid JSON
                print(f"Warning: '{self.storage_file}' is empty or corrupted. Starting with empty mappings.")
                self.mappings = {}
            except Exception as e:
                print(f"Error loading mappings: {e}. Starting with empty mappings.")
                self.mappings = {}
        else:
            # File does not exist, start with empty mappings
            self.mappings = {}

    def _save_mappings(self):
        """Saves current URL mappings to the storage file."""
        try:
            with open(self.storage_file, 'w') as f:
                json.dump(self.mappings, f, indent=4)
        except Exception as e:
            print(f"Error saving mappings: {e}")

    def _generate_short_code(self, length=6):
        """Generates a unique random short code."""
        characters = string.ascii_letters + string.digits
        current_length = length
        while True:
            short_code = ''.join(random.choices(characters, k=current_length))
            if short_code not in self.mappings:
                return short_code
            # If collision, try again with a new code, potentially increasing length
            current_length += 1 

    def shorten_url(self, long_url, custom_code=None):
        """
        Shortens a given URL.
        Args:
            long_url (str): The original long URL.
            custom_code (str, optional): A custom short code to use.
                                         If None, a random one is generated.
        Returns:
            str: The generated or custom short code if successful, None otherwise.
        """
        if not long_url or not isinstance(long_url, str):
            print("Error: Invalid long URL provided.")
            return None

        if custom_code:
            if not isinstance(custom_code, str) or not custom_code.isalnum():
                print("Error: Custom code must be an alphanumeric string.")
                return None
            if custom_code in self.mappings:
                print(f"Error: Custom code '{custom_code}' already exists. Please choose another.")
                return None
            short_code = custom_code
        else:
            short_code = self._generate_short_code()

        self.mappings[short_code] = long_url
        self._save_mappings()
        return short_code

    def retrieve_url(self, short_code):
        """
        Retrieves the original URL for a given short code.
        Args:
            short_code (str): The short code.
        Returns:
            str: The original long URL if found, None otherwise.
        """
        return self.mappings.get(short_code)

    def delete_mapping(self, short_code):
        """
        Deletes a URL mapping.
        Args:
            short_code (str): The short code to delete.
        Returns:
            bool: True if deleted, False if not found.
        """
        if short_code in self.mappings:
            del self.mappings[short_code]
            self._save_mappings()
            return True
        return False

    def list_all_mappings(self):
        """
        Lists all stored URL mappings.
        Returns:
            dict: A copy of the current mappings.
        """
        return self.mappings.copy()

def main():
    """
    Main function to run the command-line interface for the URL Shortener.
    """
    shortener = URLShortener()

    while True:
        print("\n--- URL Shortener Service ---")
        print("1. Shorten URL")
        print("2. Retrieve Original URL")
        print("3. List All Mappings")
        print("4. Delete Mapping")
        print("5. Exit")
        choice = input("Enter your choice: ")

        if choice == '1':
            long_url = input("Enter the long URL to shorten: ")
            custom_code_choice = input("Do you want to use a custom short code? (y/n): ").lower()
            custom_code = None
            if custom_code_choice == 'y':
                custom_code = input("Enter your desired custom short code (alphanumeric): ")
            
            short_code = shortener.shorten_url(long_url, custom_code)
            if short_code:
                print(f"URL shortened successfully! Short Code: {short_code}")
            else:
                print("Failed to shorten URL.")

        elif choice == '2':
            short_code = input("Enter the short code to retrieve: ")
            long_url = shortener.retrieve_url(short_code)
            if long_url:
                print(f"Original URL for '{short_code}': {long_url}")
            else:
                print(f"No URL found for short code '{short_code}'.")

        elif choice == '3':
            all_mappings = shortener.list_all_mappings()
            if all_mappings:
                print("\n--- All Mappings ---")
                for short_code, long_url in all_mappings.items():
                    print(f"  {short_code} -> {long_url}")
            else:
                print("No mappings found.")

        elif choice == '4':
            short_code = input("Enter the short code to delete: ")
            if shortener.delete_mapping(short_code):
                print(f"Mapping for '{short_code}' deleted successfully.")
            else:
                print(f"No mapping found for '{short_code}'.")

        elif choice == '5':
            print("Exiting URL Shortener. Goodbye!")
            break

        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 02:07:23
import json
import random
import string
import os

class URLShortener:
    def __init__(self, storage_file='urls.json', short_code_length=6):
        self.storage_file = storage_file
        self.short_code_length = short_code_length
        self.urls = self._load_urls()
        self.reverse_urls = {long_url: short_code for short_code, long_url in self.urls.items()}

    def _load_urls(self):
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r') as f:
                    data = json.load(f)
                    if isinstance(data, dict):
                        return data
                    else:
                        return {}
            except json.JSONDecodeError:
                return {}
            except Exception:
                return {}
        return {}

    def _save_urls(self):
        try:
            with open(self.storage_file, 'w') as f:
                json.dump(self.urls, f, indent=4)
        except IOError:
            pass

    def _generate_short_code(self):
        characters = string.ascii_letters + string.digits
        max_attempts = 1000
        attempts = 0
        while attempts < max_attempts:
            short_code = ''.join(random.choice(characters) for _ in range(self.short_code_length))
            if short_code not in self.urls:
                return short_code
            attempts += 1
        raise RuntimeError("Could not generate a unique short code after multiple attempts. Consider increasing short_code_length.")

    def shorten_url(self, long_url):
        if not isinstance(long_url, str) or not long_url.strip():
            raise ValueError("Long URL must be a non-empty string.")

        if long_url in self.reverse_urls:
            return self.reverse_urls[long_url]

        short_code = self._generate_short_code()
        self.urls[short_code] = long_url
        self.reverse_urls[long_url] = short_code
        self._save_urls()
        return short_code

    def retrieve_long_url(self, short_code):
        if not isinstance(short_code, str) or not short_code.strip():
            return None
        return self.urls.get(short_code)

    def get_all_mappings(self):
        return dict(self.urls)

    def delete_mapping(self, short_code):
        if short_code in self.urls:
            long_url = self.urls.pop(short_code)
            if long_url in self.reverse_urls:
                self.reverse_urls.pop(long_url)
            self._save_urls()
            return True
        return False

    def update_long_url(self, short_code, new_long_url):
        if not isinstance(new_long_url, str) or not new_long_url.strip():
            raise ValueError("New long URL must be a non-empty string.")

        if short_code in self.urls:
            old_long_url = self.urls[short_code]
            if old_long_url == new_long_url:
                return True

            if old_long_url in self.reverse_urls and self.reverse_urls[old_long_url] == short_code:
                del self.reverse_urls[old_long_url]

            self.urls[short_code] = new_long_url
            self.reverse_urls[new_long_url] = short_code
            self._save_urls()
            return True
        return False

if __name__ == '__main__':
    shortener = URLShortener(storage_file='my_urls.json')

    print("--- URL Shortener Service ---")

    url1 = "https://www.example.com/very/long/url/path/to/resource/1"
    url2 = "https://www.google.com/search?q=python+url+shortener&oq=python+url+shortener"
    url3 = "https://docs.python.org/3/library/json.html"
    url4 = "https://www.example.com/another/long/url"
    url5 = "http://localhost:8000/some/local/path"

    print(f"\nShortening '{url1}'...")
    try:
        short_code1 = shortener.shorten_url(url1)
        print(f"Shortened to: {short_code1}")
    except ValueError as e:
        print(f"Error: {e}")

    print(f"\nShortening '{url2}'...")
    try:
        short_code2 = shortener.shorten_url(url2)
        print(f"Shortened to: {short_code2}")
    except ValueError as e:
        print(f"Error: {e}")

    print(f"\nShortening '{url3}'...")
    try:
        short_code3 = shortener.shorten_url(url3)
        print(f"Shortened to: {short_code3}")
    except ValueError as e:
        print(f"Error: {e}")

    print(f"\nShortening '{url1}' again (should return existing):")
    try:
        short_code1_again = shortener.shorten_url(url1)
        print(f"Shortened to: {short_code1_again} (Matches original: {short_code1_again == short_code1})")
    except ValueError as e:
        print(f"Error: {e}")

    print(f"\nShortening '{url5}'...")
    try:
        short_code5 = shortener.shorten_url(url5)
        print(f"Shortened to: {short_code5}")
    except ValueError as e:
        print(f"Error: {e}")

    print(f"\nRetrieving original URL for '{short_code1}':")
    retrieved_url1 = shortener.retrieve_long_url(short_code1)
    print(f"Original URL: {retrieved_url1}")

    print(f"\nRetrieving original URL for '{short_code2}':")
    retrieved_url2 = shortener.retrieve_long_url(short_code2)
    print(f"Original URL: {retrieved_url2}")

    print(f"\nRetrieving original URL for a non-existent code 'XYZ123':")
    non_existent_url = shortener.retrieve_long_url("XYZ123")
    print(f"Original URL: {non_existent_url}")

    print("\n--- All Current Mappings ---")
    all_mappings = shortener.get_all_mappings()
    for sc, lu in all_mappings.items():
        print(f"  {sc} -> {lu}")

    print(f"\nAttempting to delete mapping for '{short_code2}'...")
    if shortener.delete_mapping(short_code2):
        print(f"Mapping for '{short_code2}' deleted successfully.")
    else:
        print(f"Mapping for '{short_code2}' not found.")

    print("\n--- Mappings After Deletion ---")
    all_mappings_after_delete = shortener.get_all_mappings()
    for sc, lu in all_mappings_after_delete.items():
        print(f"  {sc} -> {lu}")

    print(f"\nAttempting to update long URL for '{short_code1}' to '{url4}'...")
    try:
        if shortener.update_long_url(short_code1, url4):
            print(f"Mapping for '{short_code1}' updated successfully.")
        else:
            print(f"Mapping for '{short_code1}' not found for update.")
    except ValueError as e:
        print(f"Error updating: {e}")

    print("\n--- Mappings After Update ---")
    all_mappings_after_update = shortener.get_all_mappings()
    for sc, lu in all_mappings_after_update.items():
        print(f"  {sc} -> {lu}")

    print(f"\nRetrieving original URL for '{short_code1}' after update:")
    retrieved_url1_after_update = shortener.retrieve_long_url(short_code1)
    print(f"Original URL: {retrieved_url1_after_update}")

    print("\n--- Demonstrating Persistence (Loading from file) ---")
    new_shortener_instance = URLShortener