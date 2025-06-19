import json
import secrets
import os

class URLShortener:
    def __init__(self, storage_file='urls.json'):
        self.storage_file = storage_file
        self.url_map = self._load_data()
        self.short_code_length = 7

    def _load_data(self):
        if os.path.exists(self.storage_file):
            with open(self.storage_file, 'r') as f:
                try:
                    return json.load(f)
                except json.JSONDecodeError:
                    return {}
        return {}

    def _save_data(self):
        with open(self.storage_file, 'w') as f:
            json.dump(self.url_map, f, indent=4)

    def _generate_short_code(self):
        while True:
            # Generate a URL-safe string from random bytes and take the first `short_code_length` characters
            # secrets.token_urlsafe(n_bytes) generates a string of length approx (n_bytes * 4 / 3)
            # We ensure a fixed length by slicing.
            short_code = secrets.token_urlsafe(self.short_code_length + 2)[:self.short_code_length]
            if short_code not in self.url_map:
                return short_code

    def shorten_url(self, long_url):
        # Check if the long_url already has a short code
        for short_code, existing_long_url in self.url_map.items():
            if existing_long_url == long_url:
                return short_code

        new_short_code = self._generate_short_code()
        self.url_map[new_short_code] = long_url
        self._save_data()
        return new_short_code

    def get_long_url(self, short_code):
        return self.url_map.get(short_code)

if __name__ == "__main__":
    shortener = URLShortener()

    # Example Usage
    print("--- URL Shortener Service Demonstration ---")

    # Shorten some URLs
    url1 = "https://www.example.com/very/long/url/path/to/resource/12345"
    short_code1 = shortener.shorten_url(url1)
    print(f"Shortened '{url1}' to: {short_code1}")

    url2 = "https://www.another-example.org/some/other/page/index.html"
    short_code2 = shortener.shorten_url(url2)
    print(f"Shortened '{url2}' to: {short_code2}")

    url3 = "https://www.google.com"
    short_code3 = shortener.shorten_url(url3)
    print(f"Shortened '{url3}' to: {short_code3}")

    # Try to shorten the same URL again (should return the existing short code)
    short_code1_again = shortener.shorten_url(url1)
    print(f"Shortened '{url1}' again (should be same): {short_code1_again}")

    print("\n--- Retrieving Long URLs ---")
    # Retrieve long URLs using their short codes
    retrieved_url1 = shortener.get_long_url(short_code1)
    print(f"Retrieved '{short_code1}': {retrieved_url1}")

    retrieved_url2 = shortener.get_long_url(short_code2)
    print(f"Retrieved '{short_code2}': {retrieved_url2}")

    retrieved_url3 = shortener.get_long_url(short_code3)
    print(f"Retrieved '{short_code3}': {retrieved_url3}")

    # Try to retrieve a non-existent short code
    non_existent_code = "invalid"
    retrieved_non_existent = shortener.get_long_url(non_existent_code)
    print(f"Retrieved '{non_existent_code}': {retrieved_non_existent} (Expected None)")

    print("\n--- Demonstrating Persistence (reloading service) ---")
    # Create a new shortener instance to demonstrate data persistence
    # The new instance will load data from the same 'urls.json' file
    new_shortener_instance = URLShortener()
    print(f"Retrieved '{short_code1}' with new instance: {new_shortener_instance.get_long_url(short_code1)}")
    print(f"Retrieved '{short_code2}' with new instance: {new_shortener_instance.get_long_url(short_code2)}")

    # Add a new URL with the new instance
    url4 = "https://www.github.com/my/awesome/repository/project"
    short_code4 = new_shortener_instance.shorten_url(url4)
    print(f"Shortened '{url4}' with new instance: {short_code4}")

    # Verify that the newly added URL is also persistent by creating yet another instance
    print("\n--- Verifying newly added URL in another reloaded instance ---")
    another_shortener_instance = URLShortener()
    print(f"Retrieved '{short_code4}' with another instance: {another_shortener_instance.get_long_url(short_code4)}")

    # Optional: Clean up the storage file after demonstration
    # if os.path.exists('urls.json'):
    #     os.remove('urls.json')
    #     print("\nCleaned up urls.json for next run.")

# Additional implementation at 2025-06-19 05:27:32
import json
import random
import string
import os

