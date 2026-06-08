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

import "testing"

func TestProgrammingErrorSentinels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{name: "nil context", err: ErrNilContext, msg: "retry: nil context"},
		{name: "nil operation", err: ErrNilOperation, msg: "retry: nil operation"},
		{name: "nil value operation", err: ErrNilValueOperation, msg: "retry: nil value operation"},
		{name: "nil clock", err: ErrNilClock, msg: "retry: nil clock"},
		{name: "nil delay schedule", err: ErrNilDelaySchedule, msg: "retry: nil delay schedule"},
		{name: "nil delay sequence", err: ErrNilDelaySequence, msg: "retry: delay schedule returned nil Sequence"},
		{name: "negative delay", err: ErrNegativeDelay, msg: "retry: delay sequence returned negative delay"},
		{name: "nil classifier", err: ErrNilClassifier, msg: "retry: nil classifier"},
		{name: "nil classifier func", err: ErrNilClassifierFunc, msg: "retry: nil classifier function"},
		{name: "zero max attempts", err: ErrZeroMaxAttempts, msg: "retry: zero max attempts"},
		{name: "negative max elapsed", err: ErrNegativeMaxElapsed, msg: "retry: negative max elapsed"},
		{name: "nil observer", err: ErrNilObserver, msg: "retry: nil observer"},
		{name: "nil observer func", err: ErrNilObserverFunc, msg: "retry: nil observer function"},
		{name: "nil option", err: ErrNilOption, msg: "retry: nil option"},
		{name: "invalid exhausted outcome", err: ErrInvalidExhaustedOutcome, msg: "retry: invalid exhausted outcome"},
		{name: "non exhausted outcome reason", err: ErrNonExhaustedOutcomeReason, msg: "retry: non-exhausted outcome reason"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err == nil {
				t.Fatal("sentinel error is nil")
			}
			if tt.err.Error() != tt.msg {
				t.Fatalf("Error() = %q, want %q", tt.err.Error(), tt.msg)
			}
		})
	}
}
