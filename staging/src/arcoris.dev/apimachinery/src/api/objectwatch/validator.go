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
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Validator enforces local event ordering rules for one stream.
//
// Validator does not read stores, recover streams, or prove source correctness
// beyond the events it observes. It is intended for consumers and tests that
// want to fail closed on invalid event sequences.
type Validator struct {
	// progress is the highest revision boundary accepted so far. EventChanged
	// must advance strictly beyond it; EventBookmark may equal or advance it.
	progress objectstore.Revision
	// closed is set after EventRestartRequired because restart is terminal for
	// the current stream.
	closed bool
}

// NewValidator constructs a sequence validator for start.
//
// StartAfterRevision initializes progress to the requested revision so the
// first changed event must be strictly newer. StartAtCurrent starts with zero
// progress and accepts the first non-zero changed or bookmark event.
func NewValidator(start Start) (Validator, error) {
	if err := start.Validate(); err != nil {
		return Validator{}, err
	}

	validator := Validator{}
	if start.Mode == StartAfterRevision {
		validator.progress = start.Revision
	}

	return validator, nil
}

// Accept validates event against the stream progress observed so far.
//
// Accept mutates Validator only after the event validates and passes ordering
// checks. EventRestartRequired is accepted once and then closes the validator.
func (v *Validator) Accept(event Event) error {
	if v == nil {
		return continuityLostError(fmt.Errorf("nil validator"))
	}
	if v.closed {
		return closedError(fmt.Errorf("stream is terminal"))
	}
	if err := event.Validate(); err != nil {
		return err
	}

	switch event.Kind {
	case EventChanged:
		return v.acceptChanged(event)
	case EventBookmark:
		return v.acceptBookmark(event)
	case EventRestartRequired:
		return v.acceptRestart(event)
	default:
		return invalidEventError("watch.event", fmt.Errorf("event kind %s is invalid", event.Kind.String()))
	}
}

// acceptChanged requires strict progress advancement.
func (v *Validator) acceptChanged(event Event) error {
	if !v.progress.Before(event.Revision) {
		return continuityLostError(
			fmt.Errorf("changed revision %s is not after progress %s", event.Revision, v.progress),
		)
	}
	v.progress = event.Revision

	return nil
}

// acceptBookmark records a monotonic progress boundary.
func (v *Validator) acceptBookmark(event Event) error {
	if event.Revision.Before(v.progress) {
		return continuityLostError(
			fmt.Errorf("bookmark revision %s is before progress %s", event.Revision, v.progress),
		)
	}
	v.progress = event.Revision

	return nil
}

// acceptRestart records terminal restart-required state.
func (v *Validator) acceptRestart(event Event) error {
	if !event.Revision.IsZero() && event.Revision.Before(v.progress) {
		return continuityLostError(
			fmt.Errorf("restart revision %s is before progress %s", event.Revision, v.progress),
		)
	}
	if v.progress.Before(event.Revision) {
		v.progress = event.Revision
	}
	v.closed = true

	return nil
}
