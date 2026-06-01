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

package stamp

import (
	"time"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// MarshalText validates and encodes the timestamp as RFC3339Nano text.
func (t Timestamp) MarshalText() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	if t.IsZero() {
		return []byte(""), nil
	}

	return []byte(t.Time.UTC().Round(0).Format(time.RFC3339Nano)), nil
}

// UnmarshalText decodes and validates RFC3339Nano timestamp text.
func (t *Timestamp) UnmarshalText(data []byte) error {
	if t == nil {
		return nilReceiver("timestamp")
	}

	if len(data) == 0 {
		*t = Timestamp{}
		return nil
	}

	parsed, err := time.Parse(time.RFC3339Nano, string(data))
	if err != nil {
		return &Error{
			Record: diagnostic.WrapRecord(
				"timestamp",
				ErrInvalidTimestamp,
				ErrorReasonInvalidForm,
				"expected RFC3339Nano timestamp",
				err,
			),
		}
	}

	*t = NewTimestamp(parsed)
	return nil
}

// MarshalJSON validates and encodes the timestamp as one JSON string.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	text, err := t.MarshalText()
	if err != nil {
		return nil, err
	}

	return marshalJSONString(string(text), t.Validate)
}

// UnmarshalJSON decodes and validates a JSON string timestamp.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nilReceiver("timestamp")
	}

	value, err := unmarshalJSONString("timestamp", data)
	if err != nil {
		return err
	}

	return t.UnmarshalText([]byte(value))
}
