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

package retry

import (
	"errors"
	"testing"
)

func TestRetryErrorKindError(t *testing.T) {
	tests := []struct {
		name string
		kind retryErrorKind
		want string
	}{
		{
			name: "exhausted",
			kind: retryErrorKindExhausted,
			want: errExhaustedMessage,
		},
		{
			name: "interrupted",
			kind: retryErrorKindInterrupted,
			want: errInterruptedMessage,
		},
		{
			name: "unknown",
			kind: retryErrorKind(255),
			want: "retry: unknown error",
		},
		{
			name: "zero",
			kind: 0,
			want: "retry: unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.Error(); got != tt.want {
				t.Fatalf("retryErrorKind.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRetryErrorKindIs(t *testing.T) {
	if !errors.Is(retryErrorKindExhausted, ErrExhausted) {
		t.Fatalf("retryErrorKindExhausted does not match ErrExhausted")
	}
	if errors.Is(retryErrorKindExhausted, ErrInterrupted) {
		t.Fatalf("retryErrorKindExhausted matched ErrInterrupted")
	}

	if !errors.Is(retryErrorKindInterrupted, ErrInterrupted) {
		t.Fatalf("retryErrorKindInterrupted does not match ErrInterrupted")
	}
	if errors.Is(retryErrorKindInterrupted, ErrExhausted) {
		t.Fatalf("retryErrorKindInterrupted matched ErrExhausted")
	}

	if errors.Is(retryErrorKind(255), ErrExhausted) {
		t.Fatalf("unknown retryErrorKind matched ErrExhausted")
	}
	if errors.Is(retryErrorKind(255), ErrInterrupted) {
		t.Fatalf("unknown retryErrorKind matched ErrInterrupted")
	}
}

func TestRetryErrorMessage(t *testing.T) {
	errBoom := errors.New("boom")

	if got := retryErrorMessage(ErrExhausted, nil); got != errExhaustedMessage {
		t.Fatalf("retryErrorMessage without cause = %q, want %q", got, errExhaustedMessage)
	}

	want := errExhaustedMessage + ": " + errBoom.Error()
	if got := retryErrorMessage(ErrExhausted, errBoom); got != want {
		t.Fatalf("retryErrorMessage with cause = %q, want %q", got, want)
	}
}
