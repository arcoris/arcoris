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

// Status identifies the observed operational health state of a component,
// subsystem, or health check.
//
// A Status describes the local health observation itself. It intentionally does
// not encode transport behavior, HTTP status codes, gRPC serving states,
// readiness outcomes, restart decisions, admission policy, scheduler policy, or
// overload-control actions. Those decisions belong to target policies,
// transport adapters, and higher-level runtime owners.
//
// The zero value is StatusUnknown. This makes Status safe to embed in larger
// runtime structs before the first health observation is available. Callers MUST
// treat StatusUnknown as a valid but non-affirmative health state: it means that
// health has not been determined, not that the component is healthy.
//
// Status values are suitable for in-memory runtime state, snapshots, reports,
// tests, diagnostics, and policy evaluation. They are not a versioned wire
// format. Transport packages MUST define their own compatibility contract when
// exposing health state outside the process.
type Status uint8

const (
	// StatusUnknown means that the component's health could not be determined.
	//
	// Unknown is the zero value. It is used before the first observation, after an
	// inconclusive check, when a check is canceled or times out before producing a
	// reliable result, or when caller-controlled input cannot be mapped to a known
	// health state.
	//
	// StatusUnknown MUST NOT be treated as StatusHealthy. Target policies may
	// decide whether unknown is tolerated for a specific use case, but the core
	// status model keeps it distinct from both success and failure.
	StatusUnknown Status = iota

	// StatusStarting means that the component is bootstrapping and has not yet
	// completed the initialization required for normal operation.
	//
	// Starting is not a terminal failure. A starting component may be alive and
	// making progress, but it MUST NOT be assumed ready for normal workload by
	// default. Startup policy, readiness policy, and lifecycle integration decide
	// how this state affects externally observable targets.
	StatusStarting

	// StatusHealthy means that the component is operating normally for the scope
	// being checked.
	//
	// Healthy is a positive health observation. It does not by itself grant
	// admission, scheduling, traffic routing, or workload acceptance. Higher-level
	// policy may still reject work for reasons outside the local health check.
	StatusHealthy

	// StatusDegraded means that the component is still operating, but with reduced
	// capability, reduced confidence, or active protective behavior.
	//
	// Degraded is a first-class state because ARCORIS components may intentionally
	// reduce admission, apply backpressure, shed load, use partial capacity, or
	// operate with stale-but-usable state without being terminally unhealthy.
	//
	// StatusDegraded MUST remain distinct from StatusUnhealthy. A degraded
	// component may still be useful to the runtime, scheduler, or control plane,
	// depending on the target and policy that interpret the report.
	StatusDegraded

	// StatusUnhealthy means that the component is not healthy for the scope being
	// checked.
	//
	// Unhealthy is the strongest negative health observation. It still does not
	// prescribe a direct action such as restart, removal from traffic, or admission
	// closure. Those decisions depend on the target being evaluated and the owner
	// interpreting the report.
	StatusUnhealthy
)

// String returns the canonical lower-case name of s.
//
// The returned value is intended for diagnostics, tests, logs, reports, and
// human-facing messages. It is not a versioned serialization format. Unknown
// numeric values return "invalid" so callers never accidentally render an
// unknown status as a valid health state.
func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "unknown"
	case StatusStarting:
		return "starting"
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "invalid"
	}
}

// IsValid reports whether s is one of the health statuses defined by this
// package.
//
// StatusUnknown is valid because it is the intended zero value and represents a
// real observation state: health has not been determined. Any value outside the
// declared status set is invalid and MUST NOT be treated as StatusUnknown,
// StatusHealthy, or StatusUnhealthy.
func (s Status) IsValid() bool {
	switch s {
	case StatusUnknown,
		StatusStarting,
		StatusHealthy,
		StatusDegraded,
		StatusUnhealthy:
		return true
	default:
		return false
	}
}

// IsAffirmative reports whether s is a positive health observation.
//
// Only StatusHealthy is affirmative. StatusDegraded is deliberately not
// affirmative because it requires policy-specific interpretation: a degraded
// component may still be live, but it may be unsuitable for readiness,
// admission, scheduling, or normal workload.
func (s Status) IsAffirmative() bool {
	return s == StatusHealthy
}

// IsNegative reports whether s is a negative health observation.
//
// StatusUnhealthy is the only strictly negative status. StatusUnknown means the
// component could not be evaluated, StatusStarting means initialization is still
// in progress, and StatusDegraded means the component is operating with reduced
// capability. Target policies may still choose to fail those states for a
// specific target, but the core status model does not classify them as strictly
// negative.
func (s Status) IsNegative() bool {
	return s == StatusUnhealthy
}

// IsKnown reports whether s represents a determined operational state.
//
// StatusUnknown is valid but not known. Invalid numeric values are also not
// known. Starting, healthy, degraded, and unhealthy all represent explicit
// observations that can be interpreted by target policy.
func (s Status) IsKnown() bool {
	switch s {
	case StatusStarting,
		StatusHealthy,
		StatusDegraded,
		StatusUnhealthy:
		return true
	default:
		return false
	}
}

// IsOperational reports whether s describes a component that is still operating
// or making progress.
//
// Operational does not mean ready, admitted, routable, or fully healthy.
// StatusStarting is operational because startup may be progressing.
// StatusHealthy is operational by definition. StatusDegraded is operational
// because the component still has usable capability. StatusUnknown and
// StatusUnhealthy are not considered operational by the core status model.
func (s Status) IsOperational() bool {
	switch s {
	case StatusStarting,
		StatusHealthy,
		StatusDegraded:
		return true
	default:
		return false
	}
}

// MoreSevereThan reports whether s is more severe than other under the core
// health severity ordering.
//
// The ordering is intentionally conservative and transport-neutral:
//
//   - StatusHealthy
//   - StatusStarting
//   - StatusDegraded
//   - StatusUnknown
//   - StatusUnhealthy
//
// StatusStarting is more severe than StatusHealthy because startup is not a
// normal operating state. StatusUnknown is more severe than StatusDegraded
// because an unknown result provides no reliable operational signal. Invalid
// statuses are treated as more severe than every valid status so defensive
// aggregation does not accidentally hide corrupted or caller-controlled input.
func (s Status) MoreSevereThan(other Status) bool {
	return s.severity() > other.severity()
}

// severity returns the internal conservative ordering used by status-level
// aggregation helpers.
//
// The numeric value is intentionally private. Callers that need target-specific
// behavior MUST use policy-level APIs rather than depending on the exact
// ordering values.
func (s Status) severity() uint8 {
	switch s {
	case StatusHealthy:
		return 0
	case StatusStarting:
		return 1
	case StatusDegraded:
		return 2
	case StatusUnknown:
		return 3
	case StatusUnhealthy:
		return 4
	default:
		return 5
	}
}
