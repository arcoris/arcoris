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

package signals

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

type testSignal string

func (s testSignal) Signal() {}

func (s testSignal) String() string { return string(s) }

const (
	testSIGINT  testSignal = "SIGINT"
	testSIGTERM testSignal = "SIGTERM"
	testSIGHUP  testSignal = "SIGHUP"
)

type fakeNotifier struct {
	mu       sync.Mutex
	ch       chan<- os.Signal
	notifies [][]os.Signal
	stops    int
}

func (n *fakeNotifier) notify(ch chan<- os.Signal, sigs ...os.Signal) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ch = ch
	n.notifies = append(n.notifies, Clone(sigs))
}

func (n *fakeNotifier) stop(ch chan<- os.Signal) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.stops++
}

func (n *fakeNotifier) emit(sig os.Signal) {
	n.mu.Lock()
	ch := n.ch
	n.mu.Unlock()

	if ch == nil {
		panic("signals test: nil fake notifier channel")
	}

	ch <- sig
}

func (n *fakeNotifier) stopCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.stops
}

func (n *fakeNotifier) notifiedSignals() []os.Signal {
	n.mu.Lock()
	defer n.mu.Unlock()

	if len(n.notifies) == 0 {
		return nil
	}
	return Clone(n.notifies[len(n.notifies)-1])
}

func mustPanicWith(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected panic %q", want)
		}

		got := fmt.Sprint(recovered)
		if got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	fn()
}

func mustReceiveSignal(t *testing.T, ch <-chan Event, want os.Signal) {
	t.Helper()

	select {
	case got := <-ch:
		if !sameSignal(got.Signal, want) {
			t.Fatalf("signal = %v, want %v", got.Signal, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for signal event")
	}
}

func mustClose(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}

func sameSignal(left, right os.Signal) bool {
	if left == nil || right == nil {
		return left == right
	}
	return signalKey(left) == signalKey(right)
}
