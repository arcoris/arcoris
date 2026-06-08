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

	"arcoris.dev/chrono/clock"
)

func BenchmarkOutcomeIsValidSuccess(b *testing.B) {
	b.ReportAllocs()

	outcome := retryTestSuccessOutcome(1)

	for b.Loop() {
		benchmarkBoolSink = outcome.IsValid()
	}
}

func BenchmarkOutcomeIsValidExhausted(b *testing.B) {
	b.ReportAllocs()

	outcome := retryTestFailureOutcome(1, StopReasonMaxAttempts, benchmarkErrBoom)

	for b.Loop() {
		benchmarkBoolSink = outcome.IsValid()
	}
}

func BenchmarkEventIsValidAttemptStart(b *testing.B) {
	b.ReportAllocs()

	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}

	for b.Loop() {
		benchmarkBoolSink = event.IsValid()
	}
}

func BenchmarkEventIsValidRetryStop(b *testing.B) {
	b.ReportAllocs()

	event := Event{
		Kind:    EventRetryStop,
		Attempt: retryTestAttempt(1),
		Err:     benchmarkErrBoom,
		Outcome: retryTestFailureOutcome(
			1,
			StopReasonMaxAttempts,
			benchmarkErrBoom,
		),
	}

	for b.Loop() {
		benchmarkBoolSink = event.IsValid()
	}
}

func BenchmarkStopReasonString(b *testing.B) {
	b.ReportAllocs()

	reason := StopReasonMaxAttempts

	for b.Loop() {
		benchmarkReasonSink = reason
		benchmarkStringSink = reason.String()
	}
}

func BenchmarkExhaustedError(b *testing.B) {
	b.ReportAllocs()

	outcome := retryTestFailureOutcome(1, StopReasonMaxAttempts, benchmarkErrBoom)

	for b.Loop() {
		benchmarkErrorSink = NewExhaustedError(outcome)
	}
}

func BenchmarkInterruptedError(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		benchmarkErrorSink = NewInterruptedError(context.Canceled)
	}
}

func BenchmarkWaitDelayZero(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		benchmarkErrorSink = waitDelay(context.Background(), clock.RealClock{}, 0)
	}
}
