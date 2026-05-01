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

package retry

import (
	"context"
	"errors"
)

// contextStopError converts a completed context into a retry-owned interruption.
//
// The helper is intentionally private because it encodes retry-loop ownership,
// not a general context utility. Raw context errors returned by operations must
// remain operation-owned errors; only retry code that owns ctx may classify a
// context stop as ErrInterrupted.
//
// Unlike package wait, retry does not define a separate timeout classification.
// Context cancellation and context deadline expiration observed at a retry-owned
// boundary are both retry interruptions. The underlying context sentinel and
// custom cause remain visible through errors.Is and errors.As.
func contextStopError(ctx context.Context) error {
	err := ctx.Err()
	if err == nil {
		return nil
	}

	return NewInterruptedError(contextStopCause(ctx, err))
}

// contextStopCause returns the most specific cause available for a completed
// context.
//
// context.Cause returns a cancellation cause when callers use
// context.WithCancelCause, context.WithDeadlineCause, or context.WithTimeoutCause.
// For ordinary contexts, it returns the same sentinel as ctx.Err after the
// context is done.
//
// When the caller supplied a custom cause, ctx.Err still carries the standard
// context sentinel. Joining the two preserves both diagnostics and matching:
// errors.Is can still find context.Canceled or context.DeadlineExceeded, and it
// can also find the custom cause.
//
// The fallback to err keeps the helper robust for custom Context implementations
// that may not expose a richer cause.
func contextStopCause(ctx context.Context, err error) error {
	cause := context.Cause(ctx)
	if cause == nil {
		return err
	}
	if errors.Is(cause, err) {
		return cause
	}

	return errors.Join(err, cause)
}
