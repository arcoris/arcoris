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

func TestValueCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueCandidates() length = %d; want 1", len(candidates))
	}
}

func TestObjectCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ObjectCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ObjectCandidates() length = %d; want 1", len(candidates))
	}
}

func TestObjectOwnershipCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testOwnershipByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ObjectOwnershipCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ObjectOwnershipCandidates() length = %d; want 1", len(candidates))
	}
}

func TestFullCandidates(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullCandidates() length = %d; want 1", len(candidates))
	}
}

func TestCandidatesReturnMultipleForSameMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	candidates := registry.ValueCandidates(codec.MediaTypeJSON)
	if len(candidates) != 2 {
		t.Fatalf("ValueCandidates() length = %d; want 2", len(candidates))
	}
	if candidates[0].Entry().ID() != MustEntryID("json.public") ||
		candidates[1].Entry().ID() != MustEntryID("json.storage") {
		t.Fatalf("candidate IDs = %q, %q", candidates[0].Entry().ID(), candidates[1].Entry().ID())
	}
}

func TestCandidatesFilterByCapability(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.object", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.value", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueCandidates() length = %d; want 1", len(candidates))
	}
	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("FullCandidates() length = %d; want 0", len(candidates))
	}
}

func TestCandidatesReturnEmptyForInvalidMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates("application/json; charset=utf-8"); len(candidates) != 0 {
		t.Fatalf("ValueCandidates() length = %d; want 0", len(candidates))
	}
}

func TestCandidatesReturnEmptyForMissingMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates(codec.MediaTypeYAML); len(candidates) != 0 {
		t.Fatalf("ValueCandidates() length = %d; want 0", len(candidates))
	}
}

func TestCandidateEntryAndCodecAccessors(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registry := testRegistry(t, testRegistration("json.public", c))

	candidates := registry.ValueCandidates(codec.MediaTypeJSON)
	if candidates[0].Entry().ID() != MustEntryID("json.public") {
		t.Fatalf("Entry().ID() = %q", candidates[0].Entry().ID())
	}
	if candidates[0].Codec() != c {
		t.Fatalf("Codec() = %v; want original codec", candidates[0].Codec())
	}
}

func TestCandidateSlicesDetached(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	candidates := registry.ValueCandidates(codec.MediaTypeJSON)
	candidates[0] = ValueCandidate{}

	if registry.ValueCandidates(codec.MediaTypeJSON)[0].Entry().IsZero() {
		t.Fatalf("candidate slice mutation changed registry result")
	}
}

func TestZeroRegistryByteCandidatesReturnEmpty(t *testing.T) {
	var registry Registry

	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ValueCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.ObjectCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ObjectCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.ObjectOwnershipCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("ObjectOwnershipCandidates() length = %d; want 0", len(candidates))
	}
	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 0 {
		t.Fatalf("FullCandidates() length = %d; want 0", len(candidates))
	}
}
