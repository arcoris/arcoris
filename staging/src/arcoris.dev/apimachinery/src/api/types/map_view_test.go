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

func TestMapViewAccessorsDetach(t *testing.T) {
	desc := MapOf(String().Enum("alpha")).Keys(String().Enum("key")).Descriptor()
	view, ok := desc.AsMap()
	requireEqual(t, ok, true)

	key := view.Key()
	key.string.enum[0] = "changed"
	keyAgain := view.Key()
	requireEqual(t, keyAgain.string.enum[0], "key")

	value := view.Value()
	value.string.enum[0] = "changed"
	valueAgain := view.Value()
	requireEqual(t, valueAgain.string.enum[0], "alpha")
}

func TestMapViewKeyAndValueReturnDetachedDescriptors(t *testing.T) {
	desc := MapOf(String().MinBytes(1)).
		Keys(String().MaxBytes(63)).
		Descriptor()
	view := requireMapView(t, desc)

	key := view.Key()
	key.string.maxBytes.value = 1
	value := view.Value()
	value.string.minBytes.value = 99

	originalKey := requireStringView(t, view.Key())
	maxBytes, ok := originalKey.MaxBytes()
	requireEqual(t, ok, true)
	requireEqual(t, maxBytes, 63)

	originalValue := requireStringView(t, view.Value())
	minBytes, ok := originalValue.MinBytes()
	requireEqual(t, ok, true)
	requireEqual(t, minBytes, 1)
}

func TestMapViewZeroValueValue(t *testing.T) {
	var view MapView
	requireEqual(t, view.Key().Code(), DescriptorInvalid)
	requireEqual(t, view.Value().Code(), DescriptorInvalid)
}
