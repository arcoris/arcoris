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

package objectwatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestNewValidatorRejectsInvalidRequest(t *testing.T) {
	_, err := NewValidator(Request{Start: AtCurrent()})

	requireErrorIs(t, err, ErrInvalidRequest)
}

func TestValidatorStartAfterRevisionChangedOrdering(t *testing.T) {
	tests := []struct {
		name    string
		events  []Event
		wantErr bool
	}{
		{name: "changed after start", events: []Event{mustChangedEvent(t, 11)}},
		{name: "changed equal start rejected", events: []Event{mustChangedEvent(t, 10)}, wantErr: true},
		{name: "changed before start rejected", events: []Event{mustChangedEvent(t, 9)}, wantErr: true},
		{name: "increasing accepted", events: []Event{mustChangedEvent(t, 11), mustChangedEvent(t, 12)}},
		{name: "decreasing rejected", events: []Event{mustChangedEvent(t, 12), mustChangedEvent(t, 11)}, wantErr: true},
		{name: "duplicate rejected", events: []Event{mustChangedEvent(t, 12), mustChangedEvent(t, 12)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := mustValidator(t, watchRequest(Start{Mode: StartAfterRevision, Revision: 10}, false))
			err := acceptAll(&validator, tt.events)
			if tt.wantErr {
				requireErrorIs(t, err, ErrContinuityLost)
				return
			}
			requireNoError(t, err)
		})
	}
}

func TestValidatorStartAtCurrentAcceptsFirstNonZeroChange(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	requireNoError(t, validator.Accept(mustChangedEvent(t, 1)))
}

func TestValidatorAcceptsMatchingResourceAndScope(t *testing.T) {
	validator := mustValidator(t, watchNamespaceRequest(AtCurrent(), "system", false))

	requireNoError(t, validator.Accept(mustChangedEventForKey(t, watchKeyInNamespace("system", "main"), 1)))
}

func TestValidatorRejectsDifferentResourceAndCloses(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	err := validator.Accept(mustChangedEventForKey(t, watchKeyFor(otherWatchResource(), "system", "main"), 1))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEvent(t, 2))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorRejectsDifferentNamespaceForNamespaceScopeAndCloses(t *testing.T) {
	validator := mustValidator(t, watchNamespaceRequest(AtCurrent(), "system", false))

	err := validator.Accept(mustChangedEventForKey(t, watchKeyInNamespace("other", "main"), 1))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEventForKey(t, watchKeyInNamespace("system", "next"), 2))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorAllNamespacesAcceptsMultipleNamespaces(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	requireNoError(t, validator.Accept(mustChangedEventForKey(t, watchKeyInNamespace("system", "main"), 1)))
	requireNoError(t, validator.Accept(mustChangedEventForKey(t, watchKeyInNamespace("other", "next"), 2)))
}

func TestValidatorProgress(t *testing.T) {
	tests := []struct {
		name    string
		events  []Event
		wantErr bool
	}{
		{name: "progress equal start accepted", events: []Event{mustProgressEvent(t, 10)}},
		{name: "progress before start rejected", events: []Event{mustProgressEvent(t, 9)}, wantErr: true},
		{name: "changed then same progress accepted", events: []Event{mustChangedEvent(t, 11), mustProgressEvent(t, 11)}},
		{name: "changed then earlier progress rejected", events: []Event{mustChangedEvent(t, 11), mustProgressEvent(t, 10)}, wantErr: true},
		{name: "progress then same changed rejected", events: []Event{mustProgressEvent(t, 20), mustChangedEvent(t, 20)}, wantErr: true},
		{name: "progress then later changed accepted", events: []Event{mustProgressEvent(t, 20), mustChangedEvent(t, 21)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := mustValidator(t, watchRequest(Start{Mode: StartAfterRevision, Revision: 10}, true))
			err := acceptAll(&validator, tt.events)
			if tt.wantErr {
				requireErrorIs(t, err, ErrContinuityLost)
				return
			}
			requireNoError(t, err)
		})
	}
}

func TestValidatorRejectsProgressWhenNotAllowedAndCloses(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	err := validator.Accept(mustProgressEvent(t, 1))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEvent(t, 2))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorContinuityLossClosesChangedSequence(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	requireNoError(t, validator.Accept(mustChangedEvent(t, 20)))
	err := validator.Accept(mustChangedEvent(t, 19))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEvent(t, 21))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorContinuityLossClosesProgressSequence(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), true))

	requireNoError(t, validator.Accept(mustProgressEvent(t, 20)))
	err := validator.Accept(mustProgressEvent(t, 19))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEvent(t, 21))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorRestartRequiredIsTerminal(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))
	restart := mustRestartEvent(t, RestartContinuityLost, 0)

	requireNoError(t, validator.Accept(restart))
	err := validator.Accept(mustChangedEvent(t, 1))

	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorRestartBeforeProgressCloses(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	requireNoError(t, validator.Accept(mustChangedEvent(t, 20)))
	err := validator.Accept(mustRestartEvent(t, RestartContinuityLost, 19))
	requireErrorIs(t, err, ErrContinuityLost)

	err = validator.Accept(mustChangedEvent(t, 21))
	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorRejectsInvalidRestartRequired(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	err := validator.Accept(Event{Kind: EventRestartRequired})

	requireErrorIs(t, err, ErrInvalidEvent)
}

func TestValidatorRejectsInvalidEventWithoutClosing(t *testing.T) {
	validator := mustValidator(t, watchRequest(AtCurrent(), false))

	err := validator.Accept(Event{})
	requireErrorIs(t, err, ErrInvalidEvent)

	requireNoError(t, validator.Accept(mustChangedEvent(t, 1)))
}

func TestValidatorRejectsNilReceiver(t *testing.T) {
	var validator *Validator

	err := validator.Accept(mustChangedEvent(t, 1))

	requireErrorIs(t, err, ErrContinuityLost)
}

func mustValidator(t *testing.T, request Request) Validator {
	t.Helper()

	validator, err := NewValidator(request)
	requireNoError(t, err)
	return validator
}

func mustChangedEvent(t *testing.T, revision objectstore.Revision) Event {
	t.Helper()

	event, err := Changed(watchCreatedChange(revision))
	requireNoError(t, err)
	return event
}

func mustChangedEventForKey(t *testing.T, key objectstore.Key, revision objectstore.Revision) Event {
	t.Helper()

	event, err := Changed(watchCreatedChangeForKey(key, revision))
	requireNoError(t, err)
	return event
}

func mustProgressEvent(t *testing.T, revision objectstore.Revision) Event {
	t.Helper()

	event, err := Progress(revision)
	requireNoError(t, err)
	return event
}

func mustRestartEvent(t *testing.T, reason RestartReason, revision objectstore.Revision) Event {
	t.Helper()

	event, err := RestartRequired(reason, revision)
	requireNoError(t, err)
	return event
}

func acceptAll(validator *Validator, events []Event) error {
	for _, event := range events {
		if err := validator.Accept(event); err != nil {
			return err
		}
	}

	return nil
}
