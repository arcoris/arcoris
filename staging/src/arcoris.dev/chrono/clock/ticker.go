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

// Ticker is a clock-controlled periodic wake-up source.
//
// Ticker is the periodic waiting primitive used by Clock implementations. It
// abstracts time.Ticker so runtime loops can use the same code with RealClock in
// production and FakeClock in deterministic tests.
//
// Use Ticker when a component owns a repeated time-based loop:
//
//   - scheduler dispatch loop;
//   - adaptive controller loop;
//   - metrics sampling loop;
//   - worker heartbeat loop;
//   - lease reaper loop;
//   - retry queue wake-up loop;
//   - background reconciliation loop.
//
// Ticker exposes Reset because ARCORIS controllers may need to change loop
// cadence without rebuilding the surrounding component. Typical examples include
// adaptive sampling intervals, overload recovery probes, dynamic heartbeat
// intervals, and queue wake-up loops whose period changes with runtime pressure.
//
// Reset is still a low-level lifecycle operation. It must not encode admission,
// scheduling, rate-limit, retry, lease, or controller policy. Higher layers own
// the policy that decides when a ticker interval should change.
//
// A Ticker belongs to the Clock that created it. A Ticker created by RealClock is
// driven by real runtime time. A Ticker created by FakeClock is driven only by
// explicit fake-time advancement.
//
// Implementations must make Stop and Reset safe for ordinary ticker ownership
// patterns where one goroutine owns the ticker lifecycle and another goroutine
// may receive from C. Implementations are not required to make arbitrary
// concurrent Reset/Stop calls from multiple owners deterministic. Components
// that need multi-owner lifecycle mutation must provide their own
// synchronization.
type Ticker interface {
	// C returns the ticker's delivery channel.
	//
	// The channel receives the clock time for each tick. The returned channel is
	// receive-only to callers; ticker implementations retain ownership of
	// delivery.
	//
	// Callers must not assume that the channel is closed after Stop. This mirrors
	// the standard library ticker model and avoids using channel close as a
	// lifecycle signal.
	C() <-chan time.Time

	// Stop turns off the ticker.
	//
	// Stop prevents future ticks from being delivered. Stop does not close the
	// delivery channel. A tick that has already been delivered or is already
	// available in the channel may still be observed by a receiver.
	//
	// Stop should be called when the owning component no longer needs the ticker,
	// especially in long-running controllers, workers, reapers, and background
	// loops.
	Stop()

	// Reset changes the ticker period.
	//
	// Reset schedules future ticks using the new duration according to the owning
	// clock. For RealClock-backed tickers, Reset follows the standard library
	// ticker reset semantics. For FakeClock-backed tickers, Reset uses the fake
	// clock's current time as the base for the next tick.
	//
	// The duration must be positive. Implementations must reject non-positive
	// durations consistently with time.Ticker.Reset and time.NewTicker behavior.
	//
	// Reset does not drain the delivery channel. A tick that was already
	// delivered before Reset may still be observed by a receiver. Components that
	// require strict tick ownership must coordinate receive, Stop, and Reset at
	// the component level.
	Reset(d time.Duration)
}
