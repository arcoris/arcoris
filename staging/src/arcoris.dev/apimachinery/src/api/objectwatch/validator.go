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

// Validator enforces local request and event ordering rules for one stream.
//
// Validator does not read stores, recover streams, or prove source correctness
// beyond the events it observes. It is intended for consumers and tests that
// want to fail closed on invalid event sequences. Validator is not safe for
// concurrent use.
type Validator struct {
	// request is the structural collection contract the stream is expected to
	// satisfy.
	request Request
	// progress is the highest revision boundary accepted so far. EventChanged
	// must advance strictly beyond it; EventProgress may equal or advance it.
	progress objectstore.Revision
	// closed is set after EventRestartRequired or any observed continuity
	// violation because the current stream can no longer be trusted.
	closed bool
}

// NewValidator constructs a request-aware stream validator.
//
// StartAfterRevision initializes progress to the requested revision so the
// first changed event must be strictly newer. StartAtCurrent starts with zero
// progress and accepts the first non-zero changed or progress event.
func NewValidator(request Request) (Validator, error) {
	if err := request.Validate(); err != nil {
		return Validator{}, err
	}

	validator := Validator{request: request}
	if request.Start.Mode == StartAfterRevision {
		validator.progress = request.Start.Revision
	}

	return validator, nil
}

// Accept validates event against the request and stream progress observed so far.
//
// Malformed event shapes are returned as ErrInvalidEvent without closing the
// validator. Valid events that violate request scope or revision continuity
// return ErrContinuityLost and close the validator.
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
	case EventProgress:
		return v.acceptProgress(event)
	case EventRestartRequired:
		return v.acceptRestart(event)
	default:
		return invalidEventError("watch.event", fmt.Errorf("event kind %s is invalid", event.Kind.String()))
	}
}

// acceptChanged requires strict progress advancement.
func (v *Validator) acceptChanged(event Event) error {
	if !objectstore.ChangeMatchesListRequest(event.Change, v.request.Collection) {
		return v.closeWithContinuityLoss(fmt.Errorf("changed event does not match requested collection"))
	}
	if !v.progress.Before(event.Revision) {
		return v.closeWithContinuityLoss(fmt.Errorf(
			"changed revision %s is not after progress %s",
			event.Revision,
			v.progress,
		))
	}
	v.progress = event.Revision

	return nil
}

// acceptProgress records a monotonic progress boundary.
func (v *Validator) acceptProgress(event Event) error {
	if !v.request.AllowProgress {
		return v.closeWithContinuityLoss(fmt.Errorf("progress event is not allowed by request"))
	}
	if event.Revision.Before(v.progress) {
		return v.closeWithContinuityLoss(fmt.Errorf(
			"progress revision %s is before progress %s",
			event.Revision,
			v.progress,
		))
	}
	v.progress = event.Revision

	return nil
}

// acceptRestart records terminal restart-required state.
func (v *Validator) acceptRestart(event Event) error {
	if !event.Revision.IsZero() && event.Revision.Before(v.progress) {
		return v.closeWithContinuityLoss(fmt.Errorf(
			"restart revision %s is before progress %s",
			event.Revision,
			v.progress,
		))
	}
	if v.progress.Before(event.Revision) {
		v.progress = event.Revision
	}
	v.closed = true

	return nil
}

// closeWithContinuityLoss marks the stream terminal before returning the
// continuity diagnostic so later Accept calls fail closed with ErrClosed.
func (v *Validator) closeWithContinuityLoss(cause error) error {
	v.closed = true
	return continuityLostError(cause)
}
