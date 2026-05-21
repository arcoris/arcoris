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
)

func TestWithObserverAppendsObservers(t *testing.T) {
	first := ObserverFunc(func(context.Context, Event) {})
	second := ObserverFunc(func(context.Context, Event) {})

	cfg := configOf(
		WithObserver(first),
		WithObserver(second),
	)

	if len(cfg.observers) != 2 {
		t.Fatalf("observers len = %d, want 2", len(cfg.observers))
	}
	if cfg.observers[0] == nil {
		t.Fatalf("first observer is nil")
	}
	if cfg.observers[1] == nil {
		t.Fatalf("second observer is nil")
	}
}

func TestWithObserverPanicsOnNilObserver(t *testing.T) {
	expectPanic(t, panicNilObserver, func() {
		_ = WithObserver(nil)
	})
}

func TestWithObserverFunc(t *testing.T) {
	called := false

	cfg := configOf(WithObserverFunc(func(context.Context, Event) {
		called = true
	}))

	if len(cfg.observers) != 1 {
		t.Fatalf("observers len = %d, want 1", len(cfg.observers))
	}

	cfg.observers[0].ObserveRetry(context.Background(), Event{})

	if !called {
		t.Fatalf("observer function was not called")
	}
}

func TestWithObserverFuncPanicsOnNilFunction(t *testing.T) {
	expectPanic(t, panicNilObserverFunc, func() {
		_ = WithObserverFunc(nil)
	})
}
