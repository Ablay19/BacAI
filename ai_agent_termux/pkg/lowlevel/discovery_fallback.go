package lowlevel

import (
	"os"
)

// FastScanDir fallback for Windows - use standard os.ReadDir
func FastScanDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}
