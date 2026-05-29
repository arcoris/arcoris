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

import "testing"

func TestValidateQualifiedName(t *testing.T) {
	opts := QualifiedNameOptions{
		AllowPrefix:   true,
		MaxNameLength: 20,
		AllowNameDot:  true,
	}

	valid := []string{
		"name",
		"control.arcoris.dev/name",
		"control.arcoris.dev/name.segment",
		"control.arcoris.dev/name-segment",
	}
	for _, value := range valid {
		t.Run("valid "+value, func(t *testing.T) {
			requireValid(t, ValidateQualifiedName(value, opts))
		})
	}

	invalid := []struct {
		name   string
		value  string
		reason Reason
	}{
		{name: "invalid prefix", value: "control/name", reason: ReasonInvalidForm},
		{name: "empty prefix", value: "/name", reason: ReasonInvalidForm},
		{name: "empty name", value: "control.arcoris.dev/", reason: ReasonEmptyValue},
		{name: "multiple slash", value: "control.arcoris.dev/name/extra", reason: ReasonInvalidForm},
		{name: "uppercase name", value: "control.arcoris.dev/Name", reason: ReasonInvalidCharacter},
		{name: "leading hyphen", value: "control.arcoris.dev/-name", reason: ReasonInvalidEdge},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			requireViolation(t, ValidateQualifiedName(tc.value, opts), tc.reason)
		})
	}
}

func TestValidateQualifiedNameRequiresPrefix(t *testing.T) {
	opts := QualifiedNameOptions{
		RequirePrefix: true,
		MaxNameLength: 20,
	}

	requireValid(t, ValidateQualifiedName("control.arcoris.dev/name", opts))
	requireViolation(t, ValidateQualifiedName("name", opts), ReasonInvalidForm)
}

func TestValidateQualifiedNameRejectsPrefixWhenDisallowed(t *testing.T) {
	opts := QualifiedNameOptions{MaxNameLength: 20}

	requireValid(t, ValidateQualifiedName("name", opts))
	requireViolation(t, ValidateQualifiedName("control.arcoris.dev/name", opts), ReasonInvalidForm)
}
