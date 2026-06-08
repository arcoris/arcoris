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

package bulkhead_test

import (
	"fmt"

	"arcoris.dev/resilience/bulkhead"
)

func ExampleBulkhead_TryAcquire() {
	b := bulkhead.New(1)

	lease, observation, ok := b.TryAcquire()
	fmt.Println(ok, observation.Refusal, observation.Snapshot.Value.Available)

	denied, observation, ok := b.TryAcquire()
	fmt.Println(denied == nil, ok, observation.Refusal, observation.Snapshot.Value.Available)

	lease.Release()
	_, observation, ok = b.TryAcquire()
	fmt.Println(ok, observation.Refusal, observation.Snapshot.Value.Available)

	// Output:
	// true none 0
	// true false insufficient 0
	// true none 0
}

func ExampleBulkhead_TryAcquireAmount() {
	b := bulkhead.New(3)

	lease, observation, ok := b.TryAcquireAmount(2)
	fmt.Println(ok, lease.Amount(), observation.Snapshot.Value.Available)

	lease.Release()

	// Output:
	// true 2 1
}

func ExampleBulkhead_SetLimit_debt() {
	b := bulkhead.New(2)

	first, _, _ := b.TryAcquire()
	second, _, _ := b.TryAcquire()

	debt := b.SetLimit(1)
	fmt.Println(debt.Value.Reserved, debt.Value.Debt)

	_, denied, ok := b.TryAcquire()
	fmt.Println(ok, denied.Refusal)

	first.Release()
	second.Release()

	// Output:
	// 2 1
	// false debt
}
