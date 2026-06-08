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

import (
	"context"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
)

// requireContext panics when ctx is nil.
//
// Public retry entry points must call requireContext before observing context
// state, starting attempts, creating delays, or building outcomes. Use
// context.Background when no narrower cancellation scope is available.
func requireContext(ctx context.Context) {
	if ctx == nil {
		panic(ErrNilContext)
	}
}

// requireOperation panics when op is nil.
//
// Operation is the executable unit of Do. A nil operation is a programming error
// and must not be represented as a retryable operation failure.
func requireOperation(op Operation) {
	if op == nil {
		panic(ErrNilOperation)
	}
}

// requireValueOperation panics when op is nil.
//
// ValueOperation is the executable unit of DoValue. A nil value operation is a
// programming error and must not be represented as a retryable operation failure.
func requireValueOperation[T any](op ValueOperation[T]) {
	if op == nil {
		panic(ErrNilValueOperation)
	}
}

// requireClock panics when c is nil.
//
// Retry execution uses the configured clock for attempt timestamps, terminal
// outcome timestamps, elapsed-time accounting, and delay timers. A nil clock is
// invalid configuration.
func requireClock(c clock.Clock) {
	if c == nil {
		panic(ErrNilClock)
	}
}

// requireDelaySchedule panics when sched is nil.
//
// Retry stores a delay.Schedule rather than a delay.Sequence so every
// Do/DoValue execution can create and own its own independent sequence.
func requireDelaySchedule(sched delay.Schedule) {
	if sched == nil {
		panic(ErrNilDelaySchedule)
	}
}

// requireDelaySequence panics when seq is nil.
//
// A nil sequence means the configured Schedule violated its NewSequence contract.
// Retry reports this as a programming error at the schedule boundary.
func requireDelaySequence(seq delay.Sequence) {
	if seq == nil {
		panic(ErrNilDelaySequence)
	}
}

// requireDelay panics when d is negative and ok is true.
//
// Sequence.Next returns a meaningful delay only when ok is true. Finite sequence
// exhaustion, represented by ok=false, is handled by the retry loop as
// StopReasonDelayExhausted and must not be treated as a validation failure.
func requireDelay(d time.Duration, ok bool) {
	if ok && d < 0 {
		panic(ErrNegativeDelay)
	}
}

// requireClassifier panics when classifier is nil.
//
// Classifier decides whether an operation-owned error may be retried. A nil
// classifier would make retry behavior undefined after the first failed attempt.
func requireClassifier(classifier Classifier) {
	if classifier == nil {
		panic(ErrNilClassifier)
	}
}

// requireRetryableFunc panics when fn is nil.
//
// The diagnostic reuses the same stable message as ClassifierFunc. This keeps the
// nil-function behavior identical whether the function is called through
// ClassifierFunc directly or supplied through a future WithRetryable option.
func requireRetryableFunc(fn func(error) bool) {
	if fn == nil {
		panic(ErrNilClassifierFunc)
	}
}

// requireMaxAttempts panics when n is zero.
//
// MaxAttempts includes the initial operation call. A value of one means no retry
// attempts beyond the initial call. A value of zero cannot describe a valid retry
// execution.
func requireMaxAttempts(n uint) {
	if n == 0 {
		panic(ErrZeroMaxAttempts)
	}
}

// requireMaxElapsed panics when d is negative.
//
// A zero duration disables elapsed-time limiting. Positive durations bound the
// total runtime of one retry execution.
func requireMaxElapsed(d time.Duration) {
	if d < 0 {
		panic(ErrNegativeMaxElapsed)
	}
}

// requireObserver panics when observer is nil.
//
// Observers are optional, but once configured they must be callable.
func requireObserver(observer Observer) {
	if observer == nil {
		panic(ErrNilObserver)
	}
}

// requireObserverFunc panics when fn is nil.
//
// The diagnostic reuses the same stable message as ObserverFunc. This keeps the
// nil-function behavior identical whether the function is called through
// ObserverFunc directly or supplied through a future WithObserverFunc option.
func requireObserverFunc(fn func(context.Context, Event)) {
	if fn == nil {
		panic(ErrNilObserverFunc)
	}
}

// requireOption panics when option is nil.
//
// Nil options are rejected instead of ignored. This makes option composition
// errors visible at configuration time and keeps retry defaults explicit.
func requireOption(opt Option) {
	if opt == nil {
		panic(ErrNilOption)
	}
}
