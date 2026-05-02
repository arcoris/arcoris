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

import "testing"

func TestShutdownConfigDefaults(t *testing.T) {
	config := newShutdownConfig()

	if len(config.shutdownSignals) == 0 {
		t.Fatal("default shutdown signals are empty")
	}
	if len(config.escalationSignals) == 0 {
		t.Fatal("default escalation signals are empty")
	}
	if config.escalationBuffer != 1 {
		t.Fatalf("escalation buffer = %d, want 1", config.escalationBuffer)
	}
	if !config.escalationEnabled {
		t.Fatal("default escalation is disabled")
	}
}

func TestShutdownConfigAppliesSignalOptions(t *testing.T) {
	config := newShutdownConfig(
		WithShutdownSignals(testSIGINT, testSIGTERM, testSIGINT),
		WithEscalationSignals(testSIGHUP),
		WithEscalationBuffer(3),
	)

	if len(config.shutdownSignals) != 2 {
		t.Fatalf("shutdown len = %d, want 2", len(config.shutdownSignals))
	}
	if len(config.escalationSignals) != 1 || !sameSignal(config.escalationSignals[0], testSIGHUP) {
		t.Fatalf("escalation signals = %v, want [%v]", config.escalationSignals, testSIGHUP)
	}
	if config.escalationBuffer != 3 {
		t.Fatalf("escalation buffer = %d, want 3", config.escalationBuffer)
	}
}

func TestShutdownConfigCanDisableEscalation(t *testing.T) {
	config := newShutdownConfig(WithNoEscalation())

	if config.escalationEnabled {
		t.Fatal("escalation is enabled")
	}
	if config.escalationSignals != nil {
		t.Fatal("escalation signals are not nil")
	}
}

func TestShutdownConfigRejectsInvalidOptions(t *testing.T) {
	mustPanicWith(t, errEmptyShutdownSignals, func() {
		newShutdownConfig(WithShutdownSignals())
	})
	mustPanicWith(t, errEmptyEscalationSignals, func() {
		newShutdownConfig(WithEscalationSignals())
	})
	mustPanicWith(t, errNegativeEscalationBuffer, func() {
		newShutdownConfig(WithEscalationBuffer(-1))
	})
}

func TestShutdownConfigIgnoresNilOptions(t *testing.T) {
	config := newShutdownConfig(nil)
	if len(config.shutdownSignals) == 0 {
		t.Fatal("nil option broke default config")
	}
}
