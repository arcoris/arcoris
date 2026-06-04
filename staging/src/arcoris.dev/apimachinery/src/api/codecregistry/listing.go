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

// IsEmpty reports whether r contains no registered codecs.
func (r Registry) IsEmpty() bool {
	return len(r.entries) == 0
}

// Len returns the number of registered codec implementations.
func (r Registry) Len() int {
	return len(r.entries)
}

// Entries returns detached registry entries in deterministic registry order.
func (r Registry) Entries() []Entry {
	return slices.Clone(r.entries)
}

// Codecs returns codec implementations in deterministic registry order.
//
// The returned slice is detached, but each element is the original registered
// codec implementation.
func (r Registry) Codecs() []codec.BaseCodec {
	out := make([]codec.BaseCodec, 0, len(r.entries))
	for _, entry := range r.entries {
		out = append(out, entry.codec)
	}

	return out
}

// Formats returns detached unique formats sorted by canonical format string.
func (r Registry) Formats() []codec.Format {
	out := make([]codec.Format, 0, len(r.byFormat))
	for format := range r.byFormat {
		out = append(out, format)
	}
	slices.Sort(out)

	return out
}

// MediaTypes returns detached media types sorted by canonical media type string.
func (r Registry) MediaTypes() []codec.MediaType {
	out := make([]codec.MediaType, 0, len(r.byMediaType))
	for mediaType := range r.byMediaType {
		out = append(out, mediaType)
	}
	slices.Sort(out)

	return out
}
