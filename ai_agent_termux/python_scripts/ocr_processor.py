import sys
import os
import subprocess
import json
from PIL import Image
import pytesseract

def get_scene_description(image_path):
    """Try to get a scene description using Ollama's LLaVA model if available."""
    try:
        # Check if ollama is available
        subprocess.run(["ollama", "--version"], capture_output=True, check=True)
        
        # Call ollama with llava model
        # We use a concise prompt for scene description
        prompt = "Describe this image in detail, focusing on the content, layout, and any visible objects or scenes."
        
        result = subprocess.run(
            ["ollama", "run", "llava", f"{prompt}", "--images", image_path],
            capture_output=True,
            text=True,
            check=True,
            timeout=60 # Timeout for local LLM
        )
        return result.stdout.strip()
    except Exception as e:
        # Silently fail and return empty if multimodal is not available
        return ""

def process_image(image_path, lang='eng'):
    try:
        # 1. Get Scene Description (Multimodal)
        scene_desc = get_scene_description(image_path)
        
        # 2. Get OCR text (Tesseract)
        img = Image.open(image_path)
        if img.mode != 'RGB':
            img = img.convert('RGB')
        ocr_text = pytesseract.image_to_string(img, lang=lang).strip()
        
        # Combine results
        final_output = []
        if scene_desc:
            final_output.append("--- Scene Description (LLaVA) ---")
            final_output.append(scene_desc)
            final_output.append("\n--- Extracted Text (OCR) ---")
        
        if ocr_text:
            final_output.append(ocr_text)
        elif not scene_desc:
            final_output.append("[No text or scene description could be extracted]")
            
        print("\n".join(final_output))
        
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