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

package jitter

import (
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
)

var (
	benchmarkBoolSink     bool
	benchmarkDurationSink time.Duration
	benchmarkScheduleSink delay.Schedule
	benchmarkSequenceSink delay.Sequence
)

// Construction benchmarks include schedule construction and option validation.
func BenchmarkFullConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Full(delay.Fixed(time.Second), WithRandom(fixedRandom(0)))
	})
}

func BenchmarkEqualConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Equal(delay.Fixed(time.Second), WithRandom(fixedRandom(0)))
	})
}

func BenchmarkPositiveConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Positive(delay.Fixed(time.Second), 0.2, WithRandom(fixedRandom(0)))
	})
}

func BenchmarkProportionalConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Proportional(delay.Fixed(time.Second), 0.2, WithRandom(fixedRandom(0)))
	})
}

func BenchmarkUniformConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Uniform(0, time.Second, WithRandom(fixedRandom(0)))
	})
}

func BenchmarkDecorrelatedConstruction(b *testing.B) {
	benchmarkConstruction(b, func() delay.Schedule {
		return Decorrelated(time.Second, 10*time.Second, 2, WithRandom(fixedRandom(0)))
	})
}

func benchmarkConstruction(b *testing.B, fn func() delay.Schedule) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkScheduleSink = fn()
	}
}

// NewSequence benchmarks include child sequence creation and per-sequence
// random generator creation.
func BenchmarkFullNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Full(delay.Fixed(time.Second), WithSeed(1)))
}

func BenchmarkEqualNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Equal(delay.Fixed(time.Second), WithSeed(1)))
}

func BenchmarkPositiveNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Positive(delay.Fixed(time.Second), 0.2, WithSeed(1)))
}

func BenchmarkProportionalNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Proportional(delay.Fixed(time.Second), 0.2, WithSeed(1)))
}

func BenchmarkUniformNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Uniform(0, time.Second, WithSeed(1)))
}

func BenchmarkDecorrelatedNewSequence(b *testing.B) {
	benchmarkNewSequence(b, Decorrelated(time.Second, 10*time.Second, 2, WithSeed(1)))
}

func benchmarkNewSequence(b *testing.B, sched delay.Schedule) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkSequenceSink = sched.NewSequence()
	}
}

// Next benchmarks measure steady-state sequence iteration after setup.
func BenchmarkFullNext(b *testing.B) {
	benchmarkNext(b, Full(delay.Fixed(time.Second), WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkEqualNext(b *testing.B) {
	benchmarkNext(b, Equal(delay.Fixed(time.Second), WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkPositiveNext(b *testing.B) {
	benchmarkNext(b, Positive(delay.Fixed(time.Second), 0.2, WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkProportionalNext(b *testing.B) {
	benchmarkNext(b, Proportional(delay.Fixed(time.Second), 0.2, WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkUniformNext(b *testing.B) {
	benchmarkNext(b, Uniform(0, time.Second, WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkDecorrelatedNext(b *testing.B) {
	benchmarkNext(b, Decorrelated(time.Second, 10*time.Second, 2, WithRandom(fixedRandom(0))).NewSequence())
}

func BenchmarkJitterWrapperExhaustion(b *testing.B) {
	benchmarkNext(b, Full(delay.Delays(), WithRandom(fixedRandom(0))).NewSequence())
}

func benchmarkNext(b *testing.B, seq delay.Sequence) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, ok := seq.Next()
		benchmarkDurationSink = d
		benchmarkBoolSink = ok
	}
}

// Random duration benchmarks isolate inclusive integer-nanosecond mapping.
func BenchmarkRandomDurationInclusiveSmallBound(b *testing.B) {
	benchmarkRandomDurationInclusive(b, fixedRandom(5), 10*time.Nanosecond)
}

func BenchmarkRandomDurationInclusiveLargeBound(b *testing.B) {
	benchmarkRandomDurationInclusive(b, fixedRandom(int64(time.Hour)), 24*time.Hour)
}

func BenchmarkRandomDurationInclusiveMaxDuration(b *testing.B) {
	benchmarkRandomDurationInclusive(b, fixedRandom(int64(maxDuration)), maxDuration)
}

func BenchmarkRandomDurationInclusiveRejection(b *testing.B) {
	benchmarkRandomDurationInclusive(
		b,
		&sequenceRandom{values: []int64{int64(maxDuration), 5}},
		10*time.Nanosecond,
	)
}

func benchmarkRandomDurationInclusive(b *testing.B, random RandomGenerator, maxOffset time.Duration) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkDurationSink = randomDurationInclusive(random, maxOffset)
	}
}
