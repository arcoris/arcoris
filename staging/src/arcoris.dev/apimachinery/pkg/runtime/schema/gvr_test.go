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

// GroupVersionResource tests pin the colon-based resource identity grammar and
// reject Kubernetes-style dotted and URL-like alternatives.

// TestGroupVersionResourceParseValidValues verifies the canonical GVR grammar.
func TestGroupVersionResourceParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  GroupVersionResource
	}{
		{input: "v1:pods", want: GroupVersionResource{Version: "v1", Resource: "pods"}},
		{input: "control.arcoris.dev/v1alpha1:workloadclasses", want: GroupVersionResource{Group: "control.arcoris.dev", Version: "v1alpha1", Resource: "workloadclasses"}},
	}

	for _, tc := range tests {
		got, err := ParseGroupVersionResource(tc.input)
		if err != nil {
			t.Fatalf("ParseGroupVersionResource(%q) returned error: %v", tc.input, err)
		}
		if got != tc.want || got.String() != tc.input || got.Identifier() != tc.input {
			t.Fatalf("ParseGroupVersionResource(%q) = %+v", tc.input, got)
		}
	}
}

// TestGroupVersionResourceRejectsInvalidValues verifies legacy and URL-like forms fail.
func TestGroupVersionResourceRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{
		"deployments.v1.apps",
		"apps/v1/deployments",
		"apps/v1, Resource=deployments",
		"apps/v1:",
		":deployments",
		"apps/V1:deployments",
		"apps/v1:Deployments",
	} {
		if _, err := ParseGroupVersionResource(input); err == nil {
			t.Fatalf("ParseGroupVersionResource(%q) expected error", input)
		}
	}
}

// TestGroupVersionResourceHelpers verifies GVR projection and subresource composition.
func TestGroupVersionResourceHelpers(t *testing.T) {
	gvr := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1alpha1", Resource: "workloadclasses"}
	if got := gvr.GroupVersion(); got != (GroupVersion{Group: gvr.Group, Version: gvr.Version}) {
		t.Fatalf("GroupVersion = %+v", got)
	}
	if got := gvr.GroupResource(); got != (GroupResource{Group: gvr.Group, Resource: gvr.Resource}) {
		t.Fatalf("GroupResource = %+v", got)
	}
	if got := gvr.WithSubresource("status"); got != (GroupVersionResourcePath{Group: gvr.Group, Version: gvr.Version, Resource: gvr.Resource, Subresource: "status"}) {
		t.Fatalf("WithSubresource = %+v", got)
	}
}

// TestGroupVersionResourceTextAndJSONRoundTrip verifies scalar encodings and partial rejection.
func TestGroupVersionResourceTextAndJSONRoundTrip(t *testing.T) {
	gvr := GroupVersionResource{Group: "control.arcoris.dev", Version: "v1alpha1", Resource: "workloadclasses"}
	text, err := gvr.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText GroupVersionResource
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != gvr {
		t.Fatalf("fromText = %+v, want %+v", fromText, gvr)
	}

	data, err := json.Marshal(gvr)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"control.arcoris.dev/v1alpha1:workloadclasses"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var fromJSON GroupVersionResource
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != gvr {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, gvr)
	}
	if _, err := (GroupVersionResource{Resource: "workloadclasses"}).MarshalText(); err == nil {
		t.Fatalf("MarshalText on partial identity expected error")
	}
	if err := json.Unmarshal([]byte(`[]`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
