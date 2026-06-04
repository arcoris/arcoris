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

package codecregistry

// ErrorReason gives stable machine-readable detail inside broad registry errors.
type ErrorReason string

const (
	// ErrorReasonInvalidCodec reports a nil or unusable codec implementation.
	ErrorReasonInvalidCodec ErrorReason = "invalid_codec"

	// ErrorReasonInvalidInfo reports malformed codec metadata.
	ErrorReasonInvalidInfo ErrorReason = "invalid_info"

	// ErrorReasonDuplicateMediaType reports an ambiguous media type index.
	ErrorReasonDuplicateMediaType ErrorReason = "duplicate_media_type"

	// ErrorReasonCapabilityMismatch reports disagreement between metadata
	// targets and implemented codec capability interfaces.
	ErrorReasonCapabilityMismatch ErrorReason = "capability_mismatch"
)
