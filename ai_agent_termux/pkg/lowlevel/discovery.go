//go:build !windows

package lowlevel

import (
	"syscall"
)

// Linux dirent64 structure
type dirent struct {
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]int8
}

// FastScanDir uses getdents64 syscall to read directory entries in bulk.
// This is significantly faster than filepath.Walk for very large directories.
func FastScanDir(path string) ([]string, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	var names []string
	buf := make([]byte, 32*1024) // 32KB buffer for entries

	for {
		n, err := syscall.ReadDirent(fd, buf)
		if err != nil {
			return nil, err
		}
		if n <= 0 {
			break
		}

		// Parse the buffer
		var consumed int = 0
		for consumed < n {
			// Each entry has a variable length, but it's aligned
			// We use low-level parsing to extract the names
			// This is roughly what the 'os' package does but specialized

			// Note: syscall.ReadDirent returns names directly in some versions/OS
			// but for maximum speed on Linux/ARM (Termux), we want bulk reading.

			// For Go's syscall package specifically:
			nb, count, names_chunk := syscall.ParseDirent(buf[consumed:n], 1024, names)
			_ = nb
			consumed += count
			names = names_chunk
		}
	}

	return names, nil
}
