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

func TestBuildEntryRejectsNilCodec(t *testing.T) {
	entry, err := buildEntry(0, Register(MustEntryID("json.public"), nil))

	if !entry.IsZero() {
		t.Fatalf("entry = %#v; want zero", entry)
	}
	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestBuildEntryNormalizesInfo(t *testing.T) {
	c := &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(" JSON ", " Application/JSON ", codec.TargetValue),
	}}

	entry, err := buildEntry(0, testRegistration("json.public", c))
	requireNoError(t, err)

	if entry.ID() != MustEntryID("json.public") {
		t.Fatalf("ID() = %q", entry.ID())
	}
	if entry.Format() != codec.FormatJSON {
		t.Fatalf("Format() = %q", entry.Format())
	}
	if got := entry.MediaTypes()[0]; got != codec.MediaTypeJSON {
		t.Fatalf("MediaTypes()[0] = %q", got)
	}
}

func TestBuildEntryStoresDetachedInfo(t *testing.T) {
	mediaTypes := []codec.MediaType{codec.MediaTypeJSON}
	targets := []codec.Target{codec.TargetValue}
	c := &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: codec.Info{
			Format:     codec.FormatJSON,
			MediaTypes: mediaTypes,
			Targets:    targets,
		},
	}}

	entry, err := buildEntry(0, testRegistration("json.public", c))
	requireNoError(t, err)

	mediaTypes[0] = codec.MediaTypeYAML
	targets[0] = codec.TargetObject

	if got := entry.MediaTypes()[0]; got != codec.MediaTypeJSON {
		t.Fatalf("MediaTypes()[0] = %q", got)
	}
	if got := entry.Targets()[0]; got != codec.TargetValue {
		t.Fatalf("Targets()[0] = %q", got)
	}
}

func TestBuildRegistrySortsEntriesAndIndexes(t *testing.T) {
	yamlEntry, err := buildEntry(
		0,
		testValueByteRegistration("yaml.public", codec.FormatYAML, codec.MediaTypeYAML),
	)
	requireNoError(t, err)
	jsonEntry, err := buildEntry(
		1,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	registry := buildRegistry([]Entry{yamlEntry, jsonEntry})

	if got := registry.entries[0].ID(); got != MustEntryID("json.public") {
		t.Fatalf("entries[0].ID() = %q", got)
	}
	entries := registry.EntriesByFormat(codec.FormatYAML)
	if len(entries) != 1 || entries[0].Format() != codec.FormatYAML {
		t.Fatalf("EntriesByFormat(yaml) = %#v", entries)
	}
	entries = registry.EntriesByMediaType(codec.MediaTypeJSON)
	if len(entries) != 1 || entries[0].Format() != codec.FormatJSON {
		t.Fatalf("EntriesByMediaType(json) = %#v", entries)
	}
	if entry, ok := registry.LookupID(MustEntryID("json.public")); !ok || entry.Format() != codec.FormatJSON {
		t.Fatalf("LookupID(json.public) = %#v, %v", entry, ok)
	}
}

func TestNewBuildsIDMediaTypeAndFormatIndexesAfterSorting(t *testing.T) {
	zEntry, err := buildEntry(
		0,
		testValueByteRegistration("json.z", codec.FormatJSON, "application/vnd.z+json"),
	)
	requireNoError(t, err)
	jsonEntry, err := buildEntry(
		1,
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	registry := buildRegistry([]Entry{zEntry, jsonEntry})

	indexes := registry.byFormat[codec.FormatJSON]
	if len(indexes) != 2 {
		t.Fatalf("format group length = %d; want 2", len(indexes))
	}
	if indexes[0] != 0 || indexes[1] != 1 {
		t.Fatalf("format group indexes = %#v; want [0 1]", indexes)
	}
	if got := registry.entries[indexes[0]].MediaTypes()[0]; got != codec.MediaTypeJSON {
		t.Fatalf("first grouped media type = %q; want %q", got, codec.MediaTypeJSON)
	}
	if _, ok := registry.byID[MustEntryID("json.public")]; !ok {
		t.Fatalf("byID missing json.public")
	}
	if indexes := registry.byMediaType[codec.MediaTypeJSON]; len(indexes) != 1 {
		t.Fatalf("byMediaType length = %d; want 1", len(indexes))
	}
}

func TestMediaTypeCount(t *testing.T) {
	entry := Entry{
		id: MustEntryID("json.public"),
		info: codec.Info{
			MediaTypes: []codec.MediaType{
				codec.MediaTypeJSON,
				"application/vnd.arcoris.object+json",
			},
		},
	}

	if got := mediaTypeCount([]Entry{entry, entry}); got != 4 {
		t.Fatalf("mediaTypeCount() = %d", got)
	}
}

func TestCompareEntriesByIDThenFormatThenMediaType(t *testing.T) {
	a := Entry{id: MustEntryID("json.a"), info: codec.Info{Format: codec.FormatYAML, MediaTypes: []codec.MediaType{codec.MediaTypeYAML}}}
	b := Entry{id: MustEntryID("json.b"), info: codec.Info{Format: codec.FormatJSON, MediaTypes: []codec.MediaType{codec.MediaTypeJSON}}}

	if got := compareEntries(a, b); got >= 0 {
		t.Fatalf("compareEntries(id a, id b) = %d", got)
	}

	a.id = MustEntryID("json.same")
	b.id = MustEntryID("json.same")
	a.info.Format = codec.FormatJSON
	b.info.Format = codec.FormatYAML
	if got := compareEntries(a, b); got >= 0 {
		t.Fatalf("compareEntries(json, yaml) = %d", got)
	}

	a.info.Format = codec.FormatJSON
	b.info.Format = codec.FormatJSON
	if got := compareEntries(a, b); got <= 0 {
		t.Fatalf("compareEntries(media yaml, media json) = %d", got)
	}
}

func TestFirstMediaType(t *testing.T) {
	if got := firstMediaType(codec.Info{}); got != "" {
		t.Fatalf("firstMediaType(empty) = %q", got)
	}

	info := codec.Info{MediaTypes: []codec.MediaType{codec.MediaTypeJSON}}
	if got := firstMediaType(info); got != codec.MediaTypeJSON {
		t.Fatalf("firstMediaType() = %q", got)
	}
}

func TestCompareText(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{name: "less", a: "a", b: "b", want: -1},
		{name: "equal", a: "a", b: "a", want: 0},
		{name: "greater", a: "b", b: "a", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareText(tt.a, tt.b); got != tt.want {
				t.Fatalf("compareText(%q, %q) = %d", tt.a, tt.b, got)
			}
		})
	}
}
