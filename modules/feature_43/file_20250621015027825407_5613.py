import collections

def sliding_window_maximum(nums, k):
    if not nums or k <= 0:
        return []

    result = []
    # Deque stores indices
    # The front of the deque stores the index of the maximum element in the current window
    # Elements in the deque are in decreasing order of their corresponding values
    dq = collections.deque()

    for i, num in enumerate(nums):
        # Remove indices that are out of the current window (i.e., too old)
        if dq and dq[0] == i - k:
            dq.popleft()

        # Remove indices from the back whose corresponding values are less than or equal to the current number
        # This ensures that the deque maintains elements in decreasing order
        while dq and nums[dq[-1]] <= num:
            dq.pop()

        # Add the current index to the back of the deque
        dq.append(i)

        # If the window has fully formed (i.e., we have processed at least k elements)
        # The maximum for the current window is at the front of the deque
        if i >= k - 1:
            result.append(nums[dq[0]])

    return result

# Additional implementation at 2025-06-21 01:51:27
import collections

class SlidingWindowMaximumCalculator:
    def __init__(self):
        pass

    def _calculate_extremes(self, nums, k, extreme_type):
        if not isinstance(nums, list):
            raise TypeError("Input 'nums' must be a list.")
        if not all(isinstance(x, (int, float)) for x in nums):
            raise ValueError("All elements in 'nums' must be numbers.")
        if not isinstance(k, int) or k <= 0:
            raise ValueError("Window size 'k' must be a positive integer.")
        if k > len(nums) and len(nums) > 0:
            raise ValueError("Window size 'k' cannot be greater than the length of 'nums'.")
        if not nums:
            return []

        dq = collections.deque()
        result = []

        if extreme_type == 'max':
            compare_func = lambda val_in_dq, current_val: val_in_dq <= current_val
        elif extreme_type == 'min':
            compare_func = lambda val_in_dq, current_val: val_in_dq >= current_val
        else:
            raise ValueError("Internal error: extreme_type must be 'max' or 'min'.")

        for i in range(len(nums)):
            if dq and dq[0] <= i - k:
                dq.popleft()

            while dq and compare_func(nums[dq[-1]], nums[i]):
                dq.pop()

            dq.append(i)

            if i >= k - 1:
                result.append(nums[dq[0]])

        return result

    def calculate_maximums(self, nums, k):
        return self._calculate_extremes(nums, k, extreme_type='max')

    def calculate_minimums(self, nums, k):
        return self._calculate_extremes(nums, k, extreme_type='min')

# Additional implementation at 2025-06-21 01:52:34
import collections

class SlidingWindowMaximumCalculator:
    def calculate(self, nums, k):
        if not nums or k <= 0:
            return []

        n = len(nums)
        if k > n:
            return []

        dq = collections.deque()
        result = []

        for i in range(n):
            # Remove indices from the front of the deque that are outside the current window
            if dq and dq[0] <= i - k:
                dq.popleft()

            # Remove indices from the back of the deque whose corresponding values
            # are less than or equal to the current value (nums[i]).
            # This ensures that the deque stores indices of elements in decreasing order of their values.
            while dq and nums[dq[-1]] <= nums[i]:
                dq.pop()

            # Add the current index to the back of the deque
            dq.append(i)

            # If the window has fully formed (i.e., we have processed at least k elements),
            # the maximum element for the current window is at the front of the deque.
            if i >= k - 1:
                result.append(nums[dq[0]])

        return result