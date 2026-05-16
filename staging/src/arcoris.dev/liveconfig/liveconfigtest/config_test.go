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

package liveconfigtest

import (
	"errors"
	"slices"
	"testing"
	"time"
)

func TestCloneConfigIsolatesMutableFields(t *testing.T) {
	orig := NewConfig()
	clone := CloneConfig(orig)

	MutateConfig(&orig)

	want := NewConfig()
	RequireConfigEqual(t, clone, want)
}

func TestConfigMethodCloneIsolatesMutableFields(t *testing.T) {
	orig := NewConfig()
	clone := orig.Clone()

	MutateConfig(&orig)

	RequireConfigEqual(t, clone, NewConfig())
}

func TestConfigValueMethods(t *testing.T) {
	labels := map[string]string{"env": "stage"}
	cfg := NewConfig().
		WithName("custom").
		WithVersion(7).
		WithEnabled(false).
		WithTimeout(2*time.Second).
		WithLimits(4, 5).
		WithLabels(labels).
		WithLabel("owner", "team").
		WithoutLabel("env")

	labels["env"] = "mutated"

	if got, want := cfg.Name, "custom"; got != want {
		t.Fatalf("Name = %q, want %q", got, want)
	}
	if got, want := cfg.Version, 7; got != want {
		t.Fatalf("Version = %d, want %d", got, want)
	}
	if cfg.Enabled {
		t.Fatal("Enabled = true, want false")
	}
	if got, want := cfg.Timeout, 2*time.Second; got != want {
		t.Fatalf("Timeout = %s, want %s", got, want)
	}
	if got, want := cfg.Limits, []int{4, 5}; !slices.Equal(got, want) {
		t.Fatalf("Limits = %#v, want %#v", got, want)
	}
	if _, ok := cfg.Labels["env"]; ok {
		t.Fatal("env label is present, want removed")
	}
	if got, want := cfg.Labels["owner"], "team"; got != want {
		t.Fatalf("owner label = %q, want %q", got, want)
	}
	if got := cfg.Labels["env"]; got == "mutated" {
		t.Fatal("WithLabels retained caller-owned map")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !cfg.Equal(CloneConfig(cfg)) {
		t.Fatal("Equal() = false for cloned config")
	}
}

func TestMutatedConfigLeavesOriginalUnchanged(t *testing.T) {
	orig := NewConfig()
	mutated := MutatedConfig(orig)

	if !EqualConfig(orig, NewConfig()) {
		t.Fatal("MutatedConfig changed original config")
	}
	if EqualConfig(orig, mutated) {
		t.Fatal("MutatedConfig returned equal config, want changed")
	}
}

func TestEqualConfig(t *testing.T) {
	base := NewConfig()

	tests := []struct {
		name   string
		mutate func(*Config)
		want   bool
	}{
		{
			name: "equal",
			want: true,
		},
		{
			name:   "different name",
			mutate: func(cfg *Config) { cfg.Name = "other" },
		},
		{
			name:   "different limit",
			mutate: func(cfg *Config) { cfg.Limits[0] = 99 },
		},
		{
			name:   "different label",
			mutate: func(cfg *Config) { cfg.Labels["env"] = "prod" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CloneConfig(base)
			if tt.mutate != nil {
				tt.mutate(&got)
			}
			if EqualConfig(base, got) != tt.want {
				t.Fatalf("EqualConfig() = %v, want %v", EqualConfig(base, got), tt.want)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want error
	}{
		{
			name: "valid",
			cfg:  NewConfig(),
		},
		{
			name: "blank name",
			cfg:  InvalidNameConfig(),
			want: ErrInvalidName,
		},
		{
			name: "negative version",
			cfg:  InvalidVersionConfig(),
			want: ErrInvalidVersion,
		},
		{
			name: "zero timeout",
			cfg:  InvalidTimeoutConfig(),
			want: ErrInvalidTimeout,
		},
		{
			name: "negative limit",
			cfg:  InvalidLimitConfig(),
			want: ErrInvalidLimit,
		},
		{
			name: "positive timeout",
			cfg: func() Config {
				cfg := NewConfig()
				cfg.Timeout = time.Nanosecond
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if !errors.Is(err, tt.want) {
				t.Fatalf("ValidateConfig() error = %v, want %v", err, tt.want)
			}
		})
	}
}
