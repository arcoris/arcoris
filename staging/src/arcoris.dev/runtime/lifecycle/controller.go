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

package lifecycle

import (
	"sync"
	"time"
)

// Controller owns and serializes the lifecycle state of one component instance.
//
// Controller records lifecycle transitions, enforces the transition table, runs
// transition guards, publishes committed transition metadata, exposes consistent
// snapshots, and notifies observers after successful commits.
//
// Controller does not run component work. Component owners call BeginStart,
// MarkRunning, BeginStop, MarkStopped, and MarkFailed around their own startup,
// runtime, and shutdown code.
//
// The zero Controller is usable and starts in StateNew with no guards, no
// observers, and time.Now as the commit time source. A Controller must not be
// copied after first use.
type Controller struct {
	mu sync.Mutex

	state          State
	revision       uint64
	lastTransition Transition
	failureCause   error

	now       func() time.Time
	guards    []TransitionGuard
	observers []Observer

	changed chan struct{}
	done    chan struct{}
}

// NewController constructs a lifecycle Controller.
//
// The returned controller starts in StateNew with revision zero and no committed
// LastTransition. Options configure construction-time dependencies such as the
// time source, transition guards, and observers.
func NewController(options ...Option) *Controller {
	config := newControllerConfig(options...)

	if config.now == nil {
		config.now = time.Now
	}

	return &Controller{
		state:     StateNew,
		now:       config.now,
		guards:    append([]TransitionGuard(nil), config.guards...),
		observers: append([]Observer(nil), config.observers...),
		changed:   make(chan struct{}),
		done:      make(chan struct{}),
	}
}
