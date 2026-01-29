#ifndef ARENA_H
#define ARENA_H

#include <stdlib.h>
#include <string.h>

/**
 * A simple arena allocator to reduce GC pressure and allocation overhead.
 */
typedef struct {
    char* buffer;
    size_t size;
    size_t offset;
} Arena;

Arena* arena_create(size_t size);
void* arena_alloc(Arena* arena, size_t size);
void arena_reset(Arena* arena);
void arena_destroy(Arena* arena);
char* arena_strdup(Arena* arena, const char* s);

#endif
