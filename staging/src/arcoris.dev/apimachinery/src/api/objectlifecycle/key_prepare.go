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
	"arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

// preparedKeyRequest is the resolved objectstore identity shared by Get and Delete.
type preparedKeyRequest struct {
	// resolved is the resource/version selected by the request GVR.
	resolved resolvedResource

	// key is the committed-state identity derived from the resolved resource and object name.
	key objectstore.Key
}

// prepareKeyRequest resolves the requested resource and constructs its store key.
//
// Get and Delete operate from explicit GVR plus object name rather than an
// object envelope, so this helper performs only resolver and key validation.
// It does not validate a live object and it does not inspect stored state.
func (e *Executor) prepareKeyRequest(
	op Operation,
	gvr identity.GroupVersionResource,
	name metaidentity.ObjectName,
) (preparedKeyRequest, error) {
	resolved, err := e.resolveKeyResource(op, gvr)
	if err != nil {
		return preparedKeyRequest{}, err
	}

	key, err := keyFor(op, resolved, name)
	if err != nil {
		return preparedKeyRequest{}, err
	}

	return preparedKeyRequest{resolved: resolved, key: key}, nil
}
