import time
import sys
from pynput import keyboard, mouse

recorded_events = []
recording = False
playback_active = False
keyboard_controller = keyboard.Controller()
mouse_controller = mouse.Controller()

def on_press(key):
    global recording
    if recording:
        try:
            recorded_events.append({'time': time.time(), 'type': 'key_press', 'key': key.char})
        except AttributeError:
            recorded_events.append({'time': time.time(), 'type': 'key_press', 'key': key})
    if key == keyboard.Key.esc:
        if recording:
            recording = False
            print("Recording stopped.")
            return False

def on_release(key):
    global recording
    if recording:
        try:
            recorded_events.append({'time': time.time(), 'type': 'key_release', 'key': key.char})
        except AttributeError:
            recorded_events.append({'time': time.time(), 'type': 'key_release', 'key': key})

def on_click(x, y, button, pressed):
    global recording
    if recording:
        recorded_events.append({'time': time.time(), 'type': 'mouse_click', 'x': x, 'y': y, 'button': button, 'pressed': pressed})

def on_scroll(x, y, dx, dy):
    global recording
    if recording:
        recorded_events.append({'time': time.time(), 'type': 'mouse_scroll', 'x': x, 'y': y, 'dx': dx, 'dy': dy})

def on_move(x, y):
    global recording
    if recording:
        recorded_events.append({'time': time.time(), 'type': 'mouse_move', 'x': x, 'y': y})

def start_recording_session():
    global recorded_events, recording
    if recording:
        print("Already recording.")
        return
    recorded_events = []
    recording = True
    print("Recording started. Press 'Esc' to stop.")

    keyboard_listener = keyboard.Listener(on_press=on_press, on_release=on_release)
    mouse_listener = mouse.Listener(on_click=on_click, on_scroll=on_scroll, on_move=on_move)

    keyboard_listener.start()
    mouse_listener.start()

    keyboard_listener.join()
    mouse_listener.join()
    print("Recording session ended.")

def play_recorded_events():
    global playback_active
    if not recorded_events:
        print("No events recorded to play back.")
        return
    if playback_active:
        print("Playback already in progress.")
        return

    playback_active = True
    print("Playing back recorded events...")

    start_time = recorded_events[0]['time'] if recorded_events else time.time()
    for event in recorded_events:
        delay = event['time'] - start_time
        if delay > 0:
            time.sleep(delay)
        start_time = event['time']

        if event['type'] == 'key_press':
            try:
                keyboard_controller.press(event['key'])
            except ValueError:
                if isinstance(event['key'], keyboard.Key):
                    keyboard_controller.press(event['key'])
                else:
                    keyboard_controller.press(keyboard.KeyCode.from_char(event['key']))
        elif event['type'] == 'key_release':
            try:
                keyboard_controller.release(event['key'])
            except ValueError:
                if isinstance(event['key'], keyboard.Key):
                    keyboard_controller.release(event['key'])
                else:
                    keyboard_controller.release(keyboard.KeyCode.from_char(event['key']))
        elif event['type'] == 'mouse_click':
            mouse_controller.position = (event['x'], event['y'])
            if event['pressed']:
                mouse_controller.press(event['button'])
            else:
                mouse_controller.release(event['button'])
        elif event['type'] == 'mouse_scroll':
            mouse_controller.position = (event['x'], event['y'])
            mouse_controller.scroll(event['dx'], event['dy'])
        elif event['type'] == 'mouse_move':
            mouse_controller.position = (event['x'], event['y'])

    print("Playback finished.")
    playback_active = False

def main_menu():
    while True:
        print("\n--- Automation Recorder ---")
        print("1. Start Recording (Press 'Esc' to stop)")
        print("2. Play Recording")
        print("3. Exit")
        choice = input("Enter your choice: ")

        if choice == '1':
            start_recording_session()
        elif choice == '2':
            play_recorded_events()
        elif choice == '3':
            print("Exiting.")
            sys.exit(0)
        else:
            print("Invalid choice. Please try again.")

main_menu()

# Additional implementation at 2025-06-20 23:07:15
import time
import json
from pynput import keyboard, mouse

class AutomationRecorder:
    def __init__(self):
        self.events = []
        self.start_time = None
        self.keyboard_listener = None
        self.mouse_listener = None
        self.is_recording = False

    def _on_press(self, key):
        if not self.is_recording:
            return
        try:
            event_key = key.char
        except AttributeError:
            event_key = str(key)
        self._add_

