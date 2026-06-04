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

func TestLookupValue(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValue(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValue() = false")
	}
}

func TestLookupObject(t *testing.T) {
	registry, err := New(newObjectByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupObject(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupObject() = false")
	}
}

func TestLookupObjectOwnership(t *testing.T) {
	registry, err := New(newOwnershipByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupObjectOwnership(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupObjectOwnership() = false")
	}
}

func TestLookupCodec(t *testing.T) {
	registry, err := New(newFullByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupCodec() = false")
	}
}

func TestLookupWrongByteCapabilityReturnsFalse(t *testing.T) {
	registry, err := New(newObjectByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValue(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupValue() = true for object-only codec")
	}
	if _, ok := registry.LookupCodec(codec.MediaTypeJSON); ok {
		t.Fatalf("LookupCodec() = true for partial codec")
	}
}

func TestLookupMissingByteMediaTypeReturnsFalse(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValue(codec.MediaTypeYAML); ok {
		t.Fatalf("LookupValue() = true for missing media type")
	}
}
