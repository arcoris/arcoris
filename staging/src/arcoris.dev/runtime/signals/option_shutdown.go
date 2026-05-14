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

import "os"

const (
	// errEmptyShutdownSignals is the stable diagnostic text used when shutdown
	// controller construction receives an empty shutdown signal set.
	errEmptyShutdownSignals = "signals: empty shutdown signal set"

	// errEmptyEscalationSignals is the stable diagnostic text used when
	// escalation is enabled with an empty escalation signal set.
	errEmptyEscalationSignals = "signals: empty escalation signal set"

	// errNegativeEscalationBuffer is the stable diagnostic text used when
	// WithEscalationBuffer receives a negative buffer size.
	errNegativeEscalationBuffer = "signals: negative escalation buffer"
)

// ShutdownOption configures a ShutdownController during construction.
//
// Options are applied to an internal shutdownConfig before the controller starts
// its signal loop. They do not mutate an already constructed controller.
type ShutdownOption func(*shutdownConfig)

// shutdownConfig contains construction-time settings for ShutdownController.
type shutdownConfig struct {
	shutdownSignals      []os.Signal
	escalationSignals    []os.Signal
	escalationSignalsSet bool
	escalationBuffer     int
	escalationEnabled    bool
	subscribeOptions     []SubscriptionOption
}

// defaultShutdownConfig returns the default ShutdownController config.
func defaultShutdownConfig() shutdownConfig {
	return shutdownConfig{
		shutdownSignals:   ShutdownSignals(),
		escalationBuffer:  1,
		escalationEnabled: true,
	}
}

// newShutdownConfig applies options to a fresh default shutdownConfig.
//
// Nil options are ignored to keep conditional option lists easy to compose.
func newShutdownConfig(options ...ShutdownOption) shutdownConfig {
	config := defaultShutdownConfig()

	for _, option := range options {
		if option == nil {
			continue
		}
		option(&config)
	}

	config.shutdownSignals = Unique(config.shutdownSignals)
	requireNonEmptySignals(config.shutdownSignals, errEmptyShutdownSignals)

	if config.escalationEnabled {
		if !config.escalationSignalsSet {
			config.escalationSignals = Clone(config.shutdownSignals)
		}
		config.escalationSignals = Unique(config.escalationSignals)
		requireNonEmptySignals(config.escalationSignals, errEmptyEscalationSignals)
	}
	requireNonNegativeBuffer(config.escalationBuffer, errNegativeEscalationBuffer)

	return config
}

// WithShutdownSignals replaces the signal set that starts graceful shutdown.
//
// The set must be non-empty and must not contain nil signals. When escalation
// remains enabled and WithEscalationSignals is not supplied, repeated signals
// from this final shutdown set are also used as escalation signals.
func WithShutdownSignals(sigs ...os.Signal) ShutdownOption {
	return func(config *shutdownConfig) {
		copy := Unique(sigs)
		requireNonEmptySignals(copy, errEmptyShutdownSignals)
		config.shutdownSignals = copy
	}
}

// WithEscalationSignals replaces the signal set registered after shutdown
// starts and reported as escalation.
//
// The set must be non-empty and must not contain nil signals. These signals are
// not registered during NewShutdownController construction; registration is
// staged until the first shutdown signal has been recorded. Use
// WithNoEscalation to disable repeated-signal escalation entirely.
func WithEscalationSignals(sigs ...os.Signal) ShutdownOption {
	return func(config *shutdownConfig) {
		copy := Unique(sigs)
		requireNonEmptySignals(copy, errEmptyEscalationSignals)
		config.escalationSignals = copy
		config.escalationSignalsSet = true
		config.escalationEnabled = true
	}
}

// WithEscalationBuffer configures the escalation event channel buffer size.
//
// A zero buffer is valid. Escalation delivery is best-effort and non-blocking, so
// a zero buffer reports escalation only when a receiver is ready.
func WithEscalationBuffer(size int) ShutdownOption {
	return func(config *shutdownConfig) {
		requireNonNegativeBuffer(size, errNegativeEscalationBuffer)
		config.escalationBuffer = size
	}
}

// WithNoEscalation disables repeated-signal escalation delivery.
func WithNoEscalation() ShutdownOption {
	return func(config *shutdownConfig) {
		config.escalationEnabled = false
		config.escalationSignals = nil
		config.escalationSignalsSet = false
	}
}

// withShutdownSubscriptionOptions appends Subscription options used by the
// internal signal subscription.
//
// The option is intentionally unexported. Tests use it to replace os/signal with
// a fake notifier while production callers keep the standard os/signal seam.
func withShutdownSubscriptionOptions(opts ...SubscriptionOption) ShutdownOption {
	return func(config *shutdownConfig) {
		config.subscribeOptions = append(config.subscribeOptions, opts...)
	}
}
