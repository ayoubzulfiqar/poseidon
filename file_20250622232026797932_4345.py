import sys
import time

def update_progress_bar(current, total, bar_length=50):
    """
    Updates a console-based progress bar with percentage.

    Args:
        current (int): The current progress value.
        total (int): The total value for completion.
        bar_length (int): The length of the progress bar itself (excluding text).
    """
    if total == 0:
        percent = 0
    else:
        percent = (current / total) * 100

    filled_length = int(bar_length * current // total)
    bar = '█' * filled_length + '-' * (bar_length - filled_length)

    sys.stdout.write(f'\rProgress: |{bar}| {percent:.2f}% ({current}/{total})')
    sys.stdout.flush()

    if current == total:
        sys.stdout.write('\n') # Move to the next line when complete

if __name__ == "__main__":
    total_items = 100
    print("Starting task...")
    for i in range(total_items + 1):
        update_progress_bar(i, total_items)
        time.sleep(0.05) # Simulate work being done
    print("Task completed!")

    print("\nStarting another task...")
    total_steps = 50
    for step in range(total_steps + 1):
        update_progress_bar(step, total_steps, bar_length=30)
        time.sleep(0.1)
    print("Another task completed!")

# Additional implementation at 2025-06-22 23:21:10
import sys
import time

class ProgressBar:
    def __init__(self, total, length=50, fill_char='█', empty_char='░',
                 prefix='Progress:', suffix='Complete', clear_on_finish=True):
        if not isinstance(total, (int, float)) or total <= 0:
            raise ValueError("Total must be a positive number.")
        if not isinstance(length, int) or length <= 0:
            raise ValueError("Length must be a positive integer.")

        self.total = total
        self.length = length
        self.fill_char = fill_char
        self.empty_char = empty_char
        self.prefix = prefix
        self.suffix = suffix
        self.clear_on_finish = clear_on_finish
        self.current = 0
        self.start_time = None
        self._last_printed_len = 0

    def __enter__(self):
        self.start_time = time.time()
        self.update(0)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if exc_type is None:
            self.update(self.total)
        self.finish()

    def update(self, current_value):
        self.current = max(0, min(current_value, self.total)) # Clamp value between 0 and total
        
        percent = (self.current / self.total) * 100
        filled_length = int(self.length * self.current // self.total)
        bar = self.fill_char * filled_length + self.empty_char * (self.length - filled_length)

        elapsed_time = time.time() - self.start_time if self.start_time else 0
        
        eta_str = ""
        if self.current > 0 and elapsed_time > 0:
            items_per_second = self.current / elapsed_time
            if items_per_second > 0:
                remaining_items = self.total - self.current
                eta_seconds = remaining_items / items_per_second
                eta_str = f" ETA: {self._format_time(eta_seconds)}"
        
        elapsed_str = f" Elapsed: {self._format_time(elapsed_time)}"

        output = f"\r{self.prefix} |{bar}| {percent:.1f}% {self.current}/{self.total} {elapsed_str}{eta_str} {self.suffix}"
        
        sys.stdout.write(output.ljust(self._last_printed_len))
        self._last_printed_len = len(output)
        sys.stdout.flush()

    def _format_time(self, seconds):
        if seconds < 60:
            return f"{seconds:.1f}s"
        minutes = seconds / 60
        if minutes < 60:
            return f"{minutes:.1f}m"
        hours = minutes / 60
        return f"{hours:.1f}h"

    def finish(self):
        if self.clear_on_finish:
            sys.stdout.write("\r" + " " * self._last_printed_len + "\r")
        else:
            sys.stdout.write("\n")
        sys.stdout.flush()