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

package healthprobe

import (
	"sync/atomic"
	"time"

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/health"
)

// Runner periodically evaluates configured health targets and stores latest
// snapshots.
//
// Runner is the runtime object of package healthprobe. It owns a ticker loop
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

// NewRunner returns a Runner that periodically evaluates targets with evaluator.
//
// evaluator MUST be non-nil. Targets MUST be configured explicitly with
// WithTargets. NewRunner does not start background work; callers start probing by
// calling Run with an owner-controlled context.
func NewRunner(evaluator *health.Evaluator, options ...Option) (*Runner, error) {
	if evaluator == nil {
		return nil, ErrNilEvaluator
	}

	config := defaultConfig()
	if err := applyOptions(&config, options...); err != nil {
		return nil, err
	}
	if err := config.validate(); err != nil {
		return nil, err
	}

	targets := copyTargets(config.targets)

	return &Runner{
		evaluator:    evaluator,
		store:        newStore(targets),
		clock:        config.clock,
		targets:      targets,
		interval:     config.interval,
		staleAfter:   config.staleAfter,
		initialProbe: config.initialProbe,
	}, nil
}

// Snapshot returns the latest observed snapshot for target.
//
// The returned snapshot is detached from Runner internals. The boolean is false
// when target is not configured or no observation has been stored for target yet.
// Stale is computed at the read boundary using Runner's configured clock.
func (r *Runner) Snapshot(target health.Target) (Snapshot, bool) {
	if !target.IsConcrete() || !containsTarget(r.targets, target) {
		return Snapshot{}, false
	}

	snapshot, ok := r.store.snapshot(target)
	if !ok {
		return Snapshot{}, false
	}

	snapshot.Stale = isStale(r.clock.Since(snapshot.Updated), r.staleAfter)
	return snapshot, true
}

// Snapshots returns all observed snapshots in configured target order.
//
// Unobserved targets are omitted. Each returned snapshot is detached from Runner
// internals. Stale is computed at the read boundary for each snapshot.
func (r *Runner) Snapshots() []Snapshot {
	snapshots := r.store.snapshots()
	for i := range snapshots {
		snapshots[i].Stale = isStale(r.clock.Since(snapshots[i].Updated), r.staleAfter)
	}

	return snapshots
}
