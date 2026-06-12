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

	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Apply applies Desired intent to a live object, creating missing state.
func (e *Executor) Apply(ctx context.Context, req ApplyRequest) (ApplyResult, error) {
	if err := e.requireExecutor(OperationApply); err != nil {
		return ApplyResult{}, err
	}
	if err := checkContext(OperationApply, ctx); err != nil {
		return ApplyResult{}, err
	}
	if err := e.validateApplyRequest(req); err != nil {
		return ApplyResult{}, err
	}

	prepared, err := e.prepareObjectRequest(OperationApply, req.Object)
	if err != nil {
		return ApplyResult{}, err
	}

	live, ok, err := e.store.Get(ctx, prepared.key)
	if err != nil {
		return ApplyResult{}, mapStoreError(OperationApply, prepared.key, err)
	}
	if !ok {
		return e.applyCreate(ctx, req, prepared.resolved, prepared.key)
	}

	return e.applyExisting(ctx, req, prepared.resolved, prepared.key, live)
}

// applyCreate handles apply-to-missing by using the same ownership init as Create.
func (e *Executor) applyCreate(
	ctx context.Context,
	req ApplyRequest,
	resolved resolvedResource,
	key objectstore.Key,
) (ApplyResult, error) {
	ownership, err := e.initialOwnership(OperationApply, key, req.Owner, req.Object.Desired, resolved)
	if err != nil {
		return ApplyResult{}, err
	}

	committed, err := e.store.Create(ctx, key, inputState(req.Object, ownership))
	if err != nil {
		return ApplyResult{}, mapStoreError(OperationApply, key, err)
	}

	return ApplyResult{
		Result: Result{
			Operation: OperationApply,
			Effect:    EffectCreated,
			State:     committed,
			Revision:  committed.Revision,
		},
	}, nil
}

// applyExisting delegates Desired merge semantics to objectapply and commits output.
func (e *Executor) applyExisting(
	ctx context.Context,
	req ApplyRequest,
	resolved resolvedResource,
	key objectstore.Key,
	live objectstore.State,
) (ApplyResult, error) {
	ownership, err := stateOwnership(OperationApply, key, live.Ownership)
	if err != nil {
		return ApplyResult{}, err
	}

	applied, err := objectapply.Apply(
		objectapply.Request{
			Owner:     req.Owner,
			Live:      live.Object,
			Applied:   req.Object,
			Resource:  resolved.definition,
			Ownership: ownership,
		},
		e.optionsForApply(req.Force),
	)
	if err != nil {
		return ApplyResult{}, mapApplyError(OperationApply, key, err)
	}

	// Existing-object Apply always commits objectapply output in this first
	// lifecycle slice. No-op suppression needs a reliable equality decision for
	// both the object envelope and ownership state; value equality alone would
	// miss ownership-only changes such as another owner applying the same value.
	next := objectstore.State{
		Object:    applied.Object,
		Ownership: applied.Ownership,
	}
	committed, err := e.store.Update(ctx, key, live.Revision, next)
	if err != nil {
		return ApplyResult{}, mapStoreError(OperationApply, key, err)
	}

	return ApplyResult{
		Result: Result{
			Operation: OperationApply,
			Effect:    EffectUpdated,
			State:     committed,
			Revision:  committed.Revision,
		},
		Apply: applied,
	}, nil
}

// optionsForApply combines executor traversal options with request Force.
func (e *Executor) optionsForApply(force bool) objectapply.Options {
	opts := e.applyOptions
	opts.Resolver = e.resolver
	opts.Force = force

	return opts
}
