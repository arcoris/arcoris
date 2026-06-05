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

package codecregistry

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// EntryID is the stable registry identity of one configured codec candidate.
//
// EntryID is supplied by the registry owner. It is not derived from
// codec.Format, codec.MediaType, or codec.Info because several configured codec
// instances may deliberately share the same format and media type.
type EntryID string

// NewEntryID validates value and returns its canonical EntryID.
func NewEntryID(value string) (EntryID, error) {
	return EntryID(value).Normalize()
}

// MustEntryID returns a canonical EntryID or panics when value is invalid.
func MustEntryID(value string) EntryID {
	id, err := NewEntryID(value)
	if err != nil {
		panic(err)
	}

	return id
}

// IsZero reports whether id is empty.
func (id EntryID) IsZero() bool {
	return id == ""
}

// String returns id as a string.
func (id EntryID) String() string {
	return string(id)
}

// Normalize validates id and returns its canonical form.
func (id EntryID) Normalize() (EntryID, error) {
	if violation := lexical.ValidateASCIIToken(string(id), entryIDTokenOptions()); violation != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidEntryID, violation.Detail)
	}

	return id, nil
}

// entryIDTokenOptions returns the shared lexical grammar for registry entry IDs.
func entryIDTokenOptions() lexical.TokenOptions {
	return lexical.TokenOptions{
		MinLength:                1,
		AllowLower:               true,
		AllowDigit:               true,
		AllowHyphen:              true,
		AllowDot:                 true,
		AllowUnderscore:          true,
		AllowSlash:               true,
		RequireAlnumEdges:        true,
		RejectAdjacentSeparators: true,
	}
}
