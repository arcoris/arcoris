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

package fixedwindow_test

import (
	"fmt"

	"arcoris.dev/resilience/retrybudget/fixedwindow"
)

func Example() {
	budget, err := fixedwindow.New(
		fixedwindow.WithRatio(1),
		fixedwindow.WithMinRetries(0),
	)
	if err != nil {
		panic(err)
	}

	budget.RecordOriginal()
	decision := budget.TryAdmitRetry()
	snap := decision.Snapshot

	fmt.Println(decision.Allowed)
	fmt.Println(decision.Reason)
	fmt.Println(snap.Value.Kind)
	fmt.Println(snap.Value.Attempts.Original)
	fmt.Println(snap.Value.Attempts.Retry)

	// Output:
	// true
	// allowed
	// fixed_window
	// 1
	// 1
}
