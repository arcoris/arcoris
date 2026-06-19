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
	"sync"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// Reflector actively synchronizes one object collection from a ListerWatcher
// source into a Sink.
//
// A Reflector is intended to be run by at most one goroutine at a time.
// Concurrent Run calls return ErrAlreadyRunning. After Run returns, the same
// Reflector may be run again.
type Reflector struct {
	source     storewatchapi.ListerWatcher
	collection objectstore.ListRequest
	sink       Sink
	options    Options

	mu      sync.Mutex
	running bool

	lastApplied  objectstore.Revision
	lastProgress objectstore.Revision
}

// New validates dependencies and constructs a Reflector for one collection.
func New(
	source storewatchapi.ListerWatcher,
	collection objectstore.ListRequest,
	sink Sink,
	options ...Option,
) (*Reflector, error) {
	if source == nil {
		return nil, ErrNilSource
	}
	if sink == nil {
		return nil, ErrNilSink
	}
	if err := objectstore.ValidateListRequest(collection); err != nil {
		return nil, err
	}
	cfg, err := applyOptions(options)
	if err != nil {
		return nil, err
	}

	return &Reflector{
		source:     source,
		collection: collection,
		sink:       sink,
		options:    cfg,
	}, nil
}

// beginRun marks r active if no other Run call is active.
func (r *Reflector) beginRun() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return ErrAlreadyRunning
	}
	r.running = true

	return nil
}

// endRun releases the single-run guard.
func (r *Reflector) endRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.running = false
}
