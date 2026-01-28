import sys
import json
import numpy as np

try:
    from sentence_transformers import SentenceTransformer
except ImportError:
    print("sentence-transformers library not found. Please install it using: pip install sentence-transformers", file=sys.stderr)
    sys.exit(1)

def generate_embeddings(text):
    try:
        model = SentenceTransformer('all-MiniLM-L6-v2')
        embedding = model.encode(text)
        return embedding.tolist()  # Convert to list for JSON serialization
    except Exception as e:
        print(f"Error generating embeddings: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python embedding_generator.py <text>", file=sys.stderr)
        sys.exit(1)
    
    text = sys.argv[1]
    embeddings = generate_embeddings(text)
    print(json.dumps(embeddings))