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

func TestEntryZero(t *testing.T) {
	var entry Entry

	if !entry.IsZero() {
		t.Fatalf("zero entry IsZero() = false")
	}
	if entry.Codec() != nil {
		t.Fatalf("zero entry Codec() = %v", entry.Codec())
	}
	if !entry.ID().IsZero() {
		t.Fatalf("zero entry ID() = %q", entry.ID())
	}
	if !entry.Info().IsZero() {
		t.Fatalf("zero entry Info() = %#v", entry.Info())
	}
	if len(entry.MediaTypes()) != 0 || len(entry.Targets()) != 0 {
		t.Fatalf("zero entry slices are non-empty")
	}
}

func TestEntryInfoDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))
	entry := registry.Entries()[0]

	info := entry.Info()
	info.MediaTypes[0] = codec.MediaTypeYAML
	info.Targets[0] = codec.TargetObject

	got := entry.Info()
	if got.MediaTypes[0] != codec.MediaTypeJSON {
		t.Fatalf("entry media type mutated: %q", got.MediaTypes[0])
	}
	if got.Targets[0] != codec.TargetValue {
		t.Fatalf("entry target mutated: %q", got.Targets[0])
	}
}

func TestEntryID(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

	if got := registry.Entries()[0].ID(); got != MustEntryID("json.public") {
		t.Fatalf("ID() = %q", got)
	}
}

func TestEntryFormat(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))

	if got := registry.Entries()[0].Format(); got != codec.FormatJSON {
		t.Fatalf("Format() = %q", got)
	}
}

func TestEntryMediaTypesDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))
	entry := registry.Entries()[0]

	mediaTypes := entry.MediaTypes()
	mediaTypes[0] = codec.MediaTypeYAML

	if got := entry.MediaTypes()[0]; got != codec.MediaTypeJSON {
		t.Fatalf("MediaTypes()[0] = %q", got)
	}
}

func TestEntryTargetsDetached(t *testing.T) {
	registry := testRegistry(t, testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON))
	entry := registry.Entries()[0]

	targets := entry.Targets()
	targets[0] = codec.TargetObject

	if got := entry.Targets()[0]; got != codec.TargetValue {
		t.Fatalf("Targets()[0] = %q", got)
	}
}

func TestEntryCodec(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registry := testRegistry(t, testRegistration("json.public", c))

	if got := registry.Entries()[0].Codec(); got != c {
		t.Fatalf("Codec() = %v; want original codec", got)
	}
}
