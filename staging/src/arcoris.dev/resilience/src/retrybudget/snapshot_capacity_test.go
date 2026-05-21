/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retrybudget

import (
	"math"
	"testing"
)

func TestCapacitySnapshotIsValid(t *testing.T) {
	tests := []struct {
		name string
		val  CapacitySnapshot
		want bool
	}{
		{name: "available", val: CapacitySnapshot{Allowed: 4, Available: 2, Exhausted: false}, want: true},
		{name: "exhausted", val: CapacitySnapshot{Allowed: 4, Available: 0, Exhausted: true}, want: true},
		{name: "unlimited", val: maxedCapacity(), want: true},
		{name: "available greater than allowed", val: CapacitySnapshot{Allowed: 1, Available: 2}, want: false},
		{name: "exhausted with available", val: CapacitySnapshot{Allowed: 4, Available: 1, Exhausted: true}, want: false},
		{name: "not exhausted with zero available", val: CapacitySnapshot{Allowed: 4, Available: 0, Exhausted: false}, want: false},
		{name: "zero", val: CapacitySnapshot{}, want: false},
		{name: "max unavailable exhausted", val: CapacitySnapshot{Allowed: math.MaxUint64, Available: 0, Exhausted: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapacitySnapshotHasAvailable(t *testing.T) {
	if !maxedCapacity().HasAvailable() {
		t.Fatal("HasAvailable returned false for max capacity")
	}
	if (CapacitySnapshot{Allowed: 1, Available: 0, Exhausted: true}).HasAvailable() {
		t.Fatal("HasAvailable returned true for exhausted capacity")
	}
}
