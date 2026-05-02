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
	"sync"
)

const (
	// errNilNotifyContextParent is the stable diagnostic text used when
	// NotifyContext receives a nil parent context.
	errNilNotifyContextParent = "signals: nil notify context parent"
)

// StopFunc releases a signal context registration.
//
// StopFunc values returned by NotifyContext are idempotent. Calling StopFunc
// unregisters signal delivery and cancels the returned context if it has not
// already been cancelled. When the context already has a signal or parent cause,
// StopFunc does not overwrite that cause.
type StopFunc func()

// NotifyContext returns a child context cancelled when a configured signal is
// received.
//
// The returned context is cancelled when parent stops, when one configured
// signal is received, or when the returned StopFunc is called. Signal
// cancellation uses SignalError as the context cause. Parent cancellation
// preserves the parent cause when available.
//
// An empty sigs list means ShutdownSignals. NotifyContext panics when parent is
// nil or when any configured signal is nil.
func NotifyContext(parent context.Context, sigs ...os.Signal) (context.Context, StopFunc) {
	return notifyContextWithOptions(parent, sigs, nil)
}

// notifyContextWithOptions is the testable implementation behind NotifyContext.
//
// The options slice is package-local so tests can replace os/signal without
// adding a public dependency injection point to this low-level package.
func notifyContextWithOptions(parent context.Context, sigs []os.Signal, opts []SubscriptionOption) (context.Context, StopFunc) {
	requireContext(parent, errNilNotifyContextParent)

	ctx, cancel := context.WithCancelCause(parent)
	sub := SubscribeWithOptions(sigs, opts...)

	var once sync.Once
	stop := func() {
		once.Do(func() {
			sub.Stop()
			cancel(context.Canceled)
		})
	}

	go func() {
		sig, err := sub.Wait(parent)
		if err != nil {
			if errors.Is(err, ErrStopped) {
				return
			}

			cancel(err)
			sub.Stop()
			return
		}

		cancel(NewSignalError(sig))
		sub.Stop()
	}()

	return ctx, StopFunc(stop)
}