# Additional implementation at 2025-06-20 23:08:16
import time
import json
from pynput import keyboard, mouse

class AutomationRecorder:
    def __init__(self):
        self.events = []
        self.is_recording = False
        self.start_time = None
        self.keyboard_listener = None
        self.mouse_listener = None
        self.keyboard_controller = keyboard.Controller()
        self.mouse_controller = mouse.Controller()

    def _on_press(self, key):
        if not self.is_recording:
            return
        try:
            char_key = key.char
        except AttributeError:
            char_key = str(key) # Special keys like Key.space, Key.ctrl_l
        self.events.append({
            "type": "keyboard",
            "action": "press",
            "key": char_key,
            "time": time.time() - self.start_time
        })

    def _on_release(self, key):
        if not self.is_recording:
            return
        try:
            char_key = key.char
        except AttributeError:
            char_key = str(key)
        self.events.append({
            "type": "keyboard",
            "action": "release",
            "key": char_key,
            "time": time.time() - self.start_time
        })
        if key == keyboard.Key.esc: # Stop recording with ESC key
            self.stop_recording()

    def _on_mouse_move(self, x, y):
        if not self.is_recording:
            return
        self.events.append({
            "type": "mouse",
            "action": "move",
            "x": x,
            "y": y,
            "time": time.time() - self.start_time
        })

    def _on_mouse_click(self, x, y, button, pressed):
        if not self.is_recording:
            return
        self.events.append({
            "type": "mouse",
            "action": "click",
            "x": x,
            "y": y,
            "button": str(button),
            "pressed": pressed,
            "time": time.time() - self.start_time
        })

    def _on_mouse_scroll(self, x, y, dx, dy):
        if not self.is_recording:
            return
        self.events.append({
            "type": "mouse",
            "action": "scroll",
            "x": x,
            "y": y,
            "dx": dx,
            "dy": dy,
            "time": time.time() - self.start_time
        })

    def start_recording(self):
        if self.is_recording:
            print("Already recording.")
            return
        print("Recording started. Press ESC to stop.")
        self.events = []
        self.is_recording = True
        self.start_time = time.time()

        self.keyboard_listener = keyboard.Listener(
            on_press=self._on_press,
            on_release=self._on_release
        )
        self.mouse_listener = mouse.Listener(
            on_move=self._on_mouse_move,
            on_click=self._on_mouse_click,
            on_scroll=self._on_mouse_scroll
        )

        self.keyboard_listener.start()
        self.mouse_listener.start()
        self.keyboard_listener.join() # Blocks until ESC is pressed and listener stops

    def stop_recording(self):
        if not self.is_recording:
            print("Not currently recording.")
            return
        print("Recording stopped.")
        self.is_recording = False
        if self.keyboard_listener:
            self.keyboard_listener.stop()
            self.keyboard_listener = None
        if self.mouse_listener:
            self.mouse_listener.stop()
            self.mouse_listener = None

    def play_recording(self, speed_multiplier=1.0, loop=1):
        if not self.events:
            print("No events to play. Record something first.")
            return

        print(f"Playing recording (x{speed_multiplier} speed, {loop} loop(s))...")
        for _ in range(loop):
            last_event_time = 0.0
            for event in self.events:
                delay = (event["time"] - last_event_time) / speed_multiplier
                if delay > 0:
                    time.sleep(delay)
                last_event_time = event["time"]

                if event["type"] == "keyboard":
                    key_str = event["key"]
                    try:
                        # Try to convert string back to pynput Key object if it's a special key
                        if key_str.startswith("Key."):
                            key = getattr(keyboard.Key, key_str.split('.')[1])
                        else: # It's a character
                            key = key_str
                    except AttributeError:
                        key = key_str # Fallback if conversion fails

                    if event["action"] == "press":
                        self.keyboard_controller.press(key)
                    elif event["action"] == "release":
                        self.keyboard_controller.release(key)
                elif event["type"] == "mouse":
                    if event["action"] == "move":
                        self.mouse_controller.position = (event["x"], event["y"])
                    elif event["action"] == "click":
                        button_str = event["button"]
                        try:
                            # Try to convert string back to pynput Button object
                            button = getattr(mouse.Button, button_str.split('.')[1])
                        except AttributeError:
                            button = mouse.Button.left # Default to left if conversion fails

                        if event["pressed"]:
                            self.mouse_controller.press(button)
                        else:
                            self.mouse_controller.release(button)
                    elif event["action"] == "scroll":
                        self.mouse_controller.scroll(event["dx"], event["dy"])
            print(f"Loop {_ + 1}/{loop} finished.")
        print("Playback finished.")

    def save_recording(self, filename):
        if not self.events:
            print("No events to save.")
            return
        try:
            with open(filename, 'w') as f:
                json.dump(self.events, f, indent=4)
            print(f"Recording saved to {filename}")
        except IOError as e:
            print(f"Error saving recording: {e}")

    def load_recording(self, filename):
        try:
            with open(filename, 'r') as f:
                self.events = json.load(f)
            print(f"Recording loaded from {filename}")
        except FileNotFoundError:
            print(f"File not found: {filename}")
        except json.JSONDecodeError:
            print(f"Error decoding JSON from {filename}. File might be corrupted.")
        except IOError as e:
            print(f"Error loading recording: {e}")

