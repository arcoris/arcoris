// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lexical

import (
	"strings"
	"testing"
)

func TestValidateDNS1123Label(t *testing.T) {
	label63 := strings.Repeat("a", MaxDNS1123LabelLength)

	valid := []string{
		"a",
		"abc",
		"a1",
		"a-b",
		"abc-123",
		"control",
		"z",
		"0",
		"workers",
		"worker-1",
		"task01",
		label63,
	}
	for _, value := range valid {
		t.Run("valid "+value, func(t *testing.T) {
			requireValid(t, ValidateDNS1123Label(value))
		})
	}

	invalid := []struct {
		name   string
		value  string
		reason Reason
	}{
		{name: "empty", value: "", reason: ReasonEmptyValue},
		{name: "64 byte label", value: strings.Repeat("a", MaxDNS1123LabelLength+1), reason: ReasonInvalidLength},
		{name: "uppercase single letter", value: "A", reason: ReasonInvalidEdge},
		{name: "uppercase", value: "Workers", reason: ReasonInvalidEdge},
		{name: "underscore short", value: "a_b", reason: ReasonInvalidCharacter},
		{name: "underscore", value: "worker_1", reason: ReasonInvalidCharacter},
		{name: "dot short", value: "a.b", reason: ReasonInvalidCharacter},
		{name: "dot", value: "worker.main", reason: ReasonInvalidCharacter},
		{name: "slash", value: "worker/main", reason: ReasonInvalidCharacter},
		{name: "space", value: "worker 1", reason: ReasonInvalidCharacter},
		{name: "leading hyphen short", value: "-a", reason: ReasonInvalidEdge},
		{name: "leading hyphen", value: "-worker", reason: ReasonInvalidEdge},
		{name: "trailing hyphen short", value: "a-", reason: ReasonInvalidEdge},
		{name: "trailing hyphen", value: "worker-", reason: ReasonInvalidEdge},
		{name: "non ASCII", value: "воркер", reason: ReasonInvalidEdge},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			requireViolation(t, ValidateDNS1123Label(tc.value), tc.reason)
		})
	}
}

func TestValidateDNS1123Subdomain(t *testing.T) {
	total253 := strings.Repeat("a", 63) + "." +
		strings.Repeat("b", 63) + "." +
		strings.Repeat("c", 63) + "." +
		strings.Repeat("d", 61)

	valid := []string{
		"a",
		"abc",
		"a1",
		"a-b",
		"abc-123",
		"control",
		"a.b",
		"a-b.c-d",
		"control.arcoris.dev",
		"x.y-z.1a",
		total253,
	}
	for _, value := range valid {
		t.Run("valid "+value, func(t *testing.T) {
			requireValid(t, ValidateDNS1123Subdomain(value))
		})
	}

	invalid := []struct {
		name   string
		value  string
		reason Reason
	}{
		{name: "empty", value: "", reason: ReasonEmptyValue},
		{name: "single dot", value: ".", reason: ReasonInvalidForm},
		{name: "leading dot short", value: ".a", reason: ReasonInvalidForm},
		{name: "trailing dot short", value: "a.", reason: ReasonInvalidForm},
		{name: "double dot short", value: "a..b", reason: ReasonInvalidForm},
		{name: "double dot", value: "control..dev", reason: ReasonInvalidForm},
		{name: "leading dot", value: ".control.dev", reason: ReasonInvalidForm},
		{name: "trailing dot", value: "control.dev.", reason: ReasonInvalidForm},
		{name: "label too long", value: strings.Repeat("a", 64) + ".dev", reason: ReasonInvalidLength},
		{name: "total too long", value: total253 + "x", reason: ReasonInvalidLength},
		{name: "uppercase single label", value: "A", reason: ReasonInvalidEdge},
		{name: "uppercase", value: "Control.arcoris.dev", reason: ReasonInvalidEdge},
		{name: "underscore short", value: "a_b", reason: ReasonInvalidCharacter},
		{name: "underscore", value: "control_arcoris.dev", reason: ReasonInvalidCharacter},
		{name: "leading hyphen label", value: "-a.example", reason: ReasonInvalidEdge},
		{name: "trailing hyphen label", value: "a-.example", reason: ReasonInvalidEdge},
		{name: "slash", value: "control/arcoris.dev", reason: ReasonInvalidCharacter},
		{name: "space", value: "control arcoris.dev", reason: ReasonInvalidCharacter},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			requireViolation(t, ValidateDNS1123Subdomain(tc.value), tc.reason)
		})
	}
}

func TestValidateQualifiedDNS1123Subdomain(t *testing.T) {
	valid := []string{
		"a.b",
		"control.arcoris.dev",
	}
	for _, value := range valid {
		t.Run("valid "+value, func(t *testing.T) {
			requireValid(t, ValidateQualifiedDNS1123Subdomain(value))
		})
	}

	invalid := []struct {
		name   string
		value  string
		reason Reason
	}{
		{name: "no dot", value: "a", reason: ReasonInvalidForm},
		{name: "empty", value: "", reason: ReasonEmptyValue},
		{name: "single dot", value: ".", reason: ReasonInvalidForm},
		{name: "trailing dot", value: "a.", reason: ReasonInvalidForm},
		{name: "leading dot", value: ".a", reason: ReasonInvalidForm},
		{name: "double dot", value: "a..b", reason: ReasonInvalidForm},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			requireViolation(t, ValidateQualifiedDNS1123Subdomain(tc.value), tc.reason)
		})
	}
}

func requireValid(t *testing.T, err *Violation) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected violation: %v", err)
	}
}

func requireViolation(t *testing.T, err *Violation, reason Reason) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected violation %q", reason)
	}
	if err.Reason != reason {
		t.Fatalf("Violation.Reason = %q, want %q; detail=%q", err.Reason, reason, err.Detail)
	}
	if err.Detail == "" {
		t.Fatalf("Violation.Detail is empty")
	}
}
