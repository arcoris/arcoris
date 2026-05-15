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

package liveconfig

// ChangeReason classifies the outcome of one Holder.Apply attempt.
//
// ChangeReason is deliberately narrower than an error. It answers the stable
// operational question "which branch of the holder pipeline produced this
// result?" while the error returned by Apply carries the detailed diagnostic
// payload from a normalizer or validator. This separation lets reload loops,
// logs, and future adapters group outcomes without parsing or comparing error
// values.
//
// The reasons model only the root holder pipeline:
//
//	clone -> normalize -> validate -> equal check -> publish
//
// Source-oriented failures such as loading a file, reading an environment
// variable, decoding JSON, watching Kubernetes ConfigMaps, or calling a remote
// control plane are intentionally not represented here. Those failures happen in
// source/reload adapters layered above liveconfig, and those adapters can expose
// their own classifications without expanding the holder primitive.
type ChangeReason uint8

const (
	// ChangeReasonUnknown is the zero value.
	//
	// Apply does not return Unknown. It exists so a manually constructed Change
	// starts in an obviously incomplete state instead of accidentally looking like
	// a successful or rejected holder result.
	ChangeReasonUnknown ChangeReason = iota

	// ChangeReasonPublished means the candidate was accepted and published as a
	// new source revision.
	//
	// This is the only reason for which Change.Changed is true. Apply returns a
	// nil error, Current is the newly published snapshot, and Current.Revision is
	// different from Previous.Revision.
	ChangeReasonPublished

	// ChangeReasonEqual means the candidate was accepted but EqualFunc classified
	// it as equivalent to the current value, so no new revision was published.
	//
	// Equal is a successful apply attempt, not a rejection. Apply returns a nil
	// error, clears LastError, and keeps Current at the same revision as Previous.
	ChangeReasonEqual

	// ChangeReasonNormalizeFailed means normalization rejected the candidate.
	// Validation and publication did not run, and the last-good config was kept.
	//
	// Apply returns the normalizer error and records it as LastError. Current and
	// Previous refer to the same snapshot revision.
	ChangeReasonNormalizeFailed

	// ChangeReasonValidateFailed means validation rejected the normalized
	// candidate. Publication did not run, and the last-good config was kept.
	//
	// Apply returns the validator error and records it as LastError. Current and
	// Previous refer to the same snapshot revision.
	ChangeReasonValidateFailed
)

// String returns a stable diagnostic name for r.
//
// The returned names are intended for logs, tests, and low-cardinality labels in
// code layered above this package. Unknown numeric values stringify as
// "unknown" so forward-incompatible or manually constructed values fail closed.
func (r ChangeReason) String() string {
	switch r {
	case ChangeReasonPublished:
		return "published"
	case ChangeReasonEqual:
		return "equal"
	case ChangeReasonNormalizeFailed:
		return "normalize_failed"
	case ChangeReasonValidateFailed:
		return "validate_failed"
	default:
		return "unknown"
	}
}

// IsValid reports whether r is one of the concrete Apply outcomes.
//
// Unknown is not valid because it is a zero/default sentinel rather than a
// branch that Apply returns. Unknown numeric values are also invalid.
func (r ChangeReason) IsValid() bool {
	switch r {
	case ChangeReasonPublished,
		ChangeReasonEqual,
		ChangeReasonNormalizeFailed,
		ChangeReasonValidateFailed:
		return true
	default:
		return false
	}
}

// Accepted reports whether Apply accepted the candidate.
//
// Equal is accepted: it means the candidate passed the holder pipeline and was
// intentionally not published because it matched the current value.
func (r ChangeReason) Accepted() bool {
	switch r {
	case ChangeReasonPublished, ChangeReasonEqual:
		return true
	default:
		return false
	}
}

// Rejected reports whether Apply rejected the candidate and preserved last-good.
//
// Rejection reasons always correspond to a non-nil Apply error. They do not
// publish a snapshot and do not advance the holder revision.
func (r ChangeReason) Rejected() bool {
	switch r {
	case ChangeReasonNormalizeFailed, ChangeReasonValidateFailed:
		return true
	default:
		return false
	}
}

// Published reports whether Apply published a new source revision.
//
// This is the reason-level counterpart to Change.Changed.
func (r ChangeReason) Published() bool {
	return r == ChangeReasonPublished
}

// Equal reports whether Apply accepted an equal no-op candidate.
//
// Equal candidates are successful no-ops: they preserve revision but clear
// LastError because the most recent apply attempt was valid.
func (r ChangeReason) Equal() bool {
	return r == ChangeReasonEqual
}
