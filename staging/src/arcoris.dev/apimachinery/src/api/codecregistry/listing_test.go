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
	registry := testRegistry(
		t,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	if registry.Len() != 2 {
		t.Fatalf("Len() = %d", registry.Len())
	}
}

func TestRegistryIsEmpty(t *testing.T) {
	var zero Registry
	if !zero.IsEmpty() {
		t.Fatalf("zero IsEmpty() = false")
	}

	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))
	if registry.IsEmpty() {
		t.Fatalf("non-empty IsEmpty() = true")
	}
}

func TestRegistryEntriesDeterministic(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	entries := registry.Entries()
	if entries[0].ID() != MustEntryID("json.public") || entries[1].ID() != MustEntryID("yaml.public") {
		t.Fatalf("entry order = %q, %q", entries[0].ID(), entries[1].ID())
	}
}

func TestRegistryEntriesDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

	entries := registry.Entries()
	entries[0] = Entry{}

	if registry.Entries()[0].IsZero() {
		t.Fatalf("registry entry mutated through Entries()")
	}
}

func TestRegistryCodecsDeterministic(t *testing.T) {
	jsonCodec := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	yamlCodec := newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML)
	registry := testRegistry(
		t,
		testRegistration("yaml.public", yamlCodec),
		testRegistration("json.public", jsonCodec),
	)

	codecs := registry.Codecs()
	if codecs[0] != jsonCodec || codecs[1] != yamlCodec {
		t.Fatalf("codec order = %#v", codecs)
	}
}

func TestRegistryCodecsDetached(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registry := testRegistry(t, testRegistration("json.public", c))

	codecs := registry.Codecs()
	codecs[0] = nil

	if registry.Codecs()[0] != c {
		t.Fatalf("registry codecs mutated through Codecs()")
	}
}

func TestIDsSorted(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	ids := registry.IDs()
	if ids[0] != MustEntryID("json.public") || ids[1] != MustEntryID("yaml.public") {
		t.Fatalf("IDs() = %#v", ids)
	}
}

func TestIDsDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

	ids := registry.IDs()
	ids[0] = MustEntryID("json.other")

	if registry.IDs()[0] != MustEntryID("json.public") {
		t.Fatalf("registry IDs mutated through IDs()")
	}
}

func TestRegistryFormatsSorted(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	formats := registry.Formats()
	if formats[0] != codec.FormatJSON || formats[1] != codec.FormatYAML {
		t.Fatalf("Formats() = %#v", formats)
	}
}

func TestFormatsDeduplicatesAndSorts(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.z", codec.FormatJSON, "application/vnd.z+json"),
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	formats := registry.Formats()
	if len(formats) != 2 {
		t.Fatalf("Formats() length = %d; want 2", len(formats))
	}
	if formats[0] != codec.FormatJSON || formats[1] != codec.FormatYAML {
		t.Fatalf("Formats() = %#v", formats)
	}
}

func TestRegistryFormatsDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

	formats := registry.Formats()
	formats[0] = codec.FormatYAML

	if registry.Formats()[0] != codec.FormatJSON {
		t.Fatalf("registry formats mutated through Formats()")
	}
}

func TestRegistryMediaTypesSorted(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)

	mediaTypes := registry.MediaTypes()
	if mediaTypes[0] != codec.MediaTypeJSON || mediaTypes[1] != codec.MediaTypeYAML {
		t.Fatalf("MediaTypes() = %#v", mediaTypes)
	}
}

func TestMediaTypesDeduplicatesAndSorts(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
	)

	mediaTypes := registry.MediaTypes()
	if len(mediaTypes) != 2 {
		t.Fatalf("MediaTypes() length = %d; want 2", len(mediaTypes))
	}
	if mediaTypes[0] != codec.MediaTypeJSON || mediaTypes[1] != codec.MediaTypeYAML {
		t.Fatalf("MediaTypes() = %#v", mediaTypes)
	}
}

func TestRegistryMediaTypesDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

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
		len(registry.IDs()) != 0 ||
		len(registry.Formats()) != 0 ||
		len(registry.MediaTypes()) != 0 {
		t.Fatalf("zero registry lists are not empty")
	}
}
