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

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
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

func TestRunRelistPolicyReceivesAttemptsAndCauses(t *testing.T) {
	firstErr := objectwatch.ContinuityLost(errors.New("first gap"))
	secondErr := objectwatch.HistoryUnavailable(errors.New("second gap"))
	policy := &recordingRelistPolicy{}
	thirdStream := waitingStream()
	source := &fakeListerWatcher{
		listResponses: []listResponse{
			{read: testRead(t, 1)},
			{read: testRead(t, 2)},
			{read: testRead(t, 3)},
		},
		watchResponses: []watchResponse{
			{stream: terminalStream(firstErr)},
			{stream: terminalStream(secondErr)},
			{stream: thirdStream},
		},
	}
	sink := newRecordingSink(3)
	reflector, err := New(source, testCollection(), sink, WithRelistPolicy(policy))
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- reflector.Run(ctx)
	}()

	waitRead(t, sink.replaceCh)
	waitRead(t, sink.replaceCh)
	waitRead(t, sink.replaceCh)
	cancel()

	requireErrorIs(t, <-done, context.Canceled)
	calls := policy.recordedCalls()
	if len(calls) != 2 {
		t.Fatalf("policy calls = %d; want 2", len(calls))
	}
	if calls[0].attempt != 1 || calls[1].attempt != 2 {
		t.Fatalf("attempts = %d, %d; want 1, 2", calls[0].attempt, calls[1].attempt)
	}
	requireErrorIs(t, calls[0].cause, objectwatch.ErrContinuityLost)
	requireErrorIs(t, calls[1].cause, objectwatch.ErrHistoryUnavailable)
}

func TestRunRelistPolicyCalledForRestartRequired(t *testing.T) {
	policy := &recordingRelistPolicy{}
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
	reflector, err := New(source, testCollection(), sink, WithRelistPolicy(policy))
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- reflector.Run(ctx)
	}()

	waitRead(t, sink.replaceCh)
	waitRead(t, sink.replaceCh)
	cancel()

	requireErrorIs(t, <-done, context.Canceled)
	calls := policy.recordedCalls()
	if len(calls) != 1 {
		t.Fatalf("policy calls = %d; want 1", len(calls))
	}
	if calls[0].attempt != 1 {
		t.Fatalf("attempt = %d; want 1", calls[0].attempt)
	}
	if !isRelistRequired(calls[0].cause) {
		t.Fatalf("policy cause is not relist-required: %v", calls[0].cause)
	}
}

func TestRunRelistPolicyErrorStopsRun(t *testing.T) {
	policyErr := errors.New("policy stopped")
	policy := &recordingRelistPolicy{err: policyErr}
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: testRead(t, 1)}},
		watchResponses: []watchResponse{{stream: terminalStream(objectwatch.ContinuityLost(errors.New("gap")))}},
	}
	reflector, err := New(source, testCollection(), newRecordingSink(1), WithRelistPolicy(policy))
	requireNoError(t, err)

	err = reflector.Run(context.Background())

	requireErrorIs(t, err, policyErr)
	calls := policy.recordedCalls()
	if len(calls) != 1 {
		t.Fatalf("policy calls = %d; want 1", len(calls))
	}
}

func TestRunRelistPolicyWaitHonorsContextCancellation(t *testing.T) {
	policy := newBlockingRelistPolicy()
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: testRead(t, 1)}, {read: testRead(t, 2)}},
		watchResponses: []watchResponse{{stream: terminalStream(objectwatch.ContinuityLost(errors.New("gap")))}},
	}
	sink := newRecordingSink(2)
	reflector, err := New(source, testCollection(), sink, WithRelistPolicy(policy))
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- reflector.Run(ctx)
	}()

	waitRead(t, sink.replaceCh)
	<-policy.entered
	cancel()

	requireErrorIs(t, <-done, context.Canceled)
	if sink.replaceCount() != 1 {
		t.Fatalf("replace count = %d; want 1", sink.replaceCount())
	}
}

func TestRunDoesNotCallRelistPolicyForFatalErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*fakeListerWatcher, *recordingSink, error)
	}{
		{
			name: "replace error",
			setup: func() (*fakeListerWatcher, *recordingSink, error) {
				errReplace := errors.New("replace failed")
				sink := newRecordingSink(1)
				sink.replaceErr = errReplace
				return &fakeListerWatcher{
					listResponses: []listResponse{{read: testRead(t, 1)}},
				}, sink, errReplace
			},
		},
		{
			name: "apply change error",
			setup: func() (*fakeListerWatcher, *recordingSink, error) {
				errApply := errors.New("apply failed")
				sink := newRecordingSink(1)
				sink.applyErr = errApply
				return &fakeListerWatcher{
					listResponses: []listResponse{{read: testRead(t, 1)}},
					watchResponses: []watchResponse{{stream: streamWithEvents(
						changedEvent(t, createdChange(t, testKey("system", 1), 2)),
					)}},
				}, sink, errApply
			},
		},
		{
			name: "invalid event",
			setup: func() (*fakeListerWatcher, *recordingSink, error) {
				return &fakeListerWatcher{
					listResponses: []listResponse{{read: testRead(t, 1)}},
					watchResponses: []watchResponse{{stream: streamWithEvents(objectwatch.Event{
						Kind:     objectwatch.EventChanged,
						Revision: 2,
					})}},
				}, newRecordingSink(1), ErrInvalidEvent
			},
		},
		{
			name: "source contract violation",
			setup: func() (*fakeListerWatcher, *recordingSink, error) {
				return &fakeListerWatcher{
					listResponses:  []listResponse{{read: testRead(t, 1)}},
					watchResponses: []watchResponse{{}},
				}, newRecordingSink(1), ErrSourceContractViolation
			},
		},
		{
			name: "invalid collection read",
			setup: func() (*fakeListerWatcher, *recordingSink, error) {
				return &fakeListerWatcher{
					listResponses: []listResponse{{}},
				}, newRecordingSink(1), storewatchapi.ErrInvalidCollectionRead
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, sink, target := tt.setup()
			policy := &recordingRelistPolicy{}
			reflector, err := New(source, testCollection(), sink, WithRelistPolicy(policy))
			requireNoError(t, err)

			err = reflector.Run(context.Background())

			requireErrorIs(t, err, target)
			if calls := policy.recordedCalls(); len(calls) != 0 {
				t.Fatalf("policy calls = %d; want 0", len(calls))
			}
		})
	}
}
