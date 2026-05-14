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
