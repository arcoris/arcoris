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

// Operator is the finite metadata requirement operator set.
type Operator uint8

const (
	// OperatorExists matches when a metadata key is present.
	OperatorExists Operator = iota + 1

	// OperatorDoesNotExist matches when a metadata key is absent.
	OperatorDoesNotExist

	// OperatorEquals matches when a metadata key has exactly one value.
	OperatorEquals

	// OperatorNotEquals matches when a metadata key is absent or differs.
	OperatorNotEquals

	// OperatorIn matches when a metadata key value is in a finite set.
	OperatorIn

	// OperatorNotIn matches when a metadata key is absent or outside a finite set.
	OperatorNotIn
)

// IsValid reports whether op is one of the known query operators.
func (op Operator) IsValid() bool {
	switch op {
	case OperatorExists,
		OperatorDoesNotExist,
		OperatorEquals,
		OperatorNotEquals,
		OperatorIn,
		OperatorNotIn:
		return true
	default:
		return false
	}
}

// String returns the stable operator spelling.
func (op Operator) String() string {
	switch op {
	case OperatorExists:
		return "exists"
	case OperatorDoesNotExist:
		return "doesNotExist"
	case OperatorEquals:
		return "equals"
	case OperatorNotEquals:
		return "notEquals"
	case OperatorIn:
		return "in"
	case OperatorNotIn:
		return "notIn"
	default:
		return "unknown"
	}
}
