import sys

def _convert_temperature(value, unit_from, unit_to):
    # Convert to Kelvin first
    if unit_from == 'celsius':
        kelvin = value + 273.15
    elif unit_from == 'fahrenheit':
        kelvin = (value - 32) * 5/9 + 273.15
    elif unit_from == 'kelvin':
        kelvin = value
    else:
        raise ValueError("Invalid 'from' temperature unit.")

    # Convert from Kelvin to target unit
    if unit_to == 'celsius':
        return kelvin - 273.15
    elif unit_to == 'fahrenheit':
        return (kelvin - 273.15) * 9/5 + 32
    elif unit_to == 'kelvin':
        return kelvin
    else:
        raise ValueError("Invalid 'to' temperature unit.")

def _convert_length(value, unit_from, unit_to):
    # Base unit: meter
    factors_to_meter = {
        'meter': 1.0,
        'centimeter': 0.01,
        'kilometer': 1000.0,
        'inch': 0.0254,
        'foot': 0.3048,
        'yard': 0.9144,
        'mile': 1609.34
    }

    if unit_from not in factors_to_meter:
        raise ValueError("Invalid 'from' length unit.")
    if unit_to not in factors_to_meter:
        raise ValueError("Invalid 'to' length unit.")

    # Convert to meters
    meters = value * factors_to_meter[unit_from]
    # Convert from meters to target unit
    return meters / factors_to_meter[unit_to]

def _convert_weight(value, unit_from, unit_to):
    # Base unit: kilogram
    factors_to_kg = {
        'kilogram': 1.0,
        'gram': 0.001,
        'pound': 0.453592,
        'ounce': 0.0283495,
        'metric ton': 1000.0,
        'short ton': 907.185, # US ton
        'long ton': 1016.05  # UK ton
    }

    if unit_from not in factors_to_kg:
        raise ValueError("Invalid 'from' weight unit.")
    if unit_to not in factors_to_kg:
        raise ValueError("Invalid 'to' weight unit.")

    # Convert to kilograms
    kilograms = value * factors_to_kg[unit_from]
    # Convert from kilograms to target unit
    return kilograms / factors_to_kg[unit_to]

def run_converter():
    while True:
        print("\n--- Unit Converter ---")
        print("1. Temperature")
        print("2. Length")
        print("3. Weight")
        print("4. Exit")

        choice = input("Enter category number: ").strip()

        if choice == '4':
            print("Exiting converter. Goodbye!")
            sys.exit()

        try:
            value = float(input("Enter value to convert: "))
        except ValueError:
            print("Invalid value. Please enter a number.")
            continue

        unit_from = input("Enter unit to convert FROM (e.g., celsius, meter, kilogram): ").strip().lower()
        unit_to = input("Enter unit to convert TO (e.g., fahrenheit, foot, pound): ").strip().lower()

        result = None
        error_message = None

        try:
            if choice == '1':
                print("Available temperature units: celsius, fahrenheit, kelvin")
                result = _convert_temperature(value, unit_from, unit_to)
            elif choice == '2':
                print("Available length units: meter, centimeter, kilometer, inch, foot, yard, mile")
                result = _convert_length(value, unit_from, unit_to)
            elif choice == '3':
                print("Available weight units: kilogram, gram, pound, ounce, metric ton, short ton, long ton")
                result = _convert_weight(value, unit_from, unit_to)
            else:
                error_message = "Invalid category choice. Please enter 1, 2, 3, or 4."
        except ValueError as e:
            error_message = str(e)
        except Exception as e:
            error_message = f"An unexpected error occurred: {e}"

        if result is not None:
            print(f"Result: {value} {unit_from} is {result:.4f} {unit_to}")
        elif error_message:
            print(f"Error: {error_message}")

if __name__ == "__main__":
    run_converter()

# Additional implementation at 2025-08-04 06:29:45
LENGTH_FACTORS = {
    "meter": 1.0,
    "m": 1.0,
    "kilometer": 1000.0,
    "km": 1000.0,
    "centimeter": 0.01,
    "cm": 0.01,
    "millimeter": 0.001,
    "mm": 0.001,
    "mile": 1609.34,
    "mi": 1609.34,
    "yard": 0.9144,
    "yd": 0.9144,
    "foot": 0.3048,
    "ft": 0.3048,
    "inch": 0.0254,
    "in": 0.0254,
}

