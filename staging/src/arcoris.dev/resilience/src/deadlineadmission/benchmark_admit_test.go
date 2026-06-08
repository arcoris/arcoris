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

package deadlineadmission

import (
	"context"
	"testing"
	"time"

	"arcoris.dev/resilience/deadline"
)

func BenchmarkTryAdmitAllowed(b *testing.B) {
	b.ReportAllocs()

	req := Request{
		Context: context.Background(),
		Now:     testNow(),
		Min:     time.Second,
	}
	for b.Loop() {
		benchmarkResult = TryAdmit(req)
	}
}

func BenchmarkTryAdmitDeniedExpired(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(-time.Second))
	defer cancel()

	req := Request{
		Context: ctx,
		Now:     now,
	}
	for b.Loop() {
		benchmarkResult = TryAdmit(req)
	}
}

func BenchmarkTryAdmitDeniedContextDone(b *testing.B) {
	b.ReportAllocs()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := Request{
		Context: ctx,
		Now:     testNow(),
	}
	for b.Loop() {
		benchmarkResult = TryAdmit(req)
	}
}

func BenchmarkAdmissionResultAllowed(b *testing.B) {
	b.ReportAllocs()

	decision := deadline.Decision{
		Allowed:   true,
		Remaining: time.Second,
		Reason:    deadline.ReasonAllowed,
	}
	for b.Loop() {
		benchmarkResult = AdmissionResult(decision)
	}
}

func BenchmarkAdmissionResultDenied(b *testing.B) {
	b.ReportAllocs()

	decision := deadline.Decision{
		Remaining: time.Millisecond,
		Reason:    deadline.ReasonInsufficientBudget,
	}
	for b.Loop() {
		benchmarkResult = AdmissionResult(decision)
	}
}
