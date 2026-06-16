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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCacheConcurrentReaders(t *testing.T) {
	cache := mustCache(t, testListResult(40, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = cache.Get(testItems()[0].Key)
				_ = cache.Items()
				_, _ = cache.List(query)
				_ = cache.Snapshot()
				_ = cache.Revision()
				_ = cache.Len()
			}
		}()
	}
	wg.Wait()
}

func TestCacheConcurrentReadWhileApply(t *testing.T) {
	before := testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil)
	cache := mustCache(t, testListResult(1, before))

	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_, _ = cache.Get(before.Key)
				_, _ = cache.List(objectquery.Query{})
			}
		}
	}()

	current := before
	for rev := objectstore.Revision(2); rev < 40; rev++ {
		next := testItem("system", "worker-1", rev, labelsMap("env", "prod"), nil)
		change := objectstore.MustUpdatedChange(current.Key, current.State, next.State)
		requireNoError(t, cache.Apply(change))
		current = next
	}
	close(stop)
	wg.Wait()
}

func TestCacheConcurrentReadWhileReplace(t *testing.T) {
	cache := mustCache(t, testListResult(1, testItems()...))

	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_ = cache.Items()
				_, _ = cache.List(objectquery.Query{})
			}
		}
	}()

	for rev := objectstore.Revision(2); rev < 40; rev++ {
		items := []objectstore.ListItem{
			testItem("system", "worker-1", rev, labelsMap("env", "prod"), nil),
			testItem("system", "worker-2", rev+100, labelsMap("env", "qa"), nil),
		}
		requireNoError(t, cache.Replace(testListResult(rev, items...)))
	}
	close(stop)
	wg.Wait()
}
