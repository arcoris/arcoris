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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

type listResponse struct {
	read storewatchapi.CollectionRead
	err  error
}

type watchResponse struct {
	stream objectwatch.Stream
	err    error
}

// fakeListerWatcher is a scripted source for reflector cycle tests.
//
// Each call consumes one queued response. Unexpected calls return an ordinary
// error so tests fail at the behavior boundary instead of blocking.
type fakeListerWatcher struct {
	mu sync.Mutex

	listResponses  []listResponse
	watchResponses []watchResponse

	listRequests  []objectstore.ListRequest
	watchRequests []objectwatch.Request
}

func (f *fakeListerWatcher) ListCollection(
	_ context.Context,
	collection objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.listRequests = append(f.listRequests, collection)
	if len(f.listResponses) == 0 {
		return storewatchapi.CollectionRead{}, errors.New("unexpected ListCollection call")
	}
	response := f.listResponses[0]
	f.listResponses = f.listResponses[1:]

	return response.read, response.err
}

func (f *fakeListerWatcher) Watch(_ context.Context, request objectwatch.Request) (objectwatch.Stream, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.watchRequests = append(f.watchRequests, request)
	if len(f.watchResponses) == 0 {
		return nil, errors.New("unexpected Watch call")
	}
	response := f.watchResponses[0]
	f.watchResponses = f.watchResponses[1:]

	return response.stream, response.err
}

func (f *fakeListerWatcher) recordedWatchRequests() []objectwatch.Request {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]objectwatch.Request(nil), f.watchRequests...)
}

func TestRunWatchRequestHonorsProgressOption(t *testing.T) {
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: testRead(t, 1)}},
		watchResponses: []watchResponse{{stream: terminalStream(context.Canceled)}},
	}
	reflector, err := New(source, testCollection(), newRecordingSink(1), WithRequestProgress(false))
	requireNoError(t, err)

	err = reflector.runCycle(context.Background())
	requireErrorIs(t, err, context.Canceled)

	requests := source.recordedWatchRequests()
	if len(requests) != 1 {
		t.Fatalf("watch requests = %d; want 1", len(requests))
	}
	if requests[0].AllowProgress {
		t.Fatalf("AllowProgress = true; want false")
	}
}

func TestOpenWatchHandlesWatchContractFailures(t *testing.T) {
	t.Run("nil stream nil error", func(t *testing.T) {
		source := &fakeListerWatcher{
			listResponses:  []listResponse{{read: testRead(t, 1)}},
			watchResponses: []watchResponse{{}},
		}
		reflector := newTestReflector(t, source, newRecordingSink(1))

		err := reflector.runCycle(context.Background())
		requireErrorIs(t, err, ErrInvalidEvent)
	})

	t.Run("stream and error closes stream", func(t *testing.T) {
		stream := terminalStream(context.Canceled)
		sourceErr := errors.New("watch failed")
		source := &fakeListerWatcher{
			listResponses:  []listResponse{{read: testRead(t, 1)}},
			watchResponses: []watchResponse{{stream: stream, err: sourceErr}},
		}
		reflector := newTestReflector(t, source, newRecordingSink(1))

		err := reflector.runCycle(context.Background())
		requireErrorIs(t, err, sourceErr)
		if stream.closes() != 1 {
			t.Fatalf("close count = %d; want 1", stream.closes())
		}
	})
}
