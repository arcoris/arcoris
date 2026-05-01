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

// Package lifecycle provides a small lifecycle state owner for ARCORIS
// component internals.
//
// The package records component lifecycle phase, validates legal movements,
// exposes immutable snapshots, and lets callers wait for read-side lifecycle
// conditions. It belongs to component-base because many components need the
// same mechanical lifecycle vocabulary without depending on a runtime,
// supervisor, queue, scheduler, retry, health, logging, metrics, or tracing
// package.
//
// # Lifecycle model
//
// The default lifecycle graph is:
//
//	New -> Starting -> Running -> Stopping -> Stopped
//	                  \---------------------> Failed
//
// StateStarting and StateStopping are transitional states. They record that the
// component owner has begun startup or shutdown work, but the package does not
// execute that work. StateStopped and StateFailed are terminal states. Once a
// Controller reaches a terminal state, the same lifecycle instance cannot be
// started again.
//
// # Controller, Snapshot, and Transition
//
// Controller owns mutable lifecycle state. It serializes transition commits with
// an internal mutex, enforces the transition table, runs guards before commit,
// publishes committed transition metadata, closes wait signals, and notifies
// observers after commit.
//
// Snapshot is the immutable read model returned by Controller. It is copyable and
// contains no locks, observers, guards, waiters, or references back to the
// controller.
//
// Transition describes candidate or committed lifecycle movement. Candidate
// transitions are produced before commit and may not have Revision or At set.
// Committed transitions carry controller-assigned revision and time metadata.
//
// # Guards and observers
//
// TransitionGuard is pre-commit validation. Guards run after table validation
// and required failure-cause validation, while the controller lock is held. A
// guard rejection leaves state unchanged, does not advance revision, does not
// signal waiters, and does not notify observers.
//
// Observer is post-commit notification. Observers receive committed transitions
// after the controller lock has been released. Observer failures cannot roll back
// a committed transition, so observer methods do not return errors.
//
// # Wait API
//
// Wait, WaitState, WaitTerminal, and Done are read-side blocking APIs. They
// observe snapshots and controller-owned signal channels. Wait predicates must
// be fast and must not call transition methods on the same Controller.
//
// Wait accepts a nil context as context.Background. A nil predicate is a
// lifecycle wait error, not a panic. WaitState rejects invalid target states and
// uses static graph reachability to fail early when a requested state can no
// longer be observed.
//
// # Zero values and concurrency
//
// The zero Controller is usable and starts in StateNew with no guards, no
// observers, and time.Now as its commit time source. A zero Snapshot is a valid
// initial read model. The zero Transition is not valid lifecycle movement.
//
// Controller methods are safe to call concurrently unless a method comment says
// otherwise. Controller serializes transition commits. Snapshots and transitions
// are values and may be copied. Guard and observer implementations are owned by
// the caller and must be safe for the way they are configured.
//
// # Non-goals
//
// This package does not provide:
//
//   - goroutine supervision or worker lifecycle management;
//   - retry, backoff, restart, or recovery policy;
//   - health, readiness, admission, or workload policy;
//   - logging, metrics, tracing, or exporter backends;
//   - scheduler, worker, queue, or rate-limit semantics;
//   - a generic StateMachine[S, E] framework;
//   - serialization or wire API compatibility contracts.
//
// Higher-level packages may build those behaviors around Controller, but this
// package must remain a small mechanical lifecycle layer.
//
// # Dependency policy
//
// Production code in this package depends only on the Go standard library and
// the minimal time-source interface accepted by WithClock.
package lifecycle
