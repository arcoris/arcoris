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
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// New validates dependencies and constructs a Reflector for one collection.
//
// The returned Reflector is inactive until Run is called. The zero Reflector is
// intentionally unsupported because the source, collection, sink, and options
// are all part of the synchronization contract.
func New(
	source storewatchapi.ListerWatcher,
	collection objectstore.ListRequest,
	sink Sink,
	options ...Option,
) (*Reflector, error) {
	if source == nil {
		return nil, ErrNilSource
	}
	if isNilSink(sink) {
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
