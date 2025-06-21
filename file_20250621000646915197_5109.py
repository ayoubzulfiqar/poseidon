import math

def plot_ascii_function(func_str, x_range, y_range, width=80, height=20):
    x_min, x_max = x_range
    y_min, y_max = y_range

    grid = [[' ' for _ in range(width)] for _ in range(height)]

    x_scale = (width - 1) / (x_max - x_min) if x_max != x_min else 0
    y_scale = (height - 1) / (y_max - y_min) if y_max != y_min else 0

    origin_x_char = -1
    if x_min <= 0 <= x_max:
        if x_scale != 0:
            origin_x_char = int((0 - x_min) * x_scale)
        else:
            origin_x_char = 0

    origin_y_char = -1
    if y_min <= 0 <= y_max:
        if y_scale != 0:
            origin_y_char = int((0 - y_min) * y_scale)
        else:
            origin_y_char = 0
        origin_y_char = height - 1 - origin_y_char

    if origin_y_char != -1:
        for x_idx in range(width):
            grid[origin_y_char][x_idx] = '-'
        if origin_x_char != -1:
            grid[origin_y_char][origin_x_char] = '+'

    if origin_x_char != -1:
        for y_idx in range(height):
            grid[y_idx][origin_x_char] = '|'
        if origin_y_char != -1:
            grid[origin_y_char][origin_x_char] = '+'

    scope = {'x': 0, 'math': math, 'abs': abs}

    for x_char in range(width):
        x = x_min + (x_char / x_scale) if x_scale != 0 else x_min
        scope['x'] = x
        try:
            y = eval(func_str, {"__builtins__": None}, scope)
            if y_min <= y <= y_max:
                y_char = int((y - y_min) * y_scale) if y_scale != 0 else 0
                y_char_screen = height - 1 - y_char

                if 0 <= y_char_screen < height:
                    if grid[y_char_screen][x_char] in ['-', '|', '+']:
                        grid[y_char_screen][x_char] = 'o'
                    else:
                        grid[y_char_screen][x_char] = '*'
        except (NameError, TypeError, SyntaxError, ZeroDivisionError, ValueError):
            pass

    for row in grid:
        print(''.join(row))

if __name__ == '__main__':
    plot_ascii_function("x**2", (-5, 5), (0, 25), 80, 20)
    print("\n" + "="*80 + "\n")
    plot_ascii_function("math.sin(x)", (-2*math.pi, 2*math.pi), (-1.5, 1.5), 80, 20)

# Additional implementation at 2025-06-21 00:07:43
import math

PLOT_WIDTH = 80
PLOT_HEIGHT = 25
PLOT_CHAR = '*'
AXIS_CHAR_H = '-'
AXIS_CHAR_V = '|'
AXIS_INTERSECT_CHAR = '+'
GRID_CHAR = ' '

