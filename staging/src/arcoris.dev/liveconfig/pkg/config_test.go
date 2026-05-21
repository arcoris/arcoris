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

package liveconfig

import (
	"testing"

	"arcoris.dev/chrono/clock"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig[testConfig]()
	if cfg.clock == nil {
		t.Fatal("default clock is nil")
	}
	if _, ok := cfg.clock.(clock.RealClock); !ok {
		t.Fatalf("default clock type = %T, want clock.RealClock", cfg.clock)
	}
	if cfg.clone == nil {
		t.Fatal("default clone is nil")
	}
	value := testConfig{Name: "value", Limit: 3}
	if got := cfg.clone(value); got.Name != value.Name || got.Limit != value.Limit {
		t.Fatalf("default clone = %+v, want %+v", got, value)
	}
	if cfg.normalize != nil {
		t.Fatal("default normalizer is not nil")
	}
	if cfg.validate != nil {
		t.Fatal("default validator is not nil")
	}
	if cfg.equal != nil {
		t.Fatal("default equal is not nil")
	}
}

func TestNewConfigAppliesOptionsInOrder(t *testing.T) {
	first := newTestClock()
	second := testClock{now: first.now.Add(1)}

	cfg := newConfig(WithClock[testConfig](first), WithClock[testConfig](second))
	got, ok := cfg.clock.(testClock)
	if !ok {
		t.Fatalf("newConfig clock type = %T, want testClock", cfg.clock)
	}
	if !got.now.Equal(second.now) {
		t.Fatalf("newConfig clock time = %s, want %s", got.now, second.now)
	}
}
