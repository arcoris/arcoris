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

package lifecycle

import "time"

// TimeSource is the minimal clock dependency used by Controller.
//
// A richer clock may provide timers, tickers, sleeps, or fake-time advancement,
// but lifecycle only needs current time for committed Transition.At values.
type TimeSource interface {
	Now() time.Time
}

// Option configures a lifecycle Controller at construction time.
//
// Options are applied to an internal controllerConfig before Controller is
// created. They do not mutate an already constructed Controller and are not
// retained after construction.
//
// This separation keeps the public construction API stable while allowing the
// Controller implementation to evolve internally. Options should configure
// lifecycle infrastructure only: time source, transition guards, observers, and
// other controller-owned mechanics. They must not configure component business
// logic, runtime execution, retry policy, health mapping, logging backends, or
// metrics exporters directly.
type Option func(*controllerConfig)

// controllerConfig contains construction-time settings for Controller.
//
// The config is intentionally package-local. Public callers configure it through
// Option values, while Controller receives an already normalized configuration.
//
// controllerConfig must remain small. It should contain only dependencies needed
// by the lifecycle controller itself, not dependencies that belong to component
// execution, health checking, retry, logging, metrics, or scheduling layers.
type controllerConfig struct {
	// now returns the current time used for committed transition timestamps.
	//
	// The default is time.Now. Tests and deterministic controllers may replace it
	// through WithClock.
	now func() time.Time

	// guards are evaluated before a table-valid candidate transition is
	// committed.
	//
	// Guards are called in the order they are configured. A guard rejection stops
	// evaluation and prevents the transition from being committed.
	guards []TransitionGuard

	// observers are notified after a transition has been committed.
	//
	// Observers are called in the order they are configured. Observer failures are
	// not represented in the lifecycle error model because observers cannot roll
	// back an already committed transition.
	observers []Observer
}

// defaultControllerConfig returns the default Controller construction config.
//
// The default config has no guards, no observers, and uses time.Now for
// transition commit timestamps.
func defaultControllerConfig() controllerConfig {
	return controllerConfig{
		now: time.Now,
	}
}

// newControllerConfig applies options to a fresh default controllerConfig.
//
// Nil options are ignored. This makes option composition safe for callers that
// build option lists conditionally.
//
// The returned config is independent from the variadic options slice. Guards and
// observers are stored as interface values; the lifecycle package does not clone
// the concrete objects behind those interfaces.
func newControllerConfig(options ...Option) controllerConfig {
	config := defaultControllerConfig()

	for _, option := range options {
		if option == nil {
			continue
		}

		option(&config)
	}

	return config
}