def plot_function(func_str, x_min, x_max):
    grid = [[GRID_CHAR for _ in range(PLOT_WIDTH)] for _ in range(PLOT_HEIGHT)]

    sample_points = 1000
    y_values = []
    
    eval_scope = {"x": 0.0}
    for name in dir(math):
        if not name.startswith('_'):
            eval_scope[name] = getattr(math, name)

    try:
        for i in range(sample_points):
            x_sample = x_min + (x_max - x_min) * i / (sample_points - 1)
            eval_scope["x"] = x_sample
            y_val = eval(func_str, {"__builtins__": None}, eval_scope)
            if isinstance(y_val, (int, float)):
                y_values.append(y_val)
    except Exception as e:
        print(f"Error evaluating function for y-range estimation: {e}")
        print("Please ensure your function uses 'x' as the variable and math functions (e.g., math.sin, math.log).")
        return

    if not y_values:
        print("No valid y-values could be generated. Cannot plot.")
        return

    y_min = min(y_values)
    y_max = max(y_values)

    if abs(y_max - y_min) < 1e-6:
        y_min -= 0.5
        y_max += 0.5

    x_scale = (PLOT_WIDTH - 1) / (x_max - x_min)
    y_scale = (PLOT_HEIGHT - 1) / (y_max - y_min)

    x_grid_origin = int(-x_min * x_scale)
    y_grid_origin = int(PLOT_HEIGHT - 1 - (-y_min * y_scale))

    if 0 <= y_grid_origin < PLOT_HEIGHT:
        for col in range(PLOT_WIDTH):
            grid[y_grid_origin][col] = AXIS_CHAR_H
    
    if 0 <= x_grid_origin < PLOT_WIDTH:
        for row in range(PLOT_HEIGHT):
            grid[row][x_grid_origin] = AXIS_CHAR_V
    
    if 0 <= y_grid_origin < PLOT_HEIGHT and 0 <= x_grid_origin < PLOT_WIDTH:
        grid[y_grid_origin][x_grid_origin] = AXIS_INTERSECT_CHAR

    for px in range(PLOT_WIDTH):
        x = x_min + px / x_scale
        
        eval_scope["x"] = x
        try:
            y = eval(func_str, {"__builtins__": None}, eval_scope)
        except Exception:
            continue

        if not isinstance(y, (int, float)):
            continue

        py = int(PLOT_HEIGHT - 1 - (y - y_min) * y_scale)

        if 0 <= py < PLOT_HEIGHT:
            grid[py][px] = PLOT_CHAR

    print(f"Plotting: {func_str}")
    print(f"X-range: [{x_min:.2f}, {x_max:.2f}]")
    print(f"Y-range: [{y_min:.2f}, {y_max:.2f}]")

    print(f"{y_max:.2f}".ljust(5) + "┌" + "─" * PLOT_WIDTH + "┐")
    
    for r_idx, row in enumerate(grid):
        label_str = ""
        if r_idx == PLOT_HEIGHT // 2:
            label_str = f"{y_min + (y_max - y_min) / 2:.2f}"
        
        print(label_str.ljust(5) + "│" + "".join(row) + "│")

    print(f"{y_min:.2f}".ljust(5) + "└" + "─" * PLOT_WIDTH + "┘")
    
    x_min_label = f"{x_min:.2f}"
    x_max_label = f"{x_max:.2f}"
    
    padding_left_for_x_labels = 5 + 1
    
    space_between_x_labels = PLOT_WIDTH - len(x_min_label) - len(x_max_label)
    if space_between_x_labels < 0:
        print(" " * padding_left_for_x_labels + x_min_label)
    else:
        print(" " * padding_left_for_x_labels + x_min_label + " " * space_between_x_labels + x_max_label)


def main():
    print("ASCII Function Plotter")
    print("Enter a mathematical function using 'x' as the variable.")
    print("Available functions from 'math' module: sin, cos, tan, sqrt, exp, log, log10, pi, etc.")
    print("Example: math.sin(x), x**2, 1/x, math.exp(-x**2/2)")

    while True:
        func_str = input("Enter function (e.g., math.sin(x)): ")
        if not func_str:
            print("Function cannot be empty. Exiting.")
            break

        try:
            x_range_str = input("Enter x-range (e.g., -10,10): ")
            x_min_str, x_max_str = x_range_str.split(',')
            x_min = float(x_min_str.strip())
            x_max = float(x_max_str.strip())
            if x_min >= x_max:
                print("x_min must be less than x_max. Please try again.")
                continue
        except ValueError:
            print("Invalid x-range format. Please use 'min,max'.")
            continue
        except Exception as e:
            print(f"An error occurred while parsing x-range: {e}")
            continue

        plot_function(func_str, x_min, x_max)
        
        another = input("Plot another function? (yes/no): ").lower()
        if another != 'yes':
            break

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 00:09:03
import math

