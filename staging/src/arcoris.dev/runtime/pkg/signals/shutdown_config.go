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

// newShutdownConfig applies opts to a fresh default shutdownConfig.
//
// Nil options are ignored to keep conditional option lists easy to compose.
func newShutdownConfig(opts ...ShutdownOption) shutdownConfig {
	cfg := defaultShutdownConfig()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}

	cfg.shutdownSignals = Unique(cfg.shutdownSignals)
	requireNonEmptySignals(cfg.shutdownSignals, errEmptyShutdownSignals)

	if cfg.escalationEnabled {
		if !cfg.escalationSignalsSet {
			cfg.escalationSignals = Clone(cfg.shutdownSignals)
		}
		cfg.escalationSignals = Unique(cfg.escalationSignals)
		requireNonEmptySignals(cfg.escalationSignals, errEmptyEscalationSignals)
	}
	requireNonNegativeBuffer(cfg.escalationBuffer, errNegativeEscalationBuffer)

	return cfg
}
