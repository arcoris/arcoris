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
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Create validates and commits a new live object state.
func (e *Executor) Create(ctx context.Context, req CreateRequest) (Result, error) {
	if err := e.requireExecutor(OperationCreate); err != nil {
		return Result{}, err
	}
	if err := checkContext(OperationCreate, ctx); err != nil {
		return Result{}, err
	}
	if err := e.validateCreateRequest(req); err != nil {
		return Result{}, err
	}

	prepared, err := e.prepareObjectRequest(OperationCreate, req.Object)
	if err != nil {
		return Result{}, err
	}

	ownership, err := e.initialOwnership(OperationCreate, prepared.key, req.Owner, req.Object.Desired, prepared.resolved)
	if err != nil {
		return Result{}, err
	}

	committed, err := e.store.Create(ctx, prepared.key, inputState(req.Object, ownership))
	if err != nil {
		return Result{}, mapStoreError(OperationCreate, prepared.key, err)
	}

	return Result{
		Operation: OperationCreate,
		Effect:    EffectCreated,
		State:     committed,
		Revision:  committed.Revision,
	}, nil
}

// inputState builds zero-revision objectstore input from lifecycle state.
func inputState(obj objectapply.ValueObject, ownership objectownership.State) objectstore.State {
	return objectstore.State{
		Object:    obj,
		Ownership: ownership,
	}
}
