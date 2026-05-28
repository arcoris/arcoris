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

func TestInt64PayloadCloneAndEmpty(t *testing.T) {
	payload := int64Payload{min: limit[int64]{value: 1, set: true}, enum: []int64{1}}
	requireEnumPayloadCloneAndEmpty(t, enumPayloadCloneCase[int64, int64Payload]{
		payload:     payload,
		clone:       cloneInt64Payload,
		empty:       emptyInt64Payload,
		enum:        func(p int64Payload) []int64 { return p.enum },
		setEnum:     func(p *int64Payload, values []int64) { p.enum = values },
		wantFirst:   1,
		replaceWith: 2,
	})
}
