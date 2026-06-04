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

func TestLookupFormat(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	entry, ok := registry.LookupFormat(codec.FormatJSON)
	if !ok {
		t.Fatalf("LookupFormat() = false")
	}
	if entry.Format() != codec.FormatJSON {
		t.Fatalf("LookupFormat() format = %q", entry.Format())
	}
}

func TestLookupFormatNormalizesInputIfSupported(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupFormat(" JSON "); !ok {
		t.Fatalf("LookupFormat() did not normalize input")
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

func TestLookupMissingFormatReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupFormat(codec.FormatYAML); ok {
		t.Fatalf("LookupFormat() = true for missing format")
	}
}

func TestLookupMissingMediaTypeReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupMediaType(codec.MediaTypeYAML); ok {
		t.Fatalf("LookupMediaType() = true for missing media type")
	}
}

func TestLookupInvalidFormatReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupFormat("application/json"); ok {
		t.Fatalf("LookupFormat() = true for invalid format")
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

	if _, ok := registry.LookupFormat(codec.FormatJSON); ok {
		t.Fatalf("LookupFormat() = true on zero registry")
	}
	if _, ok := registry.LookupMediaType(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupMediaType() = true on zero registry")
	}
}
