// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clock

// Pending is a non-mutating snapshot of fake-clock delivery registrations.
//
// It is intended for deterministic test coordination. Waiters are registered by
// After and Sleep, timers by NewTimer, and tickers by NewTicker. Due entries
// remain counted until Set, Step, Stop, Reset, or a direct immediate-delivery
// path removes them.
//
// Pending reports internal fake-clock registrations, not user-visible channel
// contents. A waiter or timer can be absent from Pending even when its delivered
// value is still unread. A ticker remains pending while it is active, even if a
// previously delivered tick is still buffered in its channel.
//
// Pending is not a distributed coordination primitive and must not be used as a
// scheduler, lease, or runtime lifecycle protocol.
type Pending struct {
	// Waiters is the number of pending After and Sleep waiters.
	Waiters int

	// Timers is the number of active fake timers.
	Timers int

	// Tickers is the number of active fake tickers.
	Tickers int
}

// Pending reports the fake-clock registrations that currently remain pending.
//
// Pending is safe for concurrent use. It does not advance fake time, remove due
// entries, deliver channel values, or observe real runtime time. Tests can use
// it to coordinate with goroutines that register fake-time waits without
// introducing real sleeps.
func (c *FakeClock) Pending() Pending {
	c.mu.Lock()
	defer c.mu.Unlock()

	return Pending{
		Waiters: len(c.waiters),
		Timers:  len(c.timers),
		Tickers: len(c.tickers),
	}
}
