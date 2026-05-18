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

package planner

import "arcoris.dev/measure/internal/reduce/core"

// Plan selects the planning function implied by opts.
//
// Plan normalizes opts before dispatch. StrategyDynamicChunks deliberately
// produces the same inspectable fixed-size ranges as StrategyFixedChunks.
// Dynamic runner execution does not require a materialized plan and may claim
// chunks through an atomic cursor instead.
func Plan(
	n int,
	opts core.Options,
	dst []core.Range,
) []core.Range {
	opts = core.NormalizeOptions(opts)
	switch opts.Strategy {
	case core.StrategySequential:
		return Sequential(n, dst)
	case core.StrategyFixedChunks, core.StrategyDynamicChunks:
		return FixedChunks(n, opts, dst)
	default:
		return Balanced(n, opts, dst)
	}
}
