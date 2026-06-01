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

func TestLiteralCompareBool(t *testing.T) {
	requireEqual(t, BoolLiteral(false).Compare(BoolLiteral(true)), -1)
	requireEqual(t, BoolLiteral(true).Compare(BoolLiteral(false)), 1)
	requireEqual(t, BoolLiteral(true).Compare(BoolLiteral(true)), 0)
}

func TestLiteralCompareIntegerSignedUnsigned(t *testing.T) {
	requireEqual(t, Int64Literal(-1).Compare(Uint64Literal(0)), -1)
	requireEqual(t, Uint64Literal(0).Compare(Int64Literal(-1)), 1)
	requireEqual(t, Int64Literal(7).Compare(Uint64Literal(7)), 0)
	requireEqual(t, Uint64Literal(math.MaxUint64).Compare(Int64Literal(math.MaxInt64)), 1)
}

func TestLiteralCompareString(t *testing.T) {
	requireEqual(t, StringLiteral("a").Compare(StringLiteral("b")), -1)
	requireEqual(t, StringLiteral("b").Compare(StringLiteral("a")), 1)
	requireEqual(t, StringLiteral("a").Compare(StringLiteral("a")), 0)
}
