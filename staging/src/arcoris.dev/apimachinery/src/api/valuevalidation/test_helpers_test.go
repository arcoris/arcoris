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
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

type testResolver map[types.TypeName]types.Definition

func (r testResolver) Resolve(name types.TypeName) (types.Definition, bool) {
	definition, ok := r[name]
	return definition, ok
}

func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func requireError(
	t *testing.T,
	err error,
	sentinel error,
	reason valuevalidation.ErrorReason,
	path string,
) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is(%v) = false", sentinel)
	}

	validationError := findValidationError(t, err, sentinel, reason, path)
	if validationError == nil {
		t.Fatalf("missing diagnostic sentinel=%v reason=%s path=%q in %v", sentinel, reason, path, err)
	}
	if validationError.Detail == "" {
		t.Fatalf("diagnostic detail is empty: %#v", validationError)
	}
}

func findValidationError(
	t *testing.T,
	err error,
	sentinel error,
	reason valuevalidation.ErrorReason,
	path string,
) *valuevalidation.Error {
	t.Helper()

	var list valuevalidation.ErrorList
	if errors.As(err, &list) {
		for _, childErr := range list.Unwrap() {
			found := findValidationError(t, childErr, sentinel, reason, path)
			if found != nil {
				return found
			}
		}

		return nil
	}

	var validationError *valuevalidation.Error
	if errors.As(err, &validationError) &&
		errors.Is(validationError, sentinel) &&
		validationError.Reason == reason &&
		validationError.Path == path {
		return validationError
	}

	return nil
}

func requireErrorCount(t *testing.T, err error, want int) {
	t.Helper()

	var list valuevalidation.ErrorList
	if !errors.As(err, &list) {
		t.Fatalf("error is not ErrorList: %T %v", err, err)
	}
	if got := list.Len(); got != want {
		t.Fatalf("error count = %d, want %d: %v", got, want, err)
	}
}

func mustObject(t *testing.T, members ...value.RecordMember) value.Value {
	t.Helper()

	v, err := value.RecordValue(members...)
	if err != nil {
		t.Fatalf("RecordValue() error = %v", err)
	}

	return v
}

func mustList(t *testing.T, items ...value.Value) value.Value {
	t.Helper()

	v, err := value.ListValue(items...)
	if err != nil {
		t.Fatalf("ListValue() error = %v", err)
	}

	return v
}

func mustFloat(t *testing.T, f float64) value.Value {
	t.Helper()

	v, err := value.FloatValue(f)
	if err != nil {
		t.Fatalf("FloatValue() error = %v", err)
	}

	return v
}

func mustDecimal(t *testing.T, text string) value.Value {
	t.Helper()

	decimal, err := value.ParseDecimal(text)
	if err != nil {
		t.Fatalf("ParseDecimal(%q) error = %v", text, err)
	}

	return value.DecimalValue(decimal)
}

func mustDate(t *testing.T, year int, month time.Month, day int) value.Value {
	t.Helper()

	date, err := value.NewDate(year, month, day)
	if err != nil {
		t.Fatalf("NewDate() error = %v", err)
	}

	v, err := value.DateValue(date)
	if err != nil {
		t.Fatalf("DateValue() error = %v", err)
	}

	return v
}

func rootField(names ...string) fieldpath.Path {
	path := fieldpath.RootPath()
	for _, name := range names {
		path = path.Field(name)
	}

	return path
}
