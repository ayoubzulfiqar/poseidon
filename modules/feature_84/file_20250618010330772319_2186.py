import pyperclip
import time

def monitor_clipboard():
    last_clipboard_content = pyperclip.paste()
    while True:
        current_clipboard_content = pyperclip.paste()
        if current_clipboard_content != last_clipboard_content:
            print(f"Clipboard changed: {current_clipboard_content}")
            last_clipboard_content = current_clipboard_content
        time.sleep(0.5)

if __name__ == "__main__":
    monitor_clipboard()

# Additional implementation at 2025-06-18 01:04:20
import pyperclip
import time
import datetime

def monitor_clipboard(interval=1, history_limit=10):
    """
    Monitors clipboard changes and prints new content.
    Maintains a history of copied items.

    Args:
        interval (int): How often to check the clipboard in seconds.
        history_limit (int): Maximum number of items to keep in history.
    """
    last_clipboard_content = ""
    clipboard_history = []

    print("Monitoring clipboard changes... (Press Ctrl+C to stop)")

    try:
        while True:
            current_clipboard_content = pyperclip.paste()

            if current_clipboard_content != last_clipboard_content:
                timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                print(f"[{timestamp}] New clipboard content detected:")
                print("--------------------------------------------------")
                print(current_clipboard_content)
                print("--------------------------------------------------\n")

                # Update history
                clipboard_history.append({
                    "timestamp": timestamp,
                    "content": current_clipboard_content
                })

                # Keep history within limit
                if len(clipboard_history) > history_limit:
                    clipboard_history.pop(0) # Remove oldest item

                last_clipboard_content = current_clipboard_content

            time.sleep(interval)

    except KeyboardInterrupt:
        print("\nClipboard monitoring stopped.")
        if clipboard_history:
            print("\n--- Clipboard History (Last {} items) ---".format(len(clipboard_history)))
            for i, item in enumerate(clipboard_history):
                # Truncate long content for display
                display_content = item['content'][:100]
                if len(item['content']) > 100:
                    display_content += '...'
                print(f"{i+1}. [{item['timestamp']}]\n{display_content}\n")
        else:
            print("No clipboard changes recorded during this session.")

if __name__ == "__main__":
    monitor_clipboard(interval=1, history_limit=5) # Check every 1 second, keep last 5 items in history
