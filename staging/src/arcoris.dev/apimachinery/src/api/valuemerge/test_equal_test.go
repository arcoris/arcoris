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

package valuemerge

import "arcoris.dev/apimachinery/api/value"

func valuesEqual(left value.Value, right value.Value) bool {
	if left.Kind() != right.Kind() {
		return false
	}

	switch left.Kind() {
	case value.KindInvalid, value.KindNull:
		return true
	case value.KindBool:
		l, _ := left.AsBool()
		r, _ := right.AsBool()
		return l == r
	case value.KindString:
		l, _ := left.AsString()
		r, _ := right.AsString()
		return l == r
	case value.KindInteger:
		l, _ := left.AsInteger()
		r, _ := right.AsInteger()
		return l.String() == r.String()
	case value.KindRecord:
		return objectsEqual(left, right)
	case value.KindList:
		return listsEqual(left, right)
	default:
		return false
	}
}

func objectsEqual(left value.Value, right value.Value) bool {
	leftView, _ := left.AsRecord()
	rightView, _ := right.AsRecord()
	leftMembers := leftView.Members()
	rightMembers := rightView.Members()
	if len(leftMembers) != len(rightMembers) {
		return false
	}

	for i := range leftMembers {
		if leftMembers[i].Name != rightMembers[i].Name {
			return false
		}
		if !valuesEqual(leftMembers[i].Value, rightMembers[i].Value) {
			return false
		}
	}

	return true
}

func listsEqual(left value.Value, right value.Value) bool {
	leftView, _ := left.AsList()
	rightView, _ := right.AsList()
	leftItems := leftView.Items()
	rightItems := rightView.Items()
	if len(leftItems) != len(rightItems) {
		return false
	}

	for i := range leftItems {
		if !valuesEqual(leftItems[i], rightItems[i]) {
			return false
		}
	}

	return true
}
