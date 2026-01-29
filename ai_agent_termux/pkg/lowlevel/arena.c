#include "arena.h"

/**
 * A simple arena allocator to reduce GC pressure and allocation overhead.
 */

Arena* arena_create(size_t size) {
    Arena* arena = (Arena*)malloc(sizeof(Arena));
    if (!arena) return NULL;
    
    arena->buffer = (char*)malloc(size);
    if (!arena->buffer) {
        free(arena);
        return NULL;
    }
    
    arena->size = size;
    arena->offset = 0;
    return arena;
}

void* arena_alloc(Arena* arena, size_t size) {
    if (arena->offset + size > arena->size) {
        return NULL; // Out of memory in arena
    }
    
    void* ptr = arena->buffer + arena->offset;
    arena->offset += size;
    return ptr;
}

void arena_reset(Arena* arena) {
    arena->offset = 0;
}

void arena_destroy(Arena* arena) {
    if (!arena) return;
    free(arena->buffer);
    free(arena);
}

// Helper for fast string duplication into arena
char* arena_strdup(Arena* arena, const char* s) {
    size_t len = strlen(s) + 1;
    char* p = (char*)arena_alloc(arena, len);
    if (p) {
        memcpy(p, s, len);
    }
    return p;
}
