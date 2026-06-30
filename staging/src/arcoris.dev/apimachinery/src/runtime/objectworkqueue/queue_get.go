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

import "context"

// Get returns the oldest queued item and marks it processing.
//
// Get blocks while the queue is empty. After ShutDown, Get drains already
// queued items before returning ErrShutDown. A nil ctx is treated as
// context.Background.
func (q *Queue) Get(ctx context.Context) (Item, error) {
	if q == nil {
		return Item{}, ErrInvalidQueue
	}
	ctx = normalizeContext(ctx)

	for {
		q.mu.Lock()
		if item, ok := q.popLocked(); ok {
			q.mu.Unlock()
			return item, nil
		}
		if q.shutDown {
			q.mu.Unlock()
			return Item{}, ErrShutDown
		}
		ch := q.notEmpty
		q.notEmptyWaiters++
		q.mu.Unlock()

		if err := q.waitNotEmpty(ctx, ch); err != nil {
			return Item{}, err
		}
	}
}

// popLocked removes the oldest queued entry and marks it processing.
func (q *Queue) popLocked() (Item, bool) {
	e := q.head
	if e == nil {
		return Item{}, false
	}

	q.head = e.next
	if q.head == nil {
		q.tail = nil
	} else {
		q.head.prev = nil
	}
	q.queued--

	e.prev = nil
	e.next = nil
	e.state = stateProcessing

	return e.item, true
}
