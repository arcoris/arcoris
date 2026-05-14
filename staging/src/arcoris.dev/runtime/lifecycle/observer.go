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

// Observer observes committed lifecycle transitions.
//
// Observers are notified after Controller has committed a transition and after
// the committed Transition has been assigned controller metadata such as Revision
// and At. An observer receives facts that already happened; it cannot reject,
// modify, or roll back the transition.
//
// Observer is intended for diagnostics and integration with external concerns
// such as metrics, tracing, debug timelines, tests, audit-like internal event
// streams, and health or supervisor bridges. The lifecycle package itself does
// not depend on logging, metrics, tracing, or health packages.
//
// Observers MUST be safe to call after commit and MUST NOT call transition
// methods on the same Controller. Observers SHOULD be fast and non-blocking. If
// an observer needs to perform expensive work, it should hand the transition to
// a component-owned queue or runtime outside the lifecycle critical path.
//
// Controller calls observers outside its internal transition lock. This prevents
// observer code from blocking state reads, waiters, or future transition
// attempts while still guaranteeing that observers see only committed
// transitions.
//
// Observer methods do not return errors because observer failure cannot undo a
// transition that has already been committed. Observers that interact with
// fallible external systems must handle failures internally.
//
// Observers MUST NOT panic. The lifecycle package does not define a panic
// recovery policy for observers; callers that need panic isolation should wrap
// observers explicitly.
type Observer interface {
	// ObserveLifecycleTransition observes a committed lifecycle transition.
	//
	// The transition passed to this method has already been committed. It SHOULD
	// contain non-zero commit metadata such as Revision and At when emitted by
	// Controller.
	ObserveLifecycleTransition(transition Transition)
}
