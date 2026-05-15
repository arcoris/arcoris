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

package deadline

import (
	"context"
	"fmt"
	"time"
)

const panicNilContext = "deadline: nil context"

// panicNegativeDuration formats a panic message for negative duration inputs.
func panicNegativeDuration(name string) string {
	return fmt.Sprintf("deadline: negative %s", name)
}

// requireContext rejects nil contexts at public API boundaries.
func requireContext(ctx context.Context) {
	if ctx == nil {
		panic(panicNilContext)
	}
}

// requireNonNegativeDuration rejects negative duration inputs at public API
// boundaries.
func requireNonNegativeDuration(name string, d time.Duration) {
	if d < 0 {
		panic(panicNegativeDuration(name))
	}
}
