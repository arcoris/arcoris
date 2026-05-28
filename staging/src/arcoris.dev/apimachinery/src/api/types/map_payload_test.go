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

func TestMapPayloadStoresStringKeyAndValue(t *testing.T) {
	payload := MapOf(String()).Type().mapType
	requireEqual(t, payload.key, MapKeyString)
	requireEqual(t, payload.value.Code(), TypeString)
}

func TestMapPayloadCloneAndEmpty(t *testing.T) {
	payload := MapOf(String().Enum("a")).Type().mapType
	cloned := cloneMapPayload(payload)
	cloned.value.string.enum[0] = "b"

	requireEqual(t, payload.value.string.enum[0], "a")
	requireEqual(t, emptyMapPayload(mapPayload{}), true)
	requireEqual(t, emptyMapPayload(payload), false)
}
