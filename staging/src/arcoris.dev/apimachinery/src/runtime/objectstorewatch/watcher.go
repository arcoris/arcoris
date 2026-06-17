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

import (
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// watcher binds one stream to its original request.
//
// Store.mu owns the watcher registry. The stream owns its terminal state and
// bounded queue. Keeping this small object separate from stream avoids making
// stream aware of request matching rules.
type watcher struct {
	// request is immutable after registration and defines live fanout matching.
	request objectwatch.Request
	// stream is the pull-based queue exposed to the caller.
	stream *stream
}

// enqueueChange attempts to enqueue change without blocking writers.
//
// A false return means the stream queue was full and the caller must remove the
// watcher from the live registry. Invalid event construction is terminal for
// the stream but returns true because the stream has already been closed with a
// more precise continuity error.
func (w *watcher) enqueueChange(change objectstore.Change) bool {
	event, err := changedEvent(change)
	if err != nil {
		w.terminate(continuityError(err))
		return true
	}
	return w.stream.enqueue(event)
}

// terminate closes the watcher with err.
func (w *watcher) terminate(err error) {
	w.stream.closeWithError(err)
}
