def binary_search(arr, target):
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

# Additional implementation at 2025-06-20 00:51:54
def binary_search_extended(arr, target, key=None):
    if not arr:
        return -1
    low = 0
    high = len(arr) - 1
    while low <= high:
        mid = (low + high) // 2
        mid_val = arr[mid]
        compare_val = key(mid_val) if key else mid_val
        if compare_val == target:
            return mid
        elif compare_val < target:
            low = mid + 1
        else:
            high = mid - 1
    return -1

# Additional implementation at 2025-06-20 00:52:35
def binary_search_extended(arr, target, key=None):
    if not arr:
        return False, 0, 0

    low = 0
    high = len(arr) - 1
    comparisons = 0

    while low <= high:
        mid = (low + high) // 2
        
        current_val = arr[mid]
        if key:
            current_val = key(current_val)

        comparisons += 1 

        if current_val == target:
            return True, mid, comparisons
        elif current_val < target:
            low = mid + 1
        else: # current_val > target
            high = mid - 1

    return False, low, comparisons