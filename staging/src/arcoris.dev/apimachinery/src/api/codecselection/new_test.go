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

package codecselection

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func TestNewEmptyConfig(t *testing.T) {
	selector, err := New(Config{})
	requireNoError(t, err)

	if !selector.IsZero() {
		t.Fatalf("selector IsZero() = false; want true")
	}
	if len(selector.DecodeBindings()) != 0 {
		t.Fatalf("DecodeBindings length = %d; want 0", len(selector.DecodeBindings()))
	}
	if len(selector.EncodeBindings()) != 0 {
		t.Fatalf("EncodeBindings length = %d; want 0", len(selector.EncodeBindings()))
	}
}

func TestNewAcceptsValidDecodeBinding(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.public", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)

	selector := testSelector(t, Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})

	bindings := selector.DecodeBindings()
	if len(bindings) != 1 {
		t.Fatalf("DecodeBindings length = %d; want 1", len(bindings))
	}
	if bindings[0].EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("EntryID = %q; want json.public", bindings[0].EntryID)
	}
}

func TestNewAcceptsValidEncodeBinding(t *testing.T) {
	registry := testRegistry(
		t,
		testObjectByteRegistration("json.canonical", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)

	selector := testSelector(t, Config{
		Registry: registry,
		EncodeBindings: []EncodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.canonical"),
		}},
	})

	bindings := selector.EncodeBindings()
	if len(bindings) != 1 {
		t.Fatalf("EncodeBindings length = %d; want 1", len(bindings))
	}
	if bindings[0].EntryID != codecregistry.MustEntryID("json.canonical") {
		t.Fatalf("EntryID = %q; want json.canonical", bindings[0].EntryID)
	}
}

func TestNewAllowsMultipleRegistryEntriesSharingMediaType(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
		testFullByteRegistration("json.storage", codec.MediaTypeJSON),
		testFullByteRegistration("json.canonical", codec.MediaTypeJSON),
	)
	contentType := testContentType(codec.MediaTypeJSON)

	selector := testSelector(t, Config{
		Registry: registry,
		DecodeBindings: []DecodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
		EncodeBindings: []EncodeBinding{{
			ContentType: contentType,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.canonical"),
		}},
	})

	decodeSelection, _, err := selector.SelectObjectDecoder(contentType)
	requireNoError(t, err)
	if decodeSelection.EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("decode EntryID = %q; want json.public", decodeSelection.EntryID)
	}

	encodeSelection, _, err := selector.SelectObjectEncoder(
		testPreferenceSet(testPreference(contentType, int(WeightDefault))),
	)
	requireNoError(t, err)
	if encodeSelection.EntryID != codecregistry.MustEntryID("json.canonical") {
		t.Fatalf("encode EntryID = %q; want json.canonical", encodeSelection.EntryID)
	}
}

func TestNewDetachesInputBindingSlices(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeJSON),
		testFullByteRegistration("json.storage", codec.MediaTypeCBOR),
	)
	contentType := testContentType(codec.MediaTypeJSON)
	decodeBindings := []DecodeBinding{{
		ContentType: contentType,
		Target:      codec.TargetObject,
		Transport:   TransportBytes,
		EntryID:     codecregistry.MustEntryID("json.public"),
	}}
	encodeBindings := []EncodeBinding{{
		ContentType: contentType,
		Target:      codec.TargetObject,
		Transport:   TransportBytes,
		EntryID:     codecregistry.MustEntryID("json.public"),
	}}

	selector := testSelector(t, Config{
		Registry:       registry,
		DecodeBindings: decodeBindings,
		EncodeBindings: encodeBindings,
	})

	decodeBindings[0].EntryID = codecregistry.MustEntryID("json.storage")
	encodeBindings[0].EntryID = codecregistry.MustEntryID("json.storage")

	decodeSelection, _, err := selector.SelectObjectDecoder(contentType)
	requireNoError(t, err)
	if decodeSelection.EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("decode EntryID = %q; want json.public", decodeSelection.EntryID)
	}

	encodeSelection, _, err := selector.SelectObjectEncoder(
		testPreferenceSet(testPreference(contentType, int(WeightDefault))),
	)
	requireNoError(t, err)
	if encodeSelection.EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("encode EntryID = %q; want json.public", encodeSelection.EntryID)
	}
}
