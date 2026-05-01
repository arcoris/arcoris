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

package backoff

import (
	"math"
	"time"
)

const (
	// errNilValidationMessage is the stable diagnostic text used when a
	// validation helper is called without a diagnostic message.
	//
	// Validation helpers are package-local and should always be called with the
	// stable diagnostic text owned by the file that defines the rejected
	// configuration. A missing message is a package wiring error, not user input
	// failure.
	errNilValidationMessage = "backoff: nil validation message"
)

// requireValidationMessage panics when message is empty.
//
// The helper protects package-local validation helpers from being called without
// a stable diagnostic string. Public constructors should expose precise messages
// such as "backoff: negative fixed delay" instead of generic validation output.
func requireValidationMessage(message string) {
	if message == "" {
		panic(errNilValidationMessage)
	}
}

// requireSchedule panics when schedule is nil.
//
// The message argument must be the stable diagnostic text owned by the caller's
// file. For example, Cap should pass errNilCapSchedule and Limit should pass
// errNilLimitSchedule. This keeps validation behavior centralized while keeping
// ownership of diagnostic text local to the feature being validated.
func requireSchedule(schedule Schedule, message string) {
	requireValidationMessage(message)

	if schedule == nil {
		panic(message)
	}
}

// requireSequence panics when sequence is nil.
//
// Schedule.NewSequence must return a usable Sequence. A nil Sequence violates
// the Schedule contract and should be reported at the wrapper boundary that
// observes it.
func requireSequence(sequence Sequence, message string) {
	requireValidationMessage(message)

	if sequence == nil {
		panic(message)
	}
}

// requireRandomSource panics when source is nil.
//
// RandomSource values are factories for per-sequence random generators. A nil source
// cannot provide deterministic or runtime pseudo-randomness and indicates invalid
// package use.
func requireRandomSource(source RandomSource, message string) {
	requireValidationMessage(message)

	if source == nil {
		panic(message)
	}
}

// requireRandom panics when random is nil.
//
// random generators are used by random and jitter sequences to draw pseudo-random
// offsets. A nil Random would fail later in delay generation, so it is rejected
// at the random boundary.
func requireRandom(random RandomGenerator, message string) {
	requireValidationMessage(message)

	if random == nil {
		panic(message)
	}
}

// requireRandomOption panics when option is nil.
//
// RandomOption values mutate package-local random configuration. A nil option
// cannot be applied and indicates invalid schedule construction.
func requireRandomOption(option RandomOption, message string) {
	requireValidationMessage(message)

	if option == nil {
		panic(message)
	}
}

// requireRandomConfig panics when config is nil.
//
// Public callers should not normally trigger this. It protects option functions
// and package-local tests from nil configuration wiring.
func requireRandomConfig(config *randomConfig, message string) {
	requireValidationMessage(message)

	if config == nil {
		panic(message)
	}
}

// requireNonNegativeDuration panics when duration is negative.
//
// Zero is accepted. This is the correct validation for schedules where immediate
// continuation is valid, such as Fixed, Linear initial delay, Cap maximum delay,
// random minimum delay, and explicit delay sequences.
func requireNonNegativeDuration(duration time.Duration, message string) {
	requireValidationMessage(message)

	if isNegativeDuration(duration) {
		panic(message)
	}
}

// requirePositiveDuration panics when duration is zero or negative.
//
// This is the correct validation for schedules whose mathematical model requires
// a strictly positive base value, such as Exponential and Fibonacci.
func requirePositiveDuration(duration time.Duration, message string) {
	requireValidationMessage(message)

	if !isPositiveDuration(duration) {
		panic(message)
	}
}

// requireDurationNotBefore panics when upper is smaller than lower.
//
// The helper is useful for validating closed duration ranges such as
// Random(minDelay, maxDelay) and DecorrelatedJitter(initial, maxDelay,
// multiplier). Equal bounds are valid and produce a single-value range.
func requireDurationNotBefore(upper, lower time.Duration, message string) {
	requireValidationMessage(message)

	if upper < lower {
		panic(message)
	}
}

// requireNonNegativeCount panics when count is negative.
//
// Zero is accepted. This is the correct validation for finite wrappers such as
// Limit, where zero means immediate exhaustion.
func requireNonNegativeCount(count int, message string) {
	requireValidationMessage(message)

	if count < 0 {
		panic(message)
	}
}

// requireFiniteFloat panics when value is NaN or infinite.
//
// The helper accepts negative and zero finite values. Callers that need stricter
// semantics should use a more specific helper such as requireFloatGreaterThanOne,
// requireJitterFactor, or requireJitterRatio.
func requireFiniteFloat(value float64, message string) {
	requireValidationMessage(message)

	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic(message)
	}
}

// requireFloatGreaterThanOne panics when value is not finite or is less than or
// equal to one.
//
// This is the correct validation for growth multipliers used by Exponential and
// DecorrelatedJitter. A multiplier of one is fixed, not growth. A multiplier
// below one is decay, not backoff growth.
func requireFloatGreaterThanOne(value float64, message string) {
	requireValidationMessage(message)

	if value <= 1 || math.IsNaN(value) || math.IsInf(value, 0) {
		panic(message)
	}
}

// requireJitterFactor panics when factor is not a valid non-negative finite
// jitter factor.
//
// A factor of zero is valid and disables factor-based expansion. Negative, NaN,
// and infinite values cannot describe a stable runtime delay transformation.
func requireJitterFactor(factor float64) {
	if factor < 0 || math.IsNaN(factor) || math.IsInf(factor, 0) {
		panic(errInvalidJitterFactor)
	}
}

// requireJitterRatio panics when ratio is not a valid finite ratio in [0, 1].
//
// Ratio-based jitter algorithms in this package use ratios in [0, 1] to keep
// computed lower bounds non-negative. A zero ratio disables randomization. A
// ratio of one allows the lower bound to reach zero.
func requireJitterRatio(ratio float64) {
	if ratio < 0 || ratio > 1 || math.IsNaN(ratio) || math.IsInf(ratio, 0) {
		panic(errInvalidJitterRatio)
	}
}

// requireNonNegativeSequenceDelay panics when delay is negative while ok is
// true.
//
// Sequence implementations must return non-negative delays when they report an
// available value. Wrappers should call this at child boundaries before applying
// transformations such as Cap, Limit, or jitter. The message should identify the
// wrapper that observed the contract violation.
func requireNonNegativeSequenceDelay(delay time.Duration, ok bool, message string) {
	requireValidationMessage(message)

	if ok && isNegativeDuration(delay) {
		panic(message)
	}
}
