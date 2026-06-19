// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectreflector

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestRunRelistsOnRestartRequired(t *testing.T) {
	secondStream := waitingStream()
	source := &fakeListerWatcher{
		listResponses: []listResponse{
			{read: testRead(t, 1)},
			{read: testRead(t, 2)},
		},
		watchResponses: []watchResponse{
			{stream: streamWithEvents(restartEvent(t))},
			{stream: secondStream},
		},
	}
	sink := newRecordingSink(2)
	reflector := newTestReflector(t, source, sink)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)

	go func() {
		done <- reflector.Run(ctx)
	}()

	waitRead(t, sink.replaceCh)
	waitRead(t, sink.replaceCh)
	cancel()

	requireErrorIs(t, <-done, context.Canceled)
	requests := source.recordedWatchRequests()
	if len(requests) != 2 {
		t.Fatalf("watch requests = %d; want 2", len(requests))
	}
	if requests[1].Start.Revision != 2 {
		t.Fatalf("second watch revision = %s; want 2", requests[1].Start.Revision)
	}
}

func TestRunRelistsOnContinuityErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "history unavailable", err: objectwatch.HistoryUnavailable(errors.New("old history"))},
		{name: "continuity lost", err: objectwatch.ContinuityLost(errors.New("gap"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secondStream := waitingStream()
			source := &fakeListerWatcher{
				listResponses: []listResponse{
					{read: testRead(t, 1)},
					{read: testRead(t, 3)},
				},
				watchResponses: []watchResponse{
					{stream: terminalStream(tt.err)},
					{stream: secondStream},
				},
			}
			sink := newRecordingSink(2)
			reflector := newTestReflector(t, source, sink)
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)

			go func() {
				done <- reflector.Run(ctx)
			}()

			waitRead(t, sink.replaceCh)
			waitRead(t, sink.replaceCh)
			cancel()

			requireErrorIs(t, <-done, context.Canceled)
			if sink.replaceCount() != 2 {
				t.Fatalf("replace count = %d; want 2", sink.replaceCount())
			}
		})
	}
}
