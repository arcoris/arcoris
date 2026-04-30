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

// RealClock is the production Clock implementation backed by the Go standard
// library.
//
// RealClock is stateless and zero-value usable. It is the default clock for
// production ARCORIS components that need runtime time, timers, tickers, sleeps,
// controller loop cadence, dispatch timeouts, retry waits, heartbeat intervals,
// or lease-monitoring intervals.
//
// RealClock deliberately adds no admission, scheduling, retry, rate-limit,
// lease, controller, metrics, or distributed-systems policy on top of package
// time. It only adapts the Go standard library time primitives to the clock
// package contracts.
//
// RealClock reads the local process clock. It does not provide cluster-wide time
// synchronization, distributed causality, resource versions, generations, lease
// epochs, fencing tokens, or hybrid logical time. Those concepts belong to
// higher-level ARCORIS packages.
//
// RealClock values are safe to copy because they contain no state. Timers and
// tickers created by RealClock are independent runtime objects and have their
// own lifecycle ownership.
type RealClock struct{}

// Now returns the current local wall-clock time.
//
// Now delegates directly to time.Now.
//
// The returned value may include Go's process-local monotonic clock reading.
// That monotonic component can be used by time.Since and time.Sub while the value
// remains in the same process, but it is not preserved when a time.Time is
// serialized, persisted, transmitted over the network, or reconstructed from an
// API object.
//
// Callers must not treat RealClock.Now as globally synchronized cluster time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// Since returns the duration elapsed since t.
//
// Since delegates directly to time.Since. When t was obtained from time.Now in
// the same process and still carries a monotonic clock reading, the standard
// library uses that monotonic reading for elapsed-time calculation.
//
// Use Since for local elapsed-time measurements such as request latency, retry
// delay, controller cooldown, heartbeat age, queue wait age, and timeout checks.
//
// Do not rely on monotonic semantics for timestamps that were serialized,
// persisted, loaded from an API object, or received from another process. Those
// values are wall-clock values for API, audit, and interoperability purposes.
func (RealClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// After waits for d to elapse and then delivers the current time on the returned
// channel.
//
// After delegates directly to time.After. It is appropriate for simple
// fire-and-forget waits where the caller does not need to stop or reset the
// wait.
//
// Components that need cancellation, reset behavior, shutdown coordination, or
// explicit lifecycle ownership should use NewTimer instead.
//
// The duration semantics are the standard library semantics. Non-positive
// durations make the returned channel ready as soon as the runtime can deliver
// the value.
func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// NewTimer creates a one-shot production timer.
//
// NewTimer delegates to time.NewTimer and wraps the result in the package's
// Timer contract through the unexported realTimer adapter.
//
// Use NewTimer for dispatch timeouts, retry delays, controller cooldowns, queue
// wake-up deadlines, lease checks, heartbeat deadlines, and other one-shot waits
// that need explicit Stop or Reset ownership.
//
// The returned Timer is backed by a distinct *time.Timer. The caller owns its
// lifecycle and is responsible for coordinating Stop, Reset, and receives from
// the timer channel according to the component's concurrency model.
func (RealClock) NewTimer(d time.Duration) Timer {
	return &realTimer{
		timer: time.NewTimer(d),
	}
}

// NewTicker creates a production ticker.
//
// NewTicker delegates to time.NewTicker and wraps the result in the package's
// Ticker contract through the unexported realTicker adapter.
//
// Use NewTicker for scheduler loops, adaptive controller loops, metrics sampling
// loops, worker heartbeat loops, lease reaper loops, retry wake-up loops, and
// other repeated runtime wake-ups.
//
// The duration must be positive. time.NewTicker panics for non-positive
// durations; RealClock intentionally preserves that standard library behavior
// instead of translating it.
//
// The returned Ticker is backed by a distinct *time.Ticker. The caller owns its
// lifecycle and should call Stop when the ticker is no longer needed.
func (RealClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{
		ticker: time.NewTicker(d),
	}
}

// Sleep blocks the current goroutine for at least d.
//
// Sleep delegates directly to time.Sleep. It is convenient for simple blocking
// waits, but long-running runtime components should prefer timers or tickers when
// they need cancellation, reset behavior, shutdown coordination, or explicit
// lifecycle ownership.
//
// The duration semantics are the standard library semantics. Non-positive
// durations return immediately.
func (RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
