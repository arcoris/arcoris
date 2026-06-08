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

package retry

import (
	"context"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
)

func BenchmarkDoSuccessDefault(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return nil }

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op)
	}
}

func BenchmarkDoNonRetryableDefault(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return benchmarkErrBoom }

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op)
	}
}

func BenchmarkDoRetryImmediateSuccessSecondAttempt(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		calls := 0
		benchmarkErrorSink = Do(
			context.Background(),
			func(context.Context) error {
				calls++
				if calls == 1 {
					return benchmarkErrBoom
				}
				return nil
			},
			WithClassifier(RetryAll()),
			WithMaxAttempts(2),
			WithDelaySchedule(delay.Immediate()),
		)
		benchmarkIntSink = calls
	}
}

func BenchmarkDoRetryImmediateExhaustedMaxAttempts(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return benchmarkErrBoom }
	opts := []Option{
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	}

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, opts...)
	}
}

func BenchmarkDoValueSuccess(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) (int, error) { return 42, nil }

	for b.Loop() {
		benchmarkIntSink, benchmarkErrorSink = DoValue(context.Background(), op)
	}
}

func BenchmarkDoValueRetrySuccess(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		calls := 0
		benchmarkIntSink, benchmarkErrorSink = DoValue(
			context.Background(),
			func(context.Context) (int, error) {
				calls++
				if calls == 1 {
					return 0, benchmarkErrBoom
				}
				return 42, nil
			},
			WithClassifier(RetryAll()),
			WithMaxAttempts(2),
			WithDelaySchedule(delay.Immediate()),
		)
	}
}

func BenchmarkDoObservedSuccess(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return nil }

	for b.Loop() {
		benchmarkOutcomeSink, benchmarkErrorSink = DoObserved(context.Background(), op)
	}
}

func BenchmarkDoObservedNonRetryable(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return benchmarkErrBoom }

	for b.Loop() {
		benchmarkOutcomeSink, benchmarkErrorSink = DoObserved(context.Background(), op)
	}
}

func BenchmarkDoValueObservedSuccess(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) (int, error) { return 42, nil }

	for b.Loop() {
		benchmarkIntSink, benchmarkOutcomeSink, benchmarkErrorSink = DoValueObserved(context.Background(), op)
	}
}

func BenchmarkDoContextAlreadyCanceled(b *testing.B) {
	b.ReportAllocs()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	op := func(context.Context) error { return nil }

	for b.Loop() {
		benchmarkErrorSink = Do(ctx, op)
	}
}

func BenchmarkDoDelaySequenceExhausted(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return benchmarkErrBoom }
	opts := []Option{
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithDelaySchedule(delay.Limit(delay.Immediate(), 0)),
	}

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, opts...)
	}
}

func BenchmarkDoDeadlineExhausted(b *testing.B) {
	b.ReportAllocs()

	now := benchmarkNow()
	fake := clock.NewFakeClock(now)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
	defer cancel()
	op := func(context.Context) error { return benchmarkErrBoom }
	opts := []Option{
		WithClock(fake),
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithDelaySchedule(delay.Fixed(time.Second)),
	}

	for b.Loop() {
		benchmarkErrorSink = Do(ctx, op, opts...)
	}
}

func BenchmarkDoMaxElapsedExhausted(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return benchmarkErrBoom }
	opts := []Option{
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithDelaySchedule(delay.Fixed(time.Second)),
		WithMaxElapsed(time.Nanosecond),
	}

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, opts...)
	}
}
