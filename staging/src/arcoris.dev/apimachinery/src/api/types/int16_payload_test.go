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

func TestInt16PayloadCloneAndEmpty(t *testing.T) {
	payload := int16Payload{min: limit[int16]{value: 1, set: true}, enum: []int16{1}}
	requireEnumPayloadCloneAndEmpty(t, enumPayloadCloneCase[int16, int16Payload]{
		payload:     payload,
		clone:       cloneInt16Payload,
		empty:       emptyInt16Payload,
		enum:        func(p int16Payload) []int16 { return p.enum },
		setEnum:     func(p *int16Payload, values []int16) { p.enum = values },
		wantFirst:   1,
		replaceWith: 2,
	})
}
