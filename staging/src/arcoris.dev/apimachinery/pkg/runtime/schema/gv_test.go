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

// GroupVersion tests pin the apiVersion canonical form used by object payloads
// and by all versioned type/resource identities.

// TestGroupVersionParseValidValues verifies core and named apiVersion forms.
func TestGroupVersionParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupVersion
	}{
		{input: "v1", want: GroupVersion{Version: "v1"}},
		{input: "control.arcoris.dev/v1alpha1", want: GroupVersion{Group: "control.arcoris.dev", Version: "v1alpha1"}},
	}

	for _, tt := range tests {
		got, err := ParseGroupVersion(tt.input)
		if err != nil {
			t.Fatalf("ParseGroupVersion(%q) returned error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("ParseGroupVersion(%q) = %+v, want %+v", tt.input, got, tt.want)
		}
		if got.String() != tt.input || got.APIVersion() != tt.input || got.Identifier() != tt.input {
			t.Fatalf("canonical strings for %q were not stable", tt.input)
		}
	}
}

// TestParseAPIVersion verifies apiVersion field parsing is strict GroupVersion parsing.
func TestParseAPIVersion(t *testing.T) {
	got, err := ParseAPIVersion("control.arcoris.dev/v1alpha1")
	if err != nil {
		t.Fatalf("ParseAPIVersion returned error: %v", err)
	}
	want := GroupVersion{Group: "control.arcoris.dev", Version: "v1alpha1"}
	if got != want {
		t.Fatalf("ParseAPIVersion = %+v, want %+v", got, want)
	}
}

// TestGroupVersionRejectsInvalidValues verifies missing and path-like versions fail.
func TestGroupVersionRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"", "apps", "apps/", "apps/v1/extra", "apps/V1", "apps /v1"} {
		if _, err := ParseGroupVersion(input); err == nil {
			t.Fatalf("ParseGroupVersion(%q) expected error", input)
		}
	}
}

// TestGroupVersionHelpers verifies composition into versioned kind and resource IDs.
func TestGroupVersionHelpers(t *testing.T) {
	gv := GroupVersion{Group: "control.arcoris.dev", Version: "v1alpha1"}
	if got := gv.WithKind("WorkloadClass"); got != (GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: "WorkloadClass"}) {
		t.Fatalf("WithKind = %+v", got)
	}
	if got := gv.WithResource("workloadclasses"); got != (GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: "workloadclasses"}) {
		t.Fatalf("WithResource = %+v", got)
	}
}

// TestGroupVersionTextAndJSONRoundTrip verifies scalar encodings and zero rejection.
func TestGroupVersionTextAndJSONRoundTrip(t *testing.T) {
	gv := GroupVersion{Group: "control.arcoris.dev", Version: "v1alpha1"}
	text, err := gv.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupVersion
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gv {
		t.Fatalf("fromText = %+v, want %+v", fromText, gv)
	}

	data, err := json.Marshal(gv)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}

	var fromJSON GroupVersion
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gv {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gv)
	}

	if _, err := (GroupVersion{}).MarshalText(); err == nil {
		t.Fatalf("MarshalText on zero value expected error")
	}
	if err := json.Unmarshal([]byte(`123`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
