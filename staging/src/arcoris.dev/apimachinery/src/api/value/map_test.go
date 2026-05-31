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

func TestNewMapConstructsMapValue(t *testing.T) {
	value, err := NewMap(
		MapEntry("name", String("worker")),
		MapEntry("replicas", Int64(3)),
	)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindMap)
	requireEqual(t, len(value.mapValue.entries), 2)
	requireEqual(t, value.mapValue.entries[0].Key, "name")
	requireEqual(t, value.mapValue.entries[1].Key, "replicas")
}

func TestNewMapAcceptsEmptyMap(t *testing.T) {
	value, err := NewMap()
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindMap)
	requireEqual(t, len(value.mapValue.entries), 0)
}

func TestNewMapRejectsEmptyKey(t *testing.T) {
	_, err := NewMap(MapEntry("", Null()))

	requireValueError(
		t,
		err,
		ErrEmptyKey,
		mapEntryKeyPath(0),
		ErrorReasonEmptyKey,
	)

	requireErrorIs(t, err, ErrInvalidMap)
}

func TestMustMapPanicsOnMalformedEntries(t *testing.T) {
	requirePanic(t, func() {
		MustMap(Entry{Key: "", Value: Null()})
	})
}
