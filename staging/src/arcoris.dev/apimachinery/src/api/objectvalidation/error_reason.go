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

package objectvalidation

// ErrorReason identifies a precise object contract validation failure.
type ErrorReason string

// Error reasons refine broad sentinel errors with stable diagnostics.
const (
	// ErrorReasonInvalidPlan identifies malformed plan dependencies.
	ErrorReasonInvalidPlan ErrorReason = "invalid_plan"

	// ErrorReasonInvalidMetadata identifies object metadata validation failures.
	ErrorReasonInvalidMetadata ErrorReason = "invalid_metadata"

	// ErrorReasonInvalidTypeMeta identifies TypeMeta validation failures.
	ErrorReasonInvalidTypeMeta ErrorReason = "invalid_type_meta"

	// ErrorReasonInvalidObjectMeta identifies ObjectMeta validation failures.
	ErrorReasonInvalidObjectMeta ErrorReason = "invalid_object_meta"

	// ErrorReasonResourceMismatch identifies object/resource group-kind mismatches.
	ErrorReasonResourceMismatch ErrorReason = "resource_mismatch"

	// ErrorReasonVersionNotDefined identifies object versions missing from the resource.
	ErrorReasonVersionNotDefined ErrorReason = "version_not_defined"

	// ErrorReasonInvalidScope identifies namespace/scope compatibility failures.
	ErrorReasonInvalidScope ErrorReason = "invalid_scope"

	// ErrorReasonMissingValidator identifies required validator dependencies.
	ErrorReasonMissingValidator ErrorReason = "missing_validator"

	// ErrorReasonInvalidDesired identifies desired surface validator failures.
	ErrorReasonInvalidDesired ErrorReason = "invalid_desired"

	// ErrorReasonInvalidObserved identifies observed surface validator failures.
	ErrorReasonInvalidObserved ErrorReason = "invalid_observed"

	// ErrorReasonObservedNotAllowed identifies observed values without observed descriptors.
	ErrorReasonObservedNotAllowed ErrorReason = "observed_not_allowed"
)
