// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fixedwindow

import (
	"errors"
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
		WithRatio(MustRatio(1, 2)),
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
	if cfg.ratio != MustRatio(1, 2) {
		t.Fatalf("ratio = %v, want 1/2", cfg.ratio)
	}
	if cfg.minRetries != 3 {
		t.Fatalf("minRetries = %d, want 3", cfg.minRetries)
	}
}

func TestNewConfigRejectsNilOption(t *testing.T) {
	if _, err := newConfig(nil); !errors.Is(err, ErrNilOption) {
		t.Fatalf("newConfig(nil) error = %v, want %v", err, ErrNilOption)
	}
}

func TestNewConfigValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want error
	}{
		{name: "nil option", opts: []Option{nil}, want: ErrNilOption},
		{name: "nil clock", opts: []Option{WithClock(nil)}, want: ErrNilClock},
		{name: "zero window", opts: []Option{WithWindow(0)}, want: ErrInvalidWindow},
		{name: "negative window", opts: []Option{WithWindow(-time.Second)}, want: ErrInvalidWindow},
		{name: "invalid ratio", opts: []Option{WithRatio(Ratio{})}, want: ErrInvalidRatio},
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

func TestNewConfigOptionsApplyInOrder(t *testing.T) {
	cfg, err := newConfig(
		WithWindow(time.Second),
		WithWindow(2*time.Second),
		WithRatio(MustRatio(1, 4)),
		WithRatio(MustRatio(3, 4)),
		WithMinRetries(1),
		WithMinRetries(3),
	)
	if err != nil {
		t.Fatalf("newConfig() error = %v", err)
	}
	if cfg.window != 2*time.Second {
		t.Fatalf("window = %s, want 2s", cfg.window)
	}
	if cfg.ratio != MustRatio(3, 4) {
		t.Fatalf("ratio = %v, want 3/4", cfg.ratio)
	}
	if cfg.minRetries != 3 {
		t.Fatalf("minRetries = %d, want 3", cfg.minRetries)
	}
}
