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

package objectworkqueue

import (
	"context"
	"testing"
)

func TestZeroStatsIsEmptyDiagnosticObservation(t *testing.T) {
	var stats Stats

	if stats.Queued != 0 || stats.Processing != 0 || stats.Capacity != 0 || stats.ShutDown {
		t.Fatalf("zero Stats = %#v; want empty", stats)
	}
}

func TestStatsReportsCapacityQueuedProcessingAndShutDown(t *testing.T) {
	queue := newTestQueue(t, 3)
	requireNoError(t, queue.Add(context.Background(), testItem(1)))
	requireNoError(t, queue.Add(context.Background(), testItem(2)))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	queue.ShutDown()

	stats := queue.Stats()
	if stats.Capacity != 3 {
		t.Fatalf("Capacity = %d; want 3", stats.Capacity)
	}
	if stats.Queued != 1 {
		t.Fatalf("Queued = %d; want 1", stats.Queued)
	}
	if stats.Processing != 1 {
		t.Fatalf("Processing = %d; want 1", stats.Processing)
	}
	if !stats.ShutDown {
		t.Fatalf("ShutDown = false; want true")
	}
}

func TestProcessingItemsCountAgainstCapacity(t *testing.T) {
	queue := newTestQueue(t, 1)
	requireNoError(t, queue.Add(context.Background(), testItem(1)))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	requireErrorIs(t, queue.TryAdd(testItem(2)), ErrFull)
}

func TestTrackedCountNeverExceedsCapacity(t *testing.T) {
	queue := newTestQueue(t, 2)
	requireNoError(t, queue.TryAdd(testItem(1)))
	requireNoError(t, queue.TryAdd(testItem(2)))
	requireErrorIs(t, queue.TryAdd(testItem(3)), ErrFull)

	stats := queue.Stats()
	if stats.Queued+stats.Processing > stats.Capacity {
		t.Fatalf("stats = %#v; queued+processing exceeds capacity", stats)
	}
}
