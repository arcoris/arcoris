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

	"arcoris.dev/component-base/pkg/internal/contextstop"
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
// contextStopError preserves context.Cause(ctx) when it is available. When the
// context exposes both a standard sentinel and a custom cause, both remain
// visible through errors.Is or errors.As while the wait-owned classification
// layer is added.
func contextStopError(ctx context.Context) error {
	err := ctx.Err()
	if err == nil {
		return nil
	}

	cause := contextstop.Cause(ctx, err)
	if errors.Is(err, context.DeadlineExceeded) {
		return NewTimeoutError(cause)
	}

	return NewInterruptedError(cause)
}
