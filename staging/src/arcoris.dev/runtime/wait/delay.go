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

package wait

import (
	"context"
	"time"
)

// Delay waits for duration to elapse or for ctx to stop.
//
// Delay is a low-level wait primitive. It owns only one cancellable delay. It
// does not evaluate conditions, retry operations, apply backoff, add jitter,
// record metrics, recover panics, or make scheduler decisions.
//
// The returned error has wait-owned classification semantics:
//
//   - nil means duration elapsed successfully or duration was non-positive while
//     ctx was still active;
//   - ErrInterrupted means the wait-owned context was cancelled before the delay
//     completed;
//   - ErrTimeout means the wait-owned context deadline expired before the delay
//     completed. Timeout also matches ErrInterrupted.
//
// Raw context errors are not returned directly when ctx stops. Delay converts
// the stop into NewInterruptedError or NewTimeoutError so callers can
// distinguish wait-owned interruption from condition-owned context errors.
//
// Cancellation causes created with context.WithCancelCause,
// context.WithDeadlineCause, or context.WithTimeoutCause are preserved as the
// wrapped cause. Callers can inspect them with errors.Is or errors.As.
//
// Non-positive durations are treated as an immediate delay after the pre-flight
// context check. This mirrors the useful part of time.Sleep semantics and keeps
// Delay usable as a mechanical primitive for future policies that may
// intentionally produce a zero delay. Fixed-cadence loop APIs, such as Until,
// validate their interval separately when a non-positive value would produce a
// busy loop.
//
// Delay panics when ctx is nil.
func Delay(ctx context.Context, duration time.Duration) error {
	requireContext(ctx)

	if err := contextStopError(ctx); err != nil {
		return err
	}
	if duration <= 0 {
		return nil
	}

	return NewTimer(duration).Wait(ctx)
}
