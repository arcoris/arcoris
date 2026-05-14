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

const (
	// errNilScheduleFunc is the stable diagnostic text used when a
	// ScheduleFunc method is called on a nil function value.
	//
	// A nil ScheduleFunc cannot create a Sequence and indicates invalid package
	// use rather than finite schedule exhaustion. The adapter panics immediately
	// with this message so the invalid schedule boundary is visible at the point
	// where a sequence is requested.
	errNilScheduleFunc = "delay: nil ScheduleFunc"

	// errScheduleFuncReturnedNilSequence is the stable diagnostic text used when
	// a ScheduleFunc returns a nil Sequence.
	//
	// Schedule.NewSequence must return a usable Sequence. Returning nil violates
	// the Schedule contract and would otherwise move a construction-time
	// programming error into a later retry, polling, reconnect, or controller
	// loop.
	errScheduleFuncReturnedNilSequence = "delay: ScheduleFunc returned nil Sequence"
)

// ScheduleFunc adapts a function into a Schedule.
//
// ScheduleFunc is useful for small custom schedules, tests, and adapters that
// need to satisfy Schedule without declaring a named type. The wrapped function
// is called every time NewSequence is invoked and must return a non-nil
// Sequence.
//
// Example:
//
//	schedule := delay.ScheduleFunc(func() delay.Sequence {
//		return delay.Fixed(time.Second).NewSequence()
//	})
//	sequence := schedule.NewSequence()
//	_ = sequence
//
// The function follows the same ownership model as any other Schedule
// implementation:
//
//   - it should create an independent Sequence for each call;
//   - it should not share mutable sequence state across unrelated owners;
//   - it must not sleep, create timers, observe context cancellation, execute
//     operations, classify errors, retry work, log, trace, export metrics,
//     schedule queue items, rate limit callers, or make domain decisions.
//
// A nil ScheduleFunc is a programming error. NewSequence panics immediately
// instead of returning a delayed nil dereference from a runtime loop.
//
// A ScheduleFunc that returns nil violates the Schedule contract. NewSequence
// panics so the invalid adapter is detected at the boundary where the sequence
// is requested.
//
// ScheduleFunc does not recover panics raised by the wrapped function. Panic
// recovery, if required, belongs to the caller or to an explicit higher-level
// wrapper.
type ScheduleFunc func() Sequence

// NewSequence calls f and returns the Sequence produced by it.
//
// NewSequence panics when f is nil or when f returns nil. Both cases are
// programming errors, not finite delay exhaustion.
func (f ScheduleFunc) NewSequence() Sequence {
	if f == nil {
		panic(errNilScheduleFunc)
	}

	seq := f()
	if seq == nil {
		panic(errScheduleFuncReturnedNilSequence)
	}

	return seq
}
