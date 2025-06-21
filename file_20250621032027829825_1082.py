import pyperclip
import time
import threading
import collections
import keyboard

clipboard_history = collections.deque(maxlen=20)
current_clipboard_content = ""
history_index = -1
is_programmatic_paste = False

def monitor_clipboard():
    global current_clipboard_content, history_index, is_programmatic_paste
    while True:
        try:
            new_content = pyperclip.paste()
            if new_content != current_clipboard_content:
                if is_programmatic_paste:
                    is_programmatic_paste = False
                    current_clipboard_content = new_content
                elif new_content.strip() != "":
                    if not clipboard_history or clipboard_history[0] != new_content:
                        clipboard_history.appendleft(new_content)
                        print(f"Added to history: '{new_content[:50]}...'")
                    current_clipboard_content = new_content
                    history_index = -1
            time.sleep(0.5)
        except pyperclip.PyperclipException as e:
            print(f"Clipboard access error: {e}. Retrying...")
            time.sleep(1)
        except Exception as e:
            print(f"An unexpected error occurred in monitor_clipboard: {e}")
            time.sleep(1)

def paste_from_history(direction):
    global history_index, is_programmatic_paste
    if not clipboard_history:
        print("Clipboard history is empty.")
        return

    if history_index == -1:
        history_index = 0
    else:
        if direction == "next":
            history_index = (history_index + 1) % len(clipboard_history)
        elif direction == "prev":
            history_index = (history_index - 1 + len(clipboard_history)) % len(clipboard_history)

    item_to_paste = clipboard_history[history_index]
    print(f"Pasting from history (index {history_index}): '{item_to_paste[:50]}...'")

    is_programmatic_paste = True
    pyperclip.copy(item_to_paste)
    
    time.sleep(0.05) 
    
    keyboard.press_and_release('ctrl+v')

def show_history():
    if not clipboard_history:
        print("Clipboard history is empty.")
        return
    print("\n--- Clipboard History ---")
    for i, item in enumerate(clipboard_history):
        display_item = item.replace('\n', '\\n').replace('\r', '\\r')
        if len(display_item) > 70:
            display_item = display_item[:67] + "..."
        print(f"{i}: {display_item}")
    print("-------------------------\n")

