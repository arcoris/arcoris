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

package delay

import (
	"testing"
	"time"
)

var (
	benchmarkScheduleSink Schedule
	benchmarkSequenceSink Sequence
	benchmarkDurationSink time.Duration
	benchmarkBoolSink     bool
)

func BenchmarkImmediateConstruction(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Immediate()
	}
}

func BenchmarkFixedConstruction(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Fixed(time.Second)
	}
}

func BenchmarkDelaysConstructionSmall(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Delays(0, time.Millisecond, time.Second)
	}
}

func BenchmarkDelaysConstructionLarge(b *testing.B) {
	delays := benchmarkDelayList(256)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Delays(delays...)
	}
}

func BenchmarkLinearConstruction(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Linear(time.Millisecond, time.Millisecond)
	}
}

func BenchmarkExponentialConstruction(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Exponential(time.Millisecond, 2)
	}
}

func BenchmarkFibonacciConstruction(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Fibonacci(time.Millisecond)
	}
}

func BenchmarkCapConstruction(b *testing.B) {
	child := Fixed(time.Second)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Cap(child, time.Second)
	}
}

func BenchmarkLimitConstruction(b *testing.B) {
	child := Fixed(time.Second)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Limit(child, 3)
	}
}

func BenchmarkChainConstruction(b *testing.B) {
	first := Delays(0)
	second := Exponential(time.Millisecond, 2)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Chain(first, second)
	}
}

func BenchmarkChainConstructionLarge(b *testing.B) {
	schedules := benchmarkScheduleList(16)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = Chain(schedules...)
	}
}

func BenchmarkImmediateNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Immediate())
}

func BenchmarkFixedNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Fixed(time.Second))
}

func BenchmarkDelaysNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Delays(time.Second, 2*time.Second, 3*time.Second))
}

func BenchmarkLinearNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Linear(time.Millisecond, time.Millisecond))
}

func BenchmarkExponentialNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Exponential(time.Millisecond, 2))
}

func BenchmarkFibonacciNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Fibonacci(time.Millisecond))
}

func BenchmarkCapNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Cap(Fixed(time.Second), time.Second))
}

func BenchmarkLimitNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Limit(Fixed(time.Second), 3))
}

func BenchmarkChainNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Chain(Delays(0), Exponential(time.Millisecond, 2)))
}

func BenchmarkChainNewSequenceLarge(b *testing.B) {
	benchmarkNewSequence(b, Chain(benchmarkScheduleList(16)...))
}

func BenchmarkImmediateNext(b *testing.B) {
	benchmarkNext(b, Immediate().NewSequence())
}

func BenchmarkFixedNext(b *testing.B) {
	benchmarkNext(b, Fixed(time.Second).NewSequence())
}

func BenchmarkDelaysNext(b *testing.B) {
	delays := benchmarkDelayList(256)
	seq := &sequenceScheduleSequence{delays: delays}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, ok := seq.Next()
		if !ok {
			seq.next = 0
			d, ok = seq.Next()
		}

		benchmarkDurationSink = d
		benchmarkBoolSink = ok
	}
}

func BenchmarkLinearNext(b *testing.B) {
	benchmarkNext(b, Linear(time.Millisecond, time.Millisecond).NewSequence())
}

func BenchmarkExponentialNext(b *testing.B) {
	benchmarkNext(b, Exponential(time.Millisecond, 2).NewSequence())
}

func BenchmarkFibonacciNext(b *testing.B) {
	benchmarkNext(b, Fibonacci(time.Millisecond).NewSequence())
}

func BenchmarkCapNext(b *testing.B) {
	benchmarkNext(b, Cap(Fixed(2*time.Second), time.Second).NewSequence())
}

func BenchmarkLimitNext(b *testing.B) {
	child := Fixed(time.Second).NewSequence()
	seq := &limitSequence{child: child, remaining: b.N + 1}

	benchmarkNext(b, seq)
}

func BenchmarkChainNextFinitePrefix(b *testing.B) {
	first := &sequenceScheduleSequence{delays: benchmarkDelayList(256)}
	second := Fixed(time.Second).NewSequence()
	seq := &pairChainSequence{
		first:  first,
		second: second,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if first.next >= len(first.delays) {
			first.next = 0
			seq.next = 0
		}

		d, ok := seq.Next()
		benchmarkDurationSink = d
		benchmarkBoolSink = ok
	}
}

func BenchmarkChainNextInfiniteTail(b *testing.B) {
	benchmarkNext(b, Chain(Delays(0), Fixed(time.Second)).NewSequence())
}

func BenchmarkNestedCompositionNext(b *testing.B) {
	schedule := Limit(
		Cap(
			Chain(
				Delays(0),
				Exponential(time.Millisecond, 2),
			),
			time.Second,
		),
		b.N+1,
	)

	benchmarkNext(b, schedule.NewSequence())
}

func benchmarkNewSequence(b *testing.B, schedule Schedule) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkSequenceSink = schedule.NewSequence()
	}
}

func benchmarkNext(b *testing.B, sequence Sequence) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, ok := sequence.Next()
		benchmarkDurationSink = d
		benchmarkBoolSink = ok
	}
}

func benchmarkDelayList(n int) []time.Duration {
	delays := make([]time.Duration, n)
	for i := range delays {
		delays[i] = time.Duration(i) * time.Millisecond
	}

	return delays
}

func benchmarkScheduleList(n int) []Schedule {
	schedules := make([]Schedule, n)
	for i := range schedules {
		schedules[i] = Delays(time.Duration(i) * time.Millisecond)
	}

	return schedules
}
