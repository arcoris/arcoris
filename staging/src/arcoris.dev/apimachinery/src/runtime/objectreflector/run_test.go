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

func TestRunPanicsOnNilContext(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	defer func() {
		if recover() == nil {
			t.Fatalf("Run(nil) did not panic")
		}
	}()
	_ = reflector.Run(nil)
}

func TestRunInitialSyncBuildsWatchFromBoundary(t *testing.T) {
	read := testRead(t, 7)
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: read}},
		watchResponses: []watchResponse{{stream: terminalStream(context.Canceled)}},
	}
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())
	requireErrorIs(t, err, context.Canceled)

	if sink.replaceCount() != 1 {
		t.Fatalf("replace count = %d; want 1", sink.replaceCount())
	}
	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
	requests := source.recordedWatchRequests()
	if len(requests) != 1 {
		t.Fatalf("watch requests = %d; want 1", len(requests))
	}
	if requests[0].Collection != testCollection() {
		t.Fatalf("watch collection = %#v; want %#v", requests[0].Collection, testCollection())
	}
	if requests[0].Start.Revision != read.Revision() {
		t.Fatalf("watch revision = %s; want %s", requests[0].Start.Revision, read.Revision())
	}
	if !requests[0].AllowProgress {
		t.Fatalf("AllowProgress = false; want true")
	}
}

func TestRunWatchRequestHonorsProgressOption(t *testing.T) {
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: testRead(t, 1)}},
		watchResponses: []watchResponse{{stream: terminalStream(context.Canceled)}},
	}
	reflector, err := New(source, testCollection(), newRecordingSink(1), WithRequestProgress(false))
	requireNoError(t, err)

	err = reflector.runCycle(context.Background())
	requireErrorIs(t, err, context.Canceled)

	requests := source.recordedWatchRequests()
	if len(requests) != 1 {
		t.Fatalf("watch requests = %d; want 1", len(requests))
	}
	if requests[0].AllowProgress {
		t.Fatalf("AllowProgress = true; want false")
	}
}

func TestRunDoesNotWatchWhenReplaceFails(t *testing.T) {
	errReplace := errors.New("replace failed")
	source := &fakeListerWatcher{listResponses: []listResponse{{read: testRead(t, 1)}}}
	sink := newRecordingSink(1)
	sink.replaceErr = errReplace
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())

	requireErrorIs(t, err, errReplace)
	if len(source.recordedWatchRequests()) != 0 {
		t.Fatalf("Watch was called after Replace failure")
	}
}

func TestRunApplyChangeErrorIsFatal(t *testing.T) {
	errApply := errors.New("apply failed")
	source := &fakeListerWatcher{
		listResponses: []listResponse{{read: testRead(t, 1)}},
		watchResponses: []watchResponse{{stream: streamWithEvents(
			changedEvent(t, createdChange(t, testKey("system", 1), 2)),
		)}},
	}
	sink := newRecordingSink(1)
	sink.applyErr = errApply
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())

	requireErrorIs(t, err, errApply)
	if reflector.lastApplied != 1 {
		t.Fatalf("lastApplied = %s; want 1", reflector.lastApplied)
	}
}

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

func TestRunHandlesWatchContractFailures(t *testing.T) {
	t.Run("nil stream nil error", func(t *testing.T) {
		source := &fakeListerWatcher{
			listResponses:  []listResponse{{read: testRead(t, 1)}},
			watchResponses: []watchResponse{{}},
		}
		reflector := newTestReflector(t, source, newRecordingSink(1))

		err := reflector.runCycle(context.Background())
		requireErrorIs(t, err, ErrInvalidEvent)
	})

	t.Run("stream and error closes stream", func(t *testing.T) {
		stream := terminalStream(context.Canceled)
		sourceErr := errors.New("watch failed")
		source := &fakeListerWatcher{
			listResponses:  []listResponse{{read: testRead(t, 1)}},
			watchResponses: []watchResponse{{stream: stream, err: sourceErr}},
		}
		reflector := newTestReflector(t, source, newRecordingSink(1))

		err := reflector.runCycle(context.Background())
		requireErrorIs(t, err, sourceErr)
		if stream.closes() != 1 {
			t.Fatalf("close count = %d; want 1", stream.closes())
		}
	})
}
