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

package retrybudget

import "testing"

func TestPolicySnapshotIsValid(t *testing.T) {
	tests := []struct {
		name string
		val  PolicySnapshot
		want bool
	}{
		{name: "unbounded", val: PolicySnapshot{}, want: true},
		{name: "bounded zero", val: PolicySnapshot{Ratio: RatioZero, Bounded: true}, want: true},
		{name: "bounded one", val: PolicySnapshot{Ratio: RatioOne, Minimum: 10, Bounded: true}, want: true},
		{name: "bounded fraction", val: PolicySnapshot{Ratio: MustRatio(1, 5), Minimum: 3, Bounded: true}, want: true},
		{name: "bounded unset ratio", val: PolicySnapshot{Bounded: true}, want: false},
		{name: "unbounded with ratio", val: PolicySnapshot{Ratio: MustRatio(1, 5)}, want: false},
		{name: "unbounded with minimum", val: PolicySnapshot{Minimum: 1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicySnapshotHasMinimum(t *testing.T) {
	if !(PolicySnapshot{Ratio: MustRatio(1, 5), Minimum: 1, Bounded: true}).HasMinimum() {
		t.Fatal("HasMinimum returned false")
	}
	if (PolicySnapshot{Minimum: 1}).HasMinimum() {
		t.Fatal("HasMinimum returned true for unbounded policy")
	}
}

func TestPolicySnapshotIsBounded(t *testing.T) {
	tests := []struct {
		name string
		val  PolicySnapshot
		want bool
	}{
		{name: "unbounded", val: PolicySnapshot{}, want: false},
		{name: "bounded", val: PolicySnapshot{Ratio: MustRatio(1, 5), Minimum: 1, Bounded: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsBounded(); got != tt.want {
				t.Fatalf("IsBounded() = %v, want %v", got, tt.want)
			}
		})
	}
}
