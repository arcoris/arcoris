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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

// initialOwnership builds Desired ownership for a created object.
//
// The helper asks valuefieldset for descriptor-aware semantic fields instead of
// traversing values locally. That keeps create ownership aligned with
// valueapply's field-set semantics.
func (e *Executor) initialOwnership(
	op Operation,
	key objectstore.Key,
	owner fieldownership.Owner,
	desired value.Value,
	resolved resolvedResource,
) (objectownership.State, error) {
	fields, err := valuefieldset.ExtractOwnershipFieldsAt(
		fieldpath.Root(),
		desired,
		resolved.version.Desired(),
		valuefieldset.Options{
			Resolver: e.resolver,
			MaxDepth: e.applyOptions.MaxDepth,
		},
	)
	if err != nil {
		return objectownership.State{}, errorFor(op, ErrorReasonOwnershipInitFailed, key, ErrApplyFailed, err)
	}

	entry, err := fieldownership.NewEntry(owner, fields)
	if err != nil {
		return objectownership.State{}, errorFor(op, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	}

	state, err := fieldownership.NewState(entry)
	if err != nil {
		return objectownership.State{}, errorFor(op, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	}

	return objectownership.NewState(state), nil
}

// stateOwnership converts committed ownership documents back to apply state.
func stateOwnership(op Operation, key objectstore.Key, doc objectownership.Document) (objectownership.State, error) {
	state, err := objectownership.StateFromDocument(doc)
	if err != nil {
		return objectownership.State{}, errorFor(op, ErrorReasonStoreFailed, key, ErrStoreFailed, err)
	}

	return state, nil
}
