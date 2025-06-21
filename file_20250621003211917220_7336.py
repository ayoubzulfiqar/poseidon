import pyperclip
import re

def format_clipboard():
    original_text = pyperclip.paste()

    text = original_text.replace('\r\n', '\n')

    text = text.strip()

    text = re.sub(r'\n{3,}', '\n\n', text)

    lines = text.split('\n')
    stripped_lines = [line.strip() for line in lines]
    formatted_text = "\n".join(stripped_lines)

    formatted_text = re.sub(r'[ \t]+', ' ', formatted_text)

    pyperclip.copy(formatted_text)

if __name__ == '__main__':
    format_clipboard()

# Additional implementation at 2025-06-21 00:33:25
import sys
import json
import xml.etree.ElementTree as ET
from xml.dom import minidom
import urllib.parse
import base64
import re
import html

try:
    import pyperclip
except ImportError:
    print("pyperclip module not found. Please install it using 'pip install pyperclip'.", file=sys.stderr)
    sys.exit(1)

class ClipboardFormatter:
    def __init__(self):
        pass

    def _get_clipboard_content(self):
        try:
            return pyperclip.paste()
        except pyperclip.PyperclipException as e:
            print(f"Error accessing clipboard: {e}", file=sys.stderr)
            return None

    def _set_clipboard_content(self, text):
        try:
            pyperclip.copy(text)
            return True
        except pyperclip.PyperclipException as e:
            print(f"Error writing to clipboard: {e}", file=sys.stderr)
            return False

    def _format_and_copy(self, formatter_func):
        content = self._get_clipboard_content()
        if content is None:
            return False
        
        if not content.strip():
            print("Clipboard is empty or contains only whitespace. No formatting applied.")
            return False

        formatted_content = formatter_func(content)
        if formatted_content is None:
            return False
        
        if self._set_clipboard_content(formatted_content):
            print("Formatted content copied to clipboard.")
            return True
        return False

    def format_json_pretty(self, content):
        try:
            parsed_json = json.loads(content)
            return json.dumps(parsed_json, indent=4)
        except json.JSONDecodeError:
            print("Invalid JSON content on clipboard.", file=sys.stderr)
            return None

    def format_json_minify(self, content):
        try:
            parsed_json = json.loads(content)
            return json.dumps(parsed_json, separators=(',', ':'))
        except json.JSONDecodeError:
            print("Invalid JSON content on clipboard.", file=sys.stderr)
            return None

    def format_xml_pretty(self, content):
        try:
            root = ET.fromstring(content)
            rough_string = ET.tostring(root, 'utf-8')
            reparsed = minidom.parseString(rough_string)
            return reparsed.toprettyxml(indent="    ")
        except ET.ParseError:
            print("Invalid XML content on clipboard.", file=sys.stderr)
            return None
        except Exception as e:
            print(f"Error pretty printing XML: {e}", file=sys.stderr)
            return None

    def format_xml_minify(self, content):
        try:
            root = ET.fromstring(content)
            return ET.tostring(root, encoding='unicode')
        except ET.ParseError:
            print("Invalid XML content on clipboard.", file=sys.stderr)
            return None

    def format_url_encode(self, content):
        return urllib.parse.quote_plus(content)

    def format_url_decode(self, content):
        return urllib.parse.unquote_plus(content)

    def format_base64_encode(self, content):
        try:
            encoded_bytes = base64.b64encode(content.encode('utf-8'))
            return encoded_bytes.decode('utf-8')
        except Exception as e:
            print(f"Error encoding Base64: {e}", file=sys.stderr)
            return None

    def format_base64_decode(self, content):
        try:
            decoded_bytes = base64.b64decode(content.encode('utf-8'))
            return decoded_bytes.decode('utf-8')
        except (base64.binascii.Error, UnicodeDecodeError) as e:
            print(f"Invalid Base64 content or encoding: {e}", file=sys.stderr)
            return None

    def format_strip_whitespace(self, content):
        return content.strip()

    def format_to_uppercase(self, content):
        return content.upper()

    def format_to_lowercase(self, content):
        return content.lower()

    def format_to_titlecase(self, content):
        return content.title()

    def format_remove_empty_lines(self, content):
        lines = content.splitlines()
        non_empty_lines = [line for line in lines if line.strip()]
        return "\n".join(non_empty_lines)

    def format_sort_lines(self, content):
        lines = content.splitlines()
        lines.sort()
        return "\n".join(lines)

    def format_reverse_lines(self, content):
        lines = content.splitlines()
        lines.reverse()
        return "\n".join(lines)

    def format_remove_duplicate_lines(self, content):
        lines = content.splitlines()
        seen = set()
        unique_lines = []
        for line in lines:
            if line not in seen:
                unique_lines.append(line)
                seen.add(line)
        return "\n".join(unique_lines)

    def format_trim_line_whitespace(self, content):
        lines = content.splitlines()
        trimmed_lines = [line.strip() for line in lines]
        return "\n".join(trimmed_lines)

    def format_escape_html(self, content):
        return html.escape(content)

    def format_unescape_html(self, content):
        return html.unescape(content)

    def run_interactive_menu(self):
        menu_options = {
            '1': ("JSON Pretty Print", self.format_json_pretty),
            '2': ("JSON Minify", self.format_json_minify),
            '3': ("XML Pretty Print", self.format_xml_pretty),
            '4': ("XML Minify", self.format_xml_minify),
            '5': ("URL Encode", self.format_url_encode),
            '6': ("URL Decode", self.format_url_decode),
            '7': ("Base64 Encode", self.format_base64_encode),
            '8': ("Base64 Decode", self.format_base64_decode),
            '9': ("Strip Leading/Trailing Whitespace", self.format_strip_whitespace),
            '10': ("Convert to Uppercase", self.format_to_uppercase),
            '11': ("Convert to Lowercase", self.format_to_lowercase),
            '12': ("Convert to Title Case", self.format_to_titlecase),
            '13': ("Remove Empty Lines", self.format_remove_empty_lines),
            '14': ("Sort Lines Alphabetically", self.format_sort_lines),
            '15': ("Reverse Line Order", self.format_reverse_lines),
            '16': ("Remove Duplicate Lines", self.format_remove_duplicate_lines),
            '17': ("Trim Whitespace from Each Line", self.format_trim_line_whitespace),
            '18': ("Escape HTML Characters", self.format_escape_html),
            '19': ("Unescape HTML Characters", self.format_unescape_html),
            '0': ("Exit", None)
        }

        while True:
            print("\n--- Clipboard Formatter Menu ---")
            for key, (desc, _) in menu_options.items():
                print(f"{key}. {desc}")

            choice = input("Enter your choice: ").strip()

            if choice == '0':
                print("Exiting.")
                break
            elif choice in menu_options:
                description, formatter_func = menu_options[choice]
                if formatter_func:
                    print(f"Applying: {description}...")
                    self._format_and_copy(formatter_func)
            else:
                print("Invalid choice. Please try again.")

