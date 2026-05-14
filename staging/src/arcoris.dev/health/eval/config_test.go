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

package eval

import "testing"

func TestDefaultEvaluatorConfig(t *testing.T) {
	t.Parallel()

	cfg := defaultEvaluatorConfig()
	if cfg.clock == nil {
		t.Fatal("default clock is nil")
	}
	if cfg.defaultTimeout != defaultCheckTimeout {
		t.Fatalf("default timeout = %s, want %s", cfg.defaultTimeout, defaultCheckTimeout)
	}
	if cfg.targetTimeouts == nil {
		t.Fatal("target timeouts map is nil")
	}
}
