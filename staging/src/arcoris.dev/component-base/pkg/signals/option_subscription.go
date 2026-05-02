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

const (
	// errNonPositiveSubscriptionBuffer is the stable diagnostic text used when
	// WithSubscriptionBuffer receives a non-positive channel buffer size.
	//
	// Process signal delivery should not rely on an unbuffered channel. A positive
	// buffer gives the signal package room to deliver at least one notification
	// while the owner is not currently receiving.
	errNonPositiveSubscriptionBuffer = "signals: non-positive subscription buffer"
)

// SubscriptionOption configures a signal Subscription during construction.
//
// Options are applied to an internal subscribeConfig before the Subscription is
// registered. They do not mutate an already constructed Subscription and are not
// retained after construction except through normalized configuration values.
type SubscriptionOption func(*subscribeConfig)

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

// newSubscribeConfig applies options to a fresh default subscribeConfig.
//
// Nil options are ignored to keep conditional option lists easy to compose.
func newSubscribeConfig(options ...SubscriptionOption) subscribeConfig {
	config := defaultSubscribeConfig()

	for _, option := range options {
		if option == nil {
			continue
		}
		option(&config)
	}

	if config.notifier == nil {
		config.notifier = osNotifier{}
	}
	requirePositiveBuffer(config.buffer, errNonPositiveSubscriptionBuffer)

	return config
}

// WithSubscriptionBuffer configures the signal delivery channel buffer size.
//
// The buffer must be positive. A value of one is usually enough for shutdown
// coordination where only the first signal needs to be observed promptly. Larger
// values are useful only when an owner intentionally wants to retain a small
// burst of signal notifications before it can receive from C or Wait.
func WithSubscriptionBuffer(size int) SubscriptionOption {
	return func(config *subscribeConfig) {
		requirePositiveBuffer(size, errNonPositiveSubscriptionBuffer)
		config.buffer = size
	}
}

// withNotifier configures the package-local notifier seam.
//
// The option is intentionally unexported. Production callers should not replace
// os/signal registration. Tests use this hook to avoid real process signals.
func withNotifier(n notifier) SubscriptionOption {
	return func(config *subscribeConfig) {
		config.notifier = n
	}
}
