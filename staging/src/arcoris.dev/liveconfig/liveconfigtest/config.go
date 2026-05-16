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

// Config is a deterministic test fixture for live configuration packages.
//
// Config deliberately contains both value fields and mutable aggregate fields.
// Tests can use it to verify clone isolation, equality checks, validation before
// publication, and snapshot value stability after callers mutate their original
// input values.
//
// Config is a test helper type. It is not a production configuration schema.
type Config struct {
	// Name identifies the test configuration.
	Name string

	// Version is a non-negative revision-like domain value used by tests to model
	// domain-level configuration changes independently from snapshot revisions.
	Version int

	// Enabled is a simple boolean policy toggle.
	Enabled bool

	// Timeout is a positive duration used by validation tests.
	Timeout time.Duration

	// Limits contains mutable numeric policy values.
	Limits []int

	// Labels contains mutable string metadata.
	Labels map[string]string
}

// NewConfig returns a valid deterministic Config for tests.
//
// The returned value contains allocated slices and maps so clone-isolation tests
// can mutate it without nil checks.
func NewConfig() Config {
	return Config{
		Name:    "default",
		Version: 1,
		Enabled: true,
		Timeout: time.Second,
		Limits:  []int{1, 2, 3},
		Labels: map[string]string{
			"env":   "test",
			"owner": "liveconfigtest",
		},
	}
}

// NewConfigVersion returns a valid deterministic Config with version.
func NewConfigVersion(version int) Config {
	cfg := NewConfig()
	cfg.Version = version
	return cfg
}

// Clone returns a deep copy of cfg.
//
// Clone is the method form of CloneConfig. It is convenient in tests that build
// value-style variants while keeping the original fixture isolated from later
// mutations.
func (cfg Config) Clone() Config {
	return CloneConfig(cfg)
}

// Equal reports whether cfg and other are equal according to EqualConfig.
func (cfg Config) Equal(other Config) bool {
	return EqualConfig(cfg, other)
}

// Validate checks cfg according to liveconfigtest fixture rules.
func (cfg Config) Validate() error {
	return ValidateConfig(cfg)
}

// WithName returns a cloned copy of cfg with Name replaced.
func (cfg Config) WithName(name string) Config {
	out := cfg.Clone()
	out.Name = name
	return out
}

// WithVersion returns a cloned copy of cfg with Version replaced.
func (cfg Config) WithVersion(version int) Config {
	out := cfg.Clone()
	out.Version = version
	return out
}

// WithEnabled returns a cloned copy of cfg with Enabled replaced.
func (cfg Config) WithEnabled(enabled bool) Config {
	out := cfg.Clone()
	out.Enabled = enabled
	return out
}

// WithTimeout returns a cloned copy of cfg with Timeout replaced.
func (cfg Config) WithTimeout(timeout time.Duration) Config {
	out := cfg.Clone()
	out.Timeout = timeout
	return out
}

// WithLimits returns a cloned copy of cfg with Limits replaced.
//
// The supplied limits are copied so later caller-side mutation of the argument
// slice cannot affect the returned Config.
func (cfg Config) WithLimits(limits ...int) Config {
	out := cfg.Clone()
	out.Limits = append([]int(nil), limits...)
	return out
}

// WithLabels returns a cloned copy of cfg with Labels replaced.
//
// The supplied labels map is copied. A nil labels map remains nil.
func (cfg Config) WithLabels(labels map[string]string) Config {
	out := cfg.Clone()
	out.Labels = cloneLabels(labels)
	return out
}

// WithLabel returns a cloned copy of cfg with one label set.
func (cfg Config) WithLabel(key, val string) Config {
	out := cfg.Clone()
	if out.Labels == nil {
		out.Labels = make(map[string]string)
	}
	out.Labels[key] = val
	return out
}

// WithoutLabel returns a cloned copy of cfg with one label removed.
func (cfg Config) WithoutLabel(key string) Config {
	out := cfg.Clone()
	delete(out.Labels, key)
	return out
}
