import json
import os
import random
import string
from flask import Flask, request, redirect, url_for

app = Flask(__name__)

STORAGE_FILE = 'urls.json'
BASE_URL = 'http://127.0.0.1:5000/'

def load_urls():
    if os.path.exists(STORAGE_FILE):
        with open(STORAGE_FILE, 'r') as f:
            try:
                return json.load(f)
            except json.JSONDecodeError:
                return {}
    return {}

def save_urls(data):
    with open(STORAGE_FILE, 'w') as f:
        json.dump(data, f, indent=4)

def generate_short_code(length=6):
    characters = string.ascii_letters + string.digits
    urls = load_urls()
    while True:
        short_code = ''.join(random.choice(characters) for _ in range(length))
        if short_code not in urls:
            return short_code

@app.route('/', methods=['GET'])
def index():
    return """
    <!DOCTYPE html>
    <html>
    <head>
        <title>URL Shortener</title>
        <style>
            body { font-family: sans-serif; margin: 50px; }
            form { margin-top: 20px; }
            input[type="text"] { width: 400px; padding: 8px; }
            input[type="submit"] { padding: 8px 15px; cursor: pointer; }
            .result { margin-top: 20px; font-size: 1.1em; }
        </style>
    </head>
    <body>
        <h1>URL Shortener Service</h1>
        <form action="/shorten" method="post">
            <label for="long_url">Enter Long URL:</label><br>
            <input type="text" id="long_url" name="long_url" placeholder="e.g., https://www.example.com/very/long/path" required><br><br>
            <input type="submit" value="Shorten URL">
        </form>
        <div class="result" id="result"></div>
        <script>
            document.querySelector('form').addEventListener('submit', async function(event) {
                event.preventDefault();
                const longUrl = document.getElementById('long_url').value;
                const response = await fetch('/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `long_url=${encodeURIComponent(longUrl)}`
                });
                const data = await response.json();
                const resultDiv = document.getElementById('result');
                if (response.ok) {
                    resultDiv.innerHTML = `Shortened URL: <a href="${data.short_url}" target="_blank">${data.short_url}</a>`;
                } else {
                    resultDiv.innerHTML = `Error: ${data.error}`;
                }
            });
        </script>
    </body>
    </html>
    """

@app.route('/shorten', methods=['POST'])
def shorten_url():
    long_url = request.form.get('long_url')
    if not long_url:
        return {"error": "Missing long_url parameter"}, 400

    urls = load_urls()

    for short_code, existing_long_url in urls.items():
        if existing_long_url == long_url:
            return {"short_url": BASE_URL + short_code, "message": "URL already shortened"}, 200

    short_code = generate_short_code()
    urls[short_code] = long_url
    save_urls(urls)

    return {"short_url": BASE_URL + short_code}, 200

@app.route('/<short_code>')
def redirect_to_long_url(short_code):
    urls = load_urls()
    long_url = urls.get(short_code)
    if long_url:
        return redirect(long_url)
    return "URL not found", 404

if __name__ == '__main__':
    app.run(debug=True)

# Additional implementation at 2025-06-17 23:30:56
import json
import random
import string
import os
from urllib.parse import urlparse

class URLShortener:
    def __init__(self, storage_file='urls.json', short_code_length=6):
        self.storage_file = storage_file
        self.short_code_length = short_code_length
        self.short_to_long_map = {}
        self.long_to_short_map = {}
        self._load_urls()

    def _generate_short_code(self):
        characters = string.ascii_letters + string.digits
        while True:
            short_code = ''.join(random.choice(characters) for _ in range(self.short_code_length))
            if short_code not in self.short_to_long_map:
                return short_code

    def _is_valid_url(self, url):
        try:
            result = urlparse(url)
            return all([result.scheme, result.netloc])
        except ValueError:
            return False

    def _load_urls(self):
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r') as f:
                    data = json.load(f)
                    self.short_to_long_map = data
                    self.long_to_short_map = {v: k for k, v in data.items()}
            except json.JSONDecodeError:
                pass
            except Exception:
                pass

    def _save_urls(self):
        try:
            with open(self.storage_file, 'w') as f:
                json.dump(self.short_to_long_map, f, indent=4)
        except Exception:
            pass

    def shorten_url(self, long_url):
        if not self._is_valid_url(long_url):
            raise ValueError("Invalid URL provided.")

        if long_url in self.long_to_short_map:
            return self.long_to_short_map[long_url]

        short_code = self._generate_short_code()
        self.short_to_long_map[short_code] = long_url
        self.long_to_short_map[long_url] = short_code
        self._save_urls()
        return short_code

    def retrieve_url(self, short_code):
        return self.short_to_long_map.get(short_code)

    def list_all_urls(self):
        return self.short_to_long_map.copy()

    def delete_short_url(self, short_code):
        if short_code in self.short_to_long_map:
            long_url = self.short_to_long_map.pop(short_code)
            if long_url in self.long_to_short_map:
                del self.long_to_short_map[long_url]
            self._save_urls()
            return True
        return False

