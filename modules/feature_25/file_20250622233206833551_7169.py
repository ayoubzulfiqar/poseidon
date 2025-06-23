def binary_search_early_termination(arr, target):
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

if __name__ == "__main__":
    # Example usage:
    sorted_list = [2, 5, 8, 12, 16, 23, 38, 56, 72, 91]

    # Test cases for existing elements
    print(binary_search_early_termination(sorted_list, 12))
    print(binary_search_early_termination(sorted_list, 2))
    print(binary_search_early_termination(sorted_list, 91))
    print(binary_search_early_termination(sorted_list, 23))

    # Test cases for non-existing elements
    print(binary_search_early_termination(sorted_list, 30))
    print(binary_search_early_termination(sorted_list, 1))
    print(binary_search_early_termination(sorted_list, 100))

    # Test with an empty list
    print(binary_search_early_termination([], 5))

    # Test with a single element list
    print(binary_search_early_termination([7], 7))
    print(binary_search_early_termination([7], 10))

# Additional implementation at 2025-06-22 23:32:58
def binary_search_extended(arr, target):
    low = 0
    high = len(arr) - 1

    while low <= high:
        mid = low + (high - low) // 2
        if arr[mid] == target:
            # Target found. Now, find the first occurrence if duplicates exist.
            # This is the "early termination" point for finding *a* match.
            # The subsequent loop refines to the *first* match.
            first_occurrence_index = mid
            while first_occurrence_index > 0 and arr[first_occurrence_index - 1] == target:
                first_occurrence_index -= 1
            return first_occurrence_index, True
        elif arr[mid] < target:
            low = mid + 1
        else: # arr[mid] > target
            high = mid - 1

    # Target not found. 'low' indicates the insertion point
    # where the target would be placed to maintain sorted order.
    return low, False