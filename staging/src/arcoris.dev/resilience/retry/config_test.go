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
	cfg := defaultConfig()

	requireClock(cfg.clock)
	requireDelaySchedule(cfg.delay)
	requireClassifier(cfg.classifier)

	if cfg.maxAttempts != 1 {
		t.Fatalf("default maxAttempts = %d, want 1", cfg.maxAttempts)
	}
	if cfg.maxElapsed != 0 {
		t.Fatalf("default maxElapsed = %s, want 0", cfg.maxElapsed)
	}
	if len(cfg.observers) != 0 {
		t.Fatalf("default observers len = %d, want 0", len(cfg.observers))
	}
	if cfg.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("default classifier retried error, want conservative NeverRetry behavior")
	}

	seq := cfg.delay.NewSequence()
	delay, ok := seq.Next()
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

	cfg := configOf(
		WithMaxAttempts(2),
		WithMaxAttempts(3),
		WithMaxElapsed(time.Second),
		WithMaxElapsed(2*time.Second),
		WithClassifier(firstClassifier),
		WithClassifier(secondClassifier),
	)

	if cfg.maxAttempts != 3 {
		t.Fatalf("maxAttempts = %d, want 3", cfg.maxAttempts)
	}
	if cfg.maxElapsed != 2*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", cfg.maxElapsed, 2*time.Second)
	}
	if !cfg.classifier.Retryable(errors.New("boom")) {
		t.Fatalf("last classifier option did not win")
	}
}
