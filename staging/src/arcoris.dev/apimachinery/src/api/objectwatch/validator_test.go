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
			validator := mustValidator(t, Start{Mode: StartAfterRevision, Revision: 10})
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
	validator := mustValidator(t, AtCurrent())

	requireNoError(t, validator.Accept(mustChangedEvent(t, 1)))
}

func TestValidatorBookmarkProgress(t *testing.T) {
	tests := []struct {
		name    string
		events  []Event
		wantErr bool
	}{
		{name: "bookmark equal start accepted", events: []Event{mustBookmarkEvent(t, 10)}},
		{name: "bookmark before start rejected", events: []Event{mustBookmarkEvent(t, 9)}, wantErr: true},
		{name: "changed then same bookmark accepted", events: []Event{mustChangedEvent(t, 11), mustBookmarkEvent(t, 11)}},
		{name: "changed then earlier bookmark rejected", events: []Event{mustChangedEvent(t, 11), mustBookmarkEvent(t, 10)}, wantErr: true},
		{name: "bookmark then same changed rejected", events: []Event{mustBookmarkEvent(t, 20), mustChangedEvent(t, 20)}, wantErr: true},
		{name: "bookmark then later changed accepted", events: []Event{mustBookmarkEvent(t, 20), mustChangedEvent(t, 21)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := mustValidator(t, Start{Mode: StartAfterRevision, Revision: 10})
			err := acceptAll(&validator, tt.events)
			if tt.wantErr {
				requireErrorIs(t, err, ErrContinuityLost)
				return
			}
			requireNoError(t, err)
		})
	}
}

func TestValidatorRestartRequiredIsTerminal(t *testing.T) {
	validator := mustValidator(t, AtCurrent())
	restart := mustRestartEvent(t, RestartContinuityLost, 0)

	requireNoError(t, validator.Accept(restart))
	err := validator.Accept(mustChangedEvent(t, 1))

	requireErrorIs(t, err, ErrClosed)
}

func TestValidatorRejectsInvalidRestartRequired(t *testing.T) {
	validator := mustValidator(t, AtCurrent())

	err := validator.Accept(Event{Kind: EventRestartRequired})

	requireErrorIs(t, err, ErrInvalidEvent)
}

func TestValidatorRejectsInvalidEvent(t *testing.T) {
	validator := mustValidator(t, AtCurrent())

	err := validator.Accept(Event{})

	requireErrorIs(t, err, ErrInvalidEvent)
}

func TestValidatorRejectsNilReceiver(t *testing.T) {
	var validator *Validator

	err := validator.Accept(mustChangedEvent(t, 1))

	requireErrorIs(t, err, ErrContinuityLost)
}

func mustValidator(t *testing.T, start Start) Validator {
	t.Helper()

	validator, err := NewValidator(start)
	requireNoError(t, err)
	return validator
}

func mustChangedEvent(t *testing.T, revision objectstore.Revision) Event {
	t.Helper()

	event, err := Changed(watchCreatedChange(revision))
	requireNoError(t, err)
	return event
}

func mustBookmarkEvent(t *testing.T, revision objectstore.Revision) Event {
	t.Helper()

	event, err := Bookmark(revision)
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
