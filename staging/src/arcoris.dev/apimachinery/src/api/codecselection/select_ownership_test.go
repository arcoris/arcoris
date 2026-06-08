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

func TestSelectObjectOwnershipDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testOwnershipByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetObjectOwnership, TransportBytes)

	selection, selected, err := selector.SelectObjectOwnershipDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportBytes || selection.Target != codec.TargetObjectOwnership {
		t.Fatalf("selection metadata = %#v; want ownership byte decode", selection)
	}
}

func TestSelectObjectOwnershipEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testOwnershipByteRegistration("json.public", codec.MediaTypeJSON), codec.TargetObjectOwnership, TransportBytes)

	selection, selected, err := selector.SelectObjectOwnershipEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportBytes || selection.Target != codec.TargetObjectOwnership {
		t.Fatalf("selection metadata = %#v; want ownership byte encode", selection)
	}
}

func TestSelectObjectOwnershipStreamDecoder(t *testing.T) {
	selector := testSelectorForDecode(t, testOwnershipStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetObjectOwnership, TransportStream)

	selection, selected, err := selector.SelectObjectOwnershipStreamDecoder(testContentType(codec.MediaTypeJSON))
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetObjectOwnership {
		t.Fatalf("selection metadata = %#v; want ownership stream decode", selection)
	}
}

func TestSelectObjectOwnershipStreamEncoder(t *testing.T) {
	selector := testSelectorForEncode(t, testOwnershipStreamRegistration("json.public", codec.MediaTypeJSON), codec.TargetObjectOwnership, TransportStream)

	selection, selected, err := selector.SelectObjectOwnershipStreamEncoder(
		testPreferenceSet(testPreference(testContentType(codec.MediaTypeJSON), int(WeightDefault))),
	)
	requireNoError(t, err)

	if selected == nil {
		t.Fatalf("selected codec = nil")
	}
	if selection.Transport != TransportStream || selection.Target != codec.TargetObjectOwnership {
		t.Fatalf("selection metadata = %#v; want ownership stream encode", selection)
	}
}
