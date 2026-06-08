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

package retry

import "context"

// DoObserved executes op with bounded retry orchestration and returns the
// terminal Outcome.
//
// DoObserved is the observed form of Do. It is useful when callers need
// completion metadata without installing an Observer. The returned Outcome is
// valid for every non-panic terminal path and matches the Outcome delivered in
// the terminal EventRetryStop observer event.
//
// Successful execution returns StopReasonSucceeded and a nil error.
// Non-retryable operation errors are returned unchanged with
// StopReasonNonRetryable. Retry-owned exhaustion returns an ErrExhausted-
// classified error. Retry-owned context interruption returns an ErrInterrupted-
// classified error.
//
// DoObserved does not recover panics from op, classifiers, observers, clocks, or
// delay schedules. Panics mean no retry-owned terminal Outcome is produced.
func DoObserved(
	ctx context.Context,
	op Operation,
	opts ...Option,
) (Outcome, error) {
	requireContext(ctx)
	requireOperation(op)

	_, outcome, err := run(ctx, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, op(ctx)
	}, configOf(opts...))

	return outcome, err
}
