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

package contextstop

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCausePreservesOrdinaryContextSentinels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  context.Context
		want error
	}{
		{
			name: "canceled",
			ctx:  canceledContext(),
			want: context.Canceled,
		},
		{
			name: "deadline",
			ctx:  deadlineContext(),
			want: context.DeadlineExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := Cause(test.ctx, test.ctx.Err()); !errors.Is(got, test.want) {
				t.Fatalf("Cause() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCausePreservesCustomCancelCauseAndSentinel(t *testing.T) {
	t.Parallel()

	want := errors.New("owner stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	got := Cause(ctx, ctx.Err())
	if !errors.Is(got, want) {
		t.Fatalf("Cause() = %v, want custom cause", got)
	}
	if !errors.Is(got, context.Canceled) {
		t.Fatalf("Cause() = %v, want context.Canceled", got)
	}
}

func TestCausePreservesCustomDeadlineCauseAndSentinel(t *testing.T) {
	t.Parallel()

	want := errors.New("deadline budget")
	ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(-time.Second), want)
	defer cancel()

	got := Cause(ctx, ctx.Err())
	if !errors.Is(got, want) {
		t.Fatalf("Cause() = %v, want custom cause", got)
	}
	if !errors.Is(got, context.DeadlineExceeded) {
		t.Fatalf("Cause() = %v, want context.DeadlineExceeded", got)
	}
}

func TestCauseReturnsCauseWhenItAlreadyMatchesContextErr(t *testing.T) {
	t.Parallel()

	ctx := canceledContext()
	got := Cause(ctx, ctx.Err())

	if got != context.Canceled {
		t.Fatalf("Cause() = %v, want context.Canceled", got)
	}
}

func TestCauseReturnsCustomCauseWhenErrIsNil(t *testing.T) {
	t.Parallel()

	want := errors.New("owner stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if got := Cause(ctx, nil); got != want {
		t.Fatalf("Cause(ctx, nil) = %v, want %v", got, want)
	}
}

func TestCauseReturnsNilWhenNoCauseOrErrExists(t *testing.T) {
	t.Parallel()

	if got := Cause(context.Background(), nil); got != nil {
		t.Fatalf("Cause(active context, nil) = %v, want nil", got)
	}
}

func TestCauseFallsBackToErrForCustomContext(t *testing.T) {
	t.Parallel()

	want := errors.New("custom context stopped")
	ctx := contextWithoutCause{err: want}

	if got := Cause(ctx, want); got != want {
		t.Fatalf("Cause() = %v, want %v", got, want)
	}
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}

func deadlineContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	cancel()

	return ctx
}

// contextWithoutCause exercises the fallback path for non-standard contexts.
type contextWithoutCause struct {
	err error
}

func (ctx contextWithoutCause) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (ctx contextWithoutCause) Done() <-chan struct{} {
	return nil
}

func (ctx contextWithoutCause) Err() error {
	return ctx.err
}

func (ctx contextWithoutCause) Value(any) any {
	return nil
}
