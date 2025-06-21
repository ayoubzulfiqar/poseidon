import csv
import os

def convert_csv_to_fixed_width(input_csv_path, output_fw_path, column_specs):
    try:
        with open(input_csv_path, mode='r', newline='', encoding='utf-8') as infile:
            reader = csv.reader(infile)
            header = next(reader)

            header_map = {col_name: i for i, col_name in enumerate(header)}

            with open(output_fw_path, mode='w', encoding='utf-8') as outfile:
                # Write header row
                formatted_header_parts = []
                for spec in column_specs:
                    col_name = spec['name']
                    width = spec['width']
                    align = spec['align']
                    
                    header_text = col_name
                    if len(header_text) > width:
                        header_text = header_text[:width]

                    if align == 'left':
                        formatted_header_parts.append(header_text.ljust(width))
                    elif align == 'right':
                        formatted_header_parts.append(header_text.rjust(width))
                    elif align == 'center':
                        formatted_header_parts.append(header_text.center(width))
                    else:
                        formatted_header_parts.append(header_text.ljust(width))

                outfile.write("".join(formatted_header_parts) + '\n')

                # Write data rows
                for row in reader:
                    formatted_row_parts = []
                    for spec in column_specs:
                        col_name = spec['name']
                        width = spec['width']
                        align = spec['align']

                        col_index = header_map.get(col_name)
                        if col_index is not None and col_index < len(row):
                            value = str(row[col_index])
                        else:
                            value = ""

                        if len(value) > width:
                            value = value[:width]

                        if align == 'left':
                            formatted_row_parts.append(value.ljust(width))
                        elif align == 'right':
                            formatted_row_parts.append(value.rjust(width))
                        elif align == 'center':
                            formatted_row_parts.append(value.center(width))
                        else:
                            formatted_row_parts.append(value.ljust(width))

                    outfile.write("".join(formatted_row_parts) + '\n')

    except FileNotFoundError:
        print(f"Error: One of the files not found. Input: {input_csv_path}, Output: {output_fw_path}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    input_csv_file = "input.csv"
    output_fw_file = "output.txt"

    csv_content = """ID,Name,Age,City
1,Alice Wonderland,30,New York
2,Bob The Builder,25,London
3,Charlie Chaplin,40,Paris
100,David Copperfield,50,Los Angeles,Extra
1000,Eve Green,22,Sydney
"""
    with open(input_csv_file, 'w', newline='', encoding='utf-8') as f:
        f.write(csv_content)

    column_specifications = [
        {'name': 'ID', 'width': 6, 'align': 'right'},
        {'name': 'Name', 'width': 20, 'align': 'left'},
        {'name': 'Age', 'width': 5, 'align': 'center'},
        {'name': 'City', 'width': 15, 'align': 'left'},
        {'name': 'Country', 'width': 10, 'align': 'left'}
    ]

    print(f"Converting '{input_csv_file}' to fixed-width format in '{output_fw_file}'...")
    convert_csv_to_fixed_width(input_csv_file, output_fw_file, column_specifications)
    print("Conversion complete.")

    if os.path.exists(output_fw_file):
        print("\n--- Content of output.txt ---")
        with open(output_fw_file, 'r', encoding='utf-8') as f:
            print(f.read())
        print("-----------------------------")
    else:
        print(f"Output file '{output_fw_file}' was not created.")

    # Uncomment the following lines to clean up the dummy files after execution
    # os.remove(input_csv_file)
    # os.remove(output_fw_file)
    # print("Cleaned up dummy files.")

# Additional implementation at 2025-06-20 23:17:54
import csv
import os

def csv_to_fixed_width(input_csv_path, output_fw_path, column_specs, include_header=True, padding_char=' '):
    if not os.path.exists(input_csv_path):
        raise FileNotFoundError(f"Input CSV file not found: {input_csv_path}")

    try:
        with open(input_csv_path, mode='r', newline='', encoding='utf-8') as infile, \
             open(output_fw_path, mode='w', encoding='utf-8') as outfile:

            reader = csv.reader(infile)
            header = next(reader)

            column_indices = {}
            for spec in column_specs:
                try:
                    column_indices[spec['name']] = header.index(spec['name'])
                except ValueError:
                    column_indices[spec['name']] = -1

            def format_value(value, width, align, pad_char):
                s_value = str(value if value is not None else '')
                if len(s_value) > width:
                    s_value = s_value[:width]
                
                if align == 'right':
                    return s_value.rjust(width, pad_char)
                elif align == 'center':
                    return s_value.center(width, pad_char)
                else:
                    return s_value.ljust(width, pad_char)

            if include_header:
                formatted_header_parts = []
                for spec in column_specs:
                    col_name = spec['name']
                    col_width = spec['width']
                    col_align = spec.get('align', 'left')
                    formatted_header_parts.append(format_value(col_name, col_width, col_align, padding_char))
                outfile.write("".join(formatted_header_parts) + '\n')

            for row in reader:
                formatted_row_parts = []
                for spec in column_specs:
                    col_name = spec['name']
                    col_width = spec['width']
                    col_align = spec.get('align', 'left')
                    
                    col_index = column_indices.get(col_name)
                    
                    value = ''
                    if col_index != -1 and col_index < len(row):
                        value = row[col_index]
                    
                    formatted_row_parts.append(format_value(value, col_width, col_align, padding_char))
                outfile.write("".join(formatted_row_parts) + '\n')

    except Exception as e:
        raise

csv_content = """ID,Product Name,Price,Quantity,Description,Category
1,Laptop Pro X,1200.50,50,High performance laptop,Electronics
2,Mechanical Keyboard,75.99,150,RGB backlit keyboard,Electronics
3,Wireless Mouse,25.00,300,Ergonomic design,Electronics
4,USB-C Hub,40.00,100,Multi-port adapter,Accessories
5,External SSD 1TB,150.00,75,Fast portable storage,Storage
6,Gaming Headset,99.99,80,Surround sound,Audio
7,Webcam 1080p,60.00,120,Full HD video,Peripherals
8,Monitor 27-inch,300.00,40,4K IPS display,Displays
9,Desk Lamp LED,35.00,200,Adjustable brightness,Home Office
10,Ergonomic Chair,250.00,30,Lumbar support,Furniture
"""

input_csv_file = "products.csv"
output_fw_file = "products_fixed_width.txt"

with open(input_csv_file, "w", newline='', encoding='utf-8') as f:
    f.write(csv_content)

column_definitions = [
    {'name': 'ID', 'width': 5, 'align': 'right'},
    {'name': 'Product Name', 'width': 25, 'align': 'left'},
    {'name': 'Price', 'width': 12, 'align': 'right'},
    {'name': 'Quantity', 'width': 10, 'align': 'center'},
    {'name': 'SKU', 'width': 8, 'align': 'left'},
    {'name': 'Description', 'width': 35, 'align': 'left'},
]

csv_to_fixed_width(input_csv_file, output_fw_file, column_definitions, include_header=True, padding_char=' ')

output_fw_file_no_header = "products_fixed_width_no_header.txt"
csv_to_fixed_width(input_csv_file, output_fw_file_no_header, column_definitions, include_header=False, padding_char='-')

# os.remove(input_csv_file)
# os.remove(output_fw_file)
# os.remove(output_fw_file_no_header)

# Additional implementation at 2025-06-20 23:19:11
import csv
import os

def csv_to_fixed_width(
    input_csv_path,
    output_fixed_width_path,
    column_definitions,
    include_header=True,
    padding_char=' ',
    alignment='left',
    truncate=True,
    input_encoding='utf-8',
    output_encoding='utf-8'
):
    """
    Converts a CSV file to a fixed-width column text file.

    Args:
        input_csv_path (str): Path to the input CSV file.
        output_fixed_width_path (str): Path to the output fixed-width file.
        column_definitions (list of tuple): A list of tuples, where each tuple
                                            is (column_name, width).
                                            Example: [('Name', 20), ('Age', 5), ('City', 30)]
                                            The order in this list determines the output column order.
        include_header (bool): If True, the header row will be written to the
                               fixed-width file. Defaults to True.
        padding_char (str): The character used for padding. Defaults to ' '.
        alignment (str): Alignment for columns: 'left', 'right', or 'center'.
                         Defaults to 'left'.
        truncate (bool): If True, data longer than the specified width will be
                         truncated. If False, an error might occur or output
                         might be malformed if data exceeds width. Defaults to True.
        input_encoding (str): Encoding for the input CSV file. Defaults to 'utf-8'.
        output_encoding (str): Encoding for the output fixed-width file. Defaults to 'utf-8'.
    """

    if alignment not in ['left', 'right', 'center']:
        raise ValueError("Alignment must be 'left', 'right', or 'center'.")
    if not isinstance(padding_char, str) or len(padding_char) != 1:
        raise ValueError("Padding character must be a single character string.")

    try:
        with open(input_csv_path, mode='r', newline='', encoding=input_encoding) as infile:
            reader = csv.reader(infile)
            header = next(reader) # Read the header row

            # Create a mapping from desired column names to their index in the CSV
            column_name_to_index = {name: idx for idx, name in enumerate(header)}

            # Prepare the output column structure based on column_definitions
            # This will store (csv_index, width, column_name_for_header)
            output_columns_info = []
            for col_name, width in column_definitions:
                if col_name in column_name_to_index:
                    output_columns_info.append((column_name_to_index[col_name], width, col_name))
                else:
                    # If a defined column doesn't exist in the CSV, treat it as empty
                    # and still include it in the output structure with its width.
                    # We use -1 as a placeholder for index to indicate it's not found.
                    output_columns_info.append((-1, width, col_name))

            with open(output_fixed_width_path, mode='w', encoding=output_encoding) as outfile:
                # Write header if requested
                if include_header:
                    formatted_header_parts = []
                    for _, width, col_name_for_header in output_columns_info:
                        value = str(col_name_for_header)
                        if truncate and len(value) > width:
                            value = value[:width]
                        if alignment == 'left':
                            formatted_header_parts.append(value.ljust(width, padding_char))
                        elif alignment == 'right':
                            formatted_header_parts.append(value.rjust(width, padding_char))
                        else: # center
                            formatted_header_parts.append(value.center(width, padding_char))
                    outfile.write("".join(formatted_header_parts) + '\n')

                # Write data rows
                for row in reader:
                    formatted_row_parts = []
                    for csv_index, width, _ in output_columns_info:
                        value = ''
                        if csv_index != -1 and csv_index < len(row):
                            value = str(row[csv_index])

                        if truncate and len(value) > width:
                            value = value[:width]

                        if alignment == 'left':
                            formatted_row_parts.append(value.ljust(width, padding_char))
                        elif alignment == 'right':
                            formatted_row_parts.append(value.rjust(width, padding_char))
                        else: # center
                            formatted_row_parts.append(value.center(width, padding_char))
                    outfile.write("".join(formatted_row_parts) + '\n')

    except FileNotFoundError:
        print(f"Error: Input file not found at {input_csv_path}")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == '__main__':
    # Example Usage:

    # 1. Create a dummy CSV file for testing
    dummy_csv_content = """Name,Age,City,Occupation,Notes
Alice,30,New York,Engineer,Likes to read books
Bob,25,Los Angeles,Artist,Plays guitar
Charlie,35,Chicago,Doctor,Runs marathons
David,40,Houston,Teacher,Enjoys cooking and baking
Eve,28,Miami,Designer,Travels frequently
Frank,50,San Francisco,Manager,Golf enthusiast
Grace,22,Seattle,Student,Loves hiking and nature photography
"""
    input_csv_file = "input.csv"
    with open(input_csv_file, "w", newline="", encoding="utf-8") as f:
        f.write(dummy_csv_content)

    output_fixed_width_file = "output_fixed_width.txt"

    # Define column names and their desired widths and order
    # Note: 'Occupation' is included, 'Notes' is excluded.
    # 'ZipCode' is included but not in CSV, so it will be blank.
    column_definitions = [
        ('Name', 15),
        ('Age', 5),
        ('City', 20),
        ('Occupation', 15),
        ('ZipCode', 10) # This column does not exist in the CSV, will be blank
    ]

    print(f"Converting '{input_csv_file}' to '{output_fixed_width_file}'...")

    # Basic conversion
    csv_to_fixed_width(
        input_csv_file,
        output_fixed_width_file,
        column_definitions
    )
    print("Conversion complete (default settings). Check 'output_fixed_width.txt'")

    # Example with different settings: no header, right alignment, custom padding, truncation
    output_fixed_width_file_2 = "output_fixed_width_no_header_right_align.txt"
    column_definitions_2 = [
        ('Age', 5),
        ('Name', 15), # Order changed
        ('City', 10), # Shorter width for City, will truncate
        ('Occupation', 15)
    ]
    csv_to_fixed_width(
        input_csv_file,
        output_fixed_width_file_2,
        column_definitions_2,
        include_header=False,
        padding_char='.',
        alignment='right',
        truncate=True
    )
    print("Conversion complete (no header, right align, custom padding). Check 'output_fixed_width_no_header_right_align.txt'")

    # Example with center alignment
    output_fixed_width_file_3 = "output_fixed_width_center_align.txt"
    column_definitions_3 = [
        ('Name', 15),
        ('Age', 5),
        ('City', 20)
    ]
    csv_to_fixed_width(
        input_csv_file,
        output_fixed_width_file_3,
        column_definitions_3,
        include_header=True,
        padding_char=' ',
        alignment='center'
    )
    print("Conversion complete (center align). Check 'output_fixed_width_center_align.txt'")

    # Clean up dummy files
    # os.remove(input_csv_file)
    # os.remove(output_fixed_width_file)
    # os.remove(output_fixed_width_file_2)
    # os.remove(output_fixed_width_file_3)
    # print("\nCleaned up dummy files.")

# Additional implementation at 2025-06-20 23:20:19
import csv
import json
import argparse
import sys

def format_fixed_width_field(value, width, alignment='left', pad_char=' ', truncate=True):
    """
    Formats a single value into a fixed-width string.

    Args:
        value: The value to format.
        width: The desired fixed width.
        alignment: 'left', 'right', or 'center'.
        pad_char: The character to use for padding.
        truncate: If True, truncate the value if it exceeds the width.
                  If False, raise an error if it exceeds the width.

    Returns:
        A string formatted to the specified fixed width.

    Raises:
        ValueError: If truncate is False and the value exceeds the width,
                    or if an invalid alignment is specified.
    """
    s_value = str(value)

    if len(s_value) > width:
        if truncate:
            s_value = s_value[:width]
        else:
            raise ValueError(f"Value '{s_value}' (length {len(s_value)}) exceeds specified width {width} and truncation is disabled.")

    if alignment == 'left':
        return s_value.ljust(width, pad_char)
    elif alignment == 'right':
        return s_value.rjust(width, pad_char)
    elif alignment == 'center':
        return s_value.center(width, pad_char)
    else:
        raise ValueError(f"Invalid alignment specified: '{alignment}'. Must be 'left', 'right', or 'center'.")

def convert_csv_to_fixed_width(config_path):
    """
    Converts a CSV file to a fixed-width text file based on a JSON configuration.

    Args:
        config_path: Path to the JSON configuration file.
    """
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            config = json.load(f)
    except FileNotFoundError:
        print(f"Error: Configuration file not found at '{config_path}'", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in configuration file '{config_path}'", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error loading configuration: {e}", file=sys.stderr)
        sys.exit(1)

    input_csv_path = config.get('input_csv_path')
    output_fw_path = config.get('output_fw_path')
    include_header = config.get('include_header', True) # Default to true
    columns_config = config.get('columns')

    if not input_csv_path or not output_fw_path or not columns_config:
        print("Error: Configuration must include 'input_csv_path', 'output_fw_path', and 'columns'.", file=sys.stderr)
        sys.exit(1)
    if not isinstance(columns_config, list) or not all(isinstance(col, dict) for col in columns_config):
        print("Error: 'columns' in configuration must be a list of dictionaries.", file=sys.stderr)
        sys.exit(1)

    try:
        with open(input_csv_path, 'r', newline='', encoding='utf-8') as infile, \
             open(output_fw_path, 'w', encoding='utf-8') as outfile:
            reader = csv.reader(infile)

            header = next(reader) # Read header row from CSV
            csv_column_map = {col_name: i for i, col_name in enumerate(header)}

            # Validate column configurations and prepare formatting parameters
            formatted_columns_info = []
            for col_def in columns_config:
                col_name = col_def.get('name')
                col_width = col_def.get('width')
                col_alignment = col_def.get('alignment', 'left')
                col_pad_char = col_def.get('pad_char', ' ')
                col_truncate = col_def.get('truncate', True)

                if not col_name or not isinstance(col_width, int) or col_width <= 0:
                    print(f"Error: Invalid column definition in config: {col_def}. 'name' (string) and positive 'width' (integer) are required.", file=sys.stderr)
                    sys.exit(1)

                if col_name not in csv_column_map:
                    print(f"Warning: Column '{col_name}' specified in config not found in CSV header. Will output empty field for this column.", file=sys.stderr)
                    csv_col_index = -1 # Sentinel value for not found
                else:
                    csv_col_index = csv_column_map[col_name]

                formatted_columns_info.append({
                    'index': csv_col_index,
                    'width': col_width,
                    'alignment': col_alignment,
                    'pad_char': col_pad_char,
                    'truncate': col_truncate,
                    'name': col_name # Keep name for header formatting
                })

            # Write header if required
            if include_header:
                header_line_parts = []
                for col_info in formatted_columns_info:
                    # Use the column name itself for the header line, formatted according to its column's rules
                    header_line_parts.append(
                        format_fixed_width_field(
                            col_info['name'],
                            col_info['width'],
                            col_info['alignment'],
                            col_info['pad_char'],
                            col_info['truncate']
                        )
                    )
                outfile.write("".join(header_line_parts) + '\n')

            # Process data rows
            for row_num, row in enumerate(reader, start=2): # Start at 2 because header is row 1
                line_parts = []
                try:
                    for col_info in formatted_columns_info:
                        # Get value from CSV row; if column not found, use empty string
                        value = row[col_info['index']] if col_info['index'] != -1 and col_info['index'] < len(row) else ""
                        line_parts.append(
                            format_fixed_width_field(
                                value,
                                col_info['width'],
                                col_info['alignment'],
                                col_info['pad_char'],
                                col_info['truncate']
                            )
                        )
                    outfile.write("".join(line_parts) + '\n')
                except IndexError:
                    print(f"Error: Row {row_num} in CSV has fewer columns than expected. Skipping row.", file=sys.stderr)
                    # This error occurs if a column index is out of bounds for the current row
                except ValueError as ve:
                    print(f"Error processing row {row_num}: {ve}. Skipping row.", file=sys.stderr)
                    # This catches errors from format_fixed_width_field, e.g., truncation disabled for oversized data

    except FileNotFoundError:
        print(f"Error: Input CSV file not found at '{input_csv_path}'", file=sys.stderr)
        sys.exit(1)
    except IOError as e:
        print(f"Error reading/writing files: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"An unexpected error occurred: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Convert CSV to fixed-width format based on a JSON configuration.")
    parser.add_argument('config_file', type=str,
                        help="Path to the JSON configuration file.")

    args = parser.parse_args()

    convert_csv_to_fixed_width(args.config_file)