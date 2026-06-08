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

package retrybudgetadmission

import (
	"sync"
	"sync/atomic"
	"time"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

func validAllowedDecision() retrybudget.Decision {
	return retrybudget.Decision{
		Allowed:  true,
		Reason:   retrybudget.ReasonAllowed,
		Snapshot: validSnapshot(false, 1),
	}
}

func validDeniedDecision() retrybudget.Decision {
	return retrybudget.Decision{
		Allowed:  false,
		Reason:   retrybudget.ReasonExhausted,
		Snapshot: validSnapshot(true, 0),
	}
}

func validSnapshot(exhausted bool, retry uint64) snapshot.Snapshot[retrybudget.Snapshot] {
	available := uint64(1)
	if exhausted {
		available = 0
	}
	return snapshot.Snapshot[retrybudget.Snapshot]{
		Revision: snapshot.ZeroRevision.Next(),
		Value: retrybudget.Snapshot{
			Kind: retrybudget.KindFixedWindow,
			Attempts: retrybudget.AttemptsSnapshot{
				Original: 1,
				Retry:    retry,
			},
			Capacity: retrybudget.CapacitySnapshot{
				Allowed:   1,
				Available: available,
				Exhausted: exhausted,
			},
			Window: retrybudget.WindowSnapshot{
				StartedAt: time.Unix(100, 0).UTC(),
				EndsAt:    time.Unix(160, 0).UTC(),
				Duration:  time.Minute,
				Bounded:   true,
			},
			Policy: retrybudget.PolicySnapshot{
				Ratio:   retrybudget.RatioOne,
				Minimum: 0,
				Bounded: true,
			},
		},
	}
}

type scriptedBudget struct {
	mu        sync.Mutex
	calls     atomic.Uint64
	decisions []retrybudget.Decision
}

func (b *scriptedBudget) TryAdmitRetry() retrybudget.Decision {
	b.calls.Add(1)

	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.decisions) == 0 {
		return validDeniedDecision()
	}
	decision := b.decisions[0]
	b.decisions = b.decisions[1:]
	return decision
}

type countingBudget struct {
	mu    sync.Mutex
	limit uint64
	used  uint64
}

func (b *countingBudget) TryAdmitRetry() retrybudget.Decision {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.used >= b.limit {
		return validDeniedDecision()
	}
	b.used++
	return validAllowedDecision()
}
