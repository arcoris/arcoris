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

// Add records item as pending work, blocking while capacity is full.
//
// Add coalesces duplicate queued items. If item is processing, Add marks it
// dirty so Done will requeue it. Add rejects invalid items before waiting for
// capacity. A nil ctx is treated as context.Background.
func (q *Queue) Add(ctx context.Context, item Item) error {
	if q == nil {
		return ErrInvalidQueue
	}
	ctx = normalizeContext(ctx)
	if err := validateItem(item); err != nil {
		return err
	}

	for {
		q.mu.Lock()
		err, done := q.addLocked(item, true)
		if done {
			q.mu.Unlock()
			return err
		}
		ch := q.notFull
		q.mu.Unlock()

		select {
		case <-ch:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// TryAdd records item as pending work without blocking.
//
// TryAdd returns ErrFull only when item is new and the queue has no remaining
// capacity. Duplicate tracked items still succeed because they do not increase
// tracked cardinality. TryAdd rejects invalid items before inspecting queue
// capacity.
func (q *Queue) TryAdd(item Item) error {
	if q == nil {
		return ErrInvalidQueue
	}
	if err := validateItem(item); err != nil {
		return err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	err, done := q.addLocked(item, false)
	if !done {
		return ErrFull
	}
	return err
}

// addLocked applies Add and TryAdd transitions while q.mu is held.
//
// The returned bool reports whether the caller can stop immediately. Add uses a
// false value to wait for capacity on q.notFull; TryAdd converts it to ErrFull.
func (q *Queue) addLocked(item Item, wait bool) (error, bool) {
	if q.shutDown {
		return ErrShutDown, true
	}

	id := keyForItem(item)
	if e, ok := q.items[id]; ok {
		if e.state == stateProcessing {
			e.dirty = true
		}
		return nil, true
	}

	if q.trackedLocked() >= q.capacity {
		if wait {
			return nil, false
		}
		return ErrFull, true
	}

	e := &entry{item: item}
	q.items[id] = e
	q.enqueueLocked(e)
	return nil, true
}
