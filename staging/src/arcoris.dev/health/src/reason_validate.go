// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package health

// String returns the diagnostic representation of r.
//
// String is intended for diagnostics, tests, logs, reports, and human-facing
// messages. It is not a serialization contract. ReasonNone returns "none" so
// diagnostics do not render an empty string. Invalid reasons return "invalid".
func (r Reason) String() string {
	if r == ReasonNone {
		return "none"
	}
	if !r.IsValid() {
		return "invalid"
	}

	return string(r)
}

// IsNone reports whether r contains no specific reason.
//
// ReasonNone is valid. It is appropriate when a result does not need a
// machine-readable cause, most commonly for healthy observations.
func (r Reason) IsNone() bool {
	return r == ReasonNone
}

// IsValid reports whether r is empty or follows the canonical reason syntax.
//
// Valid non-empty reasons use lower_snake_case with ASCII lower-case letters,
// digits, and single underscores between name parts. They MUST start with a
// lower-case letter, MUST NOT end with an underscore, MUST NOT contain repeated
// underscores, and MUST NOT exceed the package-defined maximum reason length.
//
// The syntax is intentionally restrictive so reasons remain safe for diagnostics,
// metrics labels, logs, reports, tests, and transport adapters. Dynamic details
// belong in Result.Message only when safe, or in Result.Cause when internal.
func (r Reason) IsValid() bool {
	if r == ReasonNone {
		return true
	}

	return validLowerSnakeIdentifier(string(r), maxReasonLength)
}

// IsBuiltin reports whether r is one of the reasons defined by this package.
//
// Custom domain reasons may still be valid even when IsBuiltin returns false.
// Use IsValid to validate the reason syntax, and use IsBuiltin only when code
// specifically needs to distinguish core health reasons from domain-defined
// reasons.
func (r Reason) IsBuiltin() bool {
	switch r {
	case ReasonNone,
		ReasonNotObserved,
		ReasonTimeout,
		ReasonCanceled,
		ReasonPanic,
		ReasonStarting,
		ReasonDraining,
		ReasonShuttingDown,
		ReasonDependencyUnavailable,
		ReasonDependencyDegraded,
		ReasonOverloaded,
		ReasonBackpressured,
		ReasonRateLimited,
		ReasonAdmissionClosed,
		ReasonCapacityExhausted,
		ReasonResourceExhausted,
		ReasonStale,
		ReasonNotSynced,
		ReasonSyncFailed,
		ReasonLagging,
		ReasonPartitioned,
		ReasonMisconfigured,
		ReasonFatal:
		return true
	default:
		return false
	}
}