def main():
    print("Clipboard History Manager started.")
    print("Hotkeys:")
    print("  Ctrl+Shift+V: Paste previous history item (cycles backward).")
    print("  Ctrl+Shift+N: Paste next history item (cycles forward).")
    print("  Ctrl+Shift+H: Show clipboard history in console.")
    print("  Ctrl+Shift+Q: Quit the manager.")

    monitor_thread = threading.Thread(target=monitor_clipboard, daemon=True)
    monitor_thread.start()

    keyboard.add_hotkey('ctrl+shift+v', lambda: paste_from_history("prev"))
    keyboard.add_hotkey('ctrl+shift+n', lambda: paste_from_history("next"))
    keyboard.add_hotkey('ctrl+shift+h', show_history)
    keyboard.add_hotkey('ctrl+shift+q', lambda: keyboard.unhook_all_hotkeys() or exit())

    keyboard.wait('ctrl+shift+q')

    print("Clipboard History Manager stopped.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 03:21:45
import pyperclip
import json
import time
import threading
import os

class ClipboardManager:
    def __init__(self, history_file="clipboard_history.json", max_history=100):
        self.history_file = history_file
        self.max_history = max_history
        self.history = []
        self._load_history()
        self._last_clipboard_content = pyperclip.paste()
        self._monitor_thread = None
        self._running = False

    def _load_history(self):
        if os.path.exists(self.history_file):
            try:
                with open(self.history_file, 'r', encoding='utf-8') as f:
                    self.history = json.load(f)
            except (json.JSONDecodeError, FileNotFoundError):
                self.history = []
        else:
            self.history = []

    def _save_history(self):
        try:
            with open(self.history_file, 'w', encoding='utf-8') as f:
                json.dump(self.history, f, indent=4)
        except IOError as e:
            print(f"Error saving history: {e}")

    def add_item(self, text):
        if not text or text.strip() == "":
            return False
        
        # Remove duplicates (most recent first)
        if text in self.history:
            self.history.remove(text)
        
        self.history.insert(0, text) # Add to the beginning

        # Trim history if it exceeds max_history
        if len(self.history) > self.max_history:
            self.history = self.history[:self.max_history]
        
        self._save_history()
        return True

    def get_history(self, search_term=None):
        if search_term:
            return [item for item in self.history if search_term.lower() in item.lower()]
        return list(self.history) # Return a copy

    def get_item(self, index):
        try:
            return self.history[index]
        except IndexError:
            return None

    def delete_item(self, index):
        try:
            deleted_item = self.history.pop(index)
            self._save_history()
            return deleted_item
        except IndexError:
            return None

    def clear_history(self):
        self.history = []
        self._save_history()

    def _monitor_clipboard(self):
        while self._running:
            current_clipboard_content = pyperclip.paste()
            if current_clipboard_content != self._last_clipboard_content:
                if self.add_item(current_clipboard_content):
                    # This print statement might interfere with user input in the console
                    # For a cleaner CLI, consider logging or a non-blocking notification
                    pass 
                self._last_clipboard_content = current_clipboard_content
            time.sleep(0.5) # Check every 0.5 seconds

    def start_monitoring(self):
        if not self._running:
            self._running = True
            self._monitor_thread = threading.Thread(target=self._monitor_clipboard, daemon=True)
            self._monitor_thread.start()
            print("Clipboard monitoring started.")

    def stop_monitoring(self):
        if self._running:
            self._running = False
            if self._monitor_thread and self._monitor_thread.is_alive():
                self._monitor_thread.join(timeout=1) # Give it a moment to stop
            print("Clipboard monitoring stopped.")

def display_history(manager, search_term=None):
    history_items = manager.get_history(search_term)
    if not history_items:
        print("History is empty." if not search_term else "No items found matching your search.")
        return False
    
    print("\n--- Clipboard History ---")
    for i, item in enumerate(history_items):
        display_item = item.replace('\n', '\\n') # Replace newlines for cleaner display
        print(f"{i+1}. {display_item[:70]}{'...' if len(display_item) > 70 else ''}")
    print("-------------------------")
    return True

def main():
    manager = ClipboardManager()
    manager.start_monitoring()

    while True:
        print("\n--- Clipboard Manager Menu ---")
        print("1. Show History")
        print("2. Paste Item to Clipboard")
        print("3. Clear History")
        print("4. Search History")
        print("5. Delete Item from History")
        print("6. Exit")
        choice = input("Enter your choice: ").strip()

        if choice == '1':
            display_history(manager)
        elif choice == '2':
            if not display_history(manager):
                continue
            try:
                idx_str = input("Enter the number of the item to paste: ").strip()
                idx = int(idx_str) - 1
                item_to_paste = manager.get_item(idx)
                if item_to_paste:
                    pyperclip.copy(item_to_paste)
                    print(f"Item '{item_to_paste[:50]}...' copied to clipboard.")
                else:
                    print("Invalid item number.")
            except ValueError:
                print("Invalid input. Please enter a number.")
        elif choice == '3':
            confirm = input("Are you sure you want to clear all history? (yes/no): ").strip().lower()
            if confirm == 'yes':
                manager.clear_history()
                print("Clipboard history cleared.")
            else:
                print("Operation cancelled.")
        elif choice == '4':
            search_term = input("Enter search term: ").strip()
            display_history(manager, search_term)
        elif choice == '5':
            if not display_history(manager):
                continue
            try:
                idx_str = input("Enter the number of the item to delete: ").strip()
                idx = int(idx_str) - 1
                deleted_item = manager.delete_item(idx)
                if deleted_item:
                    print(f"Item '{deleted_item[:50]}...' deleted from history.")
                else:
                    print("Invalid item number.")
            except ValueError:
                print("Invalid input. Please enter a number.")
        elif choice == '6':
            manager.stop_monitoring()
            print("Exiting Clipboard Manager. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 03:22:25
import tkinter as tk
from tkinter import scrolledtext
import pyperclip
import threading
import time
import json
import os

HISTORY_FILE = "clipboard_history.json"
MAX_HISTORY_SIZE = 100
MONITOR_INTERVAL_SECONDS = 0.5

class ClipboardManager:
    def __init__(self, master):
        self.master = master
        master.title("Clipboard History Manager")

        self.history = []
        self.last_clipboard_content = ""
        self.monitoring = True

        self.load_history()

        self.create_widgets()
        self.update_history_display()

        self.start_monitor_thread()

        master.protocol("WM_DELETE_WINDOW", self.on_closing)

    def create_widgets(self):
        # History Listbox
        self.history_frame = tk.Frame(self.master)
        self.history_frame.pack(pady=10, padx=10, fill=tk.BOTH, expand=True)

        self.history_listbox = tk.Listbox(self.history_frame, height=20, width=80, selectmode=tk.SINGLE)
        self.history_listbox.pack(side=tk.LEFT, fill=tk.BOTH, expand=True)

        self.scrollbar = tk.Scrollbar(self.history_frame, orient="vertical", command=self.history_listbox.yview)
        self.scrollbar.pack(side=tk.RIGHT, fill="y")
        self.history_listbox.config(yscrollcommand=self.scrollbar.set)

        # Buttons
        self.button_frame = tk.Frame(self.master)
        self.button_frame.pack(pady=5)

        self.copy_button = tk.Button(self.button_frame, text="Copy Selected", command=self.copy_selected)
        self.copy_button.pack(side=tk.LEFT, padx=5)

        self.clear_button = tk.Button(self.button_frame, text="Clear History", command=self.clear_history)
        self.clear_button.pack(side=tk.LEFT, padx=5)

    def load_history(self):
        if os.path.exists(HISTORY_FILE):
            try:
                with open(HISTORY_FILE, 'r', encoding='utf-8') as f:
                    self.history = json.load(f)
            except json.JSONDecodeError:
                self.history = []
            except Exception as e:
                self.history = []
        else:
            self.history = []

    def save_history(self):
        try:
            with open(HISTORY_FILE, 'w', encoding='utf-8') as f:
                json.dump(self.history, f, indent=4)
        except Exception as e:
            pass # Handle potential errors silently or log them

    def start_monitor_thread(self):
        self.monitor_thread = threading.Thread(target=self.monitor_clipboard, daemon=True)
        self.monitor_thread.start()

    def monitor_clipboard(self):
        while self.monitoring:
            try:
                current_clipboard = pyperclip.paste()
                if current_clipboard and current_clipboard != self.last_clipboard_content:
                    self.add_to_history(current_clipboard)
                    self.last_clipboard_content = current_clipboard
            except pyperclip.PyperclipException:
                pass # Clipboard might be locked or empty
            except Exception as e:
                pass # General error handling for monitoring
            time.sleep(MONITOR_INTERVAL_SECONDS)

    def add_to_history(self, content):
        # Remove if already exists to move to top (most recent)
        if content in self.history:
            self.history.remove(content)
        
        self.history.insert(0, content) # Add to the beginning

        # Trim history if it exceeds max size
        if len(self.history) > MAX_HISTORY_SIZE:
            self.history = self.history[:MAX_HISTORY_SIZE]
        
        # Update GUI from the main thread
        self.master.after(0, self.update_history_display)

    def update_history_display(self):
        self.history_listbox.delete(0, tk.END)
        for item in self.history:
            display_item = item.replace('\n', ' ').replace('\r', '')
            if len(display_item) > 100:
                display_item = display_item[:97] + "..."
            self.history_listbox.insert(tk.END, display_item)

    def copy_selected(self):
        selected_indices = self.history_listbox.curselection()
        if selected_indices:
            index = selected_indices[0]
            selected_content = self.history[index]
            try:
                pyperclip.copy(selected_content)
                self.last_clipboard_content = selected_content # Update last content to avoid re-adding
            except pyperclip.PyperclipException:
                pass # Cannot copy to clipboard

    def clear_history(self):
        self.history = []
        self.update_history_display()
        if os.path.exists(HISTORY_FILE):
            os.remove(HISTORY_FILE)

    def on_closing(self):
        self.monitoring = False # Stop the monitoring thread
        self.save_history()
        self.master.destroy()

if __name__ == "__main__":
    root = tk.Tk()
    app = ClipboardManager(root)
    root.mainloop()

# Additional implementation at 2025-06-21 03:22:49
import tkinter as tk
from tkinter import scrolledtext, messagebox, simpledialog
import pyperclip
import threading
import time
import json
import os
from collections import deque

class ClipboardManager:
    HISTORY_FILE = "clipboard_history.json"
    MAX_HISTORY_SIZE = 100
    CHECK_INTERVAL_SECONDS = 1

    def __init__(self, master):
        self.master = master
        master.title("Clipboard History Manager")
        master.geometry("800x600")

        self.history = deque(maxlen=self.MAX_HISTORY_SIZE)
        self.current_clipboard_content = ""
        self.monitoring_active = False
        self.search_results = [] # To store filtered results when search is active

        self.load_history()

        # --- GUI Setup ---
        # Search Frame
        self.search_frame = tk.Frame(master)
        self.search_frame.pack(pady=5, fill=tk.X)

        self.search_label = tk.Label(self.search_frame, text="Search:")
        self.search_label.pack(side=tk.LEFT, padx=5)

        self.search_entry = tk.Entry(self.search_frame, width=50)
        self.search_entry.pack(side=tk.LEFT, padx=5, expand=True, fill=tk.X)
        self.search_entry.bind("<KeyRelease>", self.search_history) # Live search

        self.clear_search_button = tk.Button(self.search_frame, text="Clear Search", command=self.clear_search)
        self.clear_search_button.pack(side=tk.LEFT, padx=5)

        # History Listbox
        self.history_frame = tk.Frame(master)
        self.history_frame.pack(padx=10, pady=5, fill=tk.BOTH, expand=True)

        self.history_listbox = tk.Listbox(self.history_frame, selectmode=tk.SINGLE)
        self.history_listbox.pack(side=tk.LEFT, fill=tk.BOTH, expand=True)
        self.history_listbox.bind("<Double-1>", self.copy_selected_item) # Double click to copy
        self.history_listbox.bind("<Delete>", self.delete_selected_item) # Delete key to delete

        self.scrollbar = tk.Scrollbar(self.history_frame, orient="vertical", command=self.history_listbox.yview)
        self.scrollbar.pack(side=tk.RIGHT, fill=tk.Y)
        self.history_listbox.config(yscrollcommand=self.scrollbar.set)

        # Buttons Frame
        self.button_frame = tk.Frame(master)
        self.button_frame.pack(pady=5)

        self.copy_button = tk.Button(self.button_frame, text="Copy Selected", command=self.copy_selected_item)
        self.copy_button.pack(side=tk.LEFT, padx=5)

        self.delete_button = tk.Button(self.button_frame, text="Delete Selected", command=self.delete_selected_item)
        self.delete_button.pack(side=tk.LEFT, padx=5)

        self.clear_button = tk.Button(self.button_frame, text="Clear All History", command=self.clear_history)
        self.clear_button.pack(side=tk.LEFT, padx=5)

        # Status Bar
        self.status_bar = tk.Label(master, text="Monitoring clipboard...", bd=1, relief=tk.SUNKEN, anchor=tk.W)
        self.status_bar.pack(side=tk.BOTTOM, fill=tk.X)

        self.update_history_display()
        self.start_monitoring()

        master.protocol("WM_DELETE_WINDOW", self.on_closing)

    def load_history(self):
        if os.path.exists(self.HISTORY_FILE):
            try:
                with open(self.HISTORY_FILE, 'r', encoding='utf-8') as f:
                    loaded_history = json.load(f)
                    # Ensure loaded history respects MAX_HISTORY_SIZE
                    self.history.extend(loaded_history[-self.MAX_HISTORY_SIZE:])
                self.status_bar.config(text=f"History loaded from {self.HISTORY_FILE}")
            except Exception as e:
                messagebox.showerror("Error", f"Failed to load history: {e}")
                self.status_bar.config(text=f"Failed to load history: {e}")
        else:
            self.status_bar.config(text="No history file found. Starting fresh.")

    def save_history(self):
        try:
            with open(self.HISTORY_FILE, 'w', encoding='utf-8') as f:
                json.dump(list(self.history), f, indent=4)
            self.status_bar.config(text=f"History saved to {self.HISTORY_FILE}")
        except Exception as e:
            messagebox.showerror("Error", f"Failed to save history: {e}")
            self.status_bar.config(text=f"Failed to save history: {e}")

    def start_monitoring(self):
        self.monitoring_active = True
        self.monitor_thread = threading.Thread(target=self._monitor_clipboard_loop, daemon=True)
        self.monitor_thread.start()
        self.status_bar.config(text="Monitoring clipboard...")

    def _monitor_clipboard_loop(self):
        while self.monitoring_active:
            try:
                current_content = pyperclip.paste()
                if current_content != self.current_clipboard_content:
                    self.current_clipboard_content = current_content
                    if current_content.strip(): # Only add non-empty content
                        self.add_to_history(current_content)
                        # Update GUI from main thread
                        self.master.after(0, lambda: self.update_history_display(
                            self.search_results if self.search_entry.get().strip() else None
                        ))
                time.sleep(self.CHECK_INTERVAL_SECONDS)
            except pyperclip.PyperclipException as e:
                # This can happen if clipboard is busy or inaccessible
                self.master.after(0, lambda: self.status_bar.config(text=f"Clipboard error: {e}"))
                time.sleep(self.CHECK_INTERVAL_SECONDS * 2) # Wait longer on error
            except Exception as e:
                self.master.after(0, lambda: self.status_bar.config(text=f"An unexpected error occurred: {e}"))
                time.sleep(self.CHECK_INTERVAL_SECONDS * 2)

    def add_to_history(self, item):
        # Remove duplicates if the item is already in history
        if item in self.history:
            self.history.remove(item)
        self.history.appendleft(item) # Add to the beginning

    def update_history_display(self, history_to_display=None):
        self.history_listbox.delete(0, tk.END)
        items_to_show = history_to_display if history_to_display is not None else self.history
        for i, item in enumerate(items_to_show):
            # Truncate long items for display
            display_item = item.replace('\n', ' ').replace('\r', '')
            if len(display_item) > 100:
                display_item = display_item[:97] + "..."
            self.history_listbox.insert(tk.END, f"{i+1}. {display_item}")
        self.status_bar.config(text=f"History updated. {len(items_to_show)} items displayed.")

    def copy_selected_item(self, event=None):
        selected_indices = self.history_listbox.curselection()
        if not selected_indices:
            self.status_bar.config(text="No item selected to copy.")
            return

        selected_index_in_display = selected_indices[0]
        
        # Determine which list we are currently displaying (full history or search results)
        if self.search_entry.get().strip():
            # If search is active, use search_results
            if 0 <= selected_index_in_display < len(self.search_results):
                item_to_copy = self.search_results[selected_index_in_display]
            else:
                self.status_bar.config(text="Error: Selected index out of range for search results.")
                return
        else:
            # Otherwise, use the full history (deque)
            if 0 <= selected_index_in_display < len(self.history):
                item_to_copy = self.history[selected_index_in_display]
            else:
                self.status_bar.config(text="Error: Selected index out of range for full history.")
                return

        try:
            pyperclip.copy(item_to_copy)
            self.status_bar.config(text=f"Copied item to clipboard: {item_to_copy[:50]}...")
            # Move copied item to the top of the history
            self.add_to_history(item_to_copy)
            self.update_history_display(self.search_results if self.search_entry.get().strip() else None)
        except pyperclip.PyperclipException as e:
            messagebox.showerror("Clipboard Error", f"Could not copy to clipboard: {e}")
            self.status_bar.config(text=f"Clipboard error: {e}")

    def delete_selected_item(self, event=None):
        selected_indices = self.history_listbox.curselection()
        if not selected_indices:
            self.status_bar.config(text="No item selected to delete.")
            return

        # Get the actual item from the underlying data structure
        selected_index_in_display = selected_indices[0]
        
        if self.search_entry.get().strip():
            # If search is active, delete from search_results and then from main history
            if 0 <= selected_index_in_display < len(self.search_results):
                item_to_delete = self.search_results.pop(selected_index_in_display)
                if item_to_delete in self.history:
                    # Deque doesn't support direct indexing for deletion, convert to list, remove, convert back
                    temp_list = list(self.history)
                    if item_to_delete in temp_list: # Check again in case of duplicates or complex scenarios
                        temp_list.remove(item_to_delete)
                    self.history = deque(temp_list, maxlen=self.MAX_HISTORY_SIZE)
            else:
                self.status_bar.config(text="Error: Selected index out of range for search results.")
                return
        else:
            # Otherwise, delete directly from the full history (deque)
            if 0 <= selected_index_in_display < len(self.history):
                # Deque doesn't support direct indexing for deletion, convert to list, remove, convert back
                temp_list = list(self.history)
                item_to_delete = temp_list.pop(selected_index_in_display)
                self.history = deque(temp_list, maxlen=self.MAX_HISTORY_SIZE)
            else:
                self.status_bar.config(text="Error: Selected index out of range for full history.")
                return

        self.status_bar.config(text=f"Deleted item: {item_to_delete[:50]}...")
        self.update_history_display(self.search_results if self.search_entry.get().strip() else None) # Refresh display

    def clear_history(self):
        if messagebox.askyesno("Clear History", "Are you sure you want to clear all clipboard history? This cannot be undone."):
            self.history.clear()
            self.search_results.clear() # Also clear search results if any
            self.update_history_display()
            self.status_bar.config(text="All clipboard history cleared.")

    def search_history(self, event=None):
        query = self.search_entry.get().strip().lower()
        if not query:
            self.search_results.clear()
            self.update_history_display() # Show full history
            return

        self.search_results = [item for item in self.history if query in item.lower()]
        self.update_history_display(self.search_results)
        self.status_bar.config(text=f"Found {len(self.search_results)} items matching '{query}'.")

    def clear_search(self):
        self.search_entry.delete(0, tk.END)
        self.search_results.clear()
        self.update_history_display()
        self.status_bar.config(text="Search cleared. Showing full history.")

    def on_closing(self):
        self.monitoring_active = False
        self.save_history()
        self.master.destroy()

if __name__ == "__main__":
    try:
        # Test pyperclip availability
        pyperclip.paste()
    except pyperclip.PyperclipException as e:
        messagebox.showerror("Pyperclip Error",
                             f"Pyperclip is unable to access the clipboard. "
                             f"Please install a copy/paste mechanism for your system (e.g., xclip on Linux).\n\nDetails: {e}")
        exit()

    root = tk.Tk()
    app = ClipboardManager(root)
    root.mainloop()