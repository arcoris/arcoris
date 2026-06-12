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
	"bytes"

	"arcoris.dev/apimachinery/api/objectownership"
)

// DecodeObjectOwnership decodes one object ownership state from JSON bytes.
//
// The byte API delegates to DecodeObjectOwnershipFrom so shape checks and
// objectownership.Validate wrapping are identical for byte and stream callers.
func (c Codec) DecodeObjectOwnership(data []byte) (objectownership.State, error) {
	return c.DecodeObjectOwnershipFrom(bytes.NewReader(data))
}

// EncodeObjectOwnership encodes one object ownership state into JSON bytes.
//
// Encoding renders the current canonical ownership JSON shape. It does not
// compute ownership, mutate metadata, or perform lifecycle decisions.
func (c Codec) EncodeObjectOwnership(
	state objectownership.State,
) ([]byte, error) {
	node, err := ownershipStateToNode(rootPath(), state, c.encode)
	if err != nil {
		return nil, err
	}

	return encodeTargetBytes(node, c.encode)
}
