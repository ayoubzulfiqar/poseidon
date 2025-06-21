import curses
import os
import sys
import stat

class FileExplorer:
    def __init__(self, stdscr):
        self.stdscr = stdscr
        self.current_dir = os.path.abspath(os.getcwd())
        self.selected_index = 0
        self.files = []
        self.update_file_list()

        curses.curs_set(0)
        curses.noecho()
        curses.cbreak()
        self.stdscr.keypad(True)

        if curses.has_colors():
            curses.start_color()
            curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
            curses.init_pair(2, curses.COLOR_WHITE, curses.COLOR_BLACK)
            curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)
            curses.init_pair(4, curses.COLOR_YELLOW, curses.COLOR_BLACK)

    def update_file_list(self):
        try:
            items = os.listdir(self.current_dir)
            dirs = sorted([item for item in items if os.path.isdir(os.path.join(self.current_dir, item))], key=str.lower)
            files = sorted([item for item in items if os.path.isfile(os.path.join(self.current_dir, item))], key=str.lower)
            self.files = ['..'] + dirs + files
            self.selected_index = min(self.selected_index, len(self.files) - 1)
            if self.selected_index < 0 and len(self.files) > 0:
                self.selected_index = 0
            elif len(self.files) == 0:
                self.selected_index = -1
        except OSError as e:
            self.files = ['<Error: ' + str(e) + '>']
            self.selected_index = 0
            h, w = self.stdscr.getmaxyx()
            self.stdscr.addstr(h - 1, 0, f"Error: {e}", curses.color_pair(4))
            self.stdscr.refresh()
            curses.napms(2000)

    def draw(self):
        self.stdscr.clear()
        h, w = self.stdscr.getmaxyx()

        path_str = f"Path: {self.current_dir}"
        self.stdscr.addstr(0, 0, path_str[:w-1], curses.color_pair(4))

        start_row = 2
        display_height = h - start_row - 1

        if not self.files:
            self.stdscr.addstr(start_row, 0, "No items in this directory.", curses.color_pair(2))
            return

        scroll_offset = 0
        if self.selected_index >= display_height:
            scroll_offset = self.selected_index - display_height + 1

        for i, item in enumerate(self.files):
            if i - scroll_offset >= display_height:
                break

            row = start_row + (i - scroll_offset)
            if row >= h - 1:
                break

            is_selected = (i == self.selected_index)
            full_path = os.path.join(self.current_dir, item)

            if item == '..':
                color_pair = curses.color_pair(1)
            elif os.path.isdir(full_path):
                color_pair = curses.color_pair(1)
            else:
                color_pair = curses.color_pair(2)

            display_name = item
            if os.path.isdir(full_path):
                display_name += '/'

            display_name = display_name[:w-1]

            if is_selected:
                self.stdscr.addstr(row, 0, display_name, curses.color_pair(3))
            else:
                self.stdscr.addstr(row, 0, display_name, color_pair)

        self.stdscr.refresh()

    def handle_input(self, key):
        if key == curses.KEY_UP:
            self.selected_index = max(0, self.selected_index - 1)
        elif key == curses.KEY_DOWN:
            self.selected_index = min(len(self.files) - 1, self.selected_index + 1)
        elif key == curses.KEY_ENTER or key in [10, 13]:
            if self.selected_index == -1:
                return True
            selected_item = self.files[self.selected_index]
            if selected_item == '..':
                self.current_dir = os.path.abspath(os.path.join(self.current_dir, os.pardir))
                self.selected_index = 0
                self.update_file_list()
            else:
                full_path = os.path.join(self.current_dir, selected_item)
                if os.path.isdir(full_path):
                    self.current_dir = full_path
                    self.selected_index = 0
                    self.update_file_list()
                elif os.path.isfile(full_path):
                    self.view_file(full_path)
        elif key == ord('q'):
            return False
        elif key == ord('d'):
            if self.selected_index != -1 and self.files[self.selected_index] != '..':
                self.delete_item(os.path.join(self.current_dir, self.files[self.selected_index]))
        elif key == ord('n'):
            self.create_directory()
        elif key == ord('r'):
            if self.selected_index != -1 and self.files[self.selected_index] != '..':
                self.rename_item(os.path.join(self.current_dir, self.files[self.selected_index]))
        elif key == curses.KEY_BACKSPACE or key == 263:
            self.current_dir = os.path.abspath(os.path.join(self.current_dir, os.pardir))
            self.selected_index = 0
            self.update_file_list()

        return True

    def prompt(self, message):
        h, w = self.stdscr.getmaxyx()
        prompt_row = h - 1
        self.stdscr.addstr(prompt_row, 0, message + ": ", curses.color_pair(4))
        self.stdscr.clrtoeol()
        curses.echo()
        curses.curs_set(1)

        input_str = ""
        prompt_len = len(message) + 2
        while True:
            char = self.stdscr.getch()
            if char in [10, 13]:
                break
            elif char == curses.KEY_BACKSPACE or char == 263:
                if len(input_str) > 0:
                    input_str = input_str[:-1]
                    self.stdscr.addstr(prompt_row, prompt_len, ' ' * (len(input_str) + 1))
                    self.stdscr.addstr(prompt_row, prompt_len, input_str)
                    self.stdscr.clrtoeol()
            elif 32 <= char <= 126:
                input_str += chr(char)
                self.stdscr.addstr(prompt_row, prompt_len, input_str)
            
        curses.noecho()
        curses.curs_set(0)
        self.stdscr.addstr(prompt_row, 0, ' ' * w)
        self.stdscr.refresh()
        return input_str

    def view_file(self, filepath):
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
            
            h, w = self.stdscr.getmaxyx()
            viewer_window = curses.newwin(h, w, 0, 0)
            viewer_window.keypad(True)
            viewer_window.scrollok(True)
            viewer_window.addstr(0, 0, f"--- Viewing: {os.path.basename(filepath)} (Press 'q' to exit) ---", curses.color_pair(4))
            
            # Add content line by line, handling wrapping
            y, x = 2, 0
            for line in content.splitlines():
                while line:
                    chunk = line[:w-x]
                    viewer_window.addstr(y, x, chunk)
                    line = line[len(chunk):]
                    x = 0
                    y += 1
                    if y >= h - 1: # Prevent writing past screen
                        break
                if y >= h - 1:
                    break

            viewer_window.refresh()

            while True:
                key = viewer_window.getch()
                if key == ord('q'):
                    break
                elif key == curses.KEY_DOWN:
                    viewer_window.scroll(1)
                    viewer_window.refresh()
                elif key == curses.KEY_UP:
                    viewer_window.scroll(-1)
                    viewer_window.refresh()
                elif key == curses.KEY_NPAGE:
                    viewer_window.scroll(h - 4)
                    viewer_window.refresh()
                elif key == curses.KEY_PPAGE:
                    viewer_window.scroll(-(h - 4))
                    viewer_window.refresh()

        except Exception as e:
            h, w = self.stdscr.getmaxyx()
            self.stdscr.addstr(h - 1, 0, f"Error viewing file: {e}", curses.color_pair(4))
            self.stdscr.refresh()
            curses.napms(2000)
        finally:
            self.stdscr.clear()
            self.draw()

    def delete_item(self, path):
        h, w = self.stdscr.getmaxyx()
        confirm = self.prompt(f"Delete '{os.path.basename(path)}'? (y/N)")
        if confirm.lower() == 'y':
            try:
                if os.path.isdir(path):
                    os.rmdir(path)
                else:
                    os.remove(path)
                self.update_file_list()
            except OSError as e:
                self.stdscr.addstr(h - 1, 0, f"Error deleting: {e}", curses.color_pair(4))
                self.stdscr.refresh()
                curses.napms(2000)
        self.draw()

    def create_directory(self):
        h, w = self.stdscr.getmaxyx()
        new_dir_name = self.prompt("New directory name")
        if new_dir_name:
            new_path = os.path.join(self.current_dir, new_dir_name)
            try:
                os.mkdir(new_path)
                self.update_file_list()
            except OSError as e:
                self.stdscr.addstr(h - 1, 0, f"Error creating directory: {e}", curses.color_pair(4))
                self.stdscr.refresh()
                curses.napms(2000)
        self.draw()

    def rename_item(self, old_path):
        h, w = self.stdscr.getmaxyx()
        old_name = os.path.basename(old_path)
        new_name = self.prompt(f"Rename '{old_name}' to")
        if new_name:
            new_path = os.path.join(os.path.dirname(old_path), new_name)
            try:
                os.rename(old_path, new_path)
                self.update_file_list()
            except OSError as e:
                self

