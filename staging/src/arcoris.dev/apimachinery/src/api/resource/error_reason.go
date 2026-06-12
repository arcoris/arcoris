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

package resource

// ErrorReason identifies the precise resource invariant that failed.
//
// Reasons are intentionally more specific than the broad sentinel errors. A
// caller can use errors.Is for coarse handling and inspect ErrorReason for
// diagnostics, tests, command output, or future editor tooling.
type ErrorReason string

// Resource validation and encoding error reasons.
const (
	// ErrorReasonInvalidGroup means the resource family group is not a valid
	// api/identity Group.
	ErrorReasonInvalidGroup ErrorReason = "invalid_group"
	// ErrorReasonInvalidKind means the resource family kind is not a valid
	// api/identity Kind.
	ErrorReasonInvalidKind ErrorReason = "invalid_kind"
	// ErrorReasonInvalidResource means the resource collection name is not a
	// valid api/identity Resource.
	ErrorReasonInvalidResource ErrorReason = "invalid_resource"
	// ErrorReasonInvalidScope means a Scope value is not one of the supported
	// resource descriptor scopes.
	ErrorReasonInvalidScope ErrorReason = "invalid_scope"
	// ErrorReasonInvalidJSON means scalar JSON decoding received a non-string,
	// null, or malformed JSON value.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"
	// ErrorReasonNilReceiver means a text or JSON decoder was called on a nil
	// receiver.
	ErrorReasonNilReceiver ErrorReason = "nil_receiver"
	// ErrorReasonMissingVersion means a resource family has no version
	// descriptors.
	ErrorReasonMissingVersion ErrorReason = "missing_version"
	// ErrorReasonDuplicateVersion means one version name appears more than once
	// in a resource family.
	ErrorReasonDuplicateVersion ErrorReason = "duplicate_version"
	// ErrorReasonNoExposedVersion means no version is marked as exposed.
	ErrorReasonNoExposedVersion ErrorReason = "no_exposed_version"
	// ErrorReasonNoCanonicalVersion means no version is marked as canonical.
	ErrorReasonNoCanonicalVersion ErrorReason = "no_canonical_version"
	// ErrorReasonMultipleCanonical means more than one version is marked as
	// canonical.
	ErrorReasonMultipleCanonical ErrorReason = "multiple_canonical_versions"
	// ErrorReasonCanonicalNotExposed means the canonical version is not also
	// exposed.
	ErrorReasonCanonicalNotExposed ErrorReason = "canonical_version_not_exposed"
	// ErrorReasonInvalidVersion means a VersionDefinition uses an invalid
	// api/identity Version.
	ErrorReasonInvalidVersion ErrorReason = "invalid_version"
	// ErrorReasonMissingDesired means a VersionDefinition has no Desired
	// structural descriptor.
	ErrorReasonMissingDesired ErrorReason = "missing_desired"
	// ErrorReasonInvalidDesired means the Desired descriptor failed api/types
	// structural validation.
	ErrorReasonInvalidDesired ErrorReason = "invalid_desired"
	// ErrorReasonDesiredNotObject means Desired is neither an object nor a
	// resolver-proven reference to an object.
	ErrorReasonDesiredNotObject ErrorReason = "desired_not_object"
	// ErrorReasonInvalidObserved means the optional Observed descriptor failed
	// api/types structural validation.
	ErrorReasonInvalidObserved ErrorReason = "invalid_observed"
	// ErrorReasonObservedNotObject means Observed is neither an object nor a
	// resolver-proven reference to an object.
	ErrorReasonObservedNotObject ErrorReason = "observed_not_object"
)
