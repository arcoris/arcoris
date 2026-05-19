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

package bulkhead_test

import (
	"fmt"

	"arcoris.dev/resilience/bulkhead"
)

func ExampleLimiter_TryAcquire() {
	limiter, _ := bulkhead.New(2)

	first, firstDecision := limiter.TryAcquire()
	second, secondDecision := limiter.TryAcquire()
	third, thirdDecision := limiter.TryAcquire()

	fmt.Println(firstDecision.Allowed)
	fmt.Println(secondDecision.Allowed)
	fmt.Println(third == nil)
	fmt.Println(thirdDecision.Reason)

	first.Release()

	fourth, fourthDecision := limiter.TryAcquire()
	fmt.Println(fourthDecision.Allowed)

	second.Release()
	fourth.Release()

	// Output:
	// true
	// true
	// true
	// full
	// true
}
