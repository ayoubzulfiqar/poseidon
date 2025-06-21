import os
import time
import logging

def monitor_folder(folder_path, interval_seconds=5):
    """
    Monitors a specified folder for new files and logs changes.

    Args:
        folder_path (str): The path to the folder to monitor.
        interval_seconds (int): The time interval (in seconds) between checks.
    """
    if not os.path.isdir(folder_path):
        logging.error(f"Error: Folder '{folder_path}' does not exist or is not a directory.")
        return

    logging.basicConfig(level=logging.INFO,
                        format='%(asctime)s - %(levelname)s - %(message)s',
                        handlers=[
                            logging.FileHandler("folder_monitor.log"),
                            logging.StreamHandler()
                        ])

    logging.info(f"Starting to monitor folder: {folder_path}")
    logging.info(f"Checking every {interval_seconds} seconds.")

    try:
        previous_files = set(os.listdir(folder_path))
        logging.info(f"Initial files in '{folder_path}': {previous_files}")
    except OSError as e:
        logging.error(f"Error accessing folder '{folder_path}': {e}")
        return

    while True:
        try:
            current_files = set(os.listdir(folder_path))

            new_files = current_files - previous_files
            deleted_files = previous_files - current_files

            if new_files:
                for file_name in new_files:
                    logging.info(f"New file detected: {file_name}")

            if deleted_files:
                for file_name in deleted_files:
                    logging.info(f"File deleted: {file_name}")

            if not new_files and not deleted_files:
                logging.debug("No changes detected.")

            previous_files = current_files

        except OSError as e:
            logging.error(f"Error accessing folder '{folder_path}': {e}")
        except Exception as e:
            logging.error(f"An unexpected error occurred: {e}")

        time.sleep(interval_seconds)

if __name__ == "__main__":
    # IMPORTANT: Replace 'path/to/your/folder' with the actual path you want to monitor.
    # For example, to monitor a 'test_folder' in the same directory as the script:
    # monitored_folder = os.path.join(os.path.dirname(__file__), "test_folder")
    # os.makedirs(monitored_folder, exist_ok=True)

    monitored_folder = "C:/temp/monitor_me" # Example path for Windows
    # monitored_folder = "/tmp/monitor_me" # Example path for Linux/macOS

    os.makedirs(monitored_folder, exist_ok=True)

    monitoring_interval = 5 # seconds

    monitor_folder(monitored_folder, monitoring_interval)

# Additional implementation at 2025-06-21 03:06:05
import os
import time
import logging
import sys

MONITORED_FOLDER = "/tmp/monitor_test_folder"
LOG_FILE = "folder_monitor.log"
POLLING_INTERVAL_SECONDS = 5

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler(LOG_FILE),
        logging.StreamHandler(sys.stdout)
    ]
)

def get_current_folder_state(folder_path):
    state = {}
    if not os.path.exists(folder_path):
        logging.error(f"Monitored folder does not exist: {folder_path}")
        return state
    if not os.path.isdir(folder_path):
        logging.error(f"Monitored path is not a directory: {folder_path}")
        return state

    try:
        for root, _, files in os.walk(folder_path):
            for file_name in files:
                file_path = os.path.join(root, file_name)
                try:
                    stats = os.stat(file_path)
                    state[file_path] = (stats.st_size, stats.st_mtime)
                except FileNotFoundError:
                    logging.warning(f"File disappeared during scan: {file_path}")
                except OSError as e:
                    logging.error(f"Error getting stats for {file_path}: {e}")
    except OSError as e:
        logging.error(f"Error walking directory {folder_path}: {e}")
    return state

