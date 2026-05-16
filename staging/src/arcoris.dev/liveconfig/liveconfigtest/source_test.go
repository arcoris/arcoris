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

package liveconfigtest

import (
	"sync"
	"testing"

	"arcoris.dev/snapshot"
)

func TestControlledSourcePublishesInitialSnapshot(t *testing.T) {
	cfg := NewConfig()
	src := NewControlledSource(CloneConfig(cfg))

	snap := src.Snapshot()
	RequireNonZeroRevision(t, snap)
	RequireConfigEqual(t, snap.Value, cfg)
	RequireSourceRevision(t, src, snap.Revision)
	RequireConfigSourceValue(t, src, cfg)
}

func TestControlledSourcePublishesUpdates(t *testing.T) {
	src := NewControlledSource(NewConfigVersion(1))
	prev := src.Revision()

	next := NewConfigVersion(2)
	snap := src.Publish(next)

	RequireChangedSince(t, snap, prev)
	RequireConfigEqual(t, snap.Value, next)
	RequireConfigEqual(t, src.Snapshot().Value, next)
	RequireConfigEqual(t, src.Current(), next)
}

func TestControlledSourcePublishesMany(t *testing.T) {
	src := NewEmptyControlledSource[Config]()

	snap := src.PublishMany(NewConfigVersion(1), NewConfigVersion(2), NewConfigVersion(3))

	if got, want := snap.Revision, snapshot.Revision(3); got != want {
		t.Fatalf("PublishMany revision = %d, want %d", got, want)
	}
	RequireConfigEqual(t, snap.Value, NewConfigVersion(3))
	RequireConfigEqual(t, src.Current(), NewConfigVersion(3))
}

func TestEmptyControlledSourceStartsZero(t *testing.T) {
	src := NewEmptyControlledSource[Config]()
	snap := src.Snapshot()

	if !snap.IsZeroRevision() {
		t.Fatalf("Snapshot revision = %d, want zero", snap.Revision)
	}
	if src.Revision() != snapshot.ZeroRevision {
		t.Fatalf("Revision() = %d, want zero", src.Revision())
	}
}

func TestZeroValueControlledSource(t *testing.T) {
	var src ControlledSource[Config]

	if got := src.Revision(); got != snapshot.ZeroRevision {
		t.Fatalf("zero Revision() = %d, want zero", got)
	}
	if snap := src.Snapshot(); !snap.IsZeroRevision() {
		t.Fatalf("zero Snapshot revision = %d, want zero", snap.Revision)
	}

	snap := src.Publish(NewConfigVersion(1))
	RequireNonZeroRevision(t, snap)
	RequireConfigEqual(t, src.Current(), NewConfigVersion(1))
}

func TestConfigSourceClonesPublishedValues(t *testing.T) {
	cfg := NewConfig()
	src := NewConfigSource(cfg)

	MutateConfig(&cfg)
	RequireConfigEqual(t, src.Current(), NewConfig())

	next := NewConfigVersion(2)
	PublishConfig(src, next)
	MutateConfig(&next)
	RequireConfigEqual(t, src.Current(), NewConfigVersion(2))
}

func TestConfigSourceStampedPublication(t *testing.T) {
	src := NewEmptyConfigSource()

	stamped := PublishConfigStamped(src, NewConfigVersion(1))

	RequireStampedNonZeroRevision(t, stamped)
	RequireConfigStampedValue(t, stamped, NewConfigVersion(1))
}

func TestControlledSourceSerializesConcurrentPublications(t *testing.T) {
	src := NewControlledSource(NewConfigVersion(0))

	const publications = 32
	var wg sync.WaitGroup
	wg.Add(publications)
	for idx := 0; idx < publications; idx++ {
		idx := idx
		go func() {
			defer wg.Done()
			src.Publish(NewConfigVersion(idx + 1))
		}()
	}
	wg.Wait()

	wantRev := snapshot.Revision(publications + 1)
	if src.Revision() != wantRev {
		t.Fatalf("Revision() = %d, want %d", src.Revision(), wantRev)
	}
}
