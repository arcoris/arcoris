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
	"testing"
	"time"
)

func TestWindowSnapshotIsValid(t *testing.T) {
	start := time.Unix(100, 0).UTC()
	valid := WindowSnapshot{StartedAt: start, EndsAt: start.Add(time.Minute), Duration: time.Minute, Bounded: true}
	tests := []struct {
		name string
		val  WindowSnapshot
		want bool
	}{
		{name: "unbounded", val: WindowSnapshot{}, want: true},
		{name: "bounded", val: valid, want: true},
		{name: "bounded zero duration", val: WindowSnapshot{StartedAt: start, EndsAt: start, Bounded: true}, want: false},
		{name: "bounded zero start", val: WindowSnapshot{EndsAt: start.Add(time.Minute), Duration: time.Minute, Bounded: true}, want: false},
		{name: "bounded inconsistent end", val: WindowSnapshot{StartedAt: start, EndsAt: start.Add(2 * time.Minute), Duration: time.Minute, Bounded: true}, want: false},
		{name: "unbounded with duration", val: WindowSnapshot{Duration: time.Minute}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindowSnapshotContains(t *testing.T) {
	start := time.Unix(100, 0).UTC()
	win := WindowSnapshot{StartedAt: start, EndsAt: start.Add(time.Minute), Duration: time.Minute, Bounded: true}
	if !win.Contains(start) {
		t.Fatal("Contains rejected start boundary")
	}
	if !win.Contains(start.Add(30 * time.Second)) {
		t.Fatal("Contains rejected inner time")
	}
	if win.Contains(start.Add(time.Minute)) {
		t.Fatal("Contains accepted exclusive end boundary")
	}
	if (WindowSnapshot{}).Contains(start) {
		t.Fatal("Contains accepted time for unbounded window")
	}
}

func TestWindowSnapshotIsBounded(t *testing.T) {
	start := time.Unix(100, 0).UTC()
	tests := []struct {
		name string
		val  WindowSnapshot
		want bool
	}{
		{name: "unbounded", val: WindowSnapshot{}, want: false},
		{name: "bounded", val: WindowSnapshot{StartedAt: start, EndsAt: start.Add(time.Minute), Duration: time.Minute, Bounded: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsBounded(); got != tt.want {
				t.Fatalf("IsBounded() = %v, want %v", got, tt.want)
			}
		})
	}
}
