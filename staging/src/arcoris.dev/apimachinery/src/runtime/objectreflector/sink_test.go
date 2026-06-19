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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

var _ Sink = (*recordingSink)(nil)

// recordingSink captures reflected reads and changes for behavior assertions.
//
// The helper clones every value it stores and returns cloned slices so tests can
// assert reflector ownership boundaries without depending on a concrete cache.
type recordingSink struct {
	mu sync.Mutex

	replaceErr error
	applyErr   error

	reads   []storewatchapi.CollectionRead
	changes []objectstore.Change

	replaceCh chan storewatchapi.CollectionRead
	changeCh  chan objectstore.Change
}

func newRecordingSink(buffer int) *recordingSink {
	return &recordingSink{
		replaceCh: make(chan storewatchapi.CollectionRead, buffer),
		changeCh:  make(chan objectstore.Change, buffer),
	}
}

func (s *recordingSink) Replace(_ context.Context, read storewatchapi.CollectionRead) error {
	if s.replaceErr != nil {
		return s.replaceErr
	}

	read = read.Clone()
	s.mu.Lock()
	s.reads = append(s.reads, read)
	s.mu.Unlock()

	select {
	case s.replaceCh <- read:
	default:
	}

	return nil
}

func (s *recordingSink) ApplyChange(_ context.Context, change objectstore.Change) error {
	if s.applyErr != nil {
		return s.applyErr
	}

	change = change.Clone()
	s.mu.Lock()
	s.changes = append(s.changes, change)
	s.mu.Unlock()

	select {
	case s.changeCh <- change:
	default:
	}

	return nil
}

func (s *recordingSink) replaceCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.reads)
}

func (s *recordingSink) changeCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.changes)
}

func (s *recordingSink) recordedChanges() []objectstore.Change {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]objectstore.Change, len(s.changes))
	for i, change := range s.changes {
		out[i] = change.Clone()
	}

	return out
}

func TestRecordingSinkReturnsDetachedChanges(t *testing.T) {
	sink := newRecordingSink(1)
	change := createdChange(t, testKey("system", 1), 2)

	requireNoError(t, sink.ApplyChange(context.Background(), change))
	recorded := sink.recordedChanges()
	recorded[0].After.Revision = 99

	again := sink.recordedChanges()
	if again[0].After.Revision == 99 {
		t.Fatalf("recordedChanges exposed mutable sink state")
	}
}
