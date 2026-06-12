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

import "arcoris.dev/apimachinery/api/codecjson/jsonconfig"

// resolvedConfig is the immutable runtime policy stored by Codec.
type resolvedConfig struct {
	// decode stores parser and target-decoder policy.
	decode resolvedDecodeConfig

	// encode stores writer and target-encoder policy.
	encode resolvedEncodeConfig
}

// resolvedDecodeConfig contains only decode fields used on the hot path.
type resolvedDecodeConfig struct {
	// maxDepth is the inclusive JSON document nesting limit.
	maxDepth int

	// maxDocumentBytes bounds raw input reads. Zero means unlimited.
	maxDocumentBytes int64

	// maxStringBytes bounds decoded string token length in bytes. Zero means unlimited.
	maxStringBytes int

	// maxNumberDigits bounds exact number parsing and decimal expansion.
	maxNumberDigits int

	// rejectUnknownEnvelopeFields controls fixed object envelope field checks.
	rejectUnknownEnvelopeFields bool

	// rejectUnknownOwnershipFields controls fixed ownership-state field checks.
	rejectUnknownOwnershipFields bool

	// validateOwnershipState controls objectownership.Validate on decode.
	validateOwnershipState bool
}

// resolvedEncodeConfig contains only encode fields used on the hot path.
type resolvedEncodeConfig struct {
	// pretty selects compact or indented node output.
	pretty bool

	// indent is the exact indentation unit used for pretty output.
	indent string

	// finalNewline appends a trailing newline after one JSON document.
	finalNewline bool

	// deterministic sorts value object members and may normalize ownership state.
	deterministic bool

	// escapeHTML mirrors encoding/json HTML escaping for JSON strings.
	escapeHTML bool

	// maxDepth is the inclusive JSON output nesting limit.
	maxDepth int

	// maxOutputBytes bounds encoded output bytes. Zero means unlimited.
	maxOutputBytes int64

	// maxNumberDigits bounds emitted JSON number text.
	maxNumberDigits int

	// floatFormat controls whether value.KindFloat is accepted.
	floatFormat jsonconfig.FloatFormatMode

	// negativeZero controls negative floating zero output.
	negativeZero jsonconfig.NegativeZeroMode

	// typeMeta controls object envelope apiVersion/kind emission.
	typeMeta jsonconfig.TypeMetaEncodeMode

	// metadata controls object envelope metadata emission.
	metadata jsonconfig.MetadataEncodeMode

	// observed controls absent observed payload emission.
	observed jsonconfig.ObservedEncodeMode

	// ownershipNormalize controls pre-encode ownership state normalization.
	ownershipNormalize jsonconfig.OwnershipNormalizeMode

	// emptySurfaces controls empty ownership surface emission.
	emptySurfaces jsonconfig.EmptyOwnershipSurfaceMode

	// emptyEntries controls empty ownership entries emission.
	emptyEntries jsonconfig.EmptyEntriesMode
}
