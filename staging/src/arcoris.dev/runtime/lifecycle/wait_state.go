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

// WaitState blocks until the lifecycle reaches target, the target becomes
// unreachable, or ctx is cancelled.
//
// WaitState returns the Snapshot whose State equals target on success. On
// failure it returns the latest Snapshot observed by the wait loop together with
// a WaitError.
//
// The target must be a valid State. If target is invalid, WaitState returns
// ErrInvalidWaitTarget.
//
// WaitState uses the static lifecycle transition graph to detect unreachable
// targets early. For example, waiting for StateStarting after the lifecycle is
// already StateRunning fails immediately because the lifecycle graph never moves
// backward. Waiting for StateRunning after the lifecycle reaches StateStopped or
// StateFailed also fails because terminal states have no outgoing transitions.
//
// Reachability is table-level reachability only. Guards may still prevent a
// table-reachable transition from committing. In that case WaitState keeps
// waiting until the target is reached, the lifecycle becomes terminal or
// unreachable, or ctx ends.
func (c *Controller) WaitState(ctx context.Context, target State) (Snapshot, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	snap, changed, done := c.waitSnapshot()

	if !target.IsValid() {
		return snap, newWaitStateError(snap, target, ErrInvalidWaitTarget)
	}

	for {
		if snap.State == target {
			return snap, nil
		}

		if !canReachState(snap.State, target) {
			return snap, newWaitStateError(snap, target, ErrWaitTargetUnreachable)
		}

		select {
		case <-changed:
			snap, changed, done = c.waitSnapshot()

		case <-done:
			snap, _, _ = c.waitSnapshot()
			if snap.State == target {
				return snap, nil
			}

			return snap, newWaitStateError(snap, target, ErrWaitTargetUnreachable)

		case <-ctx.Done():
			snap, _, _ = c.waitSnapshot()
			return snap, newWaitStateError(snap, target, ctx.Err())
		}
	}
}
