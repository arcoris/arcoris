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

import (
	"os"
	"testing"
)

func TestShutdownConfigDefaults(t *testing.T) {
	t.Parallel()

	config := newShutdownConfig()

	if len(config.shutdownSignals) == 0 {
		t.Fatal("default shutdown signals are empty")
	}
	assertSignalSlice(t, config.escalationSignals, config.shutdownSignals)
	if config.escalationBuffer != 1 {
		t.Fatalf("escalation buffer = %d, want 1", config.escalationBuffer)
	}
	if !config.escalationEnabled {
		t.Fatal("default escalation is disabled")
	}
}

func TestShutdownConfigAppliesSignalOptions(t *testing.T) {
	t.Parallel()

	config := newShutdownConfig(
		WithShutdownSignals(testSIGINT, testSIGTERM, testSIGINT),
		WithEscalationSignals(testSIGHUP, testSIGHUP),
		WithEscalationBuffer(3),
	)

	assertSignalSlice(t, config.shutdownSignals, []os.Signal{testSIGINT, testSIGTERM})
	assertSignalSlice(t, config.escalationSignals, []os.Signal{testSIGHUP})
	if config.escalationBuffer != 3 {
		t.Fatalf("escalation buffer = %d, want 3", config.escalationBuffer)
	}
}

func TestShutdownConfigDefaultsEscalationToFinalShutdownSignals(t *testing.T) {
	t.Parallel()

	config := newShutdownConfig(WithShutdownSignals(testSIGTERM))

	assertSignalSlice(t, config.escalationSignals, config.shutdownSignals)
}

func TestShutdownConfigOptionOrdering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		options          []ShutdownOption
		wantEscalation   []os.Signal
		wantEscalationOn bool
		wantBuffer       int
	}{
		{
			name: "explicit escalation survives later shutdown replacement",
			options: []ShutdownOption{
				WithEscalationSignals(testSIGHUP),
				WithShutdownSignals(testSIGTERM),
			},
			wantEscalation:   []os.Signal{testSIGHUP},
			wantEscalationOn: true,
			wantBuffer:       1,
		},
		{
			name: "later escalation re-enables disabled escalation",
			options: []ShutdownOption{
				WithNoEscalation(),
				WithEscalationSignals(testSIGQUIT),
			},
			wantEscalation:   []os.Signal{testSIGQUIT},
			wantEscalationOn: true,
			wantBuffer:       1,
		},
		{
			name: "later no escalation wins",
			options: []ShutdownOption{
				WithEscalationSignals(testSIGHUP),
				WithNoEscalation(),
			},
			wantEscalationOn: false,
			wantBuffer:       1,
		},
		{
			name: "later buffer wins",
			options: []ShutdownOption{
				WithEscalationBuffer(2),
				WithEscalationBuffer(5),
			},
			wantEscalationOn: true,
			wantBuffer:       5,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := newShutdownConfig(tt.options...)

			if config.escalationEnabled != tt.wantEscalationOn {
				t.Fatalf("escalation enabled = %v, want %v", config.escalationEnabled, tt.wantEscalationOn)
			}
			if config.escalationBuffer != tt.wantBuffer {
				t.Fatalf("escalation buffer = %d, want %d", config.escalationBuffer, tt.wantBuffer)
			}
			if tt.wantEscalation != nil {
				assertSignalSlice(t, config.escalationSignals, tt.wantEscalation)
			}
			if !tt.wantEscalationOn && config.escalationSignals != nil {
				t.Fatalf("escalation signals = %v, want nil", config.escalationSignals)
			}
		})
	}
}

func TestShutdownConfigRejectsInvalidOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		options []ShutdownOption
	}{
		{name: "empty shutdown", message: errEmptyShutdownSignals, options: []ShutdownOption{WithShutdownSignals()}},
		{name: "nil shutdown", message: errNilSignalSetSignal, options: []ShutdownOption{WithShutdownSignals(nil)}},
		{name: "empty escalation", message: errEmptyEscalationSignals, options: []ShutdownOption{WithEscalationSignals()}},
		{name: "nil escalation", message: errNilSignalSetSignal, options: []ShutdownOption{WithEscalationSignals(nil)}},
		{name: "negative buffer", message: errNegativeEscalationBuffer, options: []ShutdownOption{WithEscalationBuffer(-1)}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tt.message, func() {
				newShutdownConfig(tt.options...)
			})
		})
	}
}

func TestShutdownConfigIgnoresNilOptions(t *testing.T) {
	t.Parallel()

	config := newShutdownConfig(nil)

	if len(config.shutdownSignals) == 0 {
		t.Fatal("nil option broke default config")
	}
}

func TestShutdownConfigCollectsSubscriptionOptions(t *testing.T) {
	t.Parallel()

	option := withNotifier(&fakeNotifier{})
	config := newShutdownConfig(withShutdownSubscriptionOptions(option, nil))

	if len(config.subscribeOptions) != 2 {
		t.Fatalf("subscribe options len = %d, want 2", len(config.subscribeOptions))
	}
}
