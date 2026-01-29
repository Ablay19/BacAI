//go:build arm64
// +build arm64

package lowlevel

// FastHashASM computes a fast non-cryptographic hash using ARM64 NEON instructions.
// This is used for quick file change detection.
func FastHashASM(data []byte) uint64
