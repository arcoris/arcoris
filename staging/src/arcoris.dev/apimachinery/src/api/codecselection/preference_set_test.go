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

func TestPreferenceSetOrdersByDescendingWeight(t *testing.T) {
	json := testContentType(codec.MediaTypeJSON, MustParameter("profile", "public"))
	cbor := testContentType(codec.MediaTypeCBOR)

	preferences := testPreferenceSet(
		testPreference(json, 100),
		testPreference(cbor, 900),
	)

	items := preferences.Preferences()
	if len(items) != 2 {
		t.Fatalf("preferences length = %d; want 2", len(items))
	}
	if got := items[0].ContentType().MediaType(); got != codec.MediaTypeCBOR {
		t.Fatalf("first media type = %q; want %q", got, codec.MediaTypeCBOR)
	}
	if got := items[1].ContentType().MediaType(); got != codec.MediaTypeJSON {
		t.Fatalf("second media type = %q; want %q", got, codec.MediaTypeJSON)
	}
}

func TestPreferenceSetEqualWeightUsesOriginalOrder(t *testing.T) {
	json := testContentType(codec.MediaTypeJSON, MustParameter("profile", "public"))
	cbor := testContentType(codec.MediaTypeCBOR)

	preferences := testPreferenceSet(
		testPreference(json, 500),
		testPreference(cbor, 500),
	)

	items := preferences.Preferences()
	if got := items[0].ContentType().MediaType(); got != codec.MediaTypeJSON {
		t.Fatalf("first media type = %q; want original first", got)
	}
	if got := items[1].ContentType().MediaType(); got != codec.MediaTypeCBOR {
		t.Fatalf("second media type = %q; want original second", got)
	}
}

func TestPreferenceSetRejectsDuplicateContentType(t *testing.T) {
	contentType := testContentType(codec.MediaTypeJSON)

	_, err := NewPreferenceSet(
		testPreference(contentType, 900),
		testPreference(contentType, 100),
	)

	requireErrorIs(t, err, ErrInvalidPreference)
	requireSelectionError(t, err, "codecselection.preferences[1]", ErrorReasonInvalidPreference)
}

func TestPreferenceSetPreferencesReturnsDetachedSlice(t *testing.T) {
	json := testContentType(codec.MediaTypeJSON, MustParameter("profile", "public"))
	cbor := testContentType(codec.MediaTypeCBOR)
	preferences := testPreferenceSet(
		testPreference(json, 900),
		testPreference(cbor, 100),
	)

	items := preferences.Preferences()
	items[0] = testPreference(cbor, 900)

	again := preferences.Preferences()
	if got := again[0].ContentType().String(); got != json.String() {
		t.Fatalf("detached preference mutation changed source: %q", got)
	}
}
