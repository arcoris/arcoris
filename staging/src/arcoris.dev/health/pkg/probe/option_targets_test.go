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

package probe

import (
	"testing"

	"arcoris.dev/health"
)

func TestWithTargets(t *testing.T) {
	t.Parallel()

	source := []health.Target{health.TargetReady, health.TargetLive}
	cfg := defaultConfig()
	err := WithTargets(source...)(&cfg)
	if err != nil {
		t.Fatalf("WithTargets() = %v, want nil", err)
	}
	if !sameTargets(cfg.targets, []health.Target{health.TargetReady, health.TargetLive}) {
		t.Fatalf("targets = %v, want [ready live]", cfg.targets)
	}

	source[0] = health.TargetStartup
	if cfg.targets[0] != health.TargetReady {
		t.Fatalf("targets share caller backing array: %v", cfg.targets)
	}
}
