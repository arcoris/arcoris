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

package healthprobe

import (
	"testing"

	"arcoris.dev/component-base/pkg/health"
)

func TestWithTargets(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithTargets(health.TargetReady, health.TargetLive)(&cfg)
	if err != nil {
		t.Fatalf("WithTargets() = %v, want nil", err)
	}
	if !sameHealthprobeTargets(cfg.targets, []health.Target{health.TargetReady, health.TargetLive}) {
		t.Fatalf("targets = %v, want [ready live]", cfg.targets)
	}
}
