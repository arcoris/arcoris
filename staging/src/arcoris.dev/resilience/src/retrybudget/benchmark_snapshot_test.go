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

var benchmarkTotal uint64

func BenchmarkSnapshotIsValidFixedWindow(b *testing.B) {
	value := validSnapshotValue()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkValid = value.IsValid()
	}
}

func BenchmarkSnapshotIsValidNoop(b *testing.B) {
	value := Snapshot{
		Kind:     KindNoop,
		Attempts: AttemptsSnapshot{},
		Capacity: maxedCapacity(),
		Window:   WindowSnapshot{Bounded: false},
		Policy:   PolicySnapshot{Bounded: false},
	}
	b.ReportAllocs()
	for b.Loop() {
		benchmarkValid = value.IsValid()
	}
}

func BenchmarkAttemptsTotal(b *testing.B) {
	attempts := AttemptsSnapshot{Original: 123, Retry: 45}
	b.ReportAllocs()
	for b.Loop() {
		benchmarkTotal = attempts.Total()
	}
}

func BenchmarkCapacityIsValid(b *testing.B) {
	capacity := CapacitySnapshot{Allowed: 10, Available: 4, Exhausted: false}
	b.ReportAllocs()
	for b.Loop() {
		benchmarkValid = capacity.IsValid()
	}
}

func BenchmarkWindowContains(b *testing.B) {
	window := validSnapshotValue().Window
	now := window.StartedAt.Add(window.Duration / 2)
	b.ReportAllocs()
	for b.Loop() {
		benchmarkValid = window.Contains(now)
	}
}

func BenchmarkPolicyIsValid(b *testing.B) {
	policy := PolicySnapshot{Ratio: MustRatio(1, 5), Minimum: 10, Bounded: true}
	b.ReportAllocs()
	for b.Loop() {
		benchmarkValid = policy.IsValid()
	}
}
