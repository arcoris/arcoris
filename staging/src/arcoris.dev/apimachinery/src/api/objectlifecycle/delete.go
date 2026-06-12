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

import "context"

// Delete removes live state after resolving the requested resource identity.
//
// Delete requires the caller's expected store revision and does not implement
// finalizers, deletion timestamps, or graceful deletion. Store-level tombstone
// revisions remain part of the objectstore result; this lifecycle result
// exposes the deleted live state.
func (e *Executor) Delete(ctx context.Context, req DeleteRequest) (Result, error) {
	if err := e.requireExecutor(OperationDelete); err != nil {
		return Result{}, err
	}
	if err := checkContext(OperationDelete, ctx); err != nil {
		return Result{}, err
	}
	if err := e.validateDeleteRequest(req); err != nil {
		return Result{}, err
	}

	prepared, err := e.prepareKeyRequest(OperationDelete, req.Resource, req.Object)
	if err != nil {
		return Result{}, err
	}

	if err := validateExpectedRevision(OperationDelete, prepared.key, req.Expected); err != nil {
		return Result{}, err
	}

	deleted, err := e.store.Delete(ctx, prepared.key, req.Expected)
	if err != nil {
		return Result{}, mapStoreError(OperationDelete, prepared.key, err)
	}

	return Result{
		Operation: OperationDelete,
		Effect:    EffectDeleted,
		State:     deleted.Deleted,
		Revision:  deleted.Revision,
	}, nil
}
