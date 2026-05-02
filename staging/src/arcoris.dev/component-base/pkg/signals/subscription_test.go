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
	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT, testSIGTERM, testSIGINT}, withNotifier(n))
	defer sub.Stop()

	got := n.notifiedSignals()
	want := []os.Signal{testSIGINT, testSIGTERM}
	if len(got) != len(want) {
		t.Fatalf("registered len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !sameSignal(got[i], want[i]) {
			t.Fatalf("registered[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestSubscribeWithOptionsUsesShutdownWhenSignalsEmpty(t *testing.T) {
	n := &fakeNotifier{}
	sub := SubscribeWithOptions(nil, withNotifier(n))
	defer sub.Stop()

	if len(n.notifiedSignals()) == 0 {
		t.Fatal("empty subscription did not register default shutdown signals")
	}
}

func TestSubscribeRejectsNilSignal(t *testing.T) {
	n := &fakeNotifier{}
	mustPanicWith(t, errNilSignalSetSignal, func() {
		SubscribeWithOptions([]os.Signal{nil}, withNotifier(n))
	})
}

func TestSubscriptionWaitReturnsSignal(t *testing.T) {
	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	n.emit(testSIGINT)
	sig, err := sub.Wait(context.Background())
	if err != nil {
		t.Fatalf("Wait error = %v", err)
	}
	if !sameSignal(sig, testSIGINT) {
		t.Fatalf("signal = %v, want %v", sig, testSIGINT)
	}
}

func TestSubscriptionWaitReturnsContextCause(t *testing.T) {
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
	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	sub.Stop()

	_, err := sub.Wait(context.Background())
	if !errors.Is(err, ErrStopped) {
		t.Fatalf("Wait error = %v, want ErrStopped", err)
	}
}

func TestSubscriptionStopIsIdempotent(t *testing.T) {
	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))

	sub.Stop()
	sub.Stop()

	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
	mustClose(t, sub.Done())
}

func TestSubscriptionRejectsNilReceiverAndNilContext(t *testing.T) {
	mustPanicWith(t, errNilSubscription, func() {
		var sub *Subscription
		sub.Done()
	})

	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	mustPanicWith(t, errNilSubscriptionContext, func() {
		sub.Wait(nil)
	})
}

func TestSubscriptionChannelAccessors(t *testing.T) {
	n := &fakeNotifier{}
	sub := SubscribeWithOptions([]os.Signal{testSIGINT}, withNotifier(n))
	defer sub.Stop()

	if sub.C() == nil {
		t.Fatal("C returned nil")
	}
}

func TestContextCauseFallsBackToContextErr(t *testing.T) {
	if err := contextCause(context.Background()); err != nil {
		t.Fatalf("contextCause(background) = %v, want nil", err)
	}
}
