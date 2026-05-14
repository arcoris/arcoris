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
	"time"

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/delay"
)

const (
	// panicNilContext is the stable diagnostic text used when retry execution is
	// started without a context.
	//
	// A retry execution owns context observation at retry boundaries. A nil
	// context would fail later in cancellation checks or delay waits, so public
	// entry points reject it immediately.
	panicNilContext = "retry: nil context"

	// panicNilOperation is the stable diagnostic text used when Do receives a nil
	// Operation.
	//
	// A nil operation cannot be executed and indicates invalid caller wiring, not
	// a retryable runtime failure.
	panicNilOperation = "retry: nil operation"

	// panicNilValueOperation is the stable diagnostic text used when DoValue
	// receives a nil ValueOperation.
	//
	// A nil value operation cannot be executed and indicates invalid caller
	// wiring, not a retryable runtime failure.
	panicNilValueOperation = "retry: nil value operation"

	// panicNilClock is the stable diagnostic text used when retry configuration
	// receives a nil clock.
	//
	// Retry execution needs a clock for attempt timestamps, elapsed-time checks,
	// and retry delay timers. A nil clock would fail inside the runtime loop, so
	// it is rejected at configuration boundaries.
	panicNilClock = "retry: nil clock"

	// panicNilDelaySchedule is the stable diagnostic text used when retry
	// configuration receives a nil delay schedule.
	//
	// Retry stores a reusable delay.Schedule and creates a fresh Sequence for
	// each execution. A nil schedule cannot produce per-execution delay streams.
	panicNilDelaySchedule = "retry: nil delay schedule"

	// panicNilDelaySequence is the stable diagnostic text used when a configured
	// delay schedule returns a nil sequence.
	//
	// Schedule.NewSequence must return a usable Sequence. Returning nil violates
	// the delay.Schedule contract and is reported at the retry boundary that
	// observes it.
	panicNilDelaySequence = "retry: delay schedule returned nil Sequence"

	// panicNegativeDelay is the stable diagnostic text used when a delay
	// sequence returns a negative delay while reporting ok=true.
	//
	// A zero delay is valid and means immediate retry. A negative delay violates
	// the delay.Sequence contract and would make retry waiting semantics
	// ambiguous.
	panicNegativeDelay = "retry: delay sequence returned negative delay"

	// panicNilClassifier is the stable diagnostic text used when retry
	// configuration receives a nil Classifier.
	//
	// Retryability classification is required after operation-owned failures. A
	// nil classifier would fail later in the runtime loop and is rejected at
	// configuration boundaries.
	panicNilClassifier = "retry: nil classifier"

	// panicZeroMaxAttempts is the stable diagnostic text used when a caller
	// configures zero max attempts.
	//
	// Max attempts includes the initial operation call. A value of zero cannot
	// describe a valid retry execution policy; callers that want no retries
	// should use one attempt.
	panicZeroMaxAttempts = "retry: zero max attempts"

	// panicNegativeMaxElapsed is the stable diagnostic text used when a caller
	// configures a negative max elapsed duration.
	//
	// A zero max elapsed duration disables elapsed-time limiting. A negative
	// duration cannot describe a stable retry boundary.
	panicNegativeMaxElapsed = "retry: negative max elapsed"

	// panicNilObserver is the stable diagnostic text used when retry
	// configuration receives a nil Observer.
	//
	// Observers are optional, but configured observers must be callable.
	panicNilObserver = "retry: nil observer"

	// panicNilOption is the stable diagnostic text used when retry configuration
	// receives a nil Option.
	//
	// Nil options usually indicate invalid conditional option composition. Retry
	// rejects them immediately instead of silently ignoring caller mistakes.
	panicNilOption = "retry: nil option"
)

// requireContext panics when ctx is nil.
//
// Public retry entry points must call requireContext before observing context
// state, starting attempts, creating delays, or building outcomes. Use
// context.Background when no narrower cancellation scope is available.
func requireContext(ctx context.Context) {
	if ctx == nil {
		panic(panicNilContext)
	}
}

