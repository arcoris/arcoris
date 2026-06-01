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

package fieldpath

import (
	"errors"
	"testing"
)

func TestParseRoundTripExamples(t *testing.T) {
	testCases := []string{
		`$`,
		`$.spec.replicas`,
		`$.metadata.labels["app"]`,
		`$.metadata.labels["app.kubernetes.io/name"]`,
		`$.containers[0].image`,
		`$.conditions[{"type":"Ready"}].status`,
		`$.ports[{"name":"http","protocol":"TCP"}].port`,
		`$.routes[{"host":"api.example.com","port":443}].backend`,
		`$."api-version"`,
		`$."x-y-z"[{"name":"a\"b"}]`,
	}

	for _, text := range testCases {
		t.Run(text, func(t *testing.T) {
			path, err := Parse(text)
			requireNoError(t, err)
			requireEqual(t, path.String(), text)
		})
	}
}

func TestParseRejectsInvalidSyntax(t *testing.T) {
	testCases := []struct {
		name   string
		text   string
		target error
	}{
		{name: "missing root", text: `.spec`, target: ErrInvalidSyntax},
		{name: "truncated field", text: `$.`, target: ErrInvalidSyntax},
		{name: "empty brackets", text: `$[]`, target: ErrInvalidSyntax},
		{name: "negative index", text: `$.items[-1]`, target: ErrInvalidSyntax},
		{name: "unclosed key", text: `$.labels["app"`, target: ErrInvalidSyntax},
		{name: "empty selector", text: `$.conditions[{}]`, target: ErrInvalidSelector},
		{name: "truncated selector", text: `$.conditions[{"type":"Ready"]`, target: ErrInvalidSyntax},
		{name: "unsupported literal", text: `$.conditions[{"type":null}]`, target: ErrInvalidSyntax},
		{
			name:   "duplicate selector field",
			text:   `$.conditions[{"type":"Ready","type":"Scheduled"}]`,
			target: ErrDuplicateField,
		},
		{name: "invalid bool token", text: `$.conditions[{"ready":truthy}]`, target: ErrInvalidSyntax},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Parse(testCase.text)
			requireErrorIs(t, err, ErrInvalidPath)
			requireErrorIs(t, err, testCase.target)
		})
	}
}

func TestParseReturnsStructuredSyntaxError(t *testing.T) {
	_, err := Parse(`$.labels["app"`)

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireErrorIs(t, err, ErrInvalidSyntax)
	requireEqual(t, pathErr.Reason, ErrorReasonInvalidSyntax)
	requireEqual(t, pathErr.Detail != "", true)
}

func TestParseQuotedFieldAndKeyDistinction(t *testing.T) {
	path, err := Parse(`$."api-version"["api-version"]`)
	requireNoError(t, err)

	elements := path.Elements()
	requireEqual(t, len(elements), 2)
	requireEqual(t, elements[0].Kind(), ElementField)
	requireEqual(t, elements[1].Kind(), ElementKey)
}
