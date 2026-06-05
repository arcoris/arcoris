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

// ValueCandidate is one byte-based value codec candidate.
type ValueCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed value codec capability for entry.
	codec codec.ValueCodec
}

// Entry returns the candidate registry entry.
func (c ValueCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed value codec capability.
func (c ValueCandidate) Codec() codec.ValueCodec {
	return c.codec
}

// ObjectCandidate is one byte-based object codec candidate.
type ObjectCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed object codec capability for entry.
	codec codec.ObjectCodec
}

// Entry returns the candidate registry entry.
func (c ObjectCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed object codec capability.
func (c ObjectCandidate) Codec() codec.ObjectCodec {
	return c.codec
}

// ObjectOwnershipCandidate is one byte-based ownership codec candidate.
type ObjectOwnershipCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed ownership codec capability for entry.
	codec codec.ObjectOwnershipCodec
}

// Entry returns the candidate registry entry.
func (c ObjectOwnershipCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed ownership codec capability.
func (c ObjectOwnershipCandidate) Codec() codec.ObjectOwnershipCodec {
	return c.codec
}

// FullCandidate is one full byte-based codec candidate.
type FullCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed full byte codec capability for entry.
	codec codec.Codec
}

// Entry returns the candidate registry entry.
func (c FullCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed full byte codec capability.
func (c FullCandidate) Codec() codec.Codec {
	return c.codec
}

// ValueCandidates returns byte-based value codec candidates for mediaType.
func (r Registry) ValueCandidates(mediaType codec.MediaType) []ValueCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ValueCandidate, 0, len(entries))
	for _, entry := range entries {
		valueCodec, ok := entry.codec.(codec.ValueCodec)
		if ok {
			out = append(out, ValueCandidate{entry: entry, codec: valueCodec})
		}
	}

	return out
}

// ObjectCandidates returns byte-based object codec candidates for mediaType.
func (r Registry) ObjectCandidates(mediaType codec.MediaType) []ObjectCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ObjectCandidate, 0, len(entries))
	for _, entry := range entries {
		objectCodec, ok := entry.codec.(codec.ObjectCodec)
		if ok {
			out = append(out, ObjectCandidate{entry: entry, codec: objectCodec})
		}
	}

	return out
}

// ObjectOwnershipCandidates returns byte-based ownership codec candidates for mediaType.
func (r Registry) ObjectOwnershipCandidates(mediaType codec.MediaType) []ObjectOwnershipCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ObjectOwnershipCandidate, 0, len(entries))
	for _, entry := range entries {
		ownershipCodec, ok := entry.codec.(codec.ObjectOwnershipCodec)
		if ok {
			out = append(out, ObjectOwnershipCandidate{entry: entry, codec: ownershipCodec})
		}
	}

	return out
}

// FullCandidates returns full byte-based codec candidates for mediaType.
func (r Registry) FullCandidates(mediaType codec.MediaType) []FullCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]FullCandidate, 0, len(entries))
	for _, entry := range entries {
		fullCodec, ok := entry.codec.(codec.Codec)
		if ok {
			out = append(out, FullCandidate{entry: entry, codec: fullCodec})
		}
	}

	return out
}
