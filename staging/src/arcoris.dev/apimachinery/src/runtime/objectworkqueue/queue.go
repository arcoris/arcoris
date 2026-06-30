// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectworkqueue

import "sync"

// Queue is a bounded, deduplicating, object-keyed work queue.
//
// A Queue must not be copied after first use.
type Queue struct {
	// noCopy lets go vet report accidental Queue copies after first use.
	noCopy noCopy

	// mu guards all mutable queue state below.
	mu sync.Mutex

	// notEmpty is closed and replaced whenever a queued item may be available.
	notEmpty chan struct{}

	// notFull is closed and replaced whenever tracked capacity may be available.
	notFull chan struct{}

	// notEmptyWaiters counts callers currently waiting on notEmpty.
	notEmptyWaiters int

	// notFullWaiters counts callers currently waiting on notFull.
	notFullWaiters int

	// capacity is the maximum number of distinct tracked items.
	capacity int

	// shutDown records that the queue no longer accepts or schedules new work.
	shutDown bool

	// head is the oldest queued entry in the intrusive FIFO.
	head *entry

	// tail is the newest queued entry in the intrusive FIFO.
	tail *entry

	// queued is the number of entries linked into the FIFO.
	queued int

	// items stores every tracked entry, including queued and processing items.
	items map[keyID]*entry
}

// New validates opts and returns an empty Queue.
func New(opts Options) (*Queue, error) {
	if opts.Capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &Queue{
		notEmpty: make(chan struct{}),
		notFull:  make(chan struct{}),
		capacity: opts.Capacity,
		items:    make(map[keyID]*entry, opts.Capacity),
	}, nil
}

// Len returns the number of queued items.
//
// Processing items count toward capacity but are not included in Len. A nil
// Queue returns zero.
func (q *Queue) Len() int {
	if q == nil {
		return 0
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queuedLocked()
}

// IsShutDown reports whether ShutDown has been called.
//
// A nil Queue reports true.
func (q *Queue) IsShutDown() bool {
	if q == nil {
		return true
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	return q.shutDown
}

// trackedLocked returns the number of queued plus processing items.
func (q *Queue) trackedLocked() int {
	return len(q.items)
}

// queuedLocked returns the number of FIFO-waiting items.
func (q *Queue) queuedLocked() int {
	return q.queued
}

// processingLocked returns the number of in-flight items.
func (q *Queue) processingLocked() int {
	return q.trackedLocked() - q.queuedLocked()
}

// enqueueLocked appends e to the FIFO queue and wakes Get waiters.
func (q *Queue) enqueueLocked(e *entry) {
	e.state = stateQueued
	e.dirty = false
	e.prev = q.tail
	e.next = nil
	if q.tail == nil {
		q.head = e
	} else {
		q.tail.next = e
	}
	q.tail = e
	q.queued++
	q.signalNotEmptyLocked()
}
