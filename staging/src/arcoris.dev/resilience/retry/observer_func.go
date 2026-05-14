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

const panicNilObserverFunc = "retry: nil observer function"

// ObserverFunc adapts a function to the Observer interface.
//
// ObserverFunc is intended for small caller-owned observers, tests, and adapters
// that do not need a named type. Larger observers may prefer explicit types so
// they can document configuration, concurrency rules, buffering behavior, and
// external-system failure handling.
//
// A nil ObserverFunc is a programming error. Calling ObserveRetry on a nil
// ObserverFunc panics with a stable diagnostic message.
//
// ObserverFunc does not validate events before delegating to the wrapped
// function. Event construction and validation belong to retry execution and
// tests. Callers that need defensive validation can implement it inside the
// wrapped function.
type ObserverFunc func(ctx context.Context, event Event)

// ObserveRetry observes one retry event by calling f.
//
// ObserveRetry panics when f is nil. Otherwise, it forwards ctx and event
// unchanged to the wrapped function.
func (f ObserverFunc) ObserveRetry(ctx context.Context, event Event) {
	if f == nil {
		panic(panicNilObserverFunc)
	}

	f(ctx, event)
}
