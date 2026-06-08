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

func BenchmarkClampNoDeadline(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	now := testNow()
	for b.Loop() {
		benchmarkDuration, benchmarkDurationOK = Clamp(ctx, now, time.Second)
	}
}

func BenchmarkClampBoundedBelowRemaining(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkDuration, benchmarkDurationOK = Clamp(ctx, now, time.Millisecond)
	}
}

func BenchmarkClampBoundedAboveRemaining(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
	defer cancel()

	for b.Loop() {
		benchmarkDuration, benchmarkDurationOK = Clamp(ctx, now, time.Second)
	}
}

func BenchmarkClampExpired(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(-time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkDuration, benchmarkDurationOK = Clamp(ctx, now, time.Second)
	}
}

func BenchmarkReserveNoDeadline(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	now := testNow()
	for b.Loop() {
		benchmarkDuration, _, benchmarkDurationOK = Reserve(ctx, now, time.Second)
	}
}

func BenchmarkReserveBoundedAllowed(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkDuration, _, benchmarkDurationOK = Reserve(ctx, now, time.Millisecond)
	}
}

func BenchmarkReserveBoundedDenied(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
	defer cancel()

	for b.Loop() {
		benchmarkDuration, _, benchmarkDurationOK = Reserve(ctx, now, time.Second)
	}
}

func BenchmarkReserveContextDone(b *testing.B) {
	b.ReportAllocs()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	now := testNow()

	for b.Loop() {
		benchmarkDuration, _, benchmarkDurationOK = Reserve(ctx, now, time.Second)
	}
}

func BenchmarkReserveBudgetNoDeadline(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	now := testNow()
	for b.Loop() {
		benchmarkReserveResult = ReserveBudget(ctx, now, time.Second)
	}
}

func BenchmarkReserveBudgetBoundedAllowed(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
	defer cancel()

	for b.Loop() {
		benchmarkReserveResult = ReserveBudget(ctx, now, time.Millisecond)
	}
}

func BenchmarkReserveBudgetBoundedDenied(b *testing.B) {
	b.ReportAllocs()

	now := testNow()
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
	defer cancel()

	for b.Loop() {
		benchmarkReserveResult = ReserveBudget(ctx, now, time.Second)
	}
}
