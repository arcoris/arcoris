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

package objectenqueue

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestPackageContractsCompile(t *testing.T) {
	var _ Enqueuer = enqueuerFunc(func(context.Context, objectworkqueue.Item) error {
		return nil
	})
	var _ Mapper = MapperFunc(func(objectstore.Change, EmitFunc) error {
		return nil
	})
	var _ Predicate = PredicateFunc(func(objectstore.Change) (bool, error) {
		return true, nil
	})

	handler, err := New(enqueuerFunc(func(context.Context, objectworkqueue.Item) error {
		return nil
	}), ChangedObject())
	requireNoError(t, err)
	if handler == nil {
		t.Fatalf("handler is nil")
	}
}

type enqueuerFunc func(context.Context, objectworkqueue.Item) error

func (f enqueuerFunc) Add(ctx context.Context, item objectworkqueue.Item) error {
	return f(ctx, item)
}
