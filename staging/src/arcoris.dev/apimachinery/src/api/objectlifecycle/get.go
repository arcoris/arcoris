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
)

// Get resolves a resource identity and reads committed live state.
func (e *Executor) Get(ctx context.Context, req GetRequest) (Result, error) {
	if err := e.requireExecutor(OperationGet); err != nil {
		return Result{}, err
	}
	if err := checkContext(OperationGet, ctx); err != nil {
		return Result{}, err
	}
	if err := e.validateGetRequest(req); err != nil {
		return Result{}, err
	}

	prepared, err := e.prepareKeyRequest(OperationGet, req.Resource, req.Object)
	if err != nil {
		return Result{}, err
	}

	state, ok, err := e.store.Get(ctx, prepared.key)
	if err != nil {
		return Result{}, mapStoreError(OperationGet, prepared.key, err)
	}
	if !ok {
		return Result{}, errorFor(OperationGet, ErrorReasonNotFound, prepared.key, ErrNotFound, nil)
	}

	return Result{
		Operation: OperationGet,
		Effect:    EffectFound,
		State:     state,
		Revision:  state.Revision,
	}, nil
}