if __name__ == "__main__":
    formatter = ClipboardFormatter()
    formatter.run_interactive_menu()

# Additional implementation at 2025-06-21 00:34:45
import pyperclip
import textwrap

class ClipboardFormatter:
    @staticmethod
    def get_clipboard_content():
        try:
            return pyperclip.paste()
        except pyperclip.PyperclipException:
            return None

    @staticmethod
    def set_clipboard_content(text):
        try:
            pyperclip.copy(text)
        except pyperclip.PyperclipException:
            pass

    @staticmethod
    def strip_whitespace(text):
        lines = text.splitlines()
        processed_lines = []
        for line in lines:
            stripped_line = line.strip()
            cleaned_line = ' '.join(stripped_line.split())
            processed_lines.append(cleaned_line)
        return '\n'.join(processed_lines)

    @staticmethod
    def to_uppercase(text):
        return text.upper()

    @staticmethod
    def to_lowercase(text):
        return text.lower()

    @staticmethod
    def to_titlecase(text):
        return text.title()

    @staticmethod
    def remove_empty_lines(text):
        lines = text.splitlines()
        non_empty_lines = [line for line in lines if line.strip()]
        return '\n'.join(non_empty_lines)

    @staticmethod
    def sort_lines_alphabetically(text):
        lines = text.splitlines()
        lines.sort()
        return '\n'.join(lines)

    @staticmethod
    def add_prefix_to_lines(text, prefix=""):
        lines = text.splitlines()
        prefixed_lines = [prefix + line for line in lines]
        return '\n'.join(prefixed_lines)

    @staticmethod
    def add_suffix_to_lines(text, suffix=""):
        lines = text.splitlines()
        suffixed_lines = [line + suffix for line in lines]
        return '\n'.join(suffixed_lines)

    @staticmethod
    def replace_text(text, old_str, new_str):
        return text.replace(old_str, new_str)

    @staticmethod
    def wrap_lines(text, width=80):
        lines = text.splitlines()
        wrapped_lines = []
        for line in lines:
            wrapped_lines.extend(textwrap.wrap(line, width=width))
        return '\n'.join(wrapped_lines)

    def format_and_update(self, operations):
        original_content = self.get_clipboard_content()
        if original_content is None:
            return

        processed_content = original_content
        for op_name, op_kwargs in operations:
            method = getattr(self, op_name, None)
            if method and callable(method):
                try:
                    processed_content = method(processed_content, **op_kwargs)
                except Exception:
                    pass
        
        if processed_content != original_content:
            self.set_clipboard_content(processed_content)

