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

	"arcoris.dev/component-base/pkg/health"
)

func TestRegistryHelpers(t *testing.T) {
	t.Parallel()

	group := ForTarget(health.TargetReady, HealthyChecker("storage"), HealthyChecker("database"))
	registry := NewRegistry(t, group)

	if registry.Len(health.TargetReady) != 2 {
		t.Fatalf("Len(ready) = %d, want 2", registry.Len(health.TargetReady))
	}
	if !registry.Has(health.TargetReady, "storage") {
		t.Fatal("registry missing storage check")
	}

	group.Checks[0] = HealthyChecker("mutated")
	if registry.Has(health.TargetReady, "mutated") {
		t.Fatal("registry was affected by group mutation")
	}
}

func TestRegisterHelper(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	Register(t, registry, health.TargetLive, HealthyChecker("live"))

	if !registry.Has(health.TargetLive, "live") {
		t.Fatal("registry missing live check")
	}
}
