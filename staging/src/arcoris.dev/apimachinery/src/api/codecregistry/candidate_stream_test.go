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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestValueStreamCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testValueStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestObjectStreamCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ObjectStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ObjectStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestObjectOwnershipStreamCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testOwnershipStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ObjectOwnershipStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ObjectOwnershipStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestFullStreamCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testFullStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestStreamCandidatesReturnMultipleForSameMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueStreamRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON)
	if len(candidates) != 2 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 2", len(candidates))
	}
	if candidates[0].Entry().ID() != MustEntryID("json.public") ||
		candidates[1].Entry().ID() != MustEntryID("json.storage") {
		t.Fatalf("candidate IDs = %q, %q", candidates[0].Entry().ID(), candidates[1].Entry().ID())
	}
}

func TestStreamCandidatesFilterByCapability(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectStreamRegistration("json.object", codec.FormatJSON, codec.MediaTypeJSON),
		testValueStreamRegistration("json.value", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 1", len(candidates))
	}
	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("FullStreamCandidates() length = %d; want 0", len(candidates))
	}
}

func TestStreamCandidateEntryAndCodecAccessors(t *testing.T) {
	c := newValueStreamCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registry := testRegistry(t, testRegistration("json.public", c))

	candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON)
	if candidates[0].Entry().ID() != MustEntryID("json.public") {
		t.Fatalf("Entry().ID() = %q", candidates[0].Entry().ID())
	}
	if candidates[0].Codec() != c {
		t.Fatalf("Codec() = %v; want original codec", candidates[0].Codec())
	}
}

func TestStreamingLookupDoesNotRequireByteCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testValueStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 1", len(candidates))
	}
	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ValueCandidates() length = %d; want 0", len(candidates))
	}
}

func TestByteLookupDoesNotRequireStreamingCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueCandidates() length = %d; want 1", len(candidates))
	}
	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 0", len(candidates))
	}
}

func TestZeroRegistryStreamCandidatesReturnEmpty(t *testing.T) {
	var registry Registry

	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.ObjectStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ObjectStreamCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.ObjectOwnershipStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ObjectOwnershipStreamCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("FullStreamCandidates() length = %d; want 0", len(candidates))
	}
}
