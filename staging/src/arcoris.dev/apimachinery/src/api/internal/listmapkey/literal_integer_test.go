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

package listmapkey

import (
	"math"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestSignedIntegerLiteral(t *testing.T) {
	gotLiteral, err := signedIntegerLiteral(conditionPath(0).Field("port"), value.Int64Value(-1))

	requireNoError(t, err)
	requireEqual(t, gotLiteral.String(), "-1")
}

func TestSignedIntegerLiteralRejectsUnsignedOverflow(t *testing.T) {
	_, err := signedIntegerLiteral(
		conditionPath(0).Field("port"),
		value.Uint64Value(math.MaxUint64),
	)

	requireErrorKind(t, err, FailureKeyIntegerRange)
}

func TestUnsignedIntegerLiteral(t *testing.T) {
	gotLiteral, err := unsignedIntegerLiteral(conditionPath(0).Field("port"), value.Uint64Value(443))

	requireNoError(t, err)
	requireEqual(t, gotLiteral.String(), "443")
}

func TestUnsignedIntegerLiteralRejectsNegativeInteger(t *testing.T) {
	_, err := unsignedIntegerLiteral(conditionPath(0).Field("port"), value.Int64Value(-1))

	requireErrorKind(t, err, FailureKeyIntegerRange)
}
