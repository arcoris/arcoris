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

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// processChanged applies one committed object transition.
//
// objectwatch.Validator has already checked request membership and revision
// ordering. This method keeps only a defensive change validation before handing
// a detached payload to Sink.
func (r *Reflector) processChanged(ctx context.Context, event objectwatch.Event) error {
	change := event.Change.Clone()
	if err := objectstore.ValidateChange(change); err != nil {
		return invalidEventError(err)
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
