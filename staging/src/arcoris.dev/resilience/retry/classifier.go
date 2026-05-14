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

// Classifier decides whether an operation-owned error may be retried.
//
// A Classifier is consulted only after an operation attempt returns a non-nil
// error. It answers whether retry execution may schedule another attempt, subject
// to retry-owned limits such as maximum attempts, maximum elapsed time, context
// state, and delay sequence availability.
//
// Classifier implementations classify operation-owned errors. They must not
// execute operations, sleep, consume delay sequences, mutate retry state,
// observe timers, emit events, or create retry-owned error wrappers.
//
// A nil error is never retryable. Retry loops normally do not call Classifier for
// nil errors because nil means the operation succeeded. Classifier
// implementations should nevertheless treat nil as non-retryable so direct calls
// and defensive code remain safe.
//
// The retry package does not infer idempotency or replay safety. Callers remain
// responsible for choosing a classifier that is valid for their operation,
// transport, storage system, transaction model, and side-effect semantics.
//
// Protocol-specific and domain-specific classifiers do not belong in this
// package. HTTP status codes, gRPC status codes, database serialization errors,
// storage conflicts, and controller reconciliation conflicts should be modeled by
// adapter packages or by caller-owned Classifier implementations.
type Classifier interface {
	// Retryable reports whether err may be retried.
	//
	// Retryable must return false for nil. It should be deterministic for the
	// same error value unless the implementation explicitly documents external
	// state, such as a budget, feature gate, or dependency health signal.
	Retryable(err error) bool
}

// NeverRetry returns a classifier that rejects every error.
//
// NeverRetry is the conservative classifier and should be the default for generic
// retry execution. Retrying is only safe when the caller has explicitly decided
// that the operation can be repeated and has configured an appropriate
// classifier.
func NeverRetry() Classifier {
	return neverRetryClassifier{}
}

// RetryAll returns a classifier that accepts every non-nil error.
//
// RetryAll is useful in tests, tightly controlled internal loops, and operations
// whose retry safety is guaranteed by the caller. It still treats nil as
// non-retryable because nil means success, not a retryable failure.
func RetryAll() Classifier {
	return retryAllClassifier{}
}

// neverRetryClassifier rejects all errors.
//
// The type is intentionally private so callers depend on the Classifier contract
// and the NeverRetry constructor rather than on a concrete implementation type.
type neverRetryClassifier struct{}

// Retryable reports false for every error, including nil.
func (neverRetryClassifier) Retryable(error) bool {
	return false
}

// retryAllClassifier accepts every non-nil error.
//
// The type is intentionally private so callers depend on the Classifier contract
// and the RetryAll constructor rather than on a concrete implementation type.
type retryAllClassifier struct{}

// Retryable reports whether err is non-nil.
func (retryAllClassifier) Retryable(err error) bool {
	return err != nil
}