# Additional implementation at 2025-06-21 01:03:03


# Additional implementation at 2025-06-21 01:04:39
import curses
import os
import sys
import shutil
import stat

class FileExplorer:
    def __init__(self, stdscr):
        self.stdscr = stdscr
        curses.curs_set(0)  # Hide cursor
        curses.noecho()     # Don't echo key presses
        self.stdscr.keypad(True) # Enable special keys (like arrow keys)

        self.current_dir = os.getcwd()
        self.files = []
        self.selected_index = 0
        self.top_visible_index = 0 # For scrolling

        self.message = "" # Status message for user feedback

        self._update_file_list()

    def _update_file_list(self):
        try:
            items = os.listdir(self.current_dir)
            # Sort directories first, then files, case-insensitively
            dirs = sorted([d for d in items if os.path.isdir(os.path.join(self.current_dir, d))], key=str.lower)
            files = sorted([f for f in items if os.path.isfile(os.path.join(self.current_dir, f))], key=str.lower)
            self.files = ['..'] + dirs + files # Add '..' for parent directory
            
            # Adjust selected_index if it's out of bounds after update
            if not self.files:
                self.selected_index = -1
            elif self.selected_index >= len(self.files):
                self.selected_index = len(self.files) - 1
            elif self.selected_index < 0:
                self.selected_index = 0

        except OSError as e:
            self.message = f"Error: {e}"
            self.files = ['..'] # Fallback to just '..' if directory is unreadable
            self.selected_index = 0

    def _draw_screen(self):
        self.stdscr.clear()
        height, width = self.stdscr.getmaxyx()

        # Header: Current Directory
        self.stdscr.addstr(0, 0, f"Current Directory: {self.current_dir}", curses.A_BOLD)

        # File List
        display_start_row = 2
        display_height = height - display_start_row - 3 # Leave space for header, status, and message

        if self.selected_index < self.top_visible_index:
            self.top_visible_index = self.selected_index
        elif self.selected_index >= self.top_visible_index + display_height:
            self.top_visible_index = self.selected_index - display_height + 1

        for i in range(display_height):
            file_index = self.top_visible_index + i
            if file_index < len(self.files):
                filename = self.files[file_index]
                display_name = filename
                if len(display_name) > width - 2: # Truncate if too long
                    display_name = display_name[:width - 5] + "..."

                is_dir = os.path.isdir(os.path.join(self.current_dir, filename)) if filename != '..' else True
                
                attr = curses.A_NORMAL
                if is_dir:
                    attr |= curses.A_BOLD # Directories are bold
                
                if file_index == self.selected_index:
                    attr |= curses.A_REVERSE # Selected item is highlighted

                self.stdscr.addstr(display_start_row + i, 0, display_name, attr)
            else:
                self.stdscr.addstr(display_start_row + i, 0, "") # Clear line

        # Status/Message Bar
        self.stdscr.addstr(height - 2, 0, self.message.ljust(width - 1), curses.A_NORMAL)
        self.message = "" # Clear message after displaying

        # Footer: Keybindings
        self.stdscr.addstr(height - 1, 0, "Q:Quit | Enter:Open/CD | C:Copy | M:Move | D:Delete | N:New File | K:New Dir | V:View", curses.A_BOLD)

        self.stdscr.refresh()

    def _get_input_dialog(self, prompt, default_value=""):
        height, width = self.stdscr.getmaxyx()
        input_win = curses.newwin(3, width - 4, height // 2 - 1, 2)
        input_win.box()
        input_win.addstr(0, 2, prompt, curses.A_BOLD)
        input_win.addstr(1, 1, default_value)
        input_win.refresh()

        curses.echo() # Enable echoing for input
        curses.curs_set(1) # Show cursor

        input_str = default_value
        while True:
            input_win.addstr(1, 1, input_str.ljust(width - 7)) # Clear previous input
            input_win.move(1, 1 + len(input_str))
            input_win.refresh()
            
            key = self.stdscr.getch()
            if key == 10: # Enter key
                break
            elif key == curses.KEY_BACKSPACE or key == 127: # Backspace
                input_str = input_str[:-1]
            elif 32 <= key <= 126: # Printable ASCII characters
                input_str += chr(key)
            elif key == 27: # ESC key
                input_str = "" # Cancel operation
                break
        
        curses.noecho()
        curses.curs_set(0)
        return input_str.strip()

    def _change_dir(self, path):
        try:
            os.chdir(path)
            self.current_dir = os.getcwd()
            self._update_file_list()
            self.selected_index = 0 # Reset selection when changing directory
            self.top_visible_index = 0
            self.message = f"Changed directory to: {self.current_dir}"
        except OSError as e:
            self.message = f"Error changing directory: {e}"

    def _view_file(self, filepath):
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
            
            height, width = self.stdscr.getmaxyx()
            viewer_win = curses.newwin(height - 4, width - 4, 2, 2)
            viewer_win.box()
            viewer_win.addstr(0, 2, f"Viewing: {os.path.basename(filepath)} (Press any key to exit)", curses.A_BOLD)
            
            # Display content, line by line, handling scrolling
            lines = content.splitlines()
            current_line = 0
            
            while True:
                viewer_win.clear()
                viewer_win.box()
                viewer_win.addstr(0, 2, f"Viewing: {os.path.basename(filepath)} (Press any key to exit)", curses.A_BOLD)
                
                for i in range(min(len(lines) - current_line, height - 6)): # -6 for borders and header/footer
                    line_to_display = lines[current_line + i]
                    if len(line_to_display) > width - 6:
                        line_to_display = line_to_display[:width - 9] + "..."
                    viewer_win.addstr(2 + i, 2, line_to_display)
                
                viewer_win.refresh()
                key = self.stdscr.getch()
                if key == curses.KEY_UP and current_line > 0:
                    current_line -= 1
                elif key == curses.KEY_DOWN and current_line + (height - 6) < len(lines):
                    current_line += 1
                else:
                    break # Any other key exits

        except Exception as e:
            self.message = f"Error viewing file: {e}"
        finally:
            self.stdscr.clear() # Clear viewer window
            self._draw_screen() # Redraw main explorer

    def _copy_file(self, src_path):
        dest_path = self._get_input_dialog(f"Copy '{os.path.basename(src_path)}' to:", self.current_dir)
        if not dest_path:
            self.message = "Copy cancelled."
            return

        try:
            if os.path.isdir(dest_path): # If destination is a directory, copy into it
                shutil.copy(src_path, dest_path)
            else: # If destination is a file path, copy and rename
                shutil.copyfile(src_path, dest_path)
            self.message = f"Copied '{os.path.basename(src_path)}' to '{dest_path}'"
            self._update_file_list()
        except Exception as e:
            self.message = f"Error copying: {e}"

    def _move_file(self, src_path):
        dest_path = self._get_input_dialog(f"Move '{os.path.basename(src_path)}' to:", self.current_dir)
        if not dest_path:
            self.message = "Move cancelled."
            return

        try:
            shutil.move(src_path, dest_path)
            self.message = f"Moved '{os.path.basename(src_path)}' to '{dest_path}'"
            self._update_file_list()
        except Exception as e:
            self.message = f"Error moving: {e}"

    def _delete_item(self, item_path):
        confirm = self._get_input_dialog(f"Are you sure you want to delete '{os.path.basename(item_path)}'? (yes/no)", "no")
        if confirm.lower() == 'yes':
            try:
                if os.path.isdir(item_path):
                    shutil.rmtree(item_path)
                else:
                    os.remove(item_path)
                self.message = f"Deleted '{os.path.basename(item_path)}'"
                self._update_file_list()
            except Exception as e:
                self.message = f"Error deleting: {e}"
        else:
            self.message = "Delete cancelled."

    def _create_dir(self):
        dir_name = self._get_input_dialog("Enter new directory name:")
        if not dir_name:
            self.message = "Create directory cancelled."
            return
        
        new_dir_path = os.path.join(self.current_dir, dir_name)
        try:
            os.makedirs(new_dir_path)
            self.message = f"Created directory: '{dir_name}'"
            self._update_file_list()
        except Exception as e:
            self.message = f"Error creating directory: {e}"

    def _create_file(self):
        file_name = self._get_input_dialog("Enter new file name:")
        if not file_name:
            self.message = "Create file cancelled."
            return
        
        new_file_path = os.path.join(self.current_dir, file_name)
        try:
            with open(new_file_path, 'w') as f:
                pass # Create empty file
            self.message = f"Created file: '{file_name}'"
            self._update_file_list()
        except Exception as e:
            self.message = f"Error creating file: {e}"

    def run(self):
        while True:
            self._draw_screen()
            key = self.stdscr.getch()

            if key == ord('q'):
                break
            elif key == curses.KEY_UP:
                if self.selected_index > 0:
                    self.selected_index -= 1
            elif key == curses.KEY_DOWN:
                if self.selected_index < len(self.files) - 1:
                    self.selected_index += 1
            elif key == 10: # Enter key
                if self.selected_index != -1:
                    selected_item = self.files[self.selected_index]
                    full_path = os.path.join(self.current_dir, selected_item)
                    
                    if selected_item == '..':
                        self._change_dir(os.path.abspath(os.path.join(self.current_dir, os.pardir)))
                    elif os.path.isdir(full_path):
                        self._change_dir(full_path)
                    elif os.path.isfile(full_path):
                        self._view_file(full_path)
                    else:
                        self.message = "Cannot open this item."
            elif key == ord('c'): # Copy
                if self.selected_index != -1 and self.files[self.selected_index] != '..':
                    selected_item = self.files[self.selected_index]
                    full_path = os.path.join(self.current_dir, selected_item)
                    self._copy_file(full_path)
                else:
                    self.message = "No item selected or '..' cannot be copied."
            elif key == ord('m'): # Move/Rename
                if self.selected_index != -1 and self.files[self.selected_index] != '..':
                    selected_item = self.files[self.selected_index]
                    full_path = os.path.join(self.current_dir, selected_item)
                    self._move_file(full_path)
                else:
                    self.message = "No item selected or '..' cannot be moved."
            elif key == ord('d'): # Delete
                if self.selected_index != -1 and self.files[self.selected_index] != '..':
                    selected_item = self.files[self.selected_index]
                    full_path = os.path.join(self.current_dir, selected_item)
                    self._delete_item(full_path)
                else:
                    self.message = "No item selected or '..' cannot be deleted."
            elif key == ord('k'): # Create Directory
                self

# Additional implementation at 2025-06-21 01:05:29
import os
import sys
import datetime
import shutil # For rmtree

# --- Platform-specific input handling ---
if os.name == 'nt': # Windows
    import msvcrt
    def read_single_key():
        # msvcrt.getch() returns bytes, decode it
        # For arrow keys, it returns two bytes: b'\xe0' followed by a second byte
        key = msvcrt.getch()
        if key == b'\xe0': # Special key (arrow, F-keys, etc.)
            key = msvcrt.getch()
            if key == b'H': return 'UP'
            if key == b'P': return 'DOWN'
            if key == b'M': return 'RIGHT' # Not used yet, but good to have
            if key == b'K': return 'LEFT'  # Not used yet
        elif key == b'\r': # Enter key
            return 'ENTER'
        elif key == b'\x08': # Backspace
            return 'BACKSPACE'
        elif key == b'\x1b': # ESC key
            return 'ESC'
        try:
            return key.decode('utf-8')
        except UnicodeDecodeError:
            return key.decode('latin-1') # Fallback for some characters
else: # Unix-like (Linux, macOS)
    import tty
    import termios
    def read_single_key():
        fd = sys.stdin.fileno()
        old_settings = termios.tcgetattr(fd)
        try:
            tty.setraw(sys.stdin.fileno())
            ch = sys.stdin.read(1)
            if ch == '\x1b': # ESC or arrow key sequence
                ch += sys.stdin.read(2) # Read next two characters for arrow keys
                if ch == '\x1b[A': return 'UP'
                if ch == '\x1b[B': return 'DOWN'
                if ch == '\x1b[C': return 'RIGHT'
                if ch == '\x1b[D': return 'LEFT'
                return 'ESC' # Just ESC
            elif ch == '\x7f': # Backspace
                return 'BACKSPACE'
            elif ch == '\r': # Enter key
                return 'ENTER'
            return ch
        finally:
            termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)

