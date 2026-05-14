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

// TransitionGuard decides whether a candidate lifecycle transition may be
// committed.
//
// A guard is evaluated after the transition has been accepted by the lifecycle
// transition table and after required transition payload has been checked, but
// before Controller commits the new state. Returning nil allows the transition
// to continue. Returning a non-nil error rejects the transition and MUST leave
// the lifecycle state unchanged.
//
// A guard is a precondition check, not an effect handler. It SHOULD be fast,
// deterministic for the state it observes, and side-effect free. It MUST NOT run
// lifecycle work such as starting components, stopping goroutines, closing
// resources, retrying operations, waiting on backoff, publishing observers, or
// calling transition methods on the same Controller.
//
// Guards are intended for owner-specific invariants that do not belong in the
// static lifecycle graph. Examples include startup barriers, dependency
// readiness, shutdown barriers, in-flight work checks, and component-local
// configuration validity.
//
// Controller should wrap guard rejections in the package guard error type so
// callers can distinguish guard rejection from invalid transition, missing
// failure cause, context cancellation, or other controller errors.
type TransitionGuard interface {
	// Allow returns nil when transition may be committed.
	//
	// The transition passed to Allow is a candidate transition. It may not have
	// commit metadata such as Revision or At because those fields are assigned
	// only after guards allow the transition.
	Allow(transition Transition) error
}
