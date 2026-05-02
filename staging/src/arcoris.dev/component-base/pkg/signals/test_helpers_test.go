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

const testTimeout = time.Second

type testSignal string

func (s testSignal) Signal() {}

func (s testSignal) String() string { return string(s) }

const (
	testSIGINT  testSignal = "SIGINT"
	testSIGTERM testSignal = "SIGTERM"
	testSIGHUP  testSignal = "SIGHUP"
	testSIGQUIT testSignal = "SIGQUIT"
)

// fakeNotifier models the os/signal behaviors that package signals relies on.
//
// It tracks registration ownership per channel, treats repeated notify calls as
// registration extension, delivers only registered signals, and drops delivery
// after stop. The fake is intentionally stricter than a simple channel sender so
// tests cannot pass when production code forgets to register a signal.
type fakeNotifier struct {
	mu sync.Mutex

	changes chan struct{}

	channels        []chan<- os.Signal
	registered      map[chan<- os.Signal]map[string]os.Signal
	registeredOrder map[chan<- os.Signal][]os.Signal
	stoppedChannels map[chan<- os.Signal]bool

	notifyCalls [][]os.Signal
	stopCalls   []chan<- os.Signal
}

func (n *fakeNotifier) notify(ch chan<- os.Signal, sigs ...os.Signal) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.initLocked()
	n.notifyCalls = append(n.notifyCalls, Clone(sigs))

	if _, ok := n.registered[ch]; !ok {
		n.channels = append(n.channels, ch)
		n.registered[ch] = make(map[string]os.Signal, len(sigs))
	}
	n.stoppedChannels[ch] = false

	for _, sig := range sigs {
		key := signalKey(sig)
		if _, ok := n.registered[ch][key]; !ok {
			n.registeredOrder[ch] = append(n.registeredOrder[ch], sig)
		}
		n.registered[ch][key] = sig
	}

	n.broadcastLocked()
}

func (n *fakeNotifier) stop(ch chan<- os.Signal) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.initLocked()
	n.stopCalls = append(n.stopCalls, ch)
	n.stoppedChannels[ch] = true
	delete(n.registered, ch)
	delete(n.registeredOrder, ch)

	n.broadcastLocked()
}

// emit sends sig to every active channel registered for sig.
//
// Delivery is non-blocking to match os/signal's expectation that callers provide
// enough buffer for the notifications they care about. The return value reports
// whether at least one channel accepted the signal.
func (n *fakeNotifier) emit(sig os.Signal) bool {
	n.mu.Lock()
	n.initLocked()

	key := signalKey(sig)
	targets := make([]chan<- os.Signal, 0, len(n.channels))
	for _, ch := range n.channels {
		if n.stoppedChannels[ch] {
			continue
		}
		registered, ok := n.registered[ch]
		if !ok {
			continue
		}
		if _, ok := registered[key]; ok {
			targets = append(targets, ch)
		}
	}
	n.mu.Unlock()

	delivered := false
	for _, ch := range targets {
		select {
		case ch <- sig:
			delivered = true
		default:
		}
	}
	return delivered
}

func (n *fakeNotifier) registeredSignals() []os.Signal {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.initLocked()
	seen := make(map[string]struct{})
	var sigs []os.Signal
	for _, ch := range n.channels {
		if n.stoppedChannels[ch] {
			continue
		}
		for _, sig := range n.registeredOrder[ch] {
			key := signalKey(sig)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			sigs = append(sigs, sig)
		}
	}
	return sigs
}

func (n *fakeNotifier) notifyCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()

	return len(n.notifyCalls)
}

func (n *fakeNotifier) stopCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()

	return len(n.stopCalls)
}

func (n *fakeNotifier) stopped() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if len(n.channels) == 0 {
		return false
	}
	return n.stoppedChannels[n.channels[len(n.channels)-1]]
}

func (n *fakeNotifier) waitNotifyCount(t *testing.T, want int) {
	t.Helper()

	n.wait(t, fmt.Sprintf("notify count >= %d", want), func() bool {
		return len(n.notifyCalls) >= want
	})
}

