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

// Resource tests protect the single-label collection-name contract and keep
// group/version/subresource segments out of resource values.

// TestResourceValidValues verifies accepted DNS-1123 resource collection names.
func TestResourceValidValues(t *testing.T) {
	for _, input := range []string{"pods", "deployments", "workloadclasses", "queue-partitions"} {
		got, err := ParseResource(input)
		if err != nil {
			t.Fatalf("ParseResource(%q) returned error: %v", input, err)
		}
		if got.String() != input {
			t.Fatalf("ParseResource(%q).String() = %q", input, got.String())
		}
	}
}

// TestResourceInvalidValues verifies resources reject groups, paths, and case drift.
func TestResourceInvalidValues(t *testing.T) {
	for _, input := range []string{"", "Pods", "pod_status", "pods.status", "pods/status", "-pods", "pods-"} {
		if _, err := ParseResource(input); err == nil {
			t.Fatalf("ParseResource(%q) expected error", input)
		}
	}
}

// TestResourceTextRoundTrip verifies canonical text encoding and strict parsing.
func TestResourceTextRoundTrip(t *testing.T) {
	resource := Resource("queue-partitions")
	data, err := resource.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var decoded Resource
	if err := decoded.UnmarshalText(data); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if decoded != resource {
		t.Fatalf("decoded = %q, want %q", decoded, resource)
	}
	if _, err := Resource("Pods").MarshalText(); err == nil {
		t.Fatalf("MarshalText on invalid direct literal expected error")
	}
}

// TestResourceJSONRoundTrip verifies scalar JSON encoding and non-string rejection.
func TestResourceJSONRoundTrip(t *testing.T) {
	resource := Resource("workloadclasses")
	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"workloadclasses"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var decoded Resource
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if decoded != resource {
		t.Fatalf("decoded = %q, want %q", decoded, resource)
	}
	if err := json.Unmarshal([]byte(`[]`), &decoded); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
