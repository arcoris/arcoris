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

package meta

import "arcoris.dev/apimachinery/api/meta/internal/metagrammar"

// maxPageTokenLength bounds opaque page tokens before transport layers exist.
const maxPageTokenLength = 2048

// PageToken is an opaque pagination continuation token.
//
// api/meta validates only token safety. It does not parse storage cursor
// structure or impose pagination policy.
type PageToken string

// ParsePageToken validates and returns an opaque pagination token.
func ParsePageToken(value string) (PageToken, error) {
	token := PageToken(value)
	if err := token.Validate(); err != nil {
		return "", err
	}

	return token, nil
}

// String returns the raw opaque token text.
func (t PageToken) String() string {
	return string(t)
}

// IsZero reports whether the token is absent.
func (t PageToken) IsZero() bool {
	return t == ""
}

// Validate checks only scalar safety and deliberately ignores token internals.
func (t PageToken) Validate() error {
	return fromGrammar(
		"pageToken",
		ErrInvalidPageToken,
		metagrammar.ValidateOpaqueToken(
			"page token",
			t.String(),
			metagrammar.OpaqueTokenOptions{
				AllowEmpty: true,
				MaxLength:  maxPageTokenLength,
			},
		),
	)
}

// MarshalText validates and encodes the token as scalar text.
func (t PageToken) MarshalText() ([]byte, error) {
	return marshalText(t.String(), t.Validate)
}

// UnmarshalText decodes and validates scalar token text.
func (t *PageToken) UnmarshalText(data []byte) error {
	if t == nil {
		return nilReceiver("pageToken")
	}

	value, err := ParsePageToken(string(data))
	if err != nil {
		return err
	}

	*t = value
	return nil
}

// MarshalJSON validates and encodes the token as one JSON string.
func (t PageToken) MarshalJSON() ([]byte, error) {
	return marshalJSONString(t.String(), t.Validate)
}

// UnmarshalJSON decodes one JSON string and rejects null or non-string input.
func (t *PageToken) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nilReceiver("pageToken")
	}

	value, err := unmarshalJSONString("pageToken", data)
	if err != nil {
		return err
	}

	parsed, err := ParsePageToken(value)
	if err != nil {
		return err
	}

	*t = parsed
	return nil
}
