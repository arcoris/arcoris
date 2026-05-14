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

func TestContextCauseReturnsCustomCause(t *testing.T) {
	t.Parallel()

	want := errors.New("owner stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if got := contextCause(ctx); got != want {
		t.Fatalf("contextCause() = %v, want %v", got, want)
	}
}

func TestContextCauseFallsBackToContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if got := contextCause(ctx); !errors.Is(got, context.Canceled) {
		t.Fatalf("contextCause() = %v, want context.Canceled", got)
	}
}
