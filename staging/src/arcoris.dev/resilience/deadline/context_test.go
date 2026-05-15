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

func TestWithBoundedTimeoutPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_, _ = WithBoundedTimeout(nil, time.Now(), time.Second)
	})
	requirePanic(t, panicNegativeDuration("timeout"), func() {
		_, _ = WithBoundedTimeout(context.Background(), time.Now(), -time.Nanosecond)
	})
}

func TestWithBoundedTimeoutUsesEarlierDeadline(t *testing.T) {
	t.Parallel()

	now := testNow()

	t.Run("no parent deadline uses timeout", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := WithBoundedTimeout(context.Background(), now, 5*time.Second)
		defer cancel()

		assertContextDeadline(t, ctx, now.Add(5*time.Second))
	})

	t.Run("parent deadline earlier than timeout wins", func(t *testing.T) {
		t.Parallel()

		parentDeadline := now.Add(2 * time.Second)
		parent, parentCancel := context.WithDeadline(context.Background(), parentDeadline)
		defer parentCancel()

		ctx, cancel := WithBoundedTimeout(parent, now, 5*time.Second)
		defer cancel()

		assertContextDeadline(t, ctx, parentDeadline)
	})

	t.Run("timeout earlier than parent wins", func(t *testing.T) {
		t.Parallel()

		parent, parentCancel := context.WithDeadline(context.Background(), now.Add(10*time.Second))
		defer parentCancel()

		ctx, cancel := WithBoundedTimeout(parent, now, 3*time.Second)
		defer cancel()

		assertContextDeadline(t, ctx, now.Add(3*time.Second))
	})

	t.Run("zero timeout uses observation time", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := WithBoundedTimeout(context.Background(), now, 0)
		defer cancel()

		assertContextDeadline(t, ctx, now)
	})

	t.Run("parent cancellation cancels child", func(t *testing.T) {
		t.Parallel()

		parent, parentCancel := context.WithCancel(context.Background())
		ctx, cancel := WithBoundedTimeout(parent, now, time.Hour)
		defer cancel()

		parentCancel()
		<-ctx.Done()
		if ctx.Err() != context.Canceled {
			t.Fatalf("child Err() = %v, want %v", ctx.Err(), context.Canceled)
		}
	})
}
