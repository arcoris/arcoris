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

package fixedwindow

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.clock == nil {
		t.Fatal("default clock is nil")
	}
	if cfg.window != DefaultWindow {
		t.Fatalf("window = %s, want %s", cfg.window, DefaultWindow)
	}
	if cfg.ratio != DefaultRatio {
		t.Fatalf("ratio = %v, want %v", cfg.ratio, DefaultRatio)
	}
	if cfg.minRetries != DefaultMinRetries {
		t.Fatalf("minRetries = %d, want %d", cfg.minRetries, DefaultMinRetries)
	}
}

func TestNewConfigAppliesOptions(t *testing.T) {
	clk := newFakeClock(fixedWindowTestNow)
	cfg, err := newConfig(
		WithClock(clk),
		WithWindow(30*time.Second),
		WithRatio(0.5),
		WithMinRetries(3),
	)
	if err != nil {
		t.Fatalf("newConfig() error = %v", err)
	}
	if cfg.clock != clk {
		t.Fatal("clock option was not applied")
	}
	if cfg.window != 30*time.Second {
		t.Fatalf("window = %s, want 30s", cfg.window)
	}
	if cfg.ratio != 0.5 {
		t.Fatalf("ratio = %v, want 0.5", cfg.ratio)
	}
	if cfg.minRetries != 3 {
		t.Fatalf("minRetries = %d, want 3", cfg.minRetries)
	}
}

func TestNewConfigIgnoresNilOption(t *testing.T) {
	if _, err := newConfig(nil); err != nil {
		t.Fatalf("newConfig(nil) error = %v", err)
	}
}

func TestNewConfigValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want error
	}{
		{name: "nil clock", opts: []Option{WithClock(nil)}, want: ErrNilClock},
		{name: "zero window", opts: []Option{WithWindow(0)}, want: ErrInvalidWindow},
		{name: "negative window", opts: []Option{WithWindow(-time.Second)}, want: ErrInvalidWindow},
		{name: "negative ratio", opts: []Option{WithRatio(-0.1)}, want: ErrInvalidRatio},
		{name: "ratio greater than one", opts: []Option{WithRatio(1.1)}, want: ErrInvalidRatio},
		{name: "ratio NaN", opts: []Option{WithRatio(math.NaN())}, want: ErrInvalidRatio},
		{name: "ratio positive infinity", opts: []Option{WithRatio(math.Inf(1))}, want: ErrInvalidRatio},
		{name: "ratio negative infinity", opts: []Option{WithRatio(math.Inf(-1))}, want: ErrInvalidRatio},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newConfig(tt.opts...)
			if !errors.Is(err, tt.want) {
				t.Fatalf("newConfig() error = %v, want %v", err, tt.want)
			}
		})
	}
}
