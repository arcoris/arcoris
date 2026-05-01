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

// State returns the current lifecycle state.
//
// State is safe to call concurrently with transition methods. It observes only
// the current State value. Call Snapshot when the caller also needs Revision,
// LastTransition, or FailureCause.
func (c *Controller) State() State {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.state
}

// Snapshot returns a consistent point-in-time view of the controller.
//
// Snapshot is safe to call concurrently with transition methods. The returned
// value is copyable and does not reference controller internals. Snapshot does
// not initialize wait channels for a zero-value Controller because it is a pure
// read model; Done and Wait initialize signals when they need them.
func (c *Controller) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.snapshotLocked()
}

// snapshotLocked returns the current controller snapshot.
//
// The caller must hold c.mu.
func (c *Controller) snapshotLocked() Snapshot {
	return Snapshot{
		State:          c.state,
		Revision:       c.revision,
		LastTransition: c.lastTransition,
		FailureCause:   c.failureCause,
	}
}
