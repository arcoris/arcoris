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

// DecodeObjectOwnershipFrom decodes one object ownership document from r.
//
// Decode returns the valid raw document rather than normalized state. Semantic
// normalization remains owned by api/objectownership.
func (c Codec) DecodeObjectOwnershipFrom(
	r io.Reader,
) (objectownership.Document, error) {
	return decodeTargetFrom(r, c.decode, func(
		path jsonPath,
		node jsonNode,
		config resolvedDecodeConfig,
	) (objectownership.Document, error) {
		return nodeToOwnershipDocument(path, node, config)
	})
}

// EncodeObjectOwnershipTo writes one object ownership document as JSON.
//
// Deterministic output normalizes the document before writing; default output
// preserves caller-provided entry and field order.
func (c Codec) EncodeObjectOwnershipTo(
	w io.Writer,
	doc objectownership.Document,
) error {
	node, err := ownershipDocumentToNode(rootPath(), doc, c.encode)
	if err != nil {
		return err
	}

	return encodeTargetTo(w, node, c.encode)
}