WEIGHT_FACTORS = {
    "kilogram": 1.0,
    "kg": 1.0,
    "gram": 0.001,
    "g": 0.001,
    "milligram": 0.000001,
    "mg": 0.000001,
    "pound": 0.453592,
    "lb": 0.453592,
    "ounce": 0.0283495,
    "oz": 0.0283495,
    "metric ton": 1000.0,
    "tonne": 1000.0,
    "us ton": 907.185,
    "short ton": 907.185,
    "long ton": 1016.05,
}

def convert_length(value, from_unit, to_unit):
    from_unit = from_unit.lower()
    to_unit = to_unit.lower()

    if from_unit not in LENGTH_FACTORS or to_unit not in LENGTH_FACTORS:
        return None, "Invalid length unit(s)."

    value_in_meters = value * LENGTH_FACTORS[from_unit]
    converted_value = value_in_meters / LENGTH_FACTORS[to_unit]
    return converted_value, None

def convert_weight(value, from_unit, to_unit):
    from_unit = from_unit.lower()
    to_unit = to_unit.lower()

    if from_unit not in WEIGHT_FACTORS or to_unit not in WEIGHT_FACTORS:
        return None, "Invalid weight unit(s)."

    value_in_kg = value * WEIGHT_FACTORS[from_unit]
    converted_value = value_in_kg / WEIGHT_FACTORS[to_unit]
    return converted_value, None

def convert_temperature(value, from_unit, to_unit):
    from_unit = from_unit.lower()
    to_unit = to_unit.lower()

    if from_unit in ["celsius", "c"]:
        celsius_value = value
    elif from_unit in ["fahrenheit", "f"]:
        celsius_value = (value - 32) * 5/9
    elif from_unit in ["kelvin", "k"]:
        celsius_value = value - 273.15
    else:
        return None, "Invalid temperature 'from' unit."

    if to_unit in ["celsius", "c"]:
        converted_value = celsius_value
    elif to_unit in ["fahrenheit", "f"]:
        converted_value = (celsius_value * 9/5) + 32
    elif to_unit in ["kelvin", "k"]:
        converted_value = celsius_value + 273.15
    else:
        return None, "Invalid temperature 'to' unit."

    return converted_value, None

def get_float_input(prompt):
    while True:
        try:
            return float(input(prompt))
        except ValueError:
            print("Invalid number. Please enter a numeric value.")

def get_unit_input(prompt, valid_units_dict):
    while True:
        unit = input(prompt).strip().lower()
        if unit in valid_units_dict:
            return unit
        else:
            print(f"Invalid unit. Please choose from: {', '.join(sorted(set(valid_units_dict.keys())))}")

def get_temperature_unit_input(prompt):
    valid_temp_units = {"celsius", "c", "fahrenheit", "f", "kelvin", "k"}
    while True:
        unit = input(prompt).strip().lower()
        if unit in valid_temp_units:
            return unit
        else:
            print(f"Invalid unit. Please choose from: {', '.join(sorted(valid_temp_units))}")

