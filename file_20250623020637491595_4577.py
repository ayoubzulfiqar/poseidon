import PyPDF2
import os

def extract_text_from_pdf(pdf_path):
    if not os.path.exists(pdf_path):
        print(f"Error: File not found at {pdf_path}")
        return ""
    try:
        with open(pdf_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)
            extracted_text = ""
            for page_num in range(len(reader.pages)):
                page = reader.pages[page_num]
                extracted_text += page.extract_text() + "\n"
            return extracted_text
    except PyPDF2.errors.PdfReadError:
        print(f"Error: Could not read PDF file {pdf_path}. It might be encrypted or corrupted.")
        return ""
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return ""

if __name__ == "__main__":
    pdf_file_path = "sample.pdf"
    extracted_content = extract_text_from_pdf(pdf_file_path)
    if extracted_content:
        print(extracted_content)

# Additional implementation at 2025-06-23 02:07:02
import PyPDF2
import argparse
import os

def extract_text_from_pdf(pdf_path, output_path=None, start_page=0, end_page=None, password=None):
    """
    Extracts text from a PDF file within a specified page range and saves it to a file.

    Args:
        pdf_path (str): The path to the input PDF file.
        output_path (str, optional): The path to save the extracted text. If None, prints to console.
        start_page (int, optional): The starting page number (0-indexed). Defaults to 0.
        end_page (int, optional): The ending page number (0-indexed, inclusive). Defaults to last page.
        password (str, optional): The password for encrypted PDFs.
    """
    extracted_text_parts = []
    try:
        with open(pdf_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)

            if reader.is_encrypted:
                if password:
                    try:
                        reader.decrypt(password)
                    except PyPDF2.errors.FileKeyError:
                        print(f"Error: Incorrect password for PDF '{pdf_path}'.")
                        return
                    except Exception as e:
                        print(f"Error decrypting PDF '{pdf_path}': {e}")
                        return
                else:
                    print(f"Error: PDF '{pdf_path}' is encrypted. Please provide a password using -p or --password.")
                    return

            total_pages = len(reader.pages)
            
            # Adjust end_page if it's None or out of bounds
            if end_page is None or end_page >= total_pages:
                end_page = total_pages - 1

            # Validate start_page and end_page
            if start_page < 0:
                start_page = 0
            if start_page > end_page:
                print(f"Warning: Start page ({start_page + 1}) is greater than end page ({end_page + 1}). No text extracted.")
                return
            if start_page >= total_pages:
                print(f"Warning: Start page ({start_page + 1}) is beyond the total number of pages ({total_pages}). No text extracted.")
                return

            print(f"Extracting text from pages {start_page + 1} to {end_page + 1} of '{pdf_path}'...")

            for page_num in range(start_page, end_page + 1):
                if page_num < total_pages: # Ensure page_num is within actual bounds
                    page = reader.pages[page_num]
                    text = page.extract_text()
                    if text:
                        extracted_text_parts.append(f"--- Page {page_num + 1} ---\n{text.strip()}\n")
                    else:
                        extracted_text_parts.append(f"--- Page {page_num + 1} (No extractable text found) ---\n")
                else:
                    # This case should ideally be caught by the end_page adjustment, but as a safeguard
                    print(f"Warning: Page {page_num + 1} is out of bounds. Stopping extraction.")
                    break

    except FileNotFoundError:
        print(f"Error: PDF file not found at '{pdf_path}'")
        return
    except PyPDF2.errors.PdfReadError as e:
        print(f"Error reading PDF '{pdf_path}': {e}")
        return
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return

    full_text = "\n".join(extracted_text_parts)

    if output_path:
        try:
            with open(output_path, 'w', encoding='utf-8') as outfile:
                outfile.write(full_text)
            print(f"Successfully extracted text to '{output_path}'")
        except IOError as e:
            print(f"Error writing to output file '{output_path}': {e}")
    else:
        print("\n--- Extracted Text ---\n")
        print(full_text)
        print("\n--- End of Extracted Text ---")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Extract text from PDF files with page range and password support."
    )
    parser.add_argument("pdf_path", help="Path to the input PDF file.")
    parser.add_argument("-o", "--output", dest="output_path",
                        help="Path to the output text file. If not provided, text is printed to console.")
    parser.add_argument("-s", "--start-page", type=int, default=1,
                        help="Starting page number (1-indexed). Defaults to 1.")
    parser.add_argument("-e", "--end-page", type=int,
                        help="Ending page number (1-indexed, inclusive). Defaults to the last page.")
    parser.add_argument("-p", "--password", help="Password for encrypted PDF files.")

    args = parser.parse_args()

    # Adjust page numbers to be 0-indexed for PyPDF2
    start_page_0_indexed = args.start_page - 1
    end_page_0_indexed = (args.end_page - 1) if args.end_page is not None else None

    extract_text_from_pdf(
        pdf_path=args.pdf_path,
        output_path=args.output_path,
        start_page=start_page_0_indexed,
        end_page=end_page_0_indexed,
        password=args.password
    )

# Additional implementation at 2025-06-23 02:07:54
import os
from pypdf import PdfReader, PdfReadError

