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

package snapshot

import (
	"runtime"
	"sync"
	"testing"
)

func TestPublisherConcurrentPublishesReachExpectedRevision(t *testing.T) {
	publisher := NewPublisher[int]()

	const publishers = 16
	const publishesPerPublisher = 100

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(publishers)

	for id := 0; id < publishers; id++ {
		id := id
		go func() {
			defer wg.Done()
			<-start
			for n := 0; n < publishesPerPublisher; n++ {
				publisher.Publish(id*publishesPerPublisher + n)
			}
		}()
	}

	close(start)
	wg.Wait()

	want := LocalRevision(publishers * publishesPerPublisher)
	if got := publisher.Revision(); got != want {
		t.Fatalf("LocalRevision() = %d, want %d", got, want)
	}
}

func TestPublisherReadersDoNotObserveRevisionRollback(t *testing.T) {
	publisher := NewPublisher[int]()

	const publishers = 8
	const publishesPerPublisher = 2000

	start := make(chan struct{})
	done := make(chan struct{})
	readerReady := make(chan struct{})
	rollback := make(chan string, 1)

	go func() {
		defer close(rollback)
		close(readerReady)

		var max LocalRevision
		for {
			select {
			case <-done:
				return
			default:
			}

			rev := publisher.Snapshot().Revision
			if rev.IsZero() {
				continue
			}
			if rev < max {
				rollback <- "observed revision rollback"
				return
			}
			max = rev
		}
	}()

	<-readerReady

	var wg sync.WaitGroup
	wg.Add(publishers)
	for id := 0; id < publishers; id++ {
		id := id
		go func() {
			defer wg.Done()
			<-start
			for n := 0; n < publishesPerPublisher; n++ {
				publisher.Publish(id*publishesPerPublisher + n)
				if n%64 == 0 {
					runtime.Gosched()
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	close(done)

	if msg := <-rollback; msg != "" {
		t.Fatal(msg)
	}

	want := LocalRevision(publishers * publishesPerPublisher)
	if got := publisher.Revision(); got != want {
		t.Fatalf("LocalRevision() = %d, want %d", got, want)
	}
}

func TestPublisherConcurrentSnapshotAndPublish(t *testing.T) {
	publisher := NewPublisher[int]()

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				publisher.Publish(id*100 + n)
			}
		}(i)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				_ = publisher.Snapshot()
			}
		}()
	}

	wg.Wait()

	if publisher.Revision().IsZero() {
		t.Fatal("publisher revision is still zero after concurrent publishes")
	}
}

func TestPublisherConcurrentStampedAndPublish(t *testing.T) {
	publisher := NewPublisher[int]()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				publisher.PublishStamped(id*100 + n)
			}
		}(i)
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				_ = publisher.Stamped()
			}
		}()
	}

	wg.Wait()
}
