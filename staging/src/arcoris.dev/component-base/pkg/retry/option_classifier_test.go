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

import (
	"errors"
	"testing"
)

func TestWithClassifier(t *testing.T) {
	classifier := RetryAll()

	config := configOf(WithClassifier(classifier))

	if !config.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("configured classifier did not retry non-nil error")
	}
}

func TestWithClassifierLastWins(t *testing.T) {
	config := configOf(
		WithClassifier(NeverRetry()),
		WithClassifier(RetryAll()),
	)

	if !config.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("last classifier option did not win")
	}
}

func TestWithClassifierPanicsOnNilClassifier(t *testing.T) {
	expectPanic(t, panicNilClassifier, func() {
		_ = WithClassifier(nil)
	})
}

func TestWithRetryable(t *testing.T) {
	errExpected := errors.New("expected")

	config := configOf(WithRetryable(func(err error) bool {
		return errors.Is(err, errExpected)
	}))

	if !config.classifier.Retryable(errExpected) {
		t.Fatalf("configured retryable func rejected expected error")
	}
	if config.classifier.Retryable(errors.New("other")) {
		t.Fatalf("configured retryable func accepted unexpected error")
	}
}

func TestWithRetryablePanicsOnNilFunction(t *testing.T) {
	expectPanic(t, panicNilClassifierFunc, func() {
		_ = WithRetryable(nil)
	})
}
