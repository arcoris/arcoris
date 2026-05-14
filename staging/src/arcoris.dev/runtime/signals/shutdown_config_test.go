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

	cfg := newShutdownConfig()

	if len(cfg.shutdownSignals) == 0 {
		t.Fatal("default shutdown signals are empty")
	}
	assertSignalSlice(t, cfg.escalationSignals, cfg.shutdownSignals)
	if cfg.escalationBuffer != 1 {
		t.Fatalf("escalation buffer = %d, want 1", cfg.escalationBuffer)
	}
	if !cfg.escalationEnabled {
		t.Fatal("default escalation is disabled")
	}
}

func TestShutdownConfigAppliesSignalOptions(t *testing.T) {
	t.Parallel()

	cfg := newShutdownConfig(
		WithShutdownSignals(testSIGINT, testSIGTERM, testSIGINT),
		WithEscalationSignals(testSIGHUP, testSIGHUP),
		WithEscalationBuffer(3),
	)

	assertSignalSlice(t, cfg.shutdownSignals, []os.Signal{testSIGINT, testSIGTERM})
	assertSignalSlice(t, cfg.escalationSignals, []os.Signal{testSIGHUP})
	if cfg.escalationBuffer != 3 {
		t.Fatalf("escalation buffer = %d, want 3", cfg.escalationBuffer)
	}
}

func TestShutdownConfigDefaultsEscalationToFinalShutdownSignals(t *testing.T) {
	t.Parallel()

	cfg := newShutdownConfig(WithShutdownSignals(testSIGTERM))

	assertSignalSlice(t, cfg.escalationSignals, cfg.shutdownSignals)
}

func TestShutdownConfigOptionOrdering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		opts             []ShutdownOption
		wantEscalation   []os.Signal
		wantEscalationOn bool
		wantBuffer       int
	}{
		{
			name: "explicit escalation survives later shutdown replacement",
			opts: []ShutdownOption{
				WithEscalationSignals(testSIGHUP),
				WithShutdownSignals(testSIGTERM),
			},
			wantEscalation:   []os.Signal{testSIGHUP},
			wantEscalationOn: true,
			wantBuffer:       1,
		},
		{
			name: "later escalation re-enables disabled escalation",
			opts: []ShutdownOption{
				WithNoEscalation(),
				WithEscalationSignals(testSIGQUIT),
			},
			wantEscalation:   []os.Signal{testSIGQUIT},
			wantEscalationOn: true,
			wantBuffer:       1,
		},
		{
			name: "later no escalation wins",
			opts: []ShutdownOption{
				WithEscalationSignals(testSIGHUP),
				WithNoEscalation(),
			},
			wantEscalationOn: false,
			wantBuffer:       1,
		},
		{
			name: "later buffer wins",
			opts: []ShutdownOption{
				WithEscalationBuffer(2),
				WithEscalationBuffer(5),
			},
			wantEscalationOn: true,
			wantBuffer:       5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := newShutdownConfig(tc.opts...)

			if cfg.escalationEnabled != tc.wantEscalationOn {
				t.Fatalf("escalation enabled = %v, want %v", cfg.escalationEnabled, tc.wantEscalationOn)
			}
			if cfg.escalationBuffer != tc.wantBuffer {
				t.Fatalf("escalation buffer = %d, want %d", cfg.escalationBuffer, tc.wantBuffer)
			}
			if tc.wantEscalation != nil {
				assertSignalSlice(t, cfg.escalationSignals, tc.wantEscalation)
			}
			if !tc.wantEscalationOn && cfg.escalationSignals != nil {
				t.Fatalf("escalation signals = %v, want nil", cfg.escalationSignals)
			}
		})
	}
}

func TestShutdownConfigRejectsInvalidOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		opts    []ShutdownOption
	}{
		{name: "empty shutdown", message: errEmptyShutdownSignals, opts: []ShutdownOption{WithShutdownSignals()}},
		{name: "nil shutdown", message: errNilSignalSetSignal, opts: []ShutdownOption{WithShutdownSignals(nil)}},
		{name: "empty escalation", message: errEmptyEscalationSignals, opts: []ShutdownOption{WithEscalationSignals()}},
		{name: "nil escalation", message: errNilSignalSetSignal, opts: []ShutdownOption{WithEscalationSignals(nil)}},
		{name: "negative buffer", message: errNegativeEscalationBuffer, opts: []ShutdownOption{WithEscalationBuffer(-1)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tc.message, func() {
				newShutdownConfig(tc.opts...)
			})
		})
	}
}

func TestShutdownConfigIgnoresNilOptions(t *testing.T) {
	t.Parallel()

	cfg := newShutdownConfig(nil)

	if len(cfg.shutdownSignals) == 0 {
		t.Fatal("nil option broke default config")
	}
}

func TestShutdownConfigCollectsSubscriptionOptions(t *testing.T) {
	t.Parallel()

	opt := withNotifier(&fakeNotifier{})
	cfg := newShutdownConfig(withShutdownSubscriptionOptions(opt, nil))

	if len(cfg.subscribeOptions) != 2 {
		t.Fatalf("subscribe options len = %d, want 2", len(cfg.subscribeOptions))
	}
}
