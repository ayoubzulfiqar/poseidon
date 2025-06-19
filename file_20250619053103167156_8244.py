import threading
import time
import sys
from pynput import keyboard

class HotkeyManager:
    def __init__(self):
        self.hotkeys = []
        self.listener = None
        self.running = False
        self._hotkey_callbacks = {}
        self._stop_event = threading.Event()

    def _process_keys_for_hotkey(self, keys_input):
        processed_keys = []
        for key_item in keys_input:
            if isinstance(key_item, str):
                key_item_lower = key_item.lower()
                if key_item_lower == 'ctrl':
                    processed_keys.append(keyboard.Key.ctrl_l)
                elif key_item_lower == 'alt':
                    processed_keys.append(keyboard.Key.alt_l)
                elif key_item_lower == 'shift':
                    processed_keys.append(keyboard.Key.shift_l)
                elif key_item_lower == 'cmd' or key_item_lower == 'super' or key_item_lower == 'win':
                    processed_keys.append(keyboard.Key.cmd_l)
                elif hasattr(keyboard.Key, key_item_lower):
                    processed_keys.append(getattr(keyboard.Key, key_item_lower))
                elif len(key_item) == 1:
                    processed_keys.append(keyboard.KeyCode(char=key_item))
                else:
                    try:
                        processed_keys.append(getattr(keyboard.Key, key_item_lower))
                    except AttributeError:
                        pass
            elif isinstance(key_item, (keyboard.Key, keyboard.KeyCode)):
                processed_keys.append(key_item)
        return processed_keys

    def _on_activate_hotkey(self, hotkey_canonical_keys):
        callback = self._hotkey_callbacks.get(hotkey_canonical_keys)
        if callback:
            threading.Thread(target=callback, daemon=True).start()

    def add_hotkey(self, keys_input, callback_function):
        processed_keys = self._process_keys_for_hotkey(keys_input)
        if not processed_keys:
            return

        canonical_keys_for_dict = frozenset(keyboard.Listener.canonical(k) for k in processed_keys)
        
        hotkey_instance = keyboard.HotKey(processed_keys, lambda: self._on_activate_hotkey(canonical_keys_for_dict))
        
        self.hotkeys.append(hotkey_instance)
        self._hotkey_callbacks[canonical_keys_for_dict] = callback_function

    def start(self):
        if self.running:
            return

        def on_press(key):
            for hotkey_instance in self.hotkeys:
                try:
                    hotkey_instance.press(self.listener.canonical(key))
                except AttributeError:
                    pass

        def on_release(key):
            for hotkey_instance in self.hotkeys:
                try:
                    hotkey_instance.release(self.listener.canonical(key))
                except AttributeError:
                    pass

        self.listener = keyboard.Listener(on_press=on_press, on_release=on_release)
        self.listener.start()
        self.

# Additional implementation at 2025-06-19 05:32:20
import keyboard
import subprocess
import sys
import os

class HotkeyManager:
    def __init__(self):
        self.hotkeys = {}
        self.enabled = True

    def _hotkey_callback_wrapper(self, hotkey_id):
        if not self.enabled:
            return

        hotkey_info = self.hotkeys.get(hotkey_id)
        if hotkey_info:
            

# Additional implementation at 2025-06-19 05:33:50
import json
import os
import subprocess
import threading
import time
from pynput import keyboard
import pyautogui

