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
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// Watch opens a pull-based stream over committed changes for request.
//
// Watch validates the request before acquiring Store.mu, then computes retained
// replay and registers the live watcher while holding Store.mu. Because writes
// use the same mutex, no committed change can occur between replay selection
// and live registration.
func (s *Store) Watch(ctx context.Context, request objectwatch.Request) (objectwatch.Stream, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	if err := s.WatchCapabilities().SupportsRequest(request); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	replay, err := s.prepareReplay(ctx, request)
	if err != nil {
		return nil, err
	}

	w := &watcher{request: request}
	w.stream = newStream(s, w, s.streamBuffer)
	for _, change := range replay {
		if !w.enqueueChange(change) {
			err := streamOverflowError()
			w.terminate(err)
			return nil, err
		}
	}
	s.watchers[w] = struct{}{}

	return w.stream, nil
}

// prepareReplay computes the retained historical changes required by request.
//
// StartAfterRevision uses retained history. StartAtCurrent performs a backend
// List under Store.mu to observe the current collection boundary without
// replaying older retained changes.
func (s *Store) prepareReplay(ctx context.Context, request objectwatch.Request) ([]objectstore.Change, error) {
	switch request.Start.Mode {
	case objectwatch.StartAfterRevision:
		return s.history.replay(request.Collection, request.Start.Revision)
	case objectwatch.StartAtCurrent:
		if _, err := s.backend.List(ctx, request.Collection); err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, request.Start.Validate()
	}
}
