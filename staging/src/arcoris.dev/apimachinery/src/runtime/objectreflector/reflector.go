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

	// mu protects the single-run guard only. The synchronization loop itself is
	// single-threaded after beginRun succeeds, so revision fields do not need a
	// lock while Run owns the reflector.
	mu      sync.Mutex
	running bool

	// lastApplied is the latest committed change revision successfully applied to
	// Sink, or the CollectionRead revision after Replace.
	lastApplied objectstore.Revision
	// lastProgress is the latest progress boundary observed from the stream. It is
	// diagnostic state only; objectwatch.Validator owns request-aware stream
	// ordering checks.
	lastProgress objectstore.Revision
}
