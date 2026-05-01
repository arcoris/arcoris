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

const (
	// errFakeClockBackwardSet is the panic value used when FakeClock.Set attempts
	// to move fake time backwards.
	//
	// FakeClock is intentionally monotonic. Allowing backwards movement would
	// make timer, ticker, waiter, retry, lease, and controller-loop tests harder
	// to reason about because already-fired deadlines could become "future"
	// deadlines again.
	errFakeClockBackwardSet = "clock.FakeClock: cannot move time backwards"

	// errFakeClockNegativeStep is the panic value used when FakeClock.Step is
	// called with a negative duration.
	//
	// FakeClock advances monotonically. Tests that need a different starting time
	// should create a new FakeClock or call Set with a future timestamp.
	errFakeClockNegativeStep = "clock.FakeClock: negative step"
)

// Now returns the current fake time.
//
// Now does not observe real runtime time. The returned value changes only when
// Set or Step advances the fake clock.
//
// Now is safe for concurrent use.
func (c *FakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

// Since returns the duration elapsed since t according to the fake clock.
//
// Since is computed as:
//
//	FakeClock.Now().Sub(t)
//
// Since does not observe real runtime time. If the fake clock has not advanced,
// Since returns a duration based on the current fake time stored in the clock.
//
// Since is safe for concurrent use.
func (c *FakeClock) Since(t time.Time) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now.Sub(t)
}

// Set moves the fake clock to next.
//
// Set is monotonic. It panics if next is before the current fake time. Moving
// fake time backwards would make deadline-based tests ambiguous and could
// invalidate already-fired timer, ticker, or waiter state.
//
// Set may be called with the current fake time. This is useful for forcing
// immediate due-delivery processing for waiters, timers, or tickers whose
// deadline is already equal to the current fake time.
//
// Set collects all due deliveries while holding the fake clock lock, releases
// the lock, and then delivers to user-visible channels. Delivery must not happen
// while holding the lock because channel sends can interact with receiver
// scheduling and should not block fake clock state mutation.
func (c *FakeClock) Set(next time.Time) {
	deliveries := c.set(next)
	deliveries.deliver()
}

// Step advances the fake clock by d.
//
// Step panics if d is negative. Step(0) is allowed and processes deliveries that
// are already due at the current fake time.
//
// Like Set, Step collects due deliveries under the fake clock lock and performs
// channel delivery after releasing the lock.
func (c *FakeClock) Step(d time.Duration) {
	deliveries := c.step(d)
	deliveries.deliver()
}

// set updates the fake time and returns due deliveries.
//
// The returned deliveries must be delivered after c.mu has been released.
func (c *FakeClock) set(next time.Time) fakeClockDeliveries {
	c.mu.Lock()
	defer c.mu.Unlock()

	if next.Before(c.now) {
		panic(errFakeClockBackwardSet)
	}

	c.now = next
	return c.collectDueDeliveriesLocked()
}

// step advances the fake time and returns due deliveries.
//
// The returned deliveries must be delivered after c.mu has been released.
func (c *FakeClock) step(d time.Duration) fakeClockDeliveries {
	if d < 0 {
		panic(errFakeClockNegativeStep)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)
	return c.collectDueDeliveriesLocked()
}

// fakeClockDeliveries groups all user-visible channel deliveries produced by one
// fake clock advancement operation.
//
// FakeClock advancement is split into two phases:
//
//   - collect due deliveries while holding FakeClock.mu;
//   - deliver to user-visible channels after releasing FakeClock.mu.
//
// This prevents internal fake clock state from being held hostage by receiver
// scheduling.
type fakeClockDeliveries struct {
	waiters []fakeWaiterDelivery
	timers  []fakeTimerDelivery
	tickers []fakeTickerDelivery
}

// collectDueDeliveriesLocked collects all waiters, timers, and tickers that are
// due at the current fake time.
//
// The caller must hold c.mu.
//
// The method does not deliver to channels. It only builds a delivery plan that
// can be executed after the fake clock lock is released.
func (c *FakeClock) collectDueDeliveriesLocked() fakeClockDeliveries {
	return fakeClockDeliveries{
		waiters: c.collectDueWaiterDeliveriesLocked(),
		timers:  c.collectDueTimerDeliveriesLocked(),
		tickers: c.collectDueTickerDeliveriesLocked(),
	}
}

// deliver performs all planned user-visible channel deliveries.
//
// Delivery order is intentionally stable:
//
//   - waiters;
//   - timers;
//   - tickers.
//
// Within each source, due entries are delivered by scheduled deadline and then
// by registration sequence. The order prevents Sleep/After waiters from being
// starved behind periodic tickers in tests and keeps one-shot delivery behavior
// easier to reason about. Higher-level components must not rely on this as a
// distributed ordering contract; it is only an implementation detail of one fake
// clock advancement.
func (d fakeClockDeliveries) deliver() {
	deliverFakeWaiterDeliveries(d.waiters)
	deliverFakeTimerDeliveries(d.timers)
	deliverFakeTickerDeliveries(d.tickers)
}
