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

package snapshot_test

import (
	"fmt"
	"slices"

	"arcoris.dev/snapshot"
)

func ExamplePublisher() {
	var publisher snapshot.Publisher[[]string]

	handlers := []string{"first", "second"}
	publisher.Publish(slices.Clone(handlers))

	snap := publisher.Snapshot()
	fmt.Println(snap.Revision)
	fmt.Println(snap.Value[0])

	// Output:
	// 1
	// first
}
