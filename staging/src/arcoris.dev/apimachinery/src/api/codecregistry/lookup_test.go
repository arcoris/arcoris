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

func TestLookupID(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entry, ok := registry.LookupID(MustEntryID("json.public"))
	if !ok {
		t.Fatalf("LookupID() = false")
	}
	if entry.ID() != MustEntryID("json.public") {
		t.Fatalf("LookupID() ID = %q", entry.ID())
	}
}

func TestLookupIDRejectsInvalidInput(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if _, ok := registry.LookupID("JSON.Public"); ok {
		t.Fatalf("LookupID() = true for invalid ID")
	}
}

func TestLookupIDMissingReturnsFalse(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if _, ok := registry.LookupID(MustEntryID("json.missing")); ok {
		t.Fatalf("LookupID() = true for missing ID")
	}
}

func TestEntriesByMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByMediaType(codec.MediaTypeJSON)
	if len(entries) != 1 {
		t.Fatalf("EntriesByMediaType() length = %d; want 1", len(entries))
	}
	if entries[0].ID() != MustEntryID("json.public") {
		t.Fatalf("EntriesByMediaType()[0].ID() = %q", entries[0].ID())
	}
}

func TestEntriesByMediaTypeAllowsMultipleEntries(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByMediaType(codec.MediaTypeJSON)
	if len(entries) != 2 {
		t.Fatalf("EntriesByMediaType() length = %d; want 2", len(entries))
	}
	if entries[0].ID() != MustEntryID("json.public") || entries[1].ID() != MustEntryID("json.storage") {
		t.Fatalf("entry IDs = %q, %q", entries[0].ID(), entries[1].ID())
	}
}

func TestEntriesByMediaTypeNormalizesInputIfSupported(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByMediaType(" Application/JSON "); len(entries) != 1 {
		t.Fatalf("EntriesByMediaType() length = %d; want 1", len(entries))
	}
}

func TestEntriesByMediaTypeReturnsEmptyForMissingMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByMediaType(codec.MediaTypeYAML); len(entries) != 0 {
		t.Fatalf("EntriesByMediaType() length = %d; want 0", len(entries))
	}
}

func TestEntriesByMediaTypeReturnsEmptyForInvalidMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByMediaType("application/json; charset=utf-8"); len(entries) != 0 {
		t.Fatalf("EntriesByMediaType() length = %d; want 0", len(entries))
	}
}

func TestEntriesByMediaTypeReturnsDetachedSlice(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByMediaType(codec.MediaTypeJSON)
	entries[0] = Entry{}

	if registry.EntriesByMediaType(codec.MediaTypeJSON)[0].IsZero() {
		t.Fatalf("registry entries mutated through EntriesByMediaType()")
	}
}

func TestEntriesByMediaTypePreservesRegistryOrder(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.z", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByMediaType(codec.MediaTypeJSON)
	if got := entries[0].ID(); got != MustEntryID("json.public") {
		t.Fatalf("entries[0].ID() = %q; want json.public", got)
	}
	if got := entries[1].ID(); got != MustEntryID("json.z") {
		t.Fatalf("entries[1].ID() = %q; want json.z", got)
	}
}

func TestEntriesByFormat(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeYAML),
	)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	if len(entries) != 2 {
		t.Fatalf("EntriesByFormat() length = %d; want 2", len(entries))
	}
	for i, entry := range entries {
		if entry.Format() != codec.FormatJSON {
			t.Fatalf("entries[%d].Format() = %q", i, entry.Format())
		}
	}
}

func TestEntriesByFormatAllowsMultipleEntries(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByFormat(codec.FormatJSON); len(entries) != 2 {
		t.Fatalf("EntriesByFormat() length = %d; want 2", len(entries))
	}
}

func TestEntriesByFormatNormalizesInputIfSupported(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByFormat(" JSON "); len(entries) != 1 {
		t.Fatalf("EntriesByFormat() length = %d; want 1", len(entries))
	}
}

func TestEntriesByFormatReturnsEmptyForMissingFormat(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByFormat(codec.FormatYAML); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
}

func TestEntriesByFormatReturnsEmptyForInvalidFormat(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if entries := registry.EntriesByFormat("application/json"); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
}

func TestEntriesByFormatReturnsDetachedSlice(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	entries[0] = Entry{}

	if registry.EntriesByFormat(codec.FormatJSON)[0].IsZero() {
		t.Fatalf("registry entries mutated through EntriesByFormat()")
	}
}

func TestEntriesByFormatPreservesRegistryOrder(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.z", codec.FormatJSON, "application/vnd.z+json"),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	if got := entries[0].ID(); got != MustEntryID("json.public") {
		t.Fatalf("entries[0].ID() = %q; want json.public", got)
	}
	if got := entries[1].ID(); got != MustEntryID("json.z") {
		t.Fatalf("entries[1].ID() = %q; want json.z", got)
	}
}

func TestZeroRegistryLookupsReturnEmpty(t *testing.T) {
	var registry Registry

	if _, ok := registry.LookupID(MustEntryID("json.public")); ok {
		t.Fatalf("LookupID() = true on zero registry")
	}
	if entries := registry.EntriesByMediaType(codec.MediaTypeJSON); len(entries) != 0 {
		t.Fatalf("EntriesByMediaType() length = %d; want 0", len(entries))
	}
	if entries := registry.EntriesByFormat(codec.FormatJSON); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
}
