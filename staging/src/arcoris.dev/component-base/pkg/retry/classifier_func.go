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

const panicNilClassifierFunc = "retry: nil classifier function"

// ClassifierFunc adapts a function to the Classifier interface.
//
// ClassifierFunc is intended for small caller-owned retryability decisions,
// tests, and adapters that do not need a named type. Larger domain-specific
// classifiers may prefer explicit types so they can document configuration,
// ownership, and invariants.
//
// A nil ClassifierFunc is a programming error. Calling Retryable on a nil
// ClassifierFunc panics with a stable diagnostic message.
//
// Nil errors are handled by the adapter and are never passed to the wrapped
// function. This preserves the package invariant that nil means success and is
// not retryable.
type ClassifierFunc func(err error) bool

// Retryable reports whether err may be retried.
//
// Retryable panics when f is nil. It returns false for nil errors without
// calling f. For non-nil errors, it delegates classification to f.
func (f ClassifierFunc) Retryable(err error) bool {
	if f == nil {
		panic(panicNilClassifierFunc)
	}
	if err == nil {
		return false
	}

	return f(err)
}
