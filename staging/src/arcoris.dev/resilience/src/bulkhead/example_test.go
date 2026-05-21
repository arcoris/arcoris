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

	lease, _, ok := b.TryAcquire()
	fmt.Println(ok)

	denied, snap, ok := b.TryAcquire()
	fmt.Println(denied == nil, ok, snap.Value.Available)

	lease.Release()
	_, snap, ok = b.TryAcquire()
	fmt.Println(ok, snap.Value.Available)

	// Output:
	// true
	// true false 0
	// true 0
}

func ExampleBulkhead_TryAdmit() {
	b := bulkhead.New(1)

	result := b.TryAdmit(bulkhead.Request{Amount: 1})
	fmt.Println(result.IsAdmitted(), result.HasGrant())

	if lease, ok := result.Grant(); ok {
		lease.Release()
	}

	// Output:
	// true true
}
