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

func TestShutDownIsIdempotentAndReportsState(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.ShutDown()
	queue.ShutDown()

	if !queue.IsShutDown() {
		t.Fatalf("IsShutDown() = false; want true")
	}
}

func TestShutDownWakesBlockedGet(t *testing.T) {
	queue := newTestQueue(t, 1)
	result := make(chan itemResult, 1)

	go func() {
		got, err := queue.Get(context.Background())
		result <- itemResult{item: got, err: err}
	}()
	queue.ShutDown()

	requireErrorIs(t, waitItem(t, result).err, ErrShutDown)
}

func TestDoneRemainsValidAfterShutDownForProcessingItem(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	queue.ShutDown()

	requireNoError(t, queue.Done(item))
	requireStats(t, queue, 0, 0)
}
