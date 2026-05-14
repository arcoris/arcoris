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
	"context"
	"errors"
	"os"
	"testing"
)

func TestShutdownControllerDoesNotSubscribeEscalationSignalsBeforeShutdown(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGHUP),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	// Escalation-only signals must not be registered before shutdown starts,
	// because os/signal registration changes process-level signal behavior.
	assertSignalSlice(t, n.registeredSignals(), []os.Signal{testSIGTERM})
	if n.notifyCount() != 1 {
		t.Fatalf("notify count = %d, want initial shutdown registration only", n.notifyCount())
	}
}

func TestShutdownControllerRegistersEscalationSignalsAfterFirstShutdownSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGHUP),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
	n.waitNotifyCount(t, 2)

	assertSignalSlice(t, n.registeredSignals(), []os.Signal{testSIGTERM, testSIGHUP})
}

func TestShutdownControllerIgnoresPreShutdownEscalationSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGHUP),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	if n.emit(testSIGHUP) {
		t.Fatal("pre-shutdown escalation signal was delivered")
	}
	select {
	case <-controller.Done():
		t.Fatal("controller cancelled before a shutdown signal")
	default:
	}
	if event, ok := controller.First(); ok || event.Signal != nil {
		t.Fatalf("First() = (%v, %v), want empty false", event, ok)
	}
}

func TestShutdownControllerCancelsOnFirstShutdownSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGTERM),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())

	event, ok := controller.First()
	if !ok {
		t.Fatal("first signal was not recorded")
	}
	if !sameSignal(event.Signal, testSIGTERM) {
		t.Fatalf("first signal = %v, want %v", event.Signal, testSIGTERM)
	}
	if !errors.Is(context.Cause(controller.Context()), ErrSignal) {
		t.Fatalf("cause = %v, want ErrSignal", context.Cause(controller.Context()))
	}
}

func TestShutdownControllerDeliversRepeatedShutdownAsEscalationByDefault(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
	n.waitNotifyCount(t, 2)

	n.emit(testSIGTERM)
	mustReceiveSignal(t, controller.Escalation(), testSIGTERM)
}

func TestShutdownControllerDeliversConfiguredEscalationAfterFirstSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGHUP),
		WithEscalationBuffer(1),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
	n.waitNotifyCount(t, 2)

	n.emit(testSIGHUP)
	mustReceiveSignal(t, controller.Escalation(), testSIGHUP)
}

func TestShutdownControllerEscalationIsBestEffort(t *testing.T) {
	t.Parallel()

	controller := &ShutdownController{escalation: make(chan Event, 1)}
	controller.escalation <- Event{Signal: testSIGINT}

	// A full escalation channel models an owner that has not consumed the
	// previous advisory event. The controller must drop instead of blocking.
	done := make(chan struct{})
	go func() {
		controller.deliverEscalation(Event{Signal: testSIGTERM})
		close(done)
	}()

	mustClose(t, done)
	mustReceiveSignal(t, controller.Escalation(), testSIGINT)
}

func TestShutdownControllerCanDisableEscalation(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithNoEscalation(),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	if controller.Escalation() != nil {
		t.Fatal("disabled escalation returned non-nil channel")
	}

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
	n.waitStopCount(t, 1)

	if n.notifyCount() != 1 {
		t.Fatalf("notify count = %d, want no escalation registration", n.notifyCount())
	}
}

func TestShutdownControllerPreservesParentCause(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	parent, cancel := context.WithCancelCause(context.Background())
	controller := NewShutdownController(parent, withShutdownSubscriptionOptions(withNotifier(n)))
	defer controller.Stop()

	cancel(context.DeadlineExceeded)
	mustClose(t, controller.Done())
	n.waitStopCount(t, 1)

	if !errors.Is(context.Cause(controller.Context()), context.DeadlineExceeded) {
		t.Fatalf("cause = %v, want deadline exceeded", context.Cause(controller.Context()))
	}
}

func TestShutdownControllerStopDoesNotOverwriteSignalCause(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())

	// Stop is cleanup. Once a signal owns the cancellation cause, Stop must not
	// replace it with context.Canceled.
	controller.Stop()

	event, ok := Cause(controller.Context())
	if !ok {
		t.Fatal("signal cause was not preserved")
	}
	if !sameSignal(event.Signal, testSIGTERM) {
		t.Fatalf("signal = %v, want %v", event.Signal, testSIGTERM)
	}
}

func TestShutdownControllerStopDoesNotOverwriteParentCause(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	parent, cancel := context.WithCancelCause(context.Background())
	controller := NewShutdownController(parent, withShutdownSubscriptionOptions(withNotifier(n)))

	cancel(context.DeadlineExceeded)
	mustClose(t, controller.Done())

	// Parent cancellation is externally owned. Stop must release registration
	// without changing the already-observed parent cause.
	controller.Stop()

	if !errors.Is(context.Cause(controller.Context()), context.DeadlineExceeded) {
		t.Fatalf("cause = %v, want deadline exceeded", context.Cause(controller.Context()))
	}
}

func TestShutdownControllerEscalationChannelClosesOnStop(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(context.Background(), withShutdownSubscriptionOptions(withNotifier(n)))

	controller.Stop()

	mustClose(t, controller.Escalation())
}

func TestShutdownControllerEscalationChannelClosesOnParentCancel(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	parent, cancel := context.WithCancelCause(context.Background())
	controller := NewShutdownController(parent, withShutdownSubscriptionOptions(withNotifier(n)))
	defer controller.Stop()

	cancel(context.Canceled)

	mustClose(t, controller.Escalation())
}

func TestShutdownControllerStopIsIdempotent(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(context.Background(), withShutdownSubscriptionOptions(withNotifier(n)))

	controller.Stop()
	controller.Stop()
	mustClose(t, controller.Done())

	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
}

func TestShutdownControllerRejectsNilParent(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilShutdownParent, func() {
		NewShutdownController(nil)
	})
}

func TestShutdownControllerRejectsNilReceiver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*ShutdownController)
	}{
		{name: "Context", fn: func(c *ShutdownController) { c.Context() }},
		{name: "Done", fn: func(c *ShutdownController) { c.Done() }},
		{name: "Stop", fn: func(c *ShutdownController) { c.Stop() }},
		{name: "First", fn: func(c *ShutdownController) { c.First() }},
		{name: "Escalation", fn: func(c *ShutdownController) { c.Escalation() }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilShutdownController, func() {
				tc.fn(nil)
			})
		})
	}
}

func TestShutdownControllerFirstEventIsImmutable(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGINT),
		WithEscalationSignals(testSIGTERM),
		withShutdownSubscriptionOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGINT)
	mustClose(t, controller.Done())
	n.waitNotifyCount(t, 2)

	n.emit(testSIGTERM)
	mustReceiveSignal(t, controller.Escalation(), testSIGTERM)

	event, ok := controller.First()
	if !ok {
		t.Fatal("first event was not recorded")
	}
	if !sameSignal(event.Signal, testSIGINT) {
		t.Fatalf("first signal = %v, want %v", event.Signal, testSIGINT)
	}
}
