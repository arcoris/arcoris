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

	"arcoris.dev/apimachinery/api/value"
)

// DecodeValue decodes one JSON value document from bytes.
//
// The byte API is a convenience transport over DecodeValueFrom. Keeping the
// stream path authoritative prevents drift between byte-slice and io.Reader
// decoding behavior.
func (c Codec) DecodeValue(data []byte) (value.Value, error) {
	return c.DecodeValueFrom(bytes.NewReader(data))
}

// EncodeValue encodes one value document into JSON bytes.
//
// Generic value encoding is descriptor-agnostic: JSON-native value kinds are
// written directly, while value kinds that require type descriptors are rejected
// before they reach the shared node writer.
func (c Codec) EncodeValue(v value.Value) ([]byte, error) {
	node, err := valueToNode(rootPath(), v, c.encode)
	if err != nil {
		return nil, err
	}

	return encodeTargetBytes(node, c.encode)
}
