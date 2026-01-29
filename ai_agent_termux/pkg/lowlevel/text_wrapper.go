package lowlevel

/*
#include "text_utils.h"
*/
import "C"

import (
	"unsafe"
)

// SanitizeTextC cleans up a string using the high-performance C engine.
// This is done in-place if possible (internally) but Go strings are immutable,
// so we typically use this with the Arena allocator.
func SanitizeTextC(arena *Arena, input string) string {
	if arena == nil {
		return input // Fallback
	}

	// Duplicate into arena
	cs := C.CString(input)
	// We don't defer free because we want to use the arena
	// Wait, C.CString uses malloc. If we want total zero-copy,
	// we should copy directly into the arena first.

	len := len(input)
	ptr := arena.Alloc(len + 1)
	if ptr == nil {
		C.free(unsafe.Pointer(cs))
		return input
	}

	// Copy input to arena
	C.memcpy(ptr, unsafe.Pointer(cs), C.size_t(len+1))
	C.free(unsafe.Pointer(cs))

	C.sanitize_text_in_place((*C.char)(ptr))

	return C.GoString((*C.char)(ptr))
}
