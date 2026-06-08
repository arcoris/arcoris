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
)

func BenchmarkDoWithClassifier(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return nil }
	classifier := ClassifierFunc(func(error) bool { return true })

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, WithClassifier(classifier))
	}
}

func BenchmarkDoWithObserver(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return nil }
	observer := ObserverFunc(func(context.Context, Event) {})

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, WithObserver(observer))
	}
}

func BenchmarkDoWithManyObservers(b *testing.B) {
	b.ReportAllocs()

	op := func(context.Context) error { return nil }
	observer := ObserverFunc(func(context.Context, Event) {})
	opts := []Option{
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
	}

	for b.Loop() {
		benchmarkErrorSink = Do(context.Background(), op, opts...)
	}
}

func BenchmarkEmitRetryEventNoObservers(b *testing.B) {
	b.ReportAllocs()

	cfg := configOf()
	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}

	for b.Loop() {
		emitRetryEvent(context.Background(), cfg, event)
	}
}

func BenchmarkEmitRetryEventOneObserver(b *testing.B) {
	b.ReportAllocs()

	cfg := configOf(WithObserverFunc(func(context.Context, Event) {}))
	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}

	for b.Loop() {
		emitRetryEvent(context.Background(), cfg, event)
	}
}

func BenchmarkEmitRetryEventManyObservers(b *testing.B) {
	b.ReportAllocs()

	observer := ObserverFunc(func(context.Context, Event) {})
	cfg := configOf(
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
		WithObserver(observer),
	)
	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}

	for b.Loop() {
		emitRetryEvent(context.Background(), cfg, event)
	}
}
