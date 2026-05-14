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

// After waits for the duration to elapse according to the fake clock and then
// delivers the current fake time on the returned channel.
//
// After does not observe real runtime time. The returned channel becomes ready
// only when the owning FakeClock is advanced far enough with Set or Step.
//
// Non-positive durations are treated as immediately due. In that case, After
// schedules delivery using the current fake time and returns a channel that is
// ready without requiring an additional fake-clock advancement.
//
// The returned channel has capacity one. Delivery is non-blocking and is
// performed outside the fake clock lock.
//
// After does not close the returned channel. Channel close is intentionally not
// used as a waiter lifecycle signal.
func (c *FakeClock) After(d time.Duration) <-chan time.Time {
	waiter, delivery, ok := c.newWaiter(d)
	if ok {
		delivery.deliver()
	}

	return waiter.ch
}

// Sleep blocks until the duration has elapsed according to the fake clock.
//
// Sleep does not observe real runtime time. It blocks until another goroutine
// advances the same FakeClock far enough with Set or Step.
//
// Non-positive durations return after immediate fake-time delivery.
//
// Sleep is intended for tests of components that use Clock.Sleep. Tests should
// generally prefer timers or explicit fake-clock advancement when they need
// precise lifecycle control.
func (c *FakeClock) Sleep(d time.Duration) {
	<-c.After(d)
}

// HasWaiters reports whether the fake clock has pending After or Sleep waiters.
//
// HasWaiters is primarily a test coordination helper. It lets tests avoid real
// sleeps when they need to wait until a goroutine has actually registered a
// fake-time wait.
//
// HasWaiters reports only pending fake waiters. It does not report active
// timers or tickers.
func (c *FakeClock) HasWaiters() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.waiters) > 0
}

// fakeWaiter is a one-shot fake-time waiter used by After and Sleep.
//
// A waiter is active until fake time reaches or passes its deadline. Once due,
// it is removed from the owning FakeClock waiter registry and delivered exactly
// once.
//
// fakeWaiter is intentionally not exported. Public code observes only the
// receive-only channel returned by After or the blocking behavior of Sleep.
//
// All fakeWaiter fields are protected by the owning FakeClock mutex while the
// waiter is registered in FakeClock.waiters.
type fakeWaiter struct {
	deadline time.Time
	sequence uint64

	ch     chan time.Time
	active bool
}

// fakeWaiterDelivery describes one waiter delivery that must be performed after
// the fake clock lock has been released.
//
// FakeClock advancement code must not send on user-visible channels while
// holding FakeClock.mu. Channel sends can interact with receiver scheduling and
// must stay outside the internal state lock.
type fakeWaiterDelivery struct {
	ch    chan time.Time
	value time.Time
}

// deliver sends the fake time to the waiter channel.
//
// The waiter channel has capacity one, and each waiter is delivered at most
// once. This send should normally succeed immediately. The non-blocking form is
// used as a defensive measure so delivery can never block fake-clock
// advancement.
func (d fakeWaiterDelivery) deliver() {
	select {
	case d.ch <- d.value:
	default:
	}
}

// newWaiter registers a one-shot waiter and returns an immediate delivery when
// the waiter is already due.
//
// The returned delivery, when ok is true, must be delivered after c.mu has been
// released.
func (c *FakeClock) newWaiter(d time.Duration) (*fakeWaiter, fakeWaiterDelivery, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureWaiterStoreLocked()

	waiter := &fakeWaiter{
		deadline: c.now.Add(d),
		sequence: c.nextSequenceLocked(),
		ch:       make(chan time.Time, 1),
		active:   true,
	}

	if !c.now.Before(waiter.deadline) {
		waiter.active = false
		return waiter, fakeWaiterDelivery{
			ch:    waiter.ch,
			value: c.now,
		}, true
	}

	c.waiters[waiter] = struct{}{}

	return waiter, fakeWaiterDelivery{}, false
}

// ensureWaiterStoreLocked initializes the fake waiter registry.
//
// The caller must hold c.mu.
func (c *FakeClock) ensureWaiterStoreLocked() {
	if c.waiters == nil {
		c.waiters = make(map[*fakeWaiter]struct{})
	}
}

// collectDueWaiterDeliveriesLocked collects all waiters that are due at the
// current fake time.
//
// The caller must hold c.mu.
//
// Due waiters are removed from the waiter registry and marked inactive before
// delivery is returned. Delivery itself must happen after c.mu has been released.
func (c *FakeClock) collectDueWaiterDeliveriesLocked() []fakeWaiterDelivery {
	if len(c.waiters) == 0 {
		return nil
	}

	due := make([]*fakeWaiter, 0, len(c.waiters))

	for waiter := range c.waiters {
		if waiter.isDueLocked(c.now) {
			due = append(due, waiter)
		}
	}

	if len(due) == 0 {
		return nil
	}

	sort.Slice(due, func(i, j int) bool {
		if due[i].deadline.Equal(due[j].deadline) {
			return due[i].sequence < due[j].sequence
		}

		return due[i].deadline.Before(due[j].deadline)
	})

	deliveries := make([]fakeWaiterDelivery, 0, len(due))

	for _, waiter := range due {
		delivery, _ := waiter.collectDueDeliveryLocked(c.now)

		delete(c.waiters, waiter)
		deliveries = append(deliveries, delivery)
	}

	return deliveries
}

// isDueLocked reports whether the waiter should be released at now.
//
// The caller must hold the owning FakeClock mutex.
func (w *fakeWaiter) isDueLocked(now time.Time) bool {
	return w.active && !now.Before(w.deadline)
}

// collectDueDeliveryLocked returns one delivery when the waiter is due.
//
// The caller must hold the owning FakeClock mutex.
func (w *fakeWaiter) collectDueDeliveryLocked(now time.Time) (fakeWaiterDelivery, bool) {
	if !w.isDueLocked(now) {
		return fakeWaiterDelivery{}, false
	}

	w.active = false

	return fakeWaiterDelivery{
		ch:    w.ch,
		value: now,
	}, true
}

// deliverFakeWaiterDeliveries delivers all collected waiter deliveries.
//
// FakeClock advancement code should call this after releasing FakeClock.mu.
func deliverFakeWaiterDeliveries(deliveries []fakeWaiterDelivery) {
	for _, delivery := range deliveries {
		delivery.deliver()
	}
}
