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
	"math"
	"testing"
)

func TestBoolLiteral(t *testing.T) {
	value := BoolLiteral(true)

	got, ok := value.AsBool()
	requireEqual(t, ok, true)
	requireEqual(t, got, true)
	requireEqual(t, value.Kind(), LiteralBool)
}

func TestStringLiteral(t *testing.T) {
	value := StringLiteral("Ready")

	got, ok := value.AsString()
	requireEqual(t, ok, true)
	requireEqual(t, got, "Ready")
	requireEqual(t, value.Kind(), LiteralString)
}

func TestStringLiteralAllowsEmptyString(t *testing.T) {
	value := StringLiteral("")

	got, ok := value.AsString()
	requireEqual(t, ok, true)
	requireEqual(t, got, "")
}

func TestInt64Literal(t *testing.T) {
	value := Int64Literal(-7)

	got, ok := value.AsInt64()
	requireEqual(t, ok, true)
	requireEqual(t, got, int64(-7))
}

func TestUint64Literal(t *testing.T) {
	value := Uint64Literal(7)

	got, ok := value.AsUint64()
	requireEqual(t, ok, true)
	requireEqual(t, got, uint64(7))
}

func TestIntegerLiteralMinInt64(t *testing.T) {
	value := Int64Literal(math.MinInt64)

	got, ok := value.AsInt64()
	requireEqual(t, ok, true)
	requireEqual(t, got, int64(math.MinInt64))
}

func TestIntegerLiteralMaxUint64(t *testing.T) {
	value := Uint64Literal(math.MaxUint64)

	got, ok := value.AsUint64()
	requireEqual(t, ok, true)
	requireEqual(t, got, uint64(math.MaxUint64))
}
