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

// WithSubscriptionBuffer configures the signal delivery channel buffer size.
//
// The buffer must be positive. A value of one is usually enough for shutdown
// coordination where only the first signal needs to be observed promptly. Larger
// values are useful only when an owner intentionally wants to retain a small
// burst of signal notifications before it can receive from C or Wait.
func WithSubscriptionBuffer(size int) SubscriptionOption {
	return func(cfg *subscribeConfig) {
		requirePositiveBuffer(size, errNonPositiveSubscriptionBuffer)
		cfg.buffer = size
	}
}

// withNotifier configures the package-local notifier seam.
//
// The option is intentionally unexported. Production callers should not replace
// os/signal registration. Tests use this hook to avoid real process signals.
func withNotifier(n notifier) SubscriptionOption {
	return func(cfg *subscribeConfig) {
		cfg.notifier = n
	}
}
