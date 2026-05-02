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
	"testing"
	"time"
)

func TestShutdownControllerCancelsOnFirstShutdownSignal(t *testing.T) {
	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGTERM),
		withShutdownSubscribeOptions(withNotifier(n)),
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

func TestShutdownControllerDeliversEscalationAfterFirstSignal(t *testing.T) {
	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGTERM),
		WithEscalationBuffer(1),
		withShutdownSubscribeOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
	n.emit(testSIGTERM)

	mustReceiveSignal(t, controller.Escalation(), testSIGTERM)
}

func TestShutdownControllerEscalationIsBestEffort(t *testing.T) {
	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithEscalationSignals(testSIGTERM),
		WithEscalationBuffer(0),
		withShutdownSubscribeOptions(withNotifier(n)),
	)
	defer controller.Stop()

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())

	done := make(chan struct{})
	go func() {
		n.emit(testSIGTERM)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("best-effort escalation send blocked")
	}
}

func TestShutdownControllerCanDisableEscalation(t *testing.T) {
	n := &fakeNotifier{}
	controller := NewShutdownController(
		context.Background(),
		WithShutdownSignals(testSIGTERM),
		WithNoEscalation(),
		withShutdownSubscribeOptions(withNotifier(n)),
	)
	defer controller.Stop()

	if controller.Escalation() != nil {
		t.Fatal("disabled escalation returned non-nil channel")
	}

	n.emit(testSIGTERM)
	mustClose(t, controller.Done())
}

func TestShutdownControllerPreservesParentCause(t *testing.T) {
	n := &fakeNotifier{}
	parent, cancel := context.WithCancelCause(context.Background())
	controller := NewShutdownController(parent, withShutdownSubscribeOptions(withNotifier(n)))
	defer controller.Stop()

	cancel(context.DeadlineExceeded)
	mustClose(t, controller.Done())

	if !errors.Is(context.Cause(controller.Context()), context.DeadlineExceeded) {
		t.Fatalf("cause = %v, want deadline exceeded", context.Cause(controller.Context()))
	}
}

func TestShutdownControllerStopIsIdempotent(t *testing.T) {
	n := &fakeNotifier{}
	controller := NewShutdownController(context.Background(), withShutdownSubscribeOptions(withNotifier(n)))

	controller.Stop()
	controller.Stop()
	mustClose(t, controller.Done())

	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
}

func TestShutdownControllerRejectsNilParentAndNilReceiver(t *testing.T) {
	mustPanicWith(t, errNilShutdownParent, func() {
		NewShutdownController(nil)
	})
	mustPanicWith(t, errNilShutdownController, func() {
		var controller *ShutdownController
		controller.Done()
	})
}

func TestShutdownControllerRecordFirstKeepsOriginalEvent(t *testing.T) {
	controller := &ShutdownController{}
	controller.recordFirst(Event{Signal: testSIGINT})
	controller.recordFirst(Event{Signal: testSIGTERM})

	event, ok := controller.First()
	if !ok {
		t.Fatal("first event was not recorded")
	}
	if !sameSignal(event.Signal, testSIGINT) {
		t.Fatalf("first signal = %v, want %v", event.Signal, testSIGINT)
	}
}
