import threading
import time
from pynput import keyboard

class HotkeyManager:
    def __init__(self):
        self.hotkeys = {}
        self.current_pressed_keys = set()
        self.active_hotkeys = set()

        self.listener = keyboard.Listener(
            on_press=self._on_press,
            on_release=self._on_release
        )

    def _normalize_key(self, key):
        if isinstance(key, keyboard.Key):
            return key
        elif isinstance(key, keyboard.KeyCode):
            return key
        elif isinstance(key, str) and len(key) == 1:
            return keyboard.KeyCode.from_char(key)
        else:
            return None

    def add_hotkey(self, combination, callback):
        normalized_combination_list = []
        for k in combination:
            normalized_k = self._normalize_key(k)
            if normalized_k is None:
                raise ValueError(f"Invalid key format in combination: {k}")
            normalized_combination_list.append(normalized_k)
        
        normalized_combination = frozenset(normalized_combination_list)
        if not normalized_combination:
            raise ValueError("Hotkey combination cannot be empty.")
        self.hotkeys[normalized_combination] = callback

    def _on_press(self, key):
        normalized_key = self._normalize_key(key)
        if normalized_key is None:
            return

        self.current_pressed_keys.add(normalized_key)

        for combo, callback in self.hotkeys.items():
            if combo.issubset(self.current_pressed_keys):
                if combo not in self.active_hotkeys:
                    callback()
                    self.active_hotkeys.add(combo)

    def _on_release(self, key):
        normalized_key = self._normalize_key(key)
        if normalized_key is None:
            return

        if normalized_key in self.current_pressed_keys:
            self.current_pressed_keys.remove(normalized_key)

        for combo in list(self.active_hotkeys):
            if not combo.issubset(self.current_pressed_keys):
                self.active_hotkeys.remove(combo)

    def start(self):
        self.listener.start()

    def stop(self):
        self.listener.stop()
        self.listener.join()

def say_hello():
    print("Hello from hotkey!")

def open_notepad():
    print("Opening notepad (simulated)...")

def exit_program():
    print("Exiting program...")
    global running
    running = False

manager = HotkeyManager()

manager.add_hotkey([keyboard.Key.ctrl_l, keyboard.KeyCode.from_char('h')], say_hello)
manager.add_hotkey([keyboard.Key.ctrl_l, keyboard.Key.alt_l, keyboard.KeyCode.from_char('n')], open_notepad)
manager.add_hotkey([keyboard.Key.ctrl_l, keyboard.Key.alt_l, keyboard.Key.shift_l, keyboard.KeyCode.from_char('q')], exit_program)
manager.add_hotkey([keyboard.Key.esc], exit_program)

manager.start()

running = True
try:
    while running:
        time.sleep(0.1)
except KeyboardInterrupt:
    pass
finally:
    manager.stop()

# Additional implementation at 2025-06-20 22:50:36
import threading
import time
from pynput import keyboard

