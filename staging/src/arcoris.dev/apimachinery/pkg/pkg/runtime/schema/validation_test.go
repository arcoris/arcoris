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
	"strings"
	"testing"
)

// Validation tests cover the shared lexical helpers directly so atomic and
// composite tests can focus on public behavior.

// TestValidateDNS1123LabelValueLengthEdges verifies the 63-byte DNS label limit.
func TestValidateDNS1123LabelValueLengthEdges(t *testing.T) {
	if err := validateDNS1123LabelValue(strings.Repeat("a", 63)); err != nil {
		t.Fatalf("63-byte label returned error: %v", err)
	}
	if err := validateDNS1123LabelValue(strings.Repeat("a", 64)); err == nil {
		t.Fatalf("64-byte label expected error")
	}
}

// TestValidateDNS1123LabelValueCharacterEdges verifies DNS label byte classes.
func TestValidateDNS1123LabelValueCharacterEdges(t *testing.T) {
	valid := []string{"a", "0", "a-b", "a0b"}
	for _, input := range valid {
		if err := validateDNS1123LabelValue(input); err != nil {
			t.Fatalf("validateDNS1123LabelValue(%q) returned error: %v", input, err)
		}
	}

	invalid := []string{"", "-a", "a-", "A", "a_b", "a.b", "a/b"}
	for _, input := range invalid {
		if err := validateDNS1123LabelValue(input); err == nil {
			t.Fatalf("validateDNS1123LabelValue(%q) expected error", input)
		}
	}
}

// TestValidateDNS1123SubdomainLengthEdge verifies the 253-byte group limit.
func TestValidateDNS1123SubdomainLengthEdge(t *testing.T) {
	label63 := strings.Repeat("a", 63)
	label61 := strings.Repeat("b", 61)
	exact253 := strings.Join([]string{label63, label63, label63, label61}, ".")
	if len(exact253) != 253 {
		t.Fatalf("test fixture length = %d, want 253", len(exact253))
	}
	if err := validateGroupValue(exact253); err != nil {
		t.Fatalf("253-byte group returned error: %v", err)
	}

	tooLong := exact253 + "a"
	if err := validateGroupValue(tooLong); err == nil {
		t.Fatalf("254-byte group expected error")
	}
}

// TestASCIIHelpers verifies low-level ASCII predicates used by validators.
func TestASCIIHelpers(t *testing.T) {
	if !isASCIIUpper('A') || isASCIIUpper('a') {
		t.Fatalf("isASCIIUpper returned unexpected result")
	}
	if !isASCIILower('z') || isASCIILower('Z') {
		t.Fatalf("isASCIILower returned unexpected result")
	}
	if !isASCIIDigit('9') || isASCIIDigit('x') {
		t.Fatalf("isASCIIDigit returned unexpected result")
	}
	if !isASCIINonZeroDigit('1') || isASCIINonZeroDigit('0') {
		t.Fatalf("isASCIINonZeroDigit returned unexpected result")
	}
	if !isDNS1123LabelChar('-') || isDNS1123LabelEdge('-') {
		t.Fatalf("DNS label helper returned unexpected result")
	}
}

// TestUnmarshalJSONStringRejectsNonStringJSON verifies shared JSON scalar rejection.
func TestUnmarshalJSONStringRejectsNonStringJSON(t *testing.T) {
	for _, input := range []string{`{}`, `[]`, `123`, `true`, `null`} {
		if _, err := unmarshalJSONString("test", []byte(input)); err == nil {
			t.Fatalf("unmarshalJSONString(%s) expected error", input)
		}
	}
}