def main():
    print("Welcome to the Unit Converter CLI!")

    while True:
        print("\nChoose a conversion type:")
        print("1. Length")
        print("2. Weight")
        print("3. Temperature")
        print("4. Exit")

        choice = input("Enter your choice (1-4): ").strip()

        if choice == '1':
            print("\n--- Length Conversion ---")
            print("Available units: meter, km, cm, mm, mile, yard, foot, inch (and their abbreviations)")
            value = get_float_input("Enter the value to convert: ")
            from_unit = get_unit_input("Enter the unit to convert FROM (e.g., meter, ft): ", LENGTH_FACTORS)
            to_unit = get_unit_input("Enter the unit to convert TO (e.g., kilometer, in): ", LENGTH_FACTORS)
            
            result, error = convert_length(value, from_unit, to_unit)
            if error:
                print(f"Error: {error}")
            else:
                print(f"{value} {from_unit} is equal to {result:.4f} {to_unit}")

        elif choice == '2':
            print("\n--- Weight Conversion ---")
            print("Available units: kg, g, mg, lb, oz, metric ton, us ton, long ton (and their abbreviations)")
            value = get_float_input("Enter the value to convert: ")
            from_unit = get_unit_input("Enter the unit to convert FROM (e.g., kg, oz): ", WEIGHT_FACTORS)
            to_unit = get_unit_input("Enter the unit to convert TO (e.g., gram, lb): ", WEIGHT_FACTORS)
            
            result, error = convert_weight(value, from_unit, to_unit)
            if error:
                print(f"Error: {error}")
            else:
                print(f"{value} {from_unit} is equal to {result:.4f} {to_unit}")

        elif choice == '3':
            print("\n--- Temperature Conversion ---")
            print("Available units: Celsius (C), Fahrenheit (F), Kelvin (K)")
            value = get_float_input("Enter the value to convert: ")
            from_unit = get_temperature_unit_input("Enter the unit to convert FROM (e.g., C, F, K): ")
            to_unit = get_temperature_unit_input("Enter the unit to convert TO (e.g., F, C, K): ")
            
            result, error = convert_temperature(value, from_unit, to_unit)
            if error:
                print(f"Error: {error}")
            else:
                print(f"{value} {from_unit} is equal to {result:.4f} {to_unit}")

        elif choice == '4':
            print("Exiting the converter. Goodbye!")
            break
        else:
            print("Invalid choice. Please enter a number between 1 and 4.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-08-04 06:30:58
LENGTH_UNITS = {
    "m": 1.0,
    "km": 1000.0,
    "cm": 0.01,
    "mm": 0.001,
    "mi": 1609.34,
    "yd": 0.9144,
    "ft": 0.3048,
    "in": 0.0254
}

WEIGHT_UNITS = {
    "kg": 1.0,
    "g": 0.001,
    "mg": 0.000001,
    "lb": 0.453592,
    "oz": 0.0283495,
    "t": 1000.0 # Metric ton
}

def convert_length(value, from_unit, to_unit):
    if from_unit not in LENGTH_UNITS or to_unit not in LENGTH_UNITS:
        return None 

    value_in_meters = value * LENGTH_UNITS[from_unit]
    converted_value = value_in_meters / LENGTH_UNITS[to_unit]
    return converted_value

def convert_weight(value, from_unit, to_unit):
    if from_unit not in WEIGHT_UNITS or to_unit not in WEIGHT_UNITS:
        return None 

    value_in_kg = value * WEIGHT_UNITS[from_unit]
    converted_value = value_in_kg / WEIGHT_UNITS[to_unit]
    return converted_value

def convert_temperature(value, from_unit, to_unit):
    # Convert to Celsius first
    if from_unit == "C":
        celsius_value = value
    elif from_unit == "F":
        celsius_value = (value - 32) * 5/9
    elif from_unit == "K":
        celsius_value = value - 273.15
    else:
        return None 

    # Convert from Celsius to target unit
    if to_unit == "C":
        return celsius_value
    elif to_unit == "F":
        return celsius_value * 9/5 + 32
    elif to_unit == "K":
        return celsius_value + 273.15
    else:
        return None 

def main():
    while True:
        print("\n--- Unit Converter ---")
        print("Select a category:")
        print("1. Length")
        print("2. Weight")
        print("3. Temperature")
        print("4. Exit")

        category_choice = input("Enter your choice (1-4): ").strip()

        if category_choice == '4':
            print("Exiting converter. Goodbye!")
            break

        if category_choice not in ['1', '2', '3']:
            print("Invalid category choice. Please enter 1, 2, 3, or 4.")
            continue

        try:
            value_str = input("Enter the value to convert: ").strip()
            value = float(value_str)
        except ValueError:
            print("Invalid value. Please enter a number.")
            continue

        from_unit = input("Enter the unit to convert FROM (e.g., m, kg, C): ").strip().lower()
        to_unit = input("Enter the unit to convert TO (e.g., km, lb, F): ").strip().lower()

        result = None
        if category_choice == '1':
            if from_unit not in LENGTH_UNITS or to_unit not in LENGTH_UNITS:
                print(f"Invalid length unit(s). Available: {', '.join(LENGTH_UNITS.keys())}")
                continue
            result = convert_length(value, from_unit, to_unit)
        elif category_choice == '2':
            if from_unit not in WEIGHT_UNITS or to_unit not in WEIGHT_UNITS:
                print(f"Invalid weight unit(s). Available: {', '.join(WEIGHT_UNITS.keys())}")
                continue
            result = convert_weight(value, from_unit, to_unit)
        elif category_choice == '3':
            temp_units = ["c", "f", "k"]
            if from_unit not in temp_units or to_unit not in temp_units:
                print(f"Invalid temperature unit(s). Available: {', '.join(temp_units).upper()}")
                continue
            result = convert_temperature(value, from_unit.upper(), to_unit.upper()) 

        if result is not None:
            print(f"Result: {value} {from_unit} is {result:.4f} {to_unit}")
        else:
            print("Conversion failed. Check your units and try again.")

if __name__ == "__main__":
    main()