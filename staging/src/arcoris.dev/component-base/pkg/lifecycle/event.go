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

// Event identifies a lifecycle input that may move a component from one State to
// another.
//
// Events are intentionally named as lifecycle facts or lifecycle requests, not
// as target states. A State describes where the component is. An Event describes
// what happened to the lifecycle controller or what the owner is attempting to
// record.
//
// The lifecycle model is deliberately two-phase for normal startup and shutdown:
//
//   - EventBeginStart records that startup has begun;
//   - EventMarkRunning records that startup completed successfully;
//   - EventBeginStop records that shutdown has begun;
//   - EventMarkStopped records that shutdown completed successfully.
//
// EventMarkFailed records an unsuccessful terminal outcome. It is separate from
// EventBeginStop and EventMarkStopped so callers cannot accidentally mask a
// failed startup, failed runtime, or failed shutdown as a normal stop.
//
// Event values are suitable for transition records, guards, wait diagnostics,
// structured errors, and tests. They are not stable external wire-format values;
// if a future API needs serialized lifecycle events, that API should define its
// own compatibility contract rather than serializing this runtime enum directly.
type Event uint8

const (
	// EventBeginStart records that the component owner has begun the startup
	// sequence.
	//
	// The event normally moves a lifecycle from StateNew to StateStarting. It
	// does not execute startup work by itself. The component owner remains
	// responsible for opening resources, starting background execution, connecting
	// dependencies, and marking the component as running only after startup has
	// completed successfully.
	EventBeginStart Event = iota

	// EventMarkRunning records that the component startup sequence has completed
	// successfully.
	//
	// The event normally moves a lifecycle from StateStarting to StateRunning.
	// It SHOULD only be applied after all startup invariants required for normal
	// operation have been established. StateRunning may accept normal workload by
	// default, but health and readiness remain separate operational signals.
	EventMarkRunning

	// EventBeginStop records that the component owner has begun shutdown.
	//
	// The event normally moves a running lifecycle from StateRunning to
	// StateStopping. It may also be valid from StateStarting when shutdown begins
	// before startup completes. A lifecycle in StateNew may use this event to
	// reach StateStopped without running a shutdown sequence, allowing Stop or
	// Close operations to be safe before Start.
	EventBeginStop

	// EventMarkStopped records that the component shutdown sequence has completed
	// successfully.
	//
	// The event normally moves a lifecycle from StateStopping to StateStopped.
	// It SHOULD only be applied after the component has stopped accepting normal
	// workload, drained or cancelled component-owned execution as appropriate, and
	// released the resources owned by the shutdown path.
	EventMarkStopped

	// EventMarkFailed records that the component lifecycle has reached an
	// unsuccessful terminal outcome.
	//
	// The event normally moves a lifecycle from StateStarting, StateRunning, or
	// StateStopping to StateFailed. A transition using EventMarkFailed MUST carry
	// a non-nil failure cause in the transition/controller layer. The cause is not
	// encoded in the Event value itself because Event is only the lifecycle input,
	// not the diagnostic payload.
	EventMarkFailed
)

// String returns the canonical lower-case name of e.
//
// The returned value is intended for diagnostics, tests, logs, and human-facing
// messages. It is not a versioned serialization format. Unknown values return
// "invalid" so callers never accidentally render an unknown numeric event as a
// valid lifecycle input.
func (e Event) String() string {
	switch e {
	case EventBeginStart:
		return "begin_start"
	case EventMarkRunning:
		return "mark_running"
	case EventBeginStop:
		return "begin_stop"
	case EventMarkStopped:
		return "mark_stopped"
	case EventMarkFailed:
		return "mark_failed"
	default:
		return "invalid"
	}
}

// IsValid reports whether e is one of the lifecycle events defined by this
// package.
//
// IsValid is useful at package boundaries, in transition tests, and in defensive
// code that receives an Event value from caller-controlled input. Any value
// outside the declared event set is invalid and must not be treated as a no-op,
// retry, stop, or failure event.
func (e Event) IsValid() bool {
	switch e {
	case EventBeginStart,
		EventMarkRunning,
		EventBeginStop,
		EventMarkStopped,
		EventMarkFailed:
		return true
	default:
		return false
	}
}

// IsStartEvent reports whether e belongs to the startup side of the lifecycle.
//
// Startup events describe the beginning or successful completion of startup.
// They do not include EventMarkFailed because failure can occur during startup,
// runtime, or shutdown and is not itself a startup event.
func (e Event) IsStartEvent() bool {
	switch e {
	case EventBeginStart,
		EventMarkRunning:
		return true
	default:
		return false
	}
}

// IsStopEvent reports whether e belongs to the normal shutdown side of the
// lifecycle.
//
// Stop events describe the beginning or successful completion of normal
// shutdown. They do not include EventMarkFailed because failure is an
// unsuccessful terminal outcome, not a normal shutdown event.
func (e Event) IsStopEvent() bool {
	switch e {
	case EventBeginStop,
		EventMarkStopped:
		return true
	default:
		return false
	}
}

// IsTerminalEvent reports whether e records a terminal lifecycle outcome.
//
// Terminal events are events that, when accepted by the transition table, move a
// lifecycle instance into a terminal state. EventMarkStopped records successful
// termination. EventMarkFailed records unsuccessful termination.
func (e Event) IsTerminalEvent() bool {
	switch e {
	case EventMarkStopped,
		EventMarkFailed:
		return true
	default:
		return false
	}
}

// RequiresCause reports whether transitions using e must carry a non-nil cause.
//
// Only EventMarkFailed requires a cause. Normal startup and shutdown events
// should not use an error cause to carry diagnostics. If a normal lifecycle phase
// cannot complete successfully, the controller should record EventMarkFailed
// with the underlying failure cause instead of attaching an error to a successful
// event.
func (e Event) RequiresCause() bool {
	return e == EventMarkFailed
}
