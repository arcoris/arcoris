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

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// buildWatchRequest converts a collection read boundary into the watch request
// for the same collection.
//
// The reflector never derives a watch boundary from individual listed items. The
// collection read is the source-provided list-to-watch boundary, so all watch
// request construction is centralized here.
func (r *Reflector) buildWatchRequest(read storewatchapi.CollectionRead) (objectwatch.Request, error) {
	return read.Boundary().WatchRequest(storewatchapi.WatchOptions{AllowProgress: r.options.RequestProgress})
}

// openWatch opens a source stream for an already-built request and
// enforces objectwatch.Source's nil-stream/error contract.
//
// A source returning both a stream and an error is treated as failed; the stream
// is closed defensively before the error is returned.
func (r *Reflector) openWatch(ctx context.Context, request objectwatch.Request) (objectwatch.Stream, error) {
	stream, err := r.source.Watch(ctx, request)
	if err != nil {
		if stream != nil {
			return nil, errors.Join(
				sourceContractError("source returned stream and error"),
				err,
				stream.Close(),
			)
		}
		return nil, err
	}
	if stream == nil {
		return nil, sourceContractError("source returned nil stream with nil error")
	}

	return stream, nil
}
