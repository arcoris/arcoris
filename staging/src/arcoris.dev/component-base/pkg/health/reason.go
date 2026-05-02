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

package health

// maxReasonLength defines the maximum allowed length of a Reason.
const maxReasonLength = 128

// Reason identifies the stable machine-readable cause of a health result.
//
// A Reason explains why a Result has its Status. It is intentionally separate
// from Result.Message and Result.Cause:
//
//   - Reason is stable and machine-readable.
//   - Message is safe and human-readable.
//   - Cause is internal and may contain low-level error details.
//
// Reason is string-based rather than enum-based because ARCORIS domain packages
// may define their own health causes without changing the health core package.
// Custom reasons MUST use the same stable lower_snake_case format as the built-in
// reasons and MUST NOT contain secrets, credentials, connection strings, raw
// errors, stack traces, resource identifiers, timestamps, or other dynamic
// values.
//
// The zero value is ReasonNone. It is valid and means that no specific reason was
// provided. Healthy results commonly use ReasonNone. Non-healthy results SHOULD
// provide a reason when the owner can classify the condition safely.
type Reason string

const (
	// ReasonNone means that no specific reason was provided.
	//
	// ReasonNone is the zero value. It is valid and is commonly used for healthy
	// results or for observations where the status is sufficient on its own.
	ReasonNone Reason = ""

	// ReasonNotObserved means that no health observation is available yet.
	//
	// This reason is useful for unknown results produced before a checker, gate,
	// cached probe, or evaluator has observed a concrete component state.
	ReasonNotObserved Reason = "not_observed"

	// ReasonStarting means that the component or subsystem is still
	// bootstrapping.
	//
	// This reason is normally paired with StatusStarting or StatusUnknown. It
	// describes initialization progress, not terminal failure.
	ReasonStarting Reason = "starting"

	// ReasonTimeout means that a health observation did not complete before its
	// owner-controlled deadline or timeout.
	//
	// Timeout does not imply that the checked component is terminally failed. The
	// target and policy interpreting the result decide whether a timeout fails
	// startup, liveness, readiness, admission, or another higher-level decision.
	ReasonTimeout Reason = "timeout"

	// ReasonCanceled means that a health observation was canceled before it could
	// produce a reliable result.
	//
	// Cancellation usually reflects owner-controlled context cancellation,
	// shutdown, request cancellation, or evaluator interruption. It should remain
	// distinct from timeout when the owner can classify the difference.
	ReasonCanceled Reason = "canceled"

	// ReasonPanic means that a checker panicked while producing a health result.
	//
	// Evaluators SHOULD recover checker panics and convert them into a health
	// result with ReasonPanic while preserving internal details only in Cause or
	// owner-controlled diagnostics.
	ReasonPanic Reason = "panic"

	// ReasonDraining means that the component is intentionally draining existing
	// work and should normally stop receiving new work.
	//
	// Draining should normally affect readiness or admission, not liveness. A
	// draining component can still be alive and making progress.
	ReasonDraining Reason = "draining"

	// ReasonShuttingDown means that the component or process is executing a
	// shutdown sequence.
	//
	// Shutdown should normally affect readiness before it affects liveness. A
	// component may be shutting down cleanly while still making progress.
	ReasonShuttingDown Reason = "shutting_down"

	// ReasonDependencyUnavailable means that a required dependency is not
	// currently available for the checked scope.
	//
	// Dependency unavailability should normally affect readiness or a
	// dependency-specific target. It should not automatically be treated as a
	// liveness failure unless the dependency is part of the component's own
	// progress model.
	ReasonDependencyUnavailable Reason = "dependency_unavailable"

	// ReasonOverloaded means that the component is applying overload protection or
	// cannot safely accept additional work under current load.
	//
	// Overload is an expected control signal in ARCORIS. It should normally affect
	// readiness, admission, scheduling, or load routing rather than liveness.
	ReasonOverloaded Reason = "overloaded"

	// ReasonAdmissionClosed means that the component or one of its owners has
	// intentionally closed admission for new work.
	//
	// Admission closure may be caused by overload control, draining, policy,
	// maintenance, capacity protection, or a higher-level runtime decision.
	ReasonAdmissionClosed Reason = "admission_closed"

	// ReasonCapacityExhausted means that the checked scope has no safe remaining
	// capacity for the relevant workload.
	//
	// Capacity exhaustion is not necessarily a terminal failure. It may represent
	// a temporary readiness or admission condition that a scheduler or controller
	// can react to.
	ReasonCapacityExhausted Reason = "capacity_exhausted"

	// ReasonStale means that the health result is based on state that is older
	// than the owner considers acceptable.
	//
	// Staleness is useful for cached probes, control-loop heartbeats, dependency
	// snapshots, and scheduler/runtime views. The target and policy interpreting
	// the result decide whether stale state is tolerated.
	ReasonStale Reason = "stale"

	// ReasonFatal means that the checked scope reached a fatal or unrecoverable
	// condition.
	//
	// Fatal should be reserved for conditions that indicate broken progress,
	// corrupted runtime state, failed critical control loops, or another
	// owner-defined terminal health failure.
	ReasonFatal Reason = "fatal"
)

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
		ReasonStarting,
		ReasonTimeout,
		ReasonCanceled,
		ReasonPanic,
		ReasonDraining,
		ReasonShuttingDown,
		ReasonDependencyUnavailable,
		ReasonOverloaded,
		ReasonAdmissionClosed,
		ReasonCapacityExhausted,
		ReasonStale,
		ReasonFatal:
		return true
	default:
		return false
	}
}

// IsExecutionReason reports whether r describes failure or interruption while
// producing a health observation.
//
// Execution reasons describe checker/evaluator execution, not the checked
// component's own operational state.
func (r Reason) IsExecutionReason() bool {
	switch r {
	case ReasonTimeout,
		ReasonCanceled,
		ReasonPanic:
		return true
	default:
		return false
	}
}

// IsLifecycleReason reports whether r describes startup, drain, or shutdown
// state.
//
// Lifecycle reasons may be derived from lifecycle or run ownership, but they are
// still health reasons. The lifecycle package remains the owner of lifecycle
// state transitions.
func (r Reason) IsLifecycleReason() bool {
	switch r {
	case ReasonStarting,
		ReasonDraining,
		ReasonShuttingDown:
		return true
	default:
		return false
	}
}

// IsControlReason reports whether r describes runtime control, capacity, or
// overload state.
//
// Control reasons are common in readiness and admission reports. They should not
// be treated as liveness failures by default.
func (r Reason) IsControlReason() bool {
	switch r {
	case ReasonOverloaded,
		ReasonAdmissionClosed,
		ReasonCapacityExhausted,
		ReasonStale:
		return true
	default:
		return false
	}
}
