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


package capacity_test

import (
	"fmt"

	"arcoris.dev/capacity"
)

func ExampleLedger_TryReserve() {
	ledger := capacity.NewLedger(2)

	first, _, ok := ledger.TryReserve(1)
	fmt.Println(ok)

	second, _, ok := ledger.TryReserve(1)
	fmt.Println(ok)

	_, denied, ok := ledger.TryReserve(1)
	fmt.Println(ok, denied.Value.Available)

	first.Release()
	second.Release()

	snap := ledger.Snapshot()
	fmt.Println(snap.Value.Available)

	// Output:
	// true
	// true
	// false 0
	// 2
}

func ExampleLedger_SetLimit() {
	ledger := capacity.NewLedger(10)
	reservation, _, _ := ledger.TryReserve(8)

	overcommitted := ledger.SetLimit(5)
	fmt.Println(overcommitted.Value.Available, overcommitted.Value.Debt)

	restored := reservation.Release()
	fmt.Println(restored.Value.Available, restored.Value.Debt)

	// Output:
	// 0 3
	// 5 0
}
