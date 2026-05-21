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

package retrybudget_test

import (
	"fmt"
	"math"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

type exampleBudget struct {
	snap snapshot.Snapshot[retrybudget.Snapshot]
}

func newExampleBudget() *exampleBudget {
	return &exampleBudget{
		snap: snapshot.Snapshot[retrybudget.Snapshot]{
			Revision: snapshot.ZeroRevision.Next(),
			Value: retrybudget.Snapshot{
				Kind: retrybudget.KindNoop,
				Capacity: retrybudget.CapacitySnapshot{
					Allowed:   math.MaxUint64,
					Available: math.MaxUint64,
				},
			},
		},
	}
}

func (b *exampleBudget) RecordOriginal() {}

func (b *exampleBudget) TryAdmitRetry() retrybudget.Decision {
	return retrybudget.Decision{Allowed: true, Reason: retrybudget.ReasonAllowed, Snapshot: b.snap}
}

func (b *exampleBudget) Snapshot() snapshot.Snapshot[retrybudget.Snapshot] {
	return b.snap
}

func ExampleBudget() {
	var budget retrybudget.Budget = newExampleBudget()

	budget.RecordOriginal()
	decision := budget.TryAdmitRetry()
	snap := budget.Snapshot()

	fmt.Println(decision.Allowed)
	fmt.Println(decision.Snapshot.Revision == snap.Revision)
	fmt.Println(snap.Value.Kind)
	// Output:
	// true
	// true
	// noop
}
