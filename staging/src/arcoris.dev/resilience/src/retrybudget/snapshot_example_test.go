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
	"math"

	"arcoris.dev/resilience/retrybudget"
)

func ExampleSnapshot() {
	snap := retrybudget.Snapshot{
		Kind: retrybudget.KindNoop,
		Capacity: retrybudget.CapacitySnapshot{
			Allowed:   math.MaxUint64,
			Available: math.MaxUint64,
			Exhausted: false,
		},
	}

	fmt.Println(snap.Kind)
	fmt.Println(snap.Exhausted())
	fmt.Println(snap.IsValid())

	// Output:
	// noop
	// false
	// true
}
