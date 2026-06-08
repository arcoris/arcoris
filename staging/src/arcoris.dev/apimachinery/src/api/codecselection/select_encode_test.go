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

func TestSelectEncodeUsesPreferenceOrder(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.json", codec.MediaTypeJSON),
		testFullByteRegistration("json.cbor", codec.MediaTypeCBOR),
	)
	json := testContentType(codec.MediaTypeJSON)
	cbor := testContentType(codec.MediaTypeCBOR)
	selector := testSelector(t, Config{
		Registry: registry,
		EncodeBindings: []EncodeBinding{
			{
				ContentType: json,
				Target:      codec.TargetObject,
				Transport:   TransportBytes,
				EntryID:     codecregistry.MustEntryID("json.json"),
			},
			{
				ContentType: cbor,
				Target:      codec.TargetObject,
				Transport:   TransportBytes,
				EntryID:     codecregistry.MustEntryID("json.cbor"),
			},
		},
	})

	selection, _, err := selector.SelectObjectEncoder(
		testPreferenceSet(
			testPreference(json, 100),
			testPreference(cbor, 900),
		),
	)
	requireNoError(t, err)

	if selection.EntryID != codecregistry.MustEntryID("json.cbor") {
		t.Fatalf("EntryID = %q; want json.cbor", selection.EntryID)
	}
}

func TestSelectEncodeFallsThroughUnsupportedPreference(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.json", codec.MediaTypeJSON),
	)
	json := testContentType(codec.MediaTypeJSON)
	cbor := testContentType(codec.MediaTypeCBOR)
	selector := testSelector(t, Config{
		Registry: registry,
		EncodeBindings: []EncodeBinding{{
			ContentType: json,
			Target:      codec.TargetObject,
			Transport:   TransportBytes,
			EntryID:     codecregistry.MustEntryID("json.json"),
		}},
	})

	selection, _, err := selector.SelectObjectEncoder(
		testPreferenceSet(
			testPreference(cbor, 900),
			testPreference(json, 100),
		),
	)
	requireNoError(t, err)

	if selection.EntryID != codecregistry.MustEntryID("json.json") {
		t.Fatalf("EntryID = %q; want json.json", selection.EntryID)
	}
}

func TestSelectEncodeMissingPreference(t *testing.T) {
	selector := testSelector(t, Config{})

	_, _, err := selector.SelectObjectEncoder(PreferenceSet{})

	requireErrorIs(t, err, ErrNoEncodePreference)
	requireSelectionError(t, err, "codecselection.encode.object.preferences", ErrorReasonNoEncodePreference)
}

func TestSelectEncodeNoSupportedPreference(t *testing.T) {
	selector := testSelector(t, Config{})

	_, _, err := selector.SelectObjectEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)

	requireErrorIs(t, err, ErrNoEncodePreference)
	requireSelectionError(t, err, "codecselection.encode.object", ErrorReasonNoEncodePreference)
}
