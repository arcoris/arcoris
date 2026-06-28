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
)

func TestPublisherZeroValueSnapshot(t *testing.T) {
	var publisher Publisher[string]

	snap := publisher.Snapshot()
	if !snap.Revision.IsZero() {
		t.Fatalf("zero publisher revision = %d, want zero", snap.Revision)
	}
	if snap.Value != "" {
		t.Fatalf("zero publisher value = %q, want empty", snap.Value)
	}
}

func TestPublisherZeroValueReadMethodsBeforePublish(t *testing.T) {
	var publisher Publisher[string]

	snap := publisher.Snapshot()
	if !snap.Revision.IsZero() {
		t.Fatalf("Snapshot().Revision = %d, want zero", snap.Revision)
	}
	if got, want := snap.Value, ""; got != want {
		t.Fatalf("Snapshot().Value = %q, want %q", got, want)
	}

	stamped := publisher.Stamped()
	if !stamped.Revision.IsZero() {
		t.Fatalf("Stamped().Revision = %d, want zero", stamped.Revision)
	}
	if !stamped.Updated.IsZero() {
		t.Fatalf("Stamped().Updated = %s, want zero", stamped.Updated)
	}
	if got, want := stamped.Value, ""; got != want {
		t.Fatalf("Stamped().Value = %q, want %q", got, want)
	}

	if got := publisher.Revision(); !got.IsZero() {
		t.Fatalf("LocalRevision() = %d, want zero", got)
	}
}

func TestPublisherStampedReturnsLatestPublishedStampedValue(t *testing.T) {
	clk := newTestClock()
	publisher := NewPublisher[string](WithClock(clk))
	clk.set(time.Unix(10, 0))

	published := publisher.PublishStamped("value")
	loaded := publisher.Stamped()
	if loaded != published {
		t.Fatalf("Stamped() = %#v, want %#v", loaded, published)
	}
}

func TestPublisherRevision(t *testing.T) {
	publisher := NewPublisher[string]()

	if got := publisher.Revision(); !got.IsZero() {
		t.Fatalf("initial LocalRevision() = %d, want zero", got)
	}

	publisher.Publish("first")
	publisher.Publish("second")

	if got, want := publisher.Revision(), LocalRevision(2); got != want {
		t.Fatalf("LocalRevision() = %d, want %d", got, want)
	}
}
