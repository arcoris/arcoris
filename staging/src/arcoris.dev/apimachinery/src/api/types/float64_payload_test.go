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

func TestFloat64PayloadCloneAndEmpty(t *testing.T) {
	payload := float64Payload{min: limit[float64]{value: 1, set: true}, enum: []float64{1}}
	requireEnumPayloadCloneAndEmpty(t, enumPayloadCloneCase[float64, float64Payload]{
		payload:     payload,
		clone:       cloneFloat64Payload,
		empty:       emptyFloat64Payload,
		enum:        func(p float64Payload) []float64 { return p.enum },
		setEnum:     func(p *float64Payload, values []float64) { p.enum = values },
		wantFirst:   1,
		replaceWith: 2,
	})
}
