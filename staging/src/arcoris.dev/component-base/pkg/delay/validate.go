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

package delay

import (
	"math"
	"time"
)

const (
	// errNilValidationMessage is the stable diagnostic text used when a
	// package-local validation helper is called without a feature-owned panic
	// message.
	//
	// Public constructors and adapters own precise messages such as "delay: nil
	// cap schedule". Validation helpers only enforce shared mechanics, so a
	// missing message is an internal package wiring error.
	errNilValidationMessage = "delay: nil validation message"
)

// requireValidationMessage panics when message is empty.
//
// The helper keeps diagnostic ownership local to the rejecting feature while
// still centralizing the mechanical "message must exist" check used by the
// package-local validation helpers.
func requireValidationMessage(msg string) {
	if msg == "" {
		panic(errNilValidationMessage)
	}
}

// requireSchedule panics when schedule is nil.
//
// Wrapper constructors call this before storing a child schedule. A nil child
// schedule is invalid programmer input, not finite sequence exhaustion, because
// no per-owner Sequence can be created from it.
func requireSchedule(sched Schedule, msg string) {
	requireValidationMessage(msg)
	if sched == nil {
		panic(msg)
	}
}

// requireSequence panics when sequence is nil.
//
// Schedule.NewSequence must return a usable Sequence. Wrappers call this at the
// boundary where they observe child sequence creation so contract violations are
// reported before runtime owners consume delay values.
func requireSequence(seq Sequence, msg string) {
	requireValidationMessage(msg)
	if seq == nil {
		panic(msg)
	}
}

// requireNonNegativeDuration panics when duration is negative.
//
// Zero is accepted because delay=0, ok=true means immediate continuation. Public
// constructors use this for concrete delay values and bounds that must never be
// negative.
func requireNonNegativeDuration(d time.Duration, msg string) {
	requireValidationMessage(msg)
	if d < 0 {
		panic(msg)
	}
}

// requirePositiveDuration panics when d is zero or negative.
//
// Exponential and Fibonacci require strictly positive base durations. Without a
// positive base, a deterministic growth schedule would either never grow or
// would start from an invalid runtime duration.
func requirePositiveDuration(d time.Duration, msg string) {
	requireValidationMessage(msg)
	if d <= 0 {
		panic(msg)
	}
}

// requireNonNegativeCount panics when count is negative.
//
// Zero is accepted for finite wrappers such as Limit, where it means an
// immediately exhausted sequence. Negative counts cannot describe a stable
// sequence boundary.
func requireNonNegativeCount(n int, msg string) {
	requireValidationMessage(msg)
	if n < 0 {
		panic(msg)
	}
}

// requireFloatGreaterThanOne panics when v is not finite and greater than one.
//
// Growth multipliers describe expansion. A multiplier of one is fixed, a
// smaller multiplier is decay, and NaN or infinity cannot produce a stable
// deterministic delay schedule.
func requireFloatGreaterThanOne(v float64, msg string) {
	requireValidationMessage(msg)
	if v <= 1 || math.IsNaN(v) || math.IsInf(v, 0) {
		panic(msg)
	}
}

// requireNonNegativeSequenceDelay panics when delay is negative and ok is true.
//
// Sequence.Next ignores delay when ok=false, so negative exhausted delay values
// are not contract violations. A negative available delay is invalid and must be
// rejected at wrapper or adapter boundaries before it can reach clocks, waits, or
// retry orchestration.
func requireNonNegativeSequenceDelay(d time.Duration, ok bool, msg string) {
	requireValidationMessage(msg)
	if ok && d < 0 {
		panic(msg)
	}
}
