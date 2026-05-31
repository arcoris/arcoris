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

func TestNewListConstructsListValue(t *testing.T) {
	value, err := NewList(String("worker"), Int64(3))
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindList)
	requireEqual(t, len(value.listValue.items), 2)
}

func TestNewListAcceptsEmptyList(t *testing.T) {
	value, err := NewList()
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindList)
	requireEqual(t, len(value.listValue.items), 0)
	requireEqual(t, value.listValue.items == nil, true)
}

func TestMustListPanicsOnMalformedItems(t *testing.T) {
	requirePanic(t, func() {
		MustList(Value{})
	})
}
