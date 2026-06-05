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
		{name: "uppercase", value: "Workers", reason: ReasonInvalidEdge},
		{name: "underscore", value: "worker_1", reason: ReasonInvalidCharacter},
		{name: "dot", value: "worker.main", reason: ReasonInvalidCharacter},
		{name: "slash", value: "worker/main", reason: ReasonInvalidCharacter},
		{name: "space", value: "worker 1", reason: ReasonInvalidCharacter},
		{name: "leading hyphen", value: "-worker", reason: ReasonInvalidEdge},
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
		"a.b",
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
		{name: "double dot", value: "control..dev", reason: ReasonInvalidForm},
		{name: "leading dot", value: ".control.dev", reason: ReasonInvalidForm},
		{name: "trailing dot", value: "control.dev.", reason: ReasonInvalidForm},
		{name: "label too long", value: strings.Repeat("a", 64) + ".dev", reason: ReasonInvalidLength},
		{name: "total too long", value: total253 + "x", reason: ReasonInvalidLength},
		{name: "uppercase", value: "Control.arcoris.dev", reason: ReasonInvalidEdge},
		{name: "underscore", value: "control_arcoris.dev", reason: ReasonInvalidCharacter},
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
	requireValid(t, ValidateQualifiedDNS1123Subdomain("control.arcoris.dev"))
	requireViolation(t, ValidateQualifiedDNS1123Subdomain("workers"), ReasonInvalidForm)
	requireViolation(t, ValidateQualifiedDNS1123Subdomain(""), ReasonEmptyValue)
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
