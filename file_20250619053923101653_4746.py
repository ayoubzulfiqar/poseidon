from pynput import mouse

def on_move(x, y):
    print(f"Mouse moved to ({x}, {y})")

if __name__ == "__main__":
    print("Tracking mouse movements... Press Ctrl+C to stop.")
    with mouse.Listener(on_move=on_move) as listener:
        listener.join()

# Additional implementation at 2025-06-19 05:40:11
