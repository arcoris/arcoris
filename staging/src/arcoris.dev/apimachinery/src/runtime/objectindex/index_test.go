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

package objectindex

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectreflectorsink"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

var _ objectreflector.Sink = (*Index)(nil)

func TestNewValidation(t *testing.T) {
	var nilFunc ExtractorFunc
	tests := []struct {
		name        string
		definitions []Definition
		target      error
	}{
		{name: "none", target: ErrInvalidDefinition},
		{
			name: "empty name",
			definitions: []Definition{
				{Name: "", Extractor: desiredExtractor()},
			},
			target: ErrInvalidDefinition,
		},
		{
			name: "duplicate name",
			definitions: []Definition{
				{Name: "desired", Extractor: desiredExtractor()},
				{Name: "desired", Extractor: nameExtractor()},
			},
			target: ErrInvalidDefinition,
		},
		{
			name: "nil extractor",
			definitions: []Definition{
				{Name: "desired"},
			},
			target: ErrNilExtractor,
		},
		{
			name: "typed nil extractor",
			definitions: []Definition{
				{Name: "desired", Extractor: nilFunc},
			},
			target: ErrNilExtractor,
		},
		{
			name: "valid",
			definitions: []Definition{
				{Name: "desired", Extractor: desiredExtractor()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := New(tt.definitions...)
			if tt.target != nil {
				requireErrorIs(t, err, tt.target)
				if index != nil {
					t.Fatalf("index = %#v; want nil", index)
				}
				return
			}

			requireNoError(t, err)
			if index == nil {
				t.Fatal("index is nil")
			}
		})
	}
}

func TestFanoutOrderMakesIndexVisibleToEnqueueMapper(t *testing.T) {
	key := testKey("task-1")
	index, err := New(Definition{Name: "worker", Extractor: fixedExtractor("worker-a")})
	requireNoError(t, err)
	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	enqueueSink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue: queue,
		Listed: objectenqueue.ListItemMapperFunc(func(_ objectstore.ListItem, emit objectenqueue.EmitFunc) error {
			keys, err := index.Lookup("worker", "worker-a")
			if err != nil {
				return err
			}
			for _, key := range keys {
				if err := emit(objectworkqueue.Item{Key: key}); err != nil {
					return err
				}
			}

			return nil
		}),
		Changed: objectenqueue.ChangedObject(),
	})
	requireNoError(t, err)
	fanout, err := objectreflectorsink.NewFanout(cache, index, enqueueSink)
	requireNoError(t, err)

	requireNoError(t, fanout.Replace(context.Background(), testRead(t, 1, testItem(key, 1, "worker-a"))))

	item, err := queue.Get(context.Background())
	requireNoError(t, err)
	if !item.Key.Equal(key) {
		t.Fatalf("queued key = %#v; want %#v", item.Key, key)
	}
}
