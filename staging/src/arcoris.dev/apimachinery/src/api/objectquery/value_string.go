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

package objectquery

import (
	"strings"

	"arcoris.dev/apimachinery/api/value"
)

// stringOperation applies string-only field operators after compile-time field
// validation has accepted them.
func stringOperation(actual value.Value, literal value.Value, op Operator) bool {
	actualString, ok := actual.AsString()
	if !ok {
		return false
	}
	literalString, ok := literal.AsString()
	if !ok {
		return false
	}

	switch op {
	case OperatorHasPrefix:
		return strings.HasPrefix(actualString, literalString)
	case OperatorHasSuffix:
		return strings.HasSuffix(actualString, literalString)
	case OperatorContains:
		return strings.Contains(actualString, literalString)
	default:
		return false
	}
}
