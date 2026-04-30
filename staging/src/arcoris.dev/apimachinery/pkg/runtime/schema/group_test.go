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
	"strings"
	"testing"
)

// Group tests pin the core-group empty value and the DNS-1123 subdomain
// contract used for every named API group.

// TestGroupValidValues verifies core and named groups accepted by the contract.
func TestGroupValidValues(t *testing.T) {
	for _, input := range []string{"", "control.arcoris.dev", "a", "a-b.c0"} {
		got, err := ParseGroup(input)
		if err != nil {
			t.Fatalf("ParseGroup(%q) returned error: %v", input, err)
		}
		if got.String() != input {
			t.Fatalf("ParseGroup(%q).String() = %q", input, got.String())
		}
		if err := got.Validate(); err != nil {
			t.Fatalf("ParseGroup(%q).Validate() returned error: %v", input, err)
		}
	}
}

// TestGroupInvalidValues verifies that named groups stay strict DNS subdomains.
func TestGroupInvalidValues(t *testing.T) {
	label := strings.Repeat("a", 64)
	longGroupLabel := strings.Repeat("a", 63)
	longGroup := strings.Join([]string{longGroupLabel, longGroupLabel, longGroupLabel, longGroupLabel}, ".")

	for _, input := range []string{
		"Control.arcoris.dev",
		"control_arcoris_dev",
		"-control.arcoris.dev",
		"control-.arcoris.dev",
		"control..arcoris.dev",
		label + ".arcoris.dev",
		longGroup,
		"control/arcoris/dev",
		"control arcoris dev",
	} {
		if _, err := ParseGroup(input); err == nil {
			t.Fatalf("ParseGroup(%q) expected error", input)
		}
	}
}

// TestGroupTextRoundTrip verifies canonical text encoding and strict decoding.
func TestGroupTextRoundTrip(t *testing.T) {
	group := Group("control.arcoris.dev")
	data, err := group.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned error: %v", err)
	}
	if string(data) != "control.arcoris.dev" {
		t.Fatalf("MarshalText = %q", string(data))
	}

	var decoded Group
	if err := decoded.UnmarshalText(data); err != nil {
		t.Fatalf("UnmarshalText returned error: %v", err)
	}
	if decoded != group {
		t.Fatalf("decoded = %q, want %q", decoded, group)
	}

	if _, err := Group("Control").MarshalText(); err == nil {
		t.Fatalf("MarshalText on invalid direct literal expected error")
	}
	if err := decoded.UnmarshalText([]byte("Control")); err == nil {
		t.Fatalf("UnmarshalText on invalid input expected error")
	}
}

// TestGroupJSONRoundTrip verifies scalar JSON encoding and strict string decoding.
func TestGroupJSONRoundTrip(t *testing.T) {
	group := Group("control.arcoris.dev")
	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("MarshalJSON returned error: %v", err)
	}
	if string(data) != `"control.arcoris.dev"` {
		t.Fatalf("MarshalJSON = %s", data)
	}

	var decoded Group
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}
	if decoded != group {
		t.Fatalf("decoded = %q, want %q", decoded, group)
	}

	if _, err := json.Marshal(Group("Control")); err == nil {
		t.Fatalf("MarshalJSON on invalid direct literal expected error")
	}
	if err := json.Unmarshal([]byte(`"Control"`), &decoded); err == nil {
		t.Fatalf("UnmarshalJSON on invalid string expected error")
	}
}

// TestGroupRejectsNonStringJSON verifies that JSON null and non-strings are rejected.
func TestGroupRejectsNonStringJSON(t *testing.T) {
	for _, input := range []string{`{}`, `[]`, `123`, `true`, `null`} {
		var group Group
		if err := json.Unmarshal([]byte(input), &group); err == nil {
			t.Fatalf("UnmarshalJSON(%s) expected error", input)
		}
	}
}
