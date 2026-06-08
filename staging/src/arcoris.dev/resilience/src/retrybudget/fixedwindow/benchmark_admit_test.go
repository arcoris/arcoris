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

package fixedwindow

import (
	"testing"
	"time"
)

func BenchmarkLimiterTryAdmitRetryAllowed(b *testing.B) {
	limiter, _ := newBenchmarkLimiter(b, WithRatio(RatioZero), WithMinRetries(^uint64(0)))
	b.ReportAllocs()
	for b.Loop() {
		benchmarkDecision = limiter.TryAdmitRetry()
	}
}

func BenchmarkLimiterTryAdmitRetryDenied(b *testing.B) {
	limiter, _ := newBenchmarkLimiter(b, WithRatio(RatioZero), WithMinRetries(0))
	b.ReportAllocs()
	for b.Loop() {
		benchmarkDecision = limiter.TryAdmitRetry()
	}
}

func BenchmarkLimiterTryAdmitRetryRotate(b *testing.B) {
	limiter, clk := newBenchmarkLimiter(b, WithWindow(time.Nanosecond), WithRatio(RatioZero), WithMinRetries(1))
	b.ReportAllocs()
	for b.Loop() {
		advanceBenchmarkWindow(clk)
		benchmarkDecision = limiter.TryAdmitRetry()
	}
}
