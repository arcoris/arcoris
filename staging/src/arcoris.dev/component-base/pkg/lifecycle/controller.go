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

import (
	"sync"
	"time"
)

// Controller owns and serializes the lifecycle state of one component instance.
//
// Controller is a state owner, not a component runner. It records lifecycle
// transitions, enforces the lifecycle transition table, runs transition guards,
// publishes committed transition metadata, exposes consistent snapshots, and
// notifies observers after successful commits.
//
// Controller does not start goroutines, stop workers, close resources, retry
// operations, compute health, export metrics, or write logs. Component owners
// perform those actions around lifecycle transitions:
//
//   - call BeginStart before executing startup work;
//   - call MarkRunning after startup has completed successfully;
//   - call BeginStop before executing shutdown work;
//   - call MarkStopped after shutdown has completed successfully;
//   - call MarkFailed when startup, runtime, or shutdown fails terminally.
//
// Controller serializes transition commits with an internal mutex. Guards are
// evaluated while the transition lock is held so the state cannot change between
// validation and commit. Observers are notified only after commit and outside the
// lock.
//
// The zero Controller is usable and starts in StateNew with no guards, no
// observers, and time.Now as the commit time source. A Controller must not be
// copied after first use.
type Controller struct {
	mu sync.Mutex

	state          State
	revision       uint64
	lastTransition Transition
	failureCause   error

	now       func() time.Time
	guards    []TransitionGuard
	observers []Observer

	changed chan struct{}
	done    chan struct{}
}

// NewController constructs a lifecycle Controller.
//
// The returned controller starts in StateNew with Revision zero and no committed
// LastTransition. Options configure construction-time dependencies such as the
// time source, transition guards, and observers.
//
// Options are applied once during construction. The resulting guard and observer
// slices are owned by the controller and are not expected to change after
// construction.
func NewController(options ...Option) *Controller {
	config := newControllerConfig(options...)

	if config.now == nil {
		config.now = time.Now
	}

	return &Controller{
		state:     StateNew,
		now:       config.now,
		guards:    append([]TransitionGuard(nil), config.guards...),
		observers: append([]Observer(nil), config.observers...),
		changed:   make(chan struct{}),
		done:      make(chan struct{}),
	}
}

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
// value is copyable and does not reference controller internals.
func (c *Controller) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.snapshotLocked()
}

// BeginStart records that the component owner has begun startup.
//
// BeginStart normally commits:
//
//	StateNew --EventBeginStart--> StateStarting
//
// BeginStart does not execute startup work. The owner should call BeginStart,
// perform startup work, and then call MarkRunning or MarkFailed depending on the
// outcome.
func (c *Controller) BeginStart() (Transition, error) {
	return c.apply(EventBeginStart, nil)
}

// MarkRunning records that startup completed successfully.
//
// MarkRunning normally commits:
//
//	StateStarting --EventMarkRunning--> StateRunning
//
// MarkRunning does not make dependencies ready, start workers, or publish
// external readiness by itself. It only records that the lifecycle has entered
// the normal running phase.
func (c *Controller) MarkRunning() (Transition, error) {
	return c.apply(EventMarkRunning, nil)
}

// BeginStop records that the component owner has begun shutdown.
//
// BeginStop normally commits one of:
//
//	StateNew      --EventBeginStop--> StateStopped
//	StateStarting --EventBeginStop--> StateStopping
//	StateRunning  --EventBeginStop--> StateStopping
//
// BeginStop does not execute shutdown work. The owner should call BeginStop,
// perform shutdown work, and then call MarkStopped or MarkFailed depending on
// the outcome.
func (c *Controller) BeginStop() (Transition, error) {
	return c.apply(EventBeginStop, nil)
}

// MarkStopped records that shutdown completed successfully.
//
// MarkStopped normally commits:
//
//	StateStopping --EventMarkStopped--> StateStopped
//
// StateStopped is terminal. Once this transition is committed, the same
// Controller instance cannot be started again.
func (c *Controller) MarkStopped() (Transition, error) {
	return c.apply(EventMarkStopped, nil)
}

// MarkFailed records that the lifecycle reached an unsuccessful terminal state.
//
// MarkFailed normally commits one of:
//
//	StateStarting --EventMarkFailed--> StateFailed
//	StateRunning  --EventMarkFailed--> StateFailed
//	StateStopping --EventMarkFailed--> StateFailed
//
// The cause argument MUST be non-nil when the transition is table-valid. The
// cause is stored in the committed Transition and in subsequent Snapshots so
// waiters, diagnostics, health mapping, and supervisors can inspect why the
// lifecycle failed.
func (c *Controller) MarkFailed(cause error) (Transition, error) {
	return c.apply(EventMarkFailed, cause)
}

// apply validates and commits one lifecycle event.
//
// apply is the common transition path for the public lifecycle methods. It keeps
// the transition algorithm in one place:
//
//   - initialize zero-value controller internals if needed;
//   - read the current state;
//   - reduce state/event/cause to a candidate transition;
//   - reject table-invalid transitions;
//   - reject missing failure causes;
//   - run guards before commit;
//   - assign revision and commit time;
//   - publish the new state;
//   - notify observers after releasing the lock.
//
// apply returns the committed Transition on success. On failure, the controller
// state is unchanged and observers are not notified.
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

	// Observers are immutable after construction. Keeping the slice value is
	// enough; observers are invoked outside the lock.
	observers = c.observers

	c.mu.Unlock()

	notifyObservers(observers, committed)

	return committed, nil
}

// commitLocked commits transition to controller state.
//
// The caller MUST hold c.mu. The transition MUST already be table-valid, carry
// required runtime payload, and pass all configured guards.
//
// commitLocked assigns monotonic commit metadata, updates the current snapshot
// fields, wakes future waiters through the change signal, and closes the done
// signal if the transition reaches a terminal state.
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

// snapshotLocked returns the current controller snapshot.
//
// The caller MUST hold c.mu.
func (c *Controller) snapshotLocked() Snapshot {
	return Snapshot{
		State:          c.state,
		Revision:       c.revision,
		LastTransition: c.lastTransition,
		FailureCause:   c.failureCause,
	}
}

// ensureInitializedLocked initializes lazy zero-value Controller internals.
//
// The caller MUST hold c.mu.
//
// NewController eagerly initializes these fields. This helper exists so a
// zero-value Controller remains usable according to the type contract.
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
// The caller MUST hold c.mu.
//
// A custom time source should return a non-zero time. If it returns the zero
// value, Controller falls back to time.Now to preserve the invariant that
// committed transitions have non-zero commit time.
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
// The caller MUST hold c.mu.
//
// changed is closed on every committed transition. For non-terminal
// transitions, Controller creates a fresh change channel for the next revision.
// For terminal transitions, Controller closes both changed and done and does not
// create another change channel because no further transitions are possible.
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
//
// The returned changed channel is closed when any later transition commits. The
// returned done channel is closed when the lifecycle reaches a terminal state.
//
// wait.go uses this helper to observe a consistent snapshot and subscribe to the
// next state change without exposing controller internals publicly.
func (c *Controller) waitSnapshot() (Snapshot, <-chan struct{}, <-chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureInitializedLocked()

	return c.snapshotLocked(), c.changed, c.done
}

// doneSignal returns the terminal lifecycle signal for wait.go.
//
// The returned channel is closed when the lifecycle reaches StateStopped or
// StateFailed. The channel is never replaced for a controller instance.
func (c *Controller) doneSignal() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureInitializedLocked()

	return c.done
}
