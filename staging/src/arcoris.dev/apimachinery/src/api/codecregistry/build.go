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

import (
	"slices"

	"arcoris.dev/apimachinery/api/codec"
)

// buildEntry validates one constructor argument and stores normalized metadata.
func buildEntry(index int, c codec.BaseCodec) (Entry, error) {
	if isNilCodec(c) {
		return Entry{}, errorAt(
			codecPath(index),
			ErrInvalidCodec,
			ErrorReasonInvalidCodec,
			"codec must be non-nil",
		)
	}

	info, err := c.Info().Normalize()
	if err != nil {
		return Entry{}, wrapAt(
			infoPath(index),
			ErrInvalidInfo,
			ErrorReasonInvalidInfo,
			"codec info is invalid",
			err,
		)
	}
	if err := validateCapabilities(capabilityPath(index), c, info); err != nil {
		return Entry{}, err
	}

	return Entry{codec: c, info: cloneInfo(info)}, nil
}

// buildRegistry sorts entries and builds immutable indexes against sorted positions.
func buildRegistry(entries []Entry) Registry {
	sorted := slices.Clone(entries)
	slices.SortFunc(sorted, compareEntries)

	registry := Registry{
		entries:     sorted,
		byMediaType: make(map[codec.MediaType]int, mediaTypeCount(sorted)),
		byFormat:    make(map[codec.Format][]int, len(sorted)),
	}
	for i, entry := range sorted {
		for _, mediaType := range entry.info.MediaTypes {
			registry.byMediaType[mediaType] = i
		}
		registry.byFormat[entry.info.Format] = append(registry.byFormat[entry.info.Format], i)
	}

	return registry
}

// mediaTypeCount returns the total declared media type count for index sizing.
func mediaTypeCount(entries []Entry) int {
	total := 0
	for _, entry := range entries {
		total += len(entry.info.MediaTypes)
	}

	return total
}

// compareEntries orders entries deterministically by format and first media type.
func compareEntries(a Entry, b Entry) int {
	if cmp := compareText(a.info.Format.String(), b.info.Format.String()); cmp != 0 {
		return cmp
	}

	return compareText(firstMediaType(a.info).String(), firstMediaType(b.info).String())
}

// firstMediaType returns the first normalized media type when present.
func firstMediaType(info codec.Info) codec.MediaType {
	if len(info.MediaTypes) == 0 {
		return ""
	}

	return info.MediaTypes[0]
}

// compareText returns the standard three-way comparison for stable sorting.
func compareText(a string, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
