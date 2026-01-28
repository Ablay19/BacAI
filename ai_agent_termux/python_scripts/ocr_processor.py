import sys
from PIL import Image
import pytesseract

def process_image(image_path, lang='eng'):
    try:
        img = Image.open(image_path)
        # Ensure image is RGB for Tesseract if it's not already
        if img.mode != 'RGB':
            img = img.convert('RGB')
        text = pytesseract.image_to_string(img, lang=lang)
        print(text.strip())
    except Exception as e:
        print(f"Error processing image {image_path}: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python ocr_processor.py <image_path> [language]", file=sys.stderr)
        sys.exit(1)
    
    image_path = sys.argv[1]
    lang = sys.argv[2] if len(sys.argv) > 2 else 'eng'
    process_image(image_path, lang)