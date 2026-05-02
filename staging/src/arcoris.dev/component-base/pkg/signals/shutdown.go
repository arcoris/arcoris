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
	// errNilShutdownParent is the stable diagnostic text used when
	// NewShutdownController receives a nil parent context.
	errNilShutdownParent = "signals: nil shutdown parent context"

	// errNilShutdownController is the stable diagnostic text used when a
	// ShutdownController method is called on a nil receiver.
	errNilShutdownController = "signals: nil shutdown controller"
)

// ShutdownController coordinates graceful shutdown from process signals.
//
// The first configured shutdown signal cancels Context with a SignalError cause
// and records First. Repeated escalation signals are delivered through
// Escalation when escalation is enabled. ShutdownController never exits the
// process and never drives component lifecycle transitions directly.
//
// ShutdownController values must be created with NewShutdownController and must
// not be copied after first use.
type ShutdownController struct {
	noCopy noCopy

	ctx    context.Context
	cancel context.CancelCauseFunc

	sub *Subscription

	shutdownSignals   []os.Signal
	escalationSignals []os.Signal
	escalationEnabled bool
	escalation        chan Event

	firstMu  sync.Mutex
	first    Event
	hasFirst bool

	stopOnce sync.Once
}

// NewShutdownController constructs and starts a shutdown controller.
//
// The controller listens for shutdown signals immediately. The owner must call
// Stop when the controller is no longer needed. NewShutdownController panics
// when parent is nil or when options produce invalid signal configuration.
func NewShutdownController(parent context.Context, opts ...ShutdownOption) *ShutdownController {
	requireContext(parent, errNilShutdownParent)

	config := newShutdownConfig(opts...)
	subscribeSet := config.shutdownSignals
	if config.escalationEnabled {
		subscribeSet = Merge(subscribeSet, config.escalationSignals)
	}

	ctx, cancel := context.WithCancelCause(parent)
	var escalation chan Event
	if config.escalationEnabled {
		escalation = make(chan Event, config.escalationBuffer)
	}

	controller := &ShutdownController{
		ctx:               ctx,
		cancel:            cancel,
		sub:               SubscribeWithOptions(subscribeSet, config.subscribeOptions...),
		shutdownSignals:   Clone(config.shutdownSignals),
		escalationSignals: Clone(config.escalationSignals),
		escalationEnabled: config.escalationEnabled,
		escalation:        escalation,
	}

	go controller.run(parent)

	return controller
}

// Context returns the shutdown context owned by the controller.
//
// The context is cancelled by the first shutdown signal, parent cancellation, or
// Stop.
func (c *ShutdownController) Context() context.Context {
	requireShutdownController(c)
	return c.ctx
}

// Done returns the shutdown context Done channel.
func (c *ShutdownController) Done() <-chan struct{} {
	requireShutdownController(c)
	return c.ctx.Done()
}

// Stop releases signal registration and cancels the controller context.
//
// Stop is idempotent.
func (c *ShutdownController) Stop() {
	requireShutdownController(c)

	c.stopOnce.Do(func() {
		c.sub.Stop()
		c.cancel(context.Canceled)
	})
}

// First returns the first shutdown signal observed by the controller.
func (c *ShutdownController) First() (Event, bool) {
	requireShutdownController(c)

	c.firstMu.Lock()
	defer c.firstMu.Unlock()

	return c.first, c.hasFirst
}

// Escalation returns the repeated-signal escalation channel.
//
// When escalation is disabled, Escalation returns nil. When escalation is
// enabled, the channel is closed after the controller's signal loop exits.
// Escalation delivery is best-effort and non-blocking after the configured
// buffer is full.
func (c *ShutdownController) Escalation() <-chan Event {
	requireShutdownController(c)
	return c.escalation
}

// run owns the controller signal loop.
func (c *ShutdownController) run(parent context.Context) {
	defer c.sub.Stop()
	if c.escalation != nil {
		defer close(c.escalation)
	}

	for {
		sig, err := c.sub.Wait(parent)
		if err != nil {
			if errors.Is(err, ErrStopped) {
				return
			}

			c.cancel(err)
			return
		}

		event := Event{Signal: sig}
		if !c.firstRecorded() && Contains(c.shutdownSignals, sig) {
			c.recordFirst(event)
			c.cancel(NewSignalError(sig))

			if !c.escalationEnabled {
				return
			}
			continue
		}

		if c.firstRecorded() && c.escalationEnabled && Contains(c.escalationSignals, sig) {
			select {
			case c.escalation <- event:
			default:
			}
		}
	}
}

// firstRecorded reports whether the first shutdown signal has been recorded.
func (c *ShutdownController) firstRecorded() bool {
	c.firstMu.Lock()
	defer c.firstMu.Unlock()

	return c.hasFirst
}

// recordFirst records event as the first shutdown signal if none has been
// recorded yet.
func (c *ShutdownController) recordFirst(event Event) {
	c.firstMu.Lock()
	defer c.firstMu.Unlock()

	if c.hasFirst {
		return
	}

	c.first = event
	c.hasFirst = true
}

// requireShutdownController panics when c is nil.
func requireShutdownController(c *ShutdownController) {
	if c == nil {
		panic(errNilShutdownController)
	}
}
