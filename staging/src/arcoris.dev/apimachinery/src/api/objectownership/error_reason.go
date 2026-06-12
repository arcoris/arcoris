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

package objectownership

// ErrorReason gives stable machine-readable detail inside broad errors.
type ErrorReason string

const (
	// ErrorReasonInvalidState reports malformed top-level ownership state.
	ErrorReasonInvalidState ErrorReason = "invalid_state"

	// ErrorReasonInvalidDesired reports malformed Desired ownership state.
	ErrorReasonInvalidDesired ErrorReason = "invalid_desired"

	// ErrorReasonInvalidObserved reports malformed Observed ownership state.
	ErrorReasonInvalidObserved ErrorReason = "invalid_observed"

	// ErrorReasonInvalidMetadataLabels reports malformed labels ownership state.
	ErrorReasonInvalidMetadataLabels ErrorReason = "invalid_metadata_labels"

	// ErrorReasonInvalidMetadataAnnotations reports malformed annotations ownership state.
	ErrorReasonInvalidMetadataAnnotations ErrorReason = "invalid_metadata_annotations"

	// ErrorReasonNotNormalized reports valid states that are not canonical.
	ErrorReasonNotNormalized ErrorReason = "not_normalized"
)
