import pyperclip
import re

def format_clipboard_content():
    try:
        original_text = pyperclip.paste()
        
        if not original_text:
            return

        cleaned_text = original_text.replace('\xa0', ' ')
        cleaned_text = re.sub(r'[ \t]+', ' ', cleaned_text)
        cleaned_text = cleaned_text.replace('\r\n', '\n').replace('\r', '\n')
        cleaned_text = re.sub(r'\n+', '\n', cleaned_text)
        cleaned_text = cleaned_text.strip()

        pyperclip.copy(cleaned_text)

    except ImportError:
        pass
    except Exception:
        pass

if __name__ == "__main__":
    format_clipboard_content()

# Additional implementation at 2025-06-21 00:19:18
import pyperclip
import json
import re

class ClipboardFormatter:
    def __init__(self):
        self.original_content = self._get_clipboard_content()
        self.current_content = self.original_content

    def _get_clipboard_content(self):
        try:
            return pyperclip.paste()
        except pyperclip.PyperclipException:
            print("Clipboard access failed. Please ensure you have a clipboard utility installed (e.g., xclip, xsel on Linux, or ensure macOS/Windows clipboard is accessible).")
            return ""

    def _set_clipboard_content(self, text):
        try:
            pyperclip.copy(text)
            print("Formatted text copied to clipboard.")
        except pyperclip.PyperclipException:
            print("Failed to copy text to clipboard.")

    def _update_content(self, new_text):
        self.current_content = new_text
        return new_text

    def trim_whitespace(self):
        """Removes leading/trailing whitespace from each line."""
        lines = self.current_content.splitlines()
        trimmed_lines = [line.strip() for line in lines]
        return self._update_content('\n'.join(trimmed_lines))

    def normalize_whitespace(self):
        """Replaces multiple spaces/tabs with single spaces and newlines with spaces."""
        text = self.current_content
        text = re.sub(r'\s+', ' ', text).strip()
        return self._update_content(text)

    def to_uppercase(self):
        """Converts all text to uppercase."""
        return self._update_content(self.current_content.upper())

    def to_lowercase(self):
        """Converts all text to lowercase."""
        return self._update_content(self.current_content.lower())

    def to_titlecase(self):
        """Converts text to title case (first letter of each word capitalized)."""
        return self._update_content(self.current_content.title())

    def add_prefix(self, prefix):
        """Adds a prefix to each line."""
        lines = self.current_content.splitlines()
        prefixed_lines = [prefix + line for line in lines]
        return self._update_content('\n'.join(prefixed_lines))

    def add_suffix(self, suffix):
        """Adds a suffix to each line."""
        lines = self.current_content.splitlines()
        suffixed_lines = [line + suffix for line in lines]
        return self._update_content('\n'.join(suffixed_lines))

    def find_replace(self, old_str, new_str):
        """Finds and replaces all occurrences of a substring."""
        return self._update_content(self.current_content.replace(old_str, new_str))

    def pretty_print_json(self):
        """Pretty prints JSON content."""
        try:
            data = json.loads(self.current_content)
            pretty_json = json.dumps(data, indent=4)
            return self._update_content(pretty_json)
        except json.JSONDecodeError:
            print("Error: Content is not valid JSON. Cannot pretty print.")
            return self.current_content # Return original content if invalid

    def minify_json(self):
        """Minifies JSON content."""
        try:
            data = json.loads(self.current_content)
            minified_json = json.dumps(data, separators=(',', ':'))
            return self._update_content(minified_json)
        except json.JSONDecodeError:
            print("Error: Content is not valid JSON. Cannot minify.")
            return self.current_content # Return original content if invalid

    def reset_content(self):
        """Resets the current content to the original clipboard content."""
        self.current_content = self.original_content
        print("Content reset to original clipboard content.")
        return self.current_content

def main():
    formatter = ClipboardFormatter()

    if not formatter.original_content:
        print("Clipboard is empty. Nothing to format.")
        return

    print("\n--- Clipboard Formatter ---")
    print("Original clipboard content (first 100 chars):")
    print(formatter.original_content[:100] + ("..." if len(formatter.original_content) > 100 else ""))

    while True:
        print("\nSelect an operation:")
        print(" 1. Trim leading/trailing whitespace (each line)")
        print(" 2. Normalize all whitespace (multiple spaces/newlines to single space)")
        print(" 3. Convert to UPPERCASE")
        print(" 4. Convert to lowercase")
        print(" 5. Convert to Title Case")
        print(" 6. Add prefix to each line")
        print(" 7. Add suffix to each line")
        print(" 8. Find and Replace text")
        print(" 9. Pretty Print JSON")
        print("10. Minify JSON")
        print("11. Reset content to original clipboard")
        print("12. Apply changes and Exit (copy to clipboard)")
        print("13. Exit without applying changes")

        choice = input("Enter choice (1-13): ")

        if choice == '1':
            formatter.trim_whitespace()
        elif choice == '2':
            formatter.normalize_whitespace()
        elif choice == '3':
            formatter.to_uppercase()
        elif choice == '4':
            formatter.to_lowercase()
        elif choice == '5':
            formatter.to_titlecase()
        elif choice == '6':
            prefix = input("Enter prefix: ")
            formatter.add_prefix(prefix)
        elif choice == '7':
            suffix = input("Enter suffix: ")
            formatter.add_suffix(suffix)
        elif choice == '8':
            old_str = input("Enter text to find: ")
            new_str = input("Enter text to replace with: ")
            formatter.find_replace(old_str, new_str)
        elif choice == '9':
            formatter.pretty_print_json()
        elif choice == '10':
            formatter.minify_json()
        elif choice == '11':
            formatter.reset_content()
        elif choice == '12':
            formatter._set_clipboard_content(formatter.current_content)
            print("Exiting.")
            break
        elif choice == '13':
            print("Exiting without applying changes.")
            break
        else:
            print("Invalid choice. Please try again.")

        print("\n--- Current Formatted Content (first 100 chars) ---")
        print(formatter.current_content[:100] + ("..." if len(formatter.current_content) > 100 else ""))

if __name__ == "__main__":
    main()