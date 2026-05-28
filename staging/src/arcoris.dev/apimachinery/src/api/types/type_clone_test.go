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

package types

import "testing"

func TestCloneHelpersDetachSlices(t *testing.T) {
	strings := cloneStrings([]string{"a"})
	strings[0] = "b"
	requireEqual(t, cloneStrings([]string{"a"})[0], "a")
	requireEqual(t, cloneInt8s([]int8{1})[0], int8(1))
	requireEqual(t, cloneInt16s([]int16{1})[0], int16(1))
	requireEqual(t, cloneInt32s([]int32{1})[0], int32(1))
	requireEqual(t, cloneInt64s([]int64{1})[0], int64(1))
	requireEqual(t, cloneUint8s([]uint8{1})[0], uint8(1))
	requireEqual(t, cloneUint16s([]uint16{1})[0], uint16(1))
	requireEqual(t, cloneUint32s([]uint32{1})[0], uint32(1))
	requireEqual(t, cloneUint64s([]uint64{1})[0], uint64(1))
	requireEqual(t, cloneFloat32s([]float32{1})[0], float32(1))
	requireEqual(t, cloneFloat64s([]float64{1})[0], float64(1))
	requireEqual(t, cloneFieldNames([]FieldName{"name"})[0], FieldName("name"))
	requireEqual(t, cloneFields([]FieldDescriptor{Field("name").String().Required().Field()})[0].Name(), FieldName("name"))
}
