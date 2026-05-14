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

package run

import "testing"

func TestGroupConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := newGroupConfig()

	if !cfg.cancelOnError {
		t.Fatal("cancelOnError default is false")
	}
	if cfg.errorMode != ErrorModeJoin {
		t.Fatalf("errorMode = %v, want ErrorModeJoin", cfg.errorMode)
	}
}

func TestGroupConfigRejectsNilOption(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilGroupOption, func() {
		newGroupConfig(nil)
	})
}

func TestNewGroupRejectsNilOption(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilGroupOption, func() {
		NewGroup(t.Context(), nil)
	})
}

func TestGroupConfigAppliesOptionsInOrder(t *testing.T) {
	t.Parallel()

	cfg := newGroupConfig(
		WithCancelOnError(false),
		WithCancelOnError(true),
		WithErrorMode(ErrorModeFirst),
		WithErrorMode(ErrorModeJoin),
	)

	if !cfg.cancelOnError {
		t.Fatal("cancelOnError = false, want later true option")
	}
	if cfg.errorMode != ErrorModeJoin {
		t.Fatalf("errorMode = %v, want later ErrorModeJoin option", cfg.errorMode)
	}
}

func TestGroupConfigRejectsInvalidErrorMode(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errInvalidErrorMode, func() {
		newGroupConfig(func(cfg *groupConfig) {
			cfg.errorMode = ErrorMode(99)
		})
	})
}
