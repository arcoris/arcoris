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
	"os"
	"sync"
)

const (
	// errNilSubscriptionContext is the stable diagnostic text used when
	// Subscription.Wait receives a nil context.
	errNilSubscriptionContext = "signals: nil subscription wait context"

	// errNilSubscription is the stable diagnostic text used when a Subscription
	// method is called on a nil receiver.
	errNilSubscription = "signals: nil subscription"
)

// Subscription owns one process signal registration.
//
// A Subscription registers an internal channel with package os/signal and
// provides owner-controlled waiting and cleanup. It is the lowest-level signal
// owner in this package: it does not cancel contexts, interpret signals as
// lifecycle transitions, exit the process, or apply shutdown policy.
//
// Subscription values must be created with Subscribe or SubscribeWithOptions.
// The owner must call Stop when delivery is no longer needed, and must not copy
// the value after construction.
type Subscription struct {
	noCopy noCopy

	ch       chan os.Signal
	done     chan struct{}
	notifier notifier
	once     sync.Once
	mu       sync.Mutex
	stopped  bool
}

// Subscribe creates a Subscription for sigs.
//
// An empty sigs list means ShutdownSignals. Subscribe panics when any signal is
// nil.
func Subscribe(sigs ...os.Signal) *Subscription {
	return SubscribeWithOptions(sigs)
}

// SubscribeWithOptions creates a Subscription for sigs with options.
//
// An empty sigs list means ShutdownSignals. Duplicate signals are removed while
// preserving first occurrence order. The returned Subscription owns its signal
// registration until Stop is called.
func SubscribeWithOptions(sigs []os.Signal, opts ...SubscriptionOption) *Subscription {
	config := newSubscribeConfig(opts...)
	registered := signalSetOrShutdownSignals(sigs)

	s := &Subscription{
		ch:       make(chan os.Signal, config.buffer),
		done:     make(chan struct{}),
		notifier: config.notifier,
	}
	s.notifier.notify(s.ch, registered...)

	return s
}

// C returns the signal delivery channel owned by the subscription.
//
// The returned channel is receive-only for callers. Callers must not close it
// and must still call Stop when the subscription is no longer needed.
//
// C and Wait receive from the same underlying channel. A direct receive from C
// consumes a signal that a concurrent or later Wait call could otherwise have
// returned, so owners should choose one receive path for each subscription
// unless they deliberately coordinate that competition.
func (s *Subscription) C() <-chan os.Signal {
	requireSubscription(s)
	return s.ch
}

// Wait blocks until a signal is received, ctx stops, or the subscription stops.
//
// Wait returns the received signal on success. If ctx stops first, Wait returns
// the context cancellation cause when available, otherwise ctx.Err. If Stop is
// called first, Wait returns ErrStopped. Wait panics when ctx is nil.
func (s *Subscription) Wait(ctx context.Context) (os.Signal, error) {
	requireSubscription(s)
	requireContext(ctx, errNilSubscriptionContext)

	select {
	case sig := <-s.ch:
		return sig, nil

	case <-ctx.Done():
		return nil, contextCause(ctx)

	case <-s.done:
		return nil, ErrStopped
	}
}

// Stop unregisters the subscription from process signal delivery.
//
// Stop is idempotent. It does not close C because package os/signal also does
// not close delivery channels. It closes Done so goroutines waiting through Wait
// can leave with ErrStopped.
func (s *Subscription) Stop() {
	requireSubscription(s)

	s.once.Do(func() {
		s.mu.Lock()
		s.stopped = true
		s.notifier.stop(s.ch)
		s.mu.Unlock()

		close(s.done)
	})
}

// Done returns a channel closed when Stop completes.
func (s *Subscription) Done() <-chan struct{} {
	requireSubscription(s)
	return s.done
}

// registerMore extends the subscription's os/signal registration.
//
// Package os/signal permits repeated Notify calls for the same channel; later
// calls add signals to the channel's registration instead of replacing the
// previous set. ShutdownController relies on that behavior to delay escalation
// registration until graceful shutdown has actually started.
//
// The method returns false when Stop has already started. In that case it leaves
// the stopped subscription alone instead of resurrecting process signal
// delivery during owner cleanup.
func (s *Subscription) registerMore(sigs []os.Signal) bool {
	requireSubscription(s)

	registered := Unique(sigs)
	if len(registered) == 0 {
		return true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return false
	}

	s.notifier.notify(s.ch, registered...)
	return true
}

// requireSubscription panics when s is nil.
func requireSubscription(s *Subscription) {
	if s == nil {
		panic(errNilSubscription)
	}
}

// contextCause returns ctx's cancellation cause with a fallback to ctx.Err.
func contextCause(ctx context.Context) error {
	if cause := context.Cause(ctx); cause != nil {
		return cause
	}
	return ctx.Err()
}
