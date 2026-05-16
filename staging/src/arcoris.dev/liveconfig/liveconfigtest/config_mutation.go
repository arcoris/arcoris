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

import "time"

// MutateConfig mutates cfg in place in ways that affect every field group.
//
// MutateConfig is useful for clone-isolation tests. Calling MutateConfig on an
// input value after publication must not change a correctly cloned or immutable
// published snapshot. A nil pointer is accepted so cleanup paths can call it
// without branching.
func MutateConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	cfg.Name = cfg.Name + "-mutated"
	cfg.Version++
	cfg.Enabled = !cfg.Enabled
	cfg.Timeout += time.Second
	cfg.Limits = append(cfg.Limits, 99)
	if cfg.Labels == nil {
		cfg.Labels = make(map[string]string)
	}
	cfg.Labels["mutated"] = "true"
}

// MutatedConfig returns a mutated clone of cfg.
//
// The original value is left unchanged. Use MutatedConfig when a test needs a
// clearly different candidate but also needs to keep the base fixture available
// for later assertions.
func MutatedConfig(cfg Config) Config {
	out := cfg.Clone()
	MutateConfig(&out)
	return out
}

// InvalidNameConfig returns a fixture rejected with ErrInvalidName.
func InvalidNameConfig() Config {
	return NewConfig().WithName(" ")
}

// InvalidVersionConfig returns a fixture rejected with ErrInvalidVersion.
func InvalidVersionConfig() Config {
	return NewConfig().WithVersion(-1)
}

// InvalidTimeoutConfig returns a fixture rejected with ErrInvalidTimeout.
func InvalidTimeoutConfig() Config {
	return NewConfig().WithTimeout(0)
}

// InvalidLimitConfig returns a fixture rejected with ErrInvalidLimit.
func InvalidLimitConfig() Config {
	return NewConfig().WithLimits(1, -1)
}
