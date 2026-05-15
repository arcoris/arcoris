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

package liveconfig

import (
	"sync"
	"testing"

	"arcoris.dev/snapshot"
)

func TestConcurrentSnapshotAndApply(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 0})

	const writers = 32
	const readers = 32
	const readsPerReader = 64

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(writers + readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < readsPerReader; j++ {
				_ = h.Snapshot()
				_ = h.Stamped()
				_ = h.Revision()
				_ = h.LastError()
			}
		}()
	}

	for i := 0; i < writers; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			_, err := h.Apply(testConfig{Name: "writer", Limit: i})
			if err != nil {
				t.Errorf("Apply() error = %v", err)
			}
		}()
	}

	close(start)
	wg.Wait()

	want := snapshot.Revision(1 + writers)
	if got := h.Revision(); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestConcurrentInvalidApplyKeepsLastGood(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	const writers = 32
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(writers)

	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			<-start
			_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})
		}()
	}

	close(start)
	wg.Wait()

	cur := h.Snapshot()
	if cur.Revision != prev.Revision {
		t.Fatalf("Revision() = %d, want %d", cur.Revision, prev.Revision)
	}
	if got, want := cur.Value.Name, "initial"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if h.LastError() == nil {
		t.Fatal("LastError() = nil, want rejected apply error")
	}
}
