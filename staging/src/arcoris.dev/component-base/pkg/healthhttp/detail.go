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
// DetailLevel controls only response detail selection. It does not affect health
// evaluation, target policy, HTTP status mapping, endpoint paths, request
// methods, response format, logging, metrics, tracing, authentication, or
// authorization.
//
// The zero value is DetailNone. This is intentional: health endpoints are often
// exposed to load balancers, orchestrators, and infrastructure probes where the
// HTTP status code is the primary signal and the response body should reveal as
// little as possible.
//
// DetailLevel never permits exposing Result.Cause, panic stacks, raw errors,
// connection strings, credentials, internal addresses, tenant identifiers, or
// other private diagnostic data. Renderers must use DetailLevel only to decide
// whether safe check names, statuses, reasons, and messages are included.
type DetailLevel uint8

const (
	// DetailNone suppresses check-level details.
	//
	// Renderers should emit only a compact target-level response such as "ok" or
	// "unhealthy". This is the safest default for public or semi-public probe
	// endpoints.
	DetailNone DetailLevel = iota

	// DetailFailed includes only checks that fail the handler's target policy.
	//
	// The exact failed-check selection depends on the target policy used by the
	// handler. For example, a degraded check may be included for readiness when
	// readiness does not allow degraded status, while the same degraded check may
	// be omitted for liveness if liveness allows degraded status.
	DetailFailed

	// DetailAll includes all check results in report order.
	//
	// DetailAll is intended for internal diagnostics and protected operational
	// endpoints. Even with DetailAll, renderers must not expose Result.Cause or
	// other private diagnostic data.
	DetailAll
)

// String returns the stable diagnostic name for level.
//
// String is intended for diagnostics, tests, and error messages. It is not a
// wire-format negotiation mechanism. Invalid values return "invalid".
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

// IncludesChecks reports whether level allows any check-level results to be
// included in a response.
//
// DetailFailed and DetailAll both include check-level output. DetailNone does
// not. This helper is useful for renderers that can skip check selection
// entirely when no check details are allowed.
func (level DetailLevel) IncludesChecks() bool {
	switch level {
	case DetailFailed, DetailAll:
		return true
	default:
		return false
	}
}

// IncludesAllChecks reports whether level allows all check results to be
// included in a response.
//
// DetailAll is the only level that includes successful, failed, degraded,
// unknown, and starting checks together. DetailFailed includes only policy-failed
// checks and DetailNone includes no checks.
func (level DetailLevel) IncludesAllChecks() bool {
	return level == DetailAll
}

// IncludesFailedChecks reports whether level allows policy-failed check results
// to be included in a response.
//
// DetailFailed and DetailAll include failed checks. DetailNone suppresses them.
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
// The helper is package-private because callers can use DetailLevel.IsValid for
// boolean checks, while option parsing and constructor code need an
// error-returning boundary.
func validateDetailLevel(level DetailLevel) error {
	if !level.IsValid() {
		return InvalidDetailLevelError{Level: level}
	}

	return nil
}
