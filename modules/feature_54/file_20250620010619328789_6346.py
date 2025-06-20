import sys
import time

def show_progress(current, total, bar_length=50):
    """
    Displays a console-based progress bar with percentage.

    Args:
        current (int): The current progress value.
        total (int): The total value for completion.
        bar_length (int): The length of the progress bar in characters.
    """
    fraction = current / total
    percent = int(fraction * 100)
    filled_length = int(bar_length * fraction)
    bar = '█' * filled_length + '-' * (bar_length - filled_length)
    sys.stdout.write(f'\rProgress: |{bar}| {percent}%')
    sys.stdout.flush()
    if current == total:
        sys.stdout.write('\n')

if __name__ == '__main__':
    total_steps = 100
    print("Starting process...")
    for i in range(total_steps + 1):
        show_progress(i, total_steps)
        time.sleep(0.05) # Simulate work being done
    print("Process finished.")

# Additional implementation at 2025-06-20 01:06:47
import sys
import time

class ProgressBar:
    def __init__(self, total, prefix='', suffix='', decimals=1, length=100, fill='█', empty='-', end_char='\r'):
        self.total = total
        self.prefix = prefix
        self.suffix = suffix
        self.decimals = decimals
        self.length = length
        self.fill = fill
        self.empty = empty
        self.end_char = end_char
        self.current_iteration = 0
        self._start_time = None

    def _format_time(self, seconds):
        if seconds is None or seconds == float('inf'):
            return "N/A"
        m, s = divmod(int(seconds), 60)
        h, m = divmod(m, 60)
        return f"{h:02d}:{m:02d}:{s:02d}"

    def _print_bar(self, iteration):
        if self.total == 0:
            percent = "0.0"
            filled_length = 0
        else:
            percent = ("{0:." + str(self.decimals) + "f}").format(100 * (iteration / float(self.total)))
            filled_length = int(self.length * iteration // self.total)
        
        bar = self.fill * filled_length + self.empty * (self.length - filled_length)
        
        elapsed_time = None
        estimated_remaining_time = None

        if self._start_time is not None:
            elapsed_time = time.time() - self._start_time
            if iteration > 0 and elapsed_time > 0:
                items_per_second = iteration / elapsed_time
                remaining_items = self.total - iteration
                if items_per_second > 0:
                    estimated_remaining_time = remaining_items / items_per_second
                else:
                    estimated_remaining_time = float('inf')
            else:
                estimated_remaining_time = float('inf')

        time_info = f"[{self._format_time(elapsed_time)}<{self._format_time(estimated_remaining_time)}]"

        line = f'{self.prefix} |{bar}| {percent}% {self.suffix} {time_info}'
        
        max_percent_len = 3 + self.decimals + (1 if self.decimals > 0 else 0)
        max_time_info_len = 19 
        
        max_line_length = len(self.prefix) + self.length + len(self.suffix) + max_percent_len + max_time_info_len + 7
        
        sys.stdout.write(line.ljust(max_line_length))
        sys.stdout.write(self.end_char)
        sys.stdout.flush()

    def update(self, iteration):
        if self._start_time is None:
            self._start_time = time.time()

        self.current_iteration = iteration
        
        display_iteration = iteration
        if self.total > 0 and iteration > self.total:
            display_iteration = self.total
        elif self.total == 0:
            display_iteration = 0

        self._print_bar(display_iteration)

    def finish(self):
        if self.total > 0:
            self._print_bar(self.total)
        else:
            self._print_bar(0)
        sys.stdout.write('\n')
        sys.stdout.flush()

    def __enter__(self):
        self._start_time = time.time()
        self.update(0

# Additional implementation at 2025-06-20 01:07:59
import sys
import time

def progress_bar(current, total, bar_length=50, fill_char='█', empty_char='-', prefix='', suffix='', start_time=None):
    percent = (current / total) * 100
    filled_length = int(bar_length * current // total)
    bar = fill_char * filled_length + empty_char * (bar_length - filled_length)

    eta_str = ""
    if start_time is not None:
        if current == 0:
            eta_str = " ETA: Estimating..."
        else:
            elapsed_time = time.time() - start_time
            if elapsed_time > 0:
                items_per_second = current / elapsed_time
                remaining_items = total - current
                if items_per_second > 0:
                    eta_seconds = remaining_items / items_per_second
                    hours = int(eta_seconds // 3600)
                    minutes = int((eta_seconds % 3600) // 60)
                    seconds = int(eta_seconds % 60)
                    eta_str = f" ETA: {hours:02}:{minutes:02}:{seconds:02}"
                else:
                    eta_str = " ETA: --:--:--"
            else:
                eta_str = " ETA: Calculating..."

    sys.stdout.write(f"\r{prefix} |{bar}| {percent:.1f}% {suffix}{eta_str}")
    sys.stdout.flush()

    if current == total:
        sys.stdout.write("\n")

if __name__ == "__main__":
    total_items_task1 = 100
    print("Task 1: Default progress bar")
    start_time_task1 = time.time()
    for i in range(total_items_task1 + 1):
        progress_bar(i, total_items_task1, start_time=start_time_task1)
        time.sleep(0.05)

    print("\nTask 2: Custom progress bar with different characters and length")
    total_items_task2 = 75
    start_time_task2 = time.time()
    for i in range(total_items_task2 + 1):
        progress_bar(i, total_items_task2, bar_length=30, fill_char='#', empty_char='.', prefix='Processing files: ', suffix=' of 75', start_time=start_time_task2)
        time.sleep(0.08)

    print("\nTask 3: Shorter task with custom prefix/suffix and faster updates")
    total_items_task3 = 50
    start_time_task3 = time.time()
    for i in range(total_items_task3 + 1):
        progress_bar(i, total_items_task3, prefix='Downloading: ', suffix=' MB', start_time=start_time_task3)
        time.sleep(0.02)

    print("\nAll demonstration tasks completed!")