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

package maybe_test

import (
	"fmt"

	"arcoris.dev/value/maybe"
)

func ExampleMaybe_Load() {
	m := maybe.Some("ready")

	val, ok := m.Load()
	fmt.Println(val, ok)

	// Output:
	// ready true
}

func ExampleNone() {
	m := maybe.None[string]()

	val, ok := m.Load()
	fmt.Println(val, ok)

	// Output:
	//  false
}
