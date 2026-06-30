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

// signalNotEmptyLocked broadcasts that queued work may now be available.
//
// The queue uses close-and-replace channels as condition notifications so Add
// and Get can wait with caller contexts without spawning cancellation bridge
// goroutines. Channels are replaced only while waiters exist, keeping hot paths
// allocation-light when no caller is blocked. Higher layers should still bound
// their worker counts because every signal wakes all waiters observing the
// previous channel.
func (q *Queue) signalNotEmptyLocked() {
	if q.notEmptyWaiters == 0 {
		return
	}
	close(q.notEmpty)
	q.notEmpty = make(chan struct{})
}

// signalNotFullLocked broadcasts that tracked capacity may now be available.
func (q *Queue) signalNotFullLocked() {
	if q.notFullWaiters == 0 {
		return
	}
	close(q.notFull)
	q.notFull = make(chan struct{})
}

// signalAllLocked broadcasts both item-availability and capacity changes.
func (q *Queue) signalAllLocked() {
	q.signalNotEmptyLocked()
	q.signalNotFullLocked()
}

// waitNotEmpty waits for a not-empty signal or context cancellation.
func (q *Queue) waitNotEmpty(ctx context.Context, ch <-chan struct{}) error {
	err := wait(ctx, ch)
	q.mu.Lock()
	q.notEmptyWaiters--
	q.mu.Unlock()
	return err
}

// waitNotFull waits for a not-full signal or context cancellation.
func (q *Queue) waitNotFull(ctx context.Context, ch <-chan struct{}) error {
	err := wait(ctx, ch)
	q.mu.Lock()
	q.notFullWaiters--
	q.mu.Unlock()
	return err
}

// wait observes one notification channel with caller-owned cancellation.
func wait(ctx context.Context, ch <-chan struct{}) error {
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
