def generate_pascals_triangle(num_rows):
    triangle = []
    for row_num in range(num_rows):
        current_row = [1] * (row_num + 1)
        if row_num > 1:
            prev_row = triangle[row_num - 1]
            for i in range(1, row_num):
                current_row[i] = prev_row[i - 1] + prev_row[i]
        triangle.append(current_row)
    return triangle

if __name__ == "__main__":
    rows = 5
    pascals_triangle = generate_pascals_triangle(rows)
    for row in pascals_triangle:
        print(row)

# Additional implementation at 2025-06-23 03:19:32
def generate_pascals_triangle(num_rows):
    triangle = []
    for i in range(num_rows):
        current_row = [1]
        if i > 0:
            previous_row = triangle[i - 1]
            for j in range(1, i):
                current_row.append(previous_row[j - 1] + previous_row[j])
            current_row.append(1)
        triangle.append(current_row)
    return triangle

def print_pascals_triangle(triangle):
    if not triangle:
        return

    last_row_str = ' '.join(map(str, triangle[-1]))
    max_width = len(last_row_str)

    for row in triangle:
        row_str = ' '.join(map(str, row))
        print(row_str.center(max_width))

if __name__ == "__main__":
    rows_to_generate = 6
    p_triangle = generate_pascals_triangle(rows_to_generate)
    print_pascals_triangle(p_triangle)

    rows_to_generate_2 = 1
    p_triangle_2 = generate_pascals_triangle(rows_to_generate_2)
    print_pascals_triangle(p_triangle_2)

    rows_to_generate_0 = 0
    p_triangle_0 = generate_pascals_triangle(rows_to_generate_0)
    print_pascals_triangle(p_triangle_0)

# Additional implementation at 2025-06-23 03:20:12
def generate_pascals_triangle(num_rows):
    triangle = []
    for i in range(num_rows):
        current_row = [1]
        if i > 0:
            previous_row = triangle[i-1]
            for j in range(len(previous_row) - 1):
                current_row.append(previous_row[j] + previous_row[j+1])
            current_row.append(1)
        triangle.append(current_row)
    return triangle

def print_pascals_triangle(triangle):
    if not triangle:
        return

    last_row_str = ' '.join(map(str, triangle[-1]))
    max_width = len(last_row_str)

    for row in triangle:
        row_str = ' '.join(map(str, row))
        print(row_str.center(max_width))

def get_row_sum(triangle_data, row_index):
    if not triangle_data or row_index < 0 or row_index >= len(triangle_data):
        return None
    return sum(triangle_data[row_index])

def get_specific_row(triangle_data, row_index):
    if not triangle_data or row_index < 0 or row_index >= len(triangle_data):
        return None
    return triangle_data[row_index]

if __name__ == "__main__":
    try:
        rows_input = int(input("Enter the number of rows for Pascal's triangle: "))
        if rows_input <= 0:
            print("Please enter a positive integer for the number of rows.")
        else:
            pascals_triangle_data = generate_pascals_triangle(rows_input)
            print("\n--- Pascal's Triangle ---")
            print_pascals_triangle(pascals_triangle_data)

            print("\n--- Additional Functionality ---")
            
            target_sum_row_index = min(rows_input - 1, 2)
            row_sum_result = get_row_sum(pascals_triangle_data, target_sum_row_index)
            if row_sum_result is not None:
                print(f"Sum of elements in row {target_sum_row_index} (0-indexed): {row_sum_result}")
            else:
                print(f"Could not get sum for row {target_sum_row_index}.")

            target_get_row_index = min(rows_input - 1, 3)
            specific_row_result = get_specific_row(pascals_triangle_data, target_get_row_index)
            if specific_row_result is not None:
                print(f"Elements of row {target_get_row_index} (0-indexed): {specific_row_result}")
            else:
                print(f"Could not retrieve row {target_get_row_index}.")

    except ValueError:
        print("Invalid input. Please enter an integer.")

# Additional implementation at 2025-06-23 03:21:16
def generate_pascals_triangle(num_rows):
    triangle = []
    for i in range(num_rows):
        current_row = [1] * (i + 1)
        if i > 1:
            prev_row = triangle[i - 1]
            for j in range(1, i):
                current_row[j] = prev_row[j - 1] + prev_row[j]
        triangle.append(current_row)
    return triangle

def print_pascals_triangle(triangle):
    if not triangle:
        return

    max_num_width = 0
    for num in triangle[-1]:
        max_num_width = max(max_num_width, len(str(num)))

    last_row_str_length = len(" ".join([str(num).center(max_num_width) for num in triangle[-1]]))

    for row in triangle:
        row_str_elements = []
        for num in row:
            row_str_elements.append(str(num).center(max_num_width))
        
        current_row_str = " ".join(row_str_elements)
        print(current_row_str.center(last_row_str_length))

if __name__ == "__main__":
    fixed_rows = 6
    pascals_triangle_fixed = generate_pascals_triangle(fixed_rows)
    print_pascals_triangle(pascals_triangle_fixed)

    try:
        user_input_rows = int(input("Enter the number of rows for Pascal's triangle: "))
        if user_input_rows < 0:
            print("Number of rows cannot be negative.")
        else:
            pascals_triangle_user = generate_pascals_triangle(user_input_rows)
            print_pascals_triangle(pascals_triangle_user)
    except ValueError:
        print("Invalid input. Please enter an integer.")