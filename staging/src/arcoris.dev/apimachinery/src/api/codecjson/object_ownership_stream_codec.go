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

package codecjson

import (
	"io"

	"arcoris.dev/apimachinery/api/objectownership"
)

// DecodeObjectOwnershipFrom decodes one object ownership state from r.
//
// Decode returns the canonical ownership model and validates it through
// api/objectownership. The codec owns JSON syntax only.
func (c Codec) DecodeObjectOwnershipFrom(
	r io.Reader,
) (objectownership.State, error) {
	return decodeTargetFrom(r, c.decode, func(
		path jsonPath,
		node jsonNode,
		config resolvedDecodeConfig,
	) (objectownership.State, error) {
		return nodeToOwnershipState(path, node, config)
	})
}

// EncodeObjectOwnershipTo writes one object ownership state as JSON.
//
// Deterministic output normalizes the state before writing. The encoder does
// not compute ownership or interpret lifecycle semantics.
func (c Codec) EncodeObjectOwnershipTo(
	w io.Writer,
	state objectownership.State,
) error {
	node, err := ownershipStateToNode(rootPath(), state, c.encode)
	if err != nil {
		return err
	}

	return encodeTargetTo(w, node, c.encode)
}
