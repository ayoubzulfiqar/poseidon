import time
import json
import os

STATS_FILE = "productivity_stats.json"
DEFAULT_FOCUS_DURATION_MIN = 25
DEFAULT_BREAK_DURATION_MIN = 5

def load_stats():
    if os.path.exists(STATS_FILE):
        with open(STATS_FILE, 'r') as f:
            try:
                return json.load(f)
            except json.JSONDecodeError:
                print("Warning: Corrupted stats file. Starting fresh.")
                return initialize_stats()
    return initialize_stats()

def initialize_stats():
    return {
        "total_focus_sessions": 0,
        "total_focus_time_seconds": 0,
        "last_updated": None
    }

def save_stats(stats):
    stats["last_updated"] = time.strftime("%Y-%m-%d %H:%M:%S")
    with open(STATS_FILE, 'w') as f:
        json.dump(stats, f, indent=4)

def display_stats(stats):
    print("\n--- Productivity Stats ---")
    print(f"Total Focus Sessions: {stats['total_focus_sessions']}")
    total_minutes = stats['total_focus_time_seconds'] // 60
    total_hours = total_minutes // 60
    remaining_minutes = total_minutes % 60
    print(f"Total Focus Time: {total_hours} hours and {remaining_minutes} minutes")
    if stats['last_updated']:
        print(f"Last Updated: {stats['last_updated']}")
    print("--------------------------\n")

def get_positive_integer_input(prompt, default_value):
    while True:
        user_input = input(f"{prompt} (default: {default_value}): ").strip()
        if not user_input:
            return default_value
        try:
            value = int(user_input)
            if value > 0:
                return value
            else:
                print("Please enter a positive number.")
        except ValueError:
            print("Invalid input. Please enter a number.")

def countdown_timer(duration_seconds, message_prefix):
    start_time = time.time()
    end_time = start_time + duration_seconds

    print(f"\n{message_prefix} started. Press Ctrl+C to stop early.")

    try:
        while time.time() < end_time:
            remaining_seconds = int(end_time - time.time())
            minutes = remaining_seconds // 60
            seconds = remaining_seconds % 60
            print(f"\r{message_prefix}: {minutes:02d}:{seconds:02d} remaining", end="", flush=True)
            time.sleep(1)
        print("\n")

        print(f"{message_prefix} finished!")
        return True, duration_seconds
    except KeyboardInterrupt:
        print("\nTimer stopped by user (Ctrl+C).")
        actual_time_spent = int(time.time() - start_time)
        return False, actual_time_spent

def start_focus_session(stats):
    duration_min = get_positive_integer_input("Enter focus duration in minutes", DEFAULT_FOCUS_DURATION_MIN)
    duration_seconds = duration_min * 60

    print(f"Starting a {duration_min}-minute focus session...")
    completed, actual_time_spent = countdown_timer(duration_seconds, "Focus Session")

    if completed:
        stats["total_focus_sessions"] += 1
        stats["total_focus_time_seconds"] += actual_time_spent
        print("Great job! Focus session completed.")
        save_stats(stats)
        
        take_break_choice = input("Would you like to take a break now? (y/n): ").lower()
        if take_break_choice == 'y':
            start_break_session()
    else:
        stats["total_focus_time_seconds"] += actual_time_spent
        print(f"Focus session ended early. You focused for {actual_time_spent // 60} minutes and {actual_time_spent % 60} seconds.")
        save_stats(stats)

def start_break_session():
    duration_min = get_positive_integer_input("Enter break duration in minutes", DEFAULT_BREAK_DURATION_MIN)
    duration_seconds = duration_min * 60
    print(f"Starting a {duration_min}-minute break...")
    completed, _ = countdown_timer(duration_seconds, "Break")
    if completed:
        print("Break finished. Time to get back to work!")
    else:
        print("Break ended early.")

