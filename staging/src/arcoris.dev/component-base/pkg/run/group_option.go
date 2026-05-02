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

// ErrorMode controls how Group.Wait reports recorded task errors.
type ErrorMode uint8

const (
	// ErrorModeJoin reports all task errors using errors.Join.
	//
	// Joined task errors are ordered by task submission order, not goroutine
	// completion order. Join mode intentionally uses errors.Join even for one
	// task error so callers can handle joined and single-failure results through
	// the same Unwrap() []error and TaskErrors paths. This is the default because
	// stable diagnostics are more useful for component runtimes than
	// race-dependent completion ordering.
	ErrorModeJoin ErrorMode = iota

	// ErrorModeFirst reports only the first observed task error.
	//
	// "First" means the first task error recorded by the Group under its mutex.
	// It is not sorted by task submission sequence and is not a wall-clock
	// ordering guarantee. This mode is useful for fail-fast component scopes
	// where sibling tasks are expected to return context cancellation after the
	// first task fails.
	ErrorModeFirst
)

// IsValid reports whether mode is a known ErrorMode value.
func (m ErrorMode) IsValid() bool {
	return m == ErrorModeJoin || m == ErrorModeFirst
}

// groupConfig contains construction-time settings for Group.
type groupConfig struct {
	cancelOnError bool
	errorMode     ErrorMode
}

// defaultGroupConfig returns the default Group configuration.
func defaultGroupConfig() groupConfig {
	return groupConfig{
		cancelOnError: true,
		errorMode:     ErrorModeJoin,
	}
}

// GroupOption configures a Group during construction.
//
// Options are applied to an internal groupConfig before the Group starts any
// tasks. They do not mutate an already constructed Group.
type GroupOption func(*groupConfig)

// newGroupConfig applies options to a fresh default groupConfig.
//
// Nil options are ignored to keep conditional option lists easy to compose.
func newGroupConfig(options ...GroupOption) groupConfig {
	config := defaultGroupConfig()

	for _, option := range options {
		if option == nil {
			continue
		}
		option(&config)
	}

	requireErrorMode(config.errorMode)

	return config
}

// WithCancelOnError configures whether task errors cancel the group context.
//
// When enabled, the first task that returns a non-nil error cancels the group
// context with that task's TaskError. Sibling tasks can observe the cancellation
// and stop. When disabled, task errors are recorded but do not cancel the group
// context.
func WithCancelOnError(enabled bool) GroupOption {
	return func(config *groupConfig) {
		config.cancelOnError = enabled
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

	return func(config *groupConfig) {
		config.errorMode = mode
	}
}
