import codecs

def detect_file_encoding(filepath):
    try:
        with open(filepath, 'rb') as f:
            full_content_bytes = f.read()

            if not full_content_bytes:
                return 'Empty File'

            if full_content_bytes.startswith(codecs.BOM_UTF8):
                return 'UTF-8-BOM'
            elif full_content_bytes.startswith(codecs.BOM_UTF16_LE):
                return 'UTF-16-LE'
            elif full_content_bytes.startswith(codecs.BOM_UTF16_BE):
                return 'UTF-16-BE'

            try:
                full_content_bytes.decode('utf-16')
                return 'UTF-16'
            except UnicodeDecodeError:
                pass

            try:
                decoded_text = full_content_bytes.decode('utf-8')
                try:
                    decoded_text.encode('ascii')
                    return 'ASCII'
                except UnicodeEncodeError:
                    return 'UTF-8'
            except UnicodeDecodeError:
                pass

            return 'Unknown/Binary'

    except FileNotFoundError:
        return 'File Not Found'
    except Exception:
        return 'Error'

# Additional implementation at 2025-06-20 00:21:55
import os

def detect_file_encoding(file_path, sample_size=4096):
    result = {
        'encoding': 'unknown',
        'confidence': 0.0,
        'bom': False,
        'error': None
    }

    if not os.path.exists(file_path):
        result['error'] = f"File not found: {file_path}"
        return result
    if not os.path.isfile(file_path):
        result['error'] = f"Path is not a file: {file_path}"
        return result

    try:
        with open(file_path, 'rb') as f:
            raw_bytes = f.read(sample_size)
            if not raw_bytes:
                result['encoding'] = 'empty'
                result['confidence'] = 1.0
                return result

            utf8_bom = b'\xef\xbb\xbf'
            utf16_le_bom = b'\xff\xfe'
            utf16_be_bom = b'\xfe\xff'

            if raw_bytes.startswith(utf8_bom):
                try:
                    raw_bytes.decode('utf-8')
                    result['encoding'] = 'utf-8'
                    result['confidence'] = 1.0
                    result['bom'] = True
                    return result
                except UnicodeDecodeError:
                    pass

            elif raw_bytes.startswith(utf16_le_bom):
                try:
                    raw_bytes.decode('utf-16-le')
                    result['encoding'] = 'utf-16-le'
                    result['confidence'] = 1.0
                    result['bom'] = True
                    return result
                except UnicodeDecodeError:
                    pass

            elif raw_bytes.startswith(utf16_be_bom):
                try:
                    raw_bytes.decode('utf-16-be')
                    result['encoding'] = 'utf-16-be'
                    result['confidence'] = 1.0
                    result['bom'] = True
                    return result
                except UnicodeDecodeError:
                    pass

            try:
                decoded_text = raw_bytes.decode('utf-8')
                try:
                    decoded_text.encode('ascii')
                    result['encoding'] = 'ascii'
                    result['confidence'] = 1.0
                    return result
                except UnicodeEncodeError:
                    result['encoding'] = 'utf-8'
                    result['confidence'] = 0.9
                    return result
            except UnicodeDecodeError:
                pass

            try:
                raw_bytes.decode('utf-16')
                result['encoding'] = 'utf-16-be'
                result['confidence'] = 0.7
                return result
            except UnicodeDecodeError:
                pass

            try:
                raw_bytes.decode('utf-16-le')
                result['encoding'] = 'utf-16-le'
                result['confidence'] = 0.7
                return result
            except UnicodeDecodeError:
                pass

    except IOError as e:
        result['error'] = f"IOError reading file: {e}"
    except Exception as e:
        result['error'] = f"An unexpected error occurred: {e}"

    return result

def create_test_file(filename, content, encoding='utf-8', add_bom=False):
    if encoding == 'utf-8' and add_bom:
        with open(filename, 'wb') as f:
            f.write(b'\xef\xbb\xbf')
            f.write(content.encode('utf-8'))
        return

    mode = 'wb' if encoding.startswith('utf-16') else 'w'
    with open(filename, mode, encoding=encoding, newline='') as f:
        f.write(content)

def main():
    test_files_data = {
        "ascii_file.txt": "Hello, this is an ASCII file.",
        "utf8_no_bom.txt": "Hello, this is a UTF-8 file with no BOM. Привет!",
        "utf8_with_bom.txt": "Hello, this is a UTF-8 file with BOM. Привет!",
        "utf16_le_file.txt": "Hello, this is a UTF-16 LE file. Привет!",
        "utf16_be_file.txt": "Hello, this is a UTF-16 BE file. Привет!",
        "empty_file.txt": "",
        "non_existent_file.txt": None,
        "binary_data.bin": b'\x00\x01\x02\x03\x80\x81\x82\x83\xff\xfe\x00\x00',
    }

    for filename, content in test_files_data.items():
        if content is None:
            continue
        if filename == "utf8_with_bom.txt":
            create_test_file(filename, content, encoding='utf-8', add_bom=True)
        elif filename == "utf16_le_file.txt":
            create_test_file(filename, content, encoding='utf-16-le')
        elif filename == "utf16_be_file.txt":
            create_test_file(filename, content, encoding='utf-16-be')
        elif filename == "binary_data.bin":
            with open(filename, 'wb') as f:
                f.write(content)
        else:
            create_test_file(filename, content, encoding='utf-8')

    results = []
    for filename in test_files_data:
        detection_result = detect_file_encoding(filename)
        results.append(f"File: '{filename}' -> {detection_result}")

    for filename in test_files_data:
        if os.path.exists(filename):
            os.remove(filename)

    for result_str in results:
        print(result_str)

if __name__ == "__main__":
    main()