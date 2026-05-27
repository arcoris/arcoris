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

package wait

import (
	"context"
	"testing"
	"time"
)

func BenchmarkDelayZero(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Delay(ctx, 0)
	}
}

func BenchmarkTimerReset(b *testing.B) {
	timer := NewTimer(time.Hour)
	defer timer.StopAndDrain()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		timer.Reset(time.Hour)
	}
}

func BenchmarkUntilImmediateSuccess(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Until(ctx, time.Hour, Satisfied)
	}
}

func BenchmarkConditionAll(b *testing.B) {
	condition := All(Satisfied, Satisfied, Satisfied)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = condition(ctx)
	}
}

func BenchmarkConditionAny(b *testing.B) {
	condition := Any(Unsatisfied, Unsatisfied, Satisfied)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = condition(ctx)
	}
}
