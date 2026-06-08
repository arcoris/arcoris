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

package deadline

import (
	"context"
	"testing"
	"time"
)

func BenchmarkCanStartNoDeadline(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	now := testNow()
	for b.Loop() {
		benchmarkDecision = CanStart(ctx, now, time.Second)
	}
}

func BenchmarkCanStartAllowed(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkDecision = CanStart(ctx, now, time.Millisecond)
	}
}

func BenchmarkCanStartInsufficient(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
	defer cancel()

	for b.Loop() {
		benchmarkDecision = CanStart(ctx, now, time.Second)
	}
}

func BenchmarkCanStartExpired(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(-time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkDecision = CanStart(ctx, now, 0)
	}
}

func BenchmarkCanStartContextDone(b *testing.B) {
	b.ReportAllocs()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	now := testNow()

	for b.Loop() {
		benchmarkDecision = CanStart(ctx, now, 0)
	}
}

func BenchmarkDecisionIsValidAllowed(b *testing.B) {
	b.ReportAllocs()

	decision := Decision{
		Allowed:   true,
		Remaining: time.Second,
		Reason:    ReasonAllowed,
	}
	for b.Loop() {
		benchmarkDurationOK = decision.IsValid()
	}
}

func BenchmarkDecisionIsValidDenied(b *testing.B) {
	b.ReportAllocs()

	decision := Decision{
		Remaining: time.Millisecond,
		Reason:    ReasonInsufficientBudget,
	}
	for b.Loop() {
		benchmarkDurationOK = decision.IsValid()
	}
}

func BenchmarkReasonString(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		benchmarkReasonString = ReasonInsufficientBudget.String()
	}
}
