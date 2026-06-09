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

func TestUint64ViewAccessors(t *testing.T) {
	view := requireUint64View(t, Uint64().Range(2000, 9000).Enum(3000, 7000).Descriptor())

	min, ok := view.Min()
	requireEqual(t, ok, true)
	requireEqual(t, min, uint64(2000))

	max, ok := view.Max()
	requireEqual(t, ok, true)
	requireEqual(t, max, uint64(9000))
}

func TestUint64ViewEnumReturnsDetachedData(t *testing.T) {
	view := requireUint64View(t, Uint64().Enum(3000, 7000).Descriptor())

	enum := view.Enum()
	enum[0] = 99000

	requireEqual(t, view.Enum()[0], uint64(3000))
}
