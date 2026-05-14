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

package jitter

import (
	"math"
	"time"

	"arcoris.dev/chrono/delay"
)

const (
	// errNilValidationMessage is the stable diagnostic text used when a
	// package-local validation helper is called without a feature-owned panic
	// message.
	//
	// Public constructors and random adapters own precise messages such as
	// "jitter: invalid ratio". Shared helpers only enforce the common mechanics,
	// so an empty message is an internal package wiring error.
	errNilValidationMessage = "jitter: nil validation message"
)

// requireValidationMessage panics when message is empty.
//
// The helper keeps panic-message ownership close to the constructor or adapter
// that rejects input while still centralizing the mechanical guard used by this
// package's private validation functions.
func requireValidationMessage(msg string) {
	if msg == "" {
		panic(errNilValidationMessage)
	}
}

// requireSchedule panics when schedule is nil.
//
// Randomized wrappers need a child delay.Schedule to transform. A nil child is a
// construction-time programming error, not a finite sequence exhaustion signal.
func requireSchedule(sched delay.Schedule, msg string) {
	requireValidationMessage(msg)
	if sched == nil {
		panic(msg)
	}
}

// requireSequence panics when sequence is nil.
//
// delay.Schedule.NewSequence must return a usable Sequence. Jitter wrappers call
// this immediately after child sequence creation so the contract violation is
// visible before any random transform is applied.
func requireSequence(seq delay.Sequence, msg string) {
	requireValidationMessage(msg)
	if seq == nil {
		panic(msg)
	}
}

// requireRandomSource panics when source is nil.
//
// Schedules store RandomSource values so each sequence can receive its own
// random generator. A nil source cannot create per-sequence state and is invalid
// configuration.
func requireRandomSource(src RandomSource, msg string) {
	requireValidationMessage(msg)
	if src == nil {
		panic(msg)
	}
}

// requireRandom panics when random is nil.
//
// Randomized sequences call RandomGenerator.Int63 while producing concrete
// delays. A nil generator would fail later and obscure the configuration
// boundary, so constructors and source adapters reject it immediately.
func requireRandom(r RandomGenerator, msg string) {
	requireValidationMessage(msg)
	if r == nil {
		panic(msg)
	}
}

// requireRandomOption panics when option is nil.
//
// Random options mutate package-local configuration during schedule
// construction. Silently ignoring a nil option would hide invalid conditional
// option assembly by the caller.
func requireRandomOption(opt RandomOption, msg string) {
	requireValidationMessage(msg)
	if opt == nil {
		panic(msg)
	}
}

// requireRandomConfig panics when config is nil.
//
// Public callers cannot pass randomConfig directly. This helper protects
// package-owned option wiring and tests from mutating a nil configuration
// pointer.
func requireRandomConfig(cfg *randomConfig, msg string) {
	requireValidationMessage(msg)
	if cfg == nil {
		panic(msg)
	}
}

// requireNonNegativeDuration panics when duration is negative.
//
// Uniform bounds and other concrete delay ranges may include zero, but negative
// runtime delays have no valid interpretation for delay streams.
func requireNonNegativeDuration(d time.Duration, msg string) {
	requireValidationMessage(msg)
	if d < 0 {
		panic(msg)
	}
}

// requirePositiveDuration panics when duration is zero or negative.
//
// Decorrelated jitter uses this for its lower bound. A zero or negative initial
// value cannot form a valid closed random range, and allowing it would move an
// invalid constructor input into sequence-time duration arithmetic.
func requirePositiveDuration(d time.Duration, msg string) {
	requireValidationMessage(msg)
	if d <= 0 {
		panic(msg)
	}
}

// requireDurationNotBefore panics when upper is smaller than lower.
//
// Randomized schedules use closed ranges. Equal bounds are valid and represent a
// single possible delay value; reversed bounds are invalid construction input.
func requireDurationNotBefore(hi, lo time.Duration, msg string) {
	requireValidationMessage(msg)
	if hi < lo {
		panic(msg)
	}
}

// requireFloatGreaterThanOne panics when value is not finite and greater than
// one.
//
// Decorrelated jitter needs a multiplier that can expand the range above the
// previous delay. One is fixed, values below one are decay, and NaN or infinity
// cannot define stable runtime behavior.
func requireFloatGreaterThanOne(v float64, msg string) {
	requireValidationMessage(msg)
	if v <= 1 || math.IsNaN(v) || math.IsInf(v, 0) {
		panic(msg)
	}
}

// requireJitterFactor panics when factor is not a finite non-negative value.
//
// Positive jitter factors describe one-sided expansion above the child delay. A
// factor of zero is valid and disables expansion.
func requireJitterFactor(f float64) {
	if f < 0 || math.IsNaN(f) || math.IsInf(f, 0) {
		panic(errInvalidJitterFactor)
	}
}

// requireJitterRatio panics when ratio is not finite and inside [0, 1].
//
// Proportional jitter uses the ratio to compute both lower and upper bounds
// around a base delay. Restricting the ratio to [0, 1] keeps the lower bound
// non-negative.
func requireJitterRatio(r float64) {
	if r < 0 || r > 1 || math.IsNaN(r) || math.IsInf(r, 0) {
		panic(errInvalidJitterRatio)
	}
}

// requireNonNegativeSequenceDelay panics when delay is negative and ok is true.
//
// Child exhaustion is represented by ok=false and preserves the child's
// availability semantics. A negative available delay violates delay.Sequence and
// must be rejected before a random transform can hide or amplify it.
func requireNonNegativeSequenceDelay(d time.Duration, ok bool, msg string) {
	requireValidationMessage(msg)
	if ok && d < 0 {
		panic(msg)
	}
}
