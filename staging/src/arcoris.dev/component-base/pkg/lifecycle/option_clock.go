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

// WithClock configures the time source used for committed transition timestamps.
//
// The lifecycle controller only needs a current-time source for Transition.At.
// For that reason, WithClock accepts TimeSource instead of depending on the full
// clock package API.
//
// If source is nil, the option leaves the default time source unchanged.
//
// The configured time source is called when Controller commits a transition. It
// is not called while reducing transitions, validating transition rules, running
// guards, or notifying observers.
func WithClock(source TimeSource) Option {
	return func(config *controllerConfig) {
		if source == nil {
			return
		}

		config.now = source.Now
	}
}
