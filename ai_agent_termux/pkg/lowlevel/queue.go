package lowlevel

/*
#include "queue.h"
*/
import "C"
import "unsafe"

// MPMCQueue is a fast, lock-free Multi-Producer Multi-Consumer queue
type MPMCQueue struct {
	ptr *C.mpmc_queue_t
}

// NewMPMCQueue creates a new lock-free queue with a power-of-two size
func NewMPMCQueue(size int) *MPMCQueue {
	return &MPMCQueue{
		ptr: C.mpmc_create(C.size_t(size)),
	}
}

// Enqueue adds an item to the queue. Returns false if full.
func (q *MPMCQueue) Enqueue(data unsafe.Pointer) bool {
	return bool(C.mpmc_enqueue(q.ptr, data))
}

// Dequeue removes an item from the queue. Returns false if empty.
func (q *MPMCQueue) Dequeue() (unsafe.Pointer, bool) {
	var data unsafe.Pointer
	success := bool(C.mpmc_dequeue(q.ptr, &data))
	return data, success
}

// Destroy frees the queue memory
func (q *MPMCQueue) Destroy() {
	C.mpmc_destroy(q.ptr)
}
