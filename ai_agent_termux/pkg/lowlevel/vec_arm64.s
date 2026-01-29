// +build arm64

#include "textflag.h"

// func L2Norm64(v []float64) float64
TEXT ·L2Norm64(SB), NOSPLIT, $0-32
    MOVD v_base+0(FP), R0
    MOVD v_len+8(FP), R1
    
    FMOVD $(0.0), F0 // Sum accumulator
    
    CBZ R1, done

loop:
    FMOVD (R0), F1
    FMULD F1, F1, F2 // F2 = val * val
    FADDD F2, F0     // sum += val*val
    ADD $8, R0
    SUB $1, R1
    CBNZ R1, loop

done:
    FSQRTD F0, F0    // result = sqrt(sum)
    FMOVD F0, ret+24(FP)
    RET

// func NormalizeVec64(v []float64)
TEXT ·NormalizeVec64(SB), NOSPLIT, $0-24
    MOVD v_base+0(FP), R0
    MOVD v_len+8(FP), R1
    MOVD R1, R2      // Save length for second loop
    MOVD R0, R3      // Save base for second loop
    
    FMOVD $(0.0), F0 // Sum accumulator
    
    CBZ R1, finished

    // First pass: calculate sum of squares
calc_sum:
    FMOVD (R0), F1
    FMULD F1, F1, F2
    FADDD F2, F0
    ADD $8, R0
    SUB $1, R1
    CBNZ R1, calc_sum

    FSQRTD F0, F0    // F0 = norm
    
    // If norm is effectively zero, exit
    FMOVD $(1e-9), F3
    FCMPD F3, F0
    BMI finished     // if F0 < 1e-9 return

    // Second pass: divide all elements by norm
    // F0 contains norm
norm_loop:
    FMOVD (R3), F1
    FDIVD F0, F1
    FMOVD F1, (R3)
    ADD $8, R3
    SUB $1, R2
    CBNZ R2, norm_loop

finished:
    RET
