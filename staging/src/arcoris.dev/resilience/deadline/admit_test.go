/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package deadline

import (
	"context"
	"testing"
	"time"

	"arcoris.dev/admission"
)

func TestTryAdmit(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name     string
		request  Request
		want     admission.Decision
		admitted bool
		denied   bool
		metadata Decision
	}{
		{
			name: "no deadline",
			request: Request{
				Context: context.Background(),
				Now:     now,
				Min:     time.Second,
			},
			want:     admission.Admit(admission.ReasonAdmitted),
			admitted: true,
			metadata: Decision{
				Allowed: true,
				Reason:  ReasonNoDeadline,
			},
		},
		{
			name: "enough budget",
			request: Request{
				Context: contextWithDeadline(t, now.Add(10*time.Second)),
				Now:     now,
				Min:     time.Second,
			},
			want:     admission.Admit(admission.ReasonAdmitted),
			admitted: true,
			metadata: Decision{
				Allowed:   true,
				Remaining: 10 * time.Second,
				Reason:    ReasonAllowed,
			},
		},
		{
			name: "expired deadline",
			request: Request{
				Context: contextWithDeadline(t, now),
				Now:     now,
			},
			want:   admission.Deny(admission.ReasonDeadlineExceeded),
			denied: true,
			metadata: Decision{
				Reason: ReasonExpired,
			},
		},
		{
			name: "insufficient budget",
			request: Request{
				Context: contextWithDeadline(t, now.Add(time.Second)),
				Now:     now,
				Min:     2 * time.Second,
			},
			want:   admission.Deny(admission.ReasonDeadlineExceeded),
			denied: true,
			metadata: Decision{
				Remaining: time.Second,
				Reason:    ReasonInsufficientBudget,
			},
		},
		{
			name: "context done",
			request: Request{
				Context: canceled,
				Now:     now,
			},
			want:   admission.Deny(admission.ReasonCanceled),
			denied: true,
			metadata: Decision{
				Reason: ReasonContextDone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := TryAdmit(test.request)
			if !result.IsValid() {
				t.Fatalf("TryAdmit result is invalid: %+v", result.Decision())
			}
			if got := result.Decision(); got != test.want {
				t.Fatalf("decision = %+v, want %+v", got, test.want)
			}
			if got := result.IsAdmitted(); got != test.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, test.admitted)
			}
			if got := result.IsDenied(); got != test.denied {
				t.Fatalf("IsDenied = %v, want %v", got, test.denied)
			}
			if result.HasSideEffect() {
				t.Fatal("TryAdmit result has side effect, want none")
			}
			if result.HasGrant() {
				t.Fatal("TryAdmit result has grant, want none")
			}
			if _, ok := result.Grant(); ok {
				t.Fatal("Grant() ok=true, want false")
			}
			if !result.HasMetadata() {
				t.Fatal("TryAdmit result has no metadata")
			}
			if metadata, ok := result.Metadata(); !ok || metadata != test.metadata {
				t.Fatalf("metadata = (%+v, %t), want (%+v, true)", metadata, ok, test.metadata)
			}
		})
	}
}

func TestTryAdmitPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_ = TryAdmit(Request{})
	})
	requirePanic(t, panicNegativeDuration("min"), func() {
		_ = TryAdmit(Request{
			Context: context.Background(),
			Min:     -time.Nanosecond,
		})
	})
	requirePanic(t, panicNilContext, func() {
		_ = TryAdmit(Request{
			Min: -time.Nanosecond,
		})
	})
}
