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

package apidocument_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestFieldNameString(t *testing.T) {
	if got := apidocument.FieldName("desired").String(); got != "desired" {
		t.Fatalf("FieldName.String() = %q; want desired", got)
	}
}

func TestFieldNameIsZero(t *testing.T) {
	if !apidocument.FieldName("").IsZero() {
		t.Fatalf("zero FieldName did not report IsZero")
	}
	if apidocument.FieldName("desired").IsZero() {
		t.Fatalf("non-zero FieldName reported IsZero")
	}
}

func assertFieldName(t *testing.T, name string, got apidocument.FieldName, want string) {
	t.Helper()

	if got.String() != want {
		t.Fatalf("%s = %q; want %q", name, got, want)
	}
}
