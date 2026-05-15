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

func TestPolicySnapshotIsValid(t *testing.T) {
	tests := []struct {
		name string
		val  PolicySnapshot
		want bool
	}{
		{name: "unbounded", val: PolicySnapshot{}, want: true},
		{name: "bounded zero", val: PolicySnapshot{Ratio: 0, Bounded: true}, want: true},
		{name: "bounded one", val: PolicySnapshot{Ratio: 1, Minimum: 10, Bounded: true}, want: true},
		{name: "bounded fraction", val: PolicySnapshot{Ratio: 0.2, Minimum: 3, Bounded: true}, want: true},
		{name: "negative", val: PolicySnapshot{Ratio: -0.1, Bounded: true}, want: false},
		{name: "greater than one", val: PolicySnapshot{Ratio: 1.1, Bounded: true}, want: false},
		{name: "nan", val: PolicySnapshot{Ratio: math.NaN(), Bounded: true}, want: false},
		{name: "inf", val: PolicySnapshot{Ratio: math.Inf(1), Bounded: true}, want: false},
		{name: "unbounded with ratio", val: PolicySnapshot{Ratio: 0.2}, want: false},
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
	if !(PolicySnapshot{Ratio: 0.2, Minimum: 1, Bounded: true}).HasMinimum() {
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
		{name: "bounded", val: PolicySnapshot{Ratio: 0.2, Minimum: 1, Bounded: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsBounded(); got != tt.want {
				t.Fatalf("IsBounded() = %v, want %v", got, tt.want)
			}
		})
	}
}
