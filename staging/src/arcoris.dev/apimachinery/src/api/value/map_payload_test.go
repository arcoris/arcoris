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

func TestNewMapPayloadPreservesOrder(t *testing.T) {
	payload, err := newMapPayload([]Entry{
		MapEntry("first", String("one")),
		MapEntry("second", String("two")),
	})
	requireNoError(t, err)

	requireEqual(t, payload.entries[0].Key, "first")
	requireEqual(t, payload.entries[1].Key, "second")
}

func TestMapPayloadCompactRemovesEmptyStorage(t *testing.T) {
	payload, err := newMapPayload(nil)
	requireNoError(t, err)

	requireEqual(t, payload.entries == nil, true)
}
