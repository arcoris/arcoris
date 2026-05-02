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

func TestErrorModeValidity(t *testing.T) {
	t.Parallel()

	if !ErrorModeJoin.IsValid() {
		t.Fatal("ErrorModeJoin is invalid")
	}
	if !ErrorModeFirst.IsValid() {
		t.Fatal("ErrorModeFirst is invalid")
	}
	if ErrorMode(99).IsValid() {
		t.Fatal("unknown ErrorMode is valid")
	}
}

func TestGroupConfigDefaults(t *testing.T) {
	t.Parallel()

	config := newGroupConfig()

	if !config.cancelOnError {
		t.Fatal("cancelOnError default is false")
	}
	if config.errorMode != ErrorModeJoin {
		t.Fatalf("errorMode = %v, want ErrorModeJoin", config.errorMode)
	}
}

func TestGroupConfigAppliesOptionsAndIgnoresNil(t *testing.T) {
	t.Parallel()

	config := newGroupConfig(nil, WithCancelOnError(false), WithErrorMode(ErrorModeFirst))

	if config.cancelOnError {
		t.Fatal("cancelOnError = true, want false")
	}
	if config.errorMode != ErrorModeFirst {
		t.Fatalf("errorMode = %v, want ErrorModeFirst", config.errorMode)
	}
}

func TestWithErrorModeRejectsInvalidMode(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errInvalidErrorMode, func() {
		WithErrorMode(ErrorMode(99))
	})
}
