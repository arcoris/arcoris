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

package codecselection

import "testing"

func TestNewParametersSortsByName(t *testing.T) {
	parameters := MustParameters(
		MustParameter("version", "v1"),
		MustParameter("profile", "canonical"),
	)

	items := parameters.Items()
	if len(items) != 2 {
		t.Fatalf("parameters length = %d; want 2", len(items))
	}
	if items[0].Name() != "profile" || items[1].Name() != "version" {
		t.Fatalf("parameter order = %s,%s; want profile,version", items[0].Name(), items[1].Name())
	}
	if parameters.Len() != 2 {
		t.Fatalf("Len() = %d; want 2", parameters.Len())
	}
}

func TestParametersEqual(t *testing.T) {
	a := MustParameters(
		MustParameter("version", "v1"),
		MustParameter("profile", "canonical"),
	)
	b := MustParameters(
		MustParameter("profile", "canonical"),
		MustParameter("version", "v1"),
	)

	if !a.Equal(b) {
		t.Fatalf("Equal() = false; want true")
	}
}

func TestMustParametersPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustParameters did not panic")
		}
	}()

	_ = MustParameters(
		MustParameter("profile", "canonical"),
		MustParameter("PROFILE", "storage"),
	)
}

func TestParametersItemsReturnsDetachedSlice(t *testing.T) {
	parameters := MustParameters(
		MustParameter("profile", "canonical"),
		MustParameter("version", "v1"),
	)

	items := parameters.Items()
	items[0] = MustParameter("profile", "storage")

	again := parameters.Items()
	if again[0].Value() != "canonical" {
		t.Fatalf("detached parameter mutation changed source: %q", again[0].Value())
	}
}
