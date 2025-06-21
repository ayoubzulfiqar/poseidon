import math

def plot_ascii(func_str, x_min, x_max, y_min, y_max, width=80, height=25):
    grid = [[' ' for _ in range(width)] for _ in range(height)]

    if x_max <= x_min or y_max <= y_min:
        print("Error: x_max must be greater than x_min, and y_max must be greater than y_min.")
        return

    x_scale = (width - 1) / (x_max - x_min)
    y_scale = (height - 1) / (y_max - y_min)

    eval_env = {"x": 0, "math": math}

    for col in range(width):
        x = x_min + col / x_scale
        
        eval_env["x"] = x
        try:
            y = eval(func_str, {"__builtins__": None}, eval_env)
            
            row = int((y_max - y) * y_scale)

            if 0 <= row < height:
                grid[row][col] = '*'
        except (ValueError, TypeError, NameError, ZeroDivisionError):
            continue

    y_origin_row = int((y_max - 0) * y_scale)
    if 0 <= y_origin_row < height:
        for col in range(width):
            if grid[y_origin_row][col] == '*':
                grid[y_origin_row][col] = '#'
            elif grid[y_origin_row][col] == ' ':
                grid[y_origin_row][col] = '-'

    x_origin_col = int((0 - x_min) * x_scale)
    if 0 <= x_origin_col < width:
        for row in range(height):
            if grid[row][x_origin_col] == '*':
                grid[row][x_origin_col] = '#'
            elif grid[row][x_origin_col] == ' ':
                grid[row][x_origin_col] = '|'

    if 0 <= y_origin_row < height and 0 <= x_origin_col < width:
        if grid[y_origin_row][x_origin_col] == '*':
            grid[y_origin_row][x_origin_col] = '#'
        elif grid[y_origin_row][x_origin_col] == '-' or grid[y_origin_row][x_origin_col] == '|':
            grid[y_origin_row][x_origin_col] = '+'

    for row in grid:
        print("".join(row))

if __name__ == "__main__":
    print("Plotting y = x^2:")
    plot_ascii("x**2", -5, 5, -1, 25)
    print("\nPlotting y = sin(x):")
    plot_ascii("math.sin(x)", -2 * math.pi, 2 * math.pi, -1.5, 1.5)
    print("\nPlotting y = 1/x:")
    plot_ascii("1/x", -5, 5, -5, 5)
    print("\nPlotting y = e^x:")
    plot_ascii("math.exp(x)", -3, 3, 0, 20)
    print("\nPlotting y = ln(x):")
    plot_ascii("math.log(x)", 0.1, 10, -3, 3)
    print("\nPlotting y = |x|:")
    plot_ascii("abs(x)", -10, 10, 0, 10)
    print("\nPlotting y = x:")
    plot_ascii("x", -10, 10, -10, 10)

# Additional implementation at 2025-06-21 04:08:51
import math

