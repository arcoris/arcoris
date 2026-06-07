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

import "testing"

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
