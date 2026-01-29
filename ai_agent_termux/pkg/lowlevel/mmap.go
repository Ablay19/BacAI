//go:build !windows

package lowlevel

import (
	"os"
	"syscall"
)

// MmapFile maps a file into memory for zero-copy reading.
// This is significantly faster for large files on Android/Linux.
func MmapFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := info.Size()
	if size == 0 {
		return []byte{}, nil
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Munmap unmaps the memory-mapped file.
func Munmap(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return syscall.Munmap(data)
}
