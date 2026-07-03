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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestListItemMapperFuncNil(t *testing.T) {
	var mapper ListItemMapperFunc
	err := mapper.MapListItem(testListItem(1, 1, "listed"), func(objectworkqueue.Item) error {
		return nil
	})

	requireErrorIs(t, err, ErrNilListItemMapper)
}

func TestListItemMapperFuncDelegates(t *testing.T) {
	item := testListItem(1, 1, "listed")
	wantErr := errors.New("delegate failed")
	called := false
	mapper := ListItemMapperFunc(func(got objectstore.ListItem, emit EmitFunc) error {
		called = true
		requireListItem(t, got, item)
		if emit == nil {
			t.Fatalf("emit is nil")
		}
		return wantErr
	})

	err := mapper.MapListItem(item, func(objectworkqueue.Item) error { return nil })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	if !called {
		t.Fatalf("mapper function was not called")
	}
}

type pointerListItemMapper struct{}

func (*pointerListItemMapper) MapListItem(objectstore.ListItem, EmitFunc) error {
	return nil
}
