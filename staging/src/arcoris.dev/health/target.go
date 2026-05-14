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

// Target identifies the health scope being evaluated.
//
// A Target describes why a health report is being produced. It intentionally
// does not describe how the report is transported. HTTP paths such as /livez,
// /readyz, and /startupz, gRPC service names, CLI output modes, and diagnostic
// endpoints belong to adapter packages outside health.
//
// The zero value is TargetUnknown. This makes Target safe to embed in larger
// runtime structs before a concrete health scope is selected. TargetUnknown is a
// valid sentinel value, but it is not an evaluable target. Registries,
// evaluators, and adapters MUST reject TargetUnknown when a concrete health
// evaluation target is required.
//
// Target values are suitable for in-memory runtime state, reports, registry
// keys, tests, diagnostics, and policy selection. They are not a versioned wire
// format. Transport packages MUST define their own compatibility contract when
// exposing target names outside the process.
type Target uint8

const (
	// TargetUnknown means that no concrete health target has been selected.
	//
	// TargetUnknown is the zero value. It is useful as an unset sentinel at
	// package boundaries and inside larger configuration or runtime structs.
	// TargetUnknown MUST NOT be registered, evaluated, or exposed as a concrete
	// health endpoint.
	TargetUnknown Target = iota

	// TargetStartup asks whether the component has completed its startup or
	// bootstrap sequence.
	//
	// Startup is concerned with initialization progress. A component may be live
	// while startup is still incomplete, but it MUST NOT be assumed ready for
	// normal workload until startup-sensitive policy has accepted the startup
	// target.
	TargetStartup

	// TargetLive asks whether the component is still able to make progress and
	// should continue running.
	//
	// Liveness is a restart-oriented health scope. Temporary dependency outages,
	// overload, backpressure, admission closure, and graceful draining should
	// normally affect readiness rather than liveness. TargetLive should be
	// reserved for fatal or progress-breaking conditions such as deadlocks,
	// unrecoverable runtime corruption, or failed critical control loops.
	TargetLive

	// TargetReady asks whether the component should receive new work.
	//
	// Readiness is an admission-oriented health scope. It may reflect startup
	// completion, graceful draining, overload protection, dependency availability,
	// capacity availability, and other local conditions that determine whether
	// the component should accept normal workload.
	TargetReady
)

// String returns the canonical lower-case name of t.
//
// The returned value is intended for diagnostics, tests, logs, reports, and
// human-facing messages. It is not a versioned serialization format. Unknown
// numeric values return "invalid" so callers never accidentally render an
// unknown target as a valid health scope.
func (t Target) String() string {
	switch t {
	case TargetUnknown:
		return "unknown"
	case TargetStartup:
		return "startup"
	case TargetLive:
		return "live"
	case TargetReady:
		return "ready"
	default:
		return "invalid"
	}
}

// IsValid reports whether t is one of the target values defined by this package.
//
// TargetUnknown is valid because it is the intended zero-value sentinel. Validity
// does not mean that a target can be evaluated. Call IsConcrete when an API
// requires a real health evaluation scope.
func (t Target) IsValid() bool {
	switch t {
	case TargetUnknown,
		TargetStartup,
		TargetLive,
		TargetReady:
		return true
	default:
		return false
	}
}

// IsConcrete reports whether t identifies a real health evaluation scope.
//
// Concrete targets may be registered, evaluated, and used for target-specific
// policy selection. TargetUnknown is not concrete because it only represents an
// unset or unspecified target.
func (t Target) IsConcrete() bool {
	switch t {
	case TargetStartup,
		TargetLive,
		TargetReady:
		return true
	default:
		return false
	}
}

// IsStartup reports whether t is the startup health target.
//
// The startup target is used to evaluate whether component bootstrap has
// completed. It should not be confused with lifecycle state: lifecycle records
// where the component is in its start-run-stop sequence, while health startup
// reports whether startup-related health checks have accepted the component.
func (t Target) IsStartup() bool {
	return t == TargetStartup
}

// IsLive reports whether t is the liveness health target.
//
// Liveness answers whether the component should continue running. It does not
// answer whether the component should receive new work.
func (t Target) IsLive() bool {
	return t == TargetLive
}

// IsReady reports whether t is the readiness health target.
//
// Readiness answers whether the component should receive new work. It does not
// answer whether the component should be restarted.
func (t Target) IsReady() bool {
	return t == TargetReady
}

// ConcreteTargets returns the built-in concrete targets in deterministic order.
//
// The returned slice is newly allocated so callers may sort, append to, or modify
// it without mutating package-level state. TargetUnknown is intentionally omitted
// because it is a sentinel value, not an evaluable target.
func ConcreteTargets() []Target {
	return []Target{
		TargetStartup,
		TargetLive,
		TargetReady,
	}
}
