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
	"context"
	"errors"
	"testing"
)

func TestWaitReturnsContextCause(t *testing.T) {
	t.Parallel()

	want := errors.New("stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if err := Wait(ctx); !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want %v", err, want)
	}
}

func TestWaitFallsBackToContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait error = %v, want context.Canceled", err)
	}
}

func TestWaitRejectsNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilWaitContext, func() {
		Wait(nil)
	})
}
