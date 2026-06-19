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

// processChanged applies one committed object transition.
//
// The source is already responsible for scoping the watch stream, but the
// reflector verifies collection membership defensively. Accepting an
// out-of-collection change would corrupt the sink while making the source bug
// harder to diagnose.
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
