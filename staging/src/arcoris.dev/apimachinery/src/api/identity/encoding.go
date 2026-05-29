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

package identity

import "encoding/json"

// marshalText validates and returns an identity's canonical text bytes.
func marshalText(value string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}

	return []byte(value), nil
}

// marshalJSONString validates and encodes an identity as one JSON string.
func marshalJSONString(value string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}

	return json.Marshal(value)
}

// unmarshalJSONString decodes one JSON string and rejects all non-string JSON.
func unmarshalJSONString(name string, data []byte) (string, error) {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", invalidJSON(
			name,
			string(data),
			ErrorReasonInvalidJSON,
			"expected JSON string",
			err,
		)
	}

	if value == nil {
		return "", invalidJSON(
			name,
			"null",
			ErrorReasonInvalidJSON,
			"expected JSON string, got null",
			nil,
		)
	}

	return *value, nil
}
