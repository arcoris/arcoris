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

package objectlifecycle

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestNewExecutorRejectsNilOption(t *testing.T) {
	_, err := NewExecutor(nil)
	requireLifecycleError(t, err, ErrInvalidExecutor, ErrorReasonInvalidExecutor)
	requireErrorIs(t, err, ErrNilOption)
}

func TestNewExecutorRejectsNilStore(t *testing.T) {
	_, err := NewExecutor(
		WithResourceResolver(testCatalog(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireLifecycleError(t, err, ErrInvalidExecutor, ErrorReasonInvalidExecutor)
	requireErrorIs(t, err, ErrNilStore)
}

func TestNewExecutorRejectsNilResourceResolver(t *testing.T) {
	_, err := NewExecutor(
		WithStore(testStore(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireLifecycleError(t, err, ErrInvalidExecutor, ErrorReasonInvalidExecutor)
	requireErrorIs(t, err, ErrNilResourceResolver)
}

func TestNewExecutorRejectsNilDesiredValidator(t *testing.T) {
	_, err := NewExecutor(
		WithStore(testStore(t)),
		WithResourceResolver(testCatalog(t)),
	)
	requireLifecycleError(t, err, ErrInvalidExecutor, ErrorReasonInvalidExecutor)
	requireErrorIs(t, err, ErrNilDesiredValidator)
}

func TestNewExecutorAcceptsValidDependencies(t *testing.T) {
	executor := testExecutor(t)
	if executor == nil {
		t.Fatalf("executor is nil")
	}
}

func TestWithApplyOptionsResolverDoesNotOverrideTypeResolver(t *testing.T) {
	lifecycleResolver := &staticTypeResolver{}
	applyResolver := &staticTypeResolver{}
	executor, err := NewExecutor(
		WithStore(testStore(t)),
		WithResourceResolver(testCatalog(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
		WithTypeResolver(lifecycleResolver),
		WithApplyOptions(objectapply.Options{
			Resolver: applyResolver,
			MaxDepth: 7,
			Force:    true,
		}),
	)
	requireNoError(t, err)

	options := executor.optionsForApply(false)
	if options.Resolver != lifecycleResolver {
		t.Fatalf("Resolver = %T; want lifecycle resolver", options.Resolver)
	}
	if options.MaxDepth != 7 {
		t.Fatalf("MaxDepth = %d; want 7", options.MaxDepth)
	}
	if options.Force {
		t.Fatalf("Force = true; want request override false")
	}

	options = executor.optionsForApply(true)
	if !options.Force {
		t.Fatalf("Force = false; want request override true")
	}
}

type staticTypeResolver struct{}

func (r *staticTypeResolver) Resolve(types.TypeName) (types.Definition, bool) {
	return types.Definition{}, false
}
