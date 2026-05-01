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

// NewTimer creates a one-shot timer controlled by the fake clock.
//
// The timer does not observe real runtime time. It fires only when the owning
// FakeClock is advanced far enough with Set or Step.
//
// Non-positive durations are treated as immediately due, matching the useful
// behavior of time.NewTimer for immediate timers. Immediate fake timers are
// delivered after the fake clock lock is released and before NewTimer returns.
//
// A timer created by FakeClock belongs to that FakeClock. It must not be moved
// between fake clocks.
func (c *FakeClock) NewTimer(d time.Duration) Timer {
	timer, delivery, ok := c.newTimer(d)
	if ok {
		delivery.deliver()
	}

	return timer
}

// fakeTimer is the FakeClock-backed implementation of Timer.
//
// fakeTimer is intentionally unexported. Callers create fake timers through
// FakeClock.NewTimer and interact with them only through the Timer interface.
//
// All mutable fakeTimer state is protected by the owning FakeClock mutex.
// Methods that mutate timer lifecycle state must acquire t.clock.mu.
//
// The timer channel is buffered with capacity one. This allows an already-fired
// timer to deliver its value without blocking fake-clock advancement when the
// receiver is scheduled later.
//
// fakeTimer follows this lifecycle model:
//
//   - active timers are registered in FakeClock.timers;
//   - due active timers are removed from FakeClock.timers and fired once;
//   - Stop removes an active timer and prevents future delivery;
//   - Reset schedules the timer again at FakeClock.Now()+d;
//   - Reset returns whether the timer was active immediately before the reset.
//
// Reset does not drain the delivery channel. A value delivered before Reset may
// still be observed by a receiver. If a previously delivered value is still
// unread, a later immediate or due delivery may be dropped because fake timer
// delivery is non-blocking. Components that require strict timer ownership must
// coordinate receives, Stop, and Reset at the component level.
type fakeTimer struct {
	clock *FakeClock

	deadline time.Time
	sequence uint64

	ch     chan time.Time
	active bool
}

// C returns the timer delivery channel.
//
// The returned channel is receive-only for callers. It is not closed after the
// timer fires and is not closed by Stop. This mirrors the standard library timer
// model and prevents channel close from becoming a lifecycle signal.
//
// C returns the same channel on every call.
func (t *fakeTimer) C() <-chan time.Time {
	return t.ch
}

// Stop prevents the fake timer from firing if it is still active.
//
// Stop returns true when the timer was active and was stopped by this call. It
// returns false when the timer had already fired, had already been stopped, or
// was otherwise inactive.
//
// Stop does not drain the delivery channel. If the timer already delivered a
// value before Stop, that value may still be observed by a receiver.
//
// Stop is safe to call more than once.
func (t *fakeTimer) Stop() bool {
	c := t.clock

	c.mu.Lock()
	defer c.mu.Unlock()

	if !t.active {
		return false
	}

	t.active = false
	delete(c.timers, t)

	return true
}

// Reset changes the fake timer to fire after d according to the owning fake
// clock.
//
// Reset schedules the timer at:
//
//	FakeClock.Now() + d
//
// Reset returns true when the timer was active immediately before the reset and
// false when it was inactive, stopped, or already fired.
//
// Reset reactivates stopped and already-fired timers. It does not drain the
// delivery channel. A value delivered before Reset may still be observed by a
// receiver, so components that reuse timers must coordinate receives according
// to their ownership model.
//
// If a previously delivered value is still unread, a later immediate or due
// delivery may be dropped because fake timer delivery is non-blocking.
// Components that reuse timers and need to observe every firing must coordinate
// receives or drain ownership before Reset.
//
// Non-positive durations are treated as immediately due. In that case Reset
// schedules an immediate delivery after releasing the fake clock lock.
func (t *fakeTimer) Reset(d time.Duration) bool {
	wasActive, delivery, ok := t.reset(d)
	if ok {
		delivery.deliver()
	}

	return wasActive
}

