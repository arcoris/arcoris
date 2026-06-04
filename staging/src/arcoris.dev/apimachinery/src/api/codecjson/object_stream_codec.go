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

	"arcoris.dev/apimachinery/api/codec"
)

// DecodeObjectFrom decodes one value-backed object envelope from r.
//
// The stream decoder checks only JSON envelope shape and lexical identity
// fields. It deliberately does not perform object validation or resource lookup.
func (c Codec) DecodeObjectFrom(r io.Reader) (codec.Object, error) {
	return decodeTargetFrom(r, c.decode, nodeToObject)
}

// EncodeObjectTo writes one value-backed object envelope as JSON.
//
// The encoder preserves the object envelope policy: stable envelope field order,
// live payload values as supplied, and no metadata mutation.
func (c Codec) EncodeObjectTo(w io.Writer, obj codec.Object) error {
	node, err := objectToNode(rootPath(), obj, c.decode, c.encode)
	if err != nil {
		return err
	}

	return encodeTargetTo(w, node, c.encode)
}
