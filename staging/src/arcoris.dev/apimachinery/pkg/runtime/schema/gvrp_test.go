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

// GroupVersionResourcePath tests the only schema-level form that may include a
// subresource together with group and version.

// TestGroupVersionResourcePathParseValidValues verifies canonical versioned resource paths.
func TestGroupVersionResourcePathParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupVersionResourcePath
	}{
		{input: "v1:pods", want: GroupVersionResourcePath{Version: "v1", Resource: "pods"}},
		{input: "v1:pods/status", want: GroupVersionResourcePath{Version: "v1", Resource: "pods", Subresource: "status"}},
		{input: "control.arcoris.dev/v1alpha1:workloadclasses/status", want: GroupVersionResourcePath{Group: "control.arcoris.dev", Version: "v1alpha1", Resource: "workloadclasses", Subresource: "status"}},
	}

	for _, tc := range tests {
		got, err := ParseGroupVersionResourcePath(tc.input)
		if err != nil {
			t.Fatalf("ParseGroupVersionResourcePath(%q) returned error: %v", tc.input, err)
		}
		if got != tc.want || got.String() != tc.input || got.Identifier() != tc.input {
			t.Fatalf("ParseGroupVersionResourcePath(%q) = %+v", tc.input, got)
		}
	}
}

// TestGroupVersionResourcePathRejectsInvalidValues verifies malformed path segments fail.
func TestGroupVersionResourcePathRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{
		"v1:Pods",
		"v1:pods/Scale",
		"v1:pods/status/extra",
		"v1:pods//status",
		"apps/v1/deployments/status",
		"apps/v1:",
		":pods",
	} {
		if _, err := ParseGroupVersionResourcePath(input); err == nil {
			t.Fatalf("ParseGroupVersionResourcePath(%q) expected error", input)
		}
	}
}

// TestGroupVersionResourcePathHelpers verifies projection helper methods.
func TestGroupVersionResourcePathHelpers(t *testing.T) {
	gvrp := GroupVersionResourcePath{Group: "control.arcoris.dev", Version: "v1alpha1", Resource: "workloadclasses", Subresource: "status"}
	if got := gvrp.GroupVersion(); got != (GroupVersion{Group: gvrp.Group, Version: gvrp.Version}) {
		t.Fatalf("GroupVersion = %+v", got)
	}
	if got := gvrp.GroupVersionResource(); got != (GroupVersionResource{Group: gvrp.Group, Version: gvrp.Version, Resource: gvrp.Resource}) {
		t.Fatalf("GroupVersionResource = %+v", got)
	}
	if got := gvrp.GroupResource(); got != (GroupResource{Group: gvrp.Group, Resource: gvrp.Resource}) {
		t.Fatalf("GroupResource = %+v", got)
	}
	if got := gvrp.ResourcePath(); got != (ResourcePath{Resource: gvrp.Resource, Subresource: gvrp.Subresource}) {
		t.Fatalf("ResourcePath = %+v", got)
	}
	if !gvrp.HasSubresource() {
		t.Fatalf("HasSubresource = false, want true")
	}
}

// TestGroupVersionResourcePathTextAndJSONRoundTrip verifies scalar encodings and partial rejection.
func TestGroupVersionResourcePathTextAndJSONRoundTrip(t *testing.T) {
	gvrp := GroupVersionResourcePath{Version: "v1", Resource: "pods", Subresource: "status"}
	text, err := gvrp.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupVersionResourcePath
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gvrp {
		t.Fatalf("fromText = %+v, want %+v", fromText, gvrp)
	}

	data, err := json.Marshal(gvrp)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"v1:pods/status"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var fromJSON GroupVersionResourcePath
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gvrp {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gvrp)
	}
	if _, err := (GroupVersionResourcePath{Resource: "pods"}).MarshalText(); err == nil {
		t.Fatalf("MarshalText on partial identity expected error")
	}
	if err := json.Unmarshal([]byte(`true`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