if __name__ == "__main__":
    recorder = AutomationRecorder()
    recording_file = "my_automation.json"

    while True:
        print("\nAutomation Recorder Menu:")
        print("1. Start Recording (Press ESC to stop)")
        print("2. Play Recording")
        print("3. Save Recording")
        print("4. Load Recording")
        print("5. Exit")

        choice = input("Enter your choice: ")

        if choice == '1':
            recorder.start_recording()
        elif choice == '2':
            if not recorder.events:
                print("No recording loaded or made. Please record or load first.")
                continue
            try:
                speed_str = input("Enter playback speed multiplier (e.g., 0.5 for half, 2 for double, 1 for normal): ")
                speed = float(speed_str)
                if speed <= 0:
                    print("Speed multiplier must be positive.")
                    continue
            except ValueError:
                print("Invalid speed multiplier. Using default (1.0).")
                speed = 1.0

            try:
                loop_str = input("Enter number of loops (default 1): ")
                loops = int(loop_str)
                if loops <= 0:
                    print("Number of loops must be positive.")
                    continue
            except ValueError:
                print("Invalid loop count. Using default (1).")
                loops = 1

            recorder.play_recording(speed_multiplier=speed, loop=loops)
        elif choice == '3':
            recorder.save_recording(recording_file)
        elif choice == '4':
            recorder.load_recording(recording_file)
        elif choice == '5':
            print("Exiting.")
            break
        else:
            print("Invalid choice. Please try again.")

# Additional implementation at 2025-06-20 23:09:11
import time
import json
import threading
from pynput import keyboard, mouse

