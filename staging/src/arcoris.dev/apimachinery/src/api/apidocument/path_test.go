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

func TestPathString(t *testing.T) {
	if got := apidocument.Path("object.metadata").String(); got != "object.metadata" {
		t.Fatalf("Path.String() = %q; want object.metadata", got)
	}
}

func TestPathIsZero(t *testing.T) {
	if !apidocument.Path("").IsZero() {
		t.Fatalf("zero Path did not report IsZero")
	}
	if apidocument.Path("object").IsZero() {
		t.Fatalf("non-zero Path reported IsZero")
	}
}

func TestPathChild(t *testing.T) {
	if got := apidocument.Path("").Child(apidocument.ObjectFieldDesired); got.String() != "desired" {
		t.Fatalf("zero child = %q; want desired", got)
	}
	if got := apidocument.Path("object").Child(apidocument.ObjectFieldDesired); got.String() != "object.desired" {
		t.Fatalf("object child = %q; want object.desired", got)
	}
}
