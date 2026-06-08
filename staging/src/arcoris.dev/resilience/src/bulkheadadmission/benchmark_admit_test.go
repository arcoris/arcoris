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

package bulkheadadmission

import (
	"testing"

	"arcoris.dev/resilience/bulkhead"
)

func BenchmarkAdmitterTryAdmitGrantedAndRelease(b *testing.B) {
	b.ReportAllocs()

	core := bulkhead.New(1)
	admitter := New(core)
	req := Request{Amount: 1}

	for b.Loop() {
		result := admitter.TryAdmit(req)
		lease, ok := result.Grant()
		if !ok || lease == nil {
			b.Fatal("TryAdmit returned no grant")
		}
		lease.Release()
		benchmarkResult = result
	}
}

func BenchmarkAdmitterTryAdmitDenied(b *testing.B) {
	b.ReportAllocs()

	admitter := New(bulkhead.New(0))
	req := Request{Amount: 1}

	for b.Loop() {
		result := admitter.TryAdmit(req)
		if !result.Decision().IsDenied() {
			b.Fatal("TryAdmit succeeded, want denied")
		}
		benchmarkResult = result
	}
}

func BenchmarkAdmitterTryAdmitDebtDenied(b *testing.B) {
	b.ReportAllocs()

	core := bulkhead.New(2)
	held, _, ok := core.TryAcquireAmount(2)
	if !ok {
		b.Fatal("initial TryAcquireAmount failed")
	}
	defer held.Release()
	core.SetLimit(1)

	admitter := New(core)
	req := Request{Amount: 1}

	for b.Loop() {
		result := admitter.TryAdmit(req)
		if !result.Decision().IsDenied() {
			b.Fatal("TryAdmit succeeded, want denied")
		}
		benchmarkResult = result
	}
}

func BenchmarkAdmitterTryAdmitGrantedAndReleaseParallel(b *testing.B) {
	b.ReportAllocs()

	core := bulkhead.New(1024)
	admitter := New(core)
	req := Request{Amount: 1}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result := admitter.TryAdmit(req)
			lease, ok := result.Grant()
			if ok && lease != nil {
				lease.Release()
			}
			benchmarkResult = result
		}
	})
}
