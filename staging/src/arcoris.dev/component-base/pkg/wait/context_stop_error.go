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
	"errors"
)

// contextStopError converts a completed context into a wait-owned stop error.
//
// The helper is intentionally private because it encodes wait-loop ownership,
// not a general context utility. Raw context errors returned by conditions must
// remain condition-owned errors; only the wait loop that owns ctx may classify a
// context stop as ErrInterrupted or ErrTimeout.
//
// Deadline expiration is classified as ErrTimeout, which also matches
// ErrInterrupted through the wait error hierarchy. All other context stops are
// classified as ErrInterrupted.
//
// contextStopError preserves context.Cause(ctx) when it is available. This keeps
// cancellation causes visible through errors.Is and errors.As while still adding
// the wait-owned classification layer.
func contextStopError(ctx context.Context) error {
	err := ctx.Err()
	if err == nil {
		return nil
	}

	cause := contextStopCause(ctx, err)
	if errors.Is(err, context.DeadlineExceeded) {
		return NewTimeoutError(cause)
	}

	return NewInterruptedError(cause)
}

// contextStopCause returns the most specific cause available for a completed
// context.
//
// context.Cause returns a cancellation cause when callers use context.WithCancelCause,
// context.WithDeadlineCause, or context.WithTimeoutCause. For ordinary contexts,
// it returns the same sentinel as ctx.Err after the context is done.
//
// The fallback to err keeps the helper robust for custom Context implementations
// that may not expose a richer cause.
func contextStopCause(ctx context.Context, err error) error {
	cause := context.Cause(ctx)
	if cause != nil {
		return cause
	}

	return err
}
