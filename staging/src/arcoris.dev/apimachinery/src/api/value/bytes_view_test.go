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

func TestBytesAccessorReturnsClone(t *testing.T) {
	value := Bytes([]byte{1, 2, 3})

	got, ok := value.Bytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, got, []byte{1, 2, 3})

	got[1] = 9
	again, ok := value.Bytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, again, []byte{1, 2, 3})
}

func TestBytesAccessorPreservesEmptySlice(t *testing.T) {
	got, ok := Bytes(nil).Bytes()

	requireEqual(t, ok, true)
	if got == nil {
		t.Fatal("Bytes(nil) accessor returned nil slice")
	}
	requireEqual(t, len(got), 0)
}
