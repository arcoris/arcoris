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

package clock

import (
	"fmt"
	"testing"
	"time"
)

var (
	benchmarkTimeSink    time.Time
	benchmarkDuration    time.Duration
	benchmarkChannelSink <-chan time.Time
	benchmarkTimerSink   Timer
	benchmarkTickerSink  Ticker
	benchmarkPendingSink Pending
)

func BenchmarkFakeClockNow(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkTimeSink = clk.Now()
	}
}

func BenchmarkFakeClockSince(b *testing.B) {
	start := fakeClockTestTime()
	clk := NewFakeClock(start.Add(time.Hour))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkDuration = clk.Since(start)
	}
}

func BenchmarkFakeClockUntil(b *testing.B) {
	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	deadline := start.Add(time.Hour)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkDuration = clk.Until(deadline)
	}
}

func BenchmarkFakeClockAfter(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := clk.After(0)
		benchmarkTimeSink = <-ch
		benchmarkChannelSink = ch
	}
}

func BenchmarkFakeClockNewTimer(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timer := clk.NewTimer(0)
		benchmarkTimeSink = <-timer.C()
		benchmarkTimerSink = timer
	}
}

func BenchmarkFakeClockNewTicker(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ticker := clk.NewTicker(time.Hour)
		ticker.Stop()
		benchmarkTickerSink = ticker
	}
}

func BenchmarkFakeClockStepNoPending(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		clk.Step(time.Nanosecond)
	}
}

func BenchmarkFakeClockWaiterRegistrationAndStep(b *testing.B) {
	for _, n := range []int{1, 8, 64, 256} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			benchmarkFakeClockWaiterRegistrationAndStep(b, n)
		})
	}
}

// benchmarkFakeClockWaiterRegistrationAndStep includes registration in each measured
// iteration because one-shot waiters are consumed by the delivery being
// benchmarked.
func benchmarkFakeClockWaiterRegistrationAndStep(b *testing.B, n int) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		clk := NewFakeClock(fakeClockTestTime())
		waiters := make([]<-chan time.Time, n)

		for j := range waiters {
			waiters[j] = clk.After(time.Second)
		}

		clk.Step(time.Second)
		benchmarkDrain(waiters)
	}
}

func BenchmarkFakeClockTimerRegistrationAndStep(b *testing.B) {
	for _, n := range []int{1, 8, 64, 256} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			benchmarkFakeClockTimerRegistrationAndStep(b, n)
		})
	}
}

// benchmarkFakeClockTimerRegistrationAndStep includes registration in each measured iteration
// because one-shot timers leave the active registry after firing.
func benchmarkFakeClockTimerRegistrationAndStep(b *testing.B, n int) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		clk := NewFakeClock(fakeClockTestTime())
		timers := make([]Timer, n)

		for j := range timers {
			timers[j] = clk.NewTimer(time.Second)
		}

		clk.Step(time.Second)

		for _, timer := range timers {
			benchmarkTimeSink = <-timer.C()
		}
	}
}

func BenchmarkFakeClockStepTickers(b *testing.B) {
	for _, n := range []int{1, 8, 64, 256} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			benchmarkFakeClockStepTickers(b, n)
		})
	}
}

// benchmarkFakeClockStepTickers measures Step delivery for already registered
// active tickers. Registration is outside the measured loop because tickers stay
// active after delivery.
func benchmarkFakeClockStepTickers(b *testing.B, n int) {
	clk := NewFakeClock(fakeClockTestTime())
	tickers := make([]Ticker, n)
	channels := make([]<-chan time.Time, n)

	for i := range tickers {
		tickers[i] = clk.NewTicker(time.Second)
		channels[i] = tickers[i].C()
	}
	defer func() {
		for _, ticker := range tickers {
			ticker.Stop()
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		clk.Step(time.Second)
		benchmarkDrain(channels)
	}
}

func BenchmarkFakeClockMixedRegistrationAndStep(b *testing.B) {
	for _, n := range []int{1, 8, 64, 256} {
		b.Run(fmt.Sprintf("N%d", n), func(b *testing.B) {
			benchmarkFakeClockMixedRegistrationAndStep(b, n)
		})
	}
}

// benchmarkFakeClockMixedRegistrationAndStep measures the complete one-shot
// registration and delivery cycle together with active ticker delivery.
func benchmarkFakeClockMixedRegistrationAndStep(b *testing.B, n int) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		clk := NewFakeClock(fakeClockTestTime())
		waiters := make([]<-chan time.Time, n)
		timers := make([]Timer, n)
		tickers := make([]Ticker, n)

		for j := 0; j < n; j++ {
			waiters[j] = clk.After(time.Second)
			timers[j] = clk.NewTimer(time.Second)
			tickers[j] = clk.NewTicker(time.Second)
		}

		clk.Step(time.Second)

		benchmarkDrain(waiters)
		for _, timer := range timers {
			benchmarkTimeSink = <-timer.C()
		}
		for _, ticker := range tickers {
			benchmarkTimeSink = <-ticker.C()
			ticker.Stop()
		}
	}
}

func BenchmarkFakeClockPending(b *testing.B) {
	clk := NewFakeClock(fakeClockTestTime())
	_ = clk.After(time.Hour)
	timer := clk.NewTimer(time.Hour)
	ticker := clk.NewTicker(time.Hour)
	defer timer.Stop()
	defer ticker.Stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkPendingSink = clk.Pending()
	}
}

func benchmarkDrain(channels []<-chan time.Time) {
	for _, ch := range channels {
		benchmarkTimeSink = <-ch
	}
}
