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

// EntriesByFormat returns codec entries registered for format.
//
// Format is a grouping attribute, not a unique codec identity. The returned
// slice is detached and preserves deterministic registry entry order.
func (r Registry) EntriesByFormat(format codec.Format) []Entry {
	normalized, ok := normalizeFormat(format)
	if !ok {
		return nil
	}

	indexes := r.byFormat[normalized]
	if len(indexes) == 0 {
		return nil
	}

	out := make([]Entry, 0, len(indexes))
	for _, index := range indexes {
		if entry, ok := r.entryAt(index, true); ok {
			out = append(out, entry)
		}
	}

	return out
}

// LookupMediaType returns the codec entry registered for mediaType.
func (r Registry) LookupMediaType(mediaType codec.MediaType) (Entry, bool) {
	normalized, ok := normalizeMediaType(mediaType)
	if !ok {
		return Entry{}, false
	}

	index, ok := r.byMediaType[normalized]
	return r.entryAt(index, ok)
}

// entryAt returns the indexed entry when ok and index are valid.
func (r Registry) entryAt(index int, ok bool) (Entry, bool) {
	if !ok || index < 0 || index >= len(r.entries) {
		return Entry{}, false
	}

	return r.entries[index], true
}

// normalizeFormat canonicalizes lookup input and treats invalid input as absent.
func normalizeFormat(format codec.Format) (codec.Format, bool) {
	normalized, err := format.Normalize()
	if err != nil {
		return "", false
	}

	return normalized, true
}

// normalizeMediaType canonicalizes lookup input and treats invalid input as absent.
func normalizeMediaType(mediaType codec.MediaType) (codec.MediaType, bool) {
	normalized, err := mediaType.Normalize()
	if err != nil {
		return "", false
	}

	return normalized, true
}
