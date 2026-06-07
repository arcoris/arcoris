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
