/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

import (
	"encoding/json"
	"testing"
)

// GroupKind tests keep unversioned type identity strict and separate from
// versioned GroupVersionKind parsing.

// TestGroupKindParseValidValues verifies core and named unversioned kind forms.
func TestGroupKindParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupKind
	}{
		{input: "Pod", want: GroupKind{Kind: "Pod"}},
		{input: "WorkloadClass.control.arcoris.dev", want: GroupKind{Group: "control.arcoris.dev", Kind: "WorkloadClass"}},
	}

	for _, tc := range tests {
		got, err := ParseGroupKind(tc.input)
		if err != nil {
			t.Fatalf("ParseGroupKind(%q) returned error: %v", tc.input, err)
		}
		if got != tc.want || got.String() != tc.input || got.Identifier() != tc.input {
			t.Fatalf("ParseGroupKind(%q) = %+v", tc.input, got)
		}
	}
}

// TestGroupKindRejectsInvalidValues verifies malformed group and kind segments fail.
func TestGroupKindRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"pod", "Pod.Foo", "Pod.control_arcoris_dev", ".control.arcoris.dev", "Pod."} {
		if _, err := ParseGroupKind(input); err == nil {
			t.Fatalf("ParseGroupKind(%q) expected error", input)
		}
	}
}

// TestGroupKindTextAndJSONRoundTrip verifies scalar encodings for GroupKind.
func TestGroupKindTextAndJSONRoundTrip(t *testing.T) {
	gk := GroupKind{Group: "control.arcoris.dev", Kind: "WorkloadClass"}
	text, err := gk.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupKind
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gk {
		t.Fatalf("fromText = %+v, want %+v", fromText, gk)
	}

	data, err := json.Marshal(gk)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	var fromJSON GroupKind
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gk {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gk)
	}
	if err := json.Unmarshal([]byte(`{}`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