// reset updates timer state and returns an immediate delivery when the new
// deadline is already due.
//
// The returned delivery, when ok is true, must be delivered after c.mu has been
// released.
func (t *fakeTimer) reset(d time.Duration) (wasActive bool, delivery fakeTimerDelivery, ok bool) {
	c := t.clock

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureTimerStoreLocked()

	wasActive = t.active

	if t.active {
		delete(c.timers, t)
	}

	t.deadline = c.now.Add(d)
	t.sequence = c.nextSequenceLocked()
	t.active = true

	if !c.now.Before(t.deadline) {
		t.active = false

		return wasActive, fakeTimerDelivery{
			ch:    t.ch,
			value: c.now,
		}, true
	}

	c.timers[t] = struct{}{}

	return wasActive, fakeTimerDelivery{}, false
}

// fakeTimerDelivery describes one timer delivery that must be performed after
// the fake clock lock has been released.
//
// FakeClock advancement code must not send on user-visible channels while
// holding FakeClock.mu. Channel sends can interact with receiver scheduling and
// must stay outside the internal state lock.
type fakeTimerDelivery struct {
	ch    chan time.Time
	value time.Time
}

// deliver sends the timer value if the timer channel can accept it.
//
// If the channel already contains an unread timer value, this delivery is
// dropped. This preserves the fake clock invariant that time advancement never
// blocks on receiver scheduling. Timer owners that reuse timers must coordinate
// channel reads before relying on another delivery.
func (d fakeTimerDelivery) deliver() {
	select {
	case d.ch <- d.value:
	default:
	}
}

// newTimer creates and registers a fake timer.
//
// The returned delivery, when ok is true, must be delivered after c.mu has been
// released.
func (c *FakeClock) newTimer(d time.Duration) (*fakeTimer, fakeTimerDelivery, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureTimerStoreLocked()

	timer := &fakeTimer{
		clock:    c,
		deadline: c.now.Add(d),
		sequence: c.nextSequenceLocked(),
		ch:       make(chan time.Time, 1),
		active:   true,
	}

	if !c.now.Before(timer.deadline) {
		timer.active = false

		return timer, fakeTimerDelivery{
			ch:    timer.ch,
			value: c.now,
		}, true
	}

	c.timers[timer] = struct{}{}

	return timer, fakeTimerDelivery{}, false
}

// ensureTimerStoreLocked initializes the fake timer registry.
//
// The caller must hold c.mu.
func (c *FakeClock) ensureTimerStoreLocked() {
	if c.timers == nil {
		c.timers = make(map[*fakeTimer]struct{})
	}
}

// collectDueTimerDeliveriesLocked collects all timers that are due at the
// current fake time.
//
// The caller must hold c.mu.
//
// Due timers are removed from the timer registry and marked inactive before
// delivery is returned. Delivery itself must happen after c.mu has been released.
func (c *FakeClock) collectDueTimerDeliveriesLocked() []fakeTimerDelivery {
	if len(c.timers) == 0 {
		return nil
	}

	due := make([]*fakeTimer, 0, len(c.timers))

	for timer := range c.timers {
		if timer.isDueLocked(c.now) {
			due = append(due, timer)
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

	deliveries := make([]fakeTimerDelivery, 0, len(due))

	for _, timer := range due {
		delivery, _ := timer.collectDueDeliveryLocked(c.now)

		delete(c.timers, timer)
		deliveries = append(deliveries, delivery)
	}

	return deliveries
}

// isDueLocked reports whether the timer should fire at now.
//
// The caller must hold the owning FakeClock mutex.
func (t *fakeTimer) isDueLocked(now time.Time) bool {
	return t.active && !now.Before(t.deadline)
}

// collectDueDeliveryLocked returns one delivery when the timer is due.
//
// The caller must hold the owning FakeClock mutex.
func (t *fakeTimer) collectDueDeliveryLocked(now time.Time) (fakeTimerDelivery, bool) {
	if !t.isDueLocked(now) {
		return fakeTimerDelivery{}, false
	}

	t.active = false

	return fakeTimerDelivery{
		ch:    t.ch,
		value: now,
	}, true
}

// deliverFakeTimerDeliveries delivers all collected timer deliveries.
//
// FakeClock advancement code should call this after releasing FakeClock.mu.
func deliverFakeTimerDeliveries(deliveries []fakeTimerDelivery) {
	for _, delivery := range deliveries {
		delivery.deliver()
	}
}
