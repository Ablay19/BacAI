import sys
import PyPDF2

def extract_text_from_pdf(pdf_path):
    try:
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
    extract_text_from_pdf(pdf_path)