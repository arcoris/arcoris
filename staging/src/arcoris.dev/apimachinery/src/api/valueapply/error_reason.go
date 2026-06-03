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

package valueapply

// ErrorReason is the stable machine-readable reason for a valueapply error.
type ErrorReason string

const (
	// ErrorReasonInvalidRequest reports a request-level invariant failure.
	ErrorReasonInvalidRequest ErrorReason = "invalid_request"

	// ErrorReasonInvalidPath reports malformed semantic path input.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonInvalidOwner reports malformed field owner input.
	ErrorReasonInvalidOwner ErrorReason = "invalid_owner"

	// ErrorReasonInvalidValue reports validation failure for Live or Applied.
	ErrorReasonInvalidValue ErrorReason = "invalid_value"

	// ErrorReasonFieldSetFailed reports failure while extracting Applied fields.
	ErrorReasonFieldSetFailed ErrorReason = "fieldset_failed"

	// ErrorReasonCompareFailed reports failure while comparing Live and Applied.
	ErrorReasonCompareFailed ErrorReason = "compare_failed"

	// ErrorReasonConflict reports ownership conflict on changed applied fields.
	ErrorReasonConflict ErrorReason = "conflict"

	// ErrorReasonUnsupportedTakeover reports a forced conflict that cannot be
	// represented without over-removing another owner's ancestor field.
	ErrorReasonUnsupportedTakeover ErrorReason = "unsupported_takeover"

	// ErrorReasonMergeFailed reports failure while running valuemerge.
	ErrorReasonMergeFailed ErrorReason = "merge_failed"
)