func (n *fakeNotifier) waitStopCount(t *testing.T, want int) {
	t.Helper()

	n.wait(t, fmt.Sprintf("stop count >= %d", want), func() bool {
		return len(n.stopCalls) >= want
	})
}

func (n *fakeNotifier) wait(t *testing.T, description string, ready func() bool) {
	t.Helper()

	timer := time.NewTimer(testTimeout)
	defer timer.Stop()

	for {
		n.mu.Lock()
		n.initLocked()
		if ready() {
			n.mu.Unlock()
			return
		}
		changes := n.changes
		n.mu.Unlock()

		select {
		case <-changes:
		case <-timer.C:
			t.Fatalf("timed out waiting for %s", description)
		}
	}
}

func (n *fakeNotifier) initLocked() {
	if n.changes == nil {
		n.changes = make(chan struct{})
	}
	if n.registered == nil {
		n.registered = make(map[chan<- os.Signal]map[string]os.Signal)
	}
	if n.registeredOrder == nil {
		n.registeredOrder = make(map[chan<- os.Signal][]os.Signal)
	}
	if n.stoppedChannels == nil {
		n.stoppedChannels = make(map[chan<- os.Signal]bool)
	}
}

func (n *fakeNotifier) broadcastLocked() {
	close(n.changes)
	n.changes = make(chan struct{})
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
	case got, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before signal event")
		}
		if !sameSignal(got.Signal, want) {
			t.Fatalf("signal = %v, want %v", got.Signal, want)
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for signal event")
	}
}

func mustReceiveOSSignal(t *testing.T, ch <-chan os.Signal, want os.Signal) {
	t.Helper()

	select {
	case got := <-ch:
		if !sameSignal(got, want) {
			t.Fatalf("signal = %v, want %v", got, want)
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for signal")
	}
}

func mustClose[T any](t *testing.T, ch <-chan T) {
	t.Helper()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("received value before channel close")
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for channel close")
	}
}

func sameSignal(left, right os.Signal) bool {
	if left == nil || right == nil {
		return left == right
	}
	return signalKey(left) == signalKey(right)
}

func TestFakeNotifierDeliversOnlyRegisteredSignals(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	ch := make(chan os.Signal, 1)
	n.notify(ch, testSIGINT)

	if n.emit(testSIGTERM) {
		t.Fatal("fake notifier delivered an unregistered signal")
	}
	select {
	case got := <-ch:
		t.Fatalf("received unregistered signal %v", got)
	default:
	}

	if !n.emit(testSIGINT) {
		t.Fatal("fake notifier did not deliver a registered signal")
	}
	mustReceiveOSSignal(t, ch, testSIGINT)
}

func TestFakeNotifierRepeatedNotifyMergesSignals(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	ch := make(chan os.Signal, 2)

	n.notify(ch, testSIGINT)
	n.notify(ch, testSIGTERM, testSIGINT)

	got := n.registeredSignals()
	want := []os.Signal{testSIGINT, testSIGTERM}
	if len(got) != len(want) {
		t.Fatalf("registered len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !sameSignal(got[i], want[i]) {
			t.Fatalf("registered[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if n.notifyCount() != 2 {
		t.Fatalf("notify count = %d, want 2", n.notifyCount())
	}

	n.emit(testSIGINT)
	n.emit(testSIGTERM)
	mustReceiveOSSignal(t, ch, testSIGINT)
	mustReceiveOSSignal(t, ch, testSIGTERM)
}

func TestFakeNotifierStopPreventsDelivery(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	ch := make(chan os.Signal, 1)

	n.notify(ch, testSIGINT)
	n.stop(ch)

	if !n.stopped() {
		t.Fatal("fake notifier did not record stopped state")
	}
	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
	if len(n.registeredSignals()) != 0 {
		t.Fatalf("registered signals after stop = %v, want none", n.registeredSignals())
	}
	if n.emit(testSIGINT) {
		t.Fatal("fake notifier delivered after stop")
	}
	select {
	case got := <-ch:
		t.Fatalf("received stopped signal %v", got)
	default:
	}
}