class HotkeyManager:
    _MODIFIER_MAP = {
        'ctrl': keyboard.Key.ctrl_l,
        'shift': keyboard.Key.shift_l,
        'alt': keyboard.Key.alt_l,
        'win': keyboard.Key.cmd_l,
    }

    def __init__(self):
        self._hotkeys = {}
        self._pressed_keys = set()
        self._listener = None
        self._listener_thread = None
        self._running = False
        self._active_hotkeys = set()

    def _normalize_key_input(self, key_input):
        if isinstance(key_input, str):
            key_input = key_input.lower()
            if key_input in self._MODIFIER_MAP:
                return self._MODIFIER_MAP[key_input]
            elif len(key_input) == 1:
                return keyboard.KeyCode(char=key_input)
            else:
                try:
                    return getattr(keyboard.Key, key_input)
                except AttributeError:
                    raise ValueError(f"Unknown key: {key_input}")
        return key_input

    def _get_canonical_key(self, key):
        if isinstance(key, keyboard.KeyCode):
            return key.char
        elif isinstance(key, keyboard.Key):
            if key in (keyboard.Key.ctrl_l, keyboard.Key.ctrl_r):
                return keyboard.Key.ctrl
            elif key in (keyboard.Key.shift_l, keyboard.Key.shift_r):
                return keyboard.Key.shift
            elif key in (keyboard.Key.alt_l, keyboard.Key.alt_r):
                return keyboard.Key.alt
            elif key in (keyboard.Key.cmd_l, keyboard.Key.cmd_r):
                return keyboard.Key.cmd
            else:
                return key
        return None

    def _on_press(self, key):
        canonical_key = self._get_canonical_key(key)
        if canonical_key:
            self._pressed_keys.add(canonical_key)

        for hotkey_keys_frozen_set, callback in self._hotkeys.items():
            if hotkey_keys_frozen_set == self._pressed_keys:
                if hotkey_keys_frozen_set not in self._active_hotkeys:
                    self._active_hotkeys.add(hotkey_keys_frozen_set)
                    callback()

    def _on_release(self, key):
        canonical_key = self._get_canonical_key(key)
        if canonical_key and canonical_key in self._pressed_keys:
            self._pressed_keys.remove(canonical_key)
        
        hotkeys_to_deactivate = []
        for active_hotkey_set in self._active_hotkeys:
            if not active_hotkey_set.issubset(self._pressed_keys):
                hotkeys_to_deactivate.append(active_hotkey_set)
        
        for hotkey_set in hotkeys_to_deactivate:
            self._active_hotkeys.remove(hotkey_set)

    def register_hotkey(self, keys, callback):
        if not isinstance(keys, (list, tuple)):
            raise TypeError("Keys must be a list or tuple of strings.")
        
        normalized_keys = frozenset(self._normalize_key_input(k) for k in keys)
        
        if any(k is None for k in normalized_keys):
            raise ValueError("One or more keys could not be normalized. Check input.")

        if normalized_keys in self._hotkeys:
            raise ValueError(f"Hotkey {keys} is already registered.")
        
        self._hotkeys[normalized_keys] = callback

    def unregister_hotkey(self, keys):
        if not isinstance(keys, (list, tuple)):
            raise TypeError("Keys must be a list or tuple of strings.")
        
        normalized_keys = frozenset(self._normalize_key_input(k) for k in keys)
        
        if normalized_keys not in self._hotkeys:
            raise ValueError(f"Hotkey {keys} is not registered.")
        
        del self._hotkeys[normalized_keys]

    def start(self):
        if self._running:
            return

        self._listener = keyboard.Listener(
            on_press=self._on_press,
            on_release=self._on_release
        )
        self._listener_thread = threading.Thread(target=self._listener.start)
        self._listener_thread.daemon = True
        self._running = True
        self._listener_thread.start()

    def stop(self):
        if not self._running:
            return

        self._running = False
        if self._listener:
            self._listener.stop()
        if self._listener_thread and self._listener_thread.is_alive():
            self._listener_thread.join(timeout=1)
        self._listener = None
        self._listener_thread = None
        self._pressed_keys.clear()
        self._active_hotkeys.clear()

# Additional implementation at 2025-06-20 22:51:19
import threading
import time
from pynput import keyboard
import subprocess

