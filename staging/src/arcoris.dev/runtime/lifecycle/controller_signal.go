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

package lifecycle

import "time"

// ensureInitializedLocked initializes lazy zero-value Controller internals.
//
// The caller must hold c.mu. NewController eagerly initializes these fields; the
// helper exists so a zero-value Controller remains usable.
func (c *Controller) ensureInitializedLocked() {
	if c.now == nil {
		c.now = time.Now
	}

	if c.changed == nil {
		c.changed = make(chan struct{})
		if c.state.IsTerminal() {
			close(c.changed)
		}
	}

	if c.done == nil {
		c.done = make(chan struct{})
		if c.state.IsTerminal() {
			close(c.done)
		}
	}
}

// commitTimeLocked returns the timestamp for a committed transition.
//
// The caller must hold c.mu. If a custom source returns zero time, Controller
// falls back to time.Now so committed transitions keep non-zero commit metadata.
func (c *Controller) commitTimeLocked() time.Time {
	now := c.now
	if now == nil {
		now = time.Now
	}

	at := now()
	if at.IsZero() {
		return time.Now()
	}

	return at
}

// signalChangeLocked wakes waiters after a committed transition.
//
// The caller must hold c.mu. changed is closed on every committed transition.
// For non-terminal transitions, Controller creates a fresh change channel for
// the next revision. For terminal transitions, Controller closes changed and
// done exactly once and does not create another change channel.
func (c *Controller) signalChangeLocked(terminal bool) {
	changed := c.changed

	if terminal {
		close(changed)
		close(c.done)
		return
	}

	c.changed = make(chan struct{})
	close(changed)
}

// waitSnapshot returns a snapshot and notification channels for wait.go.
func (c *Controller) waitSnapshot() (Snapshot, <-chan struct{}, <-chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureInitializedLocked()

	return c.snapshotLocked(), c.changed, c.done
}

// doneSignal returns the terminal lifecycle signal for Done.
func (c *Controller) doneSignal() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureInitializedLocked()

	return c.done
}