def main():
    stats = load_stats()
    print("Welcome to the Focus Timer!")

    while True:
        print("\n--- Menu ---")
        print("1. Start Focus Session")
        print("2. Start Break")
        print("3. View Productivity Stats")
        print("4. Exit")
        choice = input("Enter your choice: ").strip()

        if choice == '1':
            start_focus_session(stats)
        elif choice == '2':
            start_break_session()
        elif choice == '3':
            display_stats(stats)
        elif choice == '4':
            print("Exiting Focus Timer. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 03:40:08
import time
import json
import os

class ProductivityStats:
    def __init__(self, filename="productivity_stats.json"):
        self.filename = filename
        self.stats = self._load_stats()

    def _load_stats(self):
        if os.path.exists(self.filename):
            try:
                with open(self.filename, 'r') as f:
                    return json.load(f)
            except json.JSONDecodeError:
                print("Warning: Could not decode stats file. Starting fresh.")
                return self._default_stats()
        return self._default_stats()

    def _default_stats(self):
        return {
            "sessions": [],
            "total_focus_time_seconds": 0,
            "total_break_time_seconds": 0,
            "completed_focus_sessions": 0,
            "completed_break_sessions": 0
        }

    def add_session(self, session_type, start_time_str, end_time_str, duration_seconds):
        session_data = {
            "type": session_type,
            "start_time": start_time_str,
            "end_time": end_time_str,
            "duration_seconds": duration_seconds
        }
        self.stats["sessions"].append(session_data)

        if session_type == "focus":
            self.stats["total_focus_time_seconds"] += duration_seconds
            self.stats["completed_focus_sessions"] += 1
        elif session_type == "break":
            self.stats["total_break_time_seconds"] += duration_seconds
            self.stats["completed_break_sessions"] += 1
        self.save_stats()

    def get_total_focus_time_formatted(self):
        total_seconds = self.stats["total_focus_time_seconds"]
        hours, remainder = divmod(total_seconds, 3600)
        minutes, seconds = divmod(remainder, 60)
        return f"{int(hours)}h {int(minutes)}m {int(seconds)}s"

    def get_total_break_time_formatted(self):
        total_seconds = self.stats["total_break_time_seconds"]
        hours, remainder = divmod(total_seconds, 3600)
        minutes, seconds = divmod(remainder, 60)
        return f"{int(hours)}h {int(minutes)}m {int(seconds)}s"

    def get_average_focus_session_duration_formatted(self):
        if self.stats["completed_focus_sessions"] == 0:
            return "N/A"
        avg_seconds = self.stats["total_focus_time_seconds"] / self.stats["completed_focus_sessions"]
        minutes, seconds = divmod(avg_seconds, 60)
        return f"{int(minutes)}m {int(seconds)}s"

    def display_stats(self):
        print("\n--- Productivity Stats ---")
        print(f"Completed Focus Sessions: {self.stats['completed_focus_sessions']}")
        print(f"Completed Break Sessions: {self.stats['completed_break_sessions']}")
        print(f"Total Focus Time: {self.get_total_focus_time_formatted()}")
        print(f"Total Break Time: {self.get_total_break_time_formatted()}")
        print(f"Average Focus Session: {self.get_average_focus_session_duration_formatted()}")
        print("--------------------------\n")

    def save_stats(self):
        with open(self.filename, 'w') as f:
            json.dump(self.stats, f, indent=4)

class FocusTimer:
    def __init__(self, stats_manager, focus_duration_minutes=25, break_duration_minutes=5):
        self.stats_manager = stats_manager
        self.focus_duration_seconds = focus_duration_minutes * 60
        self.break_duration_seconds = break_duration_minutes * 60
        self.start_time_timestamp = None
        self.end_time_timestamp = None
        self.is_running = False
        self.is_paused = False
        self.remaining_time = 0
        self.current_session_type = None

    def _countdown(self, duration_seconds, session_type):
        self.current_session_type = session_type
        self.remaining_time = duration_seconds
        self.is_running = True
        self.is_paused = False
        self.start_time_timestamp = time.time()

        print(f"\n--- {session_type.capitalize()} Timer Started ({duration_seconds // 60} minutes) ---")

        while self.remaining_time > 0 and self.is_running:
            if not self.is_paused:
                mins, secs = divmod(self.remaining_time, 60)
                timer_display = f"{int(mins):02d}:{int(secs):02d}"
                print(f"\r{session_type.capitalize()} Time Remaining: {timer_display}", end="", flush=True)
                time.sleep(1)
                self.remaining_time -= 1
            else:
                time.sleep(1) # Wait while paused

        self.end_time_timestamp = time.time()
        self.is_running = False
        print("\n") # Newline after countdown finishes

        if self.remaining_time <= 0: # Timer completed naturally
            print(f"--- {session_type.capitalize()} Timer Finished! ---")
            actual_duration = self.end_time_timestamp - self.start_time_timestamp
            self.stats_manager.add_session(
                session_type,
                time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(self.start_time_timestamp)),
                time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(self.end_time_timestamp)),
                actual_duration
            )
        else: # Timer was stopped prematurely (self.is_running became False)
            print(f"--- {session_type.capitalize()} Timer Stopped Prematurely ---")

    def start_focus_session(self):
        if self.is_running:
            print("Timer is already running. Please wait or stop the current session.")
            return
        self._countdown(self.focus_duration_seconds, "focus")

    def start_break_session(self):
        if self.is_running:
            print("Timer is already running. Please wait or stop the current session.")
            return
        self._countdown(self.break_duration_seconds, "break")

    def pause_timer(self):
        if self.is_running and not self.is_paused:
            self.is_paused = True
            print("\nTimer Paused.")
        elif self.is_paused:
            print("Timer is already paused.")
        else:
            print("No timer is running to pause.")

    def resume_timer(self):
        if self.is_running and self.is_paused:
            self.is_paused = False
            print("\nTimer Resumed.")
        elif not self.is_running:
            print("No timer is running to resume.")
        else:
            print("Timer is not paused.")

    def stop_timer(self):
        if self.is_running:
            self.is_running = False
            self.is_paused = False
        else:
            print("No timer is running to stop.")

