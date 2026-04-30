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

// realTimer adapts the Go standard library timer to the clock.Timer contract.
//
// realTimer is intentionally unexported. Callers should create production
// timers through RealClock.NewTimer instead of constructing adapter values
// directly.
//
// The adapter does not add policy on top of time.Timer. It preserves the
// standard library timer lifecycle as closely as possible:
//
//   - C exposes the timer delivery channel;
//   - Stop delegates to (*time.Timer).Stop;
//   - Reset delegates to (*time.Timer).Reset.
//
// realTimer does not drain the timer channel after Stop or Reset. Channel
// draining is an ownership concern of the component using the timer, because the
// correct behavior depends on whether another goroutine may already be receiving
// from the channel.
//
// A realTimer value owns exactly one *time.Timer. It must not be copied after
// construction. Copying the adapter would duplicate the wrapper around the same
// underlying timer and obscure lifecycle ownership.
type realTimer struct {
	timer *time.Timer
}

// C returns the timer's delivery channel.
//
// The returned channel is receive-only for callers. The channel is owned by the
// underlying *time.Timer and is not closed by Stop.
//
// C intentionally returns the same channel on every call.
func (t *realTimer) C() <-chan time.Time {
	return t.timer.C
}

// Stop prevents the timer from firing if it is still active.
//
// Stop delegates directly to (*time.Timer).Stop. It returns true if the timer was
// active and was stopped by this call. It returns false if the timer had already
// expired, had already been stopped, or was otherwise inactive.
//
// Stop does not drain the timer channel. Components that need strict ownership
// over a previously delivered timer value must coordinate channel reads at the
// component level.
func (t *realTimer) Stop() bool {
	return t.timer.Stop()
}

// Reset changes the timer to expire after d.
//
// Reset delegates directly to (*time.Timer).Reset and therefore follows the Go
// standard library semantics for active, stopped, and expired timers.
//
// Reset returns true if the timer was active before the reset and false if the
// timer was inactive, stopped, or already expired before the reset.
//
// Reset does not drain the timer channel. If a component resets a timer that may
// already have delivered a value, that component is responsible for coordinating
// receives, Stop, and Reset according to its own ownership model.
func (t *realTimer) Reset(d time.Duration) bool {
	return t.timer.Reset(d)
}
