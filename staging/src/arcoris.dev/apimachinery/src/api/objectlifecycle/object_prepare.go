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
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
)

// preparedObjectRequest is the descriptor-aware input shared by Create and Apply.
type preparedObjectRequest struct {
	// resolved is the resource/version selected by the object's TypeMeta.
	resolved resolvedResource

	// key is the committed-state identity derived from the resolved resource and ObjectMeta.
	key objectstore.Key
}

// prepareObjectRequest resolves, keys, and validates a value-backed object request.
//
// This helper centralizes the descriptor-aware boundary that both Create and
// Apply must cross before they can build ownership or commit through
// objectstore. It intentionally delegates all object shape checks to
// objectvalidation rather than duplicating descriptor logic locally.
func (e *Executor) prepareObjectRequest(
	op Operation,
	obj objectapply.ValueObject,
) (preparedObjectRequest, error) {
	resolved, err := e.resolveObjectResource(op, obj)
	if err != nil {
		return preparedObjectRequest{}, err
	}

	key, err := keyFor(op, resolved, obj.ObjectName())
	if err != nil {
		return preparedObjectRequest{}, err
	}

	if err := e.validateObject(op, key, obj, resolved); err != nil {
		return preparedObjectRequest{}, err
	}

	return preparedObjectRequest{resolved: resolved, key: key}, nil
}
