import qrcode

def generate_ascii_qr(data, version=None, error_correction=qrcode.constants.ERROR_CORRECT_L, box_size=1, border=4, dark_char='██', light_char='  '):
    qr = qrcode.QRCode(
        version=version,
        error_correction=error_correction,
        box_size=box_size,
        border=border,
    )
    qr.add_data(data)
    qr.make(fit=True)

    modules = qr.modules

    ascii_output_lines = []
    for row in modules:
        line_chars = []
        for module in row:
            if module:
                line_chars.append(dark_char)
            else:
                line_chars.append(light_char)
        ascii_output_lines.append("".join(line_chars))

    return "\n".join(ascii_output_lines)

if __name__ == '__main__':
    qr_data_1 = "https://www.example.com/ascii-qr"
    ascii_qr_code_1 = generate_ascii_qr(qr_data_1, dark_char='##', light_char='  ')
    print(ascii_qr_code_1)

    print("\n")

    qr_data_2 = "Hello, ASCII QR Code!"
    ascii_qr_code_2 = generate_ascii_qr(qr_data_2, dark_char='█', light_char=' ')
    print(ascii_qr_code_2)

    print("\n")

    qr_data_3 = "This is a longer string to demonstrate a larger QR code. It will automatically adjust its version and size."
    ascii_qr_code_3 = generate_ascii_qr(qr_data_3)
    print(ascii_qr_code_3)

    print("\n")

    qr_data_4 = "Short"
    ascii_qr_code_4 = generate_ascii_qr(qr_data_4, version=1, error_correction=qrcode.constants.ERROR_CORRECT_H)
    print(ascii_qr_code_4)

# Additional implementation at 2025-06-18 00:43:21
import qrcode
from PIL import Image

def generate_ascii_qr(data, fill_char='██', empty_char='  ', box_size=1):
    qr = qrcode.QRCode(
        version=None,
        error_correction=qrcode.constants.ERROR_CORRECT_L,
        box_size=box_size,
        border=4,
    )
    qr.add_data(data)
    qr.make(fit=True)

    img = qr.make_image(fill_color="black", back_color="white")
    img = img.convert('1')

    width, height = img.size

    ascii_qr_output = []
    for y in range(height):
        row_chars = []
        for x in range(width):
            pixel = img.getpixel((x, y))
            if pixel == 0:
                row_chars.append(fill_char)
            else:
                row_chars.append(empty_char)
        ascii_qr_output.append("".join(row_chars))

    return "\n".join(ascii_qr_output)

if __name__ == "__main__":
    print("ASCII QR Code Generator")
    data_to_encode = input("Enter the data to encode: ")

    custom_fill_char = input("Enter fill character (e.g., '##', '██', default '██'): ")
    if not custom_fill_char:
        custom_fill_char = '██'
    elif len(custom_fill_char) == 1:
        custom_fill_char = custom_fill_char * 2

    custom_empty_char = input("Enter empty character (e.g., '  ', '--', default '  '): ")
    if not custom_empty_char:
        custom_empty_char = '  '
    elif len(custom_empty_char) == 1:
        custom_empty_char = custom_empty_char * 2

    try:
        custom_box_size = int(input("Enter box size (integer, default 1, larger for bigger QR): "))
    except ValueError:
        custom_box_size = 1

    if not data_to_encode:
        print("No data entered. Exiting.")
    else:
        qr_ascii = generate_ascii_qr(data_to_encode, custom_fill_char, custom_empty_char, custom_box_size)
        print("\nGenerated ASCII QR Code:\n")
        print(qr_ascii)
        print("\nScan this QR code with your phone!")

# Additional implementation at 2025-06-18 00:44:10
import qrcode
import argparse
import sys

def generate_ascii_qr(data, error_correction_level, black_char, white_char, border_size):
    qr = qrcode.QRCode(
        version=None,
        error_correction=error_correction_level,
        box_size=1,
        border=border_size,
    )
    qr.add_data(data)
    qr.make(fit=True)

    matrix = qr.get_matrix()

    for row in matrix:
        line = []
        for module in row:
            line.append(black_char if module else white_char)
        print("".join(line))

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate an ASCII art QR code from text with customizable options."
    )
    parser.add_argument(
        "text",
        type=str,
        help="The text data to encode in the QR code."
    )
    parser.add_argument(
        "--error-correction",
        type=str,
        default="L",
        choices=["L", "M", "Q", "H"],
        help="Error correction level (L, M, Q, H). Default is L."
    )
    parser.add_argument(
        "--black-char",
        type=str,
        default="██",
        help="Character(s) to represent a black module. Default is '██'."
    )
    parser.add_argument(
        "--white-char",
        type=str,
        default="  ",
        help="Character(s) to represent a white module. Default is '  '."
    )
    parser.add_argument(
        "--border-size",
        type=int,
        default=4,
        help="Size of the quiet zone border around the QR code modules. Default is 4."
    )

    args = parser.parse_args()

    ec_map = {
        "L": qrcode.constants.ERROR_CORRECT_L,
        "M": qrcode.constants.ERROR_CORRECT_M,
        "Q": qrcode.constants.ERROR_CORRECT_Q,
        "H": qrcode.constants.ERROR_CORRECT_H,
    }
    selected_ec = ec_map[args.error_correction.upper()]

    if len(args.black_char) != len(args.white_char):
        print("Error: --black-char and --white-char must have the same length.", file=sys.stderr)
        sys.exit(1)

    generate_ascii_qr(
        args.text,
        selected_ec,
        args.black_char,
        args.white_char,
        args.border_size
    )

# Additional implementation at 2025-06-18 00:45:06
import qrcode
import argparse

def generate_ascii_qr(data, scale=1, dark_char='██', light_char='  '):
    qr = qrcode.QRCode(
        version=None,
        error_correction=qrcode.constants.ERROR_CORRECT_L,
        box_size=1,
        border=4,
    )
    qr.add_data(data)
    qr.make(fit=True)

    matrix = qr.modules

    height = len(matrix)
    width = len(matrix[0])

    for r in range(height):
        for _ in range(scale):
            for c in range(width):
                char_to_print = dark_char if matrix[r][c] else light_char
                print(char_to_print * scale, end='')
            print()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate ASCII art QR codes.")
    parser.add_argument("data", type=str, help="The data to encode in the QR code.")
    parser.add_argument("-s", "--scale", type=int, default=1,
                        help="Scaling factor for the QR code. Each QR module will be (scale * len(char)) wide and (scale) lines tall. Default is 1.")
    parser.add_argument("-d", "--dark-char", type=str, default='██',
                        help="Character(s) to use for dark modules. Default is '██'.")
    parser.add_argument("-l", "--light-char", type=str, default='  ',
                        help="Character(s) to use for light modules. Default is '  '.")

    args = parser.parse_args()

    generate_ascii_qr(args.data, args.scale, args.dark_char, args.light_char)