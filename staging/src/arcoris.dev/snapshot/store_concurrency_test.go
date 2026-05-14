/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

import (
	"sync"
	"testing"
)

func TestStoreConcurrentSnapshotAndReplace(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				store.Replace([]string{"value"})
				_ = id
			}
		}(i)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				_ = store.Snapshot()
			}
		}()
	}

	wg.Wait()
}

func TestStoreConcurrentSnapshotAndUpdate(t *testing.T) {
	store := NewStore(0, Identity[int])

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				store.Update(func(v int) int { return v + 1 })
			}
		}()
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				_ = store.Snapshot()
			}
		}()
	}

	wg.Wait()

	if got, want := store.Snapshot().Value, 1600; got != want {
		t.Fatalf("final value = %d, want %d", got, want)
	}
}
