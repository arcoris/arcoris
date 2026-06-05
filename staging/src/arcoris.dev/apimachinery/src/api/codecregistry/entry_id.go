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

import "fmt"

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
	value := string(id)
	if value == "" {
		return "", fmt.Errorf("%w: entry ID must be non-empty", ErrInvalidEntryID)
	}

	previousSeparator := false
	for index, r := range value {
		if r > 0x7f {
			return "", fmt.Errorf("%w: entry ID contains non-ASCII character at byte %d", ErrInvalidEntryID, index)
		}
		if r <= 0x20 || r == 0x7f {
			return "", fmt.Errorf("%w: entry ID contains whitespace or control character at byte %d", ErrInvalidEntryID, index)
		}
		if isEntryIDUppercase(r) {
			return "", fmt.Errorf("%w: entry ID must be lowercase", ErrInvalidEntryID)
		}

		separator := isEntryIDSeparator(r)
		switch {
		case isEntryIDAlphaNumeric(r):
			previousSeparator = false
		case separator:
			if index == 0 {
				return "", fmt.Errorf("%w: entry ID must not start with a separator", ErrInvalidEntryID)
			}
			if previousSeparator {
				return "", fmt.Errorf("%w: entry ID must not contain repeated separators", ErrInvalidEntryID)
			}
			previousSeparator = true
		default:
			return "", fmt.Errorf("%w: entry ID contains invalid character %q", ErrInvalidEntryID, r)
		}
	}
	if previousSeparator {
		return "", fmt.Errorf("%w: entry ID must not end with a separator", ErrInvalidEntryID)
	}

	return id, nil
}

// isEntryIDAlphaNumeric reports whether r is an allowed lowercase letter or digit.
func isEntryIDAlphaNumeric(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
}

// isEntryIDSeparator reports whether r separates EntryID path/token segments.
func isEntryIDSeparator(r rune) bool {
	switch r {
	case '.', '-', '_', '/':
		return true
	default:
		return false
	}
}

// isEntryIDUppercase reports whether r is an ASCII uppercase letter.
func isEntryIDUppercase(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