if __name__ == '__main__':
    shortener = URLShortener()

    try:
        short_code1 = shortener.shorten_url("https://www.example.com/very/long/url/path/to/resource/1")
        print(f"Shortened: https://www.example.com/very/long/url/path/to/resource/1 -> {short_code1}")

        short_code2 = shortener.shorten_url("https://docs.python.org/3/library/urllib.parse.html")
        print(f"Shortened: https://docs.python.org/3/library/urllib.parse.html -> {short_code2}")

        short_code3 = shortener.shorten_url("https://www.example.com/very/long/url/path/to/resource/1")
        print(f"Shortened (again): https://www.example.com/very/long/url/path/to/resource/1 -> {short_code3}")

        # shortener.shorten_url("invalid-url-format") # This would raise ValueError
    except ValueError as e:
        print(f"Error shortening URL: {e}")

    print("\nRetrieving URLs:")
    retrieved_url1 = shortener.retrieve_url(short_code1)
    print(f"Retrieved {short_code1}: {retrieved_url1}")

    retrieved_url2 = shortener.retrieve_url(short_code2)
    print(f"Retrieved {short_code2}: {retrieved_url2}")

    non_existent_code = "NONEXIST"
    retrieved_non_existent = shortener.retrieve_url(non_existent_code)
    print(f"Retrieved {non_existent_code}: {retrieved_non_existent}")

    print("\nAll current short URLs:")
    for short, long in shortener.list_all_urls().items():
        print(f"  {short}: {long}")

    print(f"\nAttempting to delete {short_code1}...")
    if shortener.delete_short_url(short_code1):
        print(f"Successfully deleted {short_code1}.")
    else:
        print(f"Failed to delete {short_code1}.")

    print("\nAll short URLs after deletion:")
    for short, long in shortener.list_all_urls().items():
        print(f"  {short}: {long}")

    print("\nDemonstrating persistence (loading new instance):")
    new_shortener_instance = URLShortener()
    print("All URLs in new instance:")
    for short, long in new_shortener_instance.list_all_urls().items():
        print(f"  {short}: {long}")

    # Optional: Clean up the storage file after demonstration
    # if os.path.exists('urls.json'):
    #     os.remove('urls.json')
    #     print("\nCleaned up urls.json")

# Additional implementation at 2025-06-17 23:32:10
import json
import os
import random
import string
import re

class URLShortener:
    def __init__(self, storage_file='urls.json'):
        self.storage_file = storage_file
        self.urls = self._load_urls()
        self.short_code_length = 6

    def _load_urls(self):
        if os.path.exists(self.storage_file):
            with open(self.storage_file, 'r') as f:
                try:
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

    def is_valid_url(self, url):
        regex = re.compile(
            r'^(?:http|ftp)s?://'
            r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|'
            r'localhost|'
            r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'
            r'(?::\d+)?'
            r'(?:/?|[/?]\S+)$', re.IGNORECASE)
        return re.match(regex, url) is not None

    def shorten_url(self, long_url, custom_code=None):
        if not self.is_valid_url(long_url):
            return {"error": "Invalid URL format."}

        for short, long in self.urls.items():
            if long == long_url:
                if custom_code and short != custom_code:
                    return {"error": f"URL already shortened to '{short}'. Cannot assign new custom code."}
                return {"short_code": short, "message": "URL already exists."}

        if custom_code:
            if custom_code in self.urls:
                return {"error": f"Custom code '{custom_code}' is already in use."}
            if not re.match(r'^[a-zA-Z0-9_-]+$', custom_code):
                return {"error": "Custom code can only contain letters, numbers, underscores, and hyphens."}
            short_code = custom_code
        else:
            short_code = self._generate_short_code()

        self.urls[short_code] = long_url
        self._save_urls()
        return {"short_code": short_code, "long_url": long_url}

    def retrieve_url(self, short_code):
        return self.urls.get(short_code)

    def list_all_urls(self):
        return self.urls

if __name__ == "__main__":
    shortener = URLShortener()

    print("URL Shortener Service")
    print("---------------------")

    while True:
        print("\nOptions:")
        print("1. Shorten a URL")
        print("2. Shorten a URL with a custom code")
        print("3. Retrieve original URL")
        print("4. List all shortened URLs")
        print("5. Exit")

        choice = input("Enter your choice: ")

        if choice == '1':
            long_url = input("Enter the long URL to shorten: ")
            result = shortener.shorten_url(long_url)
            if "error" in result:
                print(f"Error: {result['error']}")
            elif "message" in result:
                print(f"Success: {result['message']} Short code: {result['short_code']}")
            else:
                print(f"URL shortened! Short code: {result['short_code']}")
        elif choice == '2':
            long_url = input("Enter the long URL to shorten: ")
            custom_code = input("Enter your desired custom short code: ")
            result = shortener.shorten_url(long_url, custom_code=custom_code)
            if "error" in result:
                print(f"Error: {result['error']}")
            elif "message" in result:
                print(f"Success: {result['message']} Short code: {result['short_code']}")
            else:
                print(f"URL shortened with custom code! Short code: {result['short_code']}")
        elif choice == '3':
            short_code = input("Enter the short code to retrieve: ")
            long_url = shortener.retrieve_url(short_code)
            if long_url:
                print(f"Original URL for '{short_code}': {long_url}")
            else:
                print(f"Short code '{short_code}' not found.")
        elif choice == '4':
            all_urls = shortener.list_all_urls()
            if all_urls:
                print("All Shortened URLs:")
                for short, long in all_urls.items():
                    print(f"  {short} -> {long}")
            else:
                print("No URLs shortened yet.")
        elif choice == '5':
            print("Exiting URL Shortener. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")