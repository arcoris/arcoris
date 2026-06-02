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

package valuevalidation_test

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestErrorListUnwrapStillSupportsErrorsIs(t *testing.T) {
	err := valuevalidation.Validate(
		mustObject(t),
		types.Object(
			types.Field("name").String().Required(),
			types.Field("replicas").Int32().Required(),
		).Type(),
		valuevalidation.Options{},
	)

	if !errors.Is(err, valuevalidation.ErrMissingField) {
		t.Fatalf("errors.Is(ErrMissingField) = false")
	}

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}
	if list.IsEmpty() {
		t.Fatalf("ErrorList is empty")
	}
	if !strings.Contains(list.Error(), "2 errors") {
		t.Fatalf("ErrorList summary %q does not include count", list.Error())
	}
}

func TestErrorListErrorsReturnsDetachedCopy(t *testing.T) {
	err := valuevalidation.Validate(
		mustObject(t),
		types.Object(types.Field("name").String().Required()).Type(),
		valuevalidation.Options{},
	)

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}

	copied := list.Errors()
	copied[0] = nil

	if list.First() == nil {
		t.Fatalf("mutating Errors() result changed original list")
	}
}

func TestErrorListFirst(t *testing.T) {
	var empty valuevalidation.ErrorList
	if got := empty.First(); got != nil {
		t.Fatalf("empty First() = %v, want nil", got)
	}

	err := valuevalidation.Validate(
		mustObject(t),
		types.Object(types.Field("name").String().Required()).Type(),
		valuevalidation.Options{},
	)

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}
	if got := list.First(); got == nil {
		t.Fatalf("First() = nil")
	}
}

func TestErrorListFormatAll(t *testing.T) {
	err := valuevalidation.Validate(
		mustObject(t),
		types.Object(
			types.Field("name").String().Required(),
			types.Field("replicas").Int32().Required(),
		).Type(),
		valuevalidation.Options{},
	)

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}

	formatted := list.FormatAll()
	if !strings.Contains(formatted, "$.name") {
		t.Fatalf("FormatAll() = %q, want $.name", formatted)
	}
	if !strings.Contains(formatted, "$.replicas") {
		t.Fatalf("FormatAll() = %q, want $.replicas", formatted)
	}
	if !strings.Contains(formatted, "\n") {
		t.Fatalf("FormatAll() = %q, want multiple lines", formatted)
	}
}

func TestErrorListFormatAllEmpty(t *testing.T) {
	var list valuevalidation.ErrorList

	if got := list.FormatAll(); got != "" {
		t.Fatalf("FormatAll() = %q, want empty string", got)
	}
}

func TestErrorListErrorRemainsCompact(t *testing.T) {
	err := valuevalidation.Validate(
		mustObject(t),
		types.Object(
			types.Field("name").String().Required(),
			types.Field("replicas").Int32().Required(),
		).Type(),
		valuevalidation.Options{},
	)

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}

	summary := list.Error()
	if !strings.Contains(summary, "2 errors; first:") {
		t.Fatalf("Error() = %q, want compact summary", summary)
	}
	if strings.Contains(summary, "$.replicas") {
		t.Fatalf("Error() = %q, want only first diagnostic in compact summary", summary)
	}
}

func TestErrorListEmptySummary(t *testing.T) {
	var list valuevalidation.ErrorList

	if got, want := list.Error(), "valuevalidation: no errors"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
	if got := list.Errors(); got != nil {
		t.Fatalf("Errors() = %#v, want nil", got)
	}
	if got := list.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %#v, want nil", got)
	}
	if got := list.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
	if !list.IsEmpty() {
		t.Fatalf("IsEmpty() = false")
	}
}

func TestErrorListSingleSummary(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("x"),
		types.Int64().Type(),
		valuevalidation.Options{},
	)

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("errors.As(ErrorList) = false")
	}
	if got := list.Error(); !strings.Contains(got, "kind mismatch") {
		t.Fatalf("Error() = %q, want kind mismatch detail", got)
	}
}
