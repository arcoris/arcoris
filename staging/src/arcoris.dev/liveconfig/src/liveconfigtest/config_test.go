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

package liveconfigtest

import (
	"slices"
	"testing"
	"time"
)

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
