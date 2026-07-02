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

func TestErrorsAreDistinctSentinels(t *testing.T) {
	sentinels := []error{
		ErrNilQueue,
		ErrNilMapper,
		ErrNilPredicate,
		ErrNilEmit,
		ErrInvalidHandler,
	}

	for i, left := range sentinels {
		for j, right := range sentinels {
			if i == j {
				requireErrorIs(t, left, right)
				continue
			}
			if errors.Is(left, right) {
				t.Fatalf("errors.Is(%v, %v) = true; want false", left, right)
			}
		}
	}
}

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}

func requireItem(t testing.TB, got objectworkqueue.Item, want objectworkqueue.Item) {
	t.Helper()

	if !got.Key.Equal(want.Key) {
		t.Fatalf("item = %s; want %s", got.Key, want.Key)
	}
}

func requireChange(t testing.TB, got objectstore.Change, want objectstore.Change) {
	t.Helper()

	if got.Kind != want.Kind || got.Revision != want.Revision || !got.Key.Equal(want.Key) {
		t.Fatalf("change identity = %#v; want %#v", got, want)
	}
}
