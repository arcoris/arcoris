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
	// ErrorReasonInvalidDocument reports malformed top-level document shape.
	ErrorReasonInvalidDocument ErrorReason = "invalid_document"

	// ErrorReasonMissingVersion reports documents without an explicit version.
	ErrorReasonMissingVersion ErrorReason = "missing_version"

	// ErrorReasonUnsupportedVersion reports an unknown document version.
	ErrorReasonUnsupportedVersion ErrorReason = "unsupported_version"

	// ErrorReasonInvalidSurface reports malformed surface ownership shape.
	ErrorReasonInvalidSurface ErrorReason = "invalid_surface"

	// ErrorReasonInvalidEntry reports malformed owner/path entry shape.
	ErrorReasonInvalidEntry ErrorReason = "invalid_entry"

	// ErrorReasonInvalidOwner reports malformed owner identity text.
	ErrorReasonInvalidOwner ErrorReason = "invalid_owner"

	// ErrorReasonInvalidPath reports malformed document path text.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonNotNormalized reports valid documents that are not canonical.
	ErrorReasonNotNormalized ErrorReason = "not_normalized"
)
