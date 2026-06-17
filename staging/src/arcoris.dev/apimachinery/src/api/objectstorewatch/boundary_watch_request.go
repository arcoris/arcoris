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

package objectstorewatch

import "arcoris.dev/apimachinery/api/objectwatch"

// WatchRequest converts b into an objectwatch request for the same collection.
//
// The generated request always uses StartAfterRevision, including for zero
// boundaries. StartAtCurrent is intentionally not used for collection-read
// continuation because it would drop historical catch-up semantics.
func (b Boundary) WatchRequest(options WatchOptions) (objectwatch.Request, error) {
	if err := b.Validate(); err != nil {
		return objectwatch.Request{}, err
	}

	start, err := objectwatch.AfterRevision(b.revision)
	if err != nil {
		return objectwatch.Request{}, errorFor("boundary.watch_request", ErrorReasonInvalidBoundary, ErrInvalidBoundary, err)
	}

	return objectwatch.Request{
		Collection:    b.collection,
		Start:         start,
		AllowProgress: options.AllowProgress,
	}, nil
}
