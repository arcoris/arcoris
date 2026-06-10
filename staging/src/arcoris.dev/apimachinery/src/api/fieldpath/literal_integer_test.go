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

func TestIntegerStoresMinInt64(t *testing.T) {
	value := newInt64(math.MinInt64)
	got, ok := value.int64Value()

	requireEqual(t, ok, true)
	requireEqual(t, got, int64(math.MinInt64))
	requireEqual(t, value.string(), "-9223372036854775808")
}

func TestIntegerRejectsNegativeZeroRepresentation(t *testing.T) {
	err := (integer{negative: true}).validate()

	requireErrorIs(t, err, ErrInvalidLiteral)
}

func TestIntegerCompareAcrossSignedAndUnsignedDomains(t *testing.T) {
	requireEqual(t, newInt64(-1).compare(newUint64(0)), -1)
	requireEqual(t, newUint64(math.MaxUint64).compare(newInt64(math.MaxInt64)), 1)
	requireEqual(t, newUint64(1).compare(newInt64(1)), 0)
}
