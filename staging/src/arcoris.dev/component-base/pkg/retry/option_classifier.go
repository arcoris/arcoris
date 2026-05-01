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

// WithClassifier configures retryability classification.
//
// The classifier is consulted only after an operation attempt returns a non-nil
// operation-owned error. If the classifier rejects the error, retry stops and
// returns the original operation error unchanged. If the classifier accepts the
// error, retry still must satisfy retry-owned limits, context state, and backoff
// sequence availability before scheduling another attempt.
//
// The retry package does not infer idempotency, replay safety, protocol
// semantics, storage semantics, or transaction safety. Callers must supply a
// classifier that is valid for the operation being retried.
//
// WithClassifier panics when classifier is nil.
func WithClassifier(classifier Classifier) Option {
	requireClassifier(classifier)

	return func(config *config) {
		config.classifier = classifier
	}
}

// WithRetryable configures retryability classification with a function.
//
// WithRetryable is a convenience wrapper around ClassifierFunc. It is intended
// for small caller-owned classification rules and tests. Larger
// protocol-specific or domain-specific classifiers should usually use named
// types so they can document configuration, ownership, and invariants.
//
// The function is never called for nil errors by ClassifierFunc. Nil means
// operation success and is not retryable.
//
// WithRetryable panics when fn is nil.
func WithRetryable(fn func(error) bool) Option {
	requireRetryableFunc(fn)

	return WithClassifier(ClassifierFunc(fn))
}
