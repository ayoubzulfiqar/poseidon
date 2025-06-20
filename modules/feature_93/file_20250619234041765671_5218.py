import sys

def _convert_length(value, from_unit, to_unit):
    meters = 0.0
    
    # Convert to meters
    if from_unit == 'm':
        meters = value
    elif from_unit == 'ft':
        meters = value * 0.3048
    elif from_unit == 'in':
        meters = value * 0.0254
    elif from_unit == 'cm':
        meters = value * 0.01
    elif from_unit == 'km':
        meters = value * 1000
    elif from_unit == 'mi':
        meters = value * 1609.34
    elif from_unit == 'yd':
        meters = value * 0.9144
    else:
        return None # Invalid from_unit

    # Convert from meters to target unit
    if to_unit == 'm':
        return meters
    elif to_unit == 'ft':
        return meters / 0.3048
    elif to_unit == 'in':
        return meters / 0.0254
    elif to_unit == 'cm':
        return meters / 0.01
    elif to_unit == 'km':
        return meters / 1000
    elif to_unit == 'mi':
        return meters / 1609.34
    elif to_unit == 'yd':
        return meters / 0.9144
    else:
        return None # Invalid to_unit

def _convert_weight(value, from_unit, to_unit):
    kilograms = 0.0

    # Convert to kilograms
    if from_unit == 'kg':
        kilograms = value
    elif from_unit == 'lb':
        kilograms = value * 0.453592
    elif from_unit == 'g':
        kilograms = value * 0.001
    elif from_unit == 'oz':
        kilograms = value * 0.0283495
    elif from_unit == 't':
        kilograms = value * 1000
    else:
        return None # Invalid from_unit

    # Convert from kilograms to target unit
    if to_unit == 'kg':
        return kilograms
    elif to_unit == 'lb':
        return kilograms / 0.453592
    elif to_unit == 'g':
        return kilograms / 0.001
    elif to_unit == 'oz':
        return kilograms / 0.0283495
    elif to_unit == 't':
        return kilograms / 1000
    else:
        return None # Invalid to_unit

def _convert_temperature(value, from_unit, to_unit):
    celsius = 0.0

    # Convert to Celsius
    if from_unit == 'C':
        celsius = value
    elif from_unit == 'F':
        celsius = (value - 32) * 5/9
    elif from_unit == 'K':
        celsius = value - 273.15
    else:
        return None # Invalid from_unit

    # Convert from Celsius to target unit
    if to_unit == 'C':
        return celsius
    elif to_unit == 'F':
        return (celsius * 9/5) + 32
    elif to_unit == 'K':
        return celsius + 273.15
    else:
        return None # Invalid to_unit

def _get_float_input(prompt):
    while True:
        try:
            return float(input(prompt))
        except ValueError:
            print("Invalid number. Please enter a numeric value.")

def _get_unit_input(prompt, valid_units):
    while True:
        unit = input(prompt).strip().lower()
        if unit in valid_units:
            return unit
        else:
            print(f"Invalid unit. Please choose from: {', '.join(valid_units)}")

def _run_length_converter():
    print("\n--- Length Converter ---")
    value = _get_float_input("Enter the value to convert: ")
    from_unit = _get_unit_input("From unit (m, ft, in, cm, km, mi, yd): ", ['m', 'ft', 'in', 'cm', 'km', 'mi', 'yd'])
    to_unit = _get_unit_input("To unit (m, ft, in, cm, km, mi, yd): ", ['m', 'ft', 'in', 'cm', 'km', 'mi', 'yd'])

    result = _convert_length(value, from_unit, to_unit)
    if result is not None:
        print(f"{value} {from_unit} is {result:.4f} {to_unit}")
    else:
        print("An error occurred during conversion. Check units.")

def _run_weight_converter():
    print("\n--- Weight Converter ---")
    value = _get_float_input("Enter the value to convert: ")
    from_unit = _get_unit_input("From unit (kg, lb, g, oz, t): ", ['kg', 'lb', 'g', 'oz', 't'])
    to_unit = _get_unit_input("To unit (kg, lb, g, oz, t): ", ['kg', 'lb', 'g', 'oz', 't'])

    result = _convert_weight(value, from_unit, to_unit)
    if result is not None:
        print(f"{value} {from_unit} is {result:.4f} {to_unit}")
    else:
        print("An error occurred during conversion. Check units.")

def _run_temperature_converter():
    print("\n--- Temperature Converter ---")
    value = _get_float_input("Enter the value to convert: ")
    from_unit = _get_unit_input("From unit (C, F, K): ", ['c', 'f', 'k'])
    to_unit = _get_unit_input("To unit (C, F, K): ", ['c', 'f', 'k'])

    result = _convert_temperature(value, from_unit, to_unit)
    if result is not None:
        print(f"{value} {from_unit} is {result:.4f} {to_unit}")
    else:
        print("An error occurred during conversion. Check units.")

