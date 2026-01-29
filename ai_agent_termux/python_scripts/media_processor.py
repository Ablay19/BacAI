import sys
import subprocess
import json

def extract_metadata(file_path):
    try:
        # Run ffprobe to get metadata in JSON format
        cmd = [
            'ffprobe', 
            '-v', 'quiet', 
            '-print_format', 'json', 
            '-show_format', 
            '-show_streams', 
            file_path
        ]
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        metadata = json.loads(result.stdout)
        
        # Simplify and format output for LLM consumption
        fmt = metadata.get('format', {})
        streams = metadata.get('streams', [])
        
        output = [
            f"--- Media Metadata for {file_path} ---",
            f"Duration: {fmt.get('duration', 'unknown')}s",
            f"Size: {fmt.get('size', 'unknown')} bytes",
            f"Format: {fmt.get('format_long_name', 'unknown')}",
            f"Bitrate: {fmt.get('bit_rate', 'unknown')} bps"
        ]
        
        for i, stream in enumerate(streams):
            output.append(f"\nStream #{i} ({stream.get('codec_type')}):")
            output.append(f"  Codec: {stream.get('codec_name')}")
            if stream.get('codec_type') == 'video':
                output.append(f"  Resolution: {stream.get('width')}x{stream.get('height')}")
                output.append(f"  Frame Rate: {stream.get('avg_frame_rate')}")
            elif stream.get('codec_type') == 'audio':
                output.append(f"  Channels: {stream.get('channels')}")
                output.append(f"  Sample Rate: {stream.get('sample_rate')} Hz")
        
        print("\n".join(output))
        
    except Exception as e:
        print(f"Error processing media {file_path}: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python media_processor.py <file_path>", file=sys.stderr)
        sys.exit(1)
    
    file_path = sys.argv[1]
    extract_metadata(file_path)
