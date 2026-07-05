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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestConcurrentLookupWhileReplaceRuns(t *testing.T) {
	extractor := newBlockingExtractor()
	index, err := New(Definition{Name: "worker", Extractor: fixedExtractor("worker-a")})
	requireNoError(t, err)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(testKey("task-1"), 1, "worker-a"))))
	index.definitions["worker"] = extractor

	done := make(chan error, 1)
	go func() {
		done <- index.Replace(context.Background(), testRead(t, 2, testItem(testKey("task-2"), 2, "worker-b")))
	}()
	<-extractor.entered

	for i := 0; i < 16; i++ {
		_, err := index.Lookup("worker", "worker-a")
		requireNoError(t, err)
	}
	close(extractor.release)
	requireNoError(t, <-done)
}

func TestConcurrentLookupWhileApplyChangeRuns(t *testing.T) {
	extractor := newBlockingExtractor()
	index, err := New(Definition{Name: "worker", Extractor: fixedExtractor("worker-a")})
	requireNoError(t, err)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(testKey("task-1"), 1, "worker-a"))))
	index.definitions["worker"] = extractor

	done := make(chan error, 1)
	go func() {
		done <- index.ApplyChange(context.Background(), createdChange(testKey("task-2"), 2, "worker-b"))
	}()
	<-extractor.entered

	for i := 0; i < 16; i++ {
		_, err := index.Lookup("worker", "worker-a")
		requireNoError(t, err)
	}
	close(extractor.release)
	requireNoError(t, <-done)
}

func TestConcurrentLookupForDifferentValues(t *testing.T) {
	index := newDesiredIndex(t)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 3,
		testItem(testKey("task-1"), 1, "worker-a"),
		testItem(testKey("task-2"), 2, "worker-b"),
		testItem(testKey("task-3"), 3, "worker-c"),
	)))

	var wg sync.WaitGroup
	for _, value := range []Value{"worker-a", "worker-b", "worker-c"} {
		value := value
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 32; i++ {
				_, err := index.Lookup("desired", value)
				requireNoError(t, err)
			}
		}()
	}
	wg.Wait()
}

type blockingExtractor struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func newBlockingExtractor() *blockingExtractor {
	return &blockingExtractor{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
}

func (e *blockingExtractor) Extract(objectstore.ListItem, EmitFunc) error {
	e.once.Do(func() { close(e.entered) })
	<-e.release

	return nil
}
