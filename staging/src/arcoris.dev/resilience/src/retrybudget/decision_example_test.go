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

package retrybudget_test

import (
	"fmt"
	"time"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

func ExampleDecision() {
	decision := retrybudget.Decision{
		Allowed: true,
		Reason:  retrybudget.ReasonAllowed,
		Snapshot: snapshot.Snapshot[retrybudget.Snapshot]{
			Revision: snapshot.ZeroRevision.Next(),
			Value: retrybudget.Snapshot{
				Kind: retrybudget.KindFixedWindow,
				Attempts: retrybudget.AttemptsSnapshot{
					Original: 10,
					Retry:    2,
				},
				Capacity: retrybudget.CapacitySnapshot{
					Allowed:   3,
					Available: 1,
					Exhausted: false,
				},
				Window: retrybudget.WindowSnapshot{
					StartedAt: time.Unix(100, 0).UTC(),
					EndsAt:    time.Unix(160, 0).UTC(),
					Duration:  time.Minute,
					Bounded:   true,
				},
				Policy: retrybudget.PolicySnapshot{
					Ratio:   0.2,
					Minimum: 1,
					Bounded: true,
				},
			},
		},
	}

	fmt.Println(decision.IsAllowed())
	fmt.Println(decision.IsValid())

	// Output:
	// true
	// true
}
