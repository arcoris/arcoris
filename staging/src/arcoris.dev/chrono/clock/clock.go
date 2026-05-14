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

package clock

import "time"

// PassiveClock provides read-only access to runtime time.
//
// PassiveClock is the smallest clock contract in this package. Components should
// depend on PassiveClock when they only need to read the current time or measure
// elapsed duration and do not own timers, tickers, sleeps, retry loops, or
// background control loops.
//
// Typical PassiveClock users include:
//
//   - metadata writers that need creation or update timestamps;
//   - admission evaluators that check request deadlines;
//   - queue eligibility checks that compare not-before or retry-after times;
//   - rate and window calculators that compute elapsed time lazily;
//   - diagnostics and sampling code that records observation time.
//
// PassiveClock intentionally does not expose waiting primitives. This keeps
// read-only components from accidentally gaining ownership over blocking,
// sleeping, or background scheduling behavior.
//
// PassiveClock is a runtime-time abstraction. It does not provide distributed
// causality, resource versions, generations, lease epochs, fencing tokens,
// logical clocks, or cluster-wide time synchronization.
type PassiveClock interface {
	// Now returns the clock's current time.
	//
	// For RealClock, Now returns the current wall-clock time from the Go standard
	// library. For FakeClock, Now returns the current deterministic fake time.
	//
	// Callers must not assume that Now provides a globally synchronized timestamp
	// across machines. Distributed ordering, causal consistency, resource
	// versions, and lease ownership must be modeled by higher-level contracts.
	Now() time.Time

	// Since returns the duration elapsed since t according to this clock.
	//
	// Since is the preferred operation for elapsed-time measurements because it
	// allows real implementations to preserve the standard library's monotonic
	// time behavior when t was obtained from the same runtime process.
	//
	// Callers must not rely on monotonic elapsed-time semantics after a time.Time
	// value has been serialized, persisted, transmitted over the network, or
	// reconstructed from an API object. Serialized timestamps are wall-clock
	// values; process-local monotonic information is not a wire-format contract.
	Since(t time.Time) time.Duration
}

// Clock provides runtime time reads and waiting primitives.
//
// Clock extends PassiveClock with the operations required by components that own
// time-based runtime behavior. Use Clock for code that creates timers, tickers,
// sleeps, retry waits, dispatch timeouts, controller loops, heartbeat loops, or
// lease-monitoring loops.
//
// Typical Clock users include:
//
//   - scheduler loops;
//   - adaptive controllers;
//   - worker heartbeat loops;
//   - dispatch timeout monitors;
//   - retry and backoff executors;
//   - lease reapers;
//   - queue wake-up loops.
//
// Components that only need Now or Since should depend on PassiveClock instead
// of Clock. Keeping the narrower interface makes ownership clearer and prevents
// read-only code from depending on blocking or loop-driving primitives.
//
// Clock is intentionally limited to physical/runtime time. It must not grow into
// a general timeutils package. Retry policy, rate limiting, EWMA windows,
// deadline semantics, lease ownership, resource versions, logical revisions, and
// distributed clocks belong to higher-level packages.
type Clock interface {
	PassiveClock

	// After waits for the duration to elapse and then delivers the current clock
	// time on the returned channel.
	//
	// For RealClock, After follows the behavior of time.After. For FakeClock,
	// delivery is controlled by explicit fake-time advancement.
	//
	// After is useful for simple one-shot waits. Components that need to cancel,
	// stop, or reset a wait should use NewTimer instead.
	After(d time.Duration) <-chan time.Time

	// NewTimer creates a one-shot timer controlled by this clock.
	//
	// Timers are appropriate for cancelable or resettable one-shot waits such as
	// dispatch timeouts, retry delays, lease deadlines, and controller cooldowns.
	//
	// The returned Timer belongs to the same clock that created it. A Timer
	// created by a FakeClock must advance only when that FakeClock advances.
	NewTimer(d time.Duration) Timer

	// NewTicker creates a periodic ticker controlled by this clock.
	//
	// Tickers are appropriate for runtime loops that need repeated wake-ups, such
	// as scheduler ticks, heartbeat loops, controller sampling loops, and lease
	// reaper loops.
	//
	// The returned Ticker belongs to the same clock that created it. A Ticker
	// created by a FakeClock must tick only when that FakeClock advances.
	NewTicker(d time.Duration) Ticker

	// Sleep blocks until the duration has elapsed according to this clock.
	//
	// Sleep is convenient for simple blocking waits, but long-running components
	// should prefer timers or tickers when they need cancellation, reset behavior,
	// shutdown coordination, or explicit lifecycle ownership.
	//
	// For FakeClock, Sleep blocks until another goroutine advances fake time far
	// enough to release the sleeper.
	Sleep(d time.Duration)
}
