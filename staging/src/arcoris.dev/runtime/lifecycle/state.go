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

package lifecycle

// State identifies the externally observable lifecycle phase of a component.
//
// A State describes where a component is in its start-run-stop lifecycle. It is
// intentionally coarse-grained: it records the component lifecycle phase, not the
// component's health, readiness, retry policy, pause mode, drain progress,
// reload status, supervisor decision, or internal execution step.
//
// The zero value is StateNew. This makes lifecycle state embeddable in larger
// runtime structs without constructor-heavy initialization while still having a
// meaningful and valid initial phase.
//
// State values are suitable for snapshots, transition records, guards, wait
// predicates, diagnostics, and simple state checks. They are not stable external
// wire-format values; if a future API needs serialized lifecycle status, that
// status should define its own compatibility contract rather than serializing
// this runtime enum directly.
type State uint8

const (
	// StateNew means the component has been created but its lifecycle has not
	// started yet.
	//
	// StateNew is the zero value and the only valid initial state for a lifecycle
	// controller. It is distinct from StateStopped: a new component has not run a
	// shutdown sequence, has not reached a terminal state, and may still begin its
	// first startup transition.
	StateNew State = iota

	// StateStarting means the component is executing its startup sequence.
	//
	// A starting component is active, but it MUST NOT be treated as ready for
	// normal workload. Startup may allocate resources, start background loops,
	// connect to dependencies, register handlers, or perform other component-owned
	// initialization before the component can be marked as running.
	StateStarting

	// StateRunning means the component has completed startup and is in its normal
	// operational lifecycle phase.
	//
	// StateRunning is the only lifecycle state that accepts normal workload by
	// default. It does not imply that the component is healthy, ready for every
	// caller, or operating without degradation. Health, readiness, and degradation
	// are higher-level operational signals and must be represented outside State.
	StateRunning

	// StateStopping means the component is executing its shutdown sequence.
	//
	// A stopping component is active, but it MUST NOT accept new normal workload by
	// default. It may still be draining in-flight work, stopping background loops,
	// closing resources, or publishing final runtime state before reaching a
	// terminal state.
	StateStopping

	// StateStopped means the component has completed a normal shutdown and reached
	// a successful terminal lifecycle state.
	//
	// StateStopped is terminal. A lifecycle controller that reaches StateStopped
	// MUST NOT be started again. Restart orchestration belongs to a supervisor or
	// owner that creates a fresh component/lifecycle instance.
	StateStopped

	// StateFailed means the component reached an unsuccessful terminal lifecycle
	// state.
	//
	// StateFailed is terminal. It records that startup, normal operation, or
	// shutdown failed in a way that ended this lifecycle instance. The failure
	// cause belongs to the transition or snapshot that moved the component into
	// StateFailed; it is intentionally not encoded in the State value itself.
	StateFailed
)

// String returns the canonical lower-case name of s.
//
// The returned value is intended for diagnostics, tests, logs, and human-facing
// messages. It is not a versioned serialization format. Unknown values return
// "invalid" so callers never accidentally render an unknown numeric state as a
// valid lifecycle phase.
func (s State) String() string {
	switch s {
	case StateNew:
		return "new"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	case StateFailed:
		return "failed"
	default:
		return "invalid"
	}
}

// IsValid reports whether s is one of the lifecycle states defined by this
// package.
//
// IsValid is useful at package boundaries, in tests, and in defensive code that
// receives a State value from caller-controlled input. StateNew is valid because
// it is the intended zero-value state. Any value outside the declared state set
// is invalid and must not be treated as StateNew, StateStopped, or StateFailed.
func (s State) IsValid() bool {
	switch s {
	case StateNew,
		StateStarting,
		StateRunning,
		StateStopping,
		StateStopped,
		StateFailed:
		return true
	default:
		return false
	}
}

// IsTerminal reports whether s ends a lifecycle instance.
//
// Terminal states are final for a controller instance. Once a component reaches a
// terminal state, the same lifecycle instance MUST NOT transition back to
// StateNew, StateStarting, or StateRunning. A restart is represented by a new
// component or a new lifecycle controller, not by reusing a terminal instance.
func (s State) IsTerminal() bool {
	switch s {
	case StateStopped,
		StateFailed:
		return true
	default:
		return false
	}
}

// IsActive reports whether s belongs to a lifecycle that has started but has not
// reached a terminal state.
//
// Active states include transitional phases as well as the normal running phase.
// StateNew is not active because startup has not begun. StateStopped and
// StateFailed are not active because the lifecycle instance is already terminal.
func (s State) IsActive() bool {
	switch s {
	case StateStarting,
		StateRunning,
		StateStopping:
		return true
	default:
		return false
	}
}

// AcceptsWork reports whether a component in s may accept normal workload by
// default.
//
// Only StateRunning accepts normal workload. StateStarting may still be creating
// resources or connecting dependencies. StateStopping may still be draining or
// closing resources. StateNew, StateStopped, and StateFailed do not represent an
// operational component.
//
// AcceptsWork is deliberately conservative. A component with stricter admission,
// readiness, health, or workload-specific policy may reject work even while its
// lifecycle state is StateRunning.
func (s State) AcceptsWork() bool {
	return s == StateRunning
}
