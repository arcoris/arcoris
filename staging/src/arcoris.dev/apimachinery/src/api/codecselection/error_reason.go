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

package codecselection

// ErrorReason gives stable machine-readable detail inside broad selection errors.
type ErrorReason string

const (
	// ErrorReasonInvalidContentType reports malformed content type key material.
	ErrorReasonInvalidContentType ErrorReason = "invalid_content_type"

	// ErrorReasonInvalidParameters reports malformed content type parameters.
	ErrorReasonInvalidParameters ErrorReason = "invalid_parameters"

	// ErrorReasonInvalidPreference reports malformed encode preferences.
	ErrorReasonInvalidPreference ErrorReason = "invalid_preference"

	// ErrorReasonInvalidBinding reports malformed decode or encode bindings.
	ErrorReasonInvalidBinding ErrorReason = "invalid_binding"

	// ErrorReasonDuplicateDecodeBinding reports duplicate decode binding keys.
	ErrorReasonDuplicateDecodeBinding ErrorReason = "duplicate_decode_binding"

	// ErrorReasonDuplicateEncodeBinding reports duplicate encode binding keys.
	ErrorReasonDuplicateEncodeBinding ErrorReason = "duplicate_encode_binding"

	// ErrorReasonUnknownEntryID reports bindings to absent registry entries.
	ErrorReasonUnknownEntryID ErrorReason = "unknown_entry_id"

	// ErrorReasonEntryMediaTypeMismatch reports a media type not declared by entry.
	ErrorReasonEntryMediaTypeMismatch ErrorReason = "entry_media_type_mismatch"

	// ErrorReasonEntryTargetMismatch reports a target not declared by entry.
	ErrorReasonEntryTargetMismatch ErrorReason = "entry_target_mismatch"

	// ErrorReasonEntryCapabilityMismatch reports a missing byte or stream capability.
	ErrorReasonEntryCapabilityMismatch ErrorReason = "entry_capability_mismatch"

	// ErrorReasonNoDecodeBinding reports runtime decode misses.
	ErrorReasonNoDecodeBinding ErrorReason = "no_decode_binding"

	// ErrorReasonNoEncodePreference reports runtime encode preference misses.
	ErrorReasonNoEncodePreference ErrorReason = "no_encode_preference"
)