class ASCIIPlotter:
    def __init__(self, width=80, height=25):
        self.width = width
        self.height = height
        self.canvas = []
        self.x_min = -10.0
        self.x_max = 10.0
        self.y_min = -10.0
        self.y_max = 10.0
        self.functions = []

    def _initialize_canvas(self):
        self.canvas = [[' ' for _ in range(self.width)] for _ in range(self.height)]

    def _map_coords(self, x, y):
        screen_x = int((x - self.x_min) / (self.x_max - self.x_min) * (self.width - 1))
        screen_y = int((self.y_max - y) / (self.y_max - self.y_min) * (self.height - 1))
        return screen_x, screen_y

    def _draw_axes(self):
        if self.y_min <= 0 <= self.y_max:
            y_axis_row = int((self.y_max - 0) / (self.y_max - self.y_min) * (self.height - 1))
            if 0 <= y_axis_row < self.height:
                for col in range(self.width):
                    if self.canvas[y_axis_row][col] == ' ':
                        self.canvas[y_axis_row][col] = '-'

        if self.x_min <= 0 <= self.x_max:
            x_axis_col = int((0 - self.x_min) / (self.x_max - self.x_min) * (self.width - 1))
            if 0 <= x_axis_col < self.width:
                for row in range(self.height):
                    if self.canvas[row][x_axis_col] == ' ':
                        self.canvas[row][x_axis_col] = '|'

        if self.x_min <= 0 <= self.x_max and self.y_min <= 0 <= self.y_max:
            x_origin_col = int((0 - self.x_min) / (self.x_max - self.x_min) * (self.width - 1))
            y_origin_row = int((self.y_max - 0) / (self.y_max - self.y_min) * (self.height - 1))
            if 0 <= x_origin_col < self.width and 0 <= y_origin_row < self.height:
                self.canvas[y_origin_row][x_origin_col] = '+'

    def add_function(self, func_str, plot_char='*'):
        self.functions.append((func_str, plot_char))

    def set_range(self, x_min, x_max, y_min, y_max):
        self.x_min = float(x_min)
        self.x_max = float(x_max)
        self.y_min = float(y_min)
        self.y_max = float(y_max)

    def plot(self):
        self._initialize_canvas()
        self._draw_axes()

        eval_context = {
            'x': 0.0,
            'math': math,
            'sin': math.sin, 'cos': math.cos, 'tan': math.tan,
            'asin': math.asin, 'acos': math.acos, 'atan': math.atan,
            'atan2': math.atan2, 'hypot': math.hypot,
            'degrees': math.degrees, 'radians': math.radians,
            'sinh': math.sinh, 'cosh': math.cosh, 'tanh': math.tanh,
            'asinh': math.asinh, 'acosh': math.acosh, 'atanh': math.atanh,
            'sqrt': math.sqrt, 'exp': math.exp, 'log': math.log, 'log10': math.log10, 'log2': math.log2,
            'ceil': math.ceil, 'floor': math.floor, 'trunc': math.trunc,
            'fabs': math.fabs, 'fmod': math.fmod,
            'pi': math.pi, 'e': math.e, 'tau': math.tau, 'inf': math.inf, 'nan': math.nan,
            'abs': abs,
            'pow': pow
        }

        for func_str, plot_char in self.functions:
            try:
                compiled_func = compile(func_str, '<string>', 'eval')
                x_step = (self.x_max - self.x_min) / (self.width - 1)
                current_x = self.x_min
                for _ in range(self.width):
                    eval_context['x'] = current_x
                    try:
                        y = eval(compiled_func, eval_context)
                        screen_x, screen_y = self._map_coords(current_x, y)

                        if 0 <= screen_x < self.width and 0 <= screen_y < self.height:
                            if self.canvas[screen_y][screen_x] in [' ', plot_char]:
                                self.canvas[screen_y][screen_x] = plot_char
                    except (TypeError, ValueError, NameError, ZeroDivisionError):
                        pass
                    current_x += x_step

            except SyntaxError:
                print(f"Error: Invalid function syntax for '{func_str}'. Skipping.")
                continue
            except Exception as e:
                print(f"An unexpected error occurred while processing function '{func_str}': {e}. Skipping.")
                continue

        for row in self.canvas:
            print("".join(row))

        print(f"X: {self.x_min:.2f} to {self.x_max:.2f}")
        print(f"Y: {self.y_min:.2f} to {self.y_max:.2f}")
        for i, (func_str, plot_char) in enumerate(self.functions):
            print(f"Function {i+1} ('{plot_char}'): {func_str}")


if __name__ == '__main__':
    plotter = ASCIIPlotter(width=80, height=25)

    while True:
        print("\n--- ASCII Function Plotter ---")
        print("1. Add function")
        print("2. Set plot range (Xmin Xmax Ymin Ymax)")
        print("3. Plot")
        print("4. Clear functions")
        print("5. Exit")

        choice = input("Enter your choice: ").strip()

        if choice == '1':
            func_str = input("Enter function (e.g., 'x**2', 'math.sin(x)' or 'sin(x)'): ").strip()
            if not func_str:
                print("Function cannot be empty.")
                continue
            plot_char = input("Enter plot character (default: '*'): ").strip()
            if not plot_char:
                plot_char = '*'
            elif len(plot_char) > 1:
                print("Plot character must be a single character. Using first character.")
                plot_char = plot_char[0]
            plotter.add_function(func_str, plot_char)
            print(f"Function '{func_str}' added with character '{plot_char}'.")

        elif choice == '2':
            try:
                range_input = input("Enter Xmin Xmax Ymin Ymax (e.g., -5 5 -10 10): ").split()
                if len(range_input) == 4:
                    xmin, xmax, ymin, ymax = map(float, range_input)
                    if xmin >= xmax or ymin >= ymax:
                        print("Error: Xmax must be greater than Xmin, and Ymax greater than Ymin.")
                    else:
                        plotter.set_range(xmin, xmax, ymin, ymax)
                        print(f"Plot range set to X:[{xmin:.2f},{xmax:.2f}], Y:[{ymin:.2f},{ymax:.2f}]")
                else:
                    print("Invalid input. Please provide 4 numbers.")
            except ValueError:
                print("Invalid input. Please enter numbers for the range.")

        elif choice == '3':
            if not plotter.functions:
                print("No functions added yet. Please add a function first.")
            else:
                plotter.plot()

        elif choice == '4':
            plotter.functions = []
            print("All functions cleared.")

        elif choice == '5':
            print("Exiting plotter.")
            break

        else:
            print("Invalid choice. Please try again.")