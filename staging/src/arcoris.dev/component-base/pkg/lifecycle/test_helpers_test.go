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