class AutomationRecorder:
    def __init__(self):
        self.events = []
        self.recording = False
        self.start_time = None
        self.keyboard_listener = None
        self.mouse_listener = None
        self.keyboard_controller = keyboard.Controller()
        self.mouse_controller = mouse.Controller()

    def _on_press(self, key):
        if not self.recording:
            return
        try:
            key_char = key.char
        except AttributeError:
            key_char = str(key)
        self.events.append({
            'type': 'keyboard_press',
            'key': key_char,
            'time_offset': time.time() - self.start_time
        })

    def _on_release(self, key):
        if not self.recording:
            return
        try:
            key_char = key.char
        except AttributeError:
            key_char = str(key)
        self.events.append({
            'type': 'keyboard_release',
            'key': key_char,
            'time_offset': time.time() - self.start_time
        })
        if key == keyboard.Key.esc:
            self.stop_recording()
            return False

    def _on_mouse_click(self, x, y, button, pressed):
        if not self.recording:
            return
        self.events.append({
            'type': 'mouse_click',
            'x': x,
            'y': y,
            'button': str(button),
            'pressed': pressed,
            'time_offset': time.time() - self.start_time
        })

    def _on_mouse_scroll(self, x, y, dx, dy):
        if not self.recording:
            return
        self.events.append({
            'type': 'mouse_scroll',
            'x': x,
            'y': y,
            'dx': dx,
            'dy': dy,
            'time_offset': time.time() - self.start_time
        })

    def _on_mouse_move(self, x, y):
        if not self.recording:
            return
        self.events.append({
            'type': 'mouse_move',
            'x': x,
            'y': y,
            'time_offset': time.time() - self.start_time
        })

    def start_recording(self):
        if self.recording:
            print("Already recording.")
            return

        print("Starting recording. Press 'Esc' to stop recording.")
        self.events = []
        self.recording = True
        self.start_time = time.time()

        self.keyboard_listener = keyboard.Listener(on_press=self._on_press, on_release=self._on_release)
        self.mouse_listener = mouse.Listener(on_click=self._on_mouse_click, on_scroll=self._on_mouse_scroll, on_move=self._on_mouse_move)

        self.keyboard_listener.start()
        self.mouse_listener.start()

        self.keyboard_listener.join()
        self.mouse_listener.join()
        print("Recording stopped.")

    def stop_recording(self):
        if not self.recording:
            print("Not currently recording.")
            return

        self.recording = False
        if self.keyboard_listener:
            self.keyboard_listener.stop()
            self.keyboard_listener = None
        if self.mouse_listener:
            self.mouse_listener.stop()
            self.mouse_listener = None
        print("Recording session ended.")

    def play_recording(self, speed_multiplier=1.0):
        if not self.events:
            print("No events to play. Record something first or load a file.")
            return

        print(f"Playing back recording with speed multiplier: {speed_multiplier}x")
        last_event_time_offset = 0.0

        for event in self.events:
            current_time_offset = event['time_offset']
            sleep_duration = (current_time_offset - last_event_time_offset) / speed_multiplier
            if sleep_duration > 0:
                time.sleep(sleep_duration)
            last_event_time_offset = current_time_offset

            event_type = event['type']
            if event_type == 'keyboard_press':
                try:
                    key = keyboard.Key[event['key'].split('.')[-1]]
                except KeyError:
                    key = event['key']
                self.keyboard_controller.press(key)
            elif event_type == 'keyboard_release':
                try:
                    key = keyboard.Key[event['key'].split('.')[-1]]
                except KeyError:
                    key = event['key']
                self.keyboard_controller.release(key)
            elif event_type == 'mouse_click':
                self.mouse_controller.position = (event['x'], event['y'])
                button = mouse.Button[event['button'].split('.')[-1]]
                if event['pressed']:
                    self.mouse_controller.press(button)
                else:
                    self.mouse_controller.release(button)
            elif event_type == 'mouse_scroll':
                self.mouse_controller.position = (event['x'], event['y'])
                self.mouse_controller.scroll(event['dx'], event['dy'])
            elif event_type == 'mouse_move':
                self.mouse_controller.position = (event['x'], event['y'])
            else:
                print(f"Unknown event type: {event_type}")

        print("Playback finished.")

    def save_recording(self, filename="recording.json"):
        if not self.events:
            print("No events to save.")
            return
        try:
            with open(filename, 'w') as f:
                json.dump(self.events, f, indent=4)
            print(f"Recording saved to {filename}")
        except IOError as e:
            print(f"Error saving recording: {e}")

    def load_recording(self, filename="recording.json"):
        try:
            with open(filename, 'r') as f:
                self.events = json.load(f)
            print(f"Recording loaded from {filename}. Total events: {len(self.events)}")
        except FileNotFoundError:
            print(f"File not found: {filename}")
            self.events = []
        except json.JSONDecodeError as e:
            print(f"Error decoding JSON from {filename}: {e}")
            self.events = []
        except IOError as e:
            print(f"Error loading recording: {e}")
            self.events = []

def main():
    recorder = AutomationRecorder()
    filename = "my_automation.json"

    while True:
        print("\n--- Automation Recorder Menu ---")
        print("1. Start Recording (Press 'Esc' to stop)")
        print("2. Play Recording")
        print("3. Save Recording")
        print("4. Load Recording")
        print("5. Exit")
        choice = input("Enter your choice: ")

        if choice == '1':
            recorder.start_recording()
        elif choice == '2':
            if not recorder.events:
                print("No recording loaded or made. Please record or load first.")
                continue
            try:
                speed_str = input("Enter playback speed multiplier (e.g., 1.0 for normal, 2.0 for double, 0.5 for half): ")
                speed = float(speed_str)
                if speed <= 0:
                    print("Speed multiplier must be positive.")
                    continue
                recorder.play_recording(speed_multiplier=speed)
            except ValueError:
                print("Invalid speed multiplier. Please enter a number.")
        elif choice == '3':
            recorder.save_recording(filename)
        elif choice == '4':
            recorder.load_recording(filename)
        elif choice == '5':
            print("Exiting recorder.")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()