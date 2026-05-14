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

package run

// GroupOption configures a Group during construction.
//
// Options are applied to an internal groupConfig before the Group starts any
// tasks. They do not mutate an already constructed Group. A nil option is a
// programming error and causes NewGroup to panic so invalid conditional option
// composition is visible at the construction boundary.
type GroupOption func(*groupConfig)

// WithCancelOnError configures whether task errors cancel the group context.
//
// When enabled, the first task that returns a non-nil error cancels the group
// context with that task's TaskError. Sibling tasks can observe the cancellation
// and stop. When disabled, task errors are recorded but do not cancel the group
// context.
func WithCancelOnError(enabled bool) GroupOption {
	return func(cfg *groupConfig) {
		cfg.cancelOnError = enabled
	}
}

// WithErrorMode configures how Wait reports recorded task errors.
//
// ErrorModeJoin keeps all task failures, orders joined TaskError values by
// submission sequence, and uses errors.Join even when only one task failed.
// ErrorModeFirst keeps only the first error recorded by the Group under its
// mutex. WithErrorMode panics when mode is unknown.
func WithErrorMode(mode ErrorMode) GroupOption {
	requireErrorMode(mode)

	return func(cfg *groupConfig) {
		cfg.errorMode = mode
	}
}
