package lowlevel

/*
#include <stdlib.h>
#include <string.h>

typedef struct {
    char* buf;
    size_t size;
    size_t offset;
} TUIArena;

TUIArena* tui_arena_create(size_t size) {
    TUIArena* a = (TUIArena*)malloc(sizeof(TUIArena));
    a->buf = (char*)malloc(size);
    a->size = size;
    a->offset = 0;
    return a;
}

void* tui_arena_alloc(TUIArena* a, size_t size) {
    if (a->offset + size > a->size) return NULL;
    void* ptr = a->buf + a->offset;
    a->offset += size;
    return ptr;
}

void tui_arena_reset(TUIArena* a) {
    a->offset = 0;
}

void tui_arena_destroy(TUIArena* a) {
    free(a->buf);
    free(a);
}
*/
import "C"
import "unsafe"

// ArenaTUI is a high-speed memory pool for TUI rendering
type ArenaTUI struct {
	ptr *C.TUIArena
}

// NewArenaTUI creates a 1MB arena for UI buffers
func NewArenaTUI() *ArenaTUI {
	return &ArenaTUI{
		ptr: C.tui_arena_create(1024 * 1024),
	}
}

// Alloc allocates raw memory from the arena
func (a *ArenaTUI) Alloc(size int) unsafe.Pointer {
	return C.tui_arena_alloc(a.ptr, C.size_t(size))
}

// Reset clears the arena for reuse
func (a *ArenaTUI) Reset() {
	C.tui_arena_reset(a.ptr)
}

// Destroy frees all allocated memory
func (a *ArenaTUI) Destroy() {
	C.tui_arena_destroy(a.ptr)
}
