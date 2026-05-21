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