def monitor_folder():
    logging.info(f"Starting folder monitor for: {MONITORED_FOLDER}")
    logging.info(f"Logging changes to: {LOG_FILE}")
    logging.info(f"Polling interval: {POLLING_INTERVAL_SECONDS} seconds")

    last_state = get_current_folder_state(MONITORED_FOLDER)
    logging.info("Initial scan complete.")

    try:
        while True:
            time.sleep(POLLING_INTERVAL_SECONDS)
            current_state = get_current_folder_state(MONITORED_FOLDER)

            new_files = set(current_state.keys()) - set(last_state.keys())
            deleted_files = set(last_state.keys()) - set(current_state.keys())
            
            modified_files = []
            for file_path, current_metadata in current_state.items():
                if file_path in last_state:
                    last_metadata = last_state[file_path]
                    if current_metadata != last_metadata:
                        modified_files.append(file_path)

            for file_path in new_files:
                logging.info(f"NEW: {file_path}")
            for file_path in deleted_files:
                logging.info(f"DELETED: {file_path}")
            for file_path in modified_files:
                logging.info(f"MODIFIED: {file_path}")
            
            if not (new_files or deleted_files or modified_files):
                logging.debug("No changes detected.")

            last_state = current_state

    except KeyboardInterrupt:
        logging.info("Monitor stopped by user (Ctrl+C).")
    except Exception as e:
        logging.critical(f"An unexpected error occurred: {e}", exc_info=True)

if __name__ == "__main__":
    if not os.path.exists(MONITORED_FOLDER):
        logging.error(f"Error: Monitored folder '{MONITORED_FOLDER}' does not exist. Please create it or update MONITORED_FOLDER.")
        sys.exit(1)
    if not os.path.isdir(MONITORED_FOLDER):
        logging.error(f"Error: Monitored path '{MONITORED_FOLDER}' is not a directory. Please update MONITORED_FOLDER.")
        sys.exit(1)

    monitor_folder()

# Additional implementation at 2025-06-21 03:07:23
import os
import time
import logging
import argparse
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

class FolderMonitorEventHandler(FileSystemEventHandler):
    def __init__(self, logger):
        super().__init__()
        self.logger = logger

    def on_created(self, event):
        if event.is_directory:
            self.logger.info(f"Directory created: {event.src_path}")
        else:
            self.logger.info(f"File created: {event.src_path}")

    def on_deleted(self, event):
        if event.is_directory:
            self.logger.info(f"Directory deleted: {event.src_path}")
        else:
            self.logger.info(f"File deleted: {event.src_path}")

    def on_modified(self, event):
        # Log modifications only for files to avoid excessive logging from directory changes
        # (e.g., a file inside a directory being modified also triggers a directory modified event)
        if not event.is_directory:
            self.logger.info(f"File modified: {event.src_path}")

    def on_moved(self, event):
        if event.is_directory:
            self.logger.info(f"Directory moved/renamed from {event.src_path} to {event.dest_path}")
        else:
            self.logger.info(f"File moved/renamed from {event.src_path} to {event.dest_path}")

def main():
    parser = argparse.ArgumentParser(description="Monitor a folder for file system changes and log them.")
    parser.add_argument("folder_path", type=str,
                        help="The path to the folder to monitor.")
    parser.add_argument("--log_file", type=str, default="folder_monitor.log",
                        help="The path to the log file. Defaults to 'folder_monitor.log'.")
    parser.add_argument("--log_level", type=str, default="INFO",
                        choices=["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
                        help="Set the logging level. Defaults to INFO.")

    args = parser.parse_args()

    if not os.path.isdir(args.folder_path):
        print(f"Error: The specified folder path '{args.folder_path}' does not exist or is not a directory.")
        return

    log_level_map = {
        "DEBUG": logging.DEBUG,
        "INFO": logging.INFO,
        "WARNING": logging.WARNING,
        "ERROR": logging.ERROR,
        "CRITICAL": logging.CRITICAL
    }
    numeric_log_level = log_level_map.get(args.log_level.upper(), logging.INFO)

    logging.basicConfig(
        level=numeric_log_level,
        format='%(asctime)s - %(levelname)s - %(message)s',
        handlers=[
            logging.FileHandler(args.log_file),
            logging.StreamHandler()
        ]
    )
    logger = logging.getLogger(__name__)

    logger.info(f"Starting folder monitor for: {args.folder_path}")
    logger.info(f"Logging changes to: {args.log_file} with level: {args.log_level}")

    event_handler = FolderMonitorEventHandler(logger)
    observer = Observer()
    observer.schedule(event_handler, args.folder_path, recursive=True)

    observer.start()
    logger.info("Monitor started. Press Ctrl+C to stop.")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        logger.info("Monitor stopping...")
    finally:
        observer.stop()
        observer.join()
        logger.info("Monitor stopped.")

if __name__ == "__main__":
    main()