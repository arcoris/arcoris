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

package healthhttp

// DetailLevel controls how much check-level information health HTTP renderers
// may include in a response body.
//
// The zero value is DetailNone. Renderers must never use detail levels to expose
// health.Result.Cause, panic stacks, raw errors, credentials, or private
// diagnostic data.
type DetailLevel uint8

const (
	// DetailNone suppresses check-level details.
	DetailNone DetailLevel = iota

	// DetailFailed includes only checks that fail the handler's target policy.
	DetailFailed

	// DetailAll includes all check results in report order.
	DetailAll
)

// String returns the stable diagnostic name for level.
func (level DetailLevel) String() string {
	switch level {
	case DetailNone:
		return "none"
	case DetailFailed:
		return "failed"
	case DetailAll:
		return "all"
	default:
		return "invalid"
	}
}

// IsValid reports whether level is a supported response detail level.
func (level DetailLevel) IsValid() bool {
	switch level {
	case DetailNone, DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// IncludesChecks reports whether level allows any check-level results.
func (level DetailLevel) IncludesChecks() bool {
	switch level {
	case DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// IncludesAllChecks reports whether level allows all check results.
func (level DetailLevel) IncludesAllChecks() bool {
	return level == DetailAll
}

// IncludesFailedChecks reports whether level allows policy-failed check results.
func (level DetailLevel) IncludesFailedChecks() bool {
	switch level {
	case DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// validateDetailLevel returns an error if level is not supported.
func validateDetailLevel(level DetailLevel) error {
	if !level.IsValid() {
		return InvalidDetailLevelError{Level: level}
	}

	return nil
}
