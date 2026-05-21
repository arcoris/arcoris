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


package channelassert

import (
	"testing"
	"time"
)

const testTimeout = time.Second

func TestRequireReceive(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		ch := make(chan int, 1)
		ch <- 42
		if got := RequireReceive(t, ch, testTimeout); got != 42 {
			t.Fatalf("RequireReceive() = %d, want 42", got)
		}
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		ch := make(chan string, 1)
		ch <- "ok"
		if got := RequireReceive(t, ch, testTimeout); got != "ok" {
			t.Fatalf("RequireReceive() = %q, want %q", got, "ok")
		}
	})
}

func TestRequireNoReceive(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	RequireNoReceive(t, ch)
}

func TestRequireClosed(t *testing.T) {
	t.Parallel()

	ch := make(chan int)
	close(ch)
	RequireClosed(t, ch, testTimeout)
}

func TestRequireSignal(t *testing.T) {
	t.Parallel()

	t.Run("send", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{}, 1)
		ch <- struct{}{}
		RequireSignal(t, ch, testTimeout)
	})

	t.Run("close", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		close(ch)
		RequireSignal(t, ch, testTimeout)
	})
}

func TestRequireNoSignal(t *testing.T) {
	t.Parallel()

	ch := make(chan struct{}, 1)
	RequireNoSignal(t, ch)
}
