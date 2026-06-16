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

import "testing"

// TestOperatorString verifies stable diagnostic names for every public
// operator.
func TestOperatorString(t *testing.T) {
	tests := map[Operator]string{
		0:                      "unknown",
		OperatorExists:         "exists",
		OperatorDoesNotExist:   "doesNotExist",
		OperatorEquals:         "equals",
		OperatorNotEquals:      "notEquals",
		OperatorIn:             "in",
		OperatorNotIn:          "notIn",
		OperatorLessThan:       "lessThan",
		OperatorLessOrEqual:    "lessOrEqual",
		OperatorGreaterThan:    "greaterThan",
		OperatorGreaterOrEqual: "greaterOrEqual",
		OperatorHasPrefix:      "hasPrefix",
		OperatorHasSuffix:      "hasSuffix",
		OperatorContains:       "contains",
		Operator(99):           "unknown",
	}

	for op, want := range tests {
		if got := op.String(); got != want {
			t.Fatalf("%v.String() = %q; want %q", uint8(op), got, want)
		}
	}
}