def main():
    stats = ProductivityStats()
    timer = FocusTimer(stats)

    while True:
        print("\n--- Focus Timer Menu ---")
        print("1. Start Focus Session")
        print("2. Start Break Session")
        print("3. Pause Timer")
        print("4. Resume Timer")
        print("5. Stop Current Timer")
        print("6. View Productivity Stats")
        print("7. Exit")
        choice = input("Enter your choice: ")

        if choice == '1':
            timer.start_focus_session()
        elif choice == '2':
            timer.start_break_session()
        elif choice == '3':
            timer.pause_timer()
        elif choice == '4':
            timer.resume_timer()
        elif choice == '5':
            timer.stop_timer()
        elif choice == '6':
            stats.display_stats()
        elif choice == '7':
            print("Exiting Focus Timer. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 03:41:26
import tkinter as tk
from tkinter import messagebox
import datetime
import math

class FocusTimerApp:
    def __init__(self, master):
        self.master = master
        master.title("Focus Timer with Productivity Stats")
        master.geometry("400x550")
        master.resizable(False, False)
        master.config(bg="#2c3e50") # Dark background

        self.initial_duration = 25 * 60 # Default 25 minutes in seconds
        self.time_left = self.initial_duration
        self.running = False
        self.job = None # To store the after() job ID

        self.current_session_start_time = None
        self.sessions = [] # List of {'start': datetime, 'end': datetime, 'duration': seconds}
        self.total_focus_time = 0 # In seconds
        self.completed_sessions_count = 0

        self._create_widgets()
        self._update_timer_display()
        self._display_stats()

    def _create_widgets(self):
        # Timer Display
        self.timer_label = tk.Label(self.master, text="00:00", font=("Arial", 60, "bold"), fg="#ecf0f1", bg="#2c3e50")
        self.timer_label.pack(pady=20)

        # Duration Input
        duration_frame = tk.Frame(self.master, bg="#2c3e50")
        duration_frame.pack(pady=10)

        tk.Label(duration_frame, text="Set Duration (minutes):", font=("Arial", 12), fg="#ecf0f1", bg="#2c3e50").pack(side=tk.LEFT, padx=5)
        self.duration_entry = tk.Entry(duration_frame, width=10, font=("Arial", 12), bg="#34495e", fg="#ecf0f1", insertbackground="#ecf0f1")
        self.duration_entry.insert(0, str(self.initial_duration // 60))
        self.duration_entry.pack(side=tk.LEFT, padx=5)
        
        set_duration_button = tk.Button(duration_frame, text="Set", command=self._set_duration, font=("Arial", 10), bg="#3498db", fg="#ffffff", activebackground="#2980b9", activeforeground="#ffffff")
        set_duration_button.pack(side=tk.LEFT, padx=5)

        # Control Buttons
        button_frame = tk.Frame(self.master, bg="#2c3e50")
        button_frame.pack(pady=20)

        self.start_button = tk.Button(button_frame, text="Start", command=self._start_timer, font=("Arial", 14, "bold"), bg="#27ae60", fg="#ffffff", activebackground="#229954", activeforeground="#ffffff", width=8)
        self.start_button.pack(side=tk.LEFT, padx=10)

        self.pause_button = tk.Button(button_frame, text="Pause", command=self._pause_timer, font=("Arial", 14, "bold"), bg="#f39c12", fg="#ffffff", activebackground="#e67e22", activeforeground="#ffffff", width=8)
        self.pause_button.pack(side=tk.LEFT, padx=10)

        self.reset_button = tk.Button(button_frame, text="Reset", command=self._reset_timer, font=("Arial", 14, "bold"), bg="#e74c3c", fg="#ffffff", activebackground="#c0392b", activeforeground="#ffffff", width=8)
        self.reset_button.pack(side=tk.LEFT, padx=10)

        # Stats Display
        stats_frame = tk.LabelFrame(self.master, text="Productivity Stats", font=("Arial", 14, "bold"), fg="#ecf0f1", bg="#34495e", bd=2, relief="groove")
        stats_frame.pack(pady=20, padx=20, fill=tk.X)

        self.total_time_label = tk.Label(stats_frame, text="Total Focus Time: 0h 0m 0s", font=("Arial", 12), fg="#ecf0f1", bg="#34495e", anchor="w")
        self.total_time_label.pack(pady=5, padx=10, fill=tk.X)

        self.completed_sessions_label = tk.Label(stats_frame, text="Completed Sessions: 0", font=("Arial", 12), fg="#ecf0f1", bg="#34495e", anchor="w")
        self.completed_sessions_label.pack(pady=5, padx=10, fill=tk.X)

        self.avg_session_label = tk.Label(stats_frame, text="Avg Session Length: 0m 0s", font=("Arial", 12), fg="#ecf0f1", bg="#34495e", anchor="w")
        self.avg_session_label.pack(pady=5, padx=10, fill=tk.X)

    def _format_time(self, seconds):
        minutes, seconds = divmod(seconds, 60)
        hours, minutes = divmod(minutes, 60)
        return f"{int(hours):02d}:{int(minutes):02d}:{int(seconds):02d}" if hours > 0 else f"{int(minutes):02d}:{int(seconds):02d}"

    def _update_timer_display(self):
        self.timer_label.config(text=self._format_time(self.time_left))

    def _set_duration(self):
        if self.running:
            messagebox.showwarning("Timer Running", "Cannot change duration while timer is running. Please pause or reset.")
            return
        try:
            minutes = int(self.duration_entry.get())
            if minutes <= 0:
                raise ValueError("Duration must be a positive number.")
            self.initial_duration = minutes * 60
            self.time_left = self.initial_duration
            self._update_timer_display()
        except ValueError as e:
            messagebox.showerror("Invalid Input", f"Please enter a valid positive integer for minutes.\n{e}")

    def _start_timer(self):
        if self.running:
            return

        if self.time_left <= 0:
            messagebox.showinfo("Timer Finished", "Please reset the timer or set a new duration to start.")
            return

        self.running = True
        self.start_button.config(state=tk.DISABLED)
        self.pause_button.config(state=tk.NORMAL)
        self.reset_button.config(state=tk.NORMAL)

        if self.current_session_start_time is None or self.time_left == self.initial_duration:
            self.current_session_start_time = datetime.datetime.now()
        
        self._tick()

    def _pause_timer(self):
        if not self.running:
            return
        
        self.running = False
        if self.job:
            self.master.after_cancel(self.job)
            self.job = None
        self.start_button.config(state=tk.NORMAL)
        self.pause_button.config(state=tk.DISABLED)

    def _reset_timer(self):
        self._pause_timer() # Stop any running timer
        self.time_left = self.initial_duration
        self._update_timer_display()
        self.start_button.config(state=tk.NORMAL)
        self.pause_button.config(state=tk.NORMAL) 
        self.reset_button.config(state=tk.DISABLED) 

        # Clear current session start time if it was an interrupted session
        self.current_session_start_time = None

    def _tick(self):
        if self.running and self.time_left > 0:
            self.time_left -= 1
            self._update_timer_display()
            self.job = self.master.after(1000, self._tick) # Schedule next tick
        elif self.running and self.time_left <= 0:
            self._end_session()

    def _end_session(self):
        self.running = False
        if self.job:
            self.master.after_cancel(self.job)
            self.job = None
        
        end_time = datetime.datetime.now()
        
        if self.current_session_start_time: # Only record if a session was actually started
            session_duration = self.initial_duration # The full duration set for the session
            
            self.sessions.append({
                'start': self.current_session_start_time,
                'end': end_time,
                'duration': session_duration
            })
            self.total_focus_time += session_duration
            self.completed_sessions_count += 1
            
            self._display_stats()
            messagebox.showinfo("Session Complete!", "Great job! Your focus session has ended.")
        
        self.time_left = 0 # Ensure display shows 00:00
        self._update_timer_display()
        self.start_button.config(state=tk.NORMAL)
        self.pause_button.config(state=tk.DISABLED)
        self.reset_button.config(state=tk.NORMAL)
        self.current_session_start_time = None # Reset for next session

    def _display_stats(self):
        # Total Focus Time
        total_hours, remainder = divmod(self.total_focus_time, 3600)
        total_minutes, total_seconds = divmod(remainder, 60)
        self.total_time_label.config(text=f"Total Focus Time: {int(total_hours)}h {int(total_minutes)}m {int(total_seconds)}s")

        # Completed Sessions
        self.completed_sessions_label.config(text=f"Completed Sessions: {self.completed_sessions_count}")

        # Average Session Length
        if self.completed_sessions_count > 0:
            avg_seconds = self.total_focus_time / self.completed_sessions_count
            avg_minutes, avg_seconds_rem = divmod(avg_seconds, 60)
            self.avg_session_label.config(text=f"Avg Session Length: {int(avg_minutes)}m {int(avg_seconds_rem)}s")
        else:
            self.avg_session_label.config(text="Avg Session Length: 0m 0s")

if __name__ == "__main__":
    root = tk.Tk()
    app = FocusTimerApp(root)
    root.mainloop()