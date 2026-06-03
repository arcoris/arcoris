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

package valuepresence

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestPresent(t *testing.T) {
	operand := Present(value.NullValue())
	val, ok := operand.ValueOK()

	if !operand.Present() {
		t.Fatalf("Present().Present() = false")
	}
	if operand.Absent() {
		t.Fatalf("Present().Absent() = true")
	}
	if !ok {
		t.Fatalf("ValueOK ok = false")
	}
	if !val.IsNull() {
		t.Fatalf("ValueOK value is not null")
	}
}

func TestAbsent(t *testing.T) {
	operand := Absent()
	val, ok := operand.ValueOK()

	if operand.Present() {
		t.Fatalf("Absent().Present() = true")
	}
	if !operand.Absent() {
		t.Fatalf("Absent().Absent() = false")
	}
	if ok {
		t.Fatalf("ValueOK ok = true")
	}
	if !val.IsZero() {
		t.Fatalf("ValueOK value is not zero")
	}
}

func TestFromPresent(t *testing.T) {
	operand := From(value.StringValue("x"), true)
	val, ok := operand.ValueOK()

	if !ok {
		t.Fatalf("ValueOK ok = false")
	}

	got, _ := val.String()
	if got != "x" {
		t.Fatalf("ValueOK value = %q, want x", got)
	}
}

func TestFromAbsentIgnoresValue(t *testing.T) {
	operand := From(value.StringValue("x"), false)

	if operand.Present() {
		t.Fatalf("From(_, false).Present() = true")
	}
	if !operand.Value().IsZero() {
		t.Fatalf("From(_, false).Value() is not zero")
	}
}

func TestValueReturnsZeroForAbsent(t *testing.T) {
	if !Absent().Value().IsZero() {
		t.Fatalf("Absent().Value() is not zero")
	}
}

func TestClonePresent(t *testing.T) {
	operand := Present(value.BytesValue([]byte{1, 2, 3}))
	cloned := operand.Clone()

	got, ok := cloned.Value().Bytes()
	if !ok {
		t.Fatalf("Clone().Value() is not bytes")
	}
	got[0] = 9

	original, _ := operand.Value().Bytes()
	if original[0] != 1 {
		t.Fatalf("Clone aliased original bytes")
	}
}

func TestCloneAbsent(t *testing.T) {
	cloned := Absent().Clone()

	if cloned.Present() {
		t.Fatalf("Absent().Clone().Present() = true")
	}
	if !cloned.Value().IsZero() {
		t.Fatalf("Absent().Clone().Value() is not zero")
	}
}
