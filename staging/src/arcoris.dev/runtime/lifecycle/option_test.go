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

import (
	"testing"
	"time"
)

func TestDefaultControllerConfigUsesNonNilClock(t *testing.T) {
	t.Parallel()

	// The construction-time config boundary guarantees Controller always has a
	// usable time source even when callers provide no options.
	config := defaultControllerConfig()
	if config.now == nil {
		t.Fatal("default now = nil, want time source")
	}
	if config.now().IsZero() {
		t.Fatal("default now returned zero time")
	}
}

func TestNewControllerConfigIgnoresNilOptions(t *testing.T) {
	t.Parallel()

	config := newControllerConfig(nil)
	if config.now == nil {
		t.Fatal("config now = nil, want default")
	}
}

func TestNewControllerConfigAppliesOptionsInOrder(t *testing.T) {
	t.Parallel()

	var order []string
	first := Option(func(*controllerConfig) { order = append(order, "first") })
	second := Option(func(*controllerConfig) { order = append(order, "second") })
	newControllerConfig(first, second)
	assertDeepEqual(t, order, []string{"first", "second"})
}

func TestNewControllerConfigIndependentFromOptionSlice(t *testing.T) {
	t.Parallel()

	// The variadic option slice is consumed during construction; mutating the
	// caller's slice afterward must not rewrite the already returned config.
	opts := []Option{
		func(config *controllerConfig) { config.now = func() time.Time { return testTime } },
	}
	config := newControllerConfig(opts...)
	opts[0] = func(config *controllerConfig) { config.now = func() time.Time { return time.Time{} } }

	if got := config.now(); !got.Equal(testTime) {
		t.Fatalf("config.now() = %v, want %v", got, testTime)
	}
}