# --- Constants and Colors ---
COLOR_RESET = "\033[0m"
COLOR_BLUE = "\033[94m" # Directories
COLOR_GREEN = "\033[92m" # Files
COLOR_YELLOW = "\033[93m" # Instructions/Warnings
COLOR_RED = "\033[91m" # Errors/Confirmation
COLOR_CYAN = "\033[96m" # Current Path/Headers
COLOR_WHITE_BG = "\033[47m\033[30m" # White background, black text for selection

# --- Utility Functions ---
def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def format_size(size_bytes):
    if size_bytes is None:
        return "N/A"
    if size_bytes == 0:
        return "0 B"
    size_name = ("B", "KB", "MB", "GB", "TB")
    i = 0
    while size_bytes >= 1024 and i < len(size_name) - 1:
        size_bytes /= 1024
        i += 1
    return f"{size_bytes:.1f} {size_name[i]}"

def get_file_info(path):
    try:
        stat_info = os.stat(path)
        is_dir = os.path.isdir(path)
        size = stat_info.st_size if not is_dir else None
        mtime = datetime.datetime.fromtimestamp(stat_info.st_mtime).strftime('%Y-%m-%d %H:%M')
        return is_dir, size, mtime
    except OSError:
        return False, None, None # File not found or permission denied

def get_dir_contents(path):
    contents = []
    try:
        items = os.listdir(path)
        # Sort directories first, then files, both alphabetically
        dirs = sorted([item for item in items if os.path.isdir(os.path.join(path, item))], key=lambda s: s.lower())
        files = sorted([item for item in items if os.path.isfile(os.path.join(path, item))], key=lambda s: s.lower())

        # Add ".." for parent directory if not at the root
        if os.path.abspath(path) != os.path.abspath(os.path.sep):
             contents.append({'name': '..', 'is_dir': True, 'size': None, 'mtime': None})

        for item_name in dirs + files:
            full_path = os.path.join(path, item_name)
            is_dir, size, mtime = get_file_info(full_path)
            contents.append({
                'name': item_name,
                'is_dir': is_dir,
                'size': size,
                'mtime': mtime
            })
    except PermissionError:
        contents.append({'name': f"{COLOR_RED}Permission Denied{COLOR_RESET}", 'is_dir': False, 'size': None, 'mtime': None})
    except FileNotFoundError:
        contents.append({'name': f"{COLOR_RED}Directory Not Found{COLOR_RESET}", 'is_dir': False, 'size': None, 'mtime': None})
    return contents

