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

// Inspect derives a Budget from ctx at now.
//
// Inspect only inspects deadline math. It does not treat an already-canceled
// context without a deadline as an expired budget. Operational decisions that
// should reject canceled contexts are provided by CanStart, Clamp, and Reserve.
func Inspect(ctx context.Context, now time.Time) Budget {
	requireContext(ctx)

	dl, ok := ctx.Deadline()
	if !ok {
		return Budget{}
	}

	remaining := dl.Sub(now)
	if remaining <= 0 {
		return Budget{
			Deadline:    dl,
			HasDeadline: true,
			Expired:     true,
		}
	}

	return Budget{
		Deadline:    dl,
		Remaining:   remaining,
		HasDeadline: true,
	}
}
