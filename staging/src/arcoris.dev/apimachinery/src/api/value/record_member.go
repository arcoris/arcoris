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

// RecordMember is one concrete record member in declaration order.
//
// RecordMember is payload data, not a descriptor field. Name is the actual key
// present in the record, and Value is the actual nested payload value. Empty
// names and invalid zero Values are rejected by RecordValue.
type RecordMember struct {
	// Name is the concrete record member name.
	Name MemberName
	// Value is the concrete member payload.
	Value Value
}

// NewRecordMember constructs one lightweight record member tuple.
//
// The supplied Value is not deep-cloned here. RecordValue owns the construction
// boundary and clones nested payloads exactly once.
func NewRecordMember(name MemberName, value Value) RecordMember {
	return RecordMember{Name: name, Value: value}
}

// MustRecordMember validates name and constructs one record member.
//
// It is intended for tests and static fixtures. Runtime construction paths that
// receive untrusted member names should use NewMemberName and RecordValue.
func MustRecordMember(name string, value Value) RecordMember {
	return NewRecordMember(MustMemberName(name), value)
}
