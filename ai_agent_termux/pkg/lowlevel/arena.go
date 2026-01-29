package lowlevel

/*
#include "arena.h"
*/
import "C"
import (
	"unsafe"
)

type Arena struct {
	ptr *C.Arena
}

func NewArena(size int) *Arena {
	p := C.arena_create(C.size_t(size))
	if p == nil {
		return nil
	}
	return &Arena{ptr: p}
}

func (a *Arena) Alloc(size int) unsafe.Pointer {
	return C.arena_alloc(a.ptr, C.size_t(size))
}

func (a *Arena) Reset() {
	C.arena_reset(a.ptr)
}

func (a *Arena) Destroy() {
	C.arena_destroy(a.ptr)
}

func (a *Arena) StrDup(s string) string {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	p := C.arena_strdup(a.ptr, cs)
	if p == nil {
		return ""
	}
	return C.GoString(p)
}
