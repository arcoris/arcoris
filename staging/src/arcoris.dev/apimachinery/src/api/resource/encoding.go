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

package resource

import "encoding/json"

// marshalJSONString encodes descriptor scalar text as JSON.
//
// Resource currently exposes only Scope scalar encoding. The helper keeps JSON
// string emission in one place so future scalar resource descriptors can reuse
// the same "validate first, then encode a string" shape.
func marshalJSONString(value string) ([]byte, error) {
	return json.Marshal(value)
}

// unmarshalJSONString decodes a non-null JSON string for descriptor scalars.
//
// JSON null and non-string values are rejected because resource descriptor
// scalars are encoded as canonical strings, not nullable or object-shaped
// payloads.
func unmarshalJSONString(path string, data []byte) (string, error) {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", invalidJSON(path, detailJSONMustBeString, err)
	}
	if value == nil {
		return "", invalidJSON(path, detailJSONMustBeNonNullString, nil)
	}
	return *value, nil
}
