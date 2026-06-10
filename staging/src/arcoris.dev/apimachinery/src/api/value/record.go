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

// RecordValue constructs a KindRecord Value from uniquely named members.
//
// Record values are concrete keyed payload nodes with caller-ordered members.
// They intentionally do not decide whether the data will satisfy a descriptor
// object or descriptor map. Member values are cloned during construction so
// later mutations of caller-owned values cannot affect the stored record.
func RecordValue(members ...RecordMember) (Value, error) {
	payload, err := newRecordPayload(members)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindRecord, recordValue: payload}, nil
}

// MustRecordValue constructs a record Value or panics when members are malformed.
//
// It is intended for tests and static fixtures where malformed member data is a
// programmer error. Runtime construction paths should use RecordValue and return
// its structured error.
func MustRecordValue(members ...RecordMember) Value {
	value, err := RecordValue(members...)
	if err != nil {
		panic(err)
	}

	return value
}
