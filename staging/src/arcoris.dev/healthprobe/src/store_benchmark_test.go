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

package probe

import (
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

func BenchmarkStoreUpdate(b *testing.B) {
	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	report := healthtest.HealthyReport(health.TargetReady)

	b.ReportAllocs()
	for b.Loop() {
		_ = s.update(health.TargetReady, report)
	}
}

func BenchmarkStoreRejectedUpdate(b *testing.B) {
	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	report := health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusHealthy,
		Observed: testNow,
	}

	b.ReportAllocs()
	for b.Loop() {
		_ = s.update(health.TargetReady, report)
	}
}

func BenchmarkStoreSnapshot(b *testing.B) {
	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	_ = s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady))

	b.ReportAllocs()
	for b.Loop() {
		_, _ = s.snapshot(health.TargetReady)
	}
}

func BenchmarkStoreSnapshots(b *testing.B) {
	s := newStore(
		[]health.Target{health.TargetStartup, health.TargetLive, health.TargetReady},
		newTestClock(),
	)
	_ = s.update(health.TargetStartup, healthtest.HealthyReport(health.TargetStartup))
	_ = s.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive))
	_ = s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady))

	b.ReportAllocs()
	for b.Loop() {
		_ = s.snapshots()
	}
}
