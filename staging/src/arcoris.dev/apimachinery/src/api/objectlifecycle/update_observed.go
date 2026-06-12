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
	"context"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

// UpdateObserved replaces the Observed surface of existing live state.
func (e *Executor) UpdateObserved(ctx context.Context, req UpdateObservedRequest) (Result, error) {
	if err := e.requireExecutor(OperationUpdateObserved); err != nil {
		return Result{}, err
	}
	if err := checkContext(OperationUpdateObserved, ctx); err != nil {
		return Result{}, err
	}
	if err := e.validateUpdateObservedRequest(req); err != nil {
		return Result{}, err
	}

	prepared, err := e.prepareKeyRequest(OperationUpdateObserved, req.Resource, req.Object)
	if err != nil {
		return Result{}, err
	}
	if err := validateExpectedRevision(OperationUpdateObserved, prepared.key, req.Expected); err != nil {
		return Result{}, err
	}

	observedDescriptor, ok := prepared.resolved.version.Observed()
	if !ok {
		return Result{}, errorFor(
			OperationUpdateObserved,
			ErrorReasonObservedNotDefined,
			prepared.key,
			ErrValidationFailed,
			nil,
		)
	}
	if err := e.observedValidator.ValidateSurface(req.Observed, observedDescriptor, e.resolver); err != nil {
		return Result{}, errorFor(
			OperationUpdateObserved,
			ErrorReasonInvalidObserved,
			prepared.key,
			ErrValidationFailed,
			err,
		)
	}

	live, ok, err := e.store.Get(ctx, prepared.key)
	if err != nil {
		return Result{}, mapStoreError(OperationUpdateObserved, prepared.key, err)
	}
	if !ok {
		return Result{}, errorFor(OperationUpdateObserved, ErrorReasonNotFound, prepared.key, ErrNotFound, nil)
	}

	ownership, err := stateOwnership(OperationUpdateObserved, prepared.key, live.Ownership)
	if err != nil {
		return Result{}, err
	}
	observedOwnership, err := e.observedOwnership(prepared.key, req.Owner, req.Observed, observedDescriptor)
	if err != nil {
		return Result{}, err
	}

	nextObject := object.New[value.Value, value.Value](
		live.Object.TypeMeta,
		live.Object.ObjectMeta,
		live.Object.Desired.Clone(),
	).WithObserved(req.Observed.Clone())
	nextOwnership := ownership.WithObserved(observedOwnership)
	next := objectstore.State{
		Object:    nextObject,
		Ownership: nextOwnership,
	}

	committed, err := e.store.Update(ctx, prepared.key, req.Expected, next)
	if err != nil {
		return Result{}, mapStoreError(OperationUpdateObserved, prepared.key, err)
	}

	return Result{
		Operation: OperationUpdateObserved,
		Effect:    EffectUpdated,
		State:     committed,
		Revision:  committed.Revision,
	}, nil
}

// observedOwnership extracts the complete owner field set for an Observed replacement.
func (e *Executor) observedOwnership(
	key objectstore.Key,
	owner fieldownership.Owner,
	observed value.Value,
	descriptor types.Descriptor,
) (fieldownership.State, error) {
	fields, err := valuefieldset.ExtractOwnershipFieldsAt(
		fieldpath.Root(),
		observed,
		descriptor,
		valuefieldset.Options{
			Resolver: e.resolver,
			MaxDepth: e.applyOptions.MaxDepth,
		},
	)
	if err != nil {
		return fieldownership.State{}, errorFor(OperationUpdateObserved, ErrorReasonOwnershipInitFailed, key, ErrApplyFailed, err)
	}

	entry, err := fieldownership.NewEntry(owner, fields)
	if err != nil {
		return fieldownership.State{}, errorFor(OperationUpdateObserved, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	}

	state, err := fieldownership.NewState(entry)
	if err != nil {
		return fieldownership.State{}, errorFor(OperationUpdateObserved, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	}

	return state, nil
}
