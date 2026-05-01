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

package lifecycle

import (
	"errors"
	"testing"
	"time"
)

var testTime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

type testClock struct {
	now time.Time
}

func (c testClock) Now() time.Time {
	return c.now
}

func mustMatch(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false, want true", err, target)
	}
}

func mustNotMatch(t *testing.T, err error, target error) {
	t.Helper()

	if errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = true, want false", err, target)
	}
}

func assertTransitionEqual(t *testing.T, got, want Transition) {
	t.Helper()

	if got.From != want.From {
		t.Fatalf("transition.From = %s, want %s", got.From, want.From)
	}
	if got.To != want.To {
		t.Fatalf("transition.To = %s, want %s", got.To, want.To)
	}
	if got.Event != want.Event {
		t.Fatalf("transition.Event = %s, want %s", got.Event, want.Event)
	}
	if got.Revision != want.Revision {
		t.Fatalf("transition.Revision = %d, want %d", got.Revision, want.Revision)
	}
	if got.At.IsZero() != want.At.IsZero() {
		t.Fatalf("transition.At zero = %v, want %v", got.At.IsZero(), want.At.IsZero())
	}
	if !got.At.IsZero() && !got.At.Equal(want.At) {
		t.Fatalf("transition.At = %v, want %v", got.At, want.At)
	}
	if want.Cause == nil {
		if got.Cause != nil {
			t.Fatalf("transition.Cause = %v, want nil", got.Cause)
		}
		return
	}
	if got.Cause == nil {
		t.Fatal("transition.Cause = nil, want non-nil")
	}
	if !errors.Is(got.Cause, want.Cause) {
		t.Fatalf("errors.Is(transition.Cause, %v) = false, want true", want.Cause)
	}
}

func assertSnapshotEqual(t *testing.T, got, want Snapshot) {
	t.Helper()

	if got.State != want.State {
		t.Fatalf("snapshot.State = %s, want %s", got.State, want.State)
	}
	if got.Revision != want.Revision {
		t.Fatalf("snapshot.Revision = %d, want %d", got.Revision, want.Revision)
	}
	assertTransitionEqual(t, got.LastTransition, want.LastTransition)
	if want.FailureCause == nil {
		if got.FailureCause != nil {
			t.Fatalf("snapshot.FailureCause = %v, want nil", got.FailureCause)
		}
		return
	}
	if got.FailureCause == nil {
		t.Fatal("snapshot.FailureCause = nil, want non-nil")
	}
	if !errors.Is(got.FailureCause, want.FailureCause) {
		t.Fatalf("errors.Is(snapshot.FailureCause, %v) = false, want true", want.FailureCause)
	}
}

func mustReceiveSnapshot(t *testing.T, ch <-chan Snapshot) Snapshot {
	t.Helper()

	select {
	case snapshot := <-ch:
		return snapshot
	case <-time.After(time.Second):
		t.Fatal("snapshot was not received before safety timeout")
		return Snapshot{}
	}
}

func mustReceiveError(t *testing.T, ch <-chan error) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
		t.Fatal("error was not received before safety timeout")
		return nil
	}
}

func mustSignalClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("signal was not closed before safety timeout")
	}
}

func mustNotSignalClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal("signal is closed, want open")
	default:
	}
}
