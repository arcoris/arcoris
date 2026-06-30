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

import "testing"

func TestSignalNotEmptyLockedSkipsChannelWithoutWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notEmpty
	queue.signalNotEmptyLocked()
	replacement := queue.notEmpty
	queue.mu.Unlock()

	if replacement == ch {
		return
	}
	t.Fatalf("notEmpty channel was replaced without waiters")
}

func TestSignalNotEmptyLockedClosesAndReplacesChannelWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notEmpty
	queue.notEmptyWaiters = 1
	queue.signalNotEmptyLocked()
	replacement := queue.notEmpty
	queue.mu.Unlock()

	requireClosed(t, ch)
	if replacement == ch {
		t.Fatalf("notEmpty channel was not replaced")
	}
}

func TestSignalNotFullLockedSkipsChannelWithoutWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notFull
	queue.signalNotFullLocked()
	replacement := queue.notFull
	queue.mu.Unlock()

	if replacement == ch {
		return
	}
	t.Fatalf("notFull channel was replaced without waiters")
}

func TestSignalNotFullLockedClosesAndReplacesChannelWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notFull
	queue.notFullWaiters = 1
	queue.signalNotFullLocked()
	replacement := queue.notFull
	queue.mu.Unlock()

	requireClosed(t, ch)
	if replacement == ch {
		t.Fatalf("notFull channel was not replaced")
	}
}

func TestSignalAllLockedClosesAndReplacesBothChannelsWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	notEmpty := queue.notEmpty
	notFull := queue.notFull
	queue.notEmptyWaiters = 1
	queue.notFullWaiters = 1
	queue.signalAllLocked()
	newNotEmpty := queue.notEmpty
	newNotFull := queue.notFull
	queue.mu.Unlock()

	requireClosed(t, notEmpty)
	requireClosed(t, notFull)
	if newNotEmpty == notEmpty || newNotFull == notFull {
		t.Fatalf("signalAllLocked did not replace both channels")
	}
}
