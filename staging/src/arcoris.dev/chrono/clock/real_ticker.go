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

// realTicker adapts the Go standard library ticker to the clock.Ticker
// contract.
//
// realTicker is intentionally unexported. Callers should create production
// tickers through RealClock.NewTicker instead of constructing adapter values
// directly.
//
// The adapter does not add scheduling, retry, controller, lease, or rate-limit
// policy on top of time.Ticker. It preserves the standard library ticker
// lifecycle as closely as possible:
//
//   - C exposes the ticker delivery channel;
//   - Stop delegates to (*time.Ticker).Stop;
//   - Reset delegates to (*time.Ticker).Reset.
//
// realTicker does not close or drain the ticker channel. Channel ownership and
// receive coordination belong to the component that owns the ticker lifecycle.
//
// A realTicker value owns exactly one *time.Ticker. It must not be copied after
// construction. Copying the adapter would duplicate the wrapper around the same
// underlying ticker and obscure lifecycle ownership.
type realTicker struct {
	ticker *time.Ticker
}

// C returns the ticker's delivery channel.
//
// The returned channel is receive-only for callers. The channel is owned by the
// underlying *time.Ticker and is not closed by Stop.
//
// C intentionally returns the same channel on every call.
func (t *realTicker) C() <-chan time.Time {
	return t.ticker.C
}

// Stop turns off the ticker.
//
// Stop delegates directly to (*time.Ticker).Stop. After Stop, no future ticks
// are delivered by the underlying ticker.
//
// Stop does not close the delivery channel. A tick that was already delivered or
// was already available in the channel before Stop may still be observed by a
// receiver. Components that need strict tick ownership must coordinate receives
// and Stop at the component level.
//
// Stop is idempotent according to the standard library ticker lifecycle: calling
// it more than once is safe.
func (t *realTicker) Stop() {
	t.ticker.Stop()
}

// Reset changes the ticker period.
//
// Reset delegates directly to (*time.Ticker).Reset and therefore follows the Go
// standard library ticker semantics.
//
// The duration must be positive. The standard library panics when Reset is
// called with a non-positive duration; realTicker intentionally preserves that
// behavior instead of wrapping or translating the panic.
//
// Reset does not drain the delivery channel. A tick that was already delivered
// before Reset may still be observed by a receiver. Components that need strict
// tick ownership must coordinate receives, Stop, and Reset at the component
// level.
func (t *realTicker) Reset(d time.Duration) {
	t.ticker.Reset(d)
}
