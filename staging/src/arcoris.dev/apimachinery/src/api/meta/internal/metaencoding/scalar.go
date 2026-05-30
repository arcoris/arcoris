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

package metaencoding

import "encoding/json"

// MarshalText validates a scalar before exposing it as text bytes.
func MarshalText(value string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}
	return []byte(value), nil
}

// MarshalJSONString validates a scalar before encoding it as one JSON string.
func MarshalJSONString(value string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// DecodeJSONString decodes one JSON string and distinguishes JSON null.
func DecodeJSONString(data []byte) (value string, isNull bool, err error) {
	var decoded *string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return "", false, err
	}
	if decoded == nil {
		return "", true, nil
	}
	return *decoded, false, nil
}
