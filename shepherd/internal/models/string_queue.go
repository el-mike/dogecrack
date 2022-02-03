package models

import "sync"

// QueueError - reusable error structure for queue errors.
type QueueError struct {
	Message string
}

// Error - Error interface implementation.
func (e *QueueError) Error() string {
	return e.Message
}

// NewStringQueueEmptyError - returns new QueueError instance.
func NewStringQueueEmptyError() *QueueError {
	return &QueueError{
		Message: "Queue is empty!",
	}
}

// StringQueue - thread-safe queue data structure for strings.
type StringQueue struct {
	sync.RWMutex

	items []string
}

// NewStringQueue - returns new StringsQueue.
func NewStringQueue() *StringQueue {
	return &StringQueue{
		items: []string{},
	}
}

// Enqueue - adds an item to the end of the queue.
func (qs *StringQueue) Enqueue(item string) {
	qs.Lock()
	defer qs.Unlock()

	qs.items = append(qs.items, item)
}

// Dequeue - removes and returns an element from the front of the queue.
func (qs *StringQueue) Dequeue() (string, error) {
	qs.Lock()
	defer qs.Unlock()

	if len(qs.items) == 0 {
		return "", NewStringQueueEmptyError()
	}

	item := qs.items[0]
	qs.items = qs.items[1:]

	return item, nil
}

// IsEmpty - returns true if queue has no items.
func (qs *StringQueue) IsEmpty() bool {
	qs.RLock()
	defer qs.RUnlock()

	return len(qs.items) == 0
}

// Peek - returns the first element of the queue, without removing it.
func (qs *StringQueue) Peek() (string, error) {
	qs.RLock()
	defer qs.RUnlock()

	if len(qs.items) == 0 {
		return "", NewStringQueueEmptyError()
	}

	return qs.items[0], nil
}
