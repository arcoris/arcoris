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
)

func testNow() time.Time {
	return time.Now().UTC().Add(24 * time.Hour)
}

func contextWithDeadline(t *testing.T, dl time.Time) context.Context {
	t.Helper()

	ctx, cancel := context.WithDeadline(context.Background(), dl)
	t.Cleanup(cancel)
	return ctx
}

func assertContextDeadline(t *testing.T, ctx context.Context, want time.Time) {
	t.Helper()

	got, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("context has no deadline, want %v", want)
	}
	if !got.Equal(want) {
		t.Fatalf("deadline = %v, want %v", got, want)
	}
}

func requirePanic(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected panic %q", want)
		}

		got, ok := recovered.(string)
		if !ok {
			t.Fatalf("panic type = %T, want string %q", recovered, want)
		}
		if got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	fn()
}
