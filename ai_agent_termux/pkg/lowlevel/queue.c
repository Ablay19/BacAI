#include <stdatomic.h>
#include <stdbool.h>
#include <stdlib.h>

// Lock-free MPMC Ring Buffer
// Based on the 1024cores implementation

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

mpmc_queue_t* mpmc_create(size_t buffer_size) {
    // buffer_size must be a power of 2
    if ((buffer_size & (buffer_size - 1)) != 0) return NULL;

    mpmc_queue_t* q = (mpmc_queue_t*)malloc(sizeof(mpmc_queue_t));
    q->buffer = (cell_t*)malloc(sizeof(cell_t) * buffer_size);
    q->buffer_mask = buffer_size - 1;
    for (size_t i = 0; i != buffer_size; i += 1) {
        atomic_init(&q->buffer[i].sequence, i);
    }
    atomic_init(&q->enqueue_pos, 0);
    atomic_init(&q->dequeue_pos, 0);
    return q;
}

bool mpmc_enqueue(mpmc_queue_t* q, void* data) {
    cell_t* cell;
    size_t pos = atomic_load_explicit(&q->enqueue_pos, memory_order_relaxed);
    for (;;) {
        cell = &q->buffer[pos & q->buffer_mask];
        size_t seq = atomic_load_explicit(&cell->sequence, memory_order_acquire);
        intptr_t dif = (intptr_t)seq - (intptr_t)pos;
        if (dif == 0) {
            if (atomic_compare_exchange_weak_explicit(&q->enqueue_pos, &pos, pos + 1, memory_order_relaxed, memory_order_relaxed)) {
                break;
            }
        } else if (dif < 0) {
            return false; // Queue full
        } else {
            pos = atomic_load_explicit(&q->enqueue_pos, memory_order_relaxed);
        }
    }
    cell->data = data;
    atomic_store_explicit(&cell->sequence, pos + 1, memory_order_release);
    return true;
}

bool mpmc_dequeue(mpmc_queue_t* q, void** data) {
    cell_t* cell;
    size_t pos = atomic_load_explicit(&q->dequeue_pos, memory_order_relaxed);
    for (;;) {
        cell = &q->buffer[pos & q->buffer_mask];
        size_t seq = atomic_load_explicit(&cell->sequence, memory_order_acquire);
        intptr_t dif = (intptr_t)seq - (intptr_t)(pos + 1);
        if (dif == 0) {
            if (atomic_compare_exchange_weak_explicit(&q->dequeue_pos, &pos, pos + 1, memory_order_relaxed, memory_order_relaxed)) {
                break;
            }
        } else if (dif < 0) {
            return false; // Queue empty
        } else {
            pos = atomic_load_explicit(&q->dequeue_pos, memory_order_relaxed);
        }
    }
    *data = cell->data;
    atomic_store_explicit(&cell->sequence, pos + q->buffer_mask + 1, memory_order_release);
    return true;
}

void mpmc_destroy(mpmc_queue_t* q) {
    free(q->buffer);
    free(q);
}
