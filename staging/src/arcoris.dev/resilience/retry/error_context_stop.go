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

	"arcoris.dev/resilience/internal/contextstop"
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

	return NewInterruptedError(contextstop.Cause(ctx, err))
}
