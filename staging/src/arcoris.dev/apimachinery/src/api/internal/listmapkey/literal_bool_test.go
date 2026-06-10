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
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestBoolLiteral(t *testing.T) {
	gotLiteral, err := boolLiteral(conditionPath(0).Field(testFieldName("enabled")), value.BoolValue(true))

	requireNoError(t, err)
	requireEqual(t, gotLiteral.String(), "true")
}

func TestBoolLiteralRejectsWrongKind(t *testing.T) {
	_, err := boolLiteral(conditionPath(0).Field(testFieldName("enabled")), value.StringValue("true"))

	requireErrorKind(t, err, FailureKeyKindMismatch)
}