def plot_function_ascii():
    PLOT_WIDTH = 80
    PLOT_HEIGHT = 25

    CHAR_EMPTY = ' '
    CHAR_POINT = '*'
    CHAR_X_AXIS = '-'
    CHAR_Y_AXIS = '|'
    CHAR_ORIGIN = '+'

    while True:
        function_str = input("Enter function (e.g., 'math.sin(x)', 'x**2', '1/x'): ")
        try:
            eval(function_str, {"x": 0.0, "math": math})
            break
        except (NameError, SyntaxError, TypeError) as e:
            print(f"Invalid function expression: {e}. Please try again.")

    while True:
        try:
            x_min = float(input("Enter x_min: "))
            x_max = float(input("Enter x_max: "))
            if x_min >= x_max:
                print("x_max must be greater than x_min. Please try again.")
            else:
                break
        except ValueError:
            print("Invalid input for x_min or x_max. Please enter a number.")

    y_values = []
    x_step_for_y_calc = (x_max - x_min) / (PLOT_WIDTH * 4)
    if x_step_for_y_calc == 0:
        x_step_for_y_calc = 1.0

    current_x = x_min
    while current_x <= x_max:
        try:
            y = eval(function_str, {"x": current_x, "math": math})
            if math.isfinite(y):
                y_values.append(y)
        except (NameError, SyntaxError, TypeError, ZeroDivisionError):
            pass
        current_x += x_step_for_y_calc

    if not y_values:
        print("No valid y-values could be calculated for the given x-range. Cannot plot.")
        return

    y_min_plot = min(y_values)
    y_max_plot = max(y_values)

    y_range_buffer = (y_max_plot - y_min_plot) * 0.1
    if y_range_buffer == 0:
        y_range_buffer = 1.0
    y_min_plot -= y_range_buffer
    y_max_plot += y_range_buffer

    if y_min_plot >= y_max_plot:
        y_min_plot = -1.0
        y_max_plot = 1.0

    grid = [[CHAR_EMPTY for _ in range(PLOT_WIDTH)] for _ in range(PLOT_HEIGHT)]

    x_axis_screen_pos = -1
    if x_min <= 0 <= x_max:
        x_axis_screen_pos = int((0 - x_min) / (x_max - x_min) * (PLOT_WIDTH - 1))

    y_axis_screen_pos = -1
    if y_min_plot <= 0 <= y_max_plot:
        y_axis_screen_pos = PLOT_HEIGHT - 1 - int((0 - y_min_plot) / (y_max_plot - y_min_plot) * (PLOT_HEIGHT - 1))

    if 0 <= y_axis_screen_pos < PLOT_HEIGHT:
        for i in range(PLOT_WIDTH):
            grid[y_axis_screen_pos][i] = CHAR_X_AXIS

    if 0 <= x_axis_screen_pos < PLOT_WIDTH:
        for i in range(PLOT_HEIGHT):
            grid[i][x_axis_screen_pos] = CHAR_Y_AXIS

    if 0 <= x_axis_screen_pos < PLOT_WIDTH and 0 <= y_axis_screen_pos < PLOT_HEIGHT:
        grid[y_axis_screen_pos][x_axis_screen_pos] = CHAR_ORIGIN

    for px in range(PLOT_WIDTH):
        x = x_min + (px / (PLOT_WIDTH - 1)) * (x_max - x_min)

        try:
            y = eval(function_str, {"x": x, "math": math})

            if math.isfinite(y):
                py = PLOT_HEIGHT - 1 - int((y - y_min_plot) / (y_max_plot - y_min_plot) * (PLOT_HEIGHT - 1))

                if 0 <= py < PLOT_HEIGHT:
                    grid[py][px] = CHAR_POINT
        except (NameError, SyntaxError, TypeError, ZeroDivisionError):
            pass

    for row in grid:
        print("".join(row))

if __name__ == "__main__":
    plot_function_ascii()

# Additional implementation at 2025-06-21 00:09:39
import math

