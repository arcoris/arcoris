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

package annotations

import (
	"encoding/json"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// MarshalJSON validates and encodes the annotation set as a JSON object.
func (s Set) MarshalJSON() ([]byte, error) {
	if err := s.ValidateLexical(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Strings())
}

// UnmarshalJSON decodes a JSON object and validates all annotation keys and values.
func (s *Set) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nilReceiver("annotations")
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return &Error{
			Record: diagnostic.WrapRecord(
				"annotations",
				ErrInvalidJSON,
				ErrorReasonInvalidJSON,
				"expected JSON object with string annotation values",
				err,
			),
		}
	}

	parsed, err := FromStrings(raw)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}
