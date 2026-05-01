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

// WithGuard appends one transition guard to the Controller configuration.
//
// Guards are evaluated before a table-valid candidate transition is committed.
// Returning a non-nil error rejects the transition and leaves lifecycle state
// unchanged.
//
// Guards are evaluated in configuration order. If several WithGuard and
// WithGuards options are provided, their guards are appended in the same order in
// which the options are applied.
//
// A nil guard is ignored. This keeps conditional option construction safe.
func WithGuard(guard TransitionGuard) Option {
	return func(config *controllerConfig) {
		if guard == nil {
			return
		}

		config.guards = append(config.guards, guard)
	}
}

// WithGuards appends several transition guards to the Controller configuration.
//
// Guards are evaluated in the order provided. The first guard that returns a
// non-nil error rejects the transition and stops guard evaluation.
//
// Nil guards are ignored. This makes it safe to build guard lists from optional
// dependencies without adding special-case filtering at call sites.
func WithGuards(guards ...TransitionGuard) Option {
	return func(config *controllerConfig) {
		for _, guard := range guards {
			if guard == nil {
				continue
			}

			config.guards = append(config.guards, guard)
		}
	}
}
