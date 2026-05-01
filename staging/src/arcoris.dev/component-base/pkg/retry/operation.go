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

// Operation is a caller-owned action executed by a retry loop.
//
// Operation is the basic unit of work for Do. A retry execution may call an
// Operation zero, one, or many times depending on context state, operation
// result, retryability classification, configured limits, and the configured
// backoff schedule.
//
// A retry execution may skip the first call when its owning context has already
// stopped before the first attempt starts. It may call the operation once when
// the first attempt succeeds, when the first failure is not retryable, or when
// the configured attempt limit allows only the initial attempt. It may call the
// operation multiple times when failures are classified as retryable and retry
// limits still allow another attempt.
//
// Operation implementations SHOULD observe ctx when execution can block, perform
// I/O, acquire resources, or depend on cancellation or deadlines. Pure
// in-memory operations MAY ignore ctx when cancellation cannot affect their
// execution.
//
// The caller owns retry safety. The retry package does not know whether an
// operation is idempotent, replayable, transactional, externally visible, or
// safe to repeat. Callers MUST configure retry only for operations whose retry
// semantics are valid for their domain.
//
// The retry package does not recover panics raised by an Operation. Panic
// recovery, if required, belongs to the operation owner, runtime supervisor, or
// an explicit wrapper outside this package.
//
// A nil Operation is a programming error and is rejected by retry entry points.
type Operation func(ctx context.Context) error

// ValueOperation is a caller-owned action that returns a value.
//
// ValueOperation is the value-returning counterpart of Operation and is the
// basic unit of work for DoValue. It follows the same retry ownership rules as
// Operation: a retry execution may call it zero, one, or many times, and the
// caller remains responsible for idempotency, replay safety, and external side
// effects.
//
// When a ValueOperation returns a nil error, DoValue returns the value produced
// by that successful attempt. When it returns a non-nil error, DoValue treats the
// attempt as failed and ignores the returned value. Failed-attempt values are not
// preserved, merged, or exposed by the retry package.
//
// ValueOperation implementations SHOULD observe ctx when execution can block,
// perform I/O, acquire resources, or depend on cancellation or deadlines. Pure
// in-memory operations MAY ignore ctx when cancellation cannot affect their
// execution.
//
// The retry package does not recover panics raised by a ValueOperation. Panic
// recovery, if required, belongs to the operation owner, runtime supervisor, or
// an explicit wrapper outside this package.
//
// A nil ValueOperation is a programming error and is rejected by retry entry
// points.
type ValueOperation[T any] func(ctx context.Context) (T, error)
