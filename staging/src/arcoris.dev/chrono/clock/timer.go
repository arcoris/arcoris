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

// Timer is a clock-controlled one-shot timer.
//
// Timer is the cancelable and resettable one-shot waiting primitive used by
// Clock implementations. It abstracts time.Timer so runtime components can use
// the same code with RealClock in production and FakeClock in deterministic
// tests.
//
// Use Timer when a component needs explicit ownership over a one-shot wait:
//
//   - dispatch timeout;
//   - retry delay;
//   - lease deadline;
//   - controller cooldown;
//   - queue wake-up deadline;
//   - shutdown-aware wait;
//   - owner-controlled timeout reset.
//
// For simple fire-and-forget waits, Clock.After may be enough. For repeated
// periodic wake-ups, use Ticker instead.
//
// A Timer belongs to the Clock that created it. A Timer created by RealClock is
// driven by real runtime time. A Timer created by FakeClock is driven only by
// explicit fake-time advancement.
//
// Timer does not define deadline, retry, lease, admission, or scheduling policy.
// It is only a low-level waiting primitive. Higher-level packages own the
// meaning of the timeout and the action taken when the timer fires.
//
// Implementations must be safe for ordinary timer ownership patterns where one
// goroutine owns Stop and Reset while another goroutine may receive from C.
// Implementations are not required to make arbitrary concurrent Reset/Stop calls
// from multiple owners deterministic. Components that need multi-owner timer
// mutation must provide their own synchronization.
type Timer interface {
	// C returns the timer's delivery channel.
	//
	// The channel receives the clock time when the timer fires. The returned
	// channel is receive-only to callers; timer implementations retain ownership
	// of delivery.
	//
	// Callers must not assume that the channel is closed after Stop or after the
	// timer fires. This mirrors the standard library timer model and avoids using
	// channel close as a lifecycle signal.
	C() <-chan time.Time

	// Stop prevents the timer from firing if it is still active.
	//
	// Stop returns true when the timer was active and was stopped by this call. It
	// returns false when the timer had already fired, had already been stopped, or
	// was otherwise no longer active.
	//
	// Stop does not drain the delivery channel. Callers that need strict ownership
	// over a previously delivered value must coordinate channel reads according to
	// the timer implementation and their component lifecycle.
	Stop() bool

	// Reset changes the timer to fire after d according to the owning clock.
	//
	// Reset returns true when the timer was active before the reset and false when
	// it was inactive, stopped, or already fired before the reset. In either case,
	// a successful Reset schedules the timer again using the new duration.
	//
	// For RealClock-backed timers, Reset follows the standard library timer
	// semantics. For FakeClock-backed timers, Reset uses the fake clock's current
	// time as the base for the new deadline.
	//
	// Components should avoid using Reset to hide unclear ownership. If timer
	// state is shared across goroutines, protect Reset/Stop/read coordination with
	// component-level synchronization.
	Reset(d time.Duration) bool
}
