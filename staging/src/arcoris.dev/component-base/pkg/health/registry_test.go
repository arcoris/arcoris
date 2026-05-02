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

package health

import (
	"context"
	"errors"
	"testing"
)

func TestRegistryZeroValueRegisterListGet(t *testing.T) {
	t.Parallel()

	var registry Registry
	first := mustCheck(t, "first", Healthy("first"))
	second := mustCheck(t, "second", Degraded("second", ReasonOverloaded, "overloaded"))

	if err := registry.Register(TargetReady, first, second); err != nil {
		t.Fatalf("Register() = %v, want nil", err)
	}
	if got := registry.Len(TargetReady); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	if !registry.Has(TargetReady, "first") || registry.Has(TargetReady, "missing") {
		t.Fatal("Has() mismatch")
	}
	if registry.Empty() {
		t.Fatal("Empty() = true, want false")
	}

	checks := registry.Checks(TargetReady)
	if len(checks) != 2 || checks[0].Name() != "first" || checks[1].Name() != "second" {
		t.Fatalf("Checks() order mismatch")
	}
	checks[0] = second
	if registry.Checks(TargetReady)[0].Name() != "first" {
		t.Fatal("Checks() did not return defensive copy")
	}

	targets := registry.Targets()
	if len(targets) != 1 || targets[0] != TargetReady {
		t.Fatalf("Targets() = %v, want ready", targets)
	}
}

func TestRegistryEmptyAndInvalidReads(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if !registry.Empty() {
		t.Fatal("new registry should be empty")
	}
	if checks := registry.Checks(TargetUnknown); checks != nil {
		t.Fatalf("Checks(invalid target) = %v, want nil", checks)
	}
	if checks := registry.Checks(TargetReady); checks != nil {
		t.Fatalf("Checks(empty target) = %v, want nil", checks)
	}
	if registry.Has(TargetUnknown, "storage") {
		t.Fatal("Has(invalid target) = true, want false")
	}
	if registry.Has(TargetReady, "bad-name") {
		t.Fatal("Has(invalid name) = true, want false")
	}
	if registry.Has(TargetReady, "storage") {
		t.Fatal("Has(empty names) = true, want false")
	}
	if got := registry.Len(TargetUnknown); got != 0 {
		t.Fatalf("Len(invalid target) = %d, want 0", got)
	}
}

func TestRegistryRejectsInvalidInputAtomically(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	valid := mustCheck(t, "valid", Healthy("valid"))

	if err := registry.Register(TargetUnknown, valid); !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("Register(invalid target) = %v, want ErrInvalidTarget", err)
	}
	if err := registry.Register(TargetReady, nil); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("Register(nil checker) = %v, want ErrNilChecker", err)
	}

	var typedNil *typedNilChecker
	if err := registry.Register(TargetReady, typedNil); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("Register(typed nil checker) = %v, want ErrNilChecker", err)
	}

	invalid := checkerFunc{name: "bad-name", fn: func(context.Context) Result { return Healthy("bad-name") }}
	if err := registry.Register(TargetReady, valid, invalid); !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("Register(invalid batch) = %v, want ErrInvalidCheckName", err)
	}
	if got := registry.Len(TargetReady); got != 0 {
		t.Fatalf("Len() = %d, want 0 after failed batch", got)
	}
}

func TestRegistryRejectsDuplicates(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	first := mustCheck(t, "storage", Healthy("storage"))
	duplicate := mustCheck(t, "storage", Healthy("storage"))

	if err := registry.Register(TargetReady, first, duplicate); !errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("Register(batch duplicate) = %v, want ErrDuplicateCheck", err)
	}
	if got := registry.Len(TargetReady); got != 0 {
		t.Fatalf("Len() = %d, want 0 after duplicate batch", got)
	}

	if err := registry.Register(TargetReady, first); err != nil {
		t.Fatalf("Register(first) = %v, want nil", err)
	}
	if err := registry.Register(TargetReady, duplicate); !errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("Register(existing duplicate) = %v, want ErrDuplicateCheck", err)
	}
	if err := registry.Register(TargetLive, duplicate); err != nil {
		t.Fatalf("Register(same name different target) = %v, want nil", err)
	}
}

func TestRegistryErrors(t *testing.T) {
	t.Parallel()

	invalid := InvalidTargetError{Target: TargetUnknown}
	duplicate := DuplicateCheckError{Target: TargetReady, Name: "storage"}

	mustErrorIs(t, invalid, ErrInvalidTarget)
	mustErrorIs(t, duplicate, ErrDuplicateCheck)
	if invalid.Error() == "" || duplicate.Error() == "" {
		t.Fatal("registry error messages must be non-empty")
	}
}
