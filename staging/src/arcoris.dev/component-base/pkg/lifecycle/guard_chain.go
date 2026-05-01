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

// GuardChain evaluates several transition guards in order.
//
// GuardChain is useful when a component wants to compose independent lifecycle
// preconditions without creating a custom aggregate guard. The first guard that
// returns a non-nil error rejects the transition and stops evaluation.
//
// Nil guards in the chain are ignored. This makes conditional guard composition
// safe and keeps option-building code simple.
type GuardChain []TransitionGuard

// Allow evaluates all guards in order and returns the first rejection.
//
// If all guards allow the transition, Allow returns nil. The chain itself does
// not wrap errors. Controller-level code should wrap guard rejections with the
// transition and package-level guard error sentinel when returning them to
// callers.
func (g GuardChain) Allow(transition Transition) error {
	for _, guard := range g {
		if guard == nil {
			continue
		}

		if err := guard.Allow(transition); err != nil {
			return err
		}
	}

	return nil
}

// allowTransition evaluates guards for transition.
//
// This package-local helper is the controller-facing form of GuardChain. It
// keeps controller code simple while allowing options to store guards as a plain
// slice.
//
// The helper returns the first guard rejection unchanged. Controller is
// responsible for wrapping the returned error with guard-specific error
// semantics before exposing it to callers.
func allowTransition(guards []TransitionGuard, transition Transition) error {
	return GuardChain(guards).Allow(transition)
}
