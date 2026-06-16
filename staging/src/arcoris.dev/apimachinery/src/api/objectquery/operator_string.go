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

// String returns a stable diagnostic operator name.
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
	case OperatorLessThan:
		return "lessThan"
	case OperatorLessOrEqual:
		return "lessOrEqual"
	case OperatorGreaterThan:
		return "greaterThan"
	case OperatorGreaterOrEqual:
		return "greaterOrEqual"
	case OperatorHasPrefix:
		return "hasPrefix"
	case OperatorHasSuffix:
		return "hasSuffix"
	case OperatorContains:
		return "contains"
	default:
		return "unknown"
	}
}
