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

type recordingObserver struct {
	ctx   context.Context
	event Event
	calls int
}

func (o *recordingObserver) ObserveRetry(ctx context.Context, event Event) {
	o.ctx = ctx
	o.event = event
	o.calls++
}

func TestObserverReceivesContextAndEvent(t *testing.T) {
	ctx := context.WithValue(context.Background(), observerTestContextKey{}, "value")
	event := Event{
		Kind: EventAttemptStart,
		Attempt: Attempt{
			Number:    1,
			StartedAt: time.Unix(1, 0),
		},
	}

	observer := &recordingObserver{}
	observer.ObserveRetry(ctx, event)

	if observer.calls != 1 {
		t.Fatalf("observer calls = %d, want 1", observer.calls)
	}
	if observer.ctx != ctx {
		t.Fatalf("observer context was not forwarded")
	}
	if observer.event != event {
		t.Fatalf("observer event = %+v, want %+v", observer.event, event)
	}
}

type observerTestContextKey struct{}
