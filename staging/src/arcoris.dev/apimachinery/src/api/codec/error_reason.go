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

package codec

// ErrorReason gives stable machine-readable detail inside broad codec errors.
type ErrorReason string

const (
	// ErrorReasonInvalidFormat reports malformed codec format identifiers.
	ErrorReasonInvalidFormat ErrorReason = "invalid_format"

	// ErrorReasonInvalidMediaType reports malformed codec media types.
	ErrorReasonInvalidMediaType ErrorReason = "invalid_media_type"

	// ErrorReasonInvalidTarget reports malformed or unknown codec targets.
	ErrorReasonInvalidTarget ErrorReason = "invalid_target"

	// ErrorReasonInvalidInfo reports malformed codec implementation metadata.
	ErrorReasonInvalidInfo ErrorReason = "invalid_info"

	// ErrorReasonUnsupportedFormat reports a format unsupported by a codec.
	ErrorReasonUnsupportedFormat ErrorReason = "unsupported_format"

	// ErrorReasonUnsupportedMediaType reports a media type unsupported by a codec.
	ErrorReasonUnsupportedMediaType ErrorReason = "unsupported_media_type"

	// ErrorReasonUnsupportedTarget reports a target unsupported by a codec.
	ErrorReasonUnsupportedTarget ErrorReason = "unsupported_target"

	// ErrorReasonUnsupportedFeature reports unavailable optional codec behavior.
	ErrorReasonUnsupportedFeature ErrorReason = "unsupported_feature"

	// ErrorReasonDecodeFailed reports failure while decoding data.
	ErrorReasonDecodeFailed ErrorReason = "decode_failed"

	// ErrorReasonEncodeFailed reports failure while encoding data.
	ErrorReasonEncodeFailed ErrorReason = "encode_failed"

	// ErrorReasonInvalidDocument reports malformed encoded document content.
	ErrorReasonInvalidDocument ErrorReason = "invalid_document"

	// ErrorReasonStrictViolation reports input rejected by strict mode.
	ErrorReasonStrictViolation ErrorReason = "strict_violation"

	// ErrorReasonMaxDepthExceeded reports input deeper than the configured limit.
	ErrorReasonMaxDepthExceeded ErrorReason = "max_depth_exceeded"

	// ErrorReasonInvalidNumber reports an encoded number outside supported shape.
	ErrorReasonInvalidNumber ErrorReason = "invalid_number"
)
