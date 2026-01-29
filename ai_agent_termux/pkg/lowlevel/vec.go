//go:build arm64

package lowlevel

// NormalizeVec64 uses ARM64 NEON instructions to normalize a float64 vector in-place.
func NormalizeVec64(v []float64)

// L2Norm64 computes the L2 norm of a float64 vector using SIMD.
func L2Norm64(v []float64) float64
