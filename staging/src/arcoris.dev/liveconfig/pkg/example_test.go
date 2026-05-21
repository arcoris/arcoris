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

package liveconfig_test

import (
	"errors"
	"fmt"

	"arcoris.dev/liveconfig"
)

type retryPolicy struct {
	MaxAttempts int
	Ratio       float64
}

func ExampleHolder() {
	holder, err := liveconfig.New(
		retryPolicy{MaxAttempts: 3, Ratio: 0.2},
		liveconfig.WithValidator(func(policy retryPolicy) error {
			if policy.MaxAttempts < 1 {
				return errors.New("max attempts must be positive")
			}
			if policy.Ratio < 0 || policy.Ratio > 1 {
				return errors.New("ratio must be in [0,1]")
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	change, err := holder.Apply(retryPolicy{MaxAttempts: 4, Ratio: 0.1})
	if err != nil {
		panic(err)
	}

	fmt.Println(change.Changed)
	fmt.Println(holder.Snapshot().Value.MaxAttempts)

	// Output:
	// true
	// 4
}
