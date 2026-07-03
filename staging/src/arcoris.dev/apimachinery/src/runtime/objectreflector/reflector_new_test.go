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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

func TestNewValidation(t *testing.T) {
	source := &fakeListerWatcher{}
	sink := newRecordingSink(1)

	tests := []struct {
		name       string
		source     storewatchapi.ListerWatcher
		collection objectstore.ListRequest
		sink       Sink
		options    []Option
		target     error
	}{
		{name: "nil source", source: nil, collection: testCollection(), sink: sink, target: ErrNilSource},
		{name: "nil sink", source: source, collection: testCollection(), sink: nil, target: ErrNilSink},
		{name: "typed nil sink", source: source, collection: testCollection(), sink: (*recordingSink)(nil), target: ErrNilSink},
		{name: "invalid collection", source: source, collection: objectstore.ListRequest{}, sink: sink, target: objectstore.ErrInvalidListRequest},
		{name: "nil option", source: source, collection: testCollection(), sink: sink, options: []Option{nil}, target: ErrInvalidOption},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.source, tt.collection, tt.sink, tt.options...)
			requireErrorIs(t, err, tt.target)
		})
	}
}

func TestNewSuccess(t *testing.T) {
	reflector, err := New(&fakeListerWatcher{}, testCollection(), newRecordingSink(1), WithRequestProgress(false))
	requireNoError(t, err)

	if reflector.options.RequestProgress {
		t.Fatalf("RequestProgress = true; want false")
	}
}
