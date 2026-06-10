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

import "testing"

func TestCloneDeepCopiesBytes(t *testing.T) {
	clone := BytesValue([]byte{1, 2, 3}).Clone()
	bytes, ok := clone.AsBytes()
	requireEqual(t, ok, true)

	bytes[0] = 9
	again, ok := clone.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, again, []byte{1, 2, 3})
}

func TestCloneDeepCopiesRecordMembers(t *testing.T) {
	original := mustRecord(t, MustRecordMember("payload", BytesValue([]byte{1, 2})))
	clone := original.Clone()

	members := clone.recordValue.members
	members[0].Value.bytesValue[0] = 9

	view, ok := original.AsRecord()
	requireEqual(t, ok, true)

	value, ok := view.Get("payload")
	requireEqual(t, ok, true)

	bytes, ok := value.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestCloneDeepCopiesListItems(t *testing.T) {
	original := mustList(t, BytesValue([]byte{1, 2}))
	clone := original.Clone()

	clone.listValue.items[0].bytesValue[0] = 9

	view, ok := original.AsList()
	requireEqual(t, ok, true)

	value, ok := view.At(0)
	requireEqual(t, ok, true)

	bytes, ok := value.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}
