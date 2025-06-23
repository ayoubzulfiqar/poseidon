import sys

def reverse_file_lines(input_filename):
    try:
        with open(input_filename, 'r') as f:
            lines = f.readlines()
        
        reversed_lines = lines[::-1]
        
        for line in reversed_lines:
            sys.stdout.write(line)
            
    except FileNotFoundError:
        sys.stderr.write(f"Error: File '{input_filename}' not found.\n")
        sys.exit(1)
    except Exception as e:
        sys.stderr.write(f"An unexpected error occurred: {e}\n")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.stderr.write("Usage: python script_name.py <input_filename>\n")
        sys.exit(1)
    
    input_file = sys.argv[1]
    reverse_file_lines(input_file)

# Additional implementation at 2025-06-23 02:01:50
import argparse
import os
import sys

def reverse_lines_in_file(input_filepath, output_filepath=None, in_place=False):
    if in_place and output_filepath:
        raise ValueError("Cannot specify both --in-place and an output file.")

    if not os.path.exists(input_filepath):
        raise FileNotFoundError(f"Input file not found: {input_filepath}")

    lines = []
    try:
        with open(input_filepath, 'r', encoding='utf-8') as infile:
            lines = infile.readlines()
    except Exception as e:
        raise IOError(f"Error reading input file {input_filepath}: {e}")

    if not lines:
        if in_place:
            return
        elif output_filepath:
            try:
                with open(output_filepath, 'w', encoding='utf-8') as outfile:
                    pass
                return
            except Exception as e:
                raise IOError(f"Error creating empty output file {output_filepath}: {e}")
        else:
            raise ValueError("Output target not specified for empty file.")

    lines.reverse()

    target_filepath = output_filepath
    temp_filepath = None

    if in_place:
        temp_filepath = input_filepath + ".tmp_reversed"
        target_filepath = temp_filepath
    elif not target_filepath:
        raise ValueError("Output file path must be specified if not in-place.")

    try:
        with open(target_filepath, 'w', encoding='utf-8') as outfile:
            outfile.writelines(lines)
    except Exception as e:
        if temp_filepath and os.path.exists(temp_filepath):
            os.remove(temp_filepath)
        raise IOError(f"Error writing to output file {target_filepath}: {e}")

    if in_place:
        try:
            os.replace(temp_filepath, input_filepath)
        except Exception as e:
            raise IOError(f"Error replacing original file {input_filepath} with reversed content: {e}")

def main():
    parser = argparse.ArgumentParser(
        description="Reverses the order of lines in a text file."
    )

    parser.add_argument(
        'input_file',
        type=str,
        help='Path to the input text file.'
    )

    output_group = parser.add_mutually_exclusive_group(required=True)

    output_group.add_argument(
        '-o', '--output',
        type=str,
        dest='output_file',
        help='Path to the output file. Cannot be used with --in-place.'
    )
    output_group.add_argument(
        '-i', '--in-place',
        action='store_true',
        dest='in_place',
        help='Reverse lines in-place, overwriting the original file. '
             'Cannot be used with --output.'
    )

    args = parser.parse_args()

    try:
        reverse_lines_in_file(
            input_filepath=args.input_file,
            output_filepath=args.output_file,
            in_place=args.in_place
        )
    except FileNotFoundError:
        sys.exit(1)
    except IOError:
        sys.exit(1)
    except ValueError:
        sys.exit(1)
    except Exception:
        sys.exit(1)

if __name__ == '__main__':
    main()