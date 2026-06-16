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
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// matchFieldTerm evaluates a registered selectable field. Missing paths are
// handled here; arbitrary payload traversal is not queryable unless registered.
func matchFieldTerm(t term, item objectstore.ListItem) bool {
	found := lookupFieldValue(item, t.fieldRef)
	switch t.operator {
	case OperatorExists:
		return found.present
	case OperatorDoesNotExist:
		return !found.present
	case OperatorEquals:
		return found.present && valueEqual(found.value, t.values[0])
	case OperatorNotEquals:
		return !found.present || !valueEqual(found.value, t.values[0])
	case OperatorIn:
		return found.present && valueIn(found.value, t.values)
	case OperatorNotIn:
		return !found.present || !valueIn(found.value, t.values)
	case OperatorLessThan,
		OperatorLessOrEqual,
		OperatorGreaterThan,
		OperatorGreaterOrEqual:
		return found.present && matchOrderedField(found.value, t.values[0], t.operator)
	case OperatorHasPrefix, OperatorHasSuffix, OperatorContains:
		return found.present && stringOperation(found.value, t.values[0], t.operator)
	default:
		return false
	}
}

// matchOrderedField performs runtime ordered comparison after the compiler has
// already accepted the field/operator/literal combination.
func matchOrderedField(actual value.Value, literal value.Value, op Operator) bool {
	if err := requireFieldMatchComparable(actual, literal); err != nil {
		return false
	}
	cmp, ok := compareOrdered(actual, literal)
	if !ok {
		return false
	}

	switch op {
	case OperatorLessThan:
		return cmp < 0
	case OperatorLessOrEqual:
		return cmp <= 0
	case OperatorGreaterThan:
		return cmp > 0
	case OperatorGreaterOrEqual:
		return cmp >= 0
	default:
		return false
	}
}
