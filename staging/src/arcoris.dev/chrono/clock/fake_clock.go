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

import (
	"sync"
	"time"
)

// FakeClock is a deterministic Clock implementation for tests.
//
// FakeClock does not observe real runtime time. Its current time changes only
// when test code explicitly moves it with Set or Step. This makes time-dependent
// ARCORIS components testable without real sleeps, timing races, or flaky
// duration assumptions.
//
// FakeClock is intended for tests of components such as:
//
//   - scheduler loops;
//   - admission and retry logic;
//   - queue wake-up behavior;
//   - dispatch timeout handling;
//   - worker heartbeat loops;
//   - lease reapers;
//   - adaptive controller ticks;
//   - cooldown and backoff behavior.
//
// FakeClock is not a distributed clock. It does not provide causal ordering,
// resource versions, lease epochs, fencing tokens, hybrid logical time, or
// cluster-wide time synchronization. Those concepts belong to higher-level
// packages.
//
// A FakeClock owns four kinds of mutable state:
//
//   - the current fake time;
//   - pending one-shot waiters used by After and Sleep;
//   - fake timers created by NewTimer;
//   - fake tickers created by NewTicker.
//
// It also keeps a small internal sequence counter so due deliveries can be
// ordered deterministically when several waiters, timers, or tickers become due
// during the same fake-time advancement.
//
// All mutable state is protected by mu. Implementation code MUST NOT read or
// mutate now, waiters, timers, tickers, or nextSequence without holding mu.
//
// Delivery to user-visible channels MUST NOT happen while holding mu. Channel
// sends can interact with receiver scheduling and may otherwise turn fake-clock
// advancement into a blocking operation. Advancement code should collect due
// deliveries under mu, release mu, and then deliver to channels.
//
// FakeClock is safe for concurrent use by tests and by the components under
// test. It must not be copied after first use. Copying a live FakeClock would
// copy its mutex and split ownership of registered timers, tickers, and waiters.
//
// Use NewFakeClock to create a clock with an explicit initial time. The zero
// value is intentionally not the preferred construction path because tests
// should choose their initial time deliberately.
type FakeClock struct {
	// mu protects all mutable fake-clock state.
	//
	// Code that holds mu may inspect and update now, waiters, timers, and
	// tickers. Code must not send on user-visible timer, ticker, or waiter
	// channels while holding mu.
	mu sync.Mutex

	// now is the current fake time.
	//
	// now is advanced only by Set and Step. Real runtime time never updates this
	// field.
	now time.Time

	// waiters contains pending one-shot waits registered by After and Sleep.
	//
	// Waiter lifecycle is owned by fake_waiter.go.
	waiters map[*fakeWaiter]struct{}

	// timers contains active fake timers created by NewTimer.
	//
	// Timer lifecycle is owned by fake_timer.go.
	timers map[*fakeTimer]struct{}

	// tickers contains active fake tickers created by NewTicker.
	//
	// Ticker lifecycle is owned by fake_ticker.go.
	tickers map[*fakeTicker]struct{}

	// nextSequence assigns deterministic tie-break positions to scheduled fake
	// waiters, timers, and tickers.
	//
	// It is advanced only while mu is held.
	nextSequence uint64
}

// NewFakeClock creates a FakeClock initialized to now.
//
// The supplied time is used exactly as provided. Callers that need stable,
// deterministic test output should usually pass an explicit time.Date value
// instead of time.Now.
//
// NewFakeClock initializes the internal registries used by fake waiters, timers,
// and tickers. The returned clock is ready for concurrent use by test code and
// the component under test.
//
// The returned FakeClock must not be copied after first use.
func NewFakeClock(now time.Time) *FakeClock {
	return &FakeClock{
		now:     now,
		waiters: make(map[*fakeWaiter]struct{}),
		timers:  make(map[*fakeTimer]struct{}),
		tickers: make(map[*fakeTicker]struct{}),
	}
}

// nextSequenceLocked returns the next deterministic delivery tie-break value.
//
// The caller must hold c.mu. Sequence values are internal implementation detail;
// callers observe their effect only through stable fake-time delivery order.
func (c *FakeClock) nextSequenceLocked() uint64 {
	c.nextSequence++
	return c.nextSequence
}
