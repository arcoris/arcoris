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
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

func TestPublisherPublish(t *testing.T) {
	publisher := NewPublisher[string]()

	snap := publisher.Publish("value")
	if got, want := snap.Revision, LocalRevision(1); got != want {
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

func TestPublisherZeroValuePublishRemainsValid(t *testing.T) {
	var publisher Publisher[string]

	stamped := publisher.PublishStamped("value")
	if got, want := stamped.Revision, LocalRevision(1); got != want {
		t.Fatalf("PublishStamped revision = %d, want %d", got, want)
	}
	if got, want := publisher.Snapshot().Value, "value"; got != want {
		t.Fatalf("Snapshot().Value = %q, want %q", got, want)
	}
}

func TestPublisherPublishPanicsOnRevisionOverflowWithoutPublication(t *testing.T) {
	publisher := NewPublisher[string]()
	publisher.nextRevision = ^LocalRevision(0)

	panicassert.RequireMessage(t, "snapshot: local revision overflow", func() {
		_ = publisher.Publish("value")
	})

	snap := publisher.Snapshot()
	if !snap.Revision.IsZero() {
		t.Fatalf("revision = %d, want zero", snap.Revision)
	}
	if got, want := snap.Value, ""; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestPublisherPublishStampedPanicsOnRevisionOverflowBeforeStore(t *testing.T) {
	publisher := NewPublisher[string]()
	publisher.Publish("old")

	before := publisher.Snapshot()
	publisher.nextRevision = ^LocalRevision(0)

	panicassert.RequireMessage(t, "snapshot: local revision overflow", func() {
		_ = publisher.PublishStamped("new")
	})

	after := publisher.Snapshot()
	if after != before {
		t.Fatalf("Snapshot() = %#v, want %#v", after, before)
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

func TestPublisherRevisionAdvancesWhenClockMovesBackward(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(20, 0))
	publisher := NewPublisher[string](WithClock(clk))

	first := publisher.PublishStamped("first")

	clk.set(time.Unix(10, 0))
	second := publisher.PublishStamped("second")

	if got, want := first.Revision, LocalRevision(1); got != want {
		t.Fatalf("first revision = %d, want %d", got, want)
	}
	if got, want := second.Revision, LocalRevision(2); got != want {
		t.Fatalf("second revision = %d, want %d", got, want)
	}
	if !second.Updated.Equal(time.Unix(10, 0)) {
		t.Fatalf("second updated = %s, want %s", second.Updated, time.Unix(10, 0))
	}
}

func TestPublisherDoesNotClonePublishedMap(t *testing.T) {
	publisher := NewPublisher[map[string]string]()
	val := map[string]string{
		"key": "value",
	}

	publisher.Publish(val)
	val["key"] = "changed"

	if got, want := publisher.Snapshot().Value["key"], "changed"; got != want {
		t.Fatalf("published map value = %q, want %q", got, want)
	}
}

func TestPublisherDoesNotClonePublishedNestedMutableValue(t *testing.T) {
	publisher := NewPublisher[mutableReadModel]()
	val := mutableReadModelValue("name", "attr", "tag")

	publisher.Publish(val)
	mutateMutableReadModel(&val, "changed-name", "changed-attr", "changed-tag")

	got := publisher.Snapshot().Value
	assertMutableReadModel(t, got, mutableReadModelValue("changed-name", "changed-attr", "changed-tag"))
}

func TestPublisherDoesNotClonePublishedValue(t *testing.T) {
	publisher := NewPublisher[[]string]()
	val := []string{"a"}

	publisher.Publish(val)
	val[0] = "changed"

	if got, want := publisher.Snapshot().Value[0], "changed"; got != want {
		t.Fatalf("published value = %q, want %q", got, want)
	}
}
