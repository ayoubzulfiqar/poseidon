import collections
import time

class RollingFrequencyCounter:
    def __init__(self, window_seconds):
        self.window_seconds = window_seconds
        self.data = collections.deque() # Stores (timestamp, item) tuples
        self.frequencies = collections.Counter()

    def _prune_old_items(self, current_time):
        # Remove items that are older than the window
        while self.data and self.data[0][0] <= current_time - self.window_seconds:
            old_timestamp, old_item = self.data.popleft()
            self.frequencies[old_item] -= 1
            if self.frequencies[old_item] == 0:
                del self.frequencies[old_item]

    def add_item(self, item, timestamp=None):
        # Add a new item with its timestamp
        current_time = timestamp if timestamp is not None else time.time()
        self.data.append((current_time, item))
        self.frequencies[item] += 1
        # Prune old items immediately after adding a new one
        self._prune_old_items(current_time)

    def get_frequencies(self):
        # Return the current frequencies within the window
        # Ensure the window is up-to-date before returning frequencies
        current_time = time.time()
        self._prune_old_items(current_time)
        return dict(self.frequencies) # Return a copy of the frequencies

# Additional implementation at 2025-08-04 07:30:13
import collections

class RollingFrequencyCounter:
    def __init__(self, window_duration_seconds):
        if not isinstance(window_duration_seconds, (int, float)) or window_duration_seconds <= 0:
            raise ValueError("Window duration must be a positive number.")
        self.window_duration = window_duration_seconds
        self.window_data = collections.deque()  # Stores (timestamp, item) tuples
        self.frequencies = collections.Counter()

    def add_item(self, item, timestamp):
        """
        Adds an item with its timestamp to the rolling window.
        Assumes timestamps are generally non-decreasing for efficient pruning.
        """
        if not isinstance(timestamp, (int, float)):
            raise TypeError("Timestamp must be a number.")

        self.window_data.append((timestamp, item))
        self.frequencies[item] += 1
        self._prune_old_items(timestamp)

    def _prune_old_items(self, current_timestamp):
        """
        Removes items from the left of the window that are older than
        current_timestamp - window_duration.
        """
        while self.window_data and self.window_data[0][0] <= current_timestamp - self.window_duration:
            old_timestamp, old_item = self.window_data.popleft()
            self.frequencies[old_item] -= 1
            if self.frequencies[old_item] == 0:
                del self.frequencies[old_item]

    def get_frequencies(self):
        """
        Returns a copy of the current frequency counter for items in the window.
        """
        return self.frequencies.copy()

    def get_count(self, item):
        """
        Returns the count of a specific item within the current window.
        """
        return self.frequencies.get(item, 0)

    def most_common(self, n=None):
        """
        Returns the n most common items and their counts in the window.
        If n is None, returns all items in order of most common to least.
        """
        return self.frequencies.most_common(n)

    def total_items_in_window(self):
        """
        Returns the total number of items (including duplicates) currently in the window.
        """
        return len(self.window_data)

    def clear(self):
        """
        Clears all items and resets the counter.
        """
        self.window_data.clear()
        self.frequencies.clear()