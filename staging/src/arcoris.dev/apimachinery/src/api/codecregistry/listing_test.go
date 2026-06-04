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

func TestRegistryLen(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	if registry.Len() != 2 {
		t.Fatalf("Len() = %d", registry.Len())
	}
}

func TestRegistryIsEmpty(t *testing.T) {
	var zero Registry
	if !zero.IsEmpty() {
		t.Fatalf("zero IsEmpty() = false")
	}

	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)
	if registry.IsEmpty() {
		t.Fatalf("non-empty IsEmpty() = true")
	}
}

func TestRegistryEntriesDeterministic(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	entries := registry.Entries()
	if entries[0].Format() != codec.FormatJSON || entries[1].Format() != codec.FormatYAML {
		t.Fatalf("entry order = %q, %q", entries[0].Format(), entries[1].Format())
	}
}

func TestRegistryEntriesDetached(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	entries := registry.Entries()
	entries[0] = Entry{}

	if registry.Entries()[0].IsZero() {
		t.Fatalf("registry entry mutated through Entries()")
	}
}

func TestRegistryCodecsDeterministic(t *testing.T) {
	jsonCodec := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	yamlCodec := newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML)
	registry, err := New(yamlCodec, jsonCodec)
	requireNoError(t, err)

	codecs := registry.Codecs()
	if codecs[0] != jsonCodec || codecs[1] != yamlCodec {
		t.Fatalf("codec order = %#v", codecs)
	}
}

func TestRegistryCodecsDetached(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registry, err := New(c)
	requireNoError(t, err)

	codecs := registry.Codecs()
	codecs[0] = nil

	if registry.Codecs()[0] != c {
		t.Fatalf("registry codecs mutated through Codecs()")
	}
}

func TestRegistryFormatsSorted(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	formats := registry.Formats()
	if formats[0] != codec.FormatJSON || formats[1] != codec.FormatYAML {
		t.Fatalf("Formats() = %#v", formats)
	}
}

func TestRegistryFormatsDetached(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	formats := registry.Formats()
	formats[0] = codec.FormatYAML

	if registry.Formats()[0] != codec.FormatJSON {
		t.Fatalf("registry formats mutated through Formats()")
	}
}

func TestRegistryMediaTypesSorted(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	mediaTypes := registry.MediaTypes()
	if mediaTypes[0] != codec.MediaTypeJSON || mediaTypes[1] != codec.MediaTypeYAML {
		t.Fatalf("MediaTypes() = %#v", mediaTypes)
	}
}

func TestRegistryMediaTypesDetached(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	mediaTypes := registry.MediaTypes()
	mediaTypes[0] = codec.MediaTypeYAML

	if registry.MediaTypes()[0] != codec.MediaTypeJSON {
		t.Fatalf("registry media types mutated through MediaTypes()")
	}
}

func TestZeroRegistryListsAreEmpty(t *testing.T) {
	var registry Registry

	if len(registry.Entries()) != 0 ||
		len(registry.Codecs()) != 0 ||
		len(registry.Formats()) != 0 ||
		len(registry.MediaTypes()) != 0 {
		t.Fatalf("zero registry lists are not empty")
	}
}