class URLShortener:
    def __init__(self, storage_file='urls.json', default_code_length=6):
        self.storage_file = storage_file
        self.default_code_length = default_code_length
        self.url_map = self._load_urls()

    def _load_urls(self):
        """Loads URL mappings from the storage file."""
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r') as f:
                    return json.load(f)
            except json.JSONDecodeError:
                # Handle empty or malformed JSON file
                print(f"Warning: {self.storage_file} is empty or malformed. Starting with an empty map.")
                return {}
        return {}

    def _save_urls(self):
        """Saves current URL mappings to the storage file."""
        with open(self.storage_file, 'w') as f:
            json.dump(self.url_map, f, indent=4)

    def _generate_short_code(self, length):
        """Generates a unique short code."""
        characters = string.ascii_letters + string.digits
        max_attempts = 1000 # Prevent infinite loop in case of extreme collision
        attempts = 0
        while attempts < max_attempts:
            code = ''.join(random.choice(characters) for _ in range(length))
            if code not in self.url_map:
                return code
            attempts += 1
        # Fallback if too many collisions for the given length
        print(f"Warning: Could not generate a unique short code of length {length} after {max_attempts} attempts. Consider increasing length or reducing collisions.")
        return None

    def shorten_url(self, long_url, custom_code=None):
        """
        Shortens a given long URL.
        Optionally accepts a custom short code.
        Returns the short code or None if custom code is taken or generation fails.
        """
        if not long_url:
            print("Error: Long URL cannot be empty.")
            return None

        # Check if the URL is already shortened
        for short_code, mapped_url in self.url_map.items():
            if mapped_url == long_url:
                print(f"URL already shortened: {long_url} -> {short_code}")
                return short_code

        if custom_code:
            if custom_code in self.url_map:
                print(f"Error: Custom short code '{custom_code}' is already taken.")
                return None
            short_code = custom_code
        else:
            short_code = self._generate_short_code(self.default_code_length)
            if short_code is None: # Failed to generate a unique code
                return None

        self.url_map[short_code] = long_url
        self._save_urls()
        print(f"URL shortened: {long_url} -> {short_code}")
        return short_code

    def retrieve_url(self, short_code):
        """
        Retrieves the original long URL for a given short code.
        Returns the long URL or None if not found.
        """
        if not short_code:
            print("Error: Short code cannot be empty.")
            return None

        long_url = self.url_map.get(short_code)
        if long_url:
            print(f"Retrieved: {short_code} -> {long_url}")
            return long_url
        else:
            print(f"Error: Short code '{short_code}' not found.")
            return None

    def list_all_urls(self):
        """Lists all stored short codes and their corresponding long URLs."""
        if not self.url_map:
            print("No URLs currently shortened.")
            return

        print("\n--- All Shortened URLs ---")
        for short_code, long_url in self.url_map.items():
            print(f"  {short_code} : {long_url}")
        print("--------------------------")

if __name__ == "__main__":
    shortener = URLShortener(storage_file='my_short_urls.json')

    print("--- URL Shortener Service Demonstration ---")

    # Shorten some URLs
    print("\n--- Shortening URLs ---")
    short_code1 = shortener.shorten_url("https://www.example.com/very/long/url/path/to/resource/1")
    short_code2 = shortener.shorten_url("https://www.google.com")
    short_code3 = shortener.shorten_url("https://www.python.org/docs/", custom_code="pyDocs")
    short_code4 = shortener.shorten_url("https://github.com/openai", custom_code="openaiGH")

    # Try to use an already taken custom code
    print("\n--- Attempting to use taken custom code ---")
    shortener.shorten_url("https://another.site.com", custom_code="pyDocs")

    # Try to shorten an already shortened URL
    print("\n--- Attempting to shorten an already shortened URL ---")
    shortener.shorten_url("https://www.google.com")

    # Retrieve URLs
    print("\n--- Retrieving URLs ---")
    shortener.retrieve_url(short_code1)
    shortener.retrieve_url("pyDocs")
    shortener.retrieve_url("openaiGH")
    shortener.retrieve_url("nonExistentCode")
    shortener.retrieve_url("") # Test empty short code

    # List all URLs
    shortener.list_all_urls()

    # Demonstrate persistence by creating a new instance
    print("\n--- Demonstrating Persistence ---")
    # The 'shortener' object will be garbage collected or go out of scope.
    # A new instance will load from the same file.
    persisted_shortener = URLShortener(storage_file='my_short_urls.json')
    persisted_shortener.list_all_urls()

    # Clean up the storage file for repeated runs (optional, uncomment to enable)
    # if os.path.exists('my_short_urls.json'):
    #     os.remove('my_short_urls.json')
    #     print("\nCleaned up 'my_short_urls.json'")
    # else:
    #     print("\nNo 'my_short_urls.json' file to clean up.")

# Additional implementation at 2025-06-19 05:28:36
import json
import random
import string
import os
import re
import datetime

