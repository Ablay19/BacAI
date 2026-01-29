// +build !arm64

package lowlevel

import (
	"hash/fnv"
)

// FastHashASM computes a fast non-cryptographic hash.
// Fallback for non-arm64 architectures.
func FastHashASM(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}