def main():
    print("Welcome to the CLI Unit Converter!")
    while True:
        print("\nSelect a conversion type:")
        print("1. Length")
        print("2. Weight")
        print("3. Temperature")
        print("4. Exit")

        choice = input("Enter your choice (1-4): ").strip()

        if choice == '1':
            _run_length_converter()
        elif choice == '2':
            _run_weight_converter()
        elif choice == '3':
            _run_temperature_converter()
        elif choice == '4':
            print("Exiting converter. Goodbye!")
            sys.exit()
        else:
            print("Invalid choice. Please enter a number between 1 and 4.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-19 23:41:19
LENGTH_UNITS = {
    "m": 1.0,
    "meter": 1.0,
    "meters": 1.0,
    "km": 1000.0,
    "kilometer": 1000.0,
    "kilometers": 1000.0,
    "cm": 0.01,
    "centimeter": 0.01,
    "centimeters": 0.01,
    "mm": 0.001,
    "millimeter": 0.001,
    "millimeters": 0.001,
    "mi": 1609.34,
    "mile": 1609.34,
    "miles": 1609.34,
    "ft": 0.3048,
    "foot": 0.3048,
    "feet": 0.3048,
    "in": 0.0254,
    "inch": 0.0254,
    "inches": 0.0254
}

WEIGHT_UNITS = {
    "kg": 1.0,
    "kilogram": 1.0,
    "kilograms": 1.0,
    "g": 0.001,
    "gram": 0.001,
    "grams": 0.001,
    "mg": 0.000001,
    "milligram": 0.000001,
    "milligrams": 0.000001,
    "lb": 0.453592,
    "pound": 0.453592,
    "pounds": 0.453592,
    "oz": 0.0283495,
    "ounce": 0.0283495,
    "ounces": 0.0283495,
    "ton": 907.185,
    "metric ton": 1000.0,
    "tonne": 1000.0
}

def convert_temperature(value, from_unit, to_unit):
    from_unit = from_unit.lower()
    to_unit = to_unit.lower()

    if from_unit not in ["celsius", "c", "fahrenheit", "f", "kelvin", "k"] or \
       to_unit not in ["celsius", "c", "fahrenheit", "f", "kelvin", "k"]:
        return None

    celsius_val = 0.0
    if from_unit in ["celsius", "c"]:
        celsius_val = value
    elif from_unit in ["fahrenheit", "f"]:
        celsius_val = (value - 32) * 5/9
    elif from_unit in ["kelvin", "k"]:
        celsius_val = value - 273.15

    if to_unit in ["celsius", "c"]:
        return celsius_val
    elif to_unit in ["fahrenheit", "f"]:
        return (celsius_val * 9/5) + 32
    elif to_unit in ["kelvin", "k"]:
        return celsius_val + 273.15
    return None

def convert_units(value, from_unit, to_unit, unit_map):
    from_unit = from_unit.lower()
    to_unit = to_unit.lower()

    if from_unit not in unit_map or to_unit not in unit_map:
        return None

    value_in_base = value * unit_map[from_unit]
    result = value_in_base / unit_map[to_unit]
    return result

def get_valid_input(prompt, valid_options):
    while True:
        user_input = input(prompt).strip().lower()
        if user_input in valid_options:
            return user_input
        print(f"Invalid input. Please choose from: {', '.join(sorted(list(valid_options)))}")

def get_float_input(prompt):
    while True:
        try:
            value = float(input(prompt).strip())
            return value
        except ValueError:
            print("Invalid number. Please enter a numerical value.")

def display_units(unit_map):
    print("Available units:")
    for unit in sorted(unit_map.keys()):
        print(f"- {unit}")

def main():
    print("Welcome to the Unit Converter CLI!")
    print("Available categories: length, weight, temperature")
    print("Type 'quit' to exit at any time.")

    while True:
        category = input("\nEnter category (length, weight, temperature, quit): ").strip().lower()

        if category == "quit":
            print("Exiting converter. Goodbye!")
            break
        elif category not in ["length", "weight", "temperature"]:
            print("Invalid category. Please choose from 'length', 'weight', 'temperature'.")
            continue

        if category == "length":
            print("\n--- Length Conversion ---")
            display_units(LENGTH_UNITS)
            from_unit = get_valid_input("Convert from unit: ", LENGTH_UNITS.keys())
            to_unit = get_valid_input("Convert to unit: ", LENGTH_UNITS.keys())
            value = get_float_input(f"Enter value in {from_unit}: ")
            
            result = convert_units(value, from_unit, to_unit, LENGTH_UNITS)
            if result is not None:
                print(f"{value} {from_unit} is {result:.4f} {to_unit}")
            else:
                print("Error during conversion. Please check units.")

        elif category == "weight":
            print("\n--- Weight Conversion ---")
            display_units(WEIGHT_UNITS)
            from_unit = get_valid_input("Convert from unit: ", WEIGHT_UNITS.keys())
            to_unit = get_valid_input("Convert to unit: ", WEIGHT_UNITS.keys())
            value = get_float_input(f"Enter value in {from_unit}: ")

            result = convert_units(value, from_unit, to_unit, WEIGHT_UNITS)
            if result is not None:
                print(f"{value} {from_unit} is {result:.4f} {to_unit}")
            else:
                print("Error during conversion. Please check units.")

        elif category == "temperature":
            print("\n--- Temperature Conversion ---")
            temp_units = {"celsius", "c", "fahrenheit", "f", "kelvin", "k"}
            print("Available units: Celsius (C), Fahrenheit (F), Kelvin (K)")
            from_unit = get_valid_input("Convert from unit (C, F, K): ", temp_units)
            to_unit = get_valid_input("Convert to unit (C, F, K): ", temp_units)
            value = get_float_input(f"Enter value in {from_unit}: ")

            result = convert_temperature(value, from_unit, to_unit)
            if result is not None:
                print(f"{value} {from_unit} is {result:.4f} {to_unit}")
            else:
                print("Error during conversion. Please check units.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-19 23:41:51
import sys

def convert_length(value, from_unit, to_unit):
    units = {
        'm': 1.0,
        'km': 1000.0,
        'cm': 0.01,
        'mm': 0.001,
        'mi': 1609.34,
        'yd': 0.9144,
        'ft': 0.3048,
        'in': 0.0254
    }
    if from_unit not in units or to_unit not in units:
        raise ValueError("Invalid length unit. Supported: m, km, cm, mm, mi, yd, ft, in")
    
    value_in_meters = value * units[from_unit]
    converted_value = value_in_meters / units[to_unit]
    return converted_value

def convert_weight(value, from_unit, to_unit):
    units = {
        'kg': 1.0,
        'g': 0.001,
        'mg': 0.000001,
        't': 1000.0, # metric ton
        'lb': 0.453592,
        'oz': 0.0283495,
        'st': 6.35029 # stone
    }
    if from_unit not in units or to_unit not in units:
        raise ValueError("Invalid weight unit. Supported: kg, g, mg, t, lb, oz, st")
    
    value_in_kg = value * units[from_unit]
    converted_value = value_in_kg / units[to_unit]
    return converted_value

def convert_temperature(value, from_unit, to_unit):
    if from_unit == to_unit:
        return value

    if from_unit == 'c':
        celsius = value
    elif from_unit == 'f':
        celsius = (value - 32) * 5/9
    elif from_unit == 'k':
        celsius = value - 273.15
    else:
        raise ValueError("Invalid temperature unit. Supported: C, F, K")

    if to_unit == 'c':
        return celsius
    elif to_unit == 'f':
        return (celsius * 9/5) + 32
    elif to_unit == 'k':
        return celsius + 273.15
    else:
        raise ValueError("Invalid temperature unit. Supported: C, F, K")

def run_converter():
    while True:
        print("\n--- Unit Converter ---")
        print("Select a category:")
        print("1. Length (m, km, cm, mm, mi, yd, ft, in)")
        print("2. Weight (kg, g, mg, t, lb, oz, st)")
        print("3. Temperature (C, F, K)")
        print("4. Exit")

        category_choice = input("Enter choice (1-4): ").strip()

        if category_choice == '4':
            print("Exiting converter. Goodbye!")
            break

        if category_choice not in ['1', '2', '3']:
            print("Invalid category choice. Please try again.")
            continue

        try:
            value_str = input("Enter the value to convert: ").strip()
            value = float(value_str)
            from_unit = input("Enter the unit to convert FROM (e.g., m, kg, C): ").strip().lower()
            to_unit = input("Enter the unit to convert TO (e.g., km, g, F): ").strip().lower()

            result = None
            if category_choice == '1':
                result = convert_length(value, from_unit, to_unit)
            elif category_choice == '2':
                result = convert_weight(value, from_unit, to_unit)
            elif category_choice == '3':
                result = convert_temperature(value, from_unit, to_unit)
            
            if result is not None:
                print(f"Result: {value} {from_unit} is {result:.4f} {to_unit}")

        except ValueError as e:
            print(f"Error: {e}. Please check your input and units.")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    run_converter()

# Additional implementation at 2025-06-19 23:42:16
import sys

class UnitConverter:
    def __init__(self):
        self.length_units = {
            'm': 1.0,
            'km': 1000.0,
            'cm': 0.01,
            'mm': 0.001,
            'mi': 1609.34,
            'yd': 0.9144,
            'ft': 0.3048,
            'in': 0.0254
        }
        self.weight_units = {
            'kg': 1.0,
            'g': 0.001,
            'mg': 0.000001,
            't': 1000.0,
            'lb': 0.453592,
            'oz': 0.0283495,
            'st': 6.35029
        }
        self.temperature_units = ['C', 'F', 'K']

    def _get_float_input(self, prompt):
        while True:
            try:
                value = float(input(prompt))
                return value
            except ValueError:
                print("Invalid input. Please enter a numeric value.")

    def _get_unit_input(self, prompt, valid_units_collection):
        if isinstance(valid_units_collection, dict):
            valid_units_set = {u.lower() for u in valid_units_collection.keys()}
            display_units = ', '.join(sorted(valid_units_collection.keys()))
        else: # Assume list for temperature
            valid_units_set = {u.lower() for u in valid_units_collection}
            display_units = ', '.join(sorted(valid_units_collection))

        while True:
            unit_input = input(prompt).strip().lower()
            if unit_input in valid_units_set:
                return unit_input
            else:
                print(f"Invalid unit. Please choose from: {display_units}")

    def convert_length(self, value, from_unit, to_unit):
        # Ensure units are lowercase for internal consistency
        lower_length_units = {k.lower(): v for k, v in self.length_units.items()}

        if from_unit not in lower_length_units or to_unit not in lower_length_units:
            raise ValueError("Invalid length unit provided.")

        value_in_meters = value * lower_length_units[from_unit]
        converted_value = value_in_meters / lower_length_units[to_unit]
        return converted_value

    def convert_weight(self, value, from_unit, to_unit):
        # Ensure units are lowercase for internal consistency
        lower_weight_units = {k.lower(): v for k, v in self.weight_units.items()}

        if from_unit not in lower_weight_units or to_unit not in lower_weight_units:
            raise ValueError("Invalid weight unit provided.")

        value_in_kg = value * lower_weight_units[from_unit]
        converted_value = value_in_kg / lower_weight_units[to_unit]
        return converted_value

    def convert_temperature(self, value, from_unit, to_unit):
        # from_unit and to_unit are already lowercase from _get_unit_input
        # Internal logic uses lowercase unit identifiers

        if from_unit == to_unit:
            return value

        # Convert to Celsius (lowercase 'c') first
        if from_unit == 'f':
            value_in_celsius = (value - 32) * 5/9
        elif from_unit == 'k':
            value_in_celsius = value - 273.15
        else: # from_unit == 'c'
            value_in_celsius = value

        # Convert from Celsius (lowercase 'c') to target unit
        if to_unit == 'f':
            converted_value = value_in_celsius * 9/5 + 32
        elif to_unit == 'k':
            converted_value = value_in_celsius + 273.15
        else: # to_unit == 'c'
            converted_value = value_in_celsius
        return converted_value

    def run(self):
        print("Welcome to the CLI Unit Converter!")
        while True:
            print("\nChoose a category:")
            print("1. Length")
            print("2. Weight")
            print("3. Temperature")
            print("4. Exit")

            choice = input("Enter your choice (1-4): ").strip()

            if choice == '1':
                print("\n--- Length Conversion ---")
                value = self._get_float_input("Enter the value to convert: ")
                from_unit = self._get_unit_input("Convert from unit: ", self.length_units)
                to_unit = self._get_unit_input("Convert to unit: ", self.length_units)
                try:
                    result = self.convert_length(value, from_unit, to_unit)
                    print(f"{value} {from_unit} is {result:.4f} {to_unit}")
                except ValueError as e:
                    print(f"Error: {e}")

            elif choice == '2':
                print("\n--- Weight Conversion ---")
                value = self._get_float_input("Enter the value to convert: ")
                from_unit = self._get_unit_input("Convert from unit: ", self.weight_units)
                to_unit = self._get_unit_input("Convert to unit: ", self.weight_units)
                try:
                    result = self.convert_weight(value, from_unit, to_unit)
                    print(f"{value} {from_unit} is {result:.4f} {to_unit}")
                except ValueError as e:
                    print(f"Error: {e}")

            elif choice == '3':
                print("\n--- Temperature Conversion ---")
                value = self._get_float_input("Enter the value to convert: ")
                from_unit = self._get_unit_input("Convert from unit (C, F, K): ", self.temperature_units)
                to_unit = self._get_unit_input("Convert to unit (C, F, K): ", self.temperature_units)
                try:
                    result = self.convert_temperature(value, from_unit, to_unit)
                    # Display temperature units in their common uppercase form
                    print(f"{value} {from_unit.upper()} is {result:.2f} {to_unit.upper()}")
                except ValueError as e:
                    print(f"Error: {e}")

            elif choice == '4':
                print("Exiting converter. Goodbye!")
                sys.exit()
            else:
                print("Invalid choice. Please enter a number between 1 and 4.")

if __name__ == "__main__":
    converter = UnitConverter()
    converter.run()