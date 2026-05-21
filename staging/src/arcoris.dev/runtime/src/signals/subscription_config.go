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

package signals

// subscribeConfig contains construction-time settings for Subscription.
//
// The config is package-local. Public callers configure it through
// SubscriptionOption values while tests may use package-local options to replace
// the os/signal notifier seam.
type subscribeConfig struct {
	// buffer is the capacity of the signal delivery channel.
	buffer int

	// notifier registers and unregisters process signal delivery.
	notifier notifier
}

// defaultSubscribeConfig returns the default Subscription construction config.
func defaultSubscribeConfig() subscribeConfig {
	return subscribeConfig{
		buffer:   1,
		notifier: osNotifier{},
	}
}

// newSubscribeConfig applies opts to a fresh default subscribeConfig.
//
// Nil options are ignored to keep conditional option lists easy to compose.
func newSubscribeConfig(opts ...SubscriptionOption) subscribeConfig {
	cfg := defaultSubscribeConfig()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}

	if cfg.notifier == nil {
		cfg.notifier = osNotifier{}
	}
	requirePositiveBuffer(cfg.buffer, errNonPositiveSubscriptionBuffer)

	return cfg
}
