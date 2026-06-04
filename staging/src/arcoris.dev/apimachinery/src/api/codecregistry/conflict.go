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

import "arcoris.dev/apimachinery/api/codec"

// checkEntryConflicts rejects ambiguous media type indexes.
//
// Formats are grouping attributes and may repeat. Media types are unique lookup
// keys, so duplicates are rejected to keep exact lookup deterministic.
func checkEntryConflicts(
	index int,
	entry Entry,
	seenMediaTypes map[codec.MediaType]int,
) error {
	for mediaTypeIndex, mediaType := range entry.info.MediaTypes {
		if previous, ok := seenMediaTypes[mediaType]; ok {
			return duplicateMediaTypeError(index, mediaTypeIndex, previous, mediaType)
		}
		seenMediaTypes[mediaType] = index
	}

	return nil
}

// duplicateMediaTypeError creates a stable duplicate media type diagnostic.
func duplicateMediaTypeError(
	index int,
	mediaTypeIndex int,
	previous int,
	mediaType codec.MediaType,
) error {
	return errorfAt(
		mediaTypePath(index, mediaTypeIndex),
		ErrDuplicateMediaType,
		ErrorReasonDuplicateMediaType,
		"codec media type %q duplicates codecs[%d]",
		mediaType,
		previous,
	)
}
