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

package fixedwindow

import (
	"sync"
	"time"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

// Limiter is a local fixed-window retry budget.
//
// Limiter records original attempts and admits retry attempts according to the
// configured ratio/minimum policy for the current fixed window. It is safe for
// concurrent use. Limiter must not be copied after first use.
type Limiter struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects windowStart, original, and retries.
	mu sync.Mutex

	// cfg contains validated limiter policy and clock dependencies.
	cfg config

	// windowStart is the inclusive start of the current accounting window.
	windowStart time.Time

	// original is the number of original, non-retry attempts observed in the
	// current window.
	original uint64

	// retries is the number of retry attempts admitted in the current window.
	retries uint64

	// published exposes immutable retry-budget snapshots to readers.
	published snapshot.Publisher[retrybudget.Snapshot]
}
