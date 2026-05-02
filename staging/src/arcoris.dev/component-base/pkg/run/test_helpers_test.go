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

package run

import (
	"fmt"
	"testing"
	"time"
)

const testTimeout = time.Second

func mustPanicWith(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected panic %q", want)
		}

		got := fmt.Sprint(recovered)
		if got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	fn()
}

func mustClose(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for channel close")
	}
}

func mustNotClose(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal("channel closed unexpectedly")
	case <-time.After(20 * time.Millisecond):
	}
}
