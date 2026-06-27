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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestConcurrentReadsWhileReplaceRuns(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "initial"))

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				_ = cache.Ready()
				_, _ = cache.Revision()
				_, _ = cache.Get(key)
				_, _ = cache.List()
				_, _ = cache.Query(objectquery.Predicate{})
				_ = cache.Len()
			}
		}()
	}

	close(start)
	for revision := objectstore.Revision(2); revision < 80; revision++ {
		requireNoError(t, cache.Replace(
			context.Background(),
			collectionRead(t, testCollection(), revision, listItem(key, revision, revision.String())),
		))
	}
	wg.Wait()
}

func TestConcurrentReadsWhileApplyChangeRuns(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "initial"))

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				_, _ = cache.Get(key)
				_, _ = cache.List()
				_, _ = cache.Query(objectquery.Predicate{})
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
}

func TestConcurrentHistoricalReadsWhileVersionRingIsOverwritten(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 3, 1, listItem(key, 1, "initial"))

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				_, _ = cache.Get(key)
				_, _ = cache.GetAt(key, 1)
				_, _ = cache.PreviousLive(key, 80)
				_, _ = cache.List()
				_, _ = cache.Query(objectquery.Predicate{})
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
}
