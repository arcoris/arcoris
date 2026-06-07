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

package healthgate

import "arcoris.dev/health"

// Set replaces the latest result stored in g.
//
// result must either have an empty Name or the same Name as the gate. Empty
// names are normalized to the gate name. Structurally invalid results are
// rejected and leave the gate unchanged.
func (g *Gate) Set(result health.Result) error {
	if g == nil {
		return health.ErrNilChecker
	}

	normalized, err := normalizeGateResult(g.name, result)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.result = normalized

	return nil
}

// Unknown stores a StatusUnknown result in g.
func (g *Gate) Unknown(reason health.Reason, message string) error {
	if g == nil {
		return health.ErrNilChecker
	}

	return g.Set(health.Unknown(g.name, reason, message))
}

// Starting stores a StatusStarting result in g.
func (g *Gate) Starting(reason health.Reason, message string) error {
	if g == nil {
		return health.ErrNilChecker
	}

	return g.Set(health.Starting(g.name, reason, message))
}

// Healthy stores a StatusHealthy result in g.
func (g *Gate) Healthy() error {
	if g == nil {
		return health.ErrNilChecker
	}

	return g.Set(health.Healthy(g.name))
}

// Degraded stores a StatusDegraded result in g.
func (g *Gate) Degraded(reason health.Reason, message string) error {
	if g == nil {
		return health.ErrNilChecker
	}

	return g.Set(health.Degraded(g.name, reason, message))
}

// Unhealthy stores a StatusUnhealthy result in g.
func (g *Gate) Unhealthy(reason health.Reason, message string) error {
	if g == nil {
		return health.ErrNilChecker
	}

	return g.Set(health.Unhealthy(g.name, reason, message))
}
