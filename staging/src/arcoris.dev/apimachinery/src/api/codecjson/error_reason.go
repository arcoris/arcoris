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

package codecjson

// ErrorReason gives stable machine-readable detail inside JSON codec errors.
type ErrorReason string

const (
	// ErrorReasonInvalidJSON reports malformed JSON syntax or token shape.
	ErrorReasonInvalidJSON ErrorReason = "invalid_json"

	// ErrorReasonDuplicateKey reports a repeated JSON object member name.
	ErrorReasonDuplicateKey ErrorReason = "duplicate_key"

	// ErrorReasonTrailingData reports data after the first JSON document.
	ErrorReasonTrailingData ErrorReason = "trailing_data"

	// ErrorReasonUnsupportedValue reports a value kind JSON cannot round-trip.
	ErrorReasonUnsupportedValue ErrorReason = "unsupported_value"

	// ErrorReasonInvalidNumber reports an unsupported JSON number literal.
	ErrorReasonInvalidNumber ErrorReason = "invalid_number"

	// ErrorReasonInvalidEnvelope reports malformed object or ownership envelopes.
	ErrorReasonInvalidEnvelope ErrorReason = "invalid_envelope"

	// ErrorReasonMaxDepthExceeded reports input deeper than configured limits.
	ErrorReasonMaxDepthExceeded ErrorReason = "max_depth_exceeded"
)