def extract_text_from_pdf(pdf_path, output_file_path=None, start_page=None, end_page=None):
    if not os.path.exists(pdf_path):
        print(f"Error: PDF file not found at '{pdf_path}'")
        return None

    try:
        reader = PdfReader(pdf_path)
        num_pages = len(reader.pages)

        actual_start_page_idx = 0
        actual_end_page_idx = num_pages

        if start_page is not None:
            if not isinstance(start_page, int) or start_page < 1:
                print("Warning: start_page must be a positive integer. Ignoring invalid start_page.")
            else:
                actual_start_page_idx = max(0, start_page - 1)

        if end_page is not None:
            if not isinstance(end_page, int) or end_page < 1:
                print("Warning: end_page must be a positive integer. Ignoring invalid end_page.")
            else:
                actual_end_page_idx = min(num_pages, end_page)

        if actual_start_page_idx >= num_pages:
            print(f"Warning: start_page ({start_page if start_page is not None else 'default'}) is beyond the last page ({num_pages}). No text will be extracted.")
            return ""
        
        if actual_end_page_idx <= actual_start_page_idx:
            print(f"Warning: end_page ({end_page if end_page is not None else 'default'}) is before or same as start_page ({start_page if start_page is not None else 'default'}). No text will be extracted.")
            return ""
        
        extracted_texts = []
        for i in range(actual_start_page_idx, actual_end_page_idx):
            try:
                page = reader.pages[i]
                text = page.extract_text()
                if text:
                    extracted_texts.append(text)
            except Exception as page_error:
                print(f"Warning: Could not extract text from page {i+1}: {page_error}")

        full_text = "\n".join(extracted_texts)

        if output_file_path:
            try:
                with open(output_file_path, "w", encoding="utf-8") as f:
                    f.write(full_text)
            except IOError as e:
                print(f"Error: Could not write to output file '{output_file_path}': {e}")
                return None
        
        return full_text

    except PdfReadError as e:
        print(f"Error: Could not read PDF file '{pdf_path}'. It might be corrupted or encrypted: {e}")
        return None
    except Exception as e:
        print(f"An unexpected error occurred while processing '{pdf_path}': {e}")
        return None

# Additional implementation at 2025-06-23 02:08:31
import pypdf
import os

def extract_text_from_pdf(
    pdf_path: str,
    start_page: int = 0,
    end_page: int = None,
    output_txt_path: str = None,
    page_separator: str = "\n--- Page {page_num} ---\n"
) -> str:
    extracted_text = []
    try:
        with open(pdf_path, 'rb') as file:
            reader = pypdf.PdfReader(file)
            total_pages = len(reader.pages)

            if end_page is None or end_page > total_pages:
                end_page = total_pages
            
            if start_page < 0:
                start_page = 0
            if start_page >= total_pages:
                return ""

            if start_page >= end_page:
                return ""

            for i in range(start_page, end_page):
                if i < total_pages:
                    page = reader.pages[i]
                    text = page.extract_text()
                    if text:
                        if page_separator and i > start_page:
                             extracted_text.append(page_separator.format(page_num=i + 1))
                        extracted_text.append(text)
                else:
                    break

    except FileNotFoundError:
        print(f"Error: PDF file not found at '{pdf_path}'")
        return ""
    except pypdf.errors.PdfReadError:
        print(f"Error: Could not read PDF file '{pdf_path}'. It might be corrupted or encrypted.")
        return ""
    except Exception as e:
        print(f"An unexpected error occurred while processing '{pdf_path}': {e}")
        return ""

    full_text = "".join(extracted_text)

    if output_txt_path:
        try:
            with open(output_txt_path, 'w', encoding='utf-8') as outfile:
                outfile.write(full_text)
            print(f"Extracted text saved to '{output_txt_path}'")
        except IOError as e:
            print(f"Error: Could not write to output file '{output_txt_path}': {e}")
        except Exception as e:
            print(f"An unexpected error occurred while saving text: {e}")

    return full_text

if __name__ == "__main__":
    # Example usage:
    # Ensure you have a PDF file named 'example.pdf' in the same directory
    # or provide a full path to a PDF file.
    # For demonstration, let's assume 'example.pdf' exists.
    # If not, this will print a FileNotFoundError.

    pdf_file = "example.pdf" # Replace with your PDF file path

    # Extract all text and print it (first 500 characters)
    extracted_all = extract_text_from_pdf(pdf_file)
    if extracted_all:
        print(f"Extracted text (first 500 chars):\n{extracted_all[:500]}...")

    # Extract text from pages 0 to 1 (first two pages) and save to a file
    # and print the result.
    output_file = "output_pages_0_to_1.txt"
    extracted_partial = extract_text_from_pdf(pdf_file, start_page=0, end_page=2, output_txt_path=output_file)
    if extracted_partial:
        print(f"\nExtracted text from pages 0-1 (first 500 chars):\n{extracted_partial[:500]}...")

    # Extract text from a single page (e.g., page 2, 0-indexed) without page separator
    # and print the result.
    extracted_single_page = extract_text_from_pdf(pdf_file, start_page=2, end_page=3, page_separator="")
    if extracted_single_page:
        print(f"\nExtracted text from page 2 (first 500 chars):\n{extracted_single_page[:500]}...")

    # Example of handling a non-existent file
    extract_text_from_pdf("non_existent_file.pdf")

    # Example of handling a potentially invalid PDF (e.g., a text file renamed)
    # Create a dummy invalid file for testing
    dummy_invalid_pdf = "invalid_dummy.pdf"
    with open(dummy_invalid_pdf, "w") as f:
        f.write("This is not a valid PDF content.")
    extract_text_from_pdf(dummy_invalid_pdf)
    os.remove(dummy_invalid_pdf) # Clean up dummy file
