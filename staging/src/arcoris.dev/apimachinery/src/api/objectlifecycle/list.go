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

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
)

// List resolves a resource identity and reads committed live collection state.
func (e *Executor) List(ctx context.Context, req ListRequest) (ListResult, error) {
	if err := e.requireExecutor(OperationList); err != nil {
		return ListResult{}, err
	}
	if err := checkContext(OperationList, ctx); err != nil {
		return ListResult{}, err
	}
	if err := e.validateListRequest(req); err != nil {
		return ListResult{}, err
	}

	resolved, err := e.resolveKeyResource(OperationList, req.Resource)
	if err != nil {
		return ListResult{}, err
	}
	if err := validateListScopeForResource(resolved, req.Scope); err != nil {
		return ListResult{}, err
	}

	storeResult, err := e.store.List(ctx, objectstore.ListRequest{
		Resource: resolved.gvr,
		Scope:    req.Scope,
	})
	if err != nil {
		return ListResult{}, mapStoreError(OperationList, objectstore.Key{}, err)
	}

	storeResult = storeResult.Clone()
	return ListResult{
		Items:    storeResult.Items,
		Revision: storeResult.Revision,
	}, nil
}

// validateListScopeForResource applies resource-scope rules above the store.
func validateListScopeForResource(resolved resolvedResource, scope objectstore.ListScope) error {
	switch resolved.definition.Scope() {
	case resource.ScopeGlobal:
		if scope.IsNamespace() {
			return errorFor(
				OperationList,
				ErrorReasonInvalidRequest,
				objectstore.Key{},
				ErrInvalidRequest,
				objectstore.ErrInvalidListRequest,
			)
		}
		return nil
	case resource.ScopeNamespaced:
		return nil
	default:
		return errorFor(
			OperationList,
			ErrorReasonInvalidResourceContract,
			objectstore.Key{},
			ErrValidationFailed,
			ErrInvalidResourceContract,
		)
	}
}
