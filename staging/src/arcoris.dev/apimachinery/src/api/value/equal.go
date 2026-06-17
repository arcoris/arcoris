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

package value

import "bytes"

// Equal reports whether left and right are the same concrete payload value.
//
// Equal is descriptor-neutral. It compares concrete Value shape and payload
// data only. Descriptor-aware object-vs-map semantics belong to valuecompare.
// Decimal equality is numeric, so 1.20 and 1.2 are equal. Record and list
// equality preserves concrete member/item order.
func Equal(left Value, right Value) bool {
	if left.Kind() != right.Kind() {
		return false
	}

	switch left.Kind() {
	case KindInvalid:
		return true
	case KindNull:
		return true
	case KindBool:
		l, _ := left.AsBool()
		r, _ := right.AsBool()
		return l == r
	case KindString:
		l, _ := left.AsString()
		r, _ := right.AsString()
		return l == r
	case KindBytes:
		l, _ := left.AsBytes()
		r, _ := right.AsBytes()
		return bytes.Equal(l, r)
	case KindInteger:
		l, _ := left.AsInteger()
		r, _ := right.AsInteger()
		return l.Equal(r)
	case KindFloat:
		l, _ := left.AsFloat()
		r, _ := right.AsFloat()
		return l == r
	case KindDecimal:
		l, _ := left.AsDecimal()
		r, _ := right.AsDecimal()
		return l.Equal(r)
	case KindTimestamp:
		l, _ := left.AsTimestamp()
		r, _ := right.AsTimestamp()
		return l.Equal(r)
	case KindDate:
		l, _ := left.AsDate()
		r, _ := right.AsDate()
		return l.Equal(r)
	case KindTimeOfDay:
		l, _ := left.AsTimeOfDay()
		r, _ := right.AsTimeOfDay()
		return l.Equal(r)
	case KindDuration:
		l, _ := left.AsDuration()
		r, _ := right.AsDuration()
		return l == r
	case KindRecord:
		return equalRecords(left, right)
	case KindList:
		return equalLists(left, right)
	default:
		return false
	}
}

// equalRecords compares concrete record payloads in preserved member order.
func equalRecords(left Value, right Value) bool {
	l, _ := left.AsRecord()
	r, _ := right.AsRecord()
	if l.Len() != r.Len() {
		return false
	}
	for i := 0; i < l.Len(); i++ {
		lm, _ := l.Member(i)
		rm, _ := r.Member(i)
		if lm.Name != rm.Name || !Equal(lm.Value, rm.Value) {
			return false
		}
	}

	return true
}

// equalLists compares concrete list payloads in preserved item order.
func equalLists(left Value, right Value) bool {
	l, _ := left.AsList()
	r, _ := right.AsList()
	if l.Len() != r.Len() {
		return false
	}
	for i := 0; i < l.Len(); i++ {
		lv, _ := l.At(i)
		rv, _ := r.At(i)
		if !Equal(lv, rv) {
			return false
		}
	}

	return true
}
