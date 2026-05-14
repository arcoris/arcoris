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
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

	requireClock(config.clock)
	requireDelaySchedule(config.delay)
	requireClassifier(config.classifier)

	if config.maxAttempts != 1 {
		t.Fatalf("default maxAttempts = %d, want 1", config.maxAttempts)
	}
	if config.maxElapsed != 0 {
		t.Fatalf("default maxElapsed = %s, want 0", config.maxElapsed)
	}
	if len(config.observers) != 0 {
		t.Fatalf("default observers len = %d, want 0", len(config.observers))
	}
	if config.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("default classifier retried error, want conservative NeverRetry behavior")
	}

	sequence := config.delay.NewSequence()
	delay, ok := sequence.Next()
	if !ok {
		t.Fatalf("default delay exhausted, want immediate sequence")
	}
	if delay != 0 {
		t.Fatalf("default delay = %s, want 0", delay)
	}
}

func TestConfigOfAppliesOptionsInOrder(t *testing.T) {
	firstClassifier := ClassifierFunc(func(error) bool {
		return false
	})
	secondClassifier := ClassifierFunc(func(error) bool {
		return true
	})

	config := configOf(
		WithMaxAttempts(2),
		WithMaxAttempts(3),
		WithMaxElapsed(time.Second),
		WithMaxElapsed(2*time.Second),
		WithClassifier(firstClassifier),
		WithClassifier(secondClassifier),
	)

	if config.maxAttempts != 3 {
		t.Fatalf("maxAttempts = %d, want 3", config.maxAttempts)
	}
	if config.maxElapsed != 2*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", config.maxElapsed, 2*time.Second)
	}
	if !config.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("last classifier option did not win")
	}
}
