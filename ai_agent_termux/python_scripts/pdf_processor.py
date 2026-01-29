import sys
import PyPDF2

try:
    import pdfplumber
    HAS_PDFPLUMBER = True
except ImportError:
    HAS_PDFPLUMBER = False

def extract_text_with_tables(pdf_path):
    try:
        if HAS_PDFPLUMBER:
            full_text = []
            with pdfplumber.open(pdf_path) as pdf:
                for page in pdf.pages:
                    # Extract plain text
                    page_text = page.extract_text()
                    if page_text:
                        full_text.append(page_text)
                    
                    # Extract tables and format them
                    tables = page.extract_tables()
                    if tables:
                        for table in tables:
                            full_text.append("\n--- Table Detected ---\n")
                            for row in table:
                                # Clean up row items
                                clean_row = [str(item).replace('\n', ' ').strip() if item else "" for item in row]
                                full_text.append(" | ".join(clean_row))
                            full_text.append("----------------------\n")
            print("\n".join(full_text))
        else:
            # Fallback to PyPDF2
            with open(pdf_path, 'rb') as file:
                reader = PyPDF2.PdfReader(file)
                text = ""
                for page in reader.pages:
                    text += page.extract_text() + "\n"
                print(text.strip())
    except Exception as e:
        print(f"Error processing PDF {pdf_path}: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python pdf_processor.py <pdf_path>", file=sys.stderr)
        sys.exit(1)
    
    pdf_path = sys.argv[1]
    extract_text_with_tables(pdf_path)