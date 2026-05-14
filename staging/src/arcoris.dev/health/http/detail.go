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
//
// DetailLevel affects exposure only. It does not change how package health
// evaluates checks or aggregates report status.
type DetailLevel uint8

const (
	// DetailNone suppresses check-level details.
	//
	// The response still reports the aggregate target status, but no per-check
	// breakdown is exposed.
	DetailNone DetailLevel = iota

	// DetailFailed includes only checks that fail the handler's target policy.
	//
	// This is often a good compromise when operators need failure visibility
	// without routinely exposing passing checks.
	DetailFailed

	// DetailAll includes all check results in report order.
	//
	// Even at this most verbose level, the adapter still emits only safe DTO
	// fields and never exposes Cause or raw panic details.
	DetailAll
)

// String returns the stable diagnostic name for level.
//
// Invalid values return "invalid" so misconfiguration is explicit in tests and
// diagnostics.
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
//
// Unsupported values are rejected during handler construction.
func (level DetailLevel) IsValid() bool {
	switch level {
	case DetailNone, DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// IncludesChecks reports whether level allows any check-level results.
//
// This helper exists to keep rendering code explicit about when it may include
// per-check output at all.
func (level DetailLevel) IncludesChecks() bool {
	switch level {
	case DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// IncludesAllChecks reports whether level allows all check results.
//
// The helper is intentionally narrow so render paths can express their exposure
// choices without repeating enum comparisons.
func (level DetailLevel) IncludesAllChecks() bool {
	return level == DetailAll
}

// IncludesFailedChecks reports whether level allows policy-failed check results.
//
// The notion of "failed" remains policy-relative: the same check can be failed
// for readiness and passing for liveness.
func (level DetailLevel) IncludesFailedChecks() bool {
	switch level {
	case DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// validateDetailLevel returns an error if level is not supported.
//
// Validation stays local to preserve the adapter's own typed error surface.
func validateDetailLevel(level DetailLevel) error {
	if !level.IsValid() {
		return InvalidDetailLevelError{Level: level}
	}

	return nil
}
