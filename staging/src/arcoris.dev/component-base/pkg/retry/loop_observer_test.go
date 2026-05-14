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

func TestEmitRetryEventCallsObserversInRegistrationOrder(t *testing.T) {
	var order []string
	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}
	config := configOf(
		WithObserverFunc(func(_ context.Context, got Event) {
			order = append(order, "first")
			if got != event {
				t.Fatalf("first observer event = %+v, want %+v", got, event)
			}
		}),
		WithObserverFunc(func(_ context.Context, got Event) {
			order = append(order, "second")
			if got != event {
				t.Fatalf("second observer event = %+v, want %+v", got, event)
			}
		}),
	)

	emitRetryEvent(context.Background(), config, event)

	want := []string{"first", "second"}
	if len(order) != len(want) {
		t.Fatalf("observer order len = %d, want %d: %v", len(order), len(want), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("observer order[%d] = %q, want %q", i, order[i], want[i])
		}
	}
}

func TestRetryExecutionEmitUsesConfiguredObservers(t *testing.T) {
	recorder := &retryObserverRecorder{}
	event := Event{
		Kind:    EventAttemptStart,
		Attempt: retryTestAttempt(1),
	}
	execution := &retryExecution{
		config: configOf(WithObserver(recorder)),
	}

	execution.emit(context.Background(), event)

	if len(recorder.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(recorder.events))
	}
	if recorder.events[0] != event {
		t.Fatalf("events[0] = %+v, want %+v", recorder.events[0], event)
	}
}
