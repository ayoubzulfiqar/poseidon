def plot_ascii(data, width=80, height=20):
    """
    Creates a basic ASCII art plot of the given data.

    Args:
        data (list): A list of numerical values to plot.
        width (int): The total width of the plot in characters (including Y-axis).
        height (int): The total height of the plot in characters (including X-axis).
    """
    if not data:
        print("No data to plot.")
        return

    # Ensure minimum dimensions for a visible plot
    if width < 5:
        width = 5
    if height < 3:
        height = 3

    min_val = min(data)
    max_val = max(data)

    # Initialize grid with spaces
    grid = [[' ' for _ in range(width)] for _ in range(height)]

    # Draw Y-axis (vertical line)
    for y in range(height):
        grid[y][0] = '|'
    
    # Draw X-axis (horizontal line)
    for x in range(width):
        grid[height - 1][x] = '-'
    
    # Draw origin (intersection of axes)
    grid[height - 1][0] = '+'

    # Plot data points
    num_points = len(data)
    
    # Handle cases where all data points are the same value (flat line)
    if min_val == max_val:
        # Plot all points on a horizontal line, roughly in the middle of the plot area
        # or just above the x-axis if height is small.
        # Calculate y_plot_coord relative to the plotting area (height-1 rows)
        y_plot_coord = (height - 1) - (height // 2) 
        
        # Adjust if the calculated position is too close to the top or bottom boundaries
        if y_plot_coord < 0: 
            y_plot_coord = 0 # Ensure it's not above the top boundary
        if y_plot_coord >= height - 1: 
            y_plot_coord = height - 2 # Ensure it's not on the x-axis itself

        for i in range(num_points):
            # Scale x-coordinate across the available plotting width (width - 1 for y-axis)
            # If only one point, place it at x=1 (after y-axis)
            x_plot_coord = int(i / (num_points - 1) * (width - 1)) if num_points > 1 else 1
            
            # Place the point, ensuring it's within bounds and not on the Y-axis itself
            if 0 < x_plot_coord < width and 0 <= y_plot_coord < height:
                grid[y_plot_coord][x_plot_coord] = '*'
    else:
        for i, value in enumerate(data):
            # Scale x-coordinate across the available plotting width (width - 1 for y-axis)
            # If only one point, place it at x=1 (after y-axis)
            x_plot_coord = int(i / (num_points - 1) * (width - 1)) if num_points > 1 else 1
            
            # Scale y-coordinate and invert for display (row 0 is top, height-1 is bottom/x-axis)
            # Map min_val to height-1 (bottom) and max_val to 0 (top)
            scaled_y = (value - min_val) / (max_val - min_val)
            y_plot_coord = (height - 1) - int(scaled_y * (height - 1))

            # Ensure coordinates are within bounds and not on the Y-axis itself
            if 0 < x_plot_coord < width and 0 <= y_plot_coord < height:
                grid[y_plot_coord][x_plot_coord] = '*'

    # Print the grid row by row
    for row in grid:
        print("".join(row))

# --- Example Usage ---

# Example 1: Simple increasing data
print("--- Plot 1: Simple increasing data (width=60, height=15) ---")
plot_ascii([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], width=60, height=15)
print("\n")

# Example 2: Decreasing and then increasing data
print("--- Plot 2: Decreasing and then increasing data (width=60, height=15) ---")
plot_ascii([10, 8, 6, 4, 2, 1, 3, 5, 7, 9], width=60, height=15)
print("\n")

# Example 3: All data points are the same (flat line)
print("--- Plot 3: All data points are the same (width=60, height=15) ---")
plot_ascii([5, 5, 5, 5, 5, 5, 5, 5, 5, 5], width=60, height=15)
print("\n")

# Example 4: Single data point
print("--- Plot 4: Single data point (width=60, height=15) ---")
plot_ascii([7], width=60, height=15)
print("\n")

# Example 5: Empty data
print("--- Plot 5: Empty data (width=60, height=15) ---")
plot_ascii([], width=60, height=15)
print("\n")

# Example 6: More complex data with varying range
print("--- Plot 6: Complex data (width=80, height=20) ---")
plot_ascii([1, 100, 2, 90, 3, 80, 4, 70, 5, 60, 6, 50, 7, 40, 8, 30, 9, 20, 10, 10], width=80, height=20)
print("\n")

# Example 7: Data with negative values
print("--- Plot 7: Data with negative values (width=60, height=15) ---")
plot_ascii([-5, -2, 0, 3, 5, 2, -1, -4], width=60, height=15)
print("\n")

# Example 8: Small plot dimensions
print("--- Plot 8: Small dimensions (width=10, height=5) ---")
plot_ascii([1, 2, 3, 4, 5], width=10, height=5)
print("\n")

# Example 9: Data with floats
print("--- Plot 9: Data with floats (width=60, height=15) ---")
plot_ascii([0.1, 0.5, 0.9, 0.3, 0.7], width=60, height=15)
print("\n")

# Example 10: Large number of points
print("--- Plot 10: Large number of points (width=100, height=25) ---")
plot_ascii([i % 20 for i in range(200)], width=100, height=25)
print("\n")

# Additional implementation at 2025-06-19 22:41:52
def plot_ascii_chart(series_data, title=""):
    PLOT_WIDTH = 70
    PLOT_HEIGHT = 20

    Y_AXIS_LABEL_WIDTH = 7
    X_AXIS_LABEL_HEIGHT = 2
    TITLE_ROW = 0

    PLOT_TOP_ROW = 1
    PLOT_BOTTOM_ROW = PLOT_HEIGHT - X_AXIS_LABEL_HEIGHT - 1
    X_AXIS_ROW = PLOT_BOTTOM_ROW + 1
    X_LABEL_ROW_1 = X_AXIS_ROW + 1
    X_LABEL_ROW_2 = X_AXIS_ROW + 2

    PLOT_LEFT_COL = Y_AXIS_LABEL_WIDTH
    PLOT_RIGHT_COL = PLOT_WIDTH - 1

    plot_area_height = PLOT_BOTTOM_ROW - PLOT_TOP_ROW + 1
    plot_area_width = PLOT_RIGHT_COL - PLOT_LEFT_COL + 1

    EPSILON = 1e-9

    all_x = []
    all_y = []
    for s in series_data:
        all_x.extend(s['x'])
        all_y.extend(s['y'])

    if not all_x or not all_y:
        min_x, max_x = 0.0, 1.0
        min_y, max_y = 0.0, 1.0
    else:
        min_x, max_x = min(all_x), max(all_x)
        min_y, max_y = min(all_y), max(all_y)

    x_range = max_x - min_x
    y_range = max_y - min_y

    if x_range < EPSILON:
        x_range = 1.0
        min_x -= 0.5
        max_x += 0.5
    if y_range < EPSILON:
        y_range = 1.0
        min_y -= 0.5
        max_y += 0.5

    grid = [[' ' for _ in range(PLOT_WIDTH)] for _ in range(PLOT_HEIGHT)]

    for r in range(PLOT_HEIGHT):
        for c in range(PLOT_WIDTH):
            if r == TITLE_ROW:
                continue
            if r == PLOT_TOP_ROW or r == PLOT_BOTTOM_ROW:
                if c >= PLOT_LEFT_COL:
                    grid[r][c] = '-'
            if c == PLOT_LEFT_COL or c == PLOT_RIGHT_COL:
                if r >= PLOT_TOP_ROW and r <= PLOT_BOTTOM_ROW:
                    grid[r][c] = '|'
            if r == X_AXIS_ROW and c >= PLOT_LEFT_COL:
                grid[r][c] = '-'

    grid[PLOT_BOTTOM_ROW][PLOT_LEFT_COL] = '+'
    grid[PLOT_TOP_ROW][PLOT_LEFT_COL] = '+'
    grid[PLOT_BOTTOM_ROW][PLOT_RIGHT_COL] = '+'
    grid[PLOT_TOP_ROW][PLOT_RIGHT_COL] = '+'

    if title:
        padded_title = title.center(PLOT_WIDTH)
        for i, char in enumerate(padded_title):
            if i < PLOT_WIDTH:
                grid[TITLE_ROW][i] = char

    for s_idx, s in enumerate(series_data):
        for i in range(len(s['x'])):
            x_val = s['x'][i]
            y_val = s['y'][i]

            scaled_x = int((x_val - min_x) / x_range * (plot_area_width - 1))
            scaled_y = int((y_val - min_y) / y_range * (plot_area_height - 1))

            grid_x = PLOT_LEFT_COL + scaled_x
            grid_y = PLOT_BOTTOM_ROW - scaled_y

            if (PLOT_TOP_ROW <= grid_y <= PLOT_BOTTOM_ROW and
                    PLOT_LEFT_COL <= grid_x <= PLOT_RIGHT_COL):
                grid[grid_y][grid_x] = s['char']

    y_label_top = "{:.1f}".format(max_y)
    y_label_bottom = "{:.1f}".format(min_y)

    for i, char in enumerate(y_label_top.rjust(Y_AXIS_LABEL_WIDTH)):
        if i < Y_AXIS_LABEL_WIDTH:
            grid[PLOT_TOP_ROW][i] = char
    for i, char in enumerate(y_label_bottom.rjust(Y_AXIS_LABEL_WIDTH)):
        if i < Y_AXIS_LABEL_WIDTH:
            grid[PLOT_BOTTOM_ROW][i] = char

    x_label_left = "{:.1f}".format(min_x)
    x_label_right = "{:.1f}".format(max_x)

    for i, char in enumerate(x_label_left):
        if PLOT_LEFT_COL + i < PLOT_WIDTH:
            grid[X_LABEL_ROW_1][PLOT_LEFT_COL + i] = char
    for i, char in enumerate(x_label_right):
        if PLOT_RIGHT_COL - len(x_label_right) + 1 + i < PLOT_WIDTH:
            grid[X_LABEL_ROW_1][PLOT_RIGHT_COL - len(x_label_right) + 1 + i] = char

    for row in grid:
        print("".join(row))

    if series_data:
        print("Legend:")
        for s in series_data:
            print(f"  {s['char']}: {s['label']}")

series1 = {'x': [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 'y': [10, 12, 8, 15, 11, 13, 9, 16, 14, 10], 'char': '*', 'label': 'Data Set A'}
series2 = {'x': [1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5], 'y': [9, 11, 13, 10, 12, 14, 11, 15, 12], 'char': 'o', 'label': 'Data Set B'}
series3 = {'x': [0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0, 5.5, 6.0, 6.5, 7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0], 'y': [5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6], 'char': '#', 'label': 'Data Set C (Sine-like)'}

plot_ascii_chart([series1, series2, series3], "Multi-Series Scatter Plot Example")

print("\n" * 2)

series_single = {'x': [10, 20, 30, 40, 50], 'y': [100, 120, 90, 150, 110], 'char': 'X', 'label': 'Single Series'}
plot_ascii_chart([series_single], "Single Series Plot")

print("\n" * 2)

series_empty = []
plot_ascii_chart(series_empty, "Empty Plot")

print("\n" * 2)

series_flat_x = {'x': [5, 5, 5, 5], 'y': [1, 2, 3, 4], 'char': 'V', 'label': 'Vertical Line'}
series_flat_y = {'x': [1, 2, 3, 4], 'y': [5, 5, 5, 5], 'char': 'H', 'label': 'Horizontal Line'}
plot_ascii_chart([series_flat_x, series_flat_y], "Flat Line Test")

print("\n" * 2)

series_single_point = {'x': [0], 'y': [0], 'char': '@', 'label': 'Origin'}
plot_ascii_chart([series_single_point], "Single Point Plot")

# Additional implementation at 2025-06-19 22:42:45
