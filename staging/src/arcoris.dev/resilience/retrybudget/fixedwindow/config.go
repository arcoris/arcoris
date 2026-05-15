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
	"time"

	"arcoris.dev/chrono/clock"
)

const (
	// DefaultWindow is the default local retry-budget accounting window.
	DefaultWindow = time.Minute

	// DefaultRatio is the default retry allowance ratio.
	//
	// A ratio of 0.2 permits one retry attempt for every five original attempts,
	// before applying DefaultMinRetries.
	DefaultRatio = 0.2

	// DefaultMinRetries is the default minimum retry allowance per window.
	DefaultMinRetries uint64 = 10
)

// config contains the validated runtime policy used by Limiter.
type config struct {
	// clock supplies observation times for fixed-window rotation.
	clock clock.PassiveClock

	// window is the fixed local accounting window duration.
	window time.Duration

	// ratio is the retry allowance multiplier applied to original attempts.
	ratio float64

	// minRetries is the minimum retry allowance available in every window.
	minRetries uint64
}

// defaultConfig returns the package baseline configuration.
func defaultConfig() config {
	return config{
		clock:      clock.RealClock{},
		window:     DefaultWindow,
		ratio:      DefaultRatio,
		minRetries: DefaultMinRetries,
	}
}

// newConfig applies opts to the default configuration and validates the result.
func newConfig(opts ...Option) (config, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if err := validateConfig(cfg); err != nil {
		return config{}, err
	}
	return cfg, nil
}
