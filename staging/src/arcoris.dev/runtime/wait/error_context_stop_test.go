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
	"testing"
)

// TestContextStopErrorReturnsNilForActiveContext verifies that active contexts
// are not converted into wait-owned stop errors.
func TestContextStopErrorReturnsNilForActiveContext(t *testing.T) {
	t.Parallel()

	if err := contextStopError(context.Background()); err != nil {
		t.Fatalf("contextStopError(active context) = %v, want nil", err)
	}
}

// TestContextStopErrorClassifiesCancellation verifies wait-owned interruption
// classification for a cancelled context.
func TestContextStopErrorClassifiesCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := contextStopError(ctx)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
}

// TestContextStopErrorPreservesCancellationCause verifies that cancellation
// causes remain visible through the wait-owned interruption wrapper.
func TestContextStopErrorPreservesCancellationCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutdown requested")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)

	err := contextStopError(ctx)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
}

// TestContextStopErrorClassifiesDeadline verifies wait-owned timeout
// classification for an expired context deadline.
func TestContextStopErrorClassifiesDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	err := contextStopError(ctx)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
}

// TestContextStopErrorPreservesDeadlineCause verifies that custom deadline
// causes remain visible through the wait-owned timeout wrapper.
func TestContextStopErrorPreservesDeadlineCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("finite wait budget exhausted")
	ctx, cancel := context.WithTimeoutCause(context.Background(), 0, cause)
	defer cancel()

	err := contextStopError(ctx)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, cause)
	mustMatch(t, err, context.DeadlineExceeded)
}
