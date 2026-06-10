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

func TestListPayloadPreservesOrder(t *testing.T) {
	payload, err := newListPayload([]Value{
		StringValue("first"),
		StringValue("second"),
	})
	requireNoError(t, err)

	first, ok := payload.items[0].AsString()
	requireEqual(t, ok, true)
	requireEqual(t, first, "first")

	second, ok := payload.items[1].AsString()
	requireEqual(t, ok, true)
	requireEqual(t, second, "second")
}

func TestListPayloadEmptyListUsesNilStorage(t *testing.T) {
	payload, err := newListPayload(nil)
	requireNoError(t, err)

	requireEqual(t, payload.items == nil, true)
}
