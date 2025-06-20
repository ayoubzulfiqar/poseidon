def format_paragraphs(text, width):
    paragraphs = text.split('\n\n')
    formatted_paragraphs = []
    for para in paragraphs:
        words = para.split()
        lines = []
        current_line_words = []
        for word in words:
            potential_line = " ".join(current_line_words + [word])
            if len(potential_line) > width and current_line_words:
                lines.append(" ".join(current_line_words))
                current_line_words = [word]
            else:
                current_line_words.append(word)
        if current_line_words:
            lines.append(" ".join(current_line_words))
        formatted_paragraphs.append("\n".join(lines))
    return "\n\n".join(formatted_paragraphs)

# Additional implementation at 2025-06-20 00:05:09


# Additional implementation at 2025-06-20 00:06:32
import re

class ParagraphFormatter:
    def __init__(self, width=80, indent=0, first_line_indent=0, alignment='left'):
        if not isinstance(width, int) or width <= 0:
            raise ValueError("Width must be a positive integer.")
        if not isinstance(indent, int) or indent < 0:
            raise ValueError("Indent must be a non-negative integer.")
        if not isinstance(first_line_indent, int) or first_line_indent < 0:
            raise ValueError("First line indent must be a non-negative integer.")
        if alignment not in ['left', 'right', 'center', 'justify']:
            raise ValueError("Alignment must be 'left', 'right', 'center', or 'justify'.")

        self.width = width
        self.indent = indent
        self.first_line_indent = first_line_indent
        self.alignment = alignment

    def _wrap_text_with_varying_indents(self, text):
        """
        Wraps text into lines, considering first_line_indent and subsequent indent.
        Returns a list of (line_content, is_first_line_of_paragraph) tuples.
        """
        words = text.split()
        if not words:
            return []

        lines_info = []
        current_line_words = []
        current_line_length = 0
        is_first_line_of_paragraph = True

        for word in words:
            current_line_effective_width = (self.width - self.first_line_indent) if is_first_line_of_paragraph else (self.width - self.indent)
            current_line_effective_width = max(0, current_line_effective_width) # Ensure non-negative width

            word_length = len(word)
            space_needed = 1 if current_line_words else 0

            # If the word itself is longer than the effective width, it will occupy its own line
            # and potentially exceed the width. This is standard text wrapping behavior.
            if current_line_length + space_needed + word_length <= current_line_effective_width:
                current_line_words.append(word)
                current_line_length += space_needed + word_length
            else:
                # Current word doesn't fit, so finalize the current line
                if current_line_words: # Only add if there's content
                    lines_info.append((" ".join(current_line_words), is_first_line_of_paragraph))
                
                is_first_line_of_paragraph = False # Subsequent lines won't be the first line of the paragraph
                current_line_words = [word]
                current_line_length = word_length

        # Add the last accumulated line
        if current_line_words:
            lines_info.append((" ".join(current_line_words), is_first_line_of_paragraph))

        return lines_info

    def _align_line(self, line_content, target_width, alignment):
        """
        Applies alignment to a single line of text.
        """
        if not line_content:
            return ""

        current_len = len(line_content)
        if current_len >= target_width:
            return line_content # Already at or over width, no padding needed

        if alignment == 'left':
            return line_content.ljust(target_width)
        elif alignment == 'right':
            return line_content.rjust(target_width)
        elif alignment == 'center':
            return line_content.center(target_width)
        elif alignment == 'justify':
            words = line_content.split()
            if len(words) <= 1: # Cannot justify a single word or empty line
                return line_content.ljust(target_width)

            num_spaces_to_add = target_width - current_len
            num_gaps = len(words) - 1

            if num_gaps == 0: # Should be handled by len(words) <= 1, but for safety
                return line_content.ljust(target_width)

            base_spaces_per_gap = num_spaces_to_add // num_gaps
            extra_spaces = num_spaces_to_add % num_gaps

            justified_line = []
            for i, word in enumerate(words):
                justified_line.append(word)
                if i < num_gaps:
                    spaces = base_spaces_per_gap + (1 if i < extra_spaces else 0)
                    justified_line.append(" " * (1 + spaces)) # 1 for original space, + spaces for added
            return "".join(justified_line)
        else:
            return line_content.ljust(target_width) # Fallback

    def format_paragraph(self, paragraph_text):
        """
        Formats a single paragraph string according to the formatter's settings.
        """
        # Normalize whitespace: replace multiple spaces with single, strip leading/trailing
        normalized_text = " ".join(paragraph_text.split()).strip()
        if not normalized_text:
            return ""

        # Get wrapped lines with info about whether they are the first line of the paragraph
        wrapped_lines_info = self._wrap_text_with_varying_indents(normalized_text)

        formatted_lines = []
        for line_content, is_first_line_of_paragraph in wrapped_lines_info:
            current_indent_spaces = self.first_line_indent if is_first_line_of_paragraph else self.indent
            current_effective_width = self.width - current_indent_spaces
            current_effective_width = max(0, current_effective_width) # Ensure non-negative

            # Apply alignment
            aligned_line = self._align_line(line_content, current_effective_width, self.alignment)

            # Apply indentation
            indented_line = " " * current_indent_spaces + aligned_line
            formatted_lines.append(indented_line)

        return "\n".join(formatted_lines)

    def format_document(self, document_text):
        """
        Formats a multi-paragraph document string.
        Paragraphs are separated by one or more blank lines.
        """
        # Split document into paragraphs. Use regex to handle multiple newlines as delimiters.
        paragraphs = re.split(r'\n\s*\n+', document_text.strip())

        formatted_paragraphs = []
        for para in paragraphs:
            if para.strip(): # Ensure it's not just empty string from split
                formatted_paragraphs.append(self.format_paragraph(para))

        return "\n\n".join(formatted_paragraphs)

