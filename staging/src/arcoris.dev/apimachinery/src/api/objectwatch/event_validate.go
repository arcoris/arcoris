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

// Validate checks the event kind-specific invariant.
//
// It deliberately delegates committed-transition validation to objectstore for
// EventChanged so objectwatch does not duplicate storage-level change rules.
func (e Event) Validate() error {
	if !e.Kind.IsValid() {
		return invalidEventError("watch.event", fmt.Errorf("event kind %s is invalid", e.Kind.String()))
	}

	switch e.Kind {
	case EventChanged:
		return validateChangedEvent(e)
	case EventProgress:
		return validateProgressEvent(e)
	case EventRestartRequired:
		return validateRestartEvent(e)
	default:
		return invalidEventError("watch.event", fmt.Errorf("event kind %s is invalid", e.Kind.String()))
	}
}

// validateChangedEvent checks the mutation-event contract: one valid
// objectstore.Change, matching Revision, and no restart payload.
func validateChangedEvent(e Event) error {
	if err := objectstore.ValidateChange(e.Change); err != nil {
		return invalidEventError("watch.event.change", err)
	}
	if e.Revision != e.Change.Revision {
		return invalidEventError(
			"watch.event.change",
			fmt.Errorf("event revision %s differs from change revision %s", e.Revision, e.Change.Revision),
		)
	}
	if e.Restart != 0 {
		return invalidEventRestartError(fmt.Errorf("changed event has restart reason %s", e.Restart.String()))
	}

	return nil
}

// validateProgressEvent checks the progress-event contract: a non-zero
// Revision and no mutation or restart payload.
func validateProgressEvent(e Event) error {
	if e.Revision.IsZero() {
		return invalidEventError("watch.event.progress", fmt.Errorf("progress revision is zero"))
	}
	if !e.Change.IsZero() {
		return invalidEventError("watch.event.change", fmt.Errorf("progress event carries change"))
	}
	if e.Restart != 0 {
		return invalidEventRestartError(fmt.Errorf("progress event has restart reason %s", e.Restart.String()))
	}

	return nil
}

// validateRestartEvent checks the terminal restart-required contract: a valid
// reason, optional progress Revision, and no mutation payload.
func validateRestartEvent(e Event) error {
	if err := e.Restart.Validate(); err != nil {
		return invalidEventRestartError(err)
	}
	if !e.Change.IsZero() {
		return invalidEventError("watch.event.change", fmt.Errorf("restart event carries change"))
	}

	return nil
}

// invalidEventRestartError preserves both event and restart classifications so
// callers can handle malformed events broadly or invalid restart payloads
// specifically.
func invalidEventRestartError(cause error) error {
	return objectWatchError(
		"watch.event.restart",
		ErrorReasonInvalidRestart,
		[]error{ErrInvalidEvent, ErrInvalidRestart},
		cause,
	)
}
