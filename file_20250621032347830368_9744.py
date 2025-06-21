def binary_search_early_termination(arr, target):
    """
    Implements a binary search algorithm with early termination.

    Args:
        arr (list): A sorted list of elements to search within.
        target: The element to search for.

    Returns:
        int: The index of the target if found, otherwise -1.
    """
    low = 0
    high = len(arr) - 1

    while low <= high:
        mid = (low + high) // 2
        mid_val = arr[mid]

        if mid_val == target:
            return mid  # Early termination: target found
        elif mid_val < target:
            low = mid + 1
        else:  # mid_val > target
            high = mid - 1
    return -1

if __name__ == '__main__':
    # Example Usage (not part of the strict output, but for testing)
    sorted_list = [2, 5, 8, 12, 16, 23, 38, 56, 72, 91]

    # Test cases
    # print(f"Searching for 12: {binary_search_early_termination(sorted_list, 12)}") # Expected: 3
    # print(f"Searching for 2: {binary_search_early_termination(sorted_list, 2)}")   # Expected: 0
    # print(f"Searching for 91: {binary_search_early_termination(sorted_list, 91)}") # Expected: 9
    # print(f"Searching for 1: {binary_search_early_termination(sorted_list, 1)}")   # Expected: -1
    # print(f"Searching for 100: {binary_search_early_termination(sorted_list, 100)}") # Expected: -1
    # print(f"Searching for 23: {binary_search_early_termination(sorted_list, 23)}") # Expected: 5
    # print(f"Searching for 56: {binary_search_early_termination(sorted_list, 56)}") # Expected: 7
    # print(f"Searching for 8: {binary_search_early_termination(sorted_list, 8)}")   # Expected: 2

    # Empty list
    # print(f"Searching in empty list for 5: {binary_search_early_termination([], 5)}") # Expected: -1

    # Single element list
    # print(f"Searching in [7] for 7: {binary_search_early_termination([7], 7)}") # Expected: 0
    # print(f"Searching in [7] for 1: {binary_search_early_termination([7], 1)}") # Expected: -1

# Additional implementation at 2025-06-21 03:24:36
def binary_search_extended(arr, target):
    low = 0
    high = len(arr) - 1

    while low <= high:
        mid = (low + high) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            low = mid + 1
        else:
            high = mid - 1
    return -1

# Additional implementation at 2025-06-21 03:25:19
def binary_search_extended(arr, target):
    low = 0
    high = len(arr) - 1
    while low <= high:
        mid = low + (high - low) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            low = mid + 1
        else:
            high = mid - 1
    return low