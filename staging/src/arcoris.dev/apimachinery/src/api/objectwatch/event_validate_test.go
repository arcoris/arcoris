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

import "testing"

func TestChangedRejectsRevisionMismatch(t *testing.T) {
	change := watchCreatedChange(1)
	event := Event{Kind: EventChanged, Revision: 2, Change: change}

	err := event.Validate()

	requireErrorIs(t, err, ErrInvalidEvent)
	requireWatchError(t, err, ErrorReasonInvalidEvent, "watch.event.change")
}

func TestChangedRejectsRestartReason(t *testing.T) {
	change := watchCreatedChange(1)
	event := Event{
		Kind:     EventChanged,
		Revision: change.Revision,
		Change:   change,
		Restart:  RestartContinuityLost,
	}

	err := event.Validate()

	requireErrorIs(t, err, ErrInvalidEvent)
	requireErrorIs(t, err, ErrInvalidRestart)
	requireWatchError(t, err, ErrorReasonInvalidRestart, "watch.event.restart")
}

func TestProgressRejectsInvalidShapes(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{name: "zero revision", event: Event{Kind: EventProgress}},
		{name: "change", event: Event{Kind: EventProgress, Revision: 1, Change: watchCreatedChange(1)}},
		{name: "restart", event: Event{Kind: EventProgress, Revision: 1, Restart: RestartSourceReset}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			requireErrorIs(t, err, ErrInvalidEvent)
		})
	}
}

func TestRestartRequiredRejectsInvalidShapes(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{name: "zero reason", event: Event{Kind: EventRestartRequired}},
		{name: "change", event: Event{Kind: EventRestartRequired, Restart: RestartContinuityLost, Change: watchCreatedChange(1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			requireErrorIs(t, err, ErrInvalidEvent)
		})
	}
}

func TestEventRejectsUnknownKind(t *testing.T) {
	err := Event{Kind: EventKind(99)}.Validate()

	requireErrorIs(t, err, ErrInvalidEvent)
}
