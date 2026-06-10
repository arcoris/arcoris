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

import "testing"

func TestLiteralAccessors(t *testing.T) {
	boolValue, boolOK := BoolLiteral(true).AsBool()
	stringValue, stringOK := StringLiteral("Ready").AsString()
	intValue, intOK := Int64Literal(-1).AsInt64()
	uintValue, uintOK := Uint64Literal(1).AsUint64()

	requireEqual(t, boolOK, true)
	requireEqual(t, boolValue, true)
	requireEqual(t, stringOK, true)
	requireEqual(t, stringValue, "Ready")
	requireEqual(t, intOK, true)
	requireEqual(t, intValue, int64(-1))
	requireEqual(t, uintOK, true)
	requireEqual(t, uintValue, uint64(1))
}

func TestLiteralAccessorsRejectWrongKind(t *testing.T) {
	value := BoolLiteral(true)

	_, ok := value.AsString()
	requireEqual(t, ok, false)

	_, ok = value.AsInt64()
	requireEqual(t, ok, false)

	_, ok = value.AsUint64()
	requireEqual(t, ok, false)
}
