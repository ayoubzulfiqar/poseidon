import PIL.Image
import PIL.ImageDraw

def generate_barcode(data, filename="barcode.png"):
    CODE39_CHARS = {
        '0': '000110100', '1': '100100001', '2': '001100001', '3': '101100000',
        '4': '000110001', '5': '100110000', '6': '001110000', '7': '000100101',
        '8': '100100100', '9': '001100100', 'A': '100001001', 'B': '001001001',
        'C': '101001000', 'D': '000011001', 'E': '100011000', 'F': '001011000',
        'G': '000001101', 'H': '100001100', 'I': '001001100', 'J': '000011100',
        'K': '100000011', 'L': '001000011', 'M': '101000010', 'N': '000010011',
        'O': '100010010', 'P': '001010010', 'Q': '000000111', 'R': '100000110',
        'S': '001000110', 'T': '000010110', 'U': '110000001', 'V': '011000001',
        'W': '111000000', 'X': '010010001', 'Y': '110010000', 'Z': '011010000',
        '-': '010000101', '.': '110000100', ' ': '011000100', '*': '010010100',
        '$': '010101000', '/': '010100010', '+': '010001010', '%': '000101010'
    }
    allowed_chars = set(CODE39_CHARS.keys())
    allowed_chars.remove('*')
    data = data.upper()
    for char in data:
        if char not in allowed_chars:
            raise ValueError(f"Invalid character '{char}' for Code 39 barcode.")
    encoded_data = '*' + data + '*'
    unit_pixel_width = 2
    bar_height = 100
    padding_x = 20
    padding_y = 20
    total_units = 0
    for i, char in enumerate(encoded_data):
        pattern = CODE39_CHARS[char]
        for bit in pattern:
            total_units += (3 if bit == '1' else 1)
        if i < len(encoded_data) - 1:
            total_units += 1
    image_width = total_units * unit_pixel_width + 2 * padding_x
    image_height = bar_height + 2 * padding_y
    img = PIL.Image.new('RGB', (image_width, image_height), 'white')
    draw = PIL.ImageDraw.Draw(img)
    current_x = padding_x
    y_start = padding_y
    y_end = padding_y + bar_height
    for i, char in enumerate(encoded_data):
        pattern = CODE39_CHARS[char]
        for j, bit in enumerate(pattern):
            width = unit_pixel_width * (3 if bit == '1' else 1)
            color = 'black' if j % 2 == 0 else 'white'
            draw.rectangle([current_x, y_start, current_x + width, y_end], fill=color)
            current_x += width
        if i < len(encoded_data) - 1:
            gap_width = unit_pixel_width * 1
            draw.rectangle([current_x, y_start, current_x + gap_width, y_end], fill='white')
            current_x += gap_width
    img.save(filename)

# Additional implementation at 2025-06-18 01:02:54
import barcode
from barcode.writer import ImageWriter
import os

def generate_barcode(data: str, barcode_type: str, output_filename: str, options: dict = None):
    if options is None:
        options = {}

    BarcodeClass = barcode.get_barcode_class(barcode_type.lower())
    barcode_instance = BarcodeClass(data, writer=ImageWriter())
    
    output_dir = os.path.dirname(output_filename)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)

    full_path = f"{output_filename}.png"
    barcode_instance.save(full_path, options=options)

if __name__ == '__main__':
    generate_barcode("590123412345", "ean13", "barcodes/ean13_product_code", 
                     options={'module_height': 10, 'font_size': 8, 'text_distance': 3})

    generate_barcode("Hello World! 123", "code128", "barcodes/code128_shipping_label",
                     options={'module_width': 0.3, 'quiet_zone': 10, 'background': 'lightgray'})

    generate_barcode("01234567890", "upca", "barcodes/upca_item_id")

    generate_barcode("12345678", "itf", "barcodes/itf_container_id",
                     options={'foreground': 'blue', 'write_text': False})

    generate_barcode("123456789012", "ean13", "barcodes/ean13_custom_text",
                     options={'text': 'My Custom Product Code', 'font_size': 12})

    # Example of a barcode that would typically raise an error (uncomment to test)
    # try:
    #     generate_barcode("12345", "ean13", "barcodes/invalid_ean13")
    # except Exception as e:
    #     pass # Handle or ignore the error as needed

    # Example of an unsupported barcode type (uncomment to test)
    # try:
    #     generate_barcode("some_data", "unsupported_type", "barcodes/unsupported_barcode")
    # except Exception as e:
    #     pass # Handle or ignore the error as needed