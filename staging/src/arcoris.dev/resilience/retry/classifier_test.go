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

func TestNeverRetryRejectsNilAndNonNilErrors(t *testing.T) {
	classifier := NeverRetry()
	if classifier == nil {
		t.Fatalf("NeverRetry() returned nil")
	}

	if classifier.Retryable(nil) {
		t.Fatalf("NeverRetry().Retryable(nil) = true, want false")
	}

	errBoom := errors.New("boom")
	if classifier.Retryable(errBoom) {
		t.Fatalf("NeverRetry().Retryable(non-nil error) = true, want false")
	}
}

func TestRetryAllAcceptsOnlyNonNilErrors(t *testing.T) {
	classifier := RetryAll()
	if classifier == nil {
		t.Fatalf("RetryAll() returned nil")
	}

	if classifier.Retryable(nil) {
		t.Fatalf("RetryAll().Retryable(nil) = true, want false")
	}

	errBoom := errors.New("boom")
	if !classifier.Retryable(errBoom) {
		t.Fatalf("RetryAll().Retryable(non-nil error) = false, want true")
	}
}

func TestBuiltInClassifiersAreReusable(t *testing.T) {
	errBoom := errors.New("boom")

	never := NeverRetry()
	for i := 0; i < 3; i++ {
		if never.Retryable(errBoom) {
			t.Fatalf("NeverRetry classifier returned true on iteration %d, want false", i)
		}
	}

	all := RetryAll()
	for i := 0; i < 3; i++ {
		if !all.Retryable(errBoom) {
			t.Fatalf("RetryAll classifier returned false on iteration %d, want true", i)
		}
	}
}
