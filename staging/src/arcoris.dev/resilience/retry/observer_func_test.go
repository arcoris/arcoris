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
	"context"
	"testing"
	"time"
)

func TestObserverFuncImplementsObserver(t *testing.T) {
	var _ Observer = ObserverFunc(func(context.Context, Event) {})
}

func TestObserverFuncPanicsWhenNil(t *testing.T) {
	var observer ObserverFunc

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("ObserverFunc.ObserveRetry did not panic")
		}
		if recovered != panicNilObserverFunc {
			t.Fatalf("panic = %v, want %q", recovered, panicNilObserverFunc)
		}
	}()

	observer.ObserveRetry(context.Background(), Event{})
}

func TestObserverFuncDelegatesContextAndEvent(t *testing.T) {
	ctx := context.WithValue(context.Background(), observerFuncTestContextKey{}, "value")
	event := Event{
		Kind: EventAttemptStart,
		Attempt: Attempt{
			Number:    1,
			StartedAt: time.Unix(1, 0),
		},
	}

	var gotCtx context.Context
	var gotEvent Event
	calls := 0

	observer := ObserverFunc(func(ctx context.Context, event Event) {
		gotCtx = ctx
		gotEvent = event
		calls++
	})

	observer.ObserveRetry(ctx, event)

	if calls != 1 {
		t.Fatalf("ObserverFunc calls = %d, want 1", calls)
	}
	if gotCtx != ctx {
		t.Fatalf("ObserverFunc did not forward context")
	}
	if gotEvent != event {
		t.Fatalf("ObserverFunc event = %+v, want %+v", gotEvent, event)
	}
}

func TestObserverFuncDoesNotValidateEvent(t *testing.T) {
	invalidEvent := Event{
		Kind: EventKind(255),
	}

	called := false
	var gotEvent Event

	observer := ObserverFunc(func(_ context.Context, event Event) {
		called = true
		gotEvent = event
	})

	observer.ObserveRetry(context.Background(), invalidEvent)

	if !called {
		t.Fatalf("ObserverFunc did not call wrapped function")
	}
	if gotEvent != invalidEvent {
		t.Fatalf("ObserverFunc event = %+v, want %+v", gotEvent, invalidEvent)
	}
}

type observerFuncTestContextKey struct{}
