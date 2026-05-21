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

// TransitionGuardFunc adapts a function to TransitionGuard.
type TransitionGuardFunc func(Transition) error

// Allow calls f with transition.
//
// A nil TransitionGuardFunc allows every transition. This behavior makes
// optional guard wiring safe when a guard function is conditionally configured.
// Production code should still avoid intentionally passing nil guard functions
// because doing so usually hides a missing precondition.
func (f TransitionGuardFunc) Allow(transition Transition) error {
	if f == nil {
		return nil
	}

	return f(transition)
}
