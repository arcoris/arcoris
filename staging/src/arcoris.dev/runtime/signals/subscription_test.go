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

func TestSubscribeWithOptionsRegistersNormalizedSignals(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT, testSIGTERM, testSIGINT}, withNotifier(n))
	defer sub.Stop()

	assertSignalSlice(t, n.registeredSignals(), []os.Signal{testSIGINT, testSIGTERM})
}

func TestSubscribeWithOptionsUsesShutdownSignalsWhenEmpty(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions(nil, withNotifier(n))
	defer sub.Stop()

	if len(n.registeredSignals()) == 0 {
		t.Fatal("empty subscription did not register default shutdown signals")
	}
}

func TestSubscribeWithOptionsRejectsNilSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}

	mustPanicWith(t, errNilSignalSetSignal, func() {
		SubscribeWithOptions([]os.Signal{nil}, withNotifier(n))
	})
}

func TestSubscriptionWaitReturnsSignal(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	if !n.emit(testSIGINT) {
		t.Fatal("registered signal was not delivered")
	}
	sig, err := sub.Wait(context.Background())
	if err != nil {
		t.Fatalf("Wait error = %v", err)
	}
	if !sameSignal(sig, testSIGINT) {
		t.Fatalf("signal = %v, want %v", sig, testSIGINT)
	}
}

func TestSubscriptionWaitReturnsContextCause(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(context.DeadlineExceeded)

	_, err := sub.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait error = %v, want context deadline", err)
	}
}

func TestSubscriptionWaitReturnsErrStopped(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	sub.Stop()

	_, err := sub.Wait(context.Background())
	if !errors.Is(err, ErrStopped) {
		t.Fatalf("Wait error = %v, want ErrStopped", err)
	}
}

func TestSubscriptionStopIsIdempotentAndClosesDone(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))

	sub.Stop()
	sub.Stop()

	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
	if !n.stopped() {
		t.Fatal("notifier was not marked stopped")
	}
	mustClose(t, sub.Done())
}

func TestSubscriptionChannelAccessor(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	if sub.C() == nil {
		t.Fatal("C returned nil")
	}
}

func TestSubscriptionDirectChannelReceiveCompetesWithWait(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	// C and Wait share one channel. The direct receive consumes the only signal,
	// so the later Wait observes the already-cancelled context instead.
	n.emit(testSIGINT)
	mustReceiveOSSignal(t, sub.C(), testSIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := sub.Wait(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait error = %v, want context.Canceled", err)
	}
}

func TestSubscriptionRegisterMoreExtendsRegistration(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	if !sub.registerMore([]os.Signal{testSIGTERM, testSIGINT}) {
		t.Fatal("registerMore returned false for active subscription")
	}

	assertSignalSlice(t, n.registeredSignals(), []os.Signal{testSIGINT, testSIGTERM})
	if n.notifyCount() != 2 {
		t.Fatalf("notify count = %d, want 2", n.notifyCount())
	}
}

func TestSubscriptionRegisterMoreDoesNotRegisterAfterStop(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	sub.Stop()

	if sub.registerMore([]os.Signal{testSIGTERM}) {
		t.Fatal("registerMore returned true after Stop")
	}
	if n.notifyCount() != 1 {
		t.Fatalf("notify count = %d, want only initial registration", n.notifyCount())
	}
}

func TestSubscriptionRejectsNilReceiver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*Subscription)
	}{
		{name: "C", fn: func(s *Subscription) { s.C() }},
		{name: "Wait", fn: func(s *Subscription) { _, _ = s.Wait(context.Background()) }},
		{name: "Stop", fn: func(s *Subscription) { s.Stop() }},
		{name: "Done", fn: func(s *Subscription) { s.Done() }},
		{name: "registerMore", fn: func(s *Subscription) { s.registerMore([]os.Signal{testSIGINT}) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilSubscription, func() {
				tc.fn(nil)
			})
		})
	}
}

func TestSubscriptionWaitRejectsNilContext(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	mustPanicWith(t, errNilSubscriptionContext, func() {
		sub.Wait(nil)
	})
}

func TestContextCauseFallsBackToContextErr(t *testing.T) {
	t.Parallel()

	if err := contextCause(context.Background()); err != nil {
		t.Fatalf("contextCause(background) = %v, want nil", err)
	}
}
