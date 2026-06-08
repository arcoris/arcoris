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

package bulkhead

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestRefusalMatchesCapacityRefusal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		bulkhead Refusal
		capacity capacity.Refusal
	}{
		{name: "none", bulkhead: RefusalNone, capacity: capacity.RefusalNone},
		{name: "insufficient", bulkhead: RefusalInsufficient, capacity: capacity.RefusalInsufficient},
		{name: "debt", bulkhead: RefusalDebt, capacity: capacity.RefusalDebt},
		{name: "unknown resource", bulkhead: RefusalUnknownResource, capacity: capacity.RefusalUnknownResource},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.bulkhead != tt.capacity {
				t.Fatalf("bulkhead refusal = %s, want capacity refusal %s", tt.bulkhead, tt.capacity)
			}
		})
	}
}