// requireOperation panics when op is nil.
//
// Operation is the executable unit of Do. A nil operation is a programming error
// and must not be represented as a retryable operation failure.
func requireOperation(op Operation) {
	if op == nil {
		panic(panicNilOperation)
	}
}

// requireValueOperation panics when op is nil.
//
// ValueOperation is the executable unit of DoValue. A nil value operation is a
// programming error and must not be represented as a retryable operation failure.
func requireValueOperation[T any](op ValueOperation[T]) {
	if op == nil {
		panic(panicNilValueOperation)
	}
}

// requireClock panics when c is nil.
//
// Retry execution uses the configured clock for attempt timestamps, terminal
// outcome timestamps, elapsed-time accounting, and delay timers. A nil clock is
// invalid configuration.
func requireClock(c clock.Clock) {
	if c == nil {
		panic(panicNilClock)
	}
}

// requireDelaySchedule panics when sched is nil.
//
// Retry stores a delay.Schedule rather than a delay.Sequence so every
// Do/DoValue execution can create and own its own independent sequence.
func requireDelaySchedule(sched delay.Schedule) {
	if sched == nil {
		panic(panicNilDelaySchedule)
	}
}

// requireDelaySequence panics when seq is nil.
//
// A nil sequence means the configured Schedule violated its NewSequence contract.
// Retry reports this as a programming error at the schedule boundary.
func requireDelaySequence(seq delay.Sequence) {
	if seq == nil {
		panic(panicNilDelaySequence)
	}
}

// requireDelay panics when d is negative and ok is true.
//
// Sequence.Next returns a meaningful delay only when ok is true. Finite sequence
// exhaustion, represented by ok=false, is handled by the retry loop as
// StopReasonDelayExhausted and must not be treated as a validation failure.
func requireDelay(d time.Duration, ok bool) {
	if ok && d < 0 {
		panic(panicNegativeDelay)
	}
}

// requireClassifier panics when classifier is nil.
//
// Classifier decides whether an operation-owned error may be retried. A nil
// classifier would make retry behavior undefined after the first failed attempt.
func requireClassifier(classifier Classifier) {
	if classifier == nil {
		panic(panicNilClassifier)
	}
}

// requireRetryableFunc panics when fn is nil.
//
// The diagnostic reuses the same stable message as ClassifierFunc. This keeps the
// nil-function behavior identical whether the function is called through
// ClassifierFunc directly or supplied through a future WithRetryable option.
func requireRetryableFunc(fn func(error) bool) {
	if fn == nil {
		panic(panicNilClassifierFunc)
	}
}

// requireMaxAttempts panics when n is zero.
//
// MaxAttempts includes the initial operation call. A value of one means no retry
// attempts beyond the initial call. A value of zero cannot describe a valid retry
// execution.
func requireMaxAttempts(n uint) {
	if n == 0 {
		panic(panicZeroMaxAttempts)
	}
}

// requireMaxElapsed panics when d is negative.
//
// A zero duration disables elapsed-time limiting. Positive durations bound the
// total runtime of one retry execution.
func requireMaxElapsed(d time.Duration) {
	if d < 0 {
		panic(panicNegativeMaxElapsed)
	}
}

// requireObserver panics when observer is nil.
//
// Observers are optional, but once configured they must be callable.
func requireObserver(observer Observer) {
	if observer == nil {
		panic(panicNilObserver)
	}
}

// requireObserverFunc panics when fn is nil.
//
// The diagnostic reuses the same stable message as ObserverFunc. This keeps the
// nil-function behavior identical whether the function is called through
// ObserverFunc directly or supplied through a future WithObserverFunc option.
func requireObserverFunc(fn func(context.Context, Event)) {
	if fn == nil {
		panic(panicNilObserverFunc)
	}
}

// requireOption panics when option is nil.
//
// Nil options are rejected instead of ignored. This makes option composition
// errors visible at configuration time and keeps retry defaults explicit.
func requireOption(option Option) {
	if option == nil {
		panic(panicNilOption)
	}
}