class HotkeyManager:
    def __init__(self, config_file="hotkeys.json"):
        self.config_file = config_file
        self.hotkeys = {}  # Stores {'key_combination_str': {'action_type': 'command', 'value': 'notepad.exe'}}
        self.listener = None
        self._current_keys = set() # To track currently pressed keys for custom combinations
        self._load_hotkeys()

    def _load_hotkeys(self):
        if os.path.exists(self.config_file):
            try:
                with open(self.config_file, 'r') as f:
                    loaded_data = json.load(f)
                    # Convert list of keys back to tuple for dictionary key if needed,
                    # but we store as string keys, so direct assignment is fine.
                    self.hotkeys = loaded_data
                print(f"Loaded {len(self.hotkeys)} hotkeys from {self.config_file}")
            except json.JSONDecodeError:
                print(f"Error decoding JSON from {self.config_file}. Starting with empty hotkeys.")
                self.hotkeys = {}
            except Exception as e:
                print(f"An error occurred loading hotkeys: {e}")
                self.hotkeys = {}
        else:
            print(f"No config file found at {self.config_file}. Starting with empty hotkeys.")

    def _save_hotkeys(self):
        try:
            # Keys are already strings, so direct dump is fine.
            with open(self.config_file, 'w') as f:
                json.dump(self.hotkeys, f, indent=4)
            print(f"Saved {len(self.hotkeys)} hotkeys to {self.config_file}")
        except Exception as e:
            print(f"An error occurred saving hotkeys: {e}")

    def _keys_to_str(self, keys):
        # Sort keys to ensure consistent representation (e.g., ('ctrl', 'a') is same as ('a', 'ctrl'))
        # Convert Key objects to their string representation if they are pynput Key objects
        processed_keys = []
        for k in keys:
            if isinstance(k, keyboard.Key):
                processed_keys.append(str(k).replace("Key.", ""))
            elif isinstance(k, str):
                processed_keys.append(k)
            else:
                processed_keys.append(str(k)) # Fallback for other types
        
        sorted_keys = sorted(processed_keys)
        return '+'.join(sorted_keys)

    def add_hotkey(self, keys, action_type, action_value):
        """
        Adds a new hotkey.
        :param keys: A tuple of strings representing the key combination (e.g., ('ctrl', 'alt', 'a')).
                     For special keys, use pynput's Key enum names (e.g., 'ctrl', 'alt', 'shift', 'f1', 'space').
                     For regular keys, use their character string (e.g., 'a', 'b', '1').
        :param action_type: Type of action ('command', 'open_file', 'type_text').
        :param action_value: The value associated with the action (e.g., 'notepad.exe', 'https://google.com', 'Hello World!').
        """
        if not isinstance(keys, tuple):
            raise ValueError("Keys must be a tuple of strings.")
        
        hotkey_str = self._keys_to_str(keys)
        self.hotkeys[hotkey_str] = {'action_type': action_type, 'value': action_value}
        self._save_hotkeys()
        print(f"Hotkey added: {hotkey_str} -> {action_type}: {action_value}")
        # No need to rebuild listener immediately, _on_press checks current_keys against self.hotkeys

    def remove_hotkey(self, keys):
        """
        Removes an existing hotkey.
        :param keys: A tuple of strings representing the key combination to remove.
        """
        if not isinstance(keys, tuple):
            raise ValueError("Keys must be a tuple of strings.")
        
        hotkey_str = self._keys_to_str(keys)
        if hotkey_str in self.hotkeys:
            del self.hotkeys[hotkey_str]
            self._save_hotkeys()
            print(f"Hotkey removed: {hotkey_str}")
        else:
            print(f"Hotkey not found: {hotkey_str}")

    def list_hotkeys(self):
        """Lists all registered hotkeys."""
        if not self.hotkeys:
            print("No hotkeys registered.")
            return
        print("\n--- Registered Hotkeys ---")
        for key_str, action in self.hotkeys.items():
            print(f"  {key_str}: Type='{action['action_type']}', Value='{action['value']}'")
        print("--------------------------\n")

    def _execute_action(self, action_type, action_value):
        print(f"Executing action: Type='{action_type}', Value='{action_value}'")
        try:
            if action_type == 'command':
                # Use shell=True for simple commands, or pass list for more control
                subprocess.Popen(action_value, shell=True) 
            elif action_type == 'open_file':
                # os.startfile is Windows-specific. For cross-platform, use subprocess.run with 'xdg-open' or 'open'.
                if os.name == 'nt': # Windows
                    os.startfile(action_value)
                elif os.name == 'posix': # Linux/macOS
                    subprocess.Popen(['xdg-open', action_value]) # Linux
                    # For macOS, use ['open', action_value]
            elif action_type == 'type_text':
                pyautogui.write(action_value)
            else:
                print(f"Unknown action type: {action_type}")
        except Exception as e:
            print(f"Error executing action '{action_value}': {e}")

    def _on_press(self, key):
        try:
            # Convert pynput Key object to its string representation (e.g., Key.ctrl_l -> 'ctrl_l')
            # For character keys, it's already a string (e.g., 'a')
            key_str_representation = str(key).replace("'", "")
            if "Key." in key_str_representation:
                key_str_representation = key_str_representation.replace("Key.", "")
            
            self._current_keys.add(key_str_representation)
            
            # Check if the current combination of pressed keys matches any hotkey
            current_combination_str = self._keys_to_str(tuple(self._current_keys))
            
            if current_combination_str in self.hotkeys:
                action = self.hotkeys[current_combination_str]
                # Execute action in a separate thread to avoid blocking the listener
                threading.Thread(target=self._execute_action, args=(action['action_type'], action['value'])).start()
        except AttributeError:
            # This can happen for special keys that don't have a simple char representation
            pass 

    def _on_release(self, key):
        try:
            key_str_representation = str(key).replace("'", "")
            if "Key." in key_str_representation:
                key_str_representation = key_str_representation.replace("Key.", "")
            
            if key_str_representation in self._current_keys:
                self._current_keys.remove(key_str_representation)
        except KeyError:
            pass # Key was not tracked, or already removed

    def start(self):
        """Starts the hotkey listener."""
        if self.listener is not None and self.listener.running:
            print("Listener is already running.")
            return

        print("Starting hotkey listener...")
        self.listener = keyboard.Listener(
            on_press=self._on_press,
            on_release=self._on_release
        )
        self.listener.start()
        print("Hotkey listener started. Press Ctrl+C to stop.")

    def stop(self):
        """Stops the hotkey listener."""
        if self.listener and self.listener.running:
            print("Stopping hotkey listener...")
            self.listener.stop()
            self.listener = None
            print("Hotkey listener stopped.")
        else:
            print("Listener is not running.")

