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

func TestEntriesByFormat(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeYAML),
	)
	requireNoError(t, err)

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

func TestEntriesByFormatNormalizesInputIfSupported(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if entries := registry.EntriesByFormat(" JSON "); len(entries) != 1 {
		t.Fatalf("EntriesByFormat() length = %d; want 1", len(entries))
	}
}

func TestLookupMediaType(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	entry, ok := registry.LookupMediaType(codec.MediaTypeJSON)
	if !ok {
		t.Fatalf("LookupMediaType() = false")
	}
	if entry.Format() != codec.FormatJSON {
		t.Fatalf("LookupMediaType() format = %q", entry.Format())
	}
}

func TestLookupMediaTypeNormalizesInputIfSupported(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupMediaType(" Application/JSON "); !ok {
		t.Fatalf("LookupMediaType() did not normalize input")
	}
}

func TestEntriesByFormatReturnsEmptyForMissingFormat(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if entries := registry.EntriesByFormat(codec.FormatYAML); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
}

func TestLookupMissingMediaTypeReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupMediaType(codec.MediaTypeYAML); ok {
		t.Fatalf("LookupMediaType() = true for missing media type")
	}
}

func TestEntriesByFormatReturnsEmptyForInvalidFormat(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if entries := registry.EntriesByFormat("application/json"); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
}

func TestEntriesByFormatReturnsDetachedSlice(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	entries[0] = Entry{}

	if registry.EntriesByFormat(codec.FormatJSON)[0].IsZero() {
		t.Fatalf("registry entries mutated through EntriesByFormat()")
	}
}

func TestEntriesByFormatPreservesRegistryOrder(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatJSON, "application/vnd.z+json"),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	entries := registry.EntriesByFormat(codec.FormatJSON)
	if got := entries[0].MediaTypes()[0]; got != codec.MediaTypeJSON {
		t.Fatalf("entries[0] media type = %q; want %q", got, codec.MediaTypeJSON)
	}
	if got := entries[1].MediaTypes()[0]; got != "application/vnd.z+json" {
		t.Fatalf("entries[1] media type = %q; want application/vnd.z+json", got)
	}
}

func TestLookupInvalidMediaTypeReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupMediaType("application/json; charset=utf-8"); ok {
		t.Fatalf("LookupMediaType() = true for invalid media type")
	}
}

func TestZeroRegistryLookupReturnsFalse(t *testing.T) {
	var registry Registry

	if entries := registry.EntriesByFormat(codec.FormatJSON); len(entries) != 0 {
		t.Fatalf("EntriesByFormat() length = %d; want 0", len(entries))
	}
	if _, ok := registry.LookupMediaType(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupMediaType() = true on zero registry")
	}
}
