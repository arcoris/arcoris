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

package objectcontrollerwiring

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestNewInputAcceptsValidConfig(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)

	input, err := newInput(InputConfig{
		Source:     &runTestListerWatcher{},
		Collection: runTestCollection(),
		Listed:     objectenqueue.ListedObject(),
		Changed:    objectenqueue.ChangedObject(),
	}, queue)
	requireNoError(t, err)

	if input.Cache() == nil {
		t.Fatal("Cache() is nil")
	}
	if input.Reflector() == nil {
		t.Fatal("Reflector() is nil")
	}
}

func TestNewInputAcceptsIndexes(t *testing.T) {
	index := newDesiredObjectIndex(t)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)

	input, err := newInput(InputConfig{
		Source:     &runTestListerWatcher{},
		Collection: runTestCollection(),
		Listed:     objectenqueue.ListedObject(),
		Changed:    objectenqueue.ChangedObject(),
		Indexes:    []*objectindex.Index{index},
	}, queue)
	requireNoError(t, err)

	indexes := input.Indexes()
	if len(indexes) != 1 || indexes[0] != index {
		t.Fatalf("Indexes() = %#v; want configured index", indexes)
	}
}

func TestNewInputRejectsNilIndex(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)

	_, err = newInput(InputConfig{
		Source:     &runTestListerWatcher{},
		Collection: runTestCollection(),
		Listed:     objectenqueue.ListedObject(),
		Changed:    objectenqueue.ChangedObject(),
		Indexes:    []*objectindex.Index{nil},
	}, queue)

	if !errors.Is(err, ErrNilIndex) {
		t.Fatalf("error = %v; want ErrNilIndex", err)
	}
}

func TestInputIndexesReturnsDetachedSlice(t *testing.T) {
	index := newDesiredObjectIndex(t)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	input, err := newInput(InputConfig{
		Source:     &runTestListerWatcher{},
		Collection: runTestCollection(),
		Listed:     objectenqueue.ListedObject(),
		Changed:    objectenqueue.ChangedObject(),
		Indexes:    []*objectindex.Index{index},
	}, queue)
	requireNoError(t, err)

	indexes := input.Indexes()
	indexes[0] = nil

	if input.Indexes()[0] == nil {
		t.Fatal("Indexes() exposed internal slice")
	}
}

func TestInputGettersReturnNilForNilReceiver(t *testing.T) {
	var input *Input

	if input.Cache() != nil {
		t.Fatal("Cache() is not nil")
	}
	if input.Reflector() != nil {
		t.Fatal("Reflector() is not nil")
	}
	if input.Indexes() != nil {
		t.Fatal("Indexes() is not nil")
	}
}
