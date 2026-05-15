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
	"testing"
	"time"
)

func TestPublisherZeroValueSnapshot(t *testing.T) {
	var publisher Publisher[string]

	snap := publisher.Snapshot()
	if !snap.IsZeroRevision() {
		t.Fatalf("zero publisher revision = %d, want zero", snap.Revision)
	}
	if snap.Value != "" {
		t.Fatalf("zero publisher value = %q, want empty", snap.Value)
	}
}

func TestPublisherPublish(t *testing.T) {
	publisher := NewPublisher[string]()

	snap := publisher.Publish("value")
	if got, want := snap.Revision, Revision(1); got != want {
		t.Fatalf("Publish revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "value"; got != want {
		t.Fatalf("Publish value = %q, want %q", got, want)
	}

	loaded := publisher.Snapshot()
	if loaded != snap {
		t.Fatalf("Snapshot = %#v, want %#v", loaded, snap)
	}
}

func TestPublisherPublishPanicsOnRevisionOverflowWithoutPublication(t *testing.T) {
	publisher := NewPublisher[string]()
	publisher.nextRevision.Store(^uint64(0))

	requirePanicWith(t, "snapshot: revision overflow", func() {
		_ = publisher.Publish("value")
	})

	snap := publisher.Snapshot()
	if !snap.IsZeroRevision() {
		t.Fatalf("revision = %d, want zero", snap.Revision)
	}
	if got, want := snap.Value, ""; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestPublisherPublishStampedUsesClock(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(10, 0))
	publisher := NewPublisher[string](WithClock(clk))

	stamped := publisher.PublishStamped("value")
	if !stamped.Updated.Equal(time.Unix(10, 0)) {
		t.Fatalf("Updated = %s, want %s", stamped.Updated, time.Unix(10, 0))
	}
}

func TestPublisherRevision(t *testing.T) {
	publisher := NewPublisher[string]()

	if got := publisher.Revision(); !got.IsZero() {
		t.Fatalf("initial Revision() = %d, want zero", got)
	}

	publisher.Publish("first")
	publisher.Publish("second")

	if got, want := publisher.Revision(), Revision(2); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestPublisherDoesNotClonePublishedValue(t *testing.T) {
	publisher := NewPublisher[[]string]()
	val := []string{"a"}

	publisher.Publish(val)
	val[0] = "changed"

	// This test documents Publisher's immutable-publication contract. Publisher
	// does not clone values; callers must publish values they will not mutate.
	if got, want := publisher.Snapshot().Value[0], "changed"; got != want {
		t.Fatalf("published value = %q, want %q", got, want)
	}
}

func TestNewPublisherPanicsOnNilOption(t *testing.T) {
	requirePanicWith(t, "snapshot: nil option", func() {
		_ = NewPublisher[string](nil)
	})
}
