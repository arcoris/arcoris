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

// Kind tests ensure object type names remain ASCII schema identifiers instead
// of accepting display names, paths, or resource-like values.

// TestKindValidValues verifies accepted PascalCase-like ASCII kind names.
func TestKindValidValues(t *testing.T) {
	for _, input := range []string{"Pod", "WorkloadClass", "APIResource"} {
		got, err := ParseKind(input)
		if err != nil {
			t.Fatalf("ParseKind(%q) returned error: %v", input, err)
		}
		if got.String() != input {
			t.Fatalf("ParseKind(%q).String() = %q", input, got.String())
		}
	}
}

// TestKindInvalidValues verifies that kind names reject separators and Unicode.
func TestKindInvalidValues(t *testing.T) {
	for _, input := range []string{"", "pod", "Pod-Set", "Pod_Set", "Pod.Set", "Pod/Scale", "1Pod", "Пакет", "Pod Scale"} {
		if _, err := ParseKind(input); err == nil {
			t.Fatalf("ParseKind(%q) expected error", input)
		}
	}
}

// TestKindTextRoundTrip verifies text encoding rejects invalid direct literals.
func TestKindTextRoundTrip(t *testing.T) {
	kind := Kind("WorkloadClass")
	data, err := kind.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}

	var decoded Kind
	if err := decoded.UnmarshalText(data); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if decoded != kind {
		t.Fatalf("decoded = %q, want %q", decoded, kind)
	}
	if _, err := Kind("pod").MarshalText(); err == nil {
		t.Fatalf("MarshalText on invalid direct literal expected error")
	}
}

// TestKindJSONRoundTrip verifies scalar JSON encoding and non-string rejection.
func TestKindJSONRoundTrip(t *testing.T) {
	kind := Kind("APIResource")
	data, err := json.Marshal(kind)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"APIResource"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var decoded Kind
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if decoded != kind {
		t.Fatalf("decoded = %q, want %q", decoded, kind)
	}
	if err := json.Unmarshal([]byte(`{}`), &decoded); err == nil {
		t.Fatalf("UnmarshalJSON on non-string expected error")
	}
}
