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

func TestUint16ViewAccessors(t *testing.T) {
	view := requireUint16View(t, Uint16().Range(20, 90).Enum(30, 70).Descriptor())

	min, ok := view.Min()
	requireEqual(t, ok, true)
	requireEqual(t, min, uint16(20))

	max, ok := view.Max()
	requireEqual(t, ok, true)
	requireEqual(t, max, uint16(90))
}

func TestUint16ViewEnumReturnsDetachedData(t *testing.T) {
	view := requireUint16View(t, Uint16().Enum(30, 70).Descriptor())

	enum := view.Enum()
	enum[0] = 990

	requireEqual(t, view.Enum()[0], uint16(30))
}
