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

import "context"

// DoValue executes op with bounded retry orchestration and returns the value
// produced by the successful attempt.
//
// DoValue is the public entry point for operations that return a value and an
// error. It applies options, delegates retry execution to the package-private
// runtime engine, and returns either the successful attempt value or the zero
// value of T with the terminal error.
//
// The returned value is meaningful only when err is nil. When retry execution
// fails, DoValue returns the zero value of T. Values returned by failed attempts
// are ignored and are not preserved, merged, or exposed by the retry package.
//
// The default configuration is conservative: DoValue calls op at most once, uses
// NeverRetry as the classifier, has no elapsed-time limit, and registers no
// observers. Callers must explicitly configure retryability and limits before
// failed operation attempts can be retried.
//
// DoValue may call op zero, one, or many times:
//
//   - zero times when ctx is already stopped before the first attempt;
//   - one time when the first attempt succeeds, fails with a non-retryable error,
//     or the configured attempt limit allows only the initial attempt;
//   - many times when operation errors are classified as retryable and retry
//     limits, context state, and delay sequence availability allow more
//     attempts.
//
// Operation errors that are not retried are returned unchanged. Retry-owned
// exhaustion is returned as ErrExhausted. Retry-owned context interruption is
// returned as ErrInterrupted.
//
// DoValue does not infer idempotency, replay safety, transaction safety,
// protocol semantics, or external side-effect safety. The caller owns the
// decision that op may be safely retried with the configured classifier and
// limits.
//
// DoValue does not recover panics raised by op or by configured observers. Panic
// recovery, if required, belongs to the operation owner, observer implementation,
// runtime supervisor, or an explicit wrapper outside this package.
//
// DoValue panics when ctx is nil, op is nil, or any supplied option is nil or
// otherwise invalid.
func DoValue[T any](
	ctx context.Context,
	op ValueOperation[T],
	opts ...Option,
) (T, error) {
	requireContext(ctx)
	requireValueOperation(op)

	return run(ctx, op, configOf(opts...))
}
