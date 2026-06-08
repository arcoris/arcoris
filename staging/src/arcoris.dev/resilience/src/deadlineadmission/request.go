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

package deadlineadmission

import (
	"context"
	"time"
)

// Request is the admission adapter request for one deadline start decision.
//
// Every deadline input remains explicit: the caller provides the parent context,
// the observation time, and the minimum budget needed to begin work. Request
// intentionally carries no tenant, priority, operation class, fallback timeout,
// retry policy, queueing policy, metric label, or tracing state.
type Request struct {
	// Context is the parent context inspected at the admission boundary.
	Context context.Context

	// Now is the caller-supplied observation time used for deterministic deadline
	// math.
	Now time.Time

	// Min is the minimum remaining budget required to start the guarded work.
	Min time.Duration
}
