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

func TestNewObjectConstructsObjectValue(t *testing.T) {
	value, err := NewObject(
		ObjectField("name", String("worker")),
		ObjectField("replicas", Int64(3)),
	)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindObject)
	requireEqual(t, len(value.objectValue.fields), 2)
	requireEqual(t, value.objectValue.fields[0].Name, "name")
	requireEqual(t, value.objectValue.fields[1].Name, "replicas")
}

func TestNewObjectAcceptsEmptyObject(t *testing.T) {
	value, err := NewObject()
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindObject)
	requireEqual(t, len(value.objectValue.fields), 0)
}

func TestMustObjectPanicsOnMalformedFields(t *testing.T) {
	requirePanic(t, func() {
		MustObject(Field{Name: "", Value: Null()})
	})
}
