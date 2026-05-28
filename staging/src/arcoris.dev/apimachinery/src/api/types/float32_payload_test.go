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

func TestFloat32PayloadCloneAndEmpty(t *testing.T) {
	payload := float32Payload{min: limit[float32]{value: 1, set: true}, enum: []float32{1}}
	requireEnumPayloadCloneAndEmpty(t, enumPayloadCloneCase[float32, float32Payload]{
		payload:     payload,
		clone:       cloneFloat32Payload,
		empty:       emptyFloat32Payload,
		enum:        func(p float32Payload) []float32 { return p.enum },
		setEnum:     func(p *float32Payload, values []float32) { p.enum = values },
		wantFirst:   1,
		replaceWith: 2,
	})
}
