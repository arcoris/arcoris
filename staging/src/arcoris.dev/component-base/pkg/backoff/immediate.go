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

import "time"

// Immediate returns a schedule that produces an infinite stream of zero delays.
//
// Immediate is the smallest possible backoff schedule. Every sequence created by
// the returned Schedule reports delay=0, ok=true for every call to Next. A zero
// delay means the owner may continue immediately without waiting.
//
// Example:
//
//	schedule := backoff.Immediate()
//	sequence := schedule.NewSequence()
//	delay, ok := sequence.Next()
//	_ = delay
//	_ = ok
//
// Immediate is useful for:
//
//   - one-shot transparent retry paths;
//   - tests that need deterministic no-wait retry behavior;
//   - local compare-and-swap or lock-acquisition loops with external bounds;
//   - composing with Limit to allow a small number of immediate retries;
//   - composing with Sequence when the first retry should be immediate and later
//     retries should use non-zero delays.
//
// Immediate does not provide overload protection by itself. An unbounded
// immediate retry loop can become a busy loop if the owner does not apply an
// external limit such as retry max attempts, a finite Limit wrapper, context
// cancellation, or another higher-level stop condition.
//
// The returned Schedule is stateless and safe to reuse. Each call to NewSequence
// returns an independent Sequence value. The concrete sequence is also stateless,
// but callers should still follow the package-wide single-owner Sequence model
// and avoid sharing one Sequence across unrelated runtime loops.
//
// Immediate does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
func Immediate() Schedule {
	return immediateSchedule{}
}

// immediateSchedule is the reusable recipe behind Immediate.
//
// The type has no fields because an immediate schedule has no configuration and
// no shared mutable state. It exists as a named implementation instead of a
// closure so the public constructor can stay allocation-friendly, explicit, and
// easy to inspect in tests and documentation.
type immediateSchedule struct{}

// NewSequence returns an independent immediate delay sequence.
//
// The sequence is stateless and infinite. It is safe to create many sequences
// from the same immediateSchedule because no sequence state is stored on the
// schedule value.
func (immediateSchedule) NewSequence() Sequence {
	return immediateSequence{}
}

// immediateSequence is the per-owner delay stream produced by Immediate.
//
// The sequence has no fields because every call to Next returns the same result:
// zero delay and ok=true. It is modeled as a Sequence rather than being handled
// specially by retry or polling code so Immediate participates in the same
// composition model as fixed, linear, exponential, capped, limited, and jittered
// schedules.
type immediateSequence struct{}

// Next returns a zero delay and reports that the sequence is still available.
//
// The sequence is intentionally infinite. Exhaustion, if required, should be
// provided by a higher-level owner such as retry max attempts or by a future
// finite wrapper such as Limit.
func (immediateSequence) Next() (time.Duration, bool) {
	return 0, true
}