def display_panel(contents, selected_index, current_path, terminal_height=20):
    clear_screen()
    print(f"{COLOR_CYAN}Current Path: {current_path}{COLOR_RESET}\n")
    print(f"{'Name':<30} {'Size':>10} {'Modified':>20}")
    print("-" * 62)

    # Calculate visible range for scrolling
    # Ensure selected_index is always visible, ideally in the middle
    start_index = max(0, selected_index - terminal_height // 2)
    end_index = min(len(contents), start_index + terminal_height)

    # Adjust start_index if there aren't enough items to fill the screen from the calculated start
    if end_index - start_index < terminal_height:
        start_index = max(0, len(contents) - terminal_height)
        end_index = min(len(contents), start_index + terminal_height)

    for i in range(start_index, end_index):
        item = contents[i]
        name = item['name']
        is_dir = item['is_dir']
        size_str = format_size(item['size'])
        mtime_str = item['mtime'] if item['mtime'] else "N/A"

        display_name = name
        color = COLOR_BLUE if is_dir else COLOR_GREEN

        if i == selected_index:
            line = f"{COLOR_WHITE_BG}{display_name:<30} {size_str:>10} {mtime_str:>20}{COLOR_RESET}"
        else:
            line = f"{color}{display_name:<30} {size_str:>10} {mtime_str:>20}{COLOR_RESET}"
        print(line)

    # Fill remaining lines if contents are fewer than terminal_height
    for _ in range(terminal_height - (end_index - start_index)):
        print("") # Print empty line

    print("\n" + "-" * 62)
    print(f"{COLOR_YELLOW}UP/DOWN: Navigate | ENTER: Open/Enter Dir | BACKSPACE/u: Go Up | v: View | d: Delete | q/ESC: Exit{COLOR_RESET}")

def view_file(file_path):
    clear_screen()
    try:
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            print(f"{COLOR_CYAN}--- Viewing: {file_path} ---{COLOR_RESET}\n")
            # Read and print first N lines or until EOF
            for i, line in enumerate(f):
                print(line.rstrip())
                if i >= 50: # Limit output to 50 lines for large files
                    print(f"\n{COLOR_YELLOW}--- (Truncated) Press any key to return ---{COLOR_RESET}")
                    break
            print(f"\n{COLOR_CYAN}--- End of file ---{COLOR_RESET}")
    except Exception as e:
        print(f"{COLOR_RED}Error viewing file: {e}{COLOR_RESET}")
        print(f"{COLOR_YELLOW}File might be binary or unreadable.{COLOR_RESET}")
    print(f"\n{COLOR_YELLOW}Press any key to return...{COLOR_RESET}")
    read_single_key() # Wait for user input

def delete_item(item_path, is_dir):
    clear_screen()
    confirm = input(f"{COLOR_RED}Are you sure you want to DELETE '{item_path}'? (y/N): {COLOR_RESET}")
    if confirm.lower() == 'y':
        try:
            if is_dir:
                shutil.rmtree(item_path)
                print(f"{COLOR_GREEN}Directory '{item_path}' deleted successfully.{COLOR_RESET}")
            else:
                os.remove(item_path)
                print(f"{COLOR_GREEN}File '{item_path}' deleted successfully.{COLOR_RESET}")
        except OSError as e:
            print(f"{COLOR_RED}Error deleting '{item_path}': {e}{COLOR_RESET}")
    else:
        print(f"{COLOR_YELLOW}Deletion cancelled.{COLOR_RESET}")
    print(f"\n{COLOR_YELLOW}Press any key to return...{COLOR_RESET}")
    read_single_key()

def main():
    current_path = os.getcwd()
    selected_index = 0
    running = True

    # Attempt to get terminal height for better display, default to 20
    try:
        terminal_height = os.get_terminal_size().lines - 7 # Account for header and footer lines
    except OSError:
        terminal_height = 20

    while running:
        contents = get_dir_contents(current_path)
        
        # Handle cases where directory is empty or inaccessible
        if not contents or (len(contents) == 1 and contents[0]['name'] == '..' and os.path.abspath(current_path) == os.path.abspath(os.path.sep)):
            display_panel(contents, selected_index, current_path, terminal_height)
            if not contents:
                print(f"{COLOR_RED}No items or permission denied in this directory.{COLOR_RESET}")
            else: # Only '..' at root
                print(f"{COLOR_YELLOW}This is the root directory.{COLOR_RESET}")
            print(f"{COLOR_YELLOW}Press BACKSPACE/u to go up (if applicable) or q/ESC to exit.{COLOR_RESET}")
            key = read_single_key()
            if key in ('BACKSPACE', 'u'):
                # Only allow going up if not at root
                if os.path.abspath(current_path) != os.path.abspath(os.path.sep):
                    current_path = os.path.dirname(current_path)
                    selected_index = 0
            elif key in ('q