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

package core

import "testing"

func TestStrategyAndMergeDefaultsRemainStable(t *testing.T) {
	if StrategyAuto != 0 {
		t.Fatalf("StrategyAuto = %d, want 0", StrategyAuto)
	}
	if MergeLinear != 0 {
		t.Fatalf("MergeLinear = %d, want 0", MergeLinear)
	}
	if StrategyDynamic <= StrategyStatic {
		t.Fatalf("strategy order changed: static=%d dynamic=%d", StrategyStatic, StrategyDynamic)
	}
}
