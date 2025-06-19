import keyboard
import sys

def hotkey_action_a():
    print("Action A triggered!")

def hotkey_action_b():
    print("Action B triggered!")

def hotkey_quit():
    print("Quitting hotkey manager.")
    keyboard.unhook_all()
    sys.exit(0)

def main():
    keyboard.add_hotkey("ctrl+alt+a", hotkey_action_a)
    keyboard.add_hotkey("ctrl+alt+b", hotkey_action_b)
    keyboard.add_hotkey("ctrl+alt+q", hotkey_quit)

    print("Hotkey manager active. Press Ctrl+Alt+A, Ctrl+Alt+B, or Ctrl+Alt+Q to exit.")
    try:
        keyboard.wait()
    except KeyboardInterrupt:
        print("Manager interrupted.")
    finally:
        keyboard.unhook_all()

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-19 00:30:10


# Additional implementation at 2025-06-19 00:31:19
import keyboard
import threading
import time
import sys
import subprocess

class HotkeyManager:
    def __init__(self):
        self._hotkeys = {}
        self._listener_thread = None
        self._running = False

    def register_hotkey(self, hotkey_combination: str, callback_function):
        if not callable(callback_function):
            raise TypeError("Callback function must be callable.")
        self._hotkeys[hotkey_combination] = callback_function
        keyboard.add_hotkey(hotkey_combination, callback_function)
        print(f"Registered hotkey: {hotkey_combination}")

    def unregister_hotkey(self, hotkey_combination: str):
        if hotkey_combination in self._hotkeys:
            keyboard.remove_hotkey(hotkey_combination)
            del self._hotkeys[hotkey_combination]
            print(f"Unregistered hotkey: {hotkey_combination}")
        else:
            print(f"Hotkey '{hotkey_combination}' not found.")

    def list_hotkeys(self):
        if not self._hotkeys:
            print("No hotkeys registered.")
            return
        print("Registered Hotkeys:")
        for hotkey, callback in self._hotkeys.items():
            print(f"- {hotkey} -> {callback.__name__}")

    def start_listening(self):
        if self._running:
            print("Hotkey listener is already running.")
            return
        print("Starting hotkey listener...")
        def _listener_loop():
            print("Hotkey listener thread started.")
            while self._running:
                time.sleep(0.1)
            print("Hotkey listener thread stopped.")
        if self._listener_thread is None or not self._listener_thread.is_alive():
            self._running = True
            self._listener_thread = threading.Thread(target=_listener_loop, daemon=True)
            self._listener_thread.start()
            print("Hotkey listener started in background.")
        else:
            print("Hotkey listener is already running.")

    def stop_listening(self):
        if not self._running:
            print("Hotkey listener is not running.")
            return
        print("Stopping hotkey listener...")
        self._running = False
        if self._listener_thread and self._listener_thread.is_alive():
            self._listener_thread.join(timeout=1)
        hotkeys_to_unregister = list(self._hotkeys.keys())
        for hotkey in hotkeys_to_unregister:
            self.unregister_hotkey(hotkey)
        keyboard.unhook_all()
        print("Hotkey listener stopped and all hotkeys unregistered.")

def say_hello():
    print("Hello from hotkey!")

def open_notepad():
    try:
        if sys.platform.startswith('win'):
            subprocess.Popen(['notepad.exe'])
            print("Opened Notepad.")
        elif sys.platform.startswith('darwin'):
            subprocess.Popen(['open', '-a', 'TextEdit'])
            print("Opened TextEdit.")
        else:
            subprocess.Popen(['gedit'])
            print("Opened Gedit.")
    except FileNotFoundError:
        print("Error: Application not found.")
    except Exception as e:
        print(f"Error opening application: {e}")

def custom_action_with_args(arg1, arg2):
    print(f"Custom action triggered with args: {arg1}, {arg2}")

def exit_program():
    print("Exiting program via hotkey...")
    global manager_instance
    if manager_instance:
        manager_instance.stop_listening()
    sys.exit(0)

if __name__ == "__main__":
    manager_instance = HotkeyManager()
    manager_instance.register_hotkey("ctrl+alt+h", say_hello)
    manager_instance.register_hotkey("ctrl+alt+n", open_notepad)
    manager_instance.register_hotkey("ctrl+alt+x", exit_program)
    manager_instance.register_hotkey("ctrl+alt+c", lambda: custom_action_with_args("value1", 123))
    manager_instance.list_hotkeys()
    manager_instance.start_listening()
    print("\nHotkey manager is active. Press 'ctrl+alt+h' for hello, 'ctrl+alt+n' for notepad, 'ctrl+alt+c' for custom action, or 'ctrl+alt+x' to exit.")
    print("You can also press 'q' and Enter in this console to quit.")
    try:
        while True:
            user_input = input("")
            if user_input.lower() == 'q':
                print("Quitting by user request.")
                break
            time.sleep(0.1)
    except KeyboardInterrupt:
        print("\nProgram interrupted by user (Ctrl+C).")
    finally:
        if manager_instance:
            manager_instance.stop_listening()
        print("Program terminated.")