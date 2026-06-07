// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthregistry

import (
	"context"
	"errors"
	"sync"
	"testing"

	"arcoris.dev/health"
)

type typedNilChecker struct{}

func (*typedNilChecker) Name() string {
	return "typed_nil"
}

func (*typedNilChecker) Check(context.Context) health.Result {
	return health.Healthy("typed_nil")
}

func TestBuilderBuildsImmutableRegistry(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	mustRegister(t, builder, health.TargetReady, checker("storage"), checker("queue"))

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	mustRegister(t, builder, health.TargetReady, checker("later"))

	set, err := registry.ResolveChecks(health.TargetReady)
	if err != nil {
		t.Fatalf("ResolveChecks() = %v, want nil", err)
	}
	if set.Len() != 2 || !set.Has("storage") || !set.Has("queue") || set.Has("later") {
		t.Fatalf("immutable set = %+v, want storage and queue only", set.Checks())
	}

	targets := registry.Targets()
	targets[0] = health.TargetStartup
	if got := registry.Targets()[0]; got != health.TargetReady {
		t.Fatalf("Targets exposed mutable storage: got %s", got)
	}
}

func TestBuilderRejectsInvalidRegistrationBatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*Builder) error
		want error
	}{
		{
			name: "invalid target",
			call: func(builder *Builder) error {
				return builder.Register(health.TargetUnknown, checker("storage"))
			},
			want: health.ErrInvalidTarget,
		},
		{
			name: "nil checker",
			call: func(builder *Builder) error {
				return builder.Register(health.TargetReady, nil)
			},
			want: health.ErrNilChecker,
		},
		{
			name: "typed nil checker",
			call: func(builder *Builder) error {
				var checker *typedNilChecker
				return builder.Register(health.TargetReady, checker)
			},
			want: health.ErrNilChecker,
		},
		{
			name: "invalid check name",
			call: func(builder *Builder) error {
				return builder.Register(health.TargetReady, namedChecker{name: "bad-name"})
			},
			want: health.ErrInvalidCheckName,
		},
		{
			name: "duplicate in batch",
			call: func(builder *Builder) error {
				return builder.Register(health.TargetReady, checker("dup"), checker("dup"))
			},
			want: ErrDuplicateCheck,
		},
		{
			name: "duplicate existing",
			call: func(builder *Builder) error {
				if err := builder.Register(health.TargetReady, checker("dup")); err != nil {
					return err
				}
				return builder.Register(health.TargetReady, checker("dup"))
			},
			want: ErrDuplicateCheck,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			builder := NewBuilder()
			if err := tc.call(builder); !errors.Is(err, tc.want) {
				t.Fatalf("Register() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestBuilderRegistrationBatchIsAtomic(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	mustRegister(t, builder, health.TargetReady, checker("existing"))

	err := builder.Register(
		health.TargetReady,
		checker("existing"),
		checker("new_check"),
	)
	if !errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("Register() = %v, want ErrDuplicateCheck", err)
	}

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}
	if registry.Has(health.TargetReady, "new_check") {
		t.Fatal("failed batch registered non-conflicting check")
	}
}

func TestBuilderAllowsSameNameOnDifferentTargets(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	mustRegister(t, builder, health.TargetLive, checker("shared"))
	mustRegister(t, builder, health.TargetReady, checker("shared"))

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}
	if !registry.Has(health.TargetLive, "shared") || !registry.Has(health.TargetReady, "shared") {
		t.Fatal("same check name should be allowed on different targets")
	}
}

func TestDuplicateCheckErrorCarriesIndexes(t *testing.T) {
	t.Parallel()

	err := NewBuilder().Register(health.TargetReady, checker("dup"), checker("dup"))

	var duplicate DuplicateCheckError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Register() = %v, want DuplicateCheckError", err)
	}
	if duplicate.Name != "dup" || duplicate.Index != 1 || duplicate.PreviousIndex != 0 {
		t.Fatalf("duplicate = %+v, want dup index 1 previous 0", duplicate)
	}
}

func TestRegistrationErrorsExposeTypedDiagnostics(t *testing.T) {
	t.Parallel()

	nilErr := NewBuilder().Register(health.TargetReady, nil)
	var nilChecker NilCheckerError
	if !errors.As(nilErr, &nilChecker) || nilChecker.Target != health.TargetReady || nilChecker.Index != 0 {
		t.Fatalf("nil checker error = %v, want typed target ready index 0", nilErr)
	}

	nameErr := NewBuilder().Register(health.TargetReady, namedChecker{name: "bad-name"})
	var invalidName InvalidCheckNameError
	if !errors.As(nameErr, &invalidName) ||
		invalidName.Target != health.TargetReady ||
		invalidName.Index != 0 ||
		invalidName.Name != "bad-name" {
		t.Fatalf("invalid name error = %v, want typed target ready index 0", nameErr)
	}
}

func TestRegistryResolveChecks(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, health.TargetReady, checker("storage"), checker("queue"))

	set, err := registry.ResolveChecks(health.TargetReady)
	if err != nil {
		t.Fatalf("ResolveChecks(ready) = %v, want nil", err)
	}
	if set.Target() != health.TargetReady || set.Len() != 2 {
		t.Fatalf("ready set = %+v, want two ready checks", set)
	}

	empty, err := registry.ResolveChecks(health.TargetLive)
	if err != nil {
		t.Fatalf("ResolveChecks(live) = %v, want nil", err)
	}
	if empty.Target() != health.TargetLive || !empty.Empty() {
		t.Fatalf("live set = %+v, want empty live set", empty)
	}

	if _, err := registry.ResolveChecks(health.TargetUnknown); !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("ResolveChecks(unknown) = %v, want ErrInvalidTarget", err)
	}
}

func TestRegistryConcurrentReads(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, health.TargetReady, checker("storage"), checker("queue"))

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			set, err := registry.ResolveChecks(health.TargetReady)
			if err != nil || set.Len() != 2 {
				t.Errorf("ResolveChecks() = %v, %d; want nil, 2", err, set.Len())
			}
			_ = registry.Targets()
			_ = registry.Len(health.TargetReady)
			_ = registry.Has(health.TargetReady, "storage")
		}()
	}
	wg.Wait()
}

func checker(name string) health.Checker {
	return namedChecker{name: name}
}

func mustRegister(t *testing.T, builder *Builder, target health.Target, checks ...health.Checker) {
	t.Helper()

	if err := builder.Register(target, checks...); err != nil {
		t.Fatalf("Register() = %v, want nil", err)
	}
}

func mustRegistry(t *testing.T, target health.Target, checks ...health.Checker) *Registry {
	t.Helper()

	builder := NewBuilder()
	mustRegister(t, builder, target, checks...)

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	return registry
}

type namedChecker struct {
	name string
}

func (checker namedChecker) Name() string {
	return checker.name
}

func (checker namedChecker) Check(context.Context) health.Result {
	return health.Healthy(checker.name)
}
