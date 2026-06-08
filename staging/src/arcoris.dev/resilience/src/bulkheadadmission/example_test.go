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

package bulkheadadmission_test

import (
	"fmt"

	"arcoris.dev/resilience/bulkhead"
	"arcoris.dev/resilience/bulkheadadmission"
)

func ExampleAdmitter_TryAdmit() {
	b := bulkhead.New(1)
	admitter := bulkheadadmission.New(b)

	result := admitter.TryAdmit(bulkheadadmission.Request{Amount: 1})
	fmt.Println(result.Decision().IsAdmitted(), result.HasGrant())

	if lease, ok := result.Grant(); ok {
		lease.Release()
	}

	// Output:
	// true true
}
