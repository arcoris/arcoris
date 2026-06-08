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

package admissioncatalog

import (
	"errors"
	"testing"
)

func TestCatalogErrorsSupportIsAndAs(t *testing.T) {
	err := NilCatalogError{Operation: "merge", Index: 0, Path: "catalogs[0]"}
	if !errors.Is(err, ErrNilCatalog) {
		t.Fatal("nil catalog error does not match sentinel")
	}
	var typed NilCatalogError
	if !errors.As(err, &typed) {
		t.Fatal("nil catalog error does not expose typed value")
	}
}

func TestCatalogErrorsExposeDetails(t *testing.T) {
	err := NilCatalogError{Operation: "merge", Index: 3, Path: "catalogs[3]"}
	if err.Operation != "merge" || err.Index != 3 || err.Path != "catalogs[3]" {
		t.Fatalf("NilCatalogError = %+v", err)
	}
}

func TestNilCatalogErrorUsesPathInMessage(t *testing.T) {
	err := NilCatalogError{Operation: "merge", Index: 2, Path: "catalogs[2]"}
	if got, want := err.Error(), "admissioncatalog: nil catalog at catalogs[2]"; got != want {
		t.Fatalf("Error = %q, want %q", got, want)
	}
}
