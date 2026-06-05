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

// ValueStreamCandidate is one streaming value codec candidate.
type ValueStreamCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed value stream codec capability for entry.
	codec codec.ValueStreamCodec
}

// Entry returns the candidate registry entry.
func (c ValueStreamCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed value stream codec capability.
func (c ValueStreamCandidate) Codec() codec.ValueStreamCodec {
	return c.codec
}

// ObjectStreamCandidate is one streaming object codec candidate.
type ObjectStreamCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed object stream codec capability for entry.
	codec codec.ObjectStreamCodec
}

// Entry returns the candidate registry entry.
func (c ObjectStreamCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed object stream codec capability.
func (c ObjectStreamCandidate) Codec() codec.ObjectStreamCodec {
	return c.codec
}

// ObjectOwnershipStreamCandidate is one streaming ownership codec candidate.
type ObjectOwnershipStreamCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed ownership stream codec capability for entry.
	codec codec.ObjectOwnershipStreamCodec
}

// Entry returns the candidate registry entry.
func (c ObjectOwnershipStreamCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed ownership stream codec capability.
func (c ObjectOwnershipStreamCandidate) Codec() codec.ObjectOwnershipStreamCodec {
	return c.codec
}

// FullStreamCandidate is one full streaming codec candidate.
type FullStreamCandidate struct {
	// entry is the registry metadata snapshot for this configured candidate.
	entry Entry

	// codec is the typed full streaming codec capability for entry.
	codec codec.StreamingCodec
}

// Entry returns the candidate registry entry.
func (c FullStreamCandidate) Entry() Entry {
	return c.entry
}

// Codec returns the typed full streaming codec capability.
func (c FullStreamCandidate) Codec() codec.StreamingCodec {
	return c.codec
}

// ValueStreamCandidates returns streaming value codec candidates for mediaType.
func (r Registry) ValueStreamCandidates(mediaType codec.MediaType) []ValueStreamCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ValueStreamCandidate, 0, len(entries))
	for _, entry := range entries {
		valueCodec, ok := entry.codec.(codec.ValueStreamCodec)
		if ok {
			out = append(out, ValueStreamCandidate{entry: entry, codec: valueCodec})
		}
	}

	return out
}

// ObjectStreamCandidates returns streaming object codec candidates for mediaType.
func (r Registry) ObjectStreamCandidates(mediaType codec.MediaType) []ObjectStreamCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ObjectStreamCandidate, 0, len(entries))
	for _, entry := range entries {
		objectCodec, ok := entry.codec.(codec.ObjectStreamCodec)
		if ok {
			out = append(out, ObjectStreamCandidate{entry: entry, codec: objectCodec})
		}
	}

	return out
}

// ObjectOwnershipStreamCandidates returns streaming ownership codec candidates for mediaType.
func (r Registry) ObjectOwnershipStreamCandidates(
	mediaType codec.MediaType,
) []ObjectOwnershipStreamCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]ObjectOwnershipStreamCandidate, 0, len(entries))
	for _, entry := range entries {
		ownershipCodec, ok := entry.codec.(codec.ObjectOwnershipStreamCodec)
		if ok {
			out = append(out, ObjectOwnershipStreamCandidate{entry: entry, codec: ownershipCodec})
		}
	}

	return out
}

// FullStreamCandidates returns full streaming codec candidates for mediaType.
func (r Registry) FullStreamCandidates(mediaType codec.MediaType) []FullStreamCandidate {
	entries := r.EntriesByMediaType(mediaType)
	if len(entries) == 0 {
		return nil
	}

	out := make([]FullStreamCandidate, 0, len(entries))
	for _, entry := range entries {
		streamingCodec, ok := entry.codec.(codec.StreamingCodec)
		if ok {
			out = append(out, FullStreamCandidate{entry: entry, codec: streamingCodec})
		}
	}

	return out
}
