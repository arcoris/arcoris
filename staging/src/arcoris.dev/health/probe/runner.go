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

package probe

import (
	"sync/atomic"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

// Runner periodically evaluates configured health targets and stores latest
// snapshots.
//
// Runner is the runtime object of package probe. It owns a ticker loop
// when Run is active, a private latest-snapshot store, and read methods for
// concurrent consumers. Runner does not own health checks, registries, evaluator
// execution policy, transport rendering, retry policy, metrics, logging,
// tracing, lifecycle transitions, restart policy, admission, routing, or
// scheduling decisions.
//
// A Runner supports at most one active Run call at a time. Snapshot and Snapshots
// may be called concurrently with Run.
type Runner struct {
	evaluator *health.Evaluator
	store     *store

	clock        clock.Clock
	targets      []health.Target
	interval     time.Duration
	staleAfter   time.Duration
	initialProbe bool

	running atomic.Bool
}
