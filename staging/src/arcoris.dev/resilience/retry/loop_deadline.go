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

package retry

import (
	"context"
	"time"

	"arcoris.dev/resilience/deadline"
)

// contextDeadlineWouldBeExceeded reports whether d would consume or exceed the
// remaining deadline budget of ctx at the retry clock's current time.
//
// The check is conservative. A delay equal to the remaining deadline budget is
// rejected because it would leave no usable budget for the next operation
// attempt.
//
// Already-stopped contexts return false because retry-owned context
// interruption is handled by the surrounding retry loop. This helper only owns
// predictive deadline-budget exhaustion for contexts that are still active at
// the retry boundary.
func (e *retryExecution) contextDeadlineWouldBeExceeded(ctx context.Context, d time.Duration) bool {
	requireContext(ctx)
	requireDelay(d, true)

	if ctx.Err() != nil {
		return false
	}

	remaining, bounded := deadline.Remaining(ctx, e.config.clock.Now())
	if !bounded {
		return false
	}

	return d >= remaining
}
