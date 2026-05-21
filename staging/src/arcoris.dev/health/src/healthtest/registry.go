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

package healthtest

import (
	"testing"

	"arcoris.dev/health"
)

// TargetChecks groups checkers under one health target.
//
// The type keeps registry setup deterministic. Tests pass explicit groups in
// order instead of relying on map iteration when report order matters.
type TargetChecks struct {
	// Target is the concrete health target to register checks under.
	Target health.Target

	// Checks are registered in the supplied order.
	Checks []health.Checker
}

// ForTarget returns a deterministic target registration group.
//
// The checks slice is copied so later caller mutations do not change the group
// that NewRegistry will register.
func ForTarget(target health.Target, checks ...health.Checker) TargetChecks {
	copied := make([]health.Checker, len(checks))
	copy(copied, checks)

	return TargetChecks{Target: target, Checks: copied}
}

// NewRegistry returns a registry populated with groups.
//
// Registration failures fail the current test immediately because this helper is
// for fixture setup. Tests that need to inspect registration errors should call
// health.Registry.Register directly.
func NewRegistry(t testing.TB, groups ...TargetChecks) *health.Registry {
	t.Helper()

	registry := health.NewRegistry()
	for _, group := range groups {
		Register(t, registry, group.Target, group.Checks...)
	}

	return registry
}

// Register adds checks to registry and fails the test on error.
//
// The helper is intentionally narrow: it owns only health-domain registration
// setup. It does not provide generic error assertions or hide registry behavior
// in production-style constructors.
func Register(t testing.TB, r *health.Registry, target health.Target, checks ...health.Checker) {
	t.Helper()

	if r == nil {
		t.Fatalf("healthtest.Register(%s) registry = nil", target)
	}
	if err := r.Register(target, checks...); err != nil {
		t.Fatalf("healthtest.Register(%s) = %v, want nil", target, err)
	}
}