# Additional implementation at 2025-06-21 00:36:05
import pyperclip
import re
import urllib.parse

class ClipboardFormatter:
    def __init__(self):
        pass

    def _get_clipboard_content(self):
        try:
            return pyperclip.paste()
        except pyperclip.PyperclipException:
            return ""

    def _set_clipboard_content(self, text):
        try:
            pyperclip.copy(text)
        except pyperclip.PyperclipException:
            pass

    def _trim_whitespace(self, text):
        lines = [line.strip() for line in text.splitlines()]
        return "\n".join(lines).strip()

    def _to_uppercase(self, text):
        return text.upper()

    def _to_lowercase(self, text):
        return text.lower()

    def _to_titlecase(self, text):
        return text.title()

    def _remove_extra_newlines(self, text):
        text = re.sub(r'\n{2,}', '\n', text)
        return text.strip('\n')

    def _tabs_to_spaces(self, text, num_spaces=4):
        return text.replace('\t', ' ' * num_spaces)

    def _remove_non_alphanumeric_except_spaces(self, text):
        return re.sub(r'[^a-zA-Z0-9\s]', '', text)

    def _url_encode(self, text):
        return urllib.parse.quote(text)

    def _url_decode(self, text):
        return urllib.parse.unquote(text)

    def format_and_copy(self, format_type, **kwargs):
        original_text = self._get_clipboard_content()
        formatted_text = original_text

        if not original_text:
            return

        if format_type == 'trim':
            formatted_text = self._trim_whitespace(original_text)
        elif format_type == 'uppercase':
            formatted_text = self._to_uppercase(original_text)
        elif format_type == 'lowercase':
            formatted_text = self._to_lowercase(original_text)
        elif format_type == 'titlecase':
            formatted_text = self._to_titlecase(original_text)
        elif format_type == 'remove_extra_newlines':
            formatted_text = self._remove_extra_newlines(original_text)
        elif format_type == 'tabs_to_spaces':
            num_spaces = kwargs.get('num_spaces', 4)
            formatted_text = self._tabs_to_spaces(original_text, num_spaces)
        elif format_type == 'remove_non_alphanumeric':
            formatted_text = self._remove_non_alphanumeric_except_spaces(original_text)
        elif format_type == 'url_encode':
            formatted_text = self._url_encode(original_text)
        elif format_type == 'url_decode':
            formatted_text = self._url_decode(original_text)
        else:
            return

        if formatted_text != original_text:
            self._set_clipboard_content(formatted_text)