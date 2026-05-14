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

// BeginStart records that the component owner has begun startup.
func (c *Controller) BeginStart() (Transition, error) {
	return c.apply(EventBeginStart, nil)
}

// MarkRunning records that startup completed successfully.
func (c *Controller) MarkRunning() (Transition, error) {
	return c.apply(EventMarkRunning, nil)
}

// BeginStop records that the component owner has begun shutdown.
func (c *Controller) BeginStop() (Transition, error) {
	return c.apply(EventBeginStop, nil)
}

// MarkStopped records that shutdown completed successfully.
func (c *Controller) MarkStopped() (Transition, error) {
	return c.apply(EventMarkStopped, nil)
}

// MarkFailed records that the lifecycle reached an unsuccessful terminal state.
//
// The cause argument must be non-nil when EventMarkFailed is table-valid for the
// current state. The cause is stored in the committed Transition and subsequent
// Snapshots.
func (c *Controller) MarkFailed(cause error) (Transition, error) {
	return c.apply(EventMarkFailed, cause)
}

// apply validates and commits one lifecycle event.
//
// On invalid transition, missing failure cause, or guard rejection, apply leaves
// state unchanged, does not increment revision, does not signal waiters, and does
// not notify observers. Guards run while c.mu is held so state cannot change
// between guard validation and commit. Observers run only after commit and
// outside the lock.
func (c *Controller) apply(event Event, cause error) (Transition, error) {
	var (
		committed Transition
		observers []Observer
	)

	c.mu.Lock()
	c.ensureInitializedLocked()

	current := c.state

	transition, ok := reduceTransition(current, event, cause)
	if !ok {
		err := ErrInvalidTransition
		if current.IsTerminal() {
			err = ErrTerminalState
		}

		c.mu.Unlock()
		return Transition{}, newTransitionError(current, event, err)
	}

	if transition.Event.RequiresCause() && transition.Cause == nil {
		c.mu.Unlock()
		return Transition{}, newTransitionError(current, event, ErrFailureCauseRequired)
	}

	if err := allowTransition(c.guards, transition); err != nil {
		c.mu.Unlock()
		return Transition{}, newGuardError(transition, err)
	}

	committed = c.commitLocked(transition)
	observers = c.observers

	c.mu.Unlock()

	notifyObservers(observers, committed)

	return committed, nil
}

// commitLocked commits transition to controller state.
//
// The caller must hold c.mu. The transition must already be table-valid, carry
// required runtime payload, and pass all configured guards.
func (c *Controller) commitLocked(transition Transition) Transition {
	c.revision++

	transition = transition.withCommitMetadata(c.revision, c.commitTimeLocked())

	c.state = transition.To
	c.lastTransition = transition

	if transition.To == StateFailed {
		c.failureCause = transition.Cause
	}

	c.signalChangeLocked(transition.To.IsTerminal())

	return transition
}
