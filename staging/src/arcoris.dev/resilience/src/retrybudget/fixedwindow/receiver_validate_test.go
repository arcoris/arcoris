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
	"errors"
	"testing"
)

func TestLimiterNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var limiter *Limiter
	tests := []struct {
		name string
		call func()
	}{
		{name: "RecordOriginal", call: func() { limiter.RecordOriginal() }},
		{name: "TryAdmitRetry", call: func() { _ = limiter.TryAdmitRetry() }},
		{name: "Snapshot", call: func() { _ = limiter.Snapshot() }},
		{name: "Revision", call: func() { _ = limiter.Revision() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requirePanicError(t, ErrNilLimiter, tt.call)
		})
	}
}

func TestLimiterZeroValueReceiverPanics(t *testing.T) {
	t.Parallel()

	var limiter Limiter
	tests := []struct {
		name string
		call func()
	}{
		{name: "RecordOriginal", call: func() { limiter.RecordOriginal() }},
		{name: "TryAdmitRetry", call: func() { _ = limiter.TryAdmitRetry() }},
		{name: "Snapshot", call: func() { _ = limiter.Snapshot() }},
		{name: "Revision", call: func() { _ = limiter.Revision() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requirePanicError(t, ErrUninitializedLimiter, tt.call)
		})
	}
}

func TestLimiterReadyReceiverDoesNotPanic(t *testing.T) {
	t.Parallel()

	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(2))
	limiter.RecordOriginal()
	_ = limiter.Snapshot()
	_ = limiter.Revision()
	_ = limiter.TryAdmitRetry()
}

func requirePanicError(t *testing.T, want error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %T(%v), want error %v", recovered, recovered, want)
		}
		if !errors.Is(err, want) {
			t.Fatalf("panic = %v, want %v", err, want)
		}
	}()

	fn()
}