if __name__ == "__main__":
    manager = HotkeyManager()

    # Add some example hotkeys
    # Key names for pynput:
    # - Regular characters: 'a', 'b', '1', '2', etc.
    # - Special keys: 'ctrl', 'alt', 'shift', 'space', 'enter', 'esc', 'tab', 'backspace', 'delete',
    #   'up', 'down', 'left', 'right', 'home', 'end', 'page_up', 'page_down', 'insert',
    #   'f1', 'f2', ..., 'f12'.
    # - Specific left/right modifiers: 'ctrl_l', 'ctrl_r', 'alt_l', 'alt_r', 'shift_l', 'shift_r'.

    # Example 1: Ctrl+Alt+N to open Notepad (Windows) or a terminal (Linux/macOS)
    if os.name == 'nt': # Windows
        manager.add_hotkey(('ctrl', 'alt', 'n'), 'command', 'notepad.exe')
    elif os.name == 'posix': # Linux/macOS
        # Adjust command for your OS: 'gnome-terminal', 'konsole', 'xterm' for Linux, 'open -a Terminal' for macOS
        manager.add_hotkey(('ctrl', 'alt', 'n'), 'command', 'xterm') 

    # Example 2: Ctrl+Shift+G to open Google in default browser
    manager.add_hotkey(('ctrl', 'shift', 'g'), 'open_file', 'https://www.google.com')

    # Example 3: Alt+T to type "Hello from Python Hotkey Manager!"
    manager.add_hotkey(('alt', 't'), 'type_text', 'Hello from Python Hotkey Manager!')

    # Example 4: Ctrl+Alt+D to remove the Alt+T hotkey (demonstrates dynamic removal)
    # This hotkey is for demonstration purposes.
    # You might want to remove this line or change its action in a production system.
    # manager.add_hotkey(('ctrl', 'alt', 'd'), 'command', 'python -c "import time; time.sleep(0.1); manager.remove_hotkey((\'alt\', \'t\'))"')
    # Note: Directly calling manager methods from a subprocess is complex.
    # A better way would be to have a separate management interface (e.g., a simple CLI or GUI).
    # For this example, we'll just show the add/remove methods.

    manager.list_hotkeys()

    manager.start()

    try:
        # Keep the main thread alive to allow the listener to run in the background.
        # In a real application, this might be a GUI event loop or a more sophisticated main loop.
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nCtrl+C detected. Shutting down.")
    finally:
        manager.stop()
        print("Application exited.")