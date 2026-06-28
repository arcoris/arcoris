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

package objectcache

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestConcurrentReadSnapshotWhileCacheMutates(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "initial"))

	start := make(chan struct{})
	errs := make(chan error, 1)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				snap, err := cache.ReadSnapshot()
				if err != nil {
					reportConcurrentError(errs, err)
					return
				}
				if snap.Revision != snap.Value.Revision() {
					reportConcurrentError(errs, fmt.Errorf("snapshot revision %s != view revision %s", snap.Revision, snap.Value.Revision()))
					return
				}
				if list := snap.Value.List(); list.Revision != snap.Revision {
					reportConcurrentError(errs, fmt.Errorf("list revision %s != snapshot revision %s", list.Revision, snap.Revision))
					return
				}
				_, _ = snap.Value.Get(key)
				_ = snap.Value.Query(objectquery.Predicate{})
			}
		}()
	}

	close(start)
	before := testState(key, 1, "initial")
	for revision := objectstore.Revision(2); revision < 80; revision++ {
		after := testState(key, revision, revision.String())
		if revision%2 == 0 {
			requireNoError(t, cache.ApplyChange(
				context.Background(),
				objectstore.MustUpdatedChange(key, before, after),
			))
		} else {
			requireNoError(t, cache.Replace(
				context.Background(),
				collectionRead(t, testCollection(), revision, objectstore.ListItem{Key: key, State: after}),
			))
		}
		before = after
	}
	wg.Wait()
	requireNoConcurrentError(t, errs)
}

func TestConcurrentViewReadsWhileCacheMutates(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "initial"))
	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	view := snap.Value

	start := make(chan struct{})
	errs := make(chan error, 1)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				got, err := view.Get(key)
				if err != nil {
					reportConcurrentError(errs, err)
					return
				}
				if !got.Found || got.Revision != 1 || desiredStringFromState(got.State) != "initial" {
					reportConcurrentError(errs, fmt.Errorf("view Get() = %#v; want initial object at revision 1", got))
					return
				}
				if list := view.List(); list.Revision != 1 || len(list.Items) != 1 {
					reportConcurrentError(errs, fmt.Errorf("view List() = %#v; want one item at revision 1", list))
					return
				}
				if query := view.Query(objectquery.Predicate{}); query.Revision != 1 || len(query.Items) != 1 {
					reportConcurrentError(errs, fmt.Errorf("view Query() = %#v; want one item at revision 1", query))
					return
				}
			}
		}()
	}

	close(start)
	before := testState(key, 1, "initial")
	for revision := objectstore.Revision(2); revision < 80; revision++ {
		after := testState(key, revision, revision.String())
		requireNoError(t, cache.ApplyChange(
			context.Background(),
			objectstore.MustUpdatedChange(key, before, after),
		))
		before = after
	}
	wg.Wait()
	requireNoConcurrentError(t, errs)
}

func reportConcurrentError(errs chan<- error, err error) {
	select {
	case errs <- err:
	default:
	}
}

func requireNoConcurrentError(t testing.TB, errs <-chan error) {
	t.Helper()

	select {
	case err := <-errs:
		t.Fatal(err)
	default:
	}
}

func desiredStringFromState(state objectstore.State) string {
	desired, _ := state.Object.Desired.AsString()
	return desired
}
