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

package health

import (
	"context"
	"testing"
)

func BenchmarkValidCheckName(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ValidCheckName("database_pool_1")
	}
}

func BenchmarkReasonIsValid(b *testing.B) {
	reason := Reason("dependency_unavailable")
	b.ReportAllocs()
	for b.Loop() {
		_ = reason.IsValid()
	}
}

func BenchmarkRegistryChecks(b *testing.B) {
	registry := NewRegistry()
	for _, name := range []string{"storage", "queue"} {
		checkName := name
		checker, err := NewCheck(checkName, func(context.Context) Result {
			return Healthy(checkName)
		})
		if err != nil {
			b.Fatalf("NewCheck() = %v", err)
		}
		if err := registry.Register(TargetReady, checker); err != nil {
			b.Fatalf("Register() = %v", err)
		}
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = registry.Checks(TargetReady)
	}
}

func BenchmarkGateCheck(b *testing.B) {
	gate, err := NewGate("ready_gate", Healthy("ready_gate"))
	if err != nil {
		b.Fatalf("NewGate() = %v", err)
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = gate.Check(context.Background())
	}
}

func BenchmarkReportReasons(b *testing.B) {
	report := Report{Target: TargetReady, Status: StatusUnhealthy, Checks: []Result{
		Healthy("storage"),
		Degraded("queue", ReasonOverloaded, "overloaded"),
		Unknown("cache", ReasonNotObserved, "unknown"),
		Unhealthy("database", ReasonFatal, "fatal"),
	}}
	b.ReportAllocs()
	for b.Loop() {
		_ = report.Reasons()
	}
}
