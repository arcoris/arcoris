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

package healthgrpc

import (
	"time"

	"arcoris.dev/chrono/clock"
)

const (
	// defaultWatchInterval is the polling cadence for each Watch stream.
	//
	// The default is intentionally conservative: frequent enough for a standard
	// health watch to converge quickly, but not so frequent that every stream
	// becomes a tight health-evaluation loop.
	defaultWatchInterval = 5 * time.Second

	// defaultMaxListServices bounds List response construction by default.
	//
	// Most component owners expose only a small fixed health surface. A
	// bound keeps accidental large service configurations from producing large
	// response maps without requiring every caller to set a limit explicitly.
	defaultMaxListServices = 100
)

// config contains normalized Server construction settings.
//
// The struct is private because callers should configure the adapter through
// Options. NewServer validates config once, then stores the immutable result on
// Server for request-time behavior.
type config struct {
	// services is the ordered service mapping list before it is indexed.
	services []ServiceMapping

	// watchInterval is the per-stream polling interval used by Watch.
	watchInterval time.Duration

	// clock creates Watch tickers and is injectable for deterministic tests.
	clock clock.Clock

	// maxListServices is the service-count guardrail enforced by List.
	maxListServices int
}

// defaultConfig returns the safe default adapter configuration.
//
// Defaults publish only the standard whole-server gRPC health service, use the
// real clock, and keep Watch/List behavior bounded without introducing
// background workers or caches.
func defaultConfig() config {
	return config{
		services:        []ServiceMapping{defaultServiceMapping()},
		watchInterval:   defaultWatchInterval,
		clock:           clock.RealClock{},
		maxListServices: defaultMaxListServices,
	}
}

// validate verifies config after all options have been applied.
//
// Final validation catches cross-option issues such as duplicate service names.
// Individual options still validate their own values early so configuration
// mistakes point at the option that introduced them.
func (config config) validate() error {
	if nilClock(config.clock) {
		return ErrNilClock
	}
	if err := validateWatchInterval(config.watchInterval); err != nil {
		return err
	}
	if err := validateMaxListServices(config.maxListServices); err != nil {
		return err
	}
	if _, err := normalizeServiceMappings(config.services); err != nil {
		return err
	}

	return nil
}