if __name__ == '__main__':
    # Example Usage:

    sample_text = """
    This is a sample paragraph that demonstrates the functionality of the paragraph formatter. It should be wrapped to a fixed width. We can also test how it handles very long words or short lines.

    This is another paragraph. It is separated from the first by a blank line. We will see how the formatter handles multiple paragraphs in a single document.

    A third, shorter paragraph.
    """

    print("--- Default Formatting (width=80, left align, no indent) ---")
    formatter_default = ParagraphFormatter()
    print(formatter_default.format_document(sample_text))
    print("\n" + "="*80 + "\n")

    print("--- Width 40, Left Align, Indent 4, First Line Indent 8 ---")
    formatter_indented = ParagraphFormatter(width=40, indent=4, first_line_indent=8, alignment='left')
    print(formatter_indented.format_document(sample_text))
    print("\n" + "="*80 + "\n")

    print("--- Width 50, Right Align, Indent 2 ---")
    formatter_right = ParagraphFormatter(width=50, indent=2, alignment='right')
    print(formatter_right.format_document(sample_text))
    print("\n" + "="*80 + "\n")

    print("--- Width 60, Center Align, No Indent ---")
    formatter_center = ParagraphFormatter(width=60, alignment='center')
    print(formatter_center.format_document(sample_text))
    print("\n" + "="*80 + "\n")

    print("--- Width 70, Justify Align, Indent 3 ---")
    formatter_justify = ParagraphFormatter(width=70, indent=3, alignment='justify')
    print(formatter_justify.format_document(sample_text))
    print("\n" + "="*80 + "\n")

    print("--- Test with very long word and narrow width ---")
    long_word_text = "This is a paragraph with an extremelylongwordthatwillnotfitonasingleline and some other words."
    formatter_narrow = ParagraphFormatter(width=20, indent=0, alignment='left')
    print(formatter_narrow.format_document(long_word_text))
    print("\n" + "="*80 + "\n")

    print("--- Test with empty input ---")
    print(f"'{formatter_default.format_document('')}'")
    print("\n" + "="*80 + "\n")

    print("--- Test with only spaces/newlines ---")
    print(f"'{formatter_default.format_document('   \n\n  \t ')}'")
    print("\n" + "="*80 + "\n")

    print("--- Test with width less than indent (should still work) ---")
    formatter_extreme_indent = ParagraphFormatter(width=10, indent=15, first_line_indent=5, alignment='left')
    print(formatter_extreme_indent.format_document("Short text to test extreme indent."))
    print("\n" + "="*80 + "\n")