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


package run

import (
	"runtime"
	"testing"
	"time"
)

const testTimeout = time.Second

func waitGroupTaskErrorCount(t *testing.T, group *Group, want int) {
	t.Helper()

	deadline := time.NewTimer(testTimeout)
	defer deadline.Stop()

	for {
		group.mu.Lock()
		got := len(group.errs)
		group.mu.Unlock()
		if got == want {
			return
		}

		select {
		case <-deadline.C:
			t.Fatalf("timed out waiting for %d task errors", want)
		default:
			runtime.Gosched()
		}
	}
}

func waitGroupClosed(t *testing.T, group *Group) {
	t.Helper()

	deadline := time.NewTimer(testTimeout)
	defer deadline.Stop()

	for {
		group.mu.Lock()
		closed := group.closed
		group.mu.Unlock()
		if closed {
			return
		}

		select {
		case <-deadline.C:
			t.Fatal("timed out waiting for group close")
		default:
			runtime.Gosched()
		}
	}
}
