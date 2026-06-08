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

func TestNewParameter(t *testing.T) {
	parameter, err := NewParameter(" Profile ", " Canonical+V1 ")
	requireNoError(t, err)

	if parameter.Name() != "profile" {
		t.Fatalf("Name() = %q; want profile", parameter.Name())
	}
	if parameter.Value() != "Canonical+V1" {
		t.Fatalf("Value() = %q; want Canonical+V1", parameter.Value())
	}
	if parameter.IsZero() {
		t.Fatalf("IsZero() = true; want false")
	}
}

func TestNewParameterRejectsInvalidName(t *testing.T) {
	_, err := NewParameter("bad name", "value")

	requireErrorIs(t, err, ErrInvalidParameters)
	requireSelectionError(t, err, "codecselection.parameter.name", ErrorReasonInvalidParameters)
}

func TestNewParameterRejectsInvalidValue(t *testing.T) {
	_, err := NewParameter("profile", "bad value")

	requireErrorIs(t, err, ErrInvalidParameters)
	requireSelectionError(t, err, "codecselection.parameter.value", ErrorReasonInvalidParameters)
}

func TestMustParameterPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustParameter did not panic")
		}
	}()

	_ = MustParameter("bad name", "value")
}
