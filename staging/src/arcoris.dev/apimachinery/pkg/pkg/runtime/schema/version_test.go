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

// Version tests lock down the ARCORIS-only version grammar and reject broader
// release naming schemes before they reach higher API layers.

// TestVersionValidValues verifies every accepted ARCORIS version token family.
func TestVersionValidValues(t *testing.T) {
	for _, input := range []string{"v0", "v1", "v10", "v1alpha1", "v1alpha2", "v1beta1", "v2beta3"} {
		got, err := ParseVersion(input)
		if err != nil {
			t.Fatalf("ParseVersion(%q) returned error: %v", input, err)
		}
		if got.String() != input {
			t.Fatalf("ParseVersion(%q).String() = %q", input, got.String())
		}
		if err := got.Validate(); err != nil {
			t.Fatalf("ParseVersion(%q).Validate() returned error: %v", input, err)
		}
	}
}

// TestVersionInvalidValues verifies rejection of legacy and malformed versions.
func TestVersionInvalidValues(t *testing.T) {
	for _, input := range []string{"", "1", "v", "v01", "v1alpha", "v1beta", "v1alpha0", "v1beta0", "v1rc1", "v1-preview1", "V1", "v1Alpha1"} {
		if _, err := ParseVersion(input); err == nil {
			t.Fatalf("ParseVersion(%q) expected error", input)
		}
	}
}

// TestVersionTextRoundTrip verifies canonical text encoding and strict parsing.
func TestVersionTextRoundTrip(t *testing.T) {
	version := Version("v1alpha1")
	data, err := version.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}
	if string(data) != "v1alpha1" {
		t.Fatalf("MarshalText = %q", string(data))
	}

	var decoded Version
	if err := decoded.UnmarshalText(data); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if decoded != version {
		t.Fatalf("decoded = %q, want %q", decoded, version)
	}

	if _, err := Version("v01").MarshalText(); err == nil {
		t.Fatalf("MarshalText on invalid direct literal expected error")
	}
	if err := decoded.UnmarshalText([]byte("v01")); err == nil {
		t.Fatalf("UnmarshalText on invalid input expected error")
	}
}

// TestVersionJSONRoundTrip verifies scalar JSON encoding and strict string parsing.
func TestVersionJSONRoundTrip(t *testing.T) {
	version := Version("v1beta1")
	data, err := json.Marshal(version)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"v1beta1"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var decoded Version
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if decoded != version {
		t.Fatalf("decoded = %q, want %q", decoded, version)
	}

	if _, err := json.Marshal(Version("v1rc1")); err == nil {
		t.Fatalf("MarshalJSON on invalid direct literal expected error")
	}
	if err := json.Unmarshal([]byte(`"v1rc1"`), &decoded); err == nil {
		t.Fatalf("UnmarshalJSON on invalid string expected error")
	}
}

// TestVersionRejectsNonStringJSON verifies that JSON null and non-strings are rejected.
func TestVersionRejectsNonStringJSON(t *testing.T) {
	for _, input := range []string{`{}`, `[]`, `123`, `true`, `null`} {
		var version Version
		if err := json.Unmarshal([]byte(input), &version); err == nil {
			t.Fatalf("UnmarshalJSON(%s) expected error", input)
		}
	}
}
