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

func TestSelectValueDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testValueByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetValue, TransportBytes)

	selection, selected, err := selector.SelectValueDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("EntryID = %q; want json.public", selection.EntryID)
	}
	if selection.Direction != DirectionDecode || selection.Transport != TransportBytes || selection.Target != codec.TargetValue {
		t.Fatalf("selection metadata = %#v; want value byte decode", selection)
	}
}

func TestSelectValueEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testValueByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetValue, TransportBytes)

	selection, selected, err := selector.SelectValueEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.EntryID != codecregistry.MustEntryID("json.public") {
		t.Fatalf("EntryID = %q; want json.public", selection.EntryID)
	}
	if selection.Direction != DirectionEncode || selection.Transport != TransportBytes || selection.Target != codec.TargetValue {
		t.Fatalf("selection metadata = %#v; want value byte encode", selection)
	}
}

func TestSelectValueStreamDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testValueStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetValue, TransportStream)

	selection, selected, err := selector.SelectValueStreamDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetValue {
		t.Fatalf("selection metadata = %#v; want value stream decode", selection)
	}
}

func TestSelectValueStreamEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testValueStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetValue, TransportStream)

	selection, selected, err := selector.SelectValueStreamEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetValue {
		t.Fatalf("selection metadata = %#v; want value stream encode", selection)
	}
}
