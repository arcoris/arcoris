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
	"fmt"

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// Run reflects until ctx is canceled, a sink operation fails, or the source
// violates its contract.
//
// Run panics on a nil context. A nil context would create an uncancellable
// active runtime component, which is a programmer error.
func (r *Reflector) Run(ctx context.Context) error {
	if ctx == nil {
		panic("nil context")
	}
	if err := r.beginRun(); err != nil {
		return err
	}
	defer r.endRun()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := r.runCycle(ctx)
		if err == nil {
			continue
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if isRelistRequired(err) {
			continue
		}

		return err
	}
}

// runCycle performs one ListCollection -> Replace -> Watch -> ApplyChange pass.
func (r *Reflector) runCycle(ctx context.Context) (err error) {
	read, err := r.source.ListCollection(ctx, r.collection)
	if err != nil {
		return err
	}
	if err := r.validateCollectionRead(read); err != nil {
		return err
	}
	if err := r.sink.Replace(ctx, read.Clone()); err != nil {
		return err
	}
	r.lastApplied = read.Revision()
	r.lastProgress = read.Revision()

	request, err := r.buildWatchRequest(read)
	if err != nil {
		return err
	}
	stream, err := r.source.Watch(ctx, request)
	if err != nil {
		if stream != nil {
			return errors.Join(err, stream.Close())
		}
		return err
	}
	if stream == nil {
		return sourceContractError("source returned nil stream with nil error")
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	for {
		event, err := stream.Next(ctx)
		if err != nil {
			return err
		}
		if err := r.processEvent(ctx, event); err != nil {
			return err
		}
	}
}

// validateCollectionRead checks that the source listed the exact collection the
// Reflector owns, then delegates shape validation to the API contract.
func (r *Reflector) validateCollectionRead(read storewatchapi.CollectionRead) error {
	if err := read.Validate(); err != nil {
		return err
	}
	if read.Collection() != r.collection {
		return sourceContractError("collection read belongs to %v, want %v", read.Collection(), r.collection)
	}

	return nil
}

// buildWatchRequest converts a collection read boundary into the watch request
// for the same collection. The helper localizes option translation so Run stays
// auditable.
func (r *Reflector) buildWatchRequest(read storewatchapi.CollectionRead) (objectwatch.Request, error) {
	return read.Boundary().WatchRequest(storewatchapi.WatchOptions{AllowProgress: r.options.RequestProgress})
}

// isRelistRequired reports whether Run should begin a new list-watch cycle.
func isRelistRequired(err error) bool {
	return errors.Is(err, objectwatch.ErrHistoryUnavailable) ||
		errors.Is(err, objectwatch.ErrContinuityLost) ||
		errors.Is(err, errRelistRequired)
}

// relistRequiredError keeps restart control flow private to this package while
// preserving a useful diagnostic for failed tests.
func relistRequiredError(reason objectwatch.RestartReason) error {
	return errors.Join(errRelistRequired, fmt.Errorf("restart required: %s", reason.String()))
}
