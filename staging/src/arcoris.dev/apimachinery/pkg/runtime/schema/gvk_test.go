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

// GroupVersionKind tests pin the only accepted type-identity string grammar and
// explicitly reject dotted legacy triplets.

// TestGroupVersionKindParseValidValues verifies the canonical GVK grammar.
func TestGroupVersionKindParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupVersionKind
	}{
		{input: "v1, Kind=Pod", want: GroupVersionKind{Version: "v1", Kind: "Pod"}},
		{input: "control.arcoris.dev/v1alpha1, Kind=WorkloadClass", want: GroupVersionKind{Group: "control.arcoris.dev", Version: "v1alpha1", Kind: "WorkloadClass"}},
	}

	for _, tt := range tests {
		got, err := ParseGroupVersionKind(tt.input)
		if err != nil {
			t.Fatalf("ParseGroupVersionKind(%q) returned error: %v", tt.input, err)
		}
		if got != tt.want || got.String() != tt.input || got.Identifier() != tt.input {
			t.Fatalf("ParseGroupVersionKind(%q) = %+v", tt.input, got)
		}
	}
}

// TestGroupVersionKindRejectsInvalidValues verifies legacy and partial GVK forms fail.
func TestGroupVersionKindRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{
		"v1 Kind=Pod",
		"WorkloadClass.v1alpha1.control.arcoris.dev",
		"v1, Kind=pod",
		"control_arcoris_dev/v1, Kind=WorkloadClass",
		"control.arcoris.dev/V1, Kind=WorkloadClass",
		"v1, Resource=pods",
	} {
		if _, err := ParseGroupVersionKind(input); err == nil {
			t.Fatalf("ParseGroupVersionKind(%q) expected error", input)
		}
	}
}

// TestParseAPIVersionAndKind verifies strict parsing of split object fields.
func TestParseAPIVersionAndKind(t *testing.T) {
	got, err := ParseAPIVersionAndKind("control.arcoris.dev/v1alpha1", "WorkloadClass")
	if err != nil {
		t.Fatalf("ParseAPIVersionAndKind returned error: %v", err)
	}
	want := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1alpha1", Kind: "WorkloadClass"}
	if got != want {
		t.Fatalf("ParseAPIVersionAndKind = %+v, want %+v", got, want)
	}
	if _, err := ParseAPIVersionAndKind("control.arcoris.dev/v1alpha1", "workloadclass"); err == nil {
		t.Fatalf("ParseAPIVersionAndKind with bad kind expected error")
	}
}

// TestGroupVersionKindHelpers verifies GVK projection helper methods.
func TestGroupVersionKindHelpers(t *testing.T) {
	gvk := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1alpha1", Kind: "WorkloadClass"}
	if got := gvk.GroupVersion(); got != (GroupVersion{Group: gvk.Group, Version: gvk.Version}) {
		t.Fatalf("GroupVersion = %+v", got)
	}
	if got := gvk.GroupKind(); got != (GroupKind{Group: gvk.Group, Kind: gvk.Kind}) {
		t.Fatalf("GroupKind = %+v", got)
	}
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	if apiVersion != "control.arcoris.dev/v1alpha1" || kind != "WorkloadClass" {
		t.Fatalf("ToAPIVersionAndKind = %q, %q", apiVersion, kind)
	}
}

// TestGroupVersionKindTextAndJSONRoundTrip verifies scalar encodings and partial rejection.
func TestGroupVersionKindTextAndJSONRoundTrip(t *testing.T) {
	gvk := GroupVersionKind{Group: "control.arcoris.dev", Version: "v1alpha1", Kind: "WorkloadClass"}
	text, err := gvk.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupVersionKind
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gvk {
		t.Fatalf("fromText = %+v, want %+v", fromText, gvk)
	}

	data, err := json.Marshal(gvk)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"control.arcoris.dev/v1alpha1, Kind=WorkloadClass"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var fromJSON GroupVersionKind
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gvk {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gvk)
	}
	if _, err := (GroupVersionKind{Kind: "WorkloadClass"}).MarshalText(); err == nil {
		t.Fatalf("MarshalText on partial identity expected error")
	}
	if err := json.Unmarshal([]byte(`{}`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
