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
)

func TestSelectObjectDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testObjectByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetObject, TransportBytes)

	selection, selected, err := selector.SelectObjectDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportBytes || selection.Target != codec.TargetObject {
		t.Fatalf("selection metadata = %#v; want object byte decode", selection)
	}
}

func TestSelectObjectEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testObjectByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetObject, TransportBytes)

	selection, selected, err := selector.SelectObjectEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportBytes || selection.Target != codec.TargetObject {
		t.Fatalf("selection metadata = %#v; want object byte encode", selection)
	}
}

func TestSelectObjectStreamDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testObjectStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetObject, TransportStream)

	selection, selected, err := selector.SelectObjectStreamDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetObject {
		t.Fatalf("selection metadata = %#v; want object stream decode", selection)
	}
}

func TestSelectObjectStreamEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testObjectStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetObject, TransportStream)

	selection, selected, err := selector.SelectObjectStreamEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetObject {
		t.Fatalf("selection metadata = %#v; want object stream encode", selection)
	}
}
