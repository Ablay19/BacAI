#ifndef QUEUE_H
#define QUEUE_H

#include <stdatomic.h>
#include <stdbool.h>
#include <stdlib.h>

typedef struct {
    atomic_size_t sequence;
    void* data;
} cell_t;

typedef struct {
    cell_t* buffer;
    size_t buffer_mask;
    char pad0[64];
    atomic_size_t enqueue_pos;
    char pad1[64];
    atomic_size_t dequeue_pos;
    char pad2[64];
} mpmc_queue_t;

mpmc_queue_t* mpmc_create(size_t buffer_size);
bool mpmc_enqueue(mpmc_queue_t* q, void* data);
bool mpmc_dequeue(mpmc_queue_t* q, void** data);
void mpmc_destroy(mpmc_queue_t* q);

#endif
