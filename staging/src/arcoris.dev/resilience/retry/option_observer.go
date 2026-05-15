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

package retry

import "context"

// WithObserver registers an observer for retry events.
//
// Observers are appended in the order options are applied. Retry execution calls
// observers synchronously in registration order.
//
// Observers are notification boundaries only. They must not be used as retry
// policy, must not change retry decisions, and cannot report errors back to the
// retry loop. Observer failures must be handled inside the observer
// implementation.
//
// WithObserver panics when observer is nil.
func WithObserver(observer Observer) Option {
	requireObserver(observer)

	return func(cfg *config) {
		cfg.observers = append(cfg.observers, observer)
	}
}

// WithObserverFunc registers a function observer for retry events.
//
// WithObserverFunc is a convenience wrapper around ObserverFunc. It is intended
// for tests, small diagnostics hooks, and simple adapters. Larger observers
// should usually use named types so they can document concurrency, buffering,
// external-system failure handling, and ownership rules.
//
// WithObserverFunc panics when fn is nil.
func WithObserverFunc(fn func(context.Context, Event)) Option {
	requireObserverFunc(fn)

	return WithObserver(ObserverFunc(fn))
}
