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

import "testing"

func TestSnapshotIsValid(t *testing.T) {
	valid := validSnapshotValue()
	invalidKind := valid
	invalidKind.Kind = KindUnknown
	invalidCapacity := valid
	invalidCapacity.Capacity.Available = invalidCapacity.Capacity.Allowed + 1
	invalidWindow := valid
	invalidWindow.Window.Duration = 0
	invalidPolicy := valid
	invalidPolicy.Policy.Ratio = 2
	tests := []struct {
		name string
		val  Snapshot
		want bool
	}{
		{name: "valid", val: valid, want: true},
		{name: "invalid kind", val: invalidKind, want: false},
		{name: "invalid capacity", val: invalidCapacity, want: false},
		{name: "invalid window", val: invalidWindow, want: false},
		{name: "invalid policy", val: invalidPolicy, want: false},
		{name: "zero", val: Snapshot{}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnapshotHelpers(t *testing.T) {
	val := validSnapshotValue()
	if val.Exhausted() {
		t.Fatal("Exhausted returned true for available capacity")
	}
	if !val.HasTraffic() {
		t.Fatal("HasTraffic returned false")
	}
	val.Capacity = CapacitySnapshot{Allowed: 4, Available: 0, Exhausted: true}
	if !val.Exhausted() {
		t.Fatal("Exhausted returned false")
	}
}
