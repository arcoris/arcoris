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
// errors, stack traces, resource identifiers, timestamps, addresses, tenant IDs,
// object IDs, or other dynamic values.
//
// The zero value is ReasonNone. It is valid and means that no specific reason was
// provided. Healthy results commonly use ReasonNone. Non-healthy results SHOULD
// provide a reason when the owner can classify the condition safely.
//
// Reason does not define severity. The same reason may be paired with different
// statuses depending on the owner and target. For example, overload may be
// degraded while admission is still partially available, or unhealthy when the
// component must stop accepting new work.
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

	// ReasonStarting means that the component or subsystem is still
	// bootstrapping.
	//
	// This reason is normally paired with StatusStarting or StatusUnknown. It
	// describes initialization progress, not terminal failure.
	ReasonStarting Reason = "starting"

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

	// ReasonDependencyDegraded means that a required dependency is reachable but
	// operating below the owner's healthy threshold.
	//
	// Dependency degradation is useful when a dependency is still usable with
	// reduced confidence, higher latency, partial errors, fallback behavior, or
	// reduced capacity. It should remain distinct from dependency unavailability.
	ReasonDependencyDegraded Reason = "dependency_degraded"

	// ReasonOverloaded means that the component is under load pressure above the
	// owner-defined safe operating threshold.
	//
	// Overload is an expected control signal in ARCORIS. It should normally affect
	// readiness, admission, scheduling, or load routing rather than liveness.
	ReasonOverloaded Reason = "overloaded"

	// ReasonBackpressured means that the component is actively applying
	// backpressure to protect itself or upstream/downstream runtime boundaries.
	//
	// Backpressure describes a protective control state. It may be caused by
	// overload, limited downstream capacity, queue pressure, or owner-defined
	// flow-control policy.
	ReasonBackpressured Reason = "backpressured"

	// ReasonRateLimited means that work is being limited by a rate policy.
	//
	// Rate limiting may be protective, tenant-specific, workload-specific, or
	// operator-configured. It is distinct from overload because the component may
	// be enforcing policy even when it still has physical capacity.
	ReasonRateLimited Reason = "rate_limited"

	// ReasonAdmissionClosed means that the component or one of its owners has
	// intentionally closed admission for new work.
	//
	// Admission closure may be caused by overload control, draining, policy,
	// maintenance, capacity protection, or a higher-level runtime decision.
	ReasonAdmissionClosed Reason = "admission_closed"

	// ReasonCapacityExhausted means that the checked scope has no safe remaining
	// capacity for the relevant workload.
	//
	// Capacity exhaustion describes workload-serving capacity such as queue slots,
	// worker concurrency, dispatch tokens, or owner-defined admission capacity. It
	// is distinct from lower-level runtime resource exhaustion.
	ReasonCapacityExhausted Reason = "capacity_exhausted"

	// ReasonResourceExhausted means that an underlying runtime or system resource
	// required by the checked scope is exhausted.
	//
	// Resource exhaustion covers generic resources such as memory, disk, file
	// descriptors, connection pools, goroutine/thread budgets, or similar
	// low-level runtime limits.
	ReasonResourceExhausted Reason = "resource_exhausted"

	// ReasonStale means that the health result is based on state that is older
	// than the owner considers acceptable.
	//
	// Staleness is useful for cached probes, control-loop heartbeats, dependency
	// snapshots, and scheduler/runtime views. The target and policy interpreting
	// the result decide whether stale state is tolerated.
	ReasonStale Reason = "stale"

	// ReasonNotSynced means that a local view, cache, snapshot, or controller
	// input has not completed its initial synchronization.
	//
	// Not-synced state often appears during startup or after ownership changes. It
	// is distinct from sync failure because synchronization may still be making
	// normal progress.
	ReasonNotSynced Reason = "not_synced"

	// ReasonSyncFailed means that synchronization was attempted and failed.
	//
	// Sync failure is useful for caches, informers, replicated views, dependency
	// snapshots, and controller inputs where the owner can distinguish a failed
	// sync attempt from a not-yet-synced initial state.
	ReasonSyncFailed Reason = "sync_failed"

	// ReasonLagging means that a component, replica, consumer, controller, or
	// local view is behind the progress expected by its owner.
	//
	// Lagging is distinct from stale: stale describes age of observed state, while
	// lagging describes progress behind another stream, source, leader, replica,
	// queue, or control-loop expectation.
	ReasonLagging Reason = "lagging"

	// ReasonPartitioned means that communication with a required peer, group, or
	// cluster segment is partitioned for the checked scope.
	//
	// Partitioned should be used only when the owner can classify an isolation or
	// reachability split at the distributed-system boundary. Ordinary dependency
	// dial failures should normally use dependency or timeout reasons instead.
	ReasonPartitioned Reason = "partitioned"

	// ReasonMisconfigured means that configuration prevents the checked scope from
	// operating correctly.
	//
	// Misconfiguration is useful for invalid required settings, incompatible local
	// options, missing required configuration, or configuration that makes startup
	// or safe operation impossible.
	ReasonMisconfigured Reason = "misconfigured"

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

// IsObservationReason reports whether r describes the health observation
// boundary rather than the checked component's own domain state.
//
// Observation reasons are produced before, during, or around a health
// observation. They are useful for distinguishing missing or interrupted
// observation from a successful observation of an unhealthy component.
func (r Reason) IsObservationReason() bool {
	switch r {
	case ReasonNotObserved,
		ReasonTimeout,
		ReasonCanceled,
		ReasonPanic:
		return true
	default:
		return false
	}
}

// IsExecutionReason reports whether r describes failure or interruption while
// producing a health observation.
//
// Execution reasons describe checker/evaluator execution, not the checked
// component's own operational state. ReasonNotObserved is observation-related but
// not execution-related because no check execution necessarily occurred.
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

// IsDependencyReason reports whether r describes dependency availability or
// dependency quality.
//
// Dependency reasons describe the relationship between the checked scope and a
// required external or internal dependency. They do not identify the dependency;
// check names, messages, diagnostics, or domain-specific reasons provide that
// additional context.
func (r Reason) IsDependencyReason() bool {
	switch r {
	case ReasonDependencyUnavailable,
		ReasonDependencyDegraded:
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
		ReasonBackpressured,
		ReasonRateLimited,
		ReasonAdmissionClosed,
		ReasonCapacityExhausted,
		ReasonResourceExhausted:
		return true
	default:
		return false
	}
}

// IsFreshnessReason reports whether r describes stale, unsynchronized, failed
// synchronization, or lagging local state.
//
// Freshness reasons are useful for cached probes, replicated views, dependency
// snapshots, control-loop inputs, event streams, and scheduler/runtime views.
func (r Reason) IsFreshnessReason() bool {
	switch r {
	case ReasonStale,
		ReasonNotSynced,
		ReasonSyncFailed,
		ReasonLagging:
		return true
	default:
		return false
	}
}

// IsConnectivityReason reports whether r describes distributed communication
// partitioning.
//
// Connectivity reasons should remain coarse in health core. Detailed network,
// TLS, DNS, or transport failures belong in domain-specific diagnostics or Cause.
func (r Reason) IsConnectivityReason() bool {
	return r == ReasonPartitioned
}

// IsConfigurationReason reports whether r describes configuration that prevents
// correct operation.
//
// Configuration reasons are useful when startup or safe operation is blocked by
// invalid, missing, or incompatible local configuration.
func (r Reason) IsConfigurationReason() bool {
	return r == ReasonMisconfigured
}
