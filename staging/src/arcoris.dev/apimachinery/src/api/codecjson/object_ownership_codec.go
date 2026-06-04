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

// DecodeObjectOwnership decodes one object ownership document from JSON bytes.
//
// The byte API delegates to DecodeObjectOwnershipFrom so document-shape checks
// and objectownership.Validate wrapping are identical for byte and stream
// callers.
func (c Codec) DecodeObjectOwnership(data []byte) (objectownership.Document, error) {
	return c.DecodeObjectOwnershipFrom(bytes.NewReader(data))
}

// EncodeObjectOwnership encodes one object ownership document into JSON bytes.
//
// Encoding operates on objectownership.Document, not objectownership.State. The
// JSON codec only renders the stable document representation and leaves
// semantic state conversion to api/objectownership.
func (c Codec) EncodeObjectOwnership(
	doc objectownership.Document,
) ([]byte, error) {
	node, err := ownershipDocumentToNode(rootPath(), doc, c.encode)
	if err != nil {
		return nil, err
	}

	return encodeTargetBytes(node, c.encode)
}
