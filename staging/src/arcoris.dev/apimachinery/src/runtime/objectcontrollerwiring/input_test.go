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
	"testing"

	"arcoris.dev/apimachinery/runtime/objectenqueue"
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

func TestInputGettersReturnNilForNilReceiver(t *testing.T) {
	var input *Input

	if input.Cache() != nil {
		t.Fatal("Cache() is not nil")
	}
	if input.Reflector() != nil {
		t.Fatal("Reflector() is not nil")
	}
}
