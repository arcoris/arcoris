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

func TestNewEmpty(t *testing.T) {
	registry, err := New()
	requireNoError(t, err)

	if !registry.IsEmpty() || registry.Len() != 0 {
		t.Fatalf("registry empty = %v len = %d", registry.IsEmpty(), registry.Len())
	}
}

func TestZeroRegistryIsUsable(t *testing.T) {
	var registry Registry

	if !registry.IsEmpty() || registry.Len() != 0 {
		t.Fatalf("zero registry empty = %v len = %d", registry.IsEmpty(), registry.Len())
	}
	if _, ok := registry.LookupID(MustEntryID("json.public")); ok {
		t.Fatalf("zero registry ID lookup returned true")
	}
	if entries := registry.EntriesByMediaType(codec.MediaTypeJSON); len(entries) != 0 {
		t.Fatalf("zero registry media type entries length = %d; want 0", len(entries))
	}
}

func TestNewRejectsZeroRegistration(t *testing.T) {
	_, err := New(Registration{})

	requireErrorIs(t, err, ErrInvalidRegistration)
	requireRegistryError(t, err, "registrations[0]", ErrorReasonInvalidRegistration)
}

func TestNewRejectsInvalidEntryID(t *testing.T) {
	_, err := New(Register("JSON.Public", newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)))

	requireErrorIs(t, err, ErrInvalidEntryID)
	requireRegistryError(t, err, "registrations[0].id", ErrorReasonInvalidEntryID)
}

func TestNewRejectsNilCodec(t *testing.T) {
	_, err := New(Register(MustEntryID("json.public"), nil))

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestNewRejectsTypedNilCodec(t *testing.T) {
	var c *fakeValueByteCodec

	_, err := New(testRegistration("json.public", c))

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestNewRejectsInvalidInfo(t *testing.T) {
	c := fakeBaseCodec{info: codec.Info{}}

	_, err := New(testRegistration("json.public", c))

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
	requireRegistryError(t, err, "registrations[0].info", ErrorReasonInvalidInfo)
}

func TestNewReturnsZeroRegistryOnError(t *testing.T) {
	registry, err := New(
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.public", codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateEntryID)
	if !registry.IsEmpty() {
		t.Fatalf("registry = %#v; want zero registry on error", registry)
	}
}

func TestNewCallsInfoOncePerCodec(t *testing.T) {
	c := &fakeCountingValueCodec{
		fakeValueByteCodec: newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	}

	_, err := New(testRegistration("json.public", c))
	requireNoError(t, err)

	if c.calls != 1 {
		t.Fatalf("Info calls = %d; want 1", c.calls)
	}
}

func TestNewAcceptsByteOnlyValueCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueCandidates() length = %d; want 1", len(candidates))
	}
}

func TestNewAcceptsStreamOnlyValueCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testValueStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.ValueStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("ValueStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestNewAcceptsFullByteCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullCandidates() length = %d; want 1", len(candidates))
	}
}

func TestNewAcceptsFullStreamingCodec(t *testing.T) {
	registry := testRegistry(
		t,
		testFullStreamRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullStreamCandidates() length = %d; want 1", len(candidates))
	}
}

func TestNewAcceptsCodecImplementingBothByteAndStream(t *testing.T) {
	registry := testRegistry(
		t,
		testRegistration("json.public", newByteAndStreamCodec(codec.FormatJSON, codec.MediaTypeJSON)),
	)

	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullCandidates() length = %d; want 1", len(candidates))
	}
	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullStreamCandidates() length = %d; want 1", len(candidates))
	}
}
