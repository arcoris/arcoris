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

// Remaining returns the non-negative duration until ctx's deadline at now.
//
// The boolean result reports whether ctx had a deadline. When ctx has no
// deadline, Remaining returns zero, false. When the deadline has expired,
// Remaining returns zero, true.
func Remaining(ctx context.Context, now time.Time) (time.Duration, bool) {
	budget := Inspect(ctx, now)
	if !budget.HasDeadline {
		return 0, false
	}
	return budget.Remaining, true
}
