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

package deadline

import (
	"context"
	"time"
)

// WithBoundedTimeout creates a child context with a deadline no later than both
// now+timeout and the parent context deadline.
//
// WithBoundedTimeout is similar to context.WithTimeout, but it computes the
// intended child deadline from the caller-supplied observation time. Passing the
// observation time explicitly keeps tests deterministic and makes the budget
// boundary visible to higher-level resilience code.
func WithBoundedTimeout(ctx context.Context, now time.Time, timeout time.Duration) (context.Context, context.CancelFunc) {
	requireContext(ctx)
	requireNonNegativeDuration("timeout", timeout)

	childDeadline := now.Add(timeout)
	if parentDeadline, ok := ctx.Deadline(); ok && parentDeadline.Before(childDeadline) {
		childDeadline = parentDeadline
	}

	return context.WithDeadline(ctx, childDeadline)
}
