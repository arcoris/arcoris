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
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

const (
	// defaultInterval is the default fixed probe cadence.
	//
	// The value is intentionally conservative: frequent enough to keep cached
	// readiness and liveness snapshots useful, but not so aggressive that a default
	// Runner should create dependency pressure in small deployments.
	defaultInterval = 5 * time.Second

	// defaultStaleAfter is the default cache freshness window.
	//
	// The default is three default probe intervals. Callers that change interval
	// should also configure staleAfter explicitly when they need a different
	// freshness relationship.
	defaultStaleAfter = 15 * time.Second
)

// config contains normalized Runner construction settings.
//
// The config is package-local. Public callers configure Runner through Option
// constructors, while NewRunner receives a complete normalized configuration.
type config struct {
	clock        clock.Clock
	interval     time.Duration
	staleAfter   time.Duration
	targets      []health.Target
	initialProbe bool
}

// defaultConfig returns the conservative Runner configuration.
//
// Targets intentionally have no default. Component owners must explicitly choose
// which health targets should be probed and cached.
func defaultConfig() config {
	return config{
		clock:        clock.RealClock{},
		interval:     defaultInterval,
		staleAfter:   defaultStaleAfter,
		initialProbe: true,
	}
}

// validate reports whether config is complete after options have been applied.
func (config config) validate() error {
	if nilClock(config.clock) {
		return ErrNilClock
	}
	if err := validateInterval(config.interval); err != nil {
		return err
	}
	if err := validateStaleAfter(config.staleAfter); err != nil {
		return err
	}
	if _, err := normalizeTargets(config.targets); err != nil {
		return err
	}

	return nil
}
