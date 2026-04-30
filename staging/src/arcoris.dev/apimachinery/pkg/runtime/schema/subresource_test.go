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

// Subresource tests document the explicit decision that the empty subresource
// is valid and serializes as "", while malformed non-empty names are rejected.

// TestSubresourceValidValues verifies empty and named subresource values.
func TestSubresourceValidValues(t *testing.T) {
	for _, input := range []string{"", "status", "scale", "heartbeat", "renew", "cancel"} {
		got, err := ParseSubresource(input)
		if err != nil {
			t.Fatalf("ParseSubresource(%q) returned error: %v", input, err)
		}
		if got.String() != input {
			t.Fatalf("ParseSubresource(%q).String() = %q", input, got.String())
		}
	}
}

// TestSubresourceInvalidValues verifies malformed non-empty subresources are rejected.
func TestSubresourceInvalidValues(t *testing.T) {
	for _, input := range []string{"Scale", "pod_status", "pods.status", "pods/status", "-status", "status-"} {
		if _, err := ParseSubresource(input); err == nil {
			t.Fatalf("ParseSubresource(%q) expected error", input)
		}
	}
}

// TestSubresourceTextRoundTrip verifies empty subresources marshal as empty text.
func TestSubresourceTextRoundTrip(t *testing.T) {
	for _, subresource := range []Subresource{"", "status"} {
		data, err := subresource.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%q) returned error: %v", subresource, err)
		}

		var decoded Subresource
		if err := decoded.UnmarshalText(data); err != nil {
			t.Fatalf("UnmarshalText(%q) returned error: %v", string(data), err)
		}
		if decoded != subresource {
			t.Fatalf("decoded = %q, want %q", decoded, subresource)
		}
	}
	if _, err := Subresource("Scale").MarshalText(); err == nil {
		t.Fatalf("MarshalText on invalid direct literal expected error")
	}
}

// TestSubresourceJSONRoundTrip verifies empty and named subresources use JSON strings.
func TestSubresourceJSONRoundTrip(t *testing.T) {
	for _, subresource := range []Subresource{"", "scale"} {
		data, err := json.Marshal(subresource)
		if err != nil {
			t.Fatalf("MarshalJSON(%q) returned error: %v", subresource, err)
		}

		var decoded Subresource
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("UnmarshalJSON(%s) returned error: %v", data, err)
		}
		if decoded != subresource {
			t.Fatalf("decoded = %q, want %q", decoded, subresource)
		}
	}
	if err := json.Unmarshal([]byte(`true`), new(Subresource)); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
