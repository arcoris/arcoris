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

package atomicx

import "fmt"

// ExampleUint64Counter demonstrates samples and deltas for lifetime counters.
func ExampleUint64Counter() {
	var counter Uint64Counter

	counter.Inc()
	prev := counter.Sample()

	counter.Add(4)
	cur := counter.Sample()

	delta := cur.DeltaSince(prev)

	fmt.Println(counter.Load())
	fmt.Println(delta.Value)
	fmt.Println(delta.Wrapped)

	// Output:
	// 5
	// 4
	// false
}

// ExampleUint64Gauge demonstrates non-negative current-state accounting.
func ExampleUint64Gauge() {
	var gauge Uint64Gauge

	gauge.Set(0)
	gauge.Add(10)
	gauge.Sub(3)

	cur, ok := gauge.TryAdd(5)

	fmt.Println(cur)
	fmt.Println(ok)
	fmt.Println(gauge.Load())

	// Output:
	// 12
	// true
	// 12
}

// ExampleInt64Gauge demonstrates signed current-state accounting.
func ExampleInt64Gauge() {
	var gauge Int64Gauge

	gauge.Set(0)
	gauge.Add(10)
	gauge.Sub(15)

	cur, ok := gauge.TryAdd(3)

	fmt.Println(cur)
	fmt.Println(ok)
	fmt.Println(gauge.Load())

	// Output:
	// -2
	// true
	// -2
}

// ExamplePaddedUint64 demonstrates raw padded arithmetic without gauge invariants.
func ExamplePaddedUint64() {
	var cell PaddedUint64

	cell.Store(^uint64(0))
	cell.Inc()

	fmt.Println(cell.Load())

	// Output:
	// 0
}

// ExamplePaddedPointer demonstrates raw pointer publication without ownership transfer.
func ExamplePaddedPointer() {
	type policy struct {
		name string
	}

	var current PaddedPointer[policy]
	published := &policy{name: "active"}

	current.Store(published)

	fmt.Println(current.Load() == published)
	fmt.Println(current.Load().name)

	// Output:
	// true
	// active
}