class URLShortener:
    def __init__(self, storage_file="urls.json", base_domain="http://short.url/"):
        self.storage_file = storage_file
        self.base_domain = base_domain
        self.urls = self._load_data()
        self.short_code_length = 6

    def _load_data(self):
        if os.path.exists(self.storage_file):
            with open(self.storage_file, 'r') as f:
                try:
                    return json.load(f)
                except json.JSONDecodeError:
                    return {}
        return {}

    def _save_data(self):
        with open(self.storage_file, 'w') as f:
            json.dump(self.urls, f, indent=4)

    def _generate_short_code(self):
        characters = string.ascii_letters + string.digits
        while True:
            short_code = ''.join(random.choice(characters) for _ in range(self.short_code_length))
            if short_code not in self.urls:
                return short_code

    def _is_valid_url(self, url):
        regex = re.compile(
            r'^(?:http|ftp)s?://'
            r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|'
            r'localhost|'
            r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'
            r'(?::\d+)?'
            r'(?:/?|[/?]\S+)$', re.IGNORECASE)
        return re.match(regex, url) is not None

    def shorten_url(self, long_url):
        if not self._is_valid_url(long_url):
            raise ValueError("Invalid URL provided.")

        for short_code, data in self.urls.items():
            if data['long_url'] == long_url:
                return self.base_domain + short_code

        short_code = self._generate_short_code()
        self.urls[short_code] = {
            'long_url': long_url,
            'clicks': 0,
            'created_at': datetime.datetime.now().isoformat()
        }
        self._save_data()
        return self.base_domain + short_code

    def custom_shorten_url(self, long_url, custom_code):
        if not self._is_valid_url(long_url):
            raise ValueError("Invalid URL provided.")
        if not custom_code or not all(c.isalnum() for c in custom_code):
            raise ValueError("Custom code must be alphanumeric and non-empty.")
        if custom_code in self.urls:
            raise ValueError(f"Custom code '{custom_code}' is already in use.")

        self.urls[custom_code] = {
            'long_url': long_url,
            'clicks': 0,
            'created_at': datetime.datetime.now().isoformat()
        }
        self._save_data()
        return self.base_domain + custom_code

    def get_long_url(self, short_code):
        if short_code in self.urls:
            self.urls[short_code]['clicks'] += 1
            self._save_data()
            return self.urls[short_code]['long_url']
        return None

    def get_stats(self, short_code):
        if short_code in self.urls:
            return self.urls[short_code]
        return None

    def get_all_urls(self):
        return self.urls

if __name__ == "__main__":
    shortener = URLShortener()

    long_url_1 = "https://www.example.com/very/long/url/path/to/resource/page"
    short_url_1 = shortener.shorten_url(long_url_1)
    print(f"Shortened URL: {short_url_1}")

    long_url_2 = "https://docs.python.org/3/library/json.html"
    short_url_2 = shortener.shorten_url(long_url_2)
    print(f"Shortened URL: {short_url_2}")

    short_url_1_again = shortener.shorten_url(long_url_1)
    print(f"Shortened URL again: {short_url_1_again}")

    try:
        custom_url = shortener.custom_shorten_url("https://github.com/python", "pygit")
        print(f"Custom URL: {custom_url}")
    except ValueError as e:
        print(f"Error creating custom URL: {e}")

    retrieved_url_1 = shortener.get_long_url(short_url_1.split('/')[-1])
    print(f"Retrieved URL for {short_url_1.split('/')[-1]}: {retrieved_url_1}")

    retrieved_url_1 = shortener.get_long_url(short_url_1.split('/')[-1])
    print(f"Retrieved URL for {short_url_1.split('/')[-1]}: {retrieved_url_1}")

    retrieved_url_custom = shortener.get_long_url(custom_url.split('/')[-1])
    print(f"Retrieved URL for {custom_url.split('/')[-1]}: {retrieved_url_custom}")

    stats_1 = shortener.get_stats(short_url_1.split('/')[-1])
    print(f"Stats for {short_url_1.split('/')[-1]}: {stats_1}")

    stats_custom = shortener.get_stats(custom_url.split('/')[-1])
    print(f"Stats for {custom_url.split('/')[-1]}: {stats_custom}")

    non_existent_url = shortener.get_long_url("nonexistent")
    print(f"Retrieved URL for 'nonexistent': {non_existent_url}")

    try:
        shortener.shorten_url("invalid-url")
    except ValueError as e:
        print(f"Error shortening invalid URL: {e}")

    try:
        shortener.custom_shorten_url("https://another.com", "pygit")
    except ValueError as e:
        print(f"Error creating custom URL with existing code: {e}")

    all_urls = shortener.get_all_urls()
    print(f"All stored URLs: {all_urls}")