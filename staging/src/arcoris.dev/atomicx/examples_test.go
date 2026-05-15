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

package atomicx

import "fmt"

// ExampleUint64Counter demonstrates snapshots and deltas for lifetime counters.
func ExampleUint64Counter() {
	var counter Uint64Counter

	counter.Inc()
	prev := counter.Snapshot()

	counter.Add(4)
	cur := counter.Snapshot()

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

// ExampleUint32Counter demonstrates explicitly bounded lifetime counter usage.
func ExampleUint32Counter() {
	var counter Uint32Counter

	counter.Add(10)
	prev := counter.Snapshot()

	counter.Add(7)
	cur := counter.Snapshot()

	delta := cur.DeltaSince(prev)

	fmt.Println(counter.Load())
	fmt.Println(delta.Value)

	// Output:
	// 17
	// 7
}
