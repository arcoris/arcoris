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

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

var (
	benchmarkDecision retrybudget.Decision
	benchmarkSnapshot snapshot.Snapshot[retrybudget.Snapshot]
	benchmarkRevision snapshot.Revision
	benchmarkAllowed  uint64
)

func newBenchmarkLimiter(b *testing.B, opts ...Option) (*Limiter, *fakeClock) {
	b.Helper()
	clk := newFakeClock(fixedWindowTestNow)
	all := append([]Option{WithClock(clk)}, opts...)
	limiter, err := New(all...)
	if err != nil {
		b.Fatalf("New() error = %v", err)
	}
	return limiter, clk
}

func advanceBenchmarkWindow(clk *fakeClock) {
	clk.Add(time.Hour)
}
