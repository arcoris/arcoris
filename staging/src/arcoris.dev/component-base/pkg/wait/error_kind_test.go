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

package wait

import (
	"errors"
	"testing"
)

// TestWaitErrorSentinelHierarchy verifies the wait-owned error classification
// hierarchy exposed through errors.Is.
func TestWaitErrorSentinelHierarchy(t *testing.T) {
	t.Parallel()

	if !errors.Is(ErrInterrupted, ErrInterrupted) {
		t.Fatal("errors.Is(ErrInterrupted, ErrInterrupted) = false, want true")
	}
	if errors.Is(ErrInterrupted, ErrTimeout) {
		t.Fatal("errors.Is(ErrInterrupted, ErrTimeout) = true, want false")
	}
	if !errors.Is(ErrTimeout, ErrTimeout) {
		t.Fatal("errors.Is(ErrTimeout, ErrTimeout) = false, want true")
	}
	if !errors.Is(ErrTimeout, ErrInterrupted) {
		t.Fatal("errors.Is(ErrTimeout, ErrInterrupted) = false, want true")
	}
}

// TestWaitErrorKindMessages verifies stable diagnostic strings for wait error
// sentinels.
func TestWaitErrorKindMessages(t *testing.T) {
	t.Parallel()

	mustHaveMessage(t, ErrInterrupted, "wait: interrupted")
	mustHaveMessage(t, ErrTimeout, "wait: timeout")
}

// TestUnknownWaitErrorKindMessage verifies that an invalid private sentinel kind
// still returns a non-empty diagnostic string.
func TestUnknownWaitErrorKindMessage(t *testing.T) {
	t.Parallel()

	kind := waitErrorKind(0)
	mustHaveMessage(t, kind, "wait: unknown error")
}

// TestWaitErrorMessage verifies shared wait error message formatting.
func TestWaitErrorMessage(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")

	if got, want := waitErrorMessage(ErrInterrupted, nil), "wait: interrupted"; got != want {
		t.Fatalf("waitErrorMessage(ErrInterrupted, nil) = %q, want %q", got, want)
	}
	if got, want := waitErrorMessage(ErrTimeout, cause), "wait: timeout: cause"; got != want {
		t.Fatalf("waitErrorMessage(ErrTimeout, cause) = %q, want %q", got, want)
	}
}
