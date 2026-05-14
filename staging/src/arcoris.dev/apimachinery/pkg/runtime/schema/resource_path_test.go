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

// ResourcePath tests keep the optional subresource segment separate from the
// required resource collection name.

// TestResourcePathParseValidValues verifies resource-only and subresource paths.
func TestResourcePathParseValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  ResourcePath
	}{
		{input: "pods", want: ResourcePath{Resource: "pods"}},
		{input: "pods/status", want: ResourcePath{Resource: "pods", Subresource: "status"}},
	}

	for _, tc := range tests {
		got, err := ParseResourcePath(tc.input)
		if err != nil {
			t.Fatalf("ParseResourcePath(%q) returned error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Fatalf("ParseResourcePath(%q) = %+v, want %+v", tc.input, got, tc.want)
		}
		if got.String() != tc.input {
			t.Fatalf("ParseResourcePath(%q).String() = %q", tc.input, got.String())
		}
	}
}

// TestResourcePathRejectsInvalidValues verifies empty and multi-segment paths fail.
func TestResourcePathRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"", "/status", "pods/", "pods/status/extra", "pods//status", "Pods/status", "pods/Scale"} {
		if _, err := ParseResourcePath(input); err == nil {
			t.Fatalf("ParseResourcePath(%q) expected error", input)
		}
	}
}

// TestResourcePathHelpers verifies identifier and subresource helper behavior.
func TestResourcePathHelpers(t *testing.T) {
	path := ResourcePath{Resource: "pods", Subresource: "status"}
	if !path.HasSubresource() {
		t.Fatalf("HasSubresource = false, want true")
	}
	if path.Identifier() != "pods/status" {
		t.Fatalf("Identifier = %q", path.Identifier())
	}
	if (ResourcePath{Resource: "pods"}).HasSubresource() {
		t.Fatalf("resource-only path reported subresource")
	}
}

// TestResourcePathTextAndJSONRoundTrip verifies scalar encodings and zero rejection.
func TestResourcePathTextAndJSONRoundTrip(t *testing.T) {
	path := ResourcePath{Resource: "pods", Subresource: "status"}
	text, err := path.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var fromText ResourcePath
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if fromText != path {
		t.Fatalf("fromText = %+v, want %+v", fromText, path)
	}

	data, err := json.Marshal(path)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"pods/status"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var fromJSON ResourcePath
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if fromJSON != path {
		t.Fatalf("fromJSON = %+v, want %+v", fromJSON, path)
	}

	if _, err := (ResourcePath{}).MarshalText(); err == nil {
		t.Fatalf("MarshalText on zero value expected error")
	}
	if err := json.Unmarshal([]byte(`null`), &fromJSON); err == nil {
		t.Fatalf("UnmarshalJSON on null expected error")
	}
}
