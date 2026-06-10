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

// NamePrefix is metadata used by higher layers to generate a concrete object name.
//
// The metadata package never expands a prefix into a concrete name. Generation
// policy belongs to higher layers that own persistence and uniqueness.
type NamePrefix string

// String returns the raw name prefix text without validating it.
func (p NamePrefix) String() string {
	return string(p)
}

// CanonicalText validates the name prefix and returns its canonical text.
func (p NamePrefix) CanonicalText() (string, error) {
	if err := p.ValidateLexical(); err != nil {
		return "", err
	}

	return p.String(), nil
}

// IsZero reports whether the prefix is absent.
func (p NamePrefix) IsZero() bool {
	return p == ""
}

// IsAbsent reports whether the name prefix is absent.
func (p NamePrefix) IsAbsent() bool {
	return p == ""
}
