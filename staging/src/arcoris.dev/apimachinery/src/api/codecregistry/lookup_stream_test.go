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

func TestLookupValueStream(t *testing.T) {
	registry, err := New(newValueStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValueStream() = false")
	}
}

func TestLookupObjectStream(t *testing.T) {
	registry, err := New(newObjectStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupObjectStream(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupObjectStream() = false")
	}
}

func TestLookupObjectOwnershipStream(t *testing.T) {
	registry, err := New(newOwnershipStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupObjectOwnershipStream(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupObjectOwnershipStream() = false")
	}
}

func TestLookupStreamingCodec(t *testing.T) {
	registry, err := New(newFullStreamingCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupStreamingCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupStreamingCodec() = false")
	}
}

func TestLookupWrongStreamingCapabilityReturnsFalse(t *testing.T) {
	registry, err := New(newObjectStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupValueStream() = true for object-stream-only codec")
	}
	if _, ok := registry.LookupStreamingCodec(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupStreamingCodec() = true for partial stream codec")
	}
}

func TestStreamingLookupDoesNotRequireByteCodec(t *testing.T) {
	registry, err := New(newValueStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValueStream() = false for stream-only codec")
	}
	if _, ok := registry.LookupValue(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupValue() = true for stream-only codec")
	}
}

func TestByteLookupDoesNotRequireStreamingCodec(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValue(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValue() = false for byte-only codec")
	}
	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupValueStream() = true for byte-only codec")
	}
}

func TestZeroRegistryStreamLookupsReturnFalse(t *testing.T) {
	var registry Registry

	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupValueStream() = true on zero registry")
	}
	if _, ok := registry.LookupObjectStream(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupObjectStream() = true on zero registry")
	}
	if _, ok := registry.LookupObjectOwnershipStream(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupObjectOwnershipStream() = true on zero registry")
	}
	if _, ok := registry.LookupStreamingCodec(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupStreamingCodec() = true on zero registry")
	}
}
