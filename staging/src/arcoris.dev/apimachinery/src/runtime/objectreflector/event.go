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
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// processEvent validates and routes one watch event. It never forwards
// progress or restart control events to Sink.
func (r *Reflector) processEvent(ctx context.Context, event objectwatch.Event) error {
	event = event.Clone()
	if err := event.Validate(); err != nil {
		return invalidEventError(err)
	}

	switch event.Kind {
	case objectwatch.EventChanged:
		return r.processChanged(ctx, event)
	case objectwatch.EventProgress:
		return r.processProgress(event)
	case objectwatch.EventRestartRequired:
		return relistRequiredError(event.Restart)
	default:
		return sourceContractError("unknown event kind %s", event.Kind.String())
	}
}

// processChanged applies one committed object transition after validating
// collection membership and strict revision ordering.
func (r *Reflector) processChanged(ctx context.Context, event objectwatch.Event) error {
	change := event.Change.Clone()
	if err := objectstore.ValidateChange(change); err != nil {
		return invalidEventError(err)
	}
	if !objectstore.ChangeMatchesListRequest(change, r.collection) {
		return changeOutsideCollectionError(fmt.Errorf("change key %s is outside reflected collection", change.Key.String()))
	}
	if !r.lastApplied.Before(change.Revision) {
		return nonMonotonicRevisionError(
			fmt.Errorf("change revision %s is not after last applied revision %s", change.Revision, r.lastApplied),
		)
	}

	if err := r.sink.ApplyChange(ctx, change.Clone()); err != nil {
		return err
	}
	r.lastApplied = change.Revision
	if r.lastProgress.Before(change.Revision) {
		r.lastProgress = change.Revision
	}

	return nil
}

// processProgress records a source progress boundary without mutating Sink.
func (r *Reflector) processProgress(event objectwatch.Event) error {
	if event.Revision.Before(r.lastApplied) {
		return nonMonotonicRevisionError(
			fmt.Errorf("progress revision %s is before last applied revision %s", event.Revision, r.lastApplied),
		)
	}
	if event.Revision.Before(r.lastProgress) {
		return nonMonotonicRevisionError(
			fmt.Errorf("progress revision %s is before last progress revision %s", event.Revision, r.lastProgress),
		)
	}
	r.lastProgress = event.Revision

	return nil
}