def plot_ascii_function():
    while True:
        try:
            func_str = input("Enter function (e.g., x**2, sin(x)): ")
            x_min_str = input("Enter x_min (e.g., -5): ")
            x_max_str = input("Enter x_max (e.g., 5): ")
            width_str = input("Enter plot width (characters, e.g., 80): ")
            height_str = input("Enter plot height (characters, e.g., 20): ")

            x_min = float(x_min_str)
            x_max = float(x_max_str)
            width = int(width_str)
            height = int(height_str)

            if x_min >= x_max:
                print("Error: x_min must be less than x_max.")
                continue
            if width < 10 or height < 5:
                print("Error: Plot width must be at least 10, height at least 5.")
                continue

            scope = {
                'x': 0.0,
                'sin': math.sin,
                'cos': math.cos,
                'tan': math.tan,
                'exp': math.exp,
                'log': math.log,
                'log10': math.log10,
                'sqrt': math.sqrt,
                'abs': abs,
                'ceil': math.ceil,
                'floor': math.floor,
                'trunc': math.trunc,
                'radians': math.radians,
                'degrees': math.degrees,
                'pi': math.pi,
                'e': math.e
            }

            try:
                scope['x'] = x_min
                eval(func_str, {"__builtins__": None}, scope)
            except (SyntaxError, NameError, TypeError) as e:
                print(f"Error in function expression: {e}")
                continue
            except Exception as e:
                print(f"An unexpected error occurred during function test: {e}")
                continue

            break
        except ValueError:
            print("Invalid input. Please enter numbers for range and dimensions.")
        except Exception as e:
            print(f"An unexpected error occurred during input: {e}")

    x_step = (x_max - x_min) / (width - 1)
    y_values = []
    
    y_min_val = float('inf')
    y_max_val = float('-inf')

    for i in range(width):
        x = x_min + i * x_step
        scope['x'] = x
        try:
            y = eval(func_str, {"__builtins__": None}, scope)
            if isinstance(y, (int, float)):
                y_values.append(y)
                if y < y_min_val:
                    y_min_val = y
                if y > y_max_val:
                    y_max_val = y
            else:
                y_values.append(float('nan'))
        except (ZeroDivisionError, ValueError, OverflowError):
            y_values.append(float('nan'))
        except Exception:
            y_values.append(float('nan'))

    if y_min_val == float('inf') or y_max_val == float('-inf') or y_min_val == y_max_val:
        if y_min_val == float('inf'):
            y_min_val = -1.0
            y_max_val = 1.0
        elif y_min_val == y_max_val:
            if y_min_val == 0:
                y_min_val = -1.0
                y_max_val = 1.0
            else:
                buffer = abs(y_min_val * 0.1) + 0.1
                y_min_val -= buffer
                y_max_val += buffer
        
    y_range = y_max_val - y_min_val
    if y_range == 0:
        y_range = 1.0

    plot_grid = [[' ' for _ in range(width)] for _ in range(height)]

    x_axis_row = -1
    y_axis_col = -1

    if y_min_val <= 0 <= y_max_val:
        x_axis_row = height - 1 - int((0 - y_min_val) * (height - 1) / y_range)
    
    if x_min <= 0 <= x_max:
        y_axis_col = int((0 - x_min) / x_step)

    for col in range(width):
        if x_axis_row != -1:
            plot_grid[x_axis_row][col] = '-'
    for row in range(height):
        if y_axis_col != -1:
            plot_grid[row][y_axis_col] = '|'

    if x_axis_row != -1 and y_axis_col != -1:
        plot_grid[x_axis_row][y_axis_col] = '+'

    for i, y_val in enumerate(y_values):
        if not math.isnan(y_val):
            row = height - 1 - int((y_val - y_min_val) * (height - 1) / y_range)

            if 0 <= row < height:
                if plot_grid[row][i] == ' ':
                    plot_grid[row][i] = '*'
                elif plot_grid[row][i] == '-' and row == x_axis_row:
                    plot_grid[row][i] = '*'
                elif plot_grid[row][i] == '|' and i == y_axis_col:
                    plot_grid[row][i] = '*'
                elif plot_grid[row][i] == '+':
                    plot_grid[row][i] = '*'

    for row in plot_grid:
        print("".join(row))

    print(f"X-range: [{x_min:.2f}, {x_max:.2f}]")
    print(f"Y-range: [{y_min_val:.2f}, {y_max_val:.2f}]")

plot_ascii_function()