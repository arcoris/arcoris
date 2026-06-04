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

	"arcoris.dev/apimachinery/api/codec"
)

// DecodeObject decodes one value-backed object envelope from JSON bytes.
//
// The byte API delegates to DecodeObjectFrom so root-envelope validation,
// duplicate-key rejection, UTF-8 checks, depth limits, and trailing-data
// handling are shared with stream decoding.
func (c Codec) DecodeObject(data []byte) (codec.Object, error) {
	return c.DecodeObjectFrom(bytes.NewReader(data))
}

// EncodeObject encodes one value-backed object envelope into JSON bytes.
//
// The object envelope is encoded through the same private node writer as every
// other target. This keeps envelope field ordering and payload deterministic
// behavior independent from Go map encoding.
func (c Codec) EncodeObject(obj codec.Object) ([]byte, error) {
	node, err := objectToNode(rootPath(), obj, c.decode, c.encode)
	if err != nil {
		return nil, err
	}

	return encodeTargetBytes(node, c.encode)
}
