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

package admission

// Outcome classifies what happened to one admission attempt.
//
// Outcome is closed because Result validation depends on the exact semantics of
// each outcome. Use Reason for open-world domain-specific detail.
type Outcome uint8

const (
	// OutcomeUnknown is the zero value and is invalid.
	OutcomeUnknown Outcome = iota

	// OutcomeAdmitted means work may proceed now.
	OutcomeAdmitted

	// OutcomeDenied means the current admission attempt was rejected.
	OutcomeDenied

	// OutcomeQueued means the system accepted ownership of waiting work.
	OutcomeQueued

	// OutcomeDeferred means the work was not accepted now and the caller retains
	// responsibility for any later retry.
	OutcomeDeferred
)

// IsValid reports whether o is a defined non-zero outcome.
func (o Outcome) IsValid() bool {
	switch o {
	case OutcomeAdmitted, OutcomeDenied, OutcomeQueued, OutcomeDeferred:
		return true
	default:
		return false
	}
}

// IsAdmitted reports whether o admits work immediately.
//
// This is the success state for admission. It says nothing about whether the
// decision also committed a side effect or returned a grant; use Effect helpers
// or Result helpers for that ownership shape.
func (o Outcome) IsAdmitted() bool {
	return o == OutcomeAdmitted
}

// IsDenied reports whether o rejects the current admission attempt.
//
// Denial is terminal for the current attempt and does not transfer waiting
// ownership to the system.
func (o Outcome) IsDenied() bool {
	return o == OutcomeDenied
}

// IsQueued reports whether o accepted system-owned waiting.
//
// Queued outcomes are intentionally not terminal: the system now owns some
// waiting state and may expose a queue handle through Result.
func (o Outcome) IsQueued() bool {
	return o == OutcomeQueued
}

// IsDeferred reports whether o leaves retry ownership with the caller.
func (o Outcome) IsDeferred() bool {
	return o == OutcomeDeferred
}

// IsTerminal reports whether o completes the current admission attempt without
// leaving system-owned waiting behind.
func (o Outcome) IsTerminal() bool {
	return o == OutcomeAdmitted || o == OutcomeDenied || o == OutcomeDeferred
}

// String returns the stable machine-readable outcome name.
//
// Undefined values format as "unknown" so diagnostics remain safe even when a
// caller constructs an invalid Outcome manually.
func (o Outcome) String() string {
	switch o {
	case OutcomeAdmitted:
		return "admitted"
	case OutcomeDenied:
		return "denied"
	case OutcomeQueued:
		return "queued"
	case OutcomeDeferred:
		return "deferred"
	default:
		return "unknown"
	}
}
