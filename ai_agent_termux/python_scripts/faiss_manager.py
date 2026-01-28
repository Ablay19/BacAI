import sys
import json
import numpy as np

try:
    import faiss
except ImportError:
    print("faiss library not found. Please install it using: pip install faiss-cpu", file=sys.stderr)
    sys.exit(1)

class FaissManager:
    def __init__(self, dimension, index_path="faiss_index.bin"):
        self.dimension = dimension
        self.index_path = index_path
        self.index = faiss.IndexFlatL2(dimension)
        
    def add_embedding(self, embedding, metadata):
        embedding_np = np.array(embedding).astype('float32').reshape(1, -1)
        self.index.add(embedding_np)
        
    def save_index(self):
        faiss.write_index(self.index, self.index_path)
        
    def load_index(self):
        self.index = faiss.read_index(self.index_path)
        
    def search(self, query_embedding, k=5):
        query_np = np.array(query_embedding).astype('float32').reshape(1, -1)
        distances, indices = self.index.search(query_np, k)
        return distances, indices

def main():
    if len(sys.argv) < 3:
        print("Usage: python faiss_manager.py <dimension> <operation> [args...]", file=sys.stderr)
        print("Operations: add, save, load, search", file=sys.stderr)
        sys.exit(1)
        
    dimension = int(sys.argv[1])
    operation = sys.argv[2]
    
    fm = FaissManager(dimension)
    
    if operation == "add":
        if len(sys.argv) < 5:
            print("Usage: python faiss_manager.py <dim> add <embedding_json> <metadata_json>", file=sys.stderr)
            sys.exit(1)
            
        embedding = json.loads(sys.argv[3])
        metadata = json.loads(sys.argv[4])
        fm.add_embedding(embedding, metadata)
        print("Embedding added successfully")
        
    elif operation == "save":
        fm.save_index()
        print("Index saved successfully")
        
    elif operation == "load":
        fm.load_index()
        print("Index loaded successfully")
        
    elif operation == "search":
        if len(sys.argv) < 4:
            print("Usage: python faiss_manager.py <dim> search <query_embedding_json>", file=sys.stderr)
            sys.exit(1)
            
        query_embedding = json.loads(sys.argv[3])
        distances, indices = fm.search(query_embedding)
        result = {
            "distances": distances.tolist(),
            "indices": indices.tolist()
        }
        print(json.dumps(result))

if __name__ == "__main__":
    main()