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
	"time"
)

func TestNotifyContextCancelsWithSignalCause(t *testing.T) {
	n := &fakeNotifier{}
	ctx, stop := notifyContextWithOptions(context.Background(), []os.Signal{testSIGTERM}, []SubscribeOption{withNotifier(n)})
	defer stop()

	n.emit(testSIGTERM)
	mustClose(t, ctx.Done())

	event, ok := Cause(ctx)
	if !ok {
		t.Fatal("signal cause was not recorded")
	}
	if !sameSignal(event.Signal, testSIGTERM) {
		t.Fatalf("signal = %v, want %v", event.Signal, testSIGTERM)
	}
}

func TestNotifyContextPreservesParentCause(t *testing.T) {
	n := &fakeNotifier{}
	parent, cancel := context.WithCancelCause(context.Background())
	ctx, stop := notifyContextWithOptions(parent, []os.Signal{testSIGTERM}, []SubscribeOption{withNotifier(n)})
	defer stop()

	cancel(context.DeadlineExceeded)
	mustClose(t, ctx.Done())

	if !errors.Is(context.Cause(ctx), context.DeadlineExceeded) {
		t.Fatalf("cause = %v, want deadline exceeded", context.Cause(ctx))
	}
}

func TestNotifyContextStopCancelsAndUnregisters(t *testing.T) {
	n := &fakeNotifier{}
	ctx, stop := notifyContextWithOptions(context.Background(), []os.Signal{testSIGTERM}, []SubscribeOption{withNotifier(n)})

	stop()
	stop()
	mustClose(t, ctx.Done())

	if !errors.Is(context.Cause(ctx), context.Canceled) {
		t.Fatalf("cause = %v, want context.Canceled", context.Cause(ctx))
	}
	if n.stopCount() != 1 {
		t.Fatalf("stop count = %d, want 1", n.stopCount())
	}
}

func TestNotifyContextUsesShutdownWhenSignalsEmpty(t *testing.T) {
	n := &fakeNotifier{}
	_, stop := notifyContextWithOptions(context.Background(), nil, []SubscribeOption{withNotifier(n)})
	defer stop()

	if len(n.notifiedSignals()) == 0 {
		t.Fatal("NotifyContext did not register shutdown signals")
	}
}

func TestNotifyContextRejectsNilParent(t *testing.T) {
	mustPanicWith(t, errNilNotifyContextParent, func() {
		NotifyContext(nil)
	})
}

func TestNotifyContextSignalStopsSubscription(t *testing.T) {
	n := &fakeNotifier{}
	ctx, stop := notifyContextWithOptions(context.Background(), []os.Signal{testSIGINT}, []SubscribeOption{withNotifier(n)})
	defer stop()

	n.emit(testSIGINT)
	mustClose(t, ctx.Done())

	deadline := time.After(time.Second)
	for n.stopCount() == 0 {
		select {
		case <-deadline:
			t.Fatal("subscription was not stopped after signal")
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
