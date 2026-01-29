package lowlevel

import "os"

// MmapFile fallback for Windows - simple file read
func MmapFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Munmap fallback for Windows - do nothing
func Munmap(data []byte) error {
	return nil
}
