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

func TestUint8PayloadCloneAndEmpty(t *testing.T) {
	payload := uint8Payload{min: limit[uint8]{value: 1, set: true}, enum: []uint8{1}}
	requireEnumPayloadCloneAndEmpty(t, enumPayloadCloneCase[uint8, uint8Payload]{
		payload:     payload,
		clone:       cloneUint8Payload,
		empty:       emptyUint8Payload,
		enum:        func(p uint8Payload) []uint8 { return p.enum },
		setEnum:     func(p *uint8Payload, values []uint8) { p.enum = values },
		wantFirst:   1,
		replaceWith: 2,
	})
}
