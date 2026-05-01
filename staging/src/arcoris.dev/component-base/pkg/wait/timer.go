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

package wait

import (
	"context"
	"time"
)

// Timer is an owner-controlled one-shot runtime timer.
//
// Timer is a low-level primitive for code that needs explicit lifecycle control
// over a single real-time wait. It wraps the Go standard library timer model and
// adds wait-owned context classification through Wait.
//
// Use Timer when a component needs to create a timer once and then stop, drain,
// reset, or wait on it under an owning context. Typical users include fixed
// runtime loops, controller cooldowns, queue wake-up delays, retry executors,
// dispatch timeout monitors, shutdown-aware waits, and other packages that need
// more lifecycle control than Delay provides.
//
// Timer does not implement retry policy, backoff growth, jitter, rate limiting,
// metrics, tracing, logging, scheduler policy, queue policy, or deadline
// ownership semantics. Higher-level packages own the meaning of the timer and
// the action taken when it fires.
//
// A Timer value owns exactly one *time.Timer. It MUST be created with NewTimer.
// The zero value is not usable because it has no underlying runtime timer.
//
// Timer must not be copied after construction. Copying the value would duplicate
// the wrapper around the same underlying runtime timer and obscure lifecycle
// ownership.
//
// Timer supports the ordinary single-owner pattern used by runtime loops: one
// goroutine owns Stop, StopAndDrain, Reset, and Wait coordination. Timer does not
// make arbitrary concurrent Stop/Reset/Wait calls from multiple owners
// deterministic. Components that share timer mutation across goroutines MUST add
// component-level synchronization.
type Timer struct {
	// timer is the underlying runtime timer owned by this wrapper.
	//
	// The field is intentionally private so the package can enforce stop/drain
	// rules around Wait and Reset while keeping direct channel access available
	// only through C.
	timer *time.Timer
}

// NewTimer returns a new one-shot runtime timer that will fire after duration.
//
// NewTimer follows the useful construction semantics of time.NewTimer:
// non-positive durations create a timer that is ready immediately. Use Delay or
// higher-level loop APIs when non-positive durations should be interpreted as
// validation errors or immediate no-ops.
//
// The returned Timer is not context-bound. Context ownership is applied by Wait,
// which converts wait-owned context stops into ErrInterrupted or ErrTimeout.
func NewTimer(duration time.Duration) *Timer {
	return &Timer{timer: time.NewTimer(duration)}
}

// C returns the timer delivery channel.
//
// The channel receives one time value when the timer fires. The returned channel
// is receive-only for callers. The channel is owned by the underlying runtime
// timer and is not closed by Stop, StopAndDrain, Reset, or Wait.
//
// C is provided for low-level loop integration. Most callers should prefer Wait
// when they need context-aware waiting and wait-owned error classification.
//
// C panics when called on a nil or zero-value Timer.
func (t *Timer) C() <-chan time.Time {
	t.requireUsable()

	return t.timer.C
}

// Wait blocks until the timer fires or ctx stops.
//
// Wait returns nil when the timer fires first. If ctx stops first, Wait stops and
// drains the timer, then returns a wait-owned error classified as ErrInterrupted
// or ErrTimeout. Cancellation and deadline causes are preserved by the returned
// error.
//
// If ctx is already stopped before Wait starts, Wait still stops and drains the
// timer before returning the wait-owned context error. This keeps timer
// ownership explicit and prevents a timer that the caller intended to wait on
// from continuing to run after the wait has already been abandoned.
//
// If the timer fires and ctx stops at approximately the same time, either result
// may win according to the select race. Callers that require a stricter policy
// must encode it at a higher level.
//
// Wait does not recover panics and does not add retry, backoff, jitter, logging,
// tracing, or metrics policy.
//
// Wait panics when ctx is nil or when called on a nil or zero-value Timer.
func (t *Timer) Wait(ctx context.Context) error {
	requireContext(ctx)
	t.requireUsable()

	if err := contextStopError(ctx); err != nil {
		t.StopAndDrain()
		return err
	}

	select {
	case <-ctx.Done():
		t.StopAndDrain()
		return contextStopError(ctx)
	case <-t.timer.C:
		return nil
	}
}

// Stop prevents the timer from firing if it is still active.
//
// Stop returns true when the timer was active and was stopped by this call. It
// returns false when the timer had already fired, had already been stopped, or
// was otherwise inactive.
//
// Stop intentionally mirrors (*time.Timer).Stop and does not drain the delivery
// channel. Use StopAndDrain when the caller owns the timer channel and wants to
// remove a pending delivered value before reusing or abandoning the timer.
//
// Stop panics when called on a nil or zero-value Timer.
func (t *Timer) Stop() bool {
	t.requireUsable()

	return t.timer.Stop()
}

// StopAndDrain stops the timer and drains one pending delivery value, if any.
//
// StopAndDrain returns true when the timer was active and was stopped by this
// call. It returns false when the timer had already fired, had already been
// stopped, or was otherwise inactive. The return value is the result of Stop;
// draining a pending value does not change it.
//
// StopAndDrain is the safer shutdown/reuse primitive for single-owner timer
// code. It prevents a previously delivered value from being observed after the
// caller believes the timer has been stopped or before Reset schedules the next
// wait.
//
// StopAndDrain assumes the caller owns timer-channel receives. If another
// goroutine may be receiving from C at the same time, draining is an ownership
// race and must be coordinated by the component using the timer.
//
// StopAndDrain panics when called on a nil or zero-value Timer.
func (t *Timer) StopAndDrain() bool {
	t.requireUsable()

	return stopAndDrainRuntimeTimer(t.timer)
}

// Reset stops, drains, and reschedules the timer to fire after duration.
//
// Reset returns true when the timer was active before this call stopped it. It
// returns false when the timer had already fired, had already been stopped, or
// was otherwise inactive before being rescheduled.
//
// Unlike direct (*time.Timer).Reset usage, this method performs a stop/drain
// step before rescheduling. That makes the method suitable for the package's
// single-owner runtime-loop pattern, where the owner wants to avoid stale timer
// values across iterations.
//
// Non-positive durations are allowed and schedule an immediately ready timer.
// Higher-level APIs that would busy-loop with non-positive intervals must reject
// those values before calling Reset.
//
// Reset panics when called on a nil or zero-value Timer.
func (t *Timer) Reset(duration time.Duration) bool {
	t.requireUsable()

	wasActive := stopAndDrainRuntimeTimer(t.timer)
	t.timer.Reset(duration)
	return wasActive
}

// stopAndDrainRuntimeTimer stops timer and drains a pending value when one is
// available.
//
// The helper centralizes the Stop+nonblocking-drain pattern used by Timer,
// Delay, and loop code. The nonblocking drain avoids hanging when the timer did
// not deliver a value or when the value was already consumed by the owner.
func stopAndDrainRuntimeTimer(timer *time.Timer) bool {
	stopped := timer.Stop()
	if stopped {
		return true
	}

	select {
	case <-timer.C:
	default:
	}

	return false
}
