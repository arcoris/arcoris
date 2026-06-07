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

func BenchmarkResultIsValid(b *testing.B) {
	result := Degraded("queue", ReasonOverloaded, "overloaded")
	b.ReportAllocs()
	for b.Loop() {
		_ = result.IsValid()
	}
}

func BenchmarkNewCheckSet(b *testing.B) {
	checks := []Checker{
		MustCheck("storage", func(ctx context.Context) Result { return Healthy("storage") }),
		MustCheck("queue", func(ctx context.Context) Result { return Healthy("queue") }),
		MustCheck("database", func(ctx context.Context) Result { return Healthy("database") }),
	}

	b.ReportAllocs()
	for b.Loop() {
		set, err := NewCheckSet(TargetReady, checks...)
		if err != nil {
			b.Fatalf("NewCheckSet() = %v", err)
		}
		_ = set
	}
}

func BenchmarkCheckSetChecks(b *testing.B) {
	set := MustCheckSet(
		TargetReady,
		MustCheck("storage", func(ctx context.Context) Result { return Healthy("storage") }),
		MustCheck("queue", func(ctx context.Context) Result { return Healthy("queue") }),
		MustCheck("database", func(ctx context.Context) Result { return Healthy("database") }),
	)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Checks()
	}
}

func BenchmarkCheckSetRange(b *testing.B) {
	set := MustCheckSet(
		TargetReady,
		MustCheck("storage", func(ctx context.Context) Result { return Healthy("storage") }),
		MustCheck("queue", func(ctx context.Context) Result { return Healthy("queue") }),
		MustCheck("database", func(ctx context.Context) Result { return Healthy("database") }),
	)

	b.ReportAllocs()
	for b.Loop() {
		set.Range(func(Checker) bool {
			return true
		})
	}
}

func BenchmarkCheckSetHas(b *testing.B) {
	set := MustCheckSet(
		TargetReady,
		MustCheck("storage", func(ctx context.Context) Result { return Healthy("storage") }),
		MustCheck("queue", func(ctx context.Context) Result { return Healthy("queue") }),
		MustCheck("database", func(ctx context.Context) Result { return Healthy("database") }),
	)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Has("database")
	}
}

func BenchmarkReportAggregateStatus(b *testing.B) {
	report := Report{Target: TargetReady, Status: StatusUnhealthy, Checks: []Result{
		Healthy("storage"),
		Degraded("queue", ReasonOverloaded, "overloaded"),
		Unknown("cache", ReasonNotObserved, "unknown"),
		Unhealthy("database", ReasonFatal, "fatal"),
	}}
	b.ReportAllocs()
	for b.Loop() {
		_ = report.AggregateStatus()
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
