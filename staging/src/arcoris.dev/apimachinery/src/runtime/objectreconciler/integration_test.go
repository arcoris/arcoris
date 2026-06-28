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

package objectreconciler

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectcache"
)

func TestReconcileOnceWithReadyObjectCache(t *testing.T) {
	cache := readyCache(t, 11)
	var got Snapshot

	result := ReconcileOnce(
		context.Background(),
		cache,
		ReconcileFunc(func(ctx context.Context, snap Snapshot) Result {
			got = snap
			return Success()
		}),
	)

	if result.Failed() {
		t.Fatalf("result = %#v; want success", result)
	}
	if got.Revision != 11 || got.View.Revision() != 11 {
		t.Fatalf("snapshot = %#v; want revision 11", got)
	}
	if got.View.Len() != 0 {
		t.Fatalf("view Len() = %d; want 0", got.View.Len())
	}
}

func TestReconcileOnceWithNotReadyObjectCache(t *testing.T) {
	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)

	result := ReconcileOnce(context.Background(), cache, ReconcileFunc(func(context.Context, Snapshot) Result {
		t.Fatal("reconciler should not be called")
		return Success()
	}))

	requireErrorIs(t, result.Err, objectcache.ErrNotReady)
}
