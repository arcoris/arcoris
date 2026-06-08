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

func TestNewPreference(t *testing.T) {
	contentType := testContentType(codec.MediaTypeJSON)

	preference, err := NewPreference(contentType, WeightDefault)
	requireNoError(t, err)

	if !preference.ContentType().Equal(contentType) {
		t.Fatalf("ContentType() = %q; want %q", preference.ContentType(), contentType)
	}
	if preference.Weight() != WeightDefault {
		t.Fatalf("Weight() = %d; want %d", preference.Weight(), WeightDefault)
	}
	if preference.IsZero() {
		t.Fatalf("IsZero() = true; want false")
	}
}

func TestNewPreferenceRejectsInvalidContentType(t *testing.T) {
	_, err := NewPreference(ContentType{mediaType: "application/json; profile=canonical"}, WeightDefault)

	requireErrorIs(t, err, ErrInvalidPreference)
	requireErrorIs(t, err, ErrInvalidContentType)
	requireSelectionError(t, err, "codecselection.preference.contentType", ErrorReasonInvalidPreference)
}

func TestNewPreferenceRejectsInvalidWeight(t *testing.T) {
	_, err := NewPreference(testContentType(codec.MediaTypeJSON), Weight(0))

	requireErrorIs(t, err, ErrInvalidPreference)
	requireSelectionError(t, err, "codecselection.preference.weight", ErrorReasonInvalidPreference)
}

func TestMustPreferencePanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustPreference did not panic")
		}
	}()

	_ = MustPreference(ContentType{}, WeightDefault)
}
