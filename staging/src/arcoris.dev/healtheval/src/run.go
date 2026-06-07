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

package eval

import (
	"context"
	"time"

	"arcoris.dev/health"
)

// evaluateCheck executes one checker and applies evaluator-owned normalization.
func (e *Evaluator) evaluateCheck(ctx context.Context, check health.Checker, timeout time.Duration) health.Result {
	started := e.clock.Now()

	name, err := health.CheckerName(check)
	if err != nil {
		finished := e.clock.Now()
		d := nonNegativeDuration(e.clock.Since(started))

		return invalidCheckerResult(name, err, finished, d)
	}

	res := e.runCheck(ctx, check, name, timeout)

	finished := e.clock.Now()
	d := nonNegativeDuration(e.clock.Since(started))

	return normalizeEvaluatedResult(res, name, finished, d)
}
