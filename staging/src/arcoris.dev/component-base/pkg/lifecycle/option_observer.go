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

// WithObserver appends one lifecycle observer to the Controller configuration.
//
// Observers are notified after a transition has been committed. They observe
// facts that already happened and cannot reject, modify, or roll back a
// transition.
//
// Observers are notified in configuration order. If several WithObserver and
// WithObservers options are provided, their observers are appended in the same
// order in which the options are applied.
//
// A nil observer is ignored. This keeps conditional option construction safe.
func WithObserver(observer Observer) Option {
	return func(config *controllerConfig) {
		if observer == nil {
			return
		}

		config.observers = append(config.observers, observer)
	}
}

// WithObservers appends several lifecycle observers to the Controller
// configuration.
//
// Observers are notified in the order provided. The same committed Transition
// value is passed to every observer.
//
// Nil observers are ignored. This makes it safe to build observer lists from
// optional diagnostics, tracing, metrics, or test integrations without adding
// special-case filtering at call sites.
func WithObservers(observers ...Observer) Option {
	return func(config *controllerConfig) {
		for _, observer := range observers {
			if observer == nil {
				continue
			}

			config.observers = append(config.observers, observer)
		}
	}
}
