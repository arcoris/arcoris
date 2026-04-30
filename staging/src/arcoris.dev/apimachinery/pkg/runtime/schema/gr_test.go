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

// GroupResource tests keep unversioned resource identity strict and separate
// from versioned GroupVersionResource parsing.

// TestGroupResourceParseValidValues verifies core and named unversioned resource forms.
func TestGroupResourceParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupResource
	}{
		{input: "pods", want: GroupResource{Resource: "pods"}},
		{input: "workloadclasses.control.arcoris.dev", want: GroupResource{Group: "control.arcoris.dev", Resource: "workloadclasses"}},
	}

	for _, tt := range tests {
		got, err := ParseGroupResource(tt.input)
		if err != nil {
			t.Fatalf("ParseGroupResource(%q) returned error: %v", tt.input, err)
		}
		if got != tt.want || got.String() != tt.input || got.Identifier() != tt.input {
			t.Fatalf("ParseGroupResource(%q) = %+v", tt.input, got)
		}
	}
}

// TestGroupResourceRejectsInvalidValues verifies malformed group and resource segments fail.
func TestGroupResourceRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"Pods", "pods.Foo", "pods.control_arcoris_dev", ".control.arcoris.dev", "pods."} {
		if _, err := ParseGroupResource(input); err == nil {
			t.Fatalf("ParseGroupResource(%q) expected error", input)
		}
	}
}

// TestGroupResourceTextAndJSONRoundTrip verifies scalar encodings for GroupResource.
func TestGroupResourceTextAndJSONRoundTrip(t *testing.T) {
	gr := GroupResource{Group: "control.arcoris.dev", Resource: "workloadclasses"}
	text, err := gr.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupResource
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gr {
		t.Fatalf("fromText = %+v, want %+v", fromText, gr)
	}

	data, err := json.Marshal(gr)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	var fromJSON GroupResource
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gr {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gr)
	}
	if err := json.Unmarshal([]byte(`[]`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
