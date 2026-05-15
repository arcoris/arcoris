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
	"sort"
	"time"
)

const (
	// errFakeTickerNonPositiveInterval is the panic value used when a fake ticker
	// is created or reset with a non-positive interval.
	//
	// The fake ticker intentionally follows the standard library ticker contract:
	// time.NewTicker and (*time.Ticker).Reset reject non-positive durations.
	errFakeTickerNonPositiveInterval = "clock.FakeClock: non-positive ticker interval"
)

// NewTicker creates a periodic ticker controlled by the fake clock.
//
// The ticker does not observe real time. It delivers ticks only when the owning
// FakeClock is advanced with Set or Step.
//
// The duration must be positive. NewTicker panics for non-positive durations,
// matching the standard library time.NewTicker behavior.
//
// Fake ticker delivery is intentionally bounded: a ticker delivers at most one
// tick per fake-clock advancement operation, even if Set or Step moves fake time
// across several ticker intervals. This avoids unbounded tick bursts in tests
// and gives controller-loop tests deterministic one-wakeup-per-advance behavior.
//
// A ticker created by FakeClock belongs to that FakeClock. It must not be moved
// between fake clocks.
func (c *FakeClock) NewTicker(d time.Duration) Ticker {
	if d <= 0 {
		panic(errFakeTickerNonPositiveInterval)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureTickerStoreLocked()

	ticker := &fakeTicker{
		clock:    c,
		interval: d,
		next:     c.now.Add(d),
		sequence: c.nextSequenceLocked(),
		ch:       make(chan time.Time, 1),
		active:   true,
	}

	c.tickers[ticker] = struct{}{}

	return ticker
}

// fakeTicker is the FakeClock-backed implementation of Ticker.
//
// fakeTicker is intentionally unexported. Callers create fake tickers through
// FakeClock.NewTicker and interact with them only through the Ticker interface.
//
// All mutable fakeTicker state is protected by the owning FakeClock mutex.
// Methods that mutate ticker lifecycle state must acquire t.clock.mu.
//
// The ticker channel is buffered with capacity one. This mirrors the useful
// behavior of production tickers for controller tests: if a receiver has not
// consumed the previous tick, a later fake-time advancement does not block the
// clock. Fake ticker delivery is best-effort with respect to an already full
// channel; the ticker remains scheduled for future intervals.
type fakeTicker struct {
	clock *FakeClock

	interval time.Duration
	next     time.Time
	sequence uint64

	ch     chan time.Time
	active bool
}

// C returns the ticker delivery channel.
//
// The returned channel is receive-only for callers. It is not closed by Stop.
// This mirrors the standard library ticker model and prevents channel close from
// becoming a lifecycle signal.
//
// C returns the same channel on every call.
func (t *fakeTicker) C() <-chan time.Time {
	return t.ch
}

// Stop turns off the fake ticker.
//
// Stop prevents future ticks from being scheduled or delivered. Stop does not
// close the delivery channel, and it does not remove a tick that has already
// been delivered into the channel before Stop.
//
// Stop is idempotent.
func (t *fakeTicker) Stop() {
	c := t.clock

	c.mu.Lock()
	defer c.mu.Unlock()

	if !t.active {
		return
	}

	t.active = false
	delete(c.tickers, t)
}

// Reset changes the fake ticker interval.
//
// Reset follows the ticker contract used by Ticker: the new duration must be
// positive. Reset panics for non-positive durations, matching the standard
// library ticker reset behavior.
//
// Reset schedules the next tick at:
//
//	FakeClock.Now() + d
//
// Reset reactivates a stopped ticker and registers it again with the owning fake
// clock. Reset does not drain the delivery channel. A tick that was already
// delivered before Reset may still be observed by a receiver.
//
// If a previously delivered tick is still unread, a later tick may be dropped
// because fake ticker delivery is non-blocking. Components that reuse tickers
// and need strict tick ownership must coordinate receives before Reset.
func (t *fakeTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic(errFakeTickerNonPositiveInterval)
	}

	c := t.clock

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureTickerStoreLocked()

	t.interval = d
	t.next = c.now.Add(d)
	t.sequence = c.nextSequenceLocked()

	if !t.active {
		t.active = true
		c.tickers[t] = struct{}{}
	}
}

// fakeTickerDelivery describes one ticker tick that must be delivered after the
// fake clock lock has been released.
//
// FakeClock advancement code must not send on user-visible channels while
// holding FakeClock.mu. Channel sends can interact with receiver scheduling and
// must stay outside the internal state lock.
type fakeTickerDelivery struct {
	ch    chan time.Time
	value time.Time
}

// deliver sends the tick if the ticker channel can accept it.
//
// Delivery is non-blocking. If the receiver has not consumed the previous tick,
// this delivery is dropped. The ticker remains active and future fake-clock
// advancements can deliver future ticks.
func (d fakeTickerDelivery) deliver() {
	select {
	case d.ch <- d.value:
	default:
	}
}

// ensureTickerStoreLocked initializes the fake ticker registry.
//
// The caller must hold c.mu.
func (c *FakeClock) ensureTickerStoreLocked() {
	if c.tickers == nil {
		c.tickers = make(map[*fakeTicker]struct{})
	}
}

// collectDueTickerDeliveriesLocked collects ticker deliveries made due by the
// fake clock's current time.
//
// The caller must hold c.mu.
//
// Each active ticker can contribute at most one delivery per call. If fake time
// advanced across several intervals, the ticker is rescheduled relative to the
// current fake time instead of emitting a burst of missed ticks.
func (c *FakeClock) collectDueTickerDeliveriesLocked() []fakeTickerDelivery {
	if len(c.tickers) == 0 {
		return nil
	}

	due := make([]*fakeTicker, 0, len(c.tickers))

	for ticker := range c.tickers {
		if ticker.isDueLocked(c.now) {
			due = append(due, ticker)
		}
	}

	if len(due) == 0 {
		return nil
	}

	sort.Slice(due, func(i, j int) bool {
		if due[i].next.Equal(due[j].next) {
			return due[i].sequence < due[j].sequence
		}

		return due[i].next.Before(due[j].next)
	})

	deliveries := make([]fakeTickerDelivery, 0, len(due))

	for _, ticker := range due {
		delivery, _ := ticker.collectDueDeliveryLocked(c.now, c.nextSequenceLocked())
		deliveries = append(deliveries, delivery)
	}

	return deliveries
}

// isDueLocked reports whether the ticker should emit a tick at now.
//
// The caller must hold the owning FakeClock mutex.
func (t *fakeTicker) isDueLocked(now time.Time) bool {
	return t.active && !now.Before(t.next)
}

// collectDueDeliveryLocked returns one delivery when the ticker is due.
//
// The caller must hold t.clock.mu.
//
// A due ticker is rescheduled to now + interval. This preserves deterministic
// one-tick-per-advance behavior and avoids unbounded catch-up loops when tests
// advance fake time by a large duration.
func (t *fakeTicker) collectDueDeliveryLocked(now time.Time, seq uint64) (fakeTickerDelivery, bool) {
	if !t.isDueLocked(now) {
		return fakeTickerDelivery{}, false
	}

	t.next = now.Add(t.interval)
	t.sequence = seq

	return fakeTickerDelivery{
		ch:    t.ch,
		value: now,
	}, true
}

// deliverFakeTickerDeliveries delivers all collected ticker ticks.
//
// FakeClock advancement code should call this after releasing FakeClock.mu.
func deliverFakeTickerDeliveries(deliveries []fakeTickerDelivery) {
	for _, delivery := range deliveries {
		delivery.deliver()
	}
}
