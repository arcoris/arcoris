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

// Done completes processing for item.
//
// Done removes clean processing items from tracking. Dirty processing items are
// requeued at the FIFO back unless the queue is shutting down. During shutdown,
// dirty processing items are removed instead of requeued. Done rejects invalid
// items before looking up processing state.
func (q *Queue) Done(item Item) error {
	if q == nil {
		return ErrInvalidQueue
	}
	if err := validateItem(item); err != nil {
		return err
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	id := keyForItem(item)
	e, ok := q.items[id]
	if !ok || e.state != stateProcessing {
		return ErrNotProcessing
	}

	if e.dirty && !q.shutDown {
		q.enqueueLocked(e)
		return nil
	}

	delete(q.items, id)
	q.signalNotFullLocked()
	return nil
}
