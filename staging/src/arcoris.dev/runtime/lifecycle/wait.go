/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package lifecycle

import "context"

// Predicate evaluates whether a lifecycle snapshot satisfies a wait condition.
//
// Predicate receives a Snapshot rather than only State because lifecycle waiters
// often need revision, last transition, terminal failure cause, or other
// read-side context. Predicates MUST be fast and side-effect free. They MUST NOT
// call transition methods on the same Controller.
//
// Predicate is evaluated synchronously by Wait. Expensive checks, blocking I/O,
// retries, sleeps, logging, metrics export, and external synchronization do not
// belong in predicates.
type Predicate func(Snapshot) bool

// Wait blocks until predicate accepts a lifecycle snapshot, the lifecycle reaches
// a terminal state before the predicate is satisfied, or ctx is cancelled.
//
// Wait returns the Snapshot that satisfied predicate on success. On failure it
// returns the latest Snapshot observed by the wait loop together with a WaitError.
//
// Wait treats terminal state as a natural boundary. If the lifecycle reaches
// StateStopped or StateFailed and predicate is still false, the predicate can no
// longer become true because terminal states have no outgoing transitions. In
// that case Wait returns ErrWaitTargetUnreachable wrapped in WaitError.
//
// A nil predicate is invalid and returns ErrInvalidWaitPredicate. A nil context
// is treated as context.Background. This avoids panics in defensive paths while
// preserving the usual Go convention that callers should pass a real context
// when they need cancellation or deadlines.
//
// Wait is safe to call concurrently with transition methods.
func (c *Controller) Wait(ctx context.Context, predicate Predicate) (Snapshot, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	snap, changed, done := c.waitSnapshot()

	if predicate == nil {
		return snap, newWaitError(snap, ErrInvalidWaitPredicate)
	}

	for {
		if predicate(snap) {
			return snap, nil
		}

		if snap.IsTerminal() {
			return snap, newWaitError(snap, ErrWaitTargetUnreachable)
		}

		select {
		case <-changed:
			snap, changed, done = c.waitSnapshot()

		case <-done:
			snap, _, _ = c.waitSnapshot()
			if predicate(snap) {
				return snap, nil
			}

			return snap, newWaitError(snap, ErrWaitTargetUnreachable)

		case <-ctx.Done():
			snap, _, _ = c.waitSnapshot()
			return snap, newWaitError(snap, ctx.Err())
		}
	}
}