class HotkeyManager:
    def __init__(self):
        self.hotkeys = {}  # Stores {frozenset_of_keys: action_function}
        self.pressed_keys = set()
        self.active_hotkeys = set() # To prevent re-triggering while keys are held

        self.modifier_keys = {
            keyboard.Key.alt, keyboard.Key.alt_l, keyboard.Key.alt_r,
            keyboard.Key.ctrl, keyboard.Key.ctrl_l, keyboard.Key.ctrl_r,
            keyboard.Key.shift, keyboard.Key.shift_l, keyboard.Key.shift_r,
            keyboard.Key.cmd, keyboard.Key.cmd_l, keyboard.Key.cmd_r,
            keyboard.Key.menu
        }

        self.listener = keyboard.Listener(
            on_press=self._on_press,
            on_release=self._on_release
        )

    def register_hotkey(self, combination, action):
        """
        Registers a hotkey combination with an action.
        :param combination: A single pynput.keyboard.Key/KeyCode or a tuple of them.
                            Example: (keyboard.Key.ctrl_l, keyboard.KeyCode(char='a'))
        :param action: A callable function to execute when the hotkey is pressed.
        """
        if not isinstance(combination, tuple):
            combination = (combination,) # Ensure it's a tuple even for single keys
        
        # Use frozenset for the dictionary key to allow for order-independent lookup
        # and to be hashable.
        self.hotkeys[frozenset(combination)] = action

    def _on_press(self, key):
        # Add the key to the set of currently pressed keys
        self.pressed_keys.add(key)
        
        # Check for hotkey matches
        for combo_set, action in self.hotkeys.items():
            # Check if all keys in the hotkey combination are currently pressed
            if combo_set.issubset(self.pressed_keys):
                # Check if this hotkey is not already active (i.e., triggered and keys still held)
                if combo_set not in self.active_hotkeys:
                    # Check if any non-modifier, non-combo key is pressed
                    extra_non_modifier_keys_pressed = False
                    for pk in self.pressed_keys:
                        if pk not in combo_set and pk not in self.modifier_keys:
                            extra_non_modifier_keys_pressed = True
                            break
                    
                    if not extra_non_modifier_keys_pressed:
                        # This hotkey combination is met and no "extra" non-modifier keys are pressed
                        action()
                        self.active_hotkeys.add(combo_set) # Mark as active

    def _on_release(self, key):
        if key in self.pressed_keys:
            self.pressed_keys.remove(key)
        
        # Check which hotkeys are no longer fully pressed and deactivate them
        hotkeys_to_deactivate = set()
        for active_combo_set in self.active_hotkeys:
            if not active_combo_set.issubset(self.pressed_keys):
                hotkeys_to_deactivate.add(active_combo_set)
        
        for combo_set in hotkeys_to_deactivate:
            self.active_hotkeys.remove(combo_set)

    def start(self):
        """Starts the hotkey listener."""
        print("Hotkey manager started. Press hotkeys or Ctrl+C to exit.")
        self.listener.start()
        self.listener.join() # Keep the main thread alive until listener stops

    def stop(self):
        """Stops the hotkey listener."""
        print("Stopping hotkey manager.")
        self.listener.stop()

def say_hello():
    print("Hello from hotkey!")

def open_notepad():
    print("Attempting to open a text editor...")
    try:
        subprocess.Popen(['notepad.exe']) # Windows
    except FileNotFoundError:
        try:
            subprocess.Popen(['gedit']) # Linux (Ubuntu)
        except FileNotFoundError:
            try:
                subprocess.Popen(['open', '-a', 'TextEdit']) # macOS
            except FileNotFoundError:
                print("Could not find a text editor to open.")

def custom_action_with_args(arg1, arg2):
    print(f"Custom action triggered with args: {arg1}, {arg2}")

def exit_program():
    print("Exiting program via hotkey.")
    global manager
    manager.stop()

if __name__ == "__main__":
    manager = HotkeyManager()

    manager.register_hotkey(
        (keyboard.Key.ctrl_l, keyboard.Key.alt_l, keyboard.KeyCode(char='h')),
        say_hello
    )

    manager.register_hotkey(
        (keyboard.Key.ctrl_l, keyboard.Key.shift_l, keyboard.KeyCode(char='n')),
        open_notepad
    )

    manager.register_hotkey(
        (keyboard.Key.alt_l, keyboard.Key.f1),
        lambda: custom_action_with_args("value1", "value2")
    )

    manager.register_hotkey(
        (keyboard.Key.ctrl_l, keyboard.KeyCode(char='q')),
        exit_program
    )
    
    manager.register_hotkey(
        keyboard.Key.f12,
        lambda: print("F12 pressed

# Additional implementation at 2025-06-20 22:52:49
